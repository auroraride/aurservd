package biz

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/payment/alipay"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/utils"
)

type riderBiz struct {
	orm            *ent.RiderClient
	ctx            context.Context
	cacheKeyPrefix string
}

func NewRiderBiz() *riderBiz {
	return &riderBiz{
		orm:            ent.Database.Rider,
		ctx:            context.Background(),
		cacheKeyPrefix: "RIDER_",
	}
}

// ChangePhone 修改手机号
func (b *riderBiz) ChangePhone(r *ent.Rider, req *definition.RiderChangePhoneReq) (err error) {
	if r.Phone == req.Phone {
		return errors.New("新手机号不能与旧手机号相同")
	}
	// 验证验证码
	service.NewSms().VerifyCodeX(req.Phone, req.SmsId, req.SmsCode)

	// 判定修改的手机号是否已经存在
	if b.orm.QueryNotDeleted().Where(rider.PhoneIn(req.Phone)).ExistX(b.ctx) {
		return errors.New("手机号已存在, 请更换手机号")
	}

	// 修改手机号
	_, err = b.orm.UpdateOne(r).SetPhone(req.Phone).Save(b.ctx)
	if err != nil {
		return err
	}

	// 记录日志
	go logging.NewOperateLog().
		SetRef(r).
		SetOperate(model.OperateProfile).
		SetDiff(r.Phone, req.Phone).
		Send()
	return
}

// Signin 登录
func (b *riderBiz) Signin(device *model.Device, req *definition.RiderSignupReq) (res *model.RiderSigninRes, err error) {
	// 兼容老版本
	if req.SigninType == nil {
		req.SigninType = silk.UInt64(model.SigninTypeSms)
	}
	switch *req.SigninType {
	case model.SigninTypeSms:
		service.NewSms().VerifyCodeX(req.Phone, req.SmsId, req.SmsCode)
	case model.SigninTypeAuth:
		// 微信授权登录
		if req.AuthType == model.AuthTypeWechat {
			req.Phone = service.NewRiderMiniProgram().GetPhoneNumber(req.AuthCode)
		}

		// 支付宝授权登录
		if req.AuthType == model.AuthTypeAlipay {
			var phone string
			phone, err = alipay.NewMiniProgram().GetPhoneNumber(req.EncryptedData)
			if err != nil {
				return nil, err
			}
			req.Phone = phone
		}
	}

	ctx := context.Background()
	orm := ent.Database.Rider
	var u *ent.Rider

	u, err = orm.QueryNotDeleted().Where(rider.Phone(req.Phone)).WithPerson().WithEnterprise().WithStation().First(ctx)
	if err != nil {
		// 创建骑手
		u, err = orm.Create().
			SetPhone(req.Phone).
			SetLastDevice(device.Serial).
			SetDeviceType(device.Type.Value()).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	// 判定用户是否被封禁
	if service.NewRider().IsBanned(u) {
		return nil, errors.New("用户已被封禁")
	}

	token := xid.New().String() + utils.RandTokenString()
	key := fmt.Sprintf("%s%d", b.cacheKeyPrefix, u.ID)

	// 删除旧的token
	if old := cache.Get(ctx, key).Val(); old != "" {
		cache.Del(ctx, key)
		cache.Del(ctx, old)
	}

	// 更新设备
	if u.LastDevice != device.Serial {
		service.NewRider().SetNewDevice(u, device)
	}

	res = service.NewRider().Profile(u, device, token)

	effectiveContract := service.NewContract().QueryEffectiveContract(u)
	if effectiveContract != nil {
		encryptDocID, err := utils.EncryptAES([]byte(ar.Config.Contract.EncryptKey), effectiveContract.DocID)
		if err != nil || encryptDocID == "" {
			zap.L().Error("加密合同编号失败", zap.Error(err))
			return nil, err
		}
		res.ContractDocID = encryptDocID
	}

	// 设置登录token
	service.NewRider().ExtendTokenTime(u.ID, token)

	return res, nil
}

// GetAlipayOpenid 获取支付宝小程序openid
func (b *riderBiz) GetAlipayOpenid(req *model.OpenidReq) (res *model.OpenidRes, err error) {
	openId, err := alipay.NewMiniProgram().GetOpenid(req.Code)
	if err != nil || openId == "" {
		return nil, errors.New("获取openid失败")
	}
	return &model.OpenidRes{Openid: openId}, nil
}

// SetMobPushId 设置骑手推送ID
func (b *riderBiz) SetMobPushId(u *ent.Rider, req *definition.RiderSetMobPushReq) (err error) {
	if u.PushID == req.PushId {
		return
	}
	if err = ent.Database.Rider.UpdateOneID(u.ID).SetPushID(req.PushId).Exec(context.Background()); err != nil {
		return err
	}
	return
}

// EncryptContractDocId 加密合同编号 合同有用户敏感信息
func (b *riderBiz) EncryptContractDocId(contractNo string) string {
	return utils.Md5Base64String(contractNo)
}

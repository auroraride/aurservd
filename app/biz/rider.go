package biz

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/golang-module/carbon/v2"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/purchaseorder"
	"github.com/auroraride/aurservd/internal/ent/purchasepayment"
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
func (b *riderBiz) Signin(device *model.Device, req *definition.RiderSignupReq) (res *definition.RiderSigninRes, err error) {
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

	// 判定用户是否被封禁 或者被拉黑
	if service.NewRider().IsBanned(u) || service.NewRider().IsBlocked(u) {
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

	res, err = b.Profile(u, device, token)
	if err != nil {
		return nil, err
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
	if u.PushID == req.PushId || req.PushId == "" {
		return
	}
	if err = ent.Database.Rider.UpdateOneID(u.ID).SetPushID(req.PushId).Exec(context.Background()); err != nil {
		return err
	}
	return
}

// Profile 获取用户资料
func (b *riderBiz) Profile(u *ent.Rider, device *model.Device, token string) (*definition.RiderSigninRes, error) {
	s := service.NewRider()
	subd, sub := service.NewSubscribe().RecentDetail(u.ID)
	profile := &definition.RiderSigninRes{
		ID:              u.ID,
		Phone:           u.Phone,
		IsNewDevice:     s.IsNewDevice(u, device),
		IsContactFilled: u.Contact != nil,
		IsAuthed:        s.IsAuthed(u),
		Contact:         u.Contact,
		Qrcode:          s.GetQrcode(u.ID),
		Token:           token,
	}

	// 站点
	stat := u.Edges.Station
	if stat != nil {
		profile.Station = &model.EnterpriseStation{
			ID:   stat.ID,
			Name: stat.Name,
		}
	}

	// 订阅
	if sub != nil && slices.Contains(model.SubscribeNotUnSubscribed(), sub.Status) {
		profile.Subscribe = subd
		profile.CabinetBusiness = u.EnterpriseID == nil && sub.BrandID == nil
	}
	if u.Edges.Person != nil {
		profile.Name = u.Edges.Person.Name
	}
	en := u.Edges.Enterprise
	if en != nil {
		profile.Enterprise = &model.Enterprise{
			ID:    en.ID,
			Name:  en.Name,
			Agent: en.Agent,
		}
		profile.UseStore = !en.Agent || en.UseStore
		if en.Agent {
			profile.EnterpriseContact = &model.EnterpriseContact{
				Name:  en.ContactName,
				Phone: en.ContactPhone,
			}
		}
		// 判断是否能退出团签
		if subd != nil && (subd.Status == model.SubscribeStatusInactive || subd.Status == model.SubscribeStatusUnSubscribed) {
			profile.ExitEnterprise = true
		}

	} else {
		profile.OrderNotActived = silk.Bool(subd != nil && subd.Status == model.SubscribeStatusInactive)
		profile.Deposit = s.Deposit(u.ID)
		profile.UseStore = true
	}

	effectiveContract := service.NewContract().QueryEffectiveContract(u)
	if effectiveContract != nil {
		encryptDocID, err := utils.EncryptAES([]byte(ar.Config.Contract.EncryptKey), effectiveContract.DocID)
		if err != nil || encryptDocID == "" {
			zap.L().Error("加密合同编号失败", zap.Error(err))
			return nil, err
		}
		profile.ContractDocID = encryptDocID
	}

	// 待支付购车订单
	profile.Purchase = b.PurchaseObligation(u.ID)

	return profile, nil
}

// PurchaseObligation 查询骑手是否含有待支付订单 (还款时间前3天APP首页提示还款)
func (b *riderBiz) PurchaseObligation(rId uint64) bool {
	p, _ := ent.Database.PurchaseOrder.QueryNotDeleted().
		Where(
			purchaseorder.ID(rId),
			purchaseorder.HasPaymentsWith(
				purchasepayment.StatusEQ(purchasepayment.StatusObligation),
				purchasepayment.BillingDateLTE(carbon.Now().StartOfDay().AddDays(3).StdTime()),
			),
		).Exist(context.Background())
	return p
}

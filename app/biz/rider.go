package biz

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/baidu"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/allocate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/snag"
)

type riderBiz struct {
	orm *ent.RiderClient
	ctx context.Context
}

func NewRiderBiz() *riderBiz {
	return &riderBiz{
		orm: ent.Database.Rider,
		ctx: context.Background(),
	}
}

// Direction 骑手路径规划
func (b *riderBiz) Direction(req *definition.RiderDirectionReq) (*definition.RiderDirectionRes, error) {
	direction, err := baidu.NewDirection().GetDirection(req.Origin, req.Destination)
	if err != nil {
		return nil, err
	}
	return &definition.RiderDirectionRes{
		Routes:      direction.Result.Routes,
		Origin:      direction.Result.Origin,
		Destination: direction.Result.Destination,
	}, nil
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

// Allocated V2测试签约写的分配信息 TODO 记得删除
func (b *riderBiz) Allocated(req *definition.RiderAllocatedReq) {
	r := service.NewRider().QueryPhoneX(req.Phone)

	// 是否有生效中套餐
	sub := service.NewSubscribe().Recent(r.ID)
	if sub == nil {
		snag.Panic("无生效中的骑行卡")
	}

	cab := service.NewCabinet().QueryOneSerialX(req.Serial)

	// 查询分配信息是否存在, 如果存在则删除
	service.NewAllocate().SubscribeDeleteIfExists(sub.ID)

	// 存储分配信息
	err := ent.Database.Allocate.Create().
		SetType(allocate.TypeBattery).
		SetSubscribe(sub).
		SetRider(r).
		SetStatus(model.AllocateStatusPending.Value()).
		SetTime(time.Now()).
		SetModel(sub.Model).
		SetCabinetID(cab.ID).
		SetRemark("用户自主扫码").
		Exec(b.ctx)
	if err != nil {
		snag.Panic("请求失败")
	}
	return
}

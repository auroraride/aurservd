package biz

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/baidu"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/rider"
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

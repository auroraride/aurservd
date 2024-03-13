package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/internal/baidu"
	"github.com/auroraride/aurservd/internal/ent"
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
func (*riderBiz) Direction(req *definition.RiderDirectionReq) (*definition.RiderDirectionRes, error) {
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

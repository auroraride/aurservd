// Created at 2024-03-04

package biz

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/branch"
	"github.com/auroraride/aurservd/internal/ent/business"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/pkg/snag"
)

type cabinetBiz struct {
	orm *ent.CabinetClient
	ctx context.Context
}

func NewCabinet() *cabinetBiz {
	return &cabinetBiz{
		orm: ent.Database.Cabinet,
		ctx: context.Background(),
	}
}

// ListByRider  查询电柜
func (s *cabinetBiz) ListByRider(rid *ent.Rider, req *definition.CabinetByRiderReq) (res []definition.CabinetByRiderRes) {
	q := s.orm.QueryNotDeleted().WithModels().WithEnterprise().WithBranch().
		Modify(func(sel *sql.Selector) {
			sel.
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(%s, ST_GeogFromText('POINT(%f %f)'))`, branch.FieldGeom, *req.Lng, *req.Lat)), "distance").
				OrderBy(sql.Asc("distance"))
			if req.Distance != nil {
				if *req.Distance > 100000 {
					snag.Panic("请求距离太远")
				}
				sel.Where(sql.P(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`ST_DWithin(%s, ST_GeogFromText('POINT(%f %f)'), %f)`, cabinet.FieldGeom, *req.Lng, *req.Lat, *req.Distance))
				}))
			}
		})
	// 默认查询骑手订阅型号的电柜
	var sub *ent.Subscribe
	var rev *model.ReserveUnfinishedRes
	if rid != nil {
		sub = service.NewSubscribeWithRider(rid).Recent(rid.ID)
		if sub != nil {
			req.Model = &sub.Model
		}
		// 预约
		rev = service.NewReserveWithRider(rid).RiderUnfinishedDetail(rid.ID)
	}

	if req.Model != nil {
		q.Where(cabinet.HasModelsWith(batterymodel.Model(*req.Model)))
	}

	if req.Keyword != nil {
		q.Where(cabinet.NameContains(*req.Keyword))
	}

	cabinets := q.AllX(s.ctx)

	// 电柜id
	var cabIDs []uint64

	// 预约数量map
	var rm map[uint64]int

	if req.Business != nil {
		for _, c := range cabinets {
			cabIDs = append(cabIDs, c.ID)
		}
		rm = service.NewReserve().CabinetCounts(cabIDs, business.Type(*req.Business))
	}

	service.NewCabinet().SyncCabinets(cabinets)
	res = make([]definition.CabinetByRiderRes, 0)
	for _, c := range cabinets {
		resvcheck := true
		if req.Business != nil && c.ReserveAble(business.Type(*req.Business), rm[c.ID]) {
			resvcheck = sub == nil || service.NewCabinet().ModelInclude(c, sub.Model)
		}

		if model.CabinetStatus(c.Status) == model.CabinetStatusNormal && resvcheck {
			cdr := definition.CabinetByRiderRes{
				CabinetDataRes: model.CabinetDataRes{
					ID:         c.ID,
					Name:       c.Name,
					Serial:     c.Serial,
					Brand:      c.Brand,
					BatteryNum: c.BatteryNum,
					BinNum:     c.Doors,
				},
				Lng:     c.Lng,
				Lat:     c.Lat,
				Address: c.Address,
				Reserve: nil,
			}

			if rev != nil && rev.CabinetID == c.ID {
				cdr.Reserve = rev
			}

			if sub != nil && service.NewCabinet().ModelInclude(c, sub.Model) {
				// 获取可办理业务
				switch sub.Status {
				case model.SubscribeStatusInactive:
					// 未激活时仅能办理激活业务
					cdr.Businesses = []string{business.TypeActive.String()}
				case model.SubscribeStatusPaused:
					// 寄存中时仅能办理取消寄存业务
					cdr.Businesses = []string{business.TypeContinue.String()}
				case model.SubscribeStatusUsing:
					// 使用中可办理寄存和退租业务
					cdr.Businesses = []string{business.TypePause.String(), business.TypeUnsubscribe.String()}
				}
			}

			if c.Edges.Branch != nil {
				cdr.BranchID = c.Edges.Branch.ID
				cdr.Fid = service.NewBranch().EncodeFacility(nil, c)
			}

			distance, err := c.Value("distance")
			if err != nil {
				snag.Panic(err)
			}
			cdr.Distance = distance.(float64)

			bms := c.Edges.Models
			if len(bms) > 0 {
				cdr.Model = regexp.MustCompile(`(?m)(\d+)V\d+AH`).ReplaceAllString(bms[0].Model, "${1}V")
			}

			// 5分钟未更新视为离线
			if c.Health == model.CabinetHealthStatusOnline && time.Since(c.UpdatedAt).Minutes() < 5 {
				cdr.Online = true
			}

			cdr.Bins = make([]model.CabinetDataBin, len(c.Bin))
			for i, bin := range c.Bin {
				if bin.Battery {
					if bin.Full {
						cdr.Bins[i].Status = model.CabinetDataBinStatusFull
						cdr.FullNum += 1
					} else {
						cdr.Bins[i].Status = model.CabinetDataBinStatusCharging
					}
				} else {
					cdr.Bins[i].Status = model.CabinetDataBinStatusEmpty
					cdr.EmptyNum += 1
				}

				if !bin.DoorHealth {
					cdr.Bins[i].Status = model.CabinetDataBinStatusLock
					cdr.Bins[i].Remark = bin.Remark
					cdr.LockNum += 1
				}
			}
			res = append(res, cdr)
		}
	}
	return
}

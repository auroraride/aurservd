// Created at 2024-03-04

package biz

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/branch"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/es"
	"github.com/auroraride/aurservd/pkg/silk"
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
func (b *cabinetBiz) ListByRider(rid *ent.Rider, req *definition.CabinetByRiderReq) (res []definition.CabinetByRiderRes, err error) {
	q := b.orm.QueryNotDeleted().WithModels().WithEnterprise().WithBranch().
		Modify(func(sel *sql.Selector) {
			sel.
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(%s, ST_GeogFromText('POINT(%f %f)'))`, branch.FieldGeom, *req.Lng, *req.Lat)), "distance").
				OrderBy(sql.Asc("distance"))
			if req.Distance != nil {
				if *req.Distance > 100000 {
					*req.Distance = 100000
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
		// 预约
		rev = service.NewReserveWithRider(rid).RiderUnfinishedDetail(rid.ID)
	}

	if req.Model != nil {
		q.Where(cabinet.HasModelsWith(batterymodel.Model(*req.Model)))
	}

	if req.Keyword != nil {
		q.Where(cabinet.NameContains(*req.Keyword))
	}

	cabinets := q.AllX(b.ctx)

	// 电柜id
	var cabIDs []uint64

	// 预约数量map
	var rm map[model.ReserveBusinessKey]int

	if req.Business != nil {
		for _, c := range cabinets {
			cabIDs = append(cabIDs, c.ID)
		}
		rm = service.NewReserve().CabinetCounts(cabIDs)
	}

	service.NewCabinet().SyncCabinets(cabinets)
	res = make([]definition.CabinetByRiderRes, 0)
	for _, c := range cabinets {
		resvcheck := req.Business == nil
		if req.Business != nil && c.ReserveAble(model.BusinessType(*req.Business), rm) {
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
					Online:     c.Health == model.CabinetHealthStatusOnline,
				},
				Lng:        c.Lng,
				Lat:        c.Lat,
				Address:    c.Address,
				Reserve:    nil,
				Businesses: make([]string, 0),
			}

			if rev != nil && rev.CabinetID == c.ID {
				cdr.Reserve = rev
			}

			// 电柜可办理业务
			reserveNum := service.NewReserve().CabinetCounts([]uint64{c.ID})
			var batteryFullNum, emptyBinNum int
			reserveActiveNum := reserveNum[model.NewReserveBusinessKey(c.ID, model.BusinessTypeActive)]
			reserveContinueNum := reserveNum[model.NewReserveBusinessKey(c.ID, model.BusinessTypeContinue)]
			reservePauseNum := reserveNum[model.NewReserveBusinessKey(c.ID, model.BusinessTypePause)]
			reserveUnsubscribeNum := reserveNum[model.NewReserveBusinessKey(c.ID, model.BusinessTypeUnsubscribe)]

			// 可用电池数
			batteryFullNum = c.BatteryFullNum - reserveActiveNum - reserveContinueNum
			// 可用空仓数
			emptyBinNum = c.EmptyBinNum - reservePauseNum - reserveUnsubscribeNum

			if batteryFullNum >= 2 {
				cdr.Businesses = append(cdr.Businesses, model.BusinessTypeActive.String(), model.BusinessTypeContinue.String())
			}
			if emptyBinNum >= 2 {
				cdr.Businesses = append(cdr.Businesses, model.BusinessTypePause.String(), model.BusinessTypeUnsubscribe.String())
			}

			if c.Edges.Branch != nil {
				cdr.BranchID = c.Edges.Branch.ID
				cdr.Fid = service.NewBranch().EncodeFacility(nil, c)
			}

			var distance ent.Value
			distance, err = c.Value("distance")
			if distance != nil || err == nil {
				distanceFloat, ok := distance.(float64)
				if ok {
					cdr.Distance = distanceFloat
				}
			}

			bms := c.Edges.Models
			if len(bms) > 0 {
				cdr.Model = regexp.MustCompile(`(?m)(\d+)V\d+AH`).ReplaceAllString(bms[0].Model, "${1}V")
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
	// 排序 当有预约数据排在最前面
	sort.Slice(res, func(i, j int) bool {
		return res[i].Reserve != nil && res[j].Reserve == nil
	})

	return
}

// DetailBySerial  通过serial获取电柜详情
func (b *cabinetBiz) DetailBySerial(serial string) (res *model.CabinetDetailRes, err error) {
	item, _ := b.orm.QueryNotDeleted().
		Where(cabinet.Serial(serial)).
		WithModels().
		WithEnterprise().
		WithStation().
		First(b.ctx)
	if item == nil {
		return nil, errors.New("电柜不存在")
	}
	// 同步电柜并返回电柜详情
	service.NewCabinet().Sync(item)
	return service.NewCabinet().Detail(item), nil
}

// ListECData 查询电柜电耗数据
func (b *cabinetBiz) ListECData(options definition.CabinetECDataSearchOptions) (start, end time.Time, data []*definition.CabinetECData) {
	cfg := ar.Config.Elastic
	s, err := es.NewSearch[definition.CabinetECData](cfg.ApiKey, cfg.EccDatastream, cfg.Addresses)
	data = make([]*definition.CabinetECData, 0)
	if err == nil || s == nil {
		return
	}
	query := &types.Query{
		Match: map[string]types.MatchQuery{},
	}
	// 电柜编号
	if options.Serial != nil {
		query.Match[cabinet.FieldSerial] = types.MatchQuery{Query: *options.Serial}
	}
	// 若开始时间和结束时间都不为空，设置时间范围
	if options.Start != nil && options.End != nil {
		start = *options.Start
		end = *options.End
	}
	// 若开始和结束时间均为空，设置默认时间为当日
	if options.Start == nil && options.End == nil {
		start = carbon.Now().StartOfDay().StdTime()
		end = time.Now()
	}
	// 若开始时间为空，结束时间不为空，设置查询为结束时间当日
	if options.Start == nil && options.End != nil {
		start = carbon.CreateFromStdTime(*options.End).StartOfDay().StdTime()
		end = *options.End
	}
	// 若开始时间不为空，结束时间为空，设置查询为开始时间当日
	if options.Start != nil && options.End == nil {
		start = *options.Start
		end = carbon.CreateFromStdTime(*options.Start).EndOfDay().StdTime()
	}
	// 查询时间范围
	query.Range = map[string]types.RangeQuery{
		es.FieldECCTimestamp: types.DateRangeQuery{
			Gte: silk.String(start.Format(time.RFC3339)),
			Lt:  silk.String(end.Format(time.RFC3339)),
		},
	}
	data = s.DoRequest(&search.Request{
		Query: query,
		Sort: []types.SortCombinations{
			types.SortOptions{SortOptions: map[string]types.FieldSort{
				es.FieldECCTimestamp: {Order: &sortorder.Asc},
			}},
		},
	})
	return
}

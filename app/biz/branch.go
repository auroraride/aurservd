// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-05-27, by Jorjan

package biz

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/speps/go-hashids/v2"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/branch"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/snag"
)

type branchBiz struct {
	orm   *ent.BranchClient
	ctx   context.Context
	rider *ent.Rider
}

func NewBranch() *branchBiz {
	return &branchBiz{
		orm: ent.Database.Branch,
		ctx: context.Background(),
	}
}

func NewBranchWithRider(r *ent.Rider) *branchBiz {
	s := NewBranch()
	if r != nil {
		s.rider = r
	}
	return s
}

// ListByDistanceRider 根据距离列出所有网点和电柜
func (s *branchBiz) ListByDistanceRider(req *definition.BranchWithDistanceReq) (items []*definition.BranchWithDistanceRes) {
	var sub *ent.Subscribe
	if s.rider != nil {
		sub, _ = service.NewSubscribeWithRider(s.rider).QueryEffective(s.rider.ID)
	}

	// TODO 业务获取限制
	// if sub != nil {
	//     if req.Business == business.TypeActive && sub.Status != model.SubscribeStatusInactive {}
	//     if req.Business == business.TypePause && sub.Status != model.SubscribeStatusUsing {}
	//     if req.Business == business.TypeContinue && sub.Status != model.SubscribeStatusPaused {}
	//     if req.Business == business.TypeUnsubscribe && sub.Status != model.SubscribeStatusUsing {}
	// }

	temps, stores, cabinets := s.ListByDistance(req)

	items = make([]*definition.BranchWithDistanceRes, 0)
	// 三种设备类别
	itemsMap := make(map[uint64]*definition.BranchWithDistanceRes, len(temps))
	for _, temp := range temps {
		itemsMap[temp.ID] = &definition.BranchWithDistanceRes{
			ID:          temp.ID,
			Distance:    temp.Distance,
			Name:        temp.Name,
			Lng:         temp.Lng,
			Lat:         temp.Lat,
			Image:       temp.Image,
			Photos:      temp.Photos,
			Address:     temp.Address,
			Facility:    make([]*definition.BranchFacility, 0),
			FacilityMap: make(map[string]*definition.BranchFacility),
			Businesses:  make([]string, 0),
		}
	}

	// 门店
	for _, es := range stores {
		// 判断营业状态
		var fState uint
		switch es.Status {
		case model.StoreStatusOpen:
			fState = model.BranchFacilityStateOnline
		case model.StoreStatusClose:
			fState = model.BranchFacilityStateOffline
		default:
			continue
		}

		s.facility(itemsMap[es.BranchID].FacilityMap, definition.BranchFacility{
			ID:    es.ID,
			Type:  model.BranchFacilityTypeStore,
			Name:  es.Name,
			State: fState,
			Num:   0,
			Fid:   s.EncodeFacility(es, nil),
		})
	}

	// 电柜id
	var cabIDs []uint64
	// 预约数量map
	var rm map[model.ReserveBusinessKey]int

	// 每个网点可用业务
	branchBusinessesMap := make(map[uint64]map[uint64][]string)

	if req.Business != "" {
		for _, c := range cabinets {
			cabIDs = append(cabIDs, c.ID)
		}

		rm = service.NewReserve().CabinetCounts(cabIDs)
	}

	// 电柜
	// 同步电柜信息
	if req.StoreStatus == nil && req.StoreBusiness == nil {
		service.NewCabinet().SyncCabinets(cabinets)
		for _, c := range cabinets {
			// 预约检查 = 非预约筛选 或 电柜可满足预约并且如果订阅非空则电柜电池型号满足订阅型号
			// resvcheck := req.Business == "" || (c.ReserveAble(business.Type(req.Business), rm[c.ID]) && (sub == nil || NewCabinet().ModelInclude(c, sub.Model)))
			resvcheck := req.Business == ""
			if c.ReserveAble(model.BusinessType(req.Business), rm) {
				resvcheck = sub == nil || service.NewCabinet().ModelInclude(c, sub.Model)
			}

			if model.CabinetStatus(c.Status) == model.CabinetStatusNormal && resvcheck {
				fa := definition.BranchFacility{
					ID:         c.ID,
					Name:       c.Name,
					State:      model.BranchFacilityStateOffline,
					Type:       model.BranchFacilityTypeV72,
					Fid:        s.EncodeFacility(nil, c),
					CabinetNum: 1,
				}
				// 获取健康状态
				// 5分钟未更新视为离线
				if c.Health == model.CabinetHealthStatusOnline && time.Since(c.UpdatedAt).Minutes() < 5 {
					fa.State = model.BranchFacilityStateOnline
				}
				// 计算可用电池数量
				for _, bin := range c.Bin {
					fa.Total += 1
					// TODO 替换
					if bin.Electricity.IsBatteryFull() {
						fa.Num += 1
					}
					if bin.Battery {
						fa.BatteryNum += 1
					}
				}

				bms := c.Edges.Models
				if len(bms) < 1 {
					continue
				}

				// 判定电池型号
				// 如果有多个电压怎么办? 忽略, 使用第一个
				str := strings.ToUpper(bms[0].Model)

				if strings.HasPrefix(str, "60V") {
					fa.Type = model.BranchFacilityTypeV60
				}

				s.facility(itemsMap[*c.BranchID].FacilityMap, fa)

				// 电柜可办理业务
				if branchBusinessesMap[*c.BranchID] == nil {
					branchBusinessesMap[*c.BranchID] = make(map[uint64][]string)
				}

				// active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
				reserveNum := service.NewReserve().CabinetCounts([]uint64{c.ID})
				// 电柜可办理业务
				var batteryFullNum, emptyBinNum int
				reserveActiveNum := reserveNum[model.NewReserveBusinessKey(c.ID, model.BusinessTypeActive)]
				reserveContinueNum := reserveNum[model.NewReserveBusinessKey(c.ID, model.BusinessTypeContinue)]
				reservePauseNum := reserveNum[model.NewReserveBusinessKey(c.ID, model.BusinessTypePause)]
				reserveUnsubscribeNum := reserveNum[model.NewReserveBusinessKey(c.ID, model.BusinessTypeUnsubscribe)]

				// 可用电池数
				batteryFullNum = c.BatteryFullNum - reserveActiveNum - reserveContinueNum
				// 可用空仓数
				emptyBinNum = c.EmptyBinNum - reservePauseNum - reserveUnsubscribeNum

				// 电柜可办业务展示规则：
				//  1）激活：电柜可用电池数 ≥ 2
				//  2）退租：电柜空仓数 ≥ 2
				//  3）寄存：电柜空仓数 ≥ 2
				//  4）结束寄存：电柜可用电池数 ≥ 2
				if batteryFullNum >= 2 {
					branchBusinessesMap[*c.BranchID][c.ID] = append(branchBusinessesMap[*c.BranchID][c.ID], model.BusinessTypeActive.String(), model.BusinessTypeContinue.String())
				}
				if emptyBinNum >= 2 {
					branchBusinessesMap[*c.BranchID][c.ID] = append(branchBusinessesMap[*c.BranchID][c.ID], model.BusinessTypePause.String(), model.BusinessTypeUnsubscribe.String())
				}
			}
		}
	}

	// 网点业务
	for k, businessesMap := range branchBusinessesMap {
		added := make(map[string]bool)
		businesses := make([]string, 0)
		// 遍历 businessesMap
		for _, b := range businessesMap {
			// 遍历每个切片中的元素
			for _, item := range b {
				// 如果元素不在 added map 中，则添加到 businesses 切片中，并将其标记为已添加
				if !added[item] {
					businesses = append(businesses, item)
					added[item] = true
				}
			}
		}
		itemsMap[k].Businesses = businesses
	}

	for _, m := range itemsMap {
		for _, fa := range m.FacilityMap {
			m.Facility = append(m.Facility, fa)
		}
		if len(m.Facility) > 0 {
			items = append(items, m)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Distance < items[j].Distance
	})
	return
}

type branchListTemp struct {
	ID       uint64   `json:"id"`
	Distance float64  `json:"distance"`
	Name     string   `json:"name"`
	Lng      float64  `json:"lng"`
	Lat      float64  `json:"lat"`
	Image    string   `json:"image"`
	Photos   []string `json:"photos"`
	Address  string   `json:"address"`
}

func (s *branchBiz) ListByDistance(req *definition.BranchWithDistanceReq) (temps []branchListTemp, stores []*ent.Store, cabinets []*ent.Cabinet) {
	if req.Distance == nil && req.CityID == nil {
		snag.Panic("距离和城市不能同时为空")
	}
	q := s.orm.QueryNotDeleted().
		Modify(func(sel *sql.Selector) {
			bt := sql.Table(branch.Table)
			sel.Select(bt.C(branch.FieldID), bt.C(branch.FieldName), bt.C(branch.FieldAddress), bt.C(branch.FieldLat), bt.C(branch.FieldLng), bt.C(branch.FieldPhotos)).
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(%s, ST_GeogFromText('POINT(%f %f)'))`, branch.FieldGeom, *req.Lng, *req.Lat)), "distance").
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`TRIM('"' FROM %s[0]::TEXT)`, branch.FieldPhotos)), "image").
				GroupBy(bt.C(branch.FieldID)).
				OrderBy(sql.Asc("distance"))
			if req.Distance != nil {
				if *req.Distance > 100000 {
					*req.Distance = 100000
				}
				sel.Where(sql.P(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`ST_DWithin(%s, ST_GeogFromText('POINT(%f %f)'), %f)`, branch.FieldGeom, *req.Lng, *req.Lat, *req.Distance))
				}))
			} else if req.CityID != nil {
				sel.Where(sql.EQ(sel.C(branch.FieldCityID), *req.CityID))
			}
		})
	if req.Model != nil {
		q.Where(branch.HasCabinetsWith(cabinet.HasModelsWith(batterymodel.Model(*req.Model))))
	}
	err := q.Scan(s.ctx, &temps)
	if err != nil || len(temps) == 0 {
		return
	}
	ids := make([]uint64, len(temps))
	for i, temp := range temps {
		ids[i] = temp.ID
	}

	// 电柜查询
	// 当没有传入门店相关参数时需查询电柜数据
	if req.StoreStatus == nil && req.StoreBusiness == nil {
		cabQuery := ent.Database.Cabinet.QueryNotDeleted().Where(cabinet.BranchIDIn(ids...)).WithModels()
		if req.Model != nil {
			cabQuery.Where(cabinet.HasModelsWith(batterymodel.Model(*req.Model)))
		}
		cabinets = cabQuery.AllX(s.ctx)
	}

	// 门店查询
	// 当没有传入电池型号、电柜业务参数时需查询门店数据
	if req.Model == nil && req.Business == "" {
		storeQuery := ent.Database.Store.QueryNotDeleted().Where(store.BranchIDIn(ids...))
		storeQuery.Where(
			store.Rest(true),
			store.StatusIn(model.StoreStatusOpen, model.StoreStatusClose),
		)

		if req.StoreStatus != nil {
			switch *req.StoreStatus {
			case model.StoreStatusOpen:
				storeQuery.Where(store.Status(model.StoreStatusOpen))
			case model.StoreStatusClose:
				storeQuery.Where(store.Status(model.StoreStatusClose))
			default:
			}
		}

		if req.StoreBusiness != nil {
			switch *req.StoreBusiness {
			case model.StoreBusinessTypeObtain:
				storeQuery.Where(store.EbikeObtain(true))
			case model.StoreBusinessTypeRepair:
				storeQuery.Where(store.EbikeRepair(true))
			case model.StoreBusinessTypeSale:
				storeQuery.Where(store.EbikeSale(true))
			case model.StoreBusinessTypeRest:
				storeQuery.Where(store.Rest(true))
			}
		}
		stores = storeQuery.AllX(s.ctx)
	}
	return
}

func (s *branchBiz) facility(mp map[string]*definition.BranchFacility, info definition.BranchFacility) {
	fa, ok := mp[info.Type]
	if ok {
		// 合并电柜满电数量
		if info.Type != model.BranchFacilityTypeStore {
			fa.Num += info.Num
			fa.Total += info.Total
			fa.CabinetNum += info.CabinetNum
		}
	} else {
		fa = &info
		mp[info.Type] = fa
	}
}

func (s *branchBiz) Hasher() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = "branch facility"
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	return h
}

// EncodeFacility 加密设施
func (s *branchBiz) EncodeFacility(sto *ent.Store, cab *ent.Cabinet) string {
	if sto == nil && cab == nil {
		return ""
	}
	if sto != nil {
		return s.EncodeStoreID(sto.ID)
	}
	return s.EncodeCabinetID(cab.ID)
}

func (s *branchBiz) EncodeStoreID(storeID uint64) (fid string) {
	fid, _ = s.Hasher().EncodeInt64([]int64{1, int64(storeID)})
	return
}

func (s *branchBiz) EncodeCabinetID(CabinetID uint64) (fid string) {
	fid, _ = s.Hasher().EncodeInt64([]int64{2, int64(CabinetID)})
	return
}

// DecodeFacility 解码设施
func (s *branchBiz) DecodeFacility(fid string) (b *ent.Branch, sto *ent.Store, cabs []*ent.Cabinet) {
	arr, _ := s.Hasher().DecodeInt64WithError(fid)
	if len(arr) != 2 {
		snag.Panic("查询失败")
	}
	switch arr[0] {
	case 1:
		sto = service.NewStore().Query(uint64(arr[1]))
		b, _ = sto.QueryBranch().First(s.ctx)
	case 2:
		cab, _ := ent.Database.Cabinet.QueryNotDeleted().
			WithModels().
			Where(cabinet.ID(uint64(arr[1]))).
			First(s.ctx)
		if cab == nil {
			break
		}
		b, _ = cab.QueryBranch().First(s.ctx)
		if b == nil {
			break
		}
		// 查询其他电柜信息
		items, _ := ent.Database.Cabinet.QueryNotDeleted().
			WithModels().
			Where(
				cabinet.BranchID(b.ID),
				cabinet.IDNEQ(cab.ID),
				cabinet.Status(model.CabinetStatusNormal.Value()),
				// 离线也展示
				// cabinet.Health(definition.CabinetHealthStatusOnline),
			).
			All(s.ctx)
		cabs = append([]*ent.Cabinet{cab}, items...)
	}
	if b == nil || (sto == nil && len(cabs) == 0) {
		snag.Panic("查询失败")
	}
	return
}

// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/LucaTheHacker/go-haversine"
	"github.com/jinzhu/copier"
	"github.com/speps/go-hashids/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/battery"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/branch"
	"github.com/auroraride/aurservd/internal/ent/branchcontract"
	"github.com/auroraride/aurservd/internal/ent/business"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type branchService struct {
	orm   *ent.BranchClient
	ctx   context.Context
	rider *ent.Rider
}

func NewBranch() *branchService {
	return &branchService{
		orm: ent.Database.Branch,
		ctx: context.Background(),
	}
}

func NewBranchWithModifier(m *model.Modifier) *branchService {
	s := NewBranch()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	return s
}

func NewBranchWithRider(r *ent.Rider) *branchService {
	s := NewBranch()
	if r != nil {
		s.rider = r
	}
	return s
}

// Query 根据ID查询网点
func (s *branchService) Query(id uint64) *ent.Branch {
	item, err := s.orm.QueryNotDeleted().Where(branch.ID(id)).First(s.ctx)
	if err != nil {
		snag.Panic("未找到有效网点")
	}
	return item
}

// Create 新增网点
func (s *branchService) Create(req *model.BranchCreateReq) {
	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) error {
		b, err := s.orm.Create().
			SetName(*req.Name).
			SetAddress(*req.Address).
			SetCityID(*req.CityID).
			SetLng(*req.Lng).
			SetLat(*req.Lat).
			SetGeom(&model.Geometry{
				Lng: *req.Lng,
				Lat: *req.Lat,
			}).
			SetPhotos(req.Photos).
			Save(s.ctx)
		if err != nil {
			return err
		}

		if len(req.Contracts) > 0 {
			for _, contract := range req.Contracts {
				s.AddContract(b.ID, contract)
			}
		}
		return nil
	})
}

// AddContract 新增合同
func (s *branchService) AddContract(id uint64, req *model.BranchContract) *ent.BranchContract {
	return ent.Database.BranchContract.Create().
		SetBranchID(id).
		SetLandlordName(req.LandlordName).
		SetIDCardNumber(req.IDCardNumber).
		SetPhone(req.Phone).
		SetBankNumber(req.BankNumber).
		SetPledge(req.Pledge).
		SetRent(req.Rent).
		SetLease(req.Lease).
		SetElectricityPledge(req.ElectricityPledge).
		SetElectricity(req.Electricity).
		SetArea(req.Area).
		SetStartTime(tools.NewTime().ParseDateStringX(req.StartTime)).
		SetEndTime(tools.NewTime().ParseDateStringX(req.EndTime)).
		SetFile(req.File).
		SetSheets(req.Sheets).
		SaveX(s.ctx)
}

// List 网点列表
func (s *branchService) List(req *model.BranchListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().
		Order(ent.Desc(branch.FieldID)).
		WithCity().
		WithStores(func(sq *ent.StoreQuery) {
			sq.Where(store.DeletedAtIsNil())
		}).
		WithCabinets(func(cq *ent.CabinetQuery) {
			cq.Where(cabinet.DeletedAtIsNil()).WithModels()
		}).
		WithContracts(func(query *ent.BranchContractQuery) {
			query.Order(ent.Desc(branchcontract.FieldID))
		})

	if req.CityID != nil {
		q.Where(branch.CityID(*req.CityID))
	}

	if req.Name != nil {
		q.Where(branch.NameContainsFold(*req.Name))
	}

	return model.ParsePaginationResponse[model.BranchItem, ent.Branch](
		q,
		req.PaginationReq,
		func(item *ent.Branch) model.BranchItem {
			var r model.BranchItem
			_ = copier.Copy(&r, item)

			cs := make([]model.BranchContract, len(item.Edges.Contracts))
			for n, contract := range item.Edges.Contracts {
				var c model.BranchContract
				_ = copier.Copy(&c, contract)
				cs[n] = c
			}

			r.Contracts = cs
			city := item.Edges.City
			r.City = model.City{
				ID:   city.ID,
				Name: city.Name,
			}
			r.StoreTotal = len(item.Edges.Stores)
			for _, c := range item.Edges.Cabinets {
				bms := c.Edges.Models
				if len(bms) > 0 {
					cm := bms[0]
					str := strings.ToUpper(cm.Model)

					if strings.HasPrefix(str, "60V") {
						r.V60Total += 1
					}

					if strings.HasPrefix(str, "72V") {
						r.V72Total += 1
					}
				}
			}
			return r
		})
}

// Modify 修改网点
func (s *branchService) Modify(req *model.BranchModifyReq) {
	b, _ := s.orm.QueryNotDeleted().Where(branch.ID(req.ID)).First(s.ctx)
	if b == nil {
		snag.Panic("网点不存在")
	}

	// 从结构体更新
	q := s.orm.ModifyOne(b, req)

	geom := b.Geom
	if req.Lng != nil {
		geom.Lng = *req.Lng
	}
	if req.Lat != nil {
		geom.Lat = *req.Lat
	}

	nb := q.SetGeom(geom).SaveX(s.ctx)

	if req.Lng != nil || req.Lat != nil || req.Address != nil {
		// 更新电柜地址和坐标
		_ = ent.Database.Cabinet.Update().Where(cabinet.BranchID(req.ID)).SetLng(nb.Lng).SetLat(nb.Lat).SetGeom(nb.Geom).SetAddress(nb.Address).Exec(s.ctx)
		// 更新门店地址和坐标
		_ = ent.Database.Store.Update().Where(store.BranchID(req.ID)).SetLng(nb.Lng).SetLat(nb.Lat).SetAddress(nb.Address).Exec(s.ctx)
	}
}

// Selector 网点选择列表
func (s *branchService) Selector() *model.ItemListRes {
	items := make([]model.BranchSampleItem, 0)
	s.orm.QueryNotDeleted().Select(branch.FieldID, branch.FieldName).ScanX(s.ctx, &items)
	res := new(model.ItemListRes)
	model.SetItemListResItems[model.BranchSampleItem](res, items)
	return res
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

func (s *branchService) ListByDistance(req *model.BranchWithDistanceReq, sub *ent.Subscribe) (temps []branchListTemp, stores []*ent.Store, cabinets []*ent.Cabinet) {
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

	storeQuery := ent.Database.Store.QueryNotDeleted().Where(store.BranchIDIn(ids...))
	filter := req.Filter
	switch filter {
	case model.BranchFacilityTypeStore:
		break
	case model.BranchFacilityFilterEbikeObtain:
		storeQuery.Where(store.EbikeObtain(true))
	case model.BranchFacilityFilterEbikeRepair:
		storeQuery.Where(store.EbikeRepair(true))
	case model.BranchFacilityFilterEbike:
		storeQuery.Where(
			store.Or(
				store.EbikeObtain(true),
				store.EbikeRepair(true),
			),
		)
	default:
		cabQuery := ent.Database.Cabinet.QueryNotDeleted().Where(cabinet.BranchIDIn(ids...)).WithModels()
		if sub != nil {
			cabQuery.Where(cabinet.Intelligent(sub.Intelligent), cabinet.HasModelsWith(batterymodel.Model(sub.Model)))
		}
		cabinets = cabQuery.AllX(s.ctx)
	}
	stores = storeQuery.AllX(s.ctx)
	return
}

func (s *branchService) ListByDistanceManager(req *model.BranchDistanceListReq) (items []*model.BranchDistanceListRes) {
	items = make([]*model.BranchDistanceListRes, 0)
	lng := req.Lng
	if lng == 0 {
		lng = 108.947713
	}
	lat := req.Lat
	if lat == 0 {
		lat = 34.231657
	}
	distance := req.Distance
	if distance == 0 {
		distance = 100000
	}
	temps, stores, cabinets := s.ListByDistance(&model.BranchWithDistanceReq{
		Lng:      &lng,
		Lat:      &lat,
		Distance: &distance,
	}, nil)
	bmap := make(map[uint64]*model.BranchDistanceListRes)
	for _, temp := range temps {
		bmap[temp.ID] = &model.BranchDistanceListRes{
			ID:       temp.ID,
			Distance: temp.Distance,
			Name:     temp.Name,
			Lng:      temp.Lng,
			Lat:      temp.Lat,
			Cabinets: make([]model.CabinetListByDistanceRes, 0),
			Stores:   make([]model.StoreWithStatus, 0),
		}
	}
	if req.Type == 0 || req.Type == 1 {
		for _, st := range stores {
			if strings.Contains(st.Name, req.Name) {
				if b, ok := bmap[st.BranchID]; ok {
					b.Stores = append(b.Stores, model.StoreWithStatus{
						Store: model.Store{
							ID:   st.ID,
							Name: st.Name,
						},
						Status: st.Status,
					})
				}
			}
		}
	}
	if req.Type == 0 || req.Type > 1 {
		var mt string
		if req.Type > 1 {
			switch req.Type {
			case 60:
				mt = "V60"
			case 72:
				mt = "V72"
			}
		}
		for _, cab := range cabinets {
			if strings.Contains(cab.Name, req.Name) {
				var inside bool
				for _, bm := range cab.Edges.Models {
					if strings.HasPrefix(bm.Model, mt) {
						inside = true
						break
					}
				}
				if inside {
					if b, ok := bmap[*cab.BranchID]; ok {
						b.Cabinets = append(b.Cabinets, model.CabinetListByDistanceRes{
							CabinetBasicInfo: model.CabinetBasicInfo{
								ID:     cab.ID,
								Brand:  cab.Brand,
								Serial: cab.Serial,
								Name:   cab.Name,
							},
							Status: cab.Status,
							Health: cab.Health,
						})
					}
				}
			}
		}
	}
	for _, b := range bmap {
		if len(b.Cabinets) > 0 || len(b.Stores) > 0 {
			items = append(items, b)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Distance < items[j].Distance
	})
	return
}

// ListByDistanceRider 根据距离列出所有网点和电柜
func (s *branchService) ListByDistanceRider(req *model.BranchWithDistanceReq, v2 bool) (items []*model.BranchWithDistanceRes) {
	var sub *ent.Subscribe
	if s.rider != nil {
		sub, _ = NewSubscribeWithRider(s.rider).QueryEffective(s.rider.ID)
	}

	// TODO 业务获取限制
	// if sub != nil {
	//     if req.Business == business.TypeActive && sub.Status != model.SubscribeStatusInactive {}
	//     if req.Business == business.TypePause && sub.Status != model.SubscribeStatusUsing {}
	//     if req.Business == business.TypeContinue && sub.Status != model.SubscribeStatusPaused {}
	//     if req.Business == business.TypeUnsubscribe && sub.Status != model.SubscribeStatusUsing {}
	// }

	temps, stores, cabinets := s.ListByDistance(req, sub)

	items = make([]*model.BranchWithDistanceRes, 0)
	// 三种设备类别
	itemsMap := make(map[uint64]*model.BranchWithDistanceRes, len(temps))
	for _, temp := range temps {
		itemsMap[temp.ID] = &model.BranchWithDistanceRes{
			ID:          temp.ID,
			Distance:    temp.Distance,
			Name:        temp.Name,
			Lng:         temp.Lng,
			Lat:         temp.Lat,
			Image:       temp.Image,
			Photos:      temp.Photos,
			Address:     temp.Address,
			Facility:    make([]*model.BranchFacility, 0),
			FacilityMap: make(map[string]*model.BranchFacility),
			Businesses:  make([]string, 0),
		}
	}

	// 进行关联查询
	// 门店
	if req.Business == "" {
		for _, es := range stores {
			if es.Status == model.StoreStatusOpen {
				s.facility(itemsMap[es.BranchID].FacilityMap, model.BranchFacility{
					ID:    es.ID,
					Type:  model.BranchFacilityTypeStore,
					Name:  es.Name,
					State: model.BranchFacilityStateOnline,
					Num:   0,
					Fid:   s.EncodeFacility(es, nil),
				})
			}
		}
	}

	// 电柜id
	var cabIDs []uint64
	// 预约数量map
	var rm map[uint64]int

	// 每个网点可用业务
	branchBusinessesMap := make(map[uint64]map[uint64][]string)

	if req.Business != "" {
		for _, c := range cabinets {
			cabIDs = append(cabIDs, c.ID)
		}

		rm = NewReserve().CabinetCounts(cabIDs, business.Type(req.Business))
	}

	// 电柜
	// 同步电柜信息
	// NewCabinet().SyncCabinets(cabinets)
	for _, c := range cabinets {
		// 预约检查 = 非预约筛选 或 电柜可满足预约并且如果订阅非空则电柜电池型号满足订阅型号
		// resvcheck := req.Business == "" || (c.ReserveAble(business.Type(req.Business), rm[c.ID]) && (sub == nil || NewCabinet().ModelInclude(c, sub.Model)))
		resvcheck := req.Business == ""
		if c.ReserveAble(business.Type(req.Business), rm[c.ID]) {
			resvcheck = sub == nil || NewCabinet().ModelInclude(c, sub.Model)
		}

		if model.CabinetStatus(c.Status) == model.CabinetStatusNormal && resvcheck {
			fa := model.BranchFacility{
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
			// 查询激活预约数
			reserveActiveNum := NewReserve().CabinetCounts([]uint64{c.ID}, "active")
			// 查询结束寄存
			reserveContinueNum := NewReserve().CabinetCounts([]uint64{c.ID}, "continue")
			// 查询寄存
			reservePauseNum := NewReserve().CabinetCounts([]uint64{c.ID}, "pause")
			// 查询退租
			reserveUnsubscribeNum := NewReserve().CabinetCounts([]uint64{c.ID}, "unsubscribe")

			var batteryFullNum, emptyBinNum int

			// 可用电池数
			batteryFullNum = c.BatteryFullNum - reserveActiveNum[c.ID] - reserveContinueNum[c.ID]
			// 可用空仓数
			emptyBinNum = c.EmptyBinNum - reservePauseNum[c.ID] - reserveUnsubscribeNum[c.ID]

			// 电柜可办业务展示规则：
			//  1）激活：电柜可用电池数 ≥ 2
			//  2）退租：电柜空仓数 ≥ 2
			//  3）寄存：电柜空仓数 ≥ 2
			//  4）结束寄存：电柜可用电池数 ≥ 2
			if batteryFullNum >= 2 {
				branchBusinessesMap[*c.BranchID][c.ID] = append(branchBusinessesMap[*c.BranchID][c.ID], business.TypeActive.String(), business.TypeContinue.String())
			}
			if emptyBinNum >= 2 {
				branchBusinessesMap[*c.BranchID][c.ID] = append(branchBusinessesMap[*c.BranchID][c.ID], business.TypePause.String(), business.TypeUnsubscribe.String())
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
			// TODO 最多显示电池数量
			if fa.Num >= 19 && !v2 {
				fa.Num = 19
			}
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

func (s *branchService) facility(mp map[string]*model.BranchFacility, info model.BranchFacility) {
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

func (s *branchService) Sheet(req *model.BranchContractSheetReq) {
	bc, _ := ent.Database.BranchContract.QueryNotDeleted().Where(branchcontract.ID(req.ID)).First(s.ctx)
	if bc == nil {
		snag.Panic("未找到合同信息")
	}
	bc.Update().SetSheets(req.Sheets).SaveX(s.ctx)
}

func (s *branchService) Hasher() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = "branch facility"
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	return h
}

// EncodeFacility 加密设施
func (s *branchService) EncodeFacility(sto *ent.Store, cab *ent.Cabinet) string {
	if sto == nil && cab == nil {
		return ""
	}
	if sto != nil {
		return s.EncodeStoreID(sto.ID)
	}
	return s.EncodeCabinetID(cab.ID)
}

func (s *branchService) EncodeStoreID(storeID uint64) (fid string) {
	fid, _ = s.Hasher().EncodeInt64([]int64{1, int64(storeID)})
	return
}

func (s *branchService) EncodeCabinetID(CabinetID uint64) (fid string) {
	fid, _ = s.Hasher().EncodeInt64([]int64{2, int64(CabinetID)})
	return
}

// DecodeFacility 解码设施
func (s *branchService) DecodeFacility(fid string) (b *ent.Branch, sto *ent.Store, cabs []*ent.Cabinet) {
	arr, _ := s.Hasher().DecodeInt64WithError(fid)
	if len(arr) != 2 {
		snag.Panic("查询失败")
	}
	switch arr[0] {
	case 1:
		sto = NewStore().Query(uint64(arr[1]))
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
				// cabinet.Health(model.CabinetHealthStatusOnline),
			).
			All(s.ctx)
		cabs = append([]*ent.Cabinet{cab}, items...)
	}
	if b == nil || (sto == nil && len(cabs) == 0) {
		snag.Panic("查询失败")
	}
	return
}

// Facility 获取设施详情
func (s *branchService) Facility(req *model.BranchFacilityReq) (data model.BranchFacilityRes) {
	b, sto, cabs := s.DecodeFacility(req.Fid)
	distance := haversine.Distance(haversine.NewCoordinates(req.Lat, req.Lng), haversine.NewCoordinates(b.Lat, b.Lng))
	data = model.BranchFacilityRes{
		Name:     b.Name,
		Address:  b.Address,
		Lng:      b.Lng,
		Lat:      b.Lat,
		Distance: distance.Kilometers() * 1000.0,
		Image:    b.Photos[0],
		Photos:   b.Photos,
	}
	if sto != nil {
		data.Type = "store"
		// 查询门店电池库存
		ins := NewStock().StoreCurrent(sto.ID)
		models := make([]string, 0)
		for _, in := range ins {
			if in.Num > 0 && in.Battery {
				models = append(models, in.Model)
			}
		}
		data.Store = &model.BranchFacilityStore{
			Models: models,
			Name:   sto.Name,
		}
	} else {
		// 同步电柜信息
		NewCabinet().SyncCabinets(cabs)
		// 订阅
		var sub *ent.Subscribe
		// 预约
		var rev *model.ReserveUnfinishedRes
		// 当骑手登录时, 获取骑手的订阅信息
		if s.rider != nil {
			sub = NewSubscribeWithRider(s.rider).Recent(s.rider.ID)
			rev = NewReserveWithRider(s.rider).RiderUnfinishedDetail(s.rider.ID)
		}

		// 设施详情 - 电柜
		data.Type = "cabinet"
		for _, cab := range cabs {
			bms := cab.Edges.Models
			if len(bms) == 0 {
				continue
			}
			bm := battery.New(bms[0].Model)
			if !bm.IsVaild() {
				continue
			}

			// 电池状态
			batInfo := model.BranchFacilityCabinetBattery{
				Voltage:  bm.Voltage,
				Capacity: bm.Capacity,
			}
			c := model.BranchFacilityCabinet{
				ID:                cab.ID,
				Name:              cab.Name,
				Serial:            cab.Serial,
				Reserve:           nil,
				Bins:              make([]model.BranchFacilityCabinetBin, len(cab.Bin)),
				Businesses:        make([]string, 0),
				CabinetBusinesses: make([]string, 0),
			}

			// 获取电柜状态
			if cab.Status == model.CabinetStatusNormal.Value() {
				if cab.Health == model.CabinetHealthStatusOnline {
					c.Status = 1
				} else {
					c.Status = 0
				}
			} else {
				c.Status = 2
			}

			// 获取仓位详情
			for bi, bin := range cab.Bin {
				// 锁仓
				if !bin.DoorHealth {
					c.Bins[bi] = model.BranchFacilityCabinetBin{
						Status: 3,
					}
				} else {
					// 有电池
					if bin.Battery {
						c.Bins[bi] = model.BranchFacilityCabinetBin{
							Electricity: silk.Pointer(bin.Electricity),
						}
						if bin.Electricity.IsBatteryFull() {
							// 满电
							c.Bins[bi].Status = 2
							batInfo.Fully += 1
						} else {
							// 充电中
							c.Bins[bi].Status = 1
							batInfo.Charging += 1
						}
					}
				}
				c.Bins[bi].BatterySN = bin.BatterySN
			}

			c.Batteries = []model.BranchFacilityCabinetBattery{batInfo}

			// 当前预约
			if rev != nil && rev.CabinetID == cab.ID {
				c.Reserve = rev
			}

			// 当订阅非空并且订阅电池型号包含在当前电柜中时, 判定可办理业务情况
			if sub != nil && NewCabinet().ModelInclude(cab, sub.Model) {
				// 获取可办理业务
				switch sub.Status {
				case model.SubscribeStatusInactive:
					// 未激活时仅能办理激活业务
					c.Businesses = []string{business.TypeActive.String()}
				case model.SubscribeStatusPaused:
					// 寄存中时仅能办理取消寄存业务
					c.Businesses = []string{business.TypeContinue.String()}
				case model.SubscribeStatusUsing:
					// 使用中可办理寄存和退租业务
					c.Businesses = []string{business.TypePause.String(), business.TypeUnsubscribe.String()}
				}
			}

			// 电柜可办理业务
			reserveActiveNum := NewReserve().CabinetCounts([]uint64{cab.ID}, "active")
			// 查询结束寄存
			reserveContinueNum := NewReserve().CabinetCounts([]uint64{cab.ID}, "continue")
			// 查询寄存
			reservePauseNum := NewReserve().CabinetCounts([]uint64{cab.ID}, "pause")
			// 查询退租
			reserveUnsubscribeNum := NewReserve().CabinetCounts([]uint64{cab.ID}, "unsubscribe")

			var batteryFullNum, emptyBinNum int

			// 可用电池数
			batteryFullNum = cab.BatteryFullNum - reserveActiveNum[c.ID] - reserveContinueNum[c.ID]
			// 可用空仓数
			emptyBinNum = cab.EmptyBinNum - reservePauseNum[c.ID] - reserveUnsubscribeNum[c.ID]

			if batteryFullNum >= 2 {
				c.CabinetBusinesses = append(c.CabinetBusinesses, business.TypeActive.String(), business.TypeContinue.String())
			}
			if emptyBinNum >= 2 {
				c.CabinetBusinesses = append(c.CabinetBusinesses, business.TypePause.String(), business.TypeUnsubscribe.String())
			}

			data.Cabinet = append(data.Cabinet, c)
		}
	}
	return
}

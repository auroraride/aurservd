// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-01
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/allocate"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetattributes"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type ebikeService struct {
	*BaseService
	orm *ent.AssetClient
}

func NewEbike(params ...any) *ebikeService {
	return &ebikeService{
		BaseService: newService(params...),
		orm:         ent.Database.Asset,
	}
}

func (s *ebikeService) Query(id uint64) (*ent.Asset, error) {
	return s.orm.QueryNotDeleted().Where(asset.ID(id)).First(s.ctx)
}

func (s *ebikeService) QueryKeyword(keyword string) (*ent.Asset, error) {
	// 查询电车属性车牌
	attributes, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.Key("plate")).First(s.ctx)
	q := s.orm.QueryNotDeleted().Where(
		asset.Type(model.AssetTypeEbike.Value()),
		asset.Or(
			asset.Sn(keyword),
		),
	)
	if attributes != nil {
		q.Where(asset.Or(
			asset.HasValuesWith(
				assetattributevalues.AttributeID(attributes.ID),
				assetattributevalues.ValueContainsFold(keyword),
			)))
	}
	res, err := q.First(s.ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AllocatableBaseFilter 可分配车辆查询条件 (不包含门店筛选)
func (s *ebikeService) AllocatableBaseFilter() *ent.AssetQuery {
	q := s.orm.Query().Where(
		asset.Enable(true),
		asset.Status(model.AssetStatusStock.Value()),
		asset.LocationsType(model.AssetTypeEbike.Value()),
	).WithValues()
	return q
}

// IsAllocated 电车是否已分配
func (s *ebikeService) IsAllocated(id uint64) bool {
	exists, _ := ent.Database.Allocate.Query().Where(
		allocate.EbikeID(id),
		allocate.Status(model.AllocateStatusPending.Value()),
		allocate.TimeGTE(carbon.Now().SubSeconds(model.AllocateExpiration).StdTime()),
	).Exist(s.ctx)
	return exists
}

//	func (s *ebikeService) listFilter(req model.EbikeListFilter) (q *ent.EbikeQuery, info ar.Map) {
//		info = make(ar.Map)
//
//		q = s.orm.Query().
//			Order(ent.Desc(ebike.FieldCreatedAt)).
//			WithRider().
//			WithStore().
//			WithBrand().
//			WithEnterprise().
//			WithStation()
//
//		// 启用状态
//		if req.Enable != nil {
//			q.Where(ebike.Enable(*req.Enable))
//			if *req.Enable {
//				info["启用"] = "是"
//			} else {
//				info["启用"] = "否"
//			}
//		}
//
//		// 状态
//		if req.Status != nil {
//			info["状态"] = req.Status.String()
//			q.Where(ebike.Status(*req.Status))
//		}
//
//		// 骑手
//		if req.RiderID != 0 {
//			info["骑手"] = ent.NewExportInfo(req.RiderID, rider.Table)
//			q.Where(ebike.RiderID(req.RiderID))
//		}
//
//		// 门店
//		if req.StoreID != 0 {
//			info["门店"] = ent.NewExportInfo(req.StoreID, store.Table)
//			q.Where(ebike.StoreID(req.StoreID))
//		}
//
//		// 品牌
//		if req.BrandID != 0 {
//			info["品牌"] = ent.NewExportInfo(req.BrandID, ebikebrand.Table)
//			q.Where(ebike.BrandID(req.BrandID))
//		}
//
//		// 关键词
//		if req.Keyword != "" {
//			info["关键词"] = req.Keyword
//			q.Where(ebike.Or(
//				ebike.Sn(req.Keyword),
//				ebike.Plate(req.Keyword),
//				ebike.Machine(req.Keyword),
//				ebike.Sim(req.Keyword),
//				ebike.HasRiderWith(rider.Or(
//					rider.Name(req.Keyword),
//					rider.Phone(req.Keyword),
//				)),
//			))
//		}
//
//		// 生产批次
//		if req.ExFactory != "" {
//			info["生产批次"] = req.ExFactory
//			q.Where(ebike.ExFactory(req.ExFactory))
//		}
//		// 归属
//		if req.OwnerType != nil {
//			var OwnerTypeName string
//			switch *req.OwnerType {
//			case 1: // 平台
//				q.Where(ebike.EnterpriseIDIsNil())
//				OwnerTypeName = "平台"
//			case 2: // 代理商
//				q.Where(ebike.EnterpriseIDNotNil())
//				OwnerTypeName = "代理商"
//			}
//			info["归属"] = OwnerTypeName
//		}
//		// 团签
//		if req.EnterpriseID != nil {
//			info["团签"] = ent.NewExportInfo(*req.EnterpriseID, enterprise.Table)
//			q.Where(ebike.EnterpriseID(*req.EnterpriseID))
//		}
//		// 站点
//		if req.StationID != nil {
//			info["站点"] = ent.NewExportInfo(*req.StationID, enterprisestation.Table)
//			q.Where(ebike.StationID(*req.StationID))
//		}
//
//		switch req.Goal {
//		case model.EbikeStation:
//			info["查询目标"] = "站点"
//			q.Where(
//				ebike.StationIDNotNil(),
//				ebike.RiderIDIsNil(),
//			)
//		case model.EbikeRider:
//			info["查询目标"] = "骑手"
//			q.Where(
//				ebike.RiderIDNotNil(),
//			)
//		}
//
//		if req.Rto != nil {
//			switch *req.Rto {
//			case true:
//				info["是否以租代购"] = true
//				q.Where(ebike.RtoRiderIDNotNil())
//			default:
//				info["是否以租代购"] = false
//				q.Where(ebike.RtoRiderIDIsNil())
//			}
//		}
//
//		return
//	}
//
//	func (s *ebikeService) List(req *model.EbikeListReq) *model.PaginationRes {
//		q, _ := s.listFilter(req.EbikeListFilter)
//		return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Ebike) model.EbikeListRes {
//			eb := item.Edges.Brand
//			er := item.Edges.Rider
//			es := item.Edges.Store
//			res := model.EbikeListRes{
//				ID:        item.ID,
//				SN:        item.Sn,
//				BrandID:   item.BrandID,
//				ExFactory: item.ExFactory,
//				Status:    item.Status.String(),
//				EbikeAttributes: model.EbikeAttributes{
//					Enable:  silk.Pointer(item.Enable),
//					Plate:   item.Plate,
//					Machine: item.Machine,
//					Sim:     item.Sim,
//					Color:   silk.Pointer(item.Color),
//				},
//				Rto: item.RtoRiderID != nil,
//			}
//			if eb != nil {
//				res.Brand = eb.Name
//			}
//			if er != nil {
//				res.Rider = fmt.Sprintf("%s-%s", er.Name, er.Phone)
//			}
//			if es != nil {
//				res.Store = es.Name
//			}
//
//			if item.Edges.Station != nil {
//				res.StationName = item.Edges.Station.Name
//				res.StationID = &item.Edges.Station.ID
//			}
//			if item.Edges.Enterprise != nil {
//				res.EnterpriseName = item.Edges.Enterprise.Name
//				res.EnterpriseID = &item.Edges.Enterprise.ID
//			}
//
//			return res
//		})
//	}
//
//	func (s *ebikeService) Create(req *model.EbikeCreateReq) {
//		b := NewEbikeBrand().QueryX(req.BrandID)
//		s.orm.Create().
//			SetBrandID(b.ID).
//			SetExFactory(req.ExFactory).
//			SetSn(req.SN).
//			SetNillableMachine(req.Machine).
//			SetNillablePlate(req.Plate).
//			SetNillableSim(req.Sim).
//			SetNillableEnable(req.Enable).
//			SetNillableColor(req.Color).
//			ExecX(s.ctx)
//	}
//
//	func (s *ebikeService) Modify(req *model.EbikeModifyReq) {
//		updater := s.QueryX(req.ID).Update()
//
//		if req.ExFactory != nil {
//			updater.SetExFactory(*req.ExFactory)
//		}
//
//		updater.
//			SetNillableMachine(req.Machine).
//			SetNillablePlate(req.Plate).
//			SetNillableSim(req.Sim).
//			SetNillableEnable(req.Enable).
//			SetNillableColor(req.Color).
//			ExecX(s.ctx)
//	}
//
// // BatchCreate 批量创建
// // 0-型号:brand(需查询) 1-车架号:sn 2-生产批次:exFactory 3-车牌号:plate 4-终端编号:machine 5-SIM卡:sim 6-颜色:color
//
//	func (s *ebikeService) BatchCreate(c echo.Context) []string {
//		rows, sns, failed, err := s.BaseService.GetXlsxRows(c, 2, 7, 1)
//		if err != nil {
//			snag.Panic(err)
//		}
//		// 获取所有型号
//		brands := NewEbikeBrand().All()
//		bm := make(map[string]uint64)
//		for _, brand := range brands {
//			bm[brand.Name] = brand.ID
//		}
//
//		arr, _ := s.orm.Query().Where(ebike.SnIn(sns...)).All(s.ctx)
//		exists := make(map[string]bool)
//		for _, a := range arr {
//			exists[a.Sn] = true
//		}
//
//		for _, columns := range rows {
//
//			bid, ok := bm[columns[0]]
//			if !ok {
//				failed = append(failed, fmt.Sprintf("型号未找到:%s", strings.Join(columns, ",")))
//				continue
//			}
//
//			if _, ok = exists[columns[1]]; ok {
//				failed = append(failed, fmt.Sprintf("车架号重复:%s", strings.Join(columns, ",")))
//				continue
//			}
//
//			creator := s.orm.Create().SetBrandID(bid).SetSn(columns[1]).SetExFactory(columns[2]).SetRemark("批量导入")
//			if len(columns) > 3 {
//				creator.SetPlate(columns[3])
//			}
//			if len(columns) > 4 {
//				creator.SetMachine(columns[4])
//			}
//			if len(columns) > 5 {
//				creator.SetSim(columns[5])
//			}
//			color := model.EbikeColorDefault
//			if len(columns) > 6 {
//				color = strings.ReplaceAll(columns[6], "色", "")
//			}
//			creator.SetColor(color)
//
//			err := creator.Exec(s.ctx)
//			if err != nil {
//				msg := "保存失败"
//				if strings.Contains(err.Error(), "duplicate key value") {
//					msg = "有重复项"
//				}
//				failed = append(failed, fmt.Sprintf("%s:%s", msg, strings.Join(columns, ",")))
//			}
//		}
//
//		return failed
//	}
func (s *ebikeService) Detail(bike *ent.Asset, brand *ent.EbikeBrand) *model.Ebike {
	if bike == nil && brand == nil {
		return nil
	}
	res := &model.Ebike{}
	if bike != nil {
		res.EbikeInfo = model.EbikeInfo{
			ID: bike.ID,
			SN: bike.Sn,
		}
		// 查询属性
		ab, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.AssetType(model.AssetTypeEbike.Value())).All(s.ctx)
		// 赋值
		values, _ := bike.QueryValues().All(s.ctx)
		for _, v := range ab {
			for _, av := range values {
				if v.ID == av.AttributeID {
					switch v.Key {
					case "plate":
						res.Plate = silk.String(av.Value)
					case "color":
						res.Color = av.Value
					case "exFactory":
						res.ExFactory = av.Value
					}
				}
			}
		}
	}
	if brand != nil {
		res.Brand = &model.EbikeBrand{
			ID:   brand.ID,
			Name: brand.Name,
		}
	}
	return res
}

func (s *ebikeService) UnallocatedX(params *model.EbikeUnallocatedParams) *model.Ebike {
	bikes := s.SearchUnallocated(params)
	if len(bikes) == 0 {
		snag.Panic("未找到有效车辆")
	}
	if len(bikes) > 1 {
		snag.Panic("存在多个车辆, 请缩小查询范围")
	}
	return bikes[0]
}

// SearchUnallocated 获取未分配车辆信息
func (s *ebikeService) SearchUnallocated(params *model.EbikeUnallocatedParams) (res []*model.Ebike) {
	q := s.AllocatableBaseFilter().WithEbikeAllocates(func(query *ent.AllocateQuery) {
		query.Where(
			allocate.Status(model.AllocateStatusPending.Value()),
			allocate.TimeGTE(carbon.Now().SubSeconds(model.AllocateExpiration).StdTime()),
		)
	})

	if (params.ID == nil && params.Keyword == nil) || (params.ID != nil && params.Keyword != nil) {
		snag.Panic("参数错误")
	}

	// 根据ID或关键字查找车辆
	if params.ID != nil {
		q.Where(asset.ID(*params.ID))
	}

	if params.Keyword != nil {
		q.Where(
			asset.Or(
				asset.SnContainsFold(*params.Keyword),
			))
		attributes, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.Key("plate")).First(s.ctx)
		if attributes != nil {
			q.Where(
				asset.Or(
					asset.HasValuesWith(
						assetattributevalues.AttributeID(attributes.ID),
						assetattributevalues.ValueContainsFold(*params.Keyword),
					),
				),
			)
		}
	}

	// 站点
	if params.StationID != nil {
		q.Where(asset.LocationsType(model.AssetLocationsTypeStation.Value()), asset.LocationsID(*params.StationID))
	}

	// 门店
	if params.StoreID != nil {
		q.Where(asset.LocationsType(model.AssetLocationsTypeStore.Value()), asset.LocationsID(*params.StoreID))
	}

	// 查找电车
	bikes, _ := q.WithBrand().WithValues().All(s.ctx)

	res = make([]*model.Ebike, len(bikes))
	for i, bike := range bikes {
		if len(bike.Edges.EbikeAllocates) > 0 {
			continue
		}
		brand := bike.Edges.Brand
		res[i] = &model.Ebike{
			EbikeInfo: model.EbikeInfo{
				ID: bike.ID,
				SN: bike.Sn,
			},
			Brand: &model.EbikeBrand{
				ID:    brand.ID,
				Name:  brand.Name,
				Cover: brand.Cover,
			},
		}
		// 查询电车属性
		ab, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.AssetType(model.AssetTypeEbike.Value())).All(s.ctx)
		if ab != nil {
			for _, v := range bike.Edges.Values {
				for _, av := range ab {
					if v.AttributeID == av.ID {
						switch av.Key {
						case "plate":
							res[i].Plate = silk.String(v.Value)
						case "color":
							res[i].Color = v.Value
						case "exFactory":
							res[i].ExFactory = v.Value
						}
					}
				}
			}
		}
	}
	return
}

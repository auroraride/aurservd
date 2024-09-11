package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/internal/ent/assetcheck"
	"github.com/auroraride/aurservd/internal/ent/assetcheckdetails"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assetCheckService struct {
	orm *ent.AssetCheckClient
}

func NewAssetCheck() *assetCheckService {
	return &assetCheckService{
		orm: ent.Database.AssetCheck,
	}
}

// CreateAssetCheck 创建资产盘点
func (s *assetCheckService) CreateAssetCheck(ctx context.Context, req *model.AssetCheckCreateReq, modifier *model.Modifier) (cID uint64, err error) {
	// 应盘资产
	var assetIDs []uint64
	// 查询应盘电车资产
	allEbike, err := s.getCheckAsset(ctx, model.GetCheckAssetReq{
		OperatorType:  req.OperatorType,
		OperatorID:    req.OperatorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
		AssetType:     model.AssetTypeEbike,
	})
	if err != nil {
		return
	}
	bulk := make([]*ent.AssetCheckDetailsCreate, 0)
	for _, v := range allEbike {
		bulk = append(bulk, ent.Database.AssetCheckDetails.Create().
			SetAsset(v).
			SetCreator(modifier).
			SetLastModifier(modifier).
			SetStatus(model.AssetCheckDetailsStatusUntreated.Value()).
			SetLocationsID(v.LocationsID).
			SetLocationsType(v.LocationsType))
		assetIDs = append(assetIDs, v.ID)
	}
	// 查询应盘电池资产
	allBattery, err := s.getCheckAsset(ctx, model.GetCheckAssetReq{
		OperatorType:  req.OperatorType,
		OperatorID:    req.OperatorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
		AssetType:     model.AssetTypeSmartBattery,
	})
	if err != nil {
		return
	}
	for _, v := range allBattery {
		bulk = append(bulk, ent.Database.AssetCheckDetails.Create().
			SetAsset(v).
			SetCreator(modifier).
			SetLastModifier(modifier).
			SetStatus(model.AssetCheckDetailsStatusUntreated.Value()).
			SetLocationsID(v.LocationsID).
			SetLocationsType(v.LocationsType))
		assetIDs = append(assetIDs, v.ID)
	}
	var realEbikeNum, realBatteryNum uint
	var realAssetIDs []uint64
	for _, v := range req.AssetCheckCreateDetail {
		// 限制只能盘点自己的资产
		err = s.checkAssetCheckOwner(ctx, model.CheckAssetCheckOwnerReq{
			AssetID:       v.AssetID,
			AssetType:     v.AssetType,
			OperatorType:  req.OperatorType,
			OperatorID:    req.OperatorID,
			LocationsType: req.LocationsType,
			LocationsID:   req.LocationsID,
		})
		if err != nil {
			return
		}
		// 获取实际盘点数量
		if v.AssetType == model.AssetTypeEbike {
			realEbikeNum++
		} else {
			realBatteryNum++
		}
		realAssetIDs = append(realAssetIDs, v.AssetID)
	}
	b, err := ent.Database.AssetCheckDetails.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return
	}

	start := tools.NewTime().ParseDateTimeStringX(req.StartAt)
	end := tools.NewTime().ParseDateTimeStringX(req.EndAt)

	c, err := s.orm.Create().
		SetCreator(modifier).
		SetLastModifier(modifier).
		SetOperateID(req.OperatorID).
		SetOperateType(req.OperatorType.Value()).
		SetLocationsID(req.LocationsID).
		SetLocationsType(req.LocationsType.Value()).
		SetEbikeNum(uint(len(allEbike))).
		SetBatteryNum(uint(len(allBattery))).
		SetBatteryNumReal(realBatteryNum).
		SetEbikeNumReal(realEbikeNum).
		SetStartAt(start).
		SetEndAt(end).
		AddCheckDetails(b...).
		Save(ctx)
	if err != nil {
		return
	}
	// 多出来的资产ID 和 未盘点的资产ID
	missingIDs, extraIDs := s.findMissingAndExtraIDs(assetIDs, realAssetIDs)
	_, err = ent.Database.AssetCheckDetails.Update().
		Where(
			assetcheckdetails.AssetIDIn(missingIDs...),
			assetcheckdetails.CheckID(c.ID),
		).
		SetResult(model.AssetCheckResultLoss.Value()).
		SetRealLocationsID(req.LocationsID).
		SetRealLocationsType(req.LocationsType.Value()).
		Save(ctx)
	if err != nil {
		return
	}
	_, err = ent.Database.AssetCheckDetails.Update().
		Where(
			assetcheckdetails.AssetIDIn(extraIDs...),
			assetcheckdetails.CheckID(c.ID),
		).
		SetResult(model.AssetCheckResultSurplus.Value()).
		Save(ctx)
	if err != nil {
		return
	}
	// 其余为正常
	_, err = ent.Database.AssetCheckDetails.Update().
		Where(
			assetcheckdetails.AssetIDNotIn(append(missingIDs, extraIDs...)...),
			assetcheckdetails.CheckID(c.ID),
		).
		SetResult(model.AssetCheckResultNormal.Value()).
		Save(ctx)
	if err != nil {
		return
	}
	return c.ID, nil
}

func (s *assetCheckService) findMissingAndExtraIDs(assetIDs []uint64, realAssetIDs []uint64) ([]uint64, []uint64) {
	// 用 map 记录实际盘点资产的 ID
	realAssetMap := make(map[uint64]bool)
	for _, assetID := range realAssetIDs {
		realAssetMap[assetID] = true
	}

	// 分别记录未被盘点和多出来的资产ID
	missingIDs := make([]uint64, 0)
	extraIDs := make([]uint64, 0)

	// 找出未被盘点的资产ID
	for _, assetID := range assetIDs {
		if _, ok := realAssetMap[assetID]; !ok {
			missingIDs = append(missingIDs, assetID)
		} else {
			delete(realAssetMap, assetID) // 从 map 中删除已经盘点的资产ID
		}
	}

	// 剩下的就是多出来的资产ID
	for assetID := range realAssetMap {
		extraIDs = append(extraIDs, assetID)
	}

	return missingIDs, extraIDs
}

// 查询应盘点资产
func (s *assetCheckService) getCheckAsset(ctx context.Context, req model.GetCheckAssetReq) (res []*ent.Asset, err error) {
	t, _, err := s.getAssetByOperateRole(ctx, model.GetAssetByOperateRole{
		OperatorType:  req.OperatorType,
		OperatorID:    req.OperatorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
	})
	if err != nil {
		return nil, err
	}
	if !t {
		return nil, fmt.Errorf("当前用户无操作权限")
	}
	items, _ := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.LocationsType(req.LocationsType.Value()),
			asset.Type(req.AssetType.Value()),
			asset.LocationsID(req.LocationsID),
		).All(ctx)
	if items == nil {
		return nil, fmt.Errorf("未找到对应的资产")
	}
	return items, nil
}

// 根据操作角色找对应能操作的资产
func (s *assetCheckService) getAssetByOperateRole(ctx context.Context, req model.GetAssetByOperateRole) (b bool, ids []uint64, err error) {
	ids = make([]uint64, 0)
	b = false

	switch req.OperatorType {
	case model.OperatorTypeAgent:
		item, _ := ent.Database.Agent.QueryNotDeleted().Where(agent.ID(req.OperatorID)).
			WithEnterprise(func(query *ent.EnterpriseQuery) {
				query.WithStations()
			}).First(ctx)
		if item == nil {
			return b, nil, fmt.Errorf("未找到对应的代理商")
		}

		if item.Edges.Enterprise == nil {
			return b, nil, fmt.Errorf("未找到对应的团签企业")
		}

		if len(item.Edges.Enterprise.Edges.Stations) == 0 {
			return b, nil, fmt.Errorf("代理商未绑定站点")
		}
		for _, v := range item.Edges.Enterprise.Edges.Stations {
			if v.ID == req.LocationsID {
				b = true
			}
			ids = append(ids, v.ID)
		}
	case model.OperatorTypeEmployee:
		// 查询是否当前位置上班
		st, _ := ent.Database.Store.QueryNotDeleted().Where(
			store.ID(req.LocationsID),
			store.HasDutyEmployeesWith(employee.ID(req.OperatorID)),
		).First(context.Background())
		if st == nil {
			return b, nil, errors.New("当年门店未存在上班信息")
		}
		b = true
	case model.OperatorTypeAssetManager:
		// 查询是否当前位置上班
		wh, _ := ent.Database.Warehouse.QueryNotDeleted().Where(
			warehouse.ID(req.LocationsID),
			warehouse.HasDutyAssetManagersWith(assetmanager.ID(req.OperatorID)),
		).First(context.Background())
		if wh == nil {
			return b, nil, errors.New("当年门店未存在上班信息")
		}
		b = true
	default:
		return b, nil, fmt.Errorf("未知的操作人类型")
	}
	return b, ids, nil
}

// 验证资产盘点是否属于当前用户
func (s *assetCheckService) checkAssetCheckOwner(ctx context.Context, req model.CheckAssetCheckOwnerReq) error {
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.ID(req.AssetID),
			asset.Status(model.AssetStatusStock.Value()),
			asset.Type(req.AssetType.Value()),
			asset.LocationsType(req.LocationsType.Value()),
			asset.LocationsID(req.LocationsID),
		)
	item, _ := q.First(ctx)
	if item == nil {
		return fmt.Errorf("该资产不属于当前目标位置")
	}
	return nil
}

// GetAssetBySN 通过sn查询资产
func (s *assetCheckService) GetAssetBySN(ctx context.Context, req *model.AssetCheckByAssetSnReq) (res *model.AssetCheckByAssetSnRes, err error) {
	item, _ := ent.Database.Asset.QueryNotDeleted().Where(asset.Sn(req.SN)).WithModel().WithBrand().WithValues().First(ctx)
	if item == nil {
		return nil, fmt.Errorf("未找到对应的资产")
	}
	// 验证资产是否属于当前用户
	err = s.checkAssetCheckOwner(ctx, model.CheckAssetCheckOwnerReq{
		AssetID:       item.ID,
		AssetType:     model.AssetType(item.Type),
		OperatorType:  req.OperatorType,
		OperatorID:    req.OperatorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
	})
	if err != nil {
		return nil, err
	}
	res = &model.AssetCheckByAssetSnRes{
		AssetID:       item.ID,
		AssetSN:       item.Sn,
		AssetType:     model.AssetType(item.Type),
		LocationsType: model.AssetLocationsType(item.LocationsType),
		LocationsID:   item.LocationsID,
	}
	if item.Edges.Model != nil {
		res.Model = item.Edges.Model.Model
	}
	if item.Edges.Brand != nil {
		res.BrandName = item.Edges.Brand.Name
	}

	// 查询属性值
	attributeValue, _ := item.QueryValues().WithAttribute().All(context.Background())
	assetAttributeMap := make(map[uint64]model.AssetAttribute)
	for _, v := range attributeValue {
		var attributeName, attributeKey string
		if v.Edges.Attribute != nil {
			attributeName = v.Edges.Attribute.Name
			attributeKey = v.Edges.Attribute.Key
		}
		assetAttributeMap[v.AttributeID] = model.AssetAttribute{
			AttributeID:      v.AttributeID,
			AttributeValue:   v.Value,
			AttributeName:    attributeName,
			AttributeKey:     attributeKey,
			AttributeValueID: v.ID,
		}
	}

	res.Attribute = assetAttributeMap

	return res, nil
}

// MarkStartOrEndCheck 标记开始盘点或结束盘点
func (s *assetCheckService) MarkStartOrEndCheck(ctx context.Context, req *model.MarkStartOrEndCheckReq) (err error) {
	t, _, err := s.getAssetByOperateRole(ctx, model.GetAssetByOperateRole{
		OperatorType:  req.OperatorType,
		OperatorID:    req.OperatorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
	})
	if err != nil {
		return err
	}
	if !t {
		return fmt.Errorf("当前用户无操作权限")
	}
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.LocationsType(req.LocationsType.Value()),
			asset.LocationsID(req.LocationsID),

			asset.Status(model.AssetStatusStock.Value()),
		)
	if req.Enable {
		// 标记开始盘点
		q.Where(asset.CheckAtIsNil())
	} else {
		// 标记结束盘点
		q.Where(asset.CheckAtNotNil())
	}
	all, _ := q.All(ctx)
	if all == nil {
		return fmt.Errorf("未找到对应的资产")
	}
	assetIDs := make([]uint64, 0)
	for _, v := range all {
		assetIDs = append(assetIDs, v.ID)
	}
	update := ent.Database.Asset.Update().Where(asset.IDIn(assetIDs...))
	if req.Enable {
		update.SetCheckAt(time.Now())
	} else {
		update.SetNillableCheckAt(nil)
	}
	_, err = update.Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

// List 盘点列表
func (s *assetCheckService) List(ctx context.Context, req *model.AssetCheckListReq) (*model.PaginationRes, error) {
	q := ent.Database.AssetCheck.QueryNotDeleted().WithStore().WithStation().WithWarehouse().
		WithOperateStore().WithOperateManager().WithOperateAgent().WithCheckDetails(func(query *ent.AssetCheckDetailsQuery) {
		query.WithAsset()
	})

	s.listFilter(q, req)

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetCheck) *model.AssetCheckListRes {
		return s.detail(item)
	}), nil
}

func (s *assetCheckService) listFilter(q *ent.AssetCheckQuery, req *model.AssetCheckListReq) {
	if req.LocationsID != nil && req.LocationsType != nil {
		q.Where(
			assetcheck.LocationsID(*req.LocationsID),
			assetcheck.LocationsType(req.LocationsType.Value()),
		)
	}

	if req.Keyword != nil {
		q.Where(
			assetcheck.Or(
				assetcheck.HasWarehouseWith(warehouse.NameContains(*req.Keyword)),
				assetcheck.HasStationWith(enterprisestation.HasAgentsWith(agent.NameContains(*req.Keyword))),
				assetcheck.HasStoreWith(store.Name(*req.Keyword)),
			),
		)
	}
	if req.CheckResult != nil {
		if *req.CheckResult {
			q.Where(
				func(selector *sql.Selector) {
					selector.Where(
						sql.And(
							sql.ColumnsEQ(assetcheck.FieldBatteryNum, assetcheck.FieldBatteryNumReal),
							sql.ColumnsEQ(assetcheck.FieldEbikeNum, assetcheck.FieldEbikeNumReal),
						),
					)
				},
			)
		} else {
			q.Where(
				func(selector *sql.Selector) {
					selector.Where(
						sql.Or(
							sql.ColumnsNEQ(assetcheck.FieldBatteryNum, assetcheck.FieldBatteryNumReal),
							sql.ColumnsNEQ(assetcheck.FieldEbikeNum, assetcheck.FieldEbikeNumReal),
						),
					)
				},
			)
		}
	}
	if req.StartAt != nil && req.EndAt != nil {
		start := tools.NewTime().ParseDateStringX(*req.StartAt)
		end := tools.NewTime().ParseNextDateStringX(*req.EndAt)
		q.Where(assetcheck.EndAtGTE(start), assetcheck.EndAtLTE(end))
	}

	if len(req.LocationsIds) != 0 && req.LocationsType != nil {
		q.Where(
			assetcheck.LocationsIDIn(req.LocationsIds...),
			assetcheck.LocationsType(req.LocationsType.Value()),
		)
	}
}

// Detail 盘点明细
func (s *assetCheckService) Detail(ctx context.Context, id uint64) (*model.AssetCheckListRes, error) {
	item, _ := ent.Database.AssetCheck.QueryNotDeleted().Where(assetcheck.ID(id)).WithStore().WithStation().WithWarehouse().
		WithOperateStore().WithOperateManager().WithOperateAgent().WithCheckDetails(func(query *ent.AssetCheckDetailsQuery) {
		query.WithAsset(func(query *ent.AssetQuery) {
			query.WithModel().WithBrand()
		})
	}).First(ctx)
	if item == nil {
		return nil, errors.New("未找到对应的盘点记录")
	}
	return s.detail(item), nil
}

func (s *assetCheckService) detail(item *ent.AssetCheck) *model.AssetCheckListRes {
	var start, end string
	if item.StartAt != nil {
		start = item.StartAt.Format("2006-01-02 15:04:05")
	}
	if item.EndAt != nil {
		end = item.EndAt.Format("2006-01-02 15:04:05")
	}
	OperatorName := ""
	switch model.OperatorType(item.OperateType) {
	case model.OperatorTypeAgent:
		if item.Edges.OperateAgent != nil {
			OperatorName = item.Edges.OperateAgent.Name
		}
	case model.OperatorTypeEmployee:
		if item.Edges.OperateStore != nil {
			OperatorName = item.Edges.OperateStore.Name
		}
	case model.OperatorTypeAssetManager:
		if item.Edges.OperateManager != nil {
			OperatorName = item.Edges.OperateManager.Name
		}
	default:
		OperatorName = ""
	}
	locationsName := ""
	switch model.AssetLocationsType(item.LocationsType) {
	case model.AssetLocationsTypeStore:
		locationsName = "[平台门店]-" + item.Edges.Store.Name
	case model.AssetLocationsTypeStation:
		locationsName = "[代理站点]-" + item.Edges.Station.Name
	case model.AssetLocationsTypeWarehouse:
		locationsName = "[平台仓库]-" + item.Edges.Warehouse.Name
	default:
		locationsName = ""
	}

	var checkResult bool
	if item.BatteryNum == item.BatteryNumReal && item.EbikeNum == item.EbikeNumReal {
		checkResult = true
	} else {
		checkResult = false
	}

	var status string
	d, _ := item.QueryCheckDetails().Where(assetcheckdetails.ResultNEQ(model.AssetCheckResultNormal.Value())).All(context.Background())
	if len(d) > 0 {
		status = model.AssetCheckStatusPending.String()
	} else {
		status = model.AssetCheckStatusProcessed.String()
	}

	abs, _ := s.ListAbnormal(context.Background(), &model.AssetCheckListAbnormalReq{ID: item.ID})

	res := &model.AssetCheckListRes{
		ID:             item.ID,
		StartAt:        start,
		EndAt:          end,
		OperatorID:     item.OperateID,
		OperatorName:   OperatorName,
		BatteryNum:     item.BatteryNum,
		BatteryNumReal: item.BatteryNumReal,
		EbikeNum:       item.EbikeNum,
		EbikeNumReal:   item.EbikeNumReal,
		LocationsID:    item.LocationsID,
		LocationsType:  item.LocationsType,
		LocationsName:  locationsName,
		CheckResult:    checkResult,
		Status:         status,
		Abnormals:      abs,
	}
	return res
}

// ListAbnormal 查询盘点异常资产
func (s *assetCheckService) ListAbnormal(ctx context.Context, req *model.AssetCheckListAbnormalReq) (res []*model.AssetCheckAbnormal, err error) {
	// 查询异常资产
	abnormalAll, _ := ent.Database.AssetCheckDetails.QueryNotDeleted().
		Where(
			assetcheckdetails.ResultNEQ(model.AssetCheckResultNormal.Value()),
			assetcheckdetails.CheckID(req.ID),
		).WithAsset(func(query *ent.AssetQuery) {
		query.WithModel().WithBrand().WithMaterial()
	}).WithWarehouse().WithStore().WithStation().WithRider().WithCabinet().WithOperator().
		WithRealWarehouse().WithRealStation().WithRealStore().WithRealCabinet().WithRealRider().WithRealOperator().
		All(context.Background())

	for _, v := range abnormalAll {
		var name, sn, modelName, brandName, locationsName, realLocationsName, operatorName string
		if v.Edges.Asset != nil {
			name = v.Edges.Asset.Name
			sn = v.Edges.Asset.Sn
			if v.Edges.Asset.Edges.Model != nil {
				modelName = v.Edges.Asset.Edges.Model.Model
			}
			if v.Edges.Asset.Edges.Brand != nil {
				brandName = v.Edges.Asset.Edges.Brand.Name
			}
		}

		switch model.AssetLocationsType(v.LocationsType) {
		case model.AssetLocationsTypeStore:
			if v.Edges.Store != nil {
				locationsName = "[平台门店]-" + v.Edges.Store.Name
			}
		case model.AssetLocationsTypeStation:
			if v.Edges.Station != nil {
				locationsName = "[代理站点]-" + v.Edges.Station.Name
			}
		case model.AssetLocationsTypeWarehouse:
			if v.Edges.Warehouse != nil {
				locationsName = "[平台仓库]-" + v.Edges.Warehouse.Name
			}
		case model.AssetLocationsTypeRider:
			if v.Edges.Rider != nil {
				locationsName = "[骑手]-" + v.Edges.Rider.Name
			}
		case model.AssetLocationsTypeCabinet:
			if v.Edges.Cabinet != nil {
				locationsName = "[电柜]-" + v.Edges.Cabinet.Name
			}
		case model.AssetLocationsTypeOperation:
			if v.Edges.Operator != nil {
				locationsName = "[运维]-" + v.Edges.Operator.Name
			}
		default:
			locationsName = ""
		}

		switch model.AssetLocationsType(v.RealLocationsType) {
		case model.AssetLocationsTypeStore:
			if v.Edges.RealStore != nil {
				realLocationsName = "[平台门店]-" + v.Edges.RealStore.Name
			}
		case model.AssetLocationsTypeStation:
			if v.Edges.RealStation != nil {
				realLocationsName = "[代理站点]-" + v.Edges.RealStation.Name
			}
		case model.AssetLocationsTypeWarehouse:
			if v.Edges.RealWarehouse != nil {
				realLocationsName = "[平台仓库]-" + v.Edges.RealWarehouse.Name
			}
		case model.AssetLocationsTypeRider:
			if v.Edges.RealRider != nil {
				realLocationsName = "[骑手]-" + v.Edges.RealRider.Name
			}
		case model.AssetLocationsTypeCabinet:
			if v.Edges.RealCabinet != nil {
				realLocationsName = "[电柜]-" + v.Edges.RealCabinet.Name
			}
		case model.AssetLocationsTypeOperation:
			if v.Edges.RealOperator != nil {
				realLocationsName = "[运维]-" + v.Edges.RealOperator.Name
			}
		default:
			realLocationsName = ""
		}

		if v.Edges.Maintainer != nil {
			operatorName = v.Edges.Maintainer.Name
		}
		var operatorAt string
		if v.OperateAt != nil {
			operatorAt = v.OperateAt.Format("2006-01-02 15:04:05")
		}

		var assetType model.AssetType
		if v.Edges.Asset != nil {
			assetType = model.AssetType(v.Edges.Asset.Type)
		}

		res = append(res, &model.AssetCheckAbnormal{
			AssetID:           v.AssetID,
			Name:              name,
			Model:             modelName,
			Brand:             brandName,
			Result:            model.AssetCheckResult(v.Result),
			SN:                sn,
			LocationsName:     locationsName,
			RealLocationsName: realLocationsName,
			Status:            model.AssetCheckDetailsStatus(v.Status).String(),
			OperatorName:      operatorName,
			OperatorAt:        operatorAt,
			AssetType:         assetType,
		})
	}
	return res, nil
}

// AssetDetailList 盘点资产明细列表
func (s *assetCheckService) AssetDetailList(ctx context.Context, req *model.AssetCheckDetailListReq) (*model.PaginationRes, error) {
	q := ent.Database.AssetCheckDetails.QueryNotDeleted().
		Where(
			assetcheckdetails.CheckID(req.ID),
			assetcheckdetails.HasAssetWith(asset.Type(req.AssetType.Value())),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithModel().WithBrand().WithMaterial()
		}).WithWarehouse().WithStore().WithStation().WithRider().WithCabinet().WithOperator().
		WithRealWarehouse().WithRealStation().WithRealStore().WithRealCabinet().WithRealRider().WithRealOperator()

	// 实际盘点资产
	if req.RealCheck {
		q.Where(assetcheckdetails.ResultNEQ(model.AssetCheckResultUntreated.Value()))
	}
	// 属性查询
	if req.Attribute != nil {
		var attributeID uint64
		var attributeValue string
		// 解析 attribute "id:value,id:value" 格式
		for _, v := range strings.Split(*req.Attribute, ",") {
			av := strings.Split(v, ":")
			if len(av) != 2 {
				continue
			}
			attributeID, _ = strconv.ParseUint(av[0], 10, 64)
			attributeValue = av[1]
			q.Where(assetcheckdetails.HasAssetWith(asset.HasValuesWith(assetattributevalues.AttributeID(attributeID), assetattributevalues.ValueContains(attributeValue))))
		}
	}
	if req.SN != nil {
		q.Where(assetcheckdetails.HasAssetWith(asset.Sn(*req.SN)))
	}
	if req.ModelID != nil {
		q.Where(assetcheckdetails.HasAssetWith(asset.ModelID(*req.ModelID)))
	}
	if req.BrandID != nil {
		q.Where(assetcheckdetails.HasAssetWith(asset.BrandID(*req.BrandID)))
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetCheckDetails) (res *model.AssetCheckDetail) {
		res = &model.AssetCheckDetail{
			ID:      item.ID,
			AssetID: item.AssetID,
		}
		if item.Edges.Asset != nil {
			res.AssetSN = item.Edges.Asset.Sn
			res.AssetStatus = item.Edges.Asset.Status
			res.AssetType = item.Edges.Asset.Type
			attributeValue, _ := item.Edges.Asset.QueryValues().WithAttribute().All(ctx)
			assetAttributeMap := make(map[uint64]model.AssetAttribute)
			for _, v := range attributeValue {
				var attributeName, attributeKey string
				if v.Edges.Attribute != nil {
					attributeName = v.Edges.Attribute.Name
					attributeKey = v.Edges.Attribute.Key
				}
				assetAttributeMap[v.AttributeID] = model.AssetAttribute{
					AttributeID:      v.AttributeID,
					AttributeValue:   v.Value,
					AttributeName:    attributeName,
					AttributeKey:     attributeKey,
					AttributeValueID: v.ID,
				}
			}
			res.Attribute = assetAttributeMap
		}
		if item.Edges.Asset.Edges.Model != nil {
			res.Model = item.Edges.Asset.Edges.Model.Model
		}
		if item.Edges.Asset.Edges.Brand != nil {
			res.BrandName = item.Edges.Asset.Edges.Brand.Name
		}

		// 返回实际位置
		switch model.AssetLocationsType(item.LocationsType) {
		case model.AssetLocationsTypeStore:
			if item.Edges.Store != nil {
				res.LocationsName = "[平台门店]-" + item.Edges.Store.Name
			}
		case model.AssetLocationsTypeStation:
			if item.Edges.Station != nil {
				res.LocationsName = "[代理站点]-" + item.Edges.Station.Name
			}
		case model.AssetLocationsTypeWarehouse:
			if item.Edges.Warehouse != nil {
				res.LocationsName = "[平台仓库]-" + item.Edges.Warehouse.Name
			}
		case model.AssetLocationsTypeRider:
			if item.Edges.Rider != nil {
				res.LocationsName = "[骑手]-" + item.Edges.Rider.Name
			}
		case model.AssetLocationsTypeCabinet:
			if item.Edges.Cabinet != nil {
				res.LocationsName = "[电柜]-" + item.Edges.Cabinet.Name
			}
		case model.AssetLocationsTypeOperation:
			if item.Edges.Operator != nil {
				res.LocationsName = "[运维]-" + item.Edges.Operator.Name
			}
		default:
			res.LocationsName = ""
		}

		// 返回实际位置
		switch model.AssetLocationsType(item.RealLocationsType) {
		case model.AssetLocationsTypeStore:
			if item.Edges.RealStore != nil {
				res.RealLocationsName = "[平台门店]-" + item.Edges.RealStore.Name
			}
		case model.AssetLocationsTypeStation:
			if item.Edges.RealStation != nil {
				res.RealLocationsName = "[代理站点]-" + item.Edges.RealStation.Name
			}
		case model.AssetLocationsTypeWarehouse:
			if item.Edges.RealWarehouse != nil {
				res.RealLocationsName = "[平台仓库]-" + item.Edges.RealWarehouse.Name
			}
		case model.AssetLocationsTypeRider:
			if item.Edges.RealRider != nil {
				res.RealLocationsName = "[骑手]-" + item.Edges.RealRider.Name
			}
		case model.AssetLocationsTypeCabinet:
			if item.Edges.RealCabinet != nil {
				res.RealLocationsName = "[电柜]-" + item.Edges.RealCabinet.Name
			}
		case model.AssetLocationsTypeOperation:
			if item.Edges.RealOperator != nil {
				res.RealLocationsName = "[运维]-" + item.Edges.RealOperator.Name
			}
		default:
			res.RealLocationsName = ""
		}
		return res
	}), nil
}

// AssetCheckAbnormalOperate 盘点异常资产操作
func (s *assetCheckService) AssetCheckAbnormalOperate(ctx context.Context, req *model.AssetCheckAbnormalOperateReq, modifier *model.Modifier) error {
	// 查询异常资产
	item, _ := ent.Database.AssetCheckDetails.QueryNotDeleted().Where(assetcheckdetails.ID(req.ID)).WithAsset().First(ctx)
	if item == nil {
		return fmt.Errorf("未找到对应的异常资产")
	}

	var status model.AssetCheckDetailsStatus

	// 判定操作类型
	if item.Result == model.AssetCheckResultSurplus.Value() {
		// 盘盈操作入库
		fromLocationType := model.AssetLocationsType(item.LocationsType)
		transferCreateReq := &model.AssetTransferCreateReq{
			FromLocationType:  &fromLocationType,
			FromLocationID:    &item.LocationsID,
			ToLocationType:    model.AssetLocationsType(item.RealLocationsType),
			ToLocationID:      item.RealLocationsID,
			Details:           make([]model.AssetTransferCreateDetail, 0),
			Reason:            "盘盈入库",
			AssetTransferType: model.AssetTransferTypeTransfer,
		}
		if item.Edges.Asset == nil {
			return fmt.Errorf("未找到对应的资产")
		}
		transferCreateReq.Details = append(transferCreateReq.Details, model.AssetTransferCreateDetail{
			AssetType: model.AssetType(item.Edges.Asset.Type),
			SN:        &item.Edges.Asset.Sn,
		})

		_, failed, err := NewAssetTransfer().Transfer(ctx, transferCreateReq, modifier)
		if err != nil {
			return err
		}
		if len(failed) > 0 {
			return errors.New(failed[0])
		}

		// 修正已经盘点过的数据
		_, err = ent.Database.AssetCheckDetails.Update().
			Where(
				// 盘亏的资产 并未处理的
				assetcheckdetails.Result(model.AssetCheckResultLoss.Value()),
				assetcheckdetails.Status(model.AssetCheckDetailsStatusUntreated.Value()),
				assetcheckdetails.AssetID(item.AssetID),
			).
			SetResult(model.AssetCheckResultNormal.Value()).
			SetLastModifier(modifier).
			SetOperateAt(time.Now()).
			SetOperateID(modifier.ID).
			SetStatus(model.AssetCheckDetailsStatusOut.Value()).
			Save(ctx)
		if err != nil {
			return err
		}
		// 处理类型
		status = model.AssetCheckDetailsStatusIn
	}
	// 盘亏操作报废
	if item.Result == model.AssetCheckResultLoss.Value() {
		err := NewAssetScrap().Scrap(ctx, &model.AssetScrapReq{
			ScrapReasonType: model.ScrapReasonLost,
			Remark:          silk.String("盘亏报废"),
			Details: []model.AssetScrapDetails{
				{
					AssetType: model.AssetType(item.Edges.Asset.Type),
					Sn:        &item.Edges.Asset.Sn,
				},
			},
		}, modifier)
		if err != nil {
			return err
		}
		// 处理类型
		status = model.AssetCheckDetailsStatusScrap
	}

	// 更新盘点详情
	_, err := ent.Database.AssetCheckDetails.Update().
		Where(assetcheckdetails.ID(req.ID)).
		SetStatus(status.Value()).
		SetLastModifier(modifier).
		SetOperateAt(time.Now()).
		SetOperateID(modifier.ID).
		SetResult(model.AssetCheckResultNormal.Value()).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

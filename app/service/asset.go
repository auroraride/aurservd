package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/auroraride/adapter"
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetattributes"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assetService struct {
	orm *ent.AssetClient
	*BaseService
}

func NewAsset(params ...interface{}) *assetService {
	return &assetService{
		orm:         ent.Database.Asset,
		BaseService: newService(params...),
	}
}

// Create 创建物资
func (s *assetService) Create(ctx context.Context, req *model.AssetCreateReq, modifier *model.Modifier) error {
	// 去重复
	if req.SN != nil {
		if b, _ := s.orm.QueryNotDeleted().Where(asset.Sn(*req.SN), asset.Type(req.AssetType.Value())).Exist(ctx); b {
			return errors.New("编号重复")
		}
	}
	// 默认启用
	enable := true
	if req.Enable != nil {
		enable = *req.Enable
	}
	// 入库状态默认为配送中
	assetStatus := model.AssetStatusDelivering.Value()

	var name string

	q := s.orm.Create()
	switch req.AssetType {
	case model.AssetTypeSmartBattery:
		// 解析电池编号
		if req.SN == nil {
			return errors.New("智能电池编号不能为空")
		}
		if req.CityID == nil {
			return errors.New("城市不能为空")
		}
		ab, err := adapter.ParseBatterySN(*req.SN)
		if err != nil {
			return errors.New("电池编号解析失败" + err.Error())
		}
		// 查询型号是否存在
		if ab.Model == "" {
			return fmt.Errorf("电池编号%s解析失败", *req.SN)
		}
		modelInfo, _ := ent.Database.BatteryModel.Query().Where(batterymodel.Model(ab.Model)).Only(ctx)
		if modelInfo == nil {
			return fmt.Errorf("电池型号%s不存在", ab.Model)
		}
		name = s.getAssetName(req.AssetType, modelInfo.ID)
		q.SetNillableModelID(&modelInfo.ID).
			SetBrandName(ab.Brand.String()).
			SetNillableCityID(req.CityID)
	case model.AssetTypeEbike:
		if req.BrandID == nil {
			return errors.New("品牌不能为空")
		}
		q.SetBrandID(*req.BrandID)
		name = s.getAssetName(req.AssetType, *req.BrandID)
	default:
		return errors.New("未知类型")
	}

	q.SetType(req.AssetType.Value()).
		SetName(name).
		SetNillableSn(req.SN).
		SetEnable(enable).
		SetStatus(assetStatus).
		SetCreator(modifier).
		SetLastModifier(modifier).
		SetLocationsType((req.LocationsType).Value()).
		SetLocationsID(req.LocationsID)
	item, err := q.Save(ctx)
	if err != nil {
		return err
	}

	bulk := make([]*ent.AssetAttributeValuesCreate, 0, len(req.Attribute))
	for _, v := range req.Attribute {
		// 判定属性值是否存在
		if b, _ := ent.Database.AssetAttributeValues.Query().Where(assetattributevalues.AttributeID(v.AttributeID), assetattributevalues.AssetID(item.ID)).Exist(ctx); b {
			return errors.New("属性值重复")
		}
		bulk = append(bulk, ent.Database.AssetAttributeValues.
			Create().
			SetValue(v.AttributeValue).
			SetAttributeID(v.AttributeID).
			SetAssetID(item.ID))
	}
	_, err = ent.Database.AssetAttributeValues.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return err
	}
	// 创建调拨单
	_, failed, err := NewAssetTransfer().Transfer(ctx, &model.AssetTransferCreateReq{
		ToLocationType:    req.LocationsType,
		ToLocationID:      req.LocationsID,
		Reason:            "初始入库",
		AssetTransferType: model.AssetTransferTypeInitial,
		Details: []model.AssetTransferCreateDetail{
			{
				AssetType: req.AssetType,
				SN:        req.SN,
			},
		},
		OperatorID:   modifier.ID,
		OperatorType: model.OperatorTypeAssetManager,
	}, modifier)
	if err != nil {
		return err
	}
	if len(failed) > 0 {
		return errors.New(failed[0])
	}
	return nil
}

// 获取资产名称
func (s *assetService) getAssetName(assetType model.AssetType, materialID uint64) string {
	var name string
	switch assetType {
	case model.AssetTypeSmartBattery, model.AssetTypeNonSmartBattery:
		only, _ := ent.Database.BatteryModel.Query().Where(batterymodel.ID(materialID)).Only(s.ctx)
		name = only.Model
	case model.AssetTypeEbike:
		only, _ := ent.Database.EbikeBrand.Query().Where(ebikebrand.ID(materialID)).Only(s.ctx)
		name = only.Name
	case model.AssetTypeCabinetAccessory, model.AssetTypeEbikeAccessory, model.AssetTypeOtherAccessory:
		only, _ := ent.Database.Material.QueryNotDeleted().Where(material.ID(materialID)).Only(s.ctx)
		if only == nil {
			return "未知"
		}
		name = only.Name
	default:
		name = "未知"
	}
	return name
}

// BatchCreate 批量创建资产
func (s *assetService) BatchCreate(ctx echo.Context, req *model.AssetBatchCreateReq, modifier *model.Modifier) (failed []string, err error) {
	switch req.AssetType {
	case model.AssetTypeEbike:
		failed, err = s.BatchCreateEbike(ctx, modifier)
	case model.AssetTypeSmartBattery:
		failed, err = s.BatchCreateBattery(ctx, modifier)
	default:
		return nil, errors.New("未知类型")
	}
	return failed, err
}

// BatchCreateEbike 批量创建电车
// 0-型号:brand(需查询) 1-车架号:sn 2-仓库:warehouse 3-生产批次:exFactory 4-车牌号:plate 5-终端编号:machine 6-SIM卡:sim 7-颜色:color
func (s *assetService) BatchCreateEbike(ctx echo.Context, modifier *model.Modifier) (failed []string, err error) {
	// 查询所有电车属性
	attr, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.AssetType(model.AssetTypeEbike.Value())).Order(ent.Asc(assetattributes.FieldCreatedAt)).All(ctx.Request().Context())
	if len(attr) == 0 {
		return nil, errors.New("电车属性未配置")
	}
	var (
		rows [][]string
		sns  []string
	)
	rows, sns, failed, err = s.GetXlsxRows(ctx, 1, len(s.getEbikeColumns(ctx.Request().Context())), 1)
	if err != nil {
		failed = append(failed, err.Error())
	}

	// 获取所有型号
	brands := NewEbikeBrand().All()
	bm := make(map[string]uint64)
	for _, brand := range brands {
		bm[brand.Name] = brand.ID
	}

	// 获取仓库
	warehouseIds := make(map[string]uint64)
	warehouses, _ := ent.Database.Warehouse.QueryNotDeleted().All(ctx.Request().Context())
	for _, v := range warehouses {
		warehouseIds[v.Name] = v.ID
	}

	assetAll, _ := ent.Database.Asset.QueryNotDeleted().Where(asset.SnIn(sns...), asset.Type(model.AssetTypeEbike.Value())).All(s.ctx)
	exists := make(map[string]bool)
	for _, a := range assetAll {
		exists[a.Sn] = true
	}
	// 获取第一行数据的标题 对应的属性id
	title := rows[0]
	attrMap := make(map[string]uint64)
	for _, v := range attr {
		attrMap[v.Name] = v.ID
	}

	var titleID []uint64
	for i := 1; i < len(title); i++ {
		if id, ok := attrMap[title[i]]; ok {
			titleID = append(titleID, id)
		}
	}

	// 仓库分组
	groupWarehouse := make(map[uint64][]string)

	// rows 去除标题
	rows = rows[1:]
	for _, columns := range rows {
		bid, ok := bm[columns[0]]
		if !ok {
			failed = append(failed, fmt.Sprintf("型号未找到:%s", strings.Join(columns, ",")))
			continue
		}

		if _, ok = exists[columns[1]]; ok {
			failed = append(failed, fmt.Sprintf("车架号重复:%s", strings.Join(columns, ",")))
			continue
		}

		wID, ok := warehouseIds[columns[2]]
		// 判定仓库
		if !ok {
			failed = append(failed, fmt.Sprintf("仓库未找到:%s", strings.Join(columns, ",")))
			continue
		}

		name := s.getAssetName(model.AssetTypeEbike, bid)
		save, _ := s.orm.Create().
			SetSn(columns[1]).
			SetType(model.AssetTypeEbike.Value()).
			SetName(name).
			SetEnable(true).
			SetStatus(model.AssetStatusDelivering.Value()).
			SetRemark("批量导入").
			SetBrandID(bid).
			SetBrandName(columns[0]).
			SetCreator(modifier).
			SetLastModifier(modifier).
			SetLocationsType(model.AssetLocationsTypeWarehouse.Value()).
			SetLocationsID(wID).
			Save(ctx.Request().Context())
		if save == nil {
			failed = append(failed, fmt.Sprintf("保存失败:%s", strings.Join(columns, ",")))
			continue
		}
		// 获取标题对应id
		for i := 0; i < len(title)-3; i++ {
			// 判定属性值是否存在
			if b, _ := ent.Database.AssetAttributeValues.Query().Where(assetattributevalues.AttributeID(titleID[i]), assetattributevalues.AssetID(save.ID)).Exist(ctx.Request().Context()); b {
				failed = append(failed, fmt.Sprintf("属性值重复:%s", strings.Join(columns, ",")))
				continue
			}
			err = ent.Database.AssetAttributeValues.Create().
				SetValue(columns[3+i]).
				SetAssetID(save.ID).
				SetAttributeID(titleID[i]).
				Exec(ctx.Request().Context())
			if err != nil {
				failed = append(failed, fmt.Sprintf("保存失败:%s", strings.Join(columns, ",")))
				continue
			}
		}
		// 仓库分组加入物资id
		groupWarehouse[wID] = append(groupWarehouse[wID], save.Sn)
	}

	// 创建调拨单
	for k, v := range groupWarehouse {
		var details []model.AssetTransferCreateDetail
		for _, vl := range v {
			details = append(details, model.AssetTransferCreateDetail{
				AssetType: model.AssetTypeEbike,
				SN:        &vl,
			})
		}
		var fs []string
		_, fs, err = NewAssetTransfer().Transfer(context.Background(), &model.AssetTransferCreateReq{
			ToLocationType:    model.AssetLocationsTypeWarehouse,
			ToLocationID:      k,
			Reason:            "初始入库",
			AssetTransferType: model.AssetTransferTypeInitial,
			Details:           details,
			OperatorID:        modifier.ID,
			OperatorType:      model.OperatorTypeAssetManager,
		}, modifier)
		failed = append(failed, fs...)
		if err != nil {
			failed = append(failed, fmt.Sprintf("调拨单创建失败: %v", err))
		}
	}

	return failed, nil
}

// BatchCreateBattery 批量创建电池
// 0-城市:city 1-电池编号:sn 2-仓库:warehouse
func (s *assetService) BatchCreateBattery(ctx echo.Context, modifier *model.Modifier) (failed []string, err error) {
	var (
		rows [][]string
		sns  []string
	)
	rows, sns, failed, err = s.BaseService.GetXlsxRows(ctx, 1, 3, 1)
	if err != nil {
		failed = append(failed, err.Error())
	}
	// 查重
	items, _ := s.orm.QueryNotDeleted().Where(asset.SnIn(sns...), asset.Type(model.AssetTypeSmartBattery.Value())).All(s.ctx)
	m := make(map[string]bool)
	for _, item := range items {
		m[item.Sn] = true
	}

	// 查询城市
	cs := make(map[string]struct{})
	for _, row := range rows {
		cs[row[0]] = struct{}{}
	}
	var (
		cids []string
		cm   = make(map[string]uint64)
	)
	for k := range cs {
		cids = append(cids, k)
	}
	cities, _ := ent.Database.City.QueryNotDeleted().Where(city.NameIn(cids...)).All(s.ctx)
	for _, ci := range cities {
		cm[ci.Name] = ci.ID
	}

	// 查询型号是否存在
	modelInfo, _ := ent.Database.BatteryModel.Query().All(s.ctx)
	if len(modelInfo) == 0 {
		return nil, errors.New("电池型号未配置")
	}
	mm := make(map[string]uint64)
	for _, md := range modelInfo {
		mm[md.Model] = md.ID
	}

	// 获取仓库
	warehouseIds := make(map[string]uint64)
	warehouses, _ := ent.Database.Warehouse.QueryNotDeleted().All(ctx.Request().Context())
	for _, v := range warehouses {
		warehouseIds[v.Name] = v.ID
	}

	// 仓库分组
	groupWarehouse := make(map[uint64][]string)

	// rows 去除标题
	rows = rows[1:]
	for _, row := range rows {
		sn := row[1]
		if m[sn] {
			failed = append(failed, fmt.Sprintf("编号%s已存在", sn))
			continue
		}

		// 解析电池编号
		ab, err := adapter.ParseBatterySN(sn)
		if err != nil || ab.Model == "" {
			failed = append(failed, fmt.Sprintf("电池编号%s解析失败", sn))
			continue
		}

		// 城市
		cid, ok := cm[row[0]]
		if !ok {
			failed = append(failed, fmt.Sprintf("城市%s查询失败", row[0]))
			continue
		}

		// 型号
		mid, ok := mm[ab.Model]
		if !ok {
			failed = append(failed, fmt.Sprintf("型号%s查询失败", ab.Model))
			continue
		}

		// 仓库
		wid, ok := warehouseIds[row[2]]
		if !ok {
			failed = append(failed, fmt.Sprintf("仓库%s查询失败", row[2]))
			continue
		}

		name := s.getAssetName(model.AssetTypeSmartBattery, mid)
		save, err := s.orm.Create().
			SetBrandName(ab.Brand.String()).
			SetSn(sn).
			SetName(name).
			SetModelID(mid).
			SetCityID(cid).
			SetLocationsType(model.AssetLocationsTypeWarehouse.Value()).
			SetLocationsID(wid).
			SetStatus(model.AssetStatusDelivering.Value()).
			SetType(model.AssetTypeSmartBattery.Value()).
			SetEnable(true).
			SetRemark("批量导入").
			SetCreator(modifier).
			SetLastModifier(modifier).
			Save(s.ctx)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s保存失败: %v", sn, err))
			continue
		}

		// 仓库分组加入物资id
		groupWarehouse[wid] = append(groupWarehouse[wid], save.Sn)
	}

	// 创建调拨单
	for k, v := range groupWarehouse {
		var details []model.AssetTransferCreateDetail
		for _, vl := range v {
			details = append(details, model.AssetTransferCreateDetail{
				AssetType: model.AssetTypeEbike,
				SN:        &vl,
			})
		}
		var fs []string
		_, fs, err = NewAssetTransfer().Transfer(context.Background(), &model.AssetTransferCreateReq{
			ToLocationType:    model.AssetLocationsTypeWarehouse,
			ToLocationID:      k,
			Reason:            "初始入库",
			AssetTransferType: model.AssetTransferTypeInitial,
			Details:           details,
			OperatorID:        modifier.ID,
			OperatorType:      model.OperatorTypeAssetManager,
		}, modifier)
		failed = append(failed, fs...)
		if err != nil {
			failed = append(failed, fmt.Sprintf("调拨单创建失败: %v", err))
		}
	}

	return failed, nil
}

// Modify 修改资产
func (s *assetService) Modify(ctx context.Context, req *model.AssetModifyReq, modifier *model.Modifier) error {
	q := s.orm.QueryNotDeleted()
	update := s.orm.
		UpdateOneID(req.ID).
		SetLastModifier(modifier)
	if req.Enable != nil {
		update.SetEnable(*req.Enable)
	}
	if req.CityID != nil {
		only, _ := ent.Database.City.QueryNotDeleted().Where(city.ID(*req.CityID)).Only(ctx)
		if only == nil {
			return errors.New("城市不存在")
		}
		update.SetCityID(*req.CityID)
	}

	if req.BrandID != nil {
		q.Where(asset.StatusIn(model.AssetStatusStock.Value()))
		assets, _ := q.First(ctx)
		if assets == nil {
			return errors.New("电车状态库存中才允许修改")
		}
		only, _ := ent.Database.EbikeBrand.QueryNotDeleted().Where(ebikebrand.ID(*req.BrandID)).Only(ctx)
		if only == nil {
			return errors.New("品牌不存在")
		}
		update.SetBrandID(*req.BrandID)
	}
	if req.Remark != nil {
		update.SetRemark(*req.Remark)
	}
	err := update.Exec(ctx)
	if err != nil {
		return fmt.Errorf("资产修改失败: %w", err)
	}

	// 修改属性
	if len(req.Attribute) > 0 {
		for _, v := range req.Attribute {
			av, _ := ent.Database.AssetAttributeValues.Query().Where(assetattributevalues.AssetID(req.ID), assetattributevalues.AttributeID(v.AttributeID)).First(ctx)
			if av != nil {
				// 存在更新
				err = ent.Database.AssetAttributeValues.Update().Where(assetattributevalues.AssetID(req.ID), assetattributevalues.AttributeID(v.AttributeID)).SetValue(v.AttributeValue).Exec(ctx)
				if err != nil {
					return fmt.Errorf("资产属性修改失败: %w", err)
				}
			} else {
				// 不存在新增
				err = ent.Database.AssetAttributeValues.Create().
					SetValue(v.AttributeValue).
					SetAttributeID(v.AttributeID).
					SetAssetID(req.ID).
					Exec(ctx)
				if err != nil {
					return fmt.Errorf("资产属性修改失败: %w", err)
				}
			}

		}
	}
	return nil
}

// DownloadTemplate 批量导入模版下载
func (s *assetService) DownloadTemplate(ctx context.Context, t model.AssetType) (path string, name string, err error) {
	switch t {
	case model.AssetTypeEbike:
		return s.downloadEbikeTemplate(ctx)
	case model.AssetTypeSmartBattery:
		return s.downloadBatteryTemplate(ctx)
	default:
		return "", "", errors.New("未知类型")
	}
}

// 下载电车模版
func (s *assetService) downloadEbikeTemplate(ctx context.Context) (path string, name string, err error) {
	var rows [][]any
	rows = append(rows, s.getEbikeColumns(ctx))
	name = "导入电车模版"
	path = filepath.Join("runtime/export", fmt.Sprintf("%s-%s.xlsx", "", "ebike"))
	tools.NewExcel(path).AddValues(rows).Done()
	return path, name, nil
}

// 获取电车导出列
func (s *assetService) getEbikeColumns(ctx context.Context) (columns []any) {
	// 查询所有电车属性
	arr, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.AssetType(model.AssetTypeEbike.Value())).Order(ent.Asc(assetattributes.FieldCreatedAt)).All(ctx)
	// 固定列
	assetEbikeAttributesColumns := model.AssetEbikeAttributesColumns
	for _, v := range arr {
		// 配置列
		assetEbikeAttributesColumns = append(assetEbikeAttributesColumns, v.Name)
	}
	return assetEbikeAttributesColumns
}

// 下载电池模版
func (s *assetService) downloadBatteryTemplate(ctx context.Context) (path string, name string, err error) {
	data := [][]any{
		{"城市", "编号", "仓库"},
	}
	name = "导入电池模版"
	path = filepath.Join("runtime/export", fmt.Sprintf("%s-%s.xlsx", "", "battery"))
	tools.NewExcel(path).AddValues(data).Done()
	return path, name, nil
}

// Export 导出资产
func (s *assetService) Export(ctx context.Context, req *model.AssetListReq, m *model.Modifier) (model.AssetExportRes, error) {
	q := s.orm.QueryNotDeleted().
		WithCabinet().
		WithCity().
		WithStation(func(query *ent.EnterpriseStationQuery) {
			query.WithEnterprise()
		}).
		WithModel().
		WithOperator().
		WithValues().
		WithStore().
		WithWarehouse().
		WithBrand().
		WithValues().
		WithRider()
	s.filter(q, &req.AssetFilter)
	q.Order(ent.Desc(asset.FieldCreatedAt))

	if req.AssetType == nil {
		return model.AssetExportRes{}, errors.New("类型不能为空")
	}
	switch *req.AssetType {
	case model.AssetTypeSmartBattery:
		s.batteryFilter(q, &req.AssetFilter)
		return s.exportBattery(ctx, req, q, m), nil
	case model.AssetTypeEbike:
		s.ebikeFilter(q, &req.AssetFilter)
		return s.exportEbike(ctx, req, q, m), nil
	default:
		return model.AssetExportRes{}, errors.New("未知类型")
	}
}

// 导出电池
func (s *assetService) exportBattery(ctx context.Context, req *model.AssetListReq, q *ent.AssetQuery, m *model.Modifier) model.AssetExportRes {
	return NewAssetExportWithModifier(m).Start("电池列表", req.AssetFilter, nil, "", func(path string) {
		items, _ := q.All(context.Background())
		var rows tools.ExcelItems
		title := []any{
			"城市",
			"归属",
			"库存位置",
			"电池品牌",
			"电池型号",
			"电池编号",
			"库存状态",
			"是否启用",
			"备注",
		}
		rows = append(rows, title)
		for _, item := range items {
			var cityName, belong, assetLocations string
			if item.Edges.City != nil {
				cityName = item.Edges.City.Name
			}
			belong = "平台"
			if item.LocationsType == model.AssetLocationsTypeStation.Value() {
				belong = "代理商"
				if item.Edges.Station != nil && item.Edges.Station.Edges.Enterprise != nil {
					belong = item.Edges.Station.Edges.Enterprise.Name
				}
			}
			switch item.LocationsType {
			case model.AssetLocationsTypeWarehouse.Value():
				assetLocations = "[仓库]"
				if item.Edges.Warehouse != nil {
					assetLocations += item.Edges.Warehouse.Name
				}
			case model.AssetLocationsTypeStore.Value():
				assetLocations = "[门店]"
				if item.Edges.Store != nil {
					assetLocations += item.Edges.Store.Name
				}
			case model.AssetLocationsTypeCabinet.Value():
				assetLocations = "[电柜]"
				if item.Edges.Cabinet != nil {
					assetLocations += item.Edges.Cabinet.Name
				}
			case model.AssetLocationsTypeStation.Value():
				assetLocations = "[站点]"
				if item.Edges.Station != nil {
					assetLocations += item.Edges.Station.Name
				}
			case model.AssetLocationsTypeRider.Value():
				assetLocations = "[骑手]"
				if item.Edges.Rider != nil {
					assetLocations += item.Edges.Rider.Name + "-" + item.Edges.Rider.Phone
				}
			case model.AssetLocationsTypeOperation.Value():
				assetLocations = "[运维]"
				if item.Edges.Operator != nil {
					assetLocations += item.Edges.Operator.Name
				}
			}

			var modelStr string
			if item.Edges.Model != nil {
				modelStr = item.Edges.Model.Model
			}

			var brandName string
			if item.Type == model.AssetTypeNonSmartBattery.Value() || item.Type == model.AssetTypeSmartBattery.Value() {
				brandName = item.BrandName
			}

			enable := "否"
			if item.Enable {
				enable = "是"
			}

			row := []any{
				cityName,
				belong,
				assetLocations,
				brandName,
				modelStr,
				item.Sn,
				model.AssetStatus(item.Status).String(),
				enable,
				item.Remark,
			}
			rows = append(rows, row)
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}

// 导出电车
func (s *assetService) exportEbike(ctx context.Context, req *model.AssetListReq, q *ent.AssetQuery, m *model.Modifier) model.AssetExportRes {
	t, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.AssetType((model.AssetTypeEbike).Value())).WithValues().All(ctx)
	return NewAssetExportWithModifier(m).Start("电车列表", req.AssetFilter, nil, "", func(path string) {
		items, _ := q.All(context.Background())
		var rows tools.ExcelItems
		title := []any{
			"归属",
			"库存位置",
			"型号",
			"电车编号",
			"库存状态",
			"是否启用",
			"赠送状态",
			"备注",
		}

		for _, v := range t {
			title = append(title, v.Name)
		}
		rows = append(rows, title)
		for _, item := range items {
			var belong, assetLocations string
			belong = "平台"
			if item.LocationsType == model.AssetLocationsTypeStation.Value() {
				belong = "代理商"
				if item.Edges.Station != nil && item.Edges.Station.Edges.Enterprise != nil {
					belong = item.Edges.Station.Edges.Enterprise.Name
				}
			}
			switch item.LocationsType {
			case model.AssetLocationsTypeWarehouse.Value():
				assetLocations = "[仓库]"
				if item.Edges.Warehouse != nil {
					assetLocations += item.Edges.Warehouse.Name
				}
			case model.AssetLocationsTypeStore.Value():
				assetLocations = "[门店]"
				if item.Edges.Store != nil {
					assetLocations += item.Edges.Store.Name
				}
			case model.AssetLocationsTypeCabinet.Value():
				assetLocations = "[电柜]"
				if item.Edges.Cabinet != nil {
					assetLocations += item.Edges.Cabinet.Name
				}
			case model.AssetLocationsTypeStation.Value():
				assetLocations = "[站点]"
				if item.Edges.Station != nil {
					assetLocations += item.Edges.Station.Name
				}
			case model.AssetLocationsTypeRider.Value():
				assetLocations = "[骑手]"
				if item.Edges.Rider != nil {
					assetLocations += item.Edges.Rider.Name + "-" + item.Edges.Rider.Phone
				}
			case model.AssetLocationsTypeOperation.Value():
				assetLocations = "[运维]"
				if item.Edges.Operator != nil {
					assetLocations += item.Edges.Operator.Name
				}
			}

			var brandName string
			if item.Edges.Brand != nil {
				brandName = item.Edges.Brand.Name
			}

			rto := "未赠送"
			if item.RtoRiderID != nil {
				rto = "[赠送]-"
				if item.Edges.Rider != nil {
					rto += item.Edges.Rider.Name
				}
			}
			enable := "禁用"
			if item.Enable {
				enable = "启用"
			}

			row := []any{
				belong,
				assetLocations,
				brandName,
				item.Sn,
				model.AssetStatus(item.Status).String(),
				enable,
				rto,
				item.Remark,
			}
			for _, v := range t {
				val, _ := item.QueryValues().Where(assetattributevalues.AssetID(item.ID), assetattributevalues.AttributeID(v.ID)).First(context.Background())
				var value string
				if val != nil {
					value = val.Value
				}
				row = append(row, value)
			}
			rows = append(rows, row)
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}

// List 资产列表
func (s *assetService) List(ctx context.Context, req *model.AssetListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().
		WithCabinet().
		WithCity().
		WithStation(func(query *ent.EnterpriseStationQuery) {
			query.WithEnterprise()
		}).
		WithModel().
		WithOperator().
		WithValues().
		WithStore().
		WithWarehouse().
		WithBrand().
		WithValues().
		WithRider()
	s.filter(q, &req.AssetFilter)
	if req.AssetType != nil {
		switch *req.AssetType {
		case model.AssetTypeSmartBattery:
			s.batteryFilter(q, &req.AssetFilter)
		case model.AssetTypeEbike:
			s.ebikeFilter(q, &req.AssetFilter)
		default:

		}
	}
	q.Order(ent.Desc(asset.FieldCreatedAt))
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Asset) *model.AssetListRes {
		return s.DetailForList(item)
	})
}

func (s *assetService) DetailForList(item *ent.Asset) *model.AssetListRes {
	var cityName, belong, assetLocations string
	var cityID uint64
	if item.Edges.City != nil {
		cityName = item.Edges.City.Name
		cityID = item.Edges.City.ID
	}
	belong = "平台"
	if item.LocationsType == model.AssetLocationsTypeStation.Value() {
		belong = "代理商"
		if item.Edges.Station != nil && item.Edges.Station.Edges.Enterprise != nil {
			belong = item.Edges.Station.Edges.Enterprise.Name
		}
	}
	switch item.LocationsType {
	case model.AssetLocationsTypeWarehouse.Value():
		assetLocations = "[仓库]"
		if item.Edges.Warehouse != nil {
			assetLocations += item.Edges.Warehouse.Name
		}
	case model.AssetLocationsTypeStore.Value():
		assetLocations = "[门店]"
		if item.Edges.Store != nil {
			assetLocations += item.Edges.Store.Name
		}
	case model.AssetLocationsTypeCabinet.Value():
		assetLocations = "[电柜]"
		if item.Edges.Cabinet != nil {
			assetLocations += item.Edges.Cabinet.Name
		}
	case model.AssetLocationsTypeStation.Value():
		assetLocations = "[站点]"
		if item.Edges.Station != nil {
			assetLocations += item.Edges.Station.Name
		}
	case model.AssetLocationsTypeRider.Value():
		assetLocations = "[骑手]"
		if item.Edges.Rider != nil {
			assetLocations += item.Edges.Rider.Name + "-" + item.Edges.Rider.Phone
		}
	case model.AssetLocationsTypeOperation.Value():
		assetLocations = "[运维]"
		if item.Edges.Operator != nil {
			assetLocations += item.Edges.Operator.Name
		}
	}

	var modelStr string
	if item.Edges.Model != nil {
		modelStr = item.Edges.Model.Model
	}

	var brandName string
	var brandID uint64
	if item.Edges.Brand != nil {
		brandName = item.Edges.Brand.Name
		brandID = item.Edges.Brand.ID
	}
	if item.Type == model.AssetTypeNonSmartBattery.Value() || item.Type == model.AssetTypeSmartBattery.Value() {
		brandName = item.BrandName
	}

	res := &model.AssetListRes{
		ID:             item.ID,
		CityName:       cityName,
		CityID:         cityID,
		Belong:         belong,
		AssetLocations: assetLocations,
		LocationsID:    item.LocationsID,
		Brand:          brandName,
		BrandID:        brandID,
		Model:          modelStr,
		SN:             item.Sn,
		AssetStatus:    model.AssetStatus(item.Status).String(),
		Status:         model.AssetStatus(item.Status),
		Enable:         item.Enable,
		Remark:         item.Remark,
	}

	// 电车赠送状态
	res.Rto = "未赠送"
	if item.RtoRiderID != nil {
		r, _ := item.QueryRtoRider().First(context.Background())
		if r != nil {
			res.Rto = "[赠送]-" + r.Name + "-" + r.Phone
		}
	}

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

	// 库存为骑手即为使用中状态
	if item.LocationsType == model.AssetLocationsTypeRider.Value() {
		res.Status = model.AssetStatusUsing
		res.AssetStatus = model.AssetStatusUsing.String()
	}

	return res
}

// 资产公共筛选
func (s *assetService) filter(q *ent.AssetQuery, req *model.AssetFilter) {
	if req.SN != nil {
		q.Where(asset.SnContainsFold(*req.SN))
	}
	if req.AssetType != nil {
		q.Where(asset.Type(req.AssetType.Value()))
	}
	if req.OwnerType != nil {
		// 平台
		if *req.OwnerType == 1 {
			q.Where(asset.LocationsTypeNEQ(model.AssetLocationsTypeStation.Value()))
		}
		// 代理商
		if *req.OwnerType == 2 {
			q.Where(
				asset.LocationsTypeEQ(model.AssetLocationsTypeStation.Value()),
			)
			if req.EnterpriseID != nil {
				q.Where(
					asset.HasStationWith(enterprisestation.HasEnterpriseWith(enterprise.ID(*req.EnterpriseID))),
				)
			}
		}
	}
	if req.LocationsType != nil {
		q.Where(asset.LocationsType(req.LocationsType.Value()))
		switch *req.LocationsType {
		case model.AssetLocationsTypeWarehouse, model.AssetLocationsTypeStore, model.AssetLocationsTypeStation, model.AssetLocationsTypeOperation:
			if req.LocationsID != nil {
				q.Where(
					asset.LocationsID(*req.LocationsID),
				)
			}
		case model.AssetLocationsTypeCabinet:
			if req.LocationsKeyword != nil {
				q.Where(
					asset.HasCabinetWith(
						cabinet.NameContains(*req.LocationsKeyword),
					),
				)
			}
			if req.LocationsID != nil {
				q.Where(
					asset.LocationsID(*req.LocationsID),
				)
			}

		case model.AssetLocationsTypeRider:
			if req.LocationsKeyword != nil {
				q.Where(
					asset.HasRiderWith(
						rider.NameContains(*req.LocationsKeyword),
					),
				)
			}

		}
	}
	if req.Status != nil {
		if *req.Status == model.AssetStatusUsing {
			q.Where(
				asset.LocationsType(model.AssetLocationsTypeRider.Value()),
				asset.Status(model.AssetStatusStock.Value()),
			)
		} else {
			q.Where(asset.Status(req.Status.Value()))
		}
	}
	if req.Enable != nil {
		q.Where(asset.Enable(*req.Enable))
	}
	if req.MaterialID != nil {
		q.Where(asset.MaterialID(*req.MaterialID))
	}
	if req.ModelID != nil {
		q.Where(asset.ModelID(*req.ModelID))
	}
	if req.Model != nil {
		q.Where(asset.HasModelWith(batterymodel.Model(*req.Model)))
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
			q.Where(asset.HasValuesWith(assetattributevalues.AttributeID(attributeID), assetattributevalues.ValueContains(attributeValue)))
		}
	}
	if req.AssetManagerID != 0 {
		// 查询库管人员配置的仓库数据
		wIds := make([]uint64, 0)
		am, _ := ent.Database.AssetManager.QueryNotDeleted().WithBelongWarehouses().
			Where(
				assetmanager.ID(req.AssetManagerID),
				assetmanager.HasBelongWarehousesWith(warehouse.DeletedAtIsNil()),
			).First(context.Background())
		if am != nil {
			for _, wh := range am.Edges.BelongWarehouses {
				wIds = append(wIds, wh.ID)
			}
		}
		q.Where(
			asset.LocationsType(model.AssetLocationsTypeWarehouse.Value()),
			asset.LocationsIDIn(wIds...),
		)
	}
	if req.EmployeeID != 0 {
		// 查询门店人员配置的门店数据
		sIds := make([]uint64, 0)
		ep, _ := ent.Database.Employee.QueryNotDeleted().WithStores().
			Where(
				employee.ID(req.EmployeeID),
				employee.HasStoresWith(store.DeletedAtIsNil()),
			).First(context.Background())
		if ep != nil {
			for _, st := range ep.Edges.Stores {
				sIds = append(sIds, st.ID)
			}
		}
		q.Where(
			asset.LocationsType(model.AssetLocationsTypeStore.Value()),
			asset.LocationsIDIn(sIds...),
		)
	}
	if req.Battery != nil {
		q.Where(
			asset.TypeIn(
				model.AssetTypeSmartBattery.Value(),
				model.AssetTypeNonSmartBattery.Value(),
			),
		)
	}
}

// 资产电池筛选条件
func (s *assetService) batteryFilter(q *ent.AssetQuery, req *model.AssetFilter) *ent.AssetQuery {
	if req.ModelID != nil {
		q.Where(asset.ModelID(*req.ModelID))
	}
	if req.CityID != nil {
		q.Where(asset.CityID(*req.CityID))
	}
	return q
}

// 资产电车筛选条件
func (s *assetService) ebikeFilter(q *ent.AssetQuery, req *model.AssetFilter) *ent.AssetQuery {
	// 电车品牌
	if req.BrandID != nil {
		q.Where(asset.BrandID(*req.BrandID))
	}
	if req.Rto != nil && *req.Rto {
		// 电车赠送
		q.Where(asset.RtoRiderIDNotNil())
	} else {
		// 电车未赠送
		q.Where(asset.RtoRiderIDIsNil())
	}
	if req.Keyword != nil {
		attributes, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.Key("plate")).First(s.ctx)
		q.Where(
			asset.Or(
				asset.SnContainsFold(*req.Keyword),
				asset.HasValuesWith(assetattributevalues.AttributeID(attributes.ID), assetattributevalues.ValueContainsFold(*req.Keyword)),
			),
		)
	}
	return q
}

// Count 查询有效的资产数量
func (s *assetService) Count(ctx context.Context, req *model.AssetFilter) (res *model.AssetNumRes) {
	q := s.orm.QueryNotDeleted().Where(asset.StatusNEQ(model.AssetStatusScrap.Value()), asset.CheckAtIsNil())
	s.filter(q, req)
	count, _ := q.Count(ctx)
	res = &model.AssetNumRes{Num: count}
	as, _ := q.First(ctx)
	if as != nil {
		ty := model.AssetType(as.Type)
		res.AssetType = &ty
		res.AssetID = &as.ID
	}
	return
}

// QuerySn 通过SN查询被未盘点资产
func (s *assetService) QuerySn(sn string) (bat *ent.Asset, err error) {
	return s.orm.Query().WithModel().Where(asset.Sn(sn), asset.CheckAtIsNil()).First(s.ctx)
}

// QueryID 通过ID查询资产
func (s *assetService) QueryID(id uint64) (*ent.Asset, error) {
	return s.orm.Query().WithModel().Where(asset.ID(id), asset.CheckAtIsNil()).First(s.ctx)
}

// QueryAssetByLocation 查询某个位置在库存中的资产
func (s *assetService) QueryAssetByLocation(req model.QueryAssetReq) (*ent.Asset, error) {
	q := s.orm.Query().WithModel().Where(
		asset.LocationsType(req.LocationsType.Value()),
		asset.LocationsID(req.LocationsID),
		asset.CheckAtIsNil(),
	)
	if req.ID != nil {
		q.Where(asset.ID(*req.ID))
	}
	if req.Sn != nil {
		q.Where(asset.Sn(*req.Sn))
	}
	return q.First(s.ctx)
}

// QueryRiderID 通过骑手ID查询电池资产
func (s *assetService) QueryRiderID(id uint64) (*ent.Asset, error) {
	return s.orm.Query().WithModel().Where(
		asset.LocationsID(id),
		asset.LocationsType(model.AssetLocationsTypeRider.Value()),
		asset.TypeIn(model.AssetTypeSmartBattery.Value(), model.AssetTypeNonSmartBattery.Value()),
		asset.CheckAtIsNil(),
	).First(s.ctx)
}

// QueryNonSmartBattery 查询一个符合条件的非智能电池
func (s *assetService) QueryNonSmartBattery(req *model.QueryAssetBatteryReq) (bat *ent.Asset, err error) {
	q := s.orm.Query().WithModel().Where(asset.Type(model.AssetTypeNonSmartBattery.Value()), asset.CheckAtIsNil(), asset.Status(model.AssetStatusStock.Value())).Limit(1)
	if req.LocationsType != nil {
		q.Where(asset.LocationsType(req.LocationsType.Value()))
	}
	if req.LocationsID != nil {
		q.Where(asset.LocationsID(*req.LocationsID))
	}
	q.Where(asset.ModelID(req.ModelID))
	item, _ := q.First(s.ctx)
	if item == nil {
		return nil, errors.New("未找到符合条件的非智能电池")
	}
	return item, nil
}

// CheckAsset 检查非智能电池数量
func (s *assetService) CheckAsset(locationType model.AssetLocationsType, locationTypeID uint64, m string) (*ent.Asset, error) {
	item, _ := s.orm.QueryNotDeleted().Where(
		asset.LocationsID(locationTypeID),
		asset.LocationsType(locationType.Value()),
		asset.Status(model.AssetStatusStock.Value()),
		asset.Type(model.AssetTypeNonSmartBattery.Value()),
		asset.HasModelWith(batterymodel.Model(m)),
	).WithModel().First(s.ctx)
	if item == nil {
		return nil, errors.New("未找到符合条件的非智能电池")
	}
	return item, nil
}

func (s *assetService) CurrentBatteryNum(ids []uint64, locationsType model.AssetLocationsType) map[uint64]int {
	var result []struct {
		TargetID uint64 `json:"target_id"`
		Sum      int    `json:"sum"`
	}
	v := make([]interface{}, len(ids))
	for i := range v {
		v[i] = ids[i]
	}
	_ = s.orm.Query().
		Where(asset.LocationsType(locationsType.Value())).
		Modify(func(sel *sql.Selector) {
			sel.Where(sql.In(sel.C(asset.FieldLocationsID), v...)).
				Select(
					sql.As(sel.C(asset.FieldLocationsID), "target_id"),
					sql.As(sql.Count(asset.FieldID), "sum"),
				).
				GroupBy(asset.FieldLocationsID)
		}).
		Scan(s.ctx, &result)
	m := make(map[uint64]int)
	for _, r := range result {
		m[r.TargetID] = r.Sum
	}
	return m
}

func (s *assetService) CurrentBattery(id uint64, locationsType model.AssetLocationsType) int {
	return s.CurrentBatteryNum([]uint64{id}, locationsType)[id]
}

// StoreCurrent 列出当前门店所有电池物资
func (s *assetService) StoreCurrent(id uint64) []model.InventoryNum {
	ins := make([]model.InventoryNum, 0)
	_ = s.orm.Query().
		Where(
			asset.LocationsType(model.AssetLocationsTypeStore.Value()),
			asset.LocationsID(id),
			asset.Status(model.AssetStatusStock.Value()),
		).
		Modify(func(sel *sql.Selector) {
			t := sql.Table(batterymodel.Table).As("model")
			sel.LeftJoin(t).On(t.C(batterymodel.FieldID), sel.C(asset.FieldModelID))
			sel.GroupBy(asset.FieldName, t.C(batterymodel.FieldModel)).
				Select(asset.FieldName, t.C(batterymodel.FieldModel)).
				AppendSelectExprAs(sql.Raw(fmt.Sprintf("%s IS NOT NULL", asset.FieldModelID)), "battery").
				AppendSelectExprAs(sql.Raw(fmt.Sprintf("Count(%s)", asset.FieldID)), "num")
		}).
		Scan(s.ctx, &ins)

	return ins
}

// RiderBusiness 骑手业务 电池 / 电车 出入库
// 此方法操作库存
func (s *assetService) RiderBusiness(req *model.StockBusinessReq) (err error) {
	if req.StoreID == nil && req.EbikeStoreID == nil && req.BatStoreID == nil && req.CabinetID == nil && req.EnterpriseID == nil && req.StationID == nil {
		err = errors.New("参数校验错误")
		return
	}

	// 如果是骑手自己操作 激活和取消寄存拿走电池 会有电柜任务已经执行调拨
	if s.operator.Type == model.OperatorTypeRider && req.CabinetID != nil && (req.AssetTransferType == model.AssetTransferTypeActive || req.AssetTransferType == model.AssetTransferTypeContinue) {
		return
	}

	// 查询资产
	var ebikeInfo *ent.Asset
	var batteryInfo *ent.Asset
	var fromLocationType model.AssetLocationsType
	var fromLocationID uint64
	var toLocationType model.AssetLocationsType
	var toLocationID uint64
	var ebiketoLocationID uint64
	var ebiketoLocationType model.AssetLocationsType
	var ebikeFromLocationType model.AssetLocationsType
	var ebikeFromLocationID uint64
	details := make([]model.AssetTransferCreateDetail, 0)
	ebikeDetails := make([]model.AssetTransferCreateDetail, 0)
	assetType := model.AssetTypeSmartBattery

	if req.Ebike != nil {
		ebikeInfo, _ = NewAsset().QuerySn(req.Ebike.Sn)
		if ebikeInfo == nil {
			err = errors.New("电车不存在")
			return
		}
	}

	if req.Battery != nil {
		if req.Battery.SN == "" {
			batteryInfo, _ = NewAsset().QueryID(req.Battery.ID)
		} else {
			batteryInfo, _ = NewAsset().QuerySn(req.Battery.SN)
		}
		if batteryInfo == nil {
			err = errors.New("电池不存在")
			return
		}
	}
	var storeID uint64
	var ebikeStoreID uint64
	if req.StoreID != nil {
		storeID = *req.StoreID
	}
	if req.BatStoreID != nil {
		storeID = *req.BatStoreID
	}
	if req.EbikeStoreID != nil {
		ebikeStoreID = *req.EbikeStoreID
	}

	// 激活和取消寄存 需要判定非智能库存
	if req.Battery == nil && (req.AssetTransferType == model.AssetTransferTypeActive || req.AssetTransferType == model.AssetTransferTypeContinue) {

		if req.StoreID != nil || req.BatStoreID != nil {
			batteryInfo, _ = s.CheckBusinessBattery(req, model.AssetLocationsTypeStore, storeID)
		}
		// 团签业务
		if req.CabinetID == nil && req.EnterpriseID != nil && req.StationID != nil {
			batteryInfo, _ = s.CheckBusinessBattery(req, model.AssetLocationsTypeStation, *req.StationID)
		}

		// 电柜业务
		if req.CabinetID != nil {
			batteryInfo, _ = s.CheckBusinessBattery(req, model.AssetLocationsTypeCabinet, *req.CabinetID)
		}
		if batteryInfo == nil {
			err = errors.New("电池库存不足")
			return
		}
	}
	if batteryInfo != nil && batteryInfo.Type == model.AssetTypeNonSmartBattery.Value() {
		assetType = model.AssetTypeNonSmartBattery
	}
	switch req.AssetTransferType {
	case model.AssetTransferTypeActive, model.AssetTransferTypeContinue:
		// 激活和取消寄存 某个位置的库存调拨到骑手
		toLocationType = model.AssetLocationsTypeRider
		toLocationID = req.RiderID
		ebiketoLocationType = model.AssetLocationsTypeRider
		ebiketoLocationID = req.RiderID
		if batteryInfo != nil {
			fromLocationType = model.AssetLocationsType(batteryInfo.LocationsType)
			fromLocationID = batteryInfo.LocationsID
		}
		if ebikeInfo != nil {
			ebikeFromLocationType = model.AssetLocationsType(ebikeInfo.LocationsType)
			ebikeFromLocationID = ebikeInfo.LocationsID
		}
	case model.AssetTransferTypePause, model.AssetTransferTypeUnSubscribe:
		// 寄存和退租 骑手的库存调拨到某个位置
		fromLocationType = model.AssetLocationsTypeRider
		fromLocationID = req.RiderID
		ebikeFromLocationType = model.AssetLocationsTypeRider
		ebikeFromLocationID = req.RiderID
		// 电池退租使用参数
		if req.BatStoreID != nil || req.StoreID != nil {
			toLocationType = model.AssetLocationsTypeStore
			toLocationID = storeID
		}
		if req.CabinetID != nil {
			toLocationType = model.AssetLocationsTypeCabinet
			toLocationID = *req.CabinetID
		}
		if req.EnterpriseID != nil && req.StationID != nil {
			toLocationType = model.AssetLocationsTypeStation
			toLocationID = *req.StationID
		}
		// 电车退租使用参数
		if req.EbikeStoreID != nil {
			ebiketoLocationType = model.AssetLocationsTypeStore
			ebiketoLocationID = ebikeStoreID
		}
	default:
		return errors.New("业务类型错误")
	}

	if ebikeInfo != nil {
		ebikeDetails = append(ebikeDetails, model.AssetTransferCreateDetail{
			AssetType: model.AssetTypeEbike,
			SN:        silk.String(ebikeInfo.Sn),
		})
	}
	if batteryInfo != nil {
		if assetType == model.AssetTypeSmartBattery {
			details = append(details, model.AssetTransferCreateDetail{
				AssetType: assetType,
				SN:        silk.String(batteryInfo.Sn),
			})
		} else {
			details = append(details, model.AssetTransferCreateDetail{
				AssetType: assetType,
				Num:       silk.UInt(1),
				ModelID:   batteryInfo.ModelID,
			})
		}
	}

	if len(details) != 0 {
		// 创建调拨单
		_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
			FromLocationType:  &fromLocationType,
			FromLocationID:    &fromLocationID,
			ToLocationType:    toLocationType,
			ToLocationID:      toLocationID,
			Details:           details,
			Reason:            req.AssetTransferType.String() + "骑手业务",
			AssetTransferType: req.AssetTransferType,
			OperatorID:        s.operator.ID,
			OperatorType:      s.operator.Type,
			AutoIn:            true,
			SkipLimit:         true,
		}, &model.Modifier{
			ID:    s.operator.ID,
			Name:  s.operator.Name,
			Phone: s.operator.Phone,
		})
		if err != nil {
			return err
		}
		if len(failed) > 0 {
			return errors.New(failed[0])
		}
	}

	if len(ebikeDetails) != 0 {
		// 电车创建调拨单
		_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
			FromLocationType:  &ebikeFromLocationType,
			FromLocationID:    &ebikeFromLocationID,
			ToLocationType:    ebiketoLocationType,
			ToLocationID:      ebiketoLocationID,
			Details:           ebikeDetails,
			Reason:            req.AssetTransferType.String() + "骑手业务",
			AssetTransferType: req.AssetTransferType,
			OperatorID:        s.operator.ID,
			OperatorType:      s.operator.Type,
			AutoIn:            true,
			SkipLimit:         true,
		}, &model.Modifier{
			ID:    s.operator.ID,
			Name:  s.operator.Name,
			Phone: s.operator.Phone,
		})
		if err != nil {
			return err
		}
		if len(failed) > 0 {
			return errors.New(failed[0])
		}
	}
	return
}

// CheckBusinessBattery 激活和取消寄存非智能电池数量判定
func (s *assetService) CheckBusinessBattery(req *model.StockBusinessReq, locationsType model.AssetLocationsType, locationsID uint64) (batteryInfo *ent.Asset, err error) {
	// 判定非智能电池库存
	batteryInfo, _ = NewAsset().CheckAsset(locationsType, locationsID, req.Model)
	if batteryInfo == nil {
		err = errors.New("电池库存不足")
		return
	}
	if batteryInfo.Edges.Model == nil {
		err = errors.New("电池型号不存在")
		return
	}
	return batteryInfo, nil
}

// 查询资产数量

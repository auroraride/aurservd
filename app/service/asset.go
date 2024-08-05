package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/auroraride/adapter"
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetattributes"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/internal/ent/batterymodelnew"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
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
	// 入库状态默认为待入库
	assetStatus := model.AssetStatusPending.Value()

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
		modelInfo, _ := ent.Database.BatteryModelNew.QueryNotDeleted().Where(batterymodelnew.Model(ab.Model)).Only(ctx)
		if modelInfo == nil {
			return fmt.Errorf("电池型号%s不存在", ab.Model)
		}
		q.SetNillableModelID(&modelInfo.ID).
			SetBrandName(ab.Brand.String()).
			SetNillableCityID(req.CityID)
	case model.AssetTypeEbike:
		if req.BrandID == nil {
			return errors.New("品牌不能为空")
		}
		q.SetBrandID(*req.BrandID)
	default:
		return errors.New("未知类型")
	}

	name := s.getAssetName(req.AssetType)
	if req.Name != nil {
		name = *req.Name
	}

	q.SetType(req.AssetType.Value()).
		SetName(name).
		SetNillableSn(req.SN).
		SetEnable(enable).
		SetStatus(assetStatus).
		SetCreator(modifier).
		SetLastModifier(modifier)
	if req.LocationsType != nil && req.LocationsID != nil {
		q.
			SetLocationsType((*req.LocationsType).Value()).
			SetLocationsID(*req.LocationsID)
	}
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
	return nil
}

// 获取资产名称
func (s *assetService) getAssetName(assetType model.AssetType) string {
	var name string
	switch assetType {
	case model.AssetTypeSmartBattery:
		name = "智能电池"
	case model.AssetTypeEbike:
		name = "电车"
	case model.AssetTypeNonSmartBattery:
		name = "非智能电池"
	case model.AssetTypeCabinetAccessory:
		name = "电柜配件"
	case model.AssetTypeOtherAccessory:
		name = "其它配件"
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
// 0-型号:brand(需查询) 1-车架号:sn 2-生产批次:exFactory 3-车牌号:plate 4-终端编号:machine 5-SIM卡:sim 6-颜色:color
func (s *assetService) BatchCreateEbike(ctx echo.Context, modifier *model.Modifier) (failed []string, err error) {
	// 查询所有电车属性
	attr, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.AssetType(model.AssetTypeEbike.Value())).Order(ent.Asc(assetattributes.FieldCreatedAt)).All(ctx.Request().Context())
	if len(attr) == 0 {
		return nil, errors.New("电车属性未配置")
	}
	rows, sns, failed, err := s.GetXlsxRows(ctx, 1, len(s.getEbikeColumns(ctx.Request().Context())), 1)
	if err != nil || len(failed) != 0 {
		return failed, err
	}

	// 获取所有型号
	brands := NewEbikeBrand().All()
	bm := make(map[string]uint64)
	for _, brand := range brands {
		bm[brand.Name] = brand.ID
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

		name := s.getAssetName(model.AssetTypeEbike)
		save, _ := s.orm.Create().
			SetSn(columns[1]).
			SetType(model.AssetTypeEbike.Value()).
			SetName(name).
			SetEnable(true).
			SetStatus(model.AssetStatusPending.Value()).
			SetRemark("批量导入").
			SetBrandID(bid).
			SetBrandName(columns[0]).
			SetCreator(modifier).
			SetLastModifier(modifier).
			Save(ctx.Request().Context())
		if save == nil {
			failed = append(failed, fmt.Sprintf("保存失败:%s", strings.Join(columns, ",")))
			continue
		}
		// 获取标题对应id
		for i := 1; i < len(title)-2; i++ {
			// 判定属性值是否存在
			if b, _ := ent.Database.AssetAttributeValues.Query().Where(assetattributevalues.AttributeID(titleID[i-1]), assetattributevalues.AssetID(save.ID)).Exist(ctx.Request().Context()); b {
				failed = append(failed, fmt.Sprintf("属性值重复:%s", strings.Join(columns, ",")))
				continue
			}
			err = ent.Database.AssetAttributeValues.Create().
				SetValue(columns[i]).
				SetAssetID(save.ID).
				SetAttributeID(titleID[i-1]).
				Exec(ctx.Request().Context())
			if err != nil {
				failed = append(failed, fmt.Sprintf("保存失败:%s", strings.Join(columns, ",")))
			}
		}
	}
	return failed, nil
}

// BatchCreateBattery 批量创建电池
// 0-城市:city 1-电池编号:sn
func (s *assetService) BatchCreateBattery(ctx echo.Context, modifier *model.Modifier) (failed []string, err error) {
	rows, sns, failed, err := s.BaseService.GetXlsxRows(ctx, 2, 2, 1)
	if err != nil || len(failed) != 0 {
		return failed, err
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
	modelInfo, _ := ent.Database.BatteryModelNew.QueryNotDeleted().Where(batterymodelnew.Type(1)).All(s.ctx)
	if len(modelInfo) == 0 {
		return nil, errors.New("电池型号未配置")
	}
	mm := make(map[string]uint64)
	for _, md := range modelInfo {
		mm[md.Model] = md.ID
	}

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

		creator := s.orm.Create()

		// 城市
		if cid, ok := cm[row[0]]; ok {
			creator.SetCityID(cid)
		} else {
			failed = append(failed, fmt.Sprintf("城市%s查询失败", row[0]))
			continue
		}

		// 型号
		if mid, ok := mm[ab.Model]; ok {
			creator.SetModelID(mid)
		} else {
			failed = append(failed, fmt.Sprintf("型号%s查询失败", ab.Model))
			continue
		}
		name := s.getAssetName(model.AssetTypeSmartBattery)
		err = creator.
			SetBrandName(ab.Brand.String()).
			SetSn(sn).
			SetName(name).
			SetStatus(model.AssetStatusPending.Value()).
			SetType(model.AssetTypeSmartBattery.Value()).
			SetEnable(true).
			SetRemark("批量导入").
			SetCreator(modifier).
			SetLastModifier(modifier).
			Exec(s.ctx)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s保存失败: %v", sn, err))
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
		q.Where(asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusPending.Value()))
		assets, _ := q.First(ctx)
		if assets == nil {
			return errors.New("电车状态须为待入库或库存中才允许修改")
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
			if err = ent.Database.AssetAttributeValues.UpdateOneID(v.AttributeValueID).SetValue(v.AttributeValue).Exec(ctx); err != nil {
				return fmt.Errorf("资产属性修改失败: %w", err)
			}
		}
	}
	return nil
}

// Delete 删除资产
func (s *assetService) Delete(ctx context.Context, id uint64) error {
	bat, _ := s.orm.QueryNotDeleted().
		Where(
			asset.ID(id),
			asset.Status(model.AssetStatusPending.Value()),
		).
		First(ctx)
	if bat == nil {
		return fmt.Errorf("资产不存在或状态不正确")
	}
	// 待入库资产直接删除
	err := s.orm.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("资产删除失败: %w", err)
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
		{"城市", "编号"},
	}
	name = "导入电池模版"
	path = filepath.Join("runtime/export", fmt.Sprintf("%s-%s.xlsx", "", "battery"))
	tools.NewExcel(path).AddValues(data).Done()
	return path, name, nil
}

// Export 导出资产
func (s *assetService) Export(ctx context.Context, req *model.AssetListReq, m *model.Modifier) (model.ExportRes, error) {
	q := s.orm.QueryNotDeleted().WithCabinet().WithCity().WithStation().WithModel().WithOperator().WithValues().WithStore().WithWarehouse().WithBrand().WithValues()
	s.filter(q, &req.AssetFilter)
	q.Order(ent.Desc(asset.FieldCreatedAt))

	if req.AssetType == nil {
		return model.ExportRes{}, errors.New("类型不能为空")
	}
	switch *req.AssetType {
	case model.AssetTypeSmartBattery:
		s.batteryFilter(q, &req.AssetFilter)
		return s.exportBattery(ctx, req, q, m), nil
	case model.AssetTypeEbike:
		s.ebikeFilter(q, &req.AssetFilter)
		return s.exportEbike(ctx, req, q, m), nil
	default:
		return model.ExportRes{}, errors.New("未知类型")
	}
}

// 导出电池
func (s *assetService) exportBattery(ctx context.Context, req *model.AssetListReq, q *ent.AssetQuery, m *model.Modifier) model.ExportRes {
	return NewExportWithModifier(m).Start("电池列表", req.AssetFilter, nil, "", func(path string) {
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
					assetLocations += item.Edges.Rider.Name
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

			row := []any{
				cityName,
				belong,
				assetLocations,
				brandName,
				modelStr,
				item.Sn,
				model.AssetStatus(item.Status).String(),
				item.Enable,
				item.Remark,
			}
			rows = append(rows, row)
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}

// 导出电车
func (s *assetService) exportEbike(ctx context.Context, req *model.AssetListReq, q *ent.AssetQuery, m *model.Modifier) model.ExportRes {
	t, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.AssetType((model.AssetTypeEbike).Value())).WithValues().All(ctx)
	return NewExportWithModifier(m).Start("电车列表", req.AssetFilter, nil, "", func(path string) {
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
					assetLocations += item.Edges.Rider.Name
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
	q := s.orm.QueryNotDeleted().WithCabinet().WithCity().WithStation().WithModel().WithOperator().WithValues().WithStore().WithWarehouse().WithBrand().WithValues()
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
		var cityName, belong, assetLocations string
		if item.Edges.City != nil {
			cityName = item.Edges.City.Name
		}
		belong = "平台"
		if item.LocationsType == model.AssetLocationsTypeStation.Value() {
			belong = "代理商"
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
				assetLocations += item.Edges.Rider.Name
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
		if item.Edges.Brand != nil {
			brandName = item.Edges.Brand.Name
		}
		if item.Type == model.AssetTypeNonSmartBattery.Value() || item.Type == model.AssetTypeSmartBattery.Value() {
			brandName = item.BrandName
		}

		res := &model.AssetListRes{
			ID:             item.ID,
			CityName:       cityName,
			Belong:         belong,
			AssetLocations: assetLocations,
			Brand:          brandName,
			Model:          modelStr,
			SN:             item.Sn,
			AssetStatus:    model.AssetStatus(item.Status).String(),
			Enable:         item.Enable,
			Remark:         item.Remark,
		}

		attributeValue, _ := item.QueryValues().WithAttribute().All(ctx)
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
		return res
	})
}

// 资产公共筛选
func (s *assetService) filter(q *ent.AssetQuery, req *model.AssetFilter) {
	if req.SN != nil {
		q.Where(asset.Sn(*req.SN))
	}
	if req.AssetType != nil {
		q.Where(asset.Type(req.AssetType.Value()))
	}
	if req.OwnerType != nil {
		// 平台
		if *req.OwnerType == 1 {
			q.Where(asset.LocationsTypeEQ(model.AssetLocationsTypeWarehouse.Value()))
		}
		// 代理商
		if *req.OwnerType == 2 && req.StationID != nil {
			q.Where(
				asset.LocationsTypeEQ(model.AssetLocationsTypeStation.Value()),
				asset.HasStationWith(enterprisestation.ID(*req.StationID)),
			)
		}
	}
	if req.LocationsType != nil && req.LocationsKeywork != nil {
		q.Where(asset.LocationsType(req.LocationsType.Value()))
		switch *req.LocationsType {
		case model.AssetLocationsTypeWarehouse:
			q.Where(
				asset.HasWarehouseWith(
					warehouse.NameContains(*req.LocationsKeywork),
				),
			)
		case model.AssetLocationsTypeStore:
			q.Where(
				asset.HasStoreWith(
					store.NameContains(*req.LocationsKeywork),
				),
			)
		case model.AssetLocationsTypeCabinet:
			q.Where(
				asset.HasCabinetWith(
					cabinet.NameContains(*req.LocationsKeywork),
				),
			)
		case model.AssetLocationsTypeStation:
			q.Where(
				asset.HasStationWith(
					enterprisestation.NameContains(*req.LocationsKeywork),
				),
			)
		case model.AssetLocationsTypeRider:
			q.Where(
				asset.HasRiderWith(
					rider.NameContains(*req.LocationsKeywork),
				),
			)
		case model.AssetLocationsTypeOperation:
			q.Where(
				asset.HasOperatorWith(
					maintainer.NameContains(*req.LocationsKeywork),
				),
			)
		}
	}
	if req.Status != nil {
		q.Where(asset.Status(*req.Status))
	}
	if req.Enable != nil {
		q.Where(asset.Enable(*req.Enable))
	}

	// 属性查询
	if req.Attribute != nil {
		for _, v := range req.Attribute {
			q.Where(asset.HasValuesWith(assetattributevalues.AttributeID(v.AttributeID), assetattributevalues.ValueContains(v.AttributeValue)))
		}
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
	return q
}

// FlowDetail 电池流转明细
func (s *assetService) FlowDetail(ctx context.Context, id uint64) []*model.AssetFlowDetail {
	return nil
}

// Count 查询有效的资产数量
func (s *assetService) Count(ctx context.Context, req *model.AssetFilter) *model.AssetNumRes {
	q := s.orm.QueryNotDeleted().Where(asset.StatusNEQ(model.AssetStatusScrap.Value()))
	s.filter(q, req)
	count, _ := q.Count(ctx)
	return &model.AssetNumRes{Num: count}
}

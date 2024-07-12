package asset

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/asset"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterynew"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type batteryService struct {
	orm *ent.BatteryNewClient
}

func NewBattery() *batteryService {
	return &batteryService{
		orm: ent.Database.BatteryNew,
	}
}

// Create 创建电池
func (s *batteryService) Create(ctx context.Context, req *model.BatteryCreateReq, modifier *model.Modifier) {
	// 解析电池编号
	ab, err := ParseBatterySN(req.SN)
	if err != nil || ab.Model == "" {
		snag.Panic("电池编号解析失败!")
	}
	_, err = s.orm.Create().
		SetSn(req.SN).
		SetBrand(ab.Brand).
		SetModel(ab.Model).
		SetCityID(req.CityID).
		SetCreator(modifier).
		Save(ctx)
	if err != nil {
		snag.Panic("电池创建失败: " + err.Error())
	}
}

// BatchCreate 批量创建电池
func (s *batteryService) BatchCreate(ctx echo.Context, modifier *model.Modifier) ([]string, error) {
	rows, err := ParseExcel(ctx)
	if err != nil {
		return nil, err
	}

	// 获取电池编号
	var sns []string
	for _, row := range rows {
		sns = append(sns, row[1])
	}

	// 查重
	items, _ := s.orm.Query().Where(batterynew.SnIn(sns...)).All(context.Background())
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
	cities, _ := ent.Database.City.Query().Where(city.NameIn(cids...)).All(ctx.Request().Context())
	for _, ci := range cities {
		cm[ci.Name] = ci.ID
	}

	var failed []string
	for _, row := range rows {
		sn := row[1]
		if m[sn] {
			failed = append(failed, fmt.Sprintf("编号%s已存在", sn))
			continue
		}

		// 解析电池编号
		var ab asset.Battery
		ab, err = ParseBatterySN(sn)
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

		_, err = creator.
			SetModel(ab.Model).
			SetBrand(ab.Brand).
			SetSn(sn).
			SetCreator(modifier).
			Save(ctx.Request().Context())
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s保存失败: %v", sn, err))
		}
	}
	return failed, nil
}

// Modify 修改电池
func (s *batteryService) Modify(ctx context.Context, req *asset.BatteryModifyReq, modifier *model.Modifier) error {
	q := s.orm.UpdateOneID(req.ID).
		SetNillableCityID(req.CityID).
		SetNillableEnable(req.Enable).
		SetNillableRemark(req.Remark).
		SetLastModifier(modifier)

	// 报废电池
	if req.ScrapReasonType != nil {
		err := s.Scrap(ctx, req.ID, *req.ScrapReasonType, modifier)
		if err != nil {
			return err
		}
	}
	err := q.Exec(ctx)
	if err != nil {
		return fmt.Errorf("电池修改失败: %w", err)
	}
	return nil
}

// Scrap 报废电池
func (s *batteryService) Scrap(ctx context.Context, id uint64, scrapReasonType asset.ScrapReasonType, modifier *model.Modifier) error {
	// 判定是否能报废
	bat, _ := s.orm.Query().Where(
		batterynew.ID(id),
		// 电池状态在库存或故障不可报废
		batterynew.AssetStatusIn(
			asset.BatteryAssetStatusStock.Value(),
			asset.BatteryAssetStatusFault.Value(),
		),
		// 资产库存位置在骑手或运维不可报废
		batterynew.AssetLocationsTypeIn(
			asset.BatteryAssetLocationTypeRider.Value(),
			asset.BatteryAssetLocationTypeOperation.Value()),
	).First(ctx)
	if bat == nil {
		return fmt.Errorf("电池不存在或状态不正确")
	}

	err := bat.Update().
		SetScrapReasonType(scrapReasonType.Value()).
		SetEnable(false).
		SetOperateID(modifier.ID).
		SetLastModifier(modifier).
		SetScrapAt(time.Now()).
		// todo 这里操作人角色不明
		SetOperateUser("角色-" + modifier.Name).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("电池报废失败: %w", err)
	}
	return nil
}

// Delete 删除电池
func (s *batteryService) Delete(ctx context.Context, id uint64) error {
	bat, _ := s.orm.QueryNotDeleted().
		Where(
			batterynew.ID(id),
			batterynew.AssetStatus(asset.BatteryAssetStatusPending.Value()),
		).
		First(ctx)
	if bat == nil {
		return fmt.Errorf("电池不存在或状态不正确")
	}
	// 待入库电池直接删除
	err := s.orm.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("电池删除失败: %w", err)
	}
	return nil
}

// DownloadTemplate 批量导入模版下载
func (s *batteryService) DownloadTemplate() (path string, err error) {
	data := [][]any{
		{"城市", "电池编号"},
		{"上海", "TB00000000000001"},
	}
	path = filepath.Join("runtime/export", fmt.Sprintf("%s-%s.xlsx", "导入电池模版", "battery"))
	tools.NewExcel(path).AddValues(data).Done()
	return path, nil
}

// ExportList 导出电池
func (s *batteryService) ExportList(ctx context.Context, req *asset.BatteryListReq) model.ExportRes {
	q := s.orm.Query()
	s.filter(q, &req.BatteryFilter)
	return service.NewExportWithModifier(nil).Start("电池列表", req.BatteryFilter, nil, "", func(path string) {
		items, _ := q.All(ctx)
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
			if item.EnterpriseID != nil {
				belong = "代理商"
			}
			switch item.AssetLocationsType {
			case asset.BatteryAssetLocationTypeWarehouse.Value():
				assetLocations = "仓库"
			case asset.BatteryAssetLocationTypeStore.Value():
				assetLocations = "门店"
			case asset.BatteryAssetLocationTypeCabinet.Value():
				assetLocations = "电柜"
			case asset.BatteryAssetLocationTypeStation.Value():
				assetLocations = "站点"
			case asset.BatteryAssetLocationTypeRider.Value():
				assetLocations = "骑手" + item.AssetLocations
			case asset.BatteryAssetLocationTypeOperation.Value():
				assetLocations = "运维"
			}

			rows = append(rows, []any{
				cityName,
				belong,
				assetLocations,
				item.Brand,
				item.Model,
				item.Sn,
				asset.BatteryAssetStatus(item.AssetStatus).String(),
				item.Enable,
				item.Remark,
			})
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}

// List 电池列表
func (s *batteryService) List(ctx context.Context, req *asset.BatteryListReq) *model.PaginationRes {
	q := s.orm.Query()
	s.filter(q, &req.BatteryFilter)
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.BatteryNew) *asset.BatteryListRes {
		var cityName, belong, assetLocations string
		if item.Edges.City != nil {
			cityName = item.Edges.City.Name
		}
		belong = "平台"
		if item.EnterpriseID != nil {
			belong = "代理商"
		}
		switch item.AssetLocationsType {
		case asset.BatteryAssetLocationTypeWarehouse.Value():
			assetLocations = "仓库"
		case asset.BatteryAssetLocationTypeStore.Value():
			assetLocations = "门店"
		case asset.BatteryAssetLocationTypeCabinet.Value():
			assetLocations = "电柜"
		case asset.BatteryAssetLocationTypeStation.Value():
			assetLocations = "站点"
		case asset.BatteryAssetLocationTypeRider.Value():
			assetLocations = "骑手" + item.AssetLocations
		case asset.BatteryAssetLocationTypeOperation.Value():
			assetLocations = "运维"
		}

		return &asset.BatteryListRes{
			ID:             item.ID,
			CityName:       cityName,
			Belong:         belong,
			AssetLocations: assetLocations,
			Brand:          item.Brand,
			Model:          item.Model,
			SN:             item.Sn,
			AssetStatus:    asset.BatteryAssetStatus(item.AssetStatus).String(),
			Enable:         item.Enable,
			Remark:         item.Remark,
		}
	})
}

// 筛选条件
func (s *batteryService) filter(q *ent.BatteryNewQuery, req *asset.BatteryFilter) {
	if req.SN != nil {
		q.Where(batterynew.Sn(*req.SN))
	}
	if req.Model != nil {
		q.Where(batterynew.Model(*req.Model))
	}
	if req.CityID != nil {
		q.Where(batterynew.CityID(*req.CityID))
	}
	if req.OwnerType != nil {
		// 平台
		if *req.OwnerType == 1 {
			q.Where(batterynew.EnterpriseIDIsNil())
		}
		// 代理商
		if *req.OwnerType == 2 {
			q.Where(batterynew.EnterpriseIDNotNil())
		}
	}
	if req.AssetLocationsType != nil && *req.AssetLocationsKeywork {
		q.Where(batterynew.AssetLocationsType(*req.AssetLocationsType))
		switch *req.AssetLocationsType {
		case asset.BatteryAssetLocationTypeWarehouse.Value():
			// 仓库
		case asset.BatteryAssetLocationTypeStore.Value():
			// 门店
		case asset.BatteryAssetLocationTypeCabinet.Value():
			// 电柜
		case asset.BatteryAssetLocationTypeStation.Value():
			// 站点
		case asset.BatteryAssetLocationTypeRider.Value():
			// 骑手
		case asset.BatteryAssetLocationTypeOperation.Value():
			// 运维
		}
	}
	if req.AssetStatus != nil {
		q.Where(batterynew.AssetStatus(*req.AssetStatus))
	}
	if req.Enable != nil {
		q.Where(batterynew.Enable(*req.Enable))
	}
}

// FlowDetail 电池流转明细
func (s *batteryService) FlowDetail(ctx context.Context, id uint64) []*asset.BatteryFlowDetail {
	return nil
}

// ScrapList 电池报废列表
func (s *batteryService) ScrapList(ctx context.Context, req *asset.BatteryScrapListReq) *model.PaginationRes {
	q := s.orm.Query().Where(
		batterynew.ScrapReasonTypeNotNil(),
		batterynew.DeletedAtIsNil(),
	)
	if req.SN != nil {
		q.Where(batterynew.Sn(*req.SN))
	}
	if req.Model != nil {
		q.Where(batterynew.Model(*req.Model))
	}
	if req.ScrapReasonType != nil {
		q.Where(batterynew.ScrapReasonType((*req.ScrapReasonType).Value()))
	}
	// 报废时间
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			batterynew.ScrapAt(start),
			batterynew.ScrapAt(end),
		)
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.BatteryNew) *asset.BatteryScrapListRes {
		return &asset.BatteryScrapListRes{
			ID: item.ID,
			SN: item.Sn,
			// todo 这里要关联物资的电池型号
			Model:       item.Model,
			Brand:       item.Brand,
			ScrapReason: asset.ScrapReasonType(item.ScrapReasonType),
			Operate:     item.OperateUser,
			Remark:      item.Remark,
			CreatedAt:   item.CreatedAt.Format(carbon.DateTimeLayout),
			ScrapAt:     item.ScrapAt.Format(carbon.DateTimeLayout),
		}
	})
}

// ScrapBatchRestore 报废批量还原
func (s *batteryService) ScrapBatchRestore(ctx context.Context, ids []uint64, modifier *model.Modifier) error {
	// 判定是否可还原 只有报废状态能还原
	items, _ := s.orm.Query().Where(batterynew.AssetStatus(asset.BatteryAssetStatusScrap.Value()), batterynew.IDIn(ids...)).All(ctx)
	if len(items) != len(ids) {
		return fmt.Errorf("电池不存在或状态不正确")
	}
	for _, item := range items {
		err := item.Update().
			SetEnable(true).
			SetLastModifier(modifier).
			SetNillableScrapAt(nil).
			SetNillableOperateUser(nil).
			SetNillableOperateID(nil).
			SetNillableScrapReasonType(nil).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("电池还原失败: %w", err)
		}
	}
	return nil
}

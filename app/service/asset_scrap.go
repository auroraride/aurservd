package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/internal/ent/assetscrap"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assetScrapService struct {
	orm *ent.AssetScrapClient
}

func NewAssetScrap() *assetScrapService {
	return &assetScrapService{
		orm: ent.Database.AssetScrap,
	}
}

// ScrapList 报废列表
func (s *assetScrapService) ScrapList(ctx context.Context, req *model.AssetScrapListReq) *model.PaginationRes {
	q := s.orm.Query().WithAsset(func(query *ent.AssetQuery) {
		query.WithBrand().WithModel()
	}).WithAgent().WithEmployee().WithMaintainer().WithManager()
	// 公共筛选条件
	s.filter(ctx, q, req)
	q.Order(ent.Desc(assetscrap.FieldCreatedAt))

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetScrap) (res *model.AssetScrapListRes) {
		var modelStr string
		if item.Edges.Asset != nil && item.Edges.Asset.Edges.Model != nil {
			modelStr = item.Edges.Asset.Edges.Model.Model
		}
		var brandName, sn string
		if item.Edges.Asset != nil {
			if item.Edges.Asset.Edges.Brand != nil {
				brandName = item.Edges.Asset.Edges.Brand.Name
			}
			if item.Edges.Asset.Type == model.AssetTypeNonSmartBattery.Value() || item.Edges.Asset.Type == model.AssetTypeSmartBattery.Value() {
				brandName = item.Edges.Asset.BrandName
			}
			sn = item.Edges.Asset.Sn
		}
		operateName := ""
		if item.ScrapOperateID != nil && item.ScrapOperateRoleType != nil {
			switch *item.ScrapOperateRoleType {
			case model.AssetOperateRoleAdmin.Value():
				if item.Edges.Manager != nil {
					operateName = "[后台]-" + item.Edges.Manager.Name
				}
			case model.AssetOperateRoleStore.Value():
				if item.Edges.Employee != nil {
					operateName = "[门店管理员]-" + item.Edges.Employee.Name
				}
			case model.AssetOperateRoleOperation.Value():
				if item.Edges.Maintainer != nil {
					operateName = "[运维]-" + item.Edges.Maintainer.Name
				}
			case model.AssetOperateRoleAgent.Value():
				if item.Edges.Agent != nil {
					operateName = "[代理管理员]-" + item.Edges.Agent.Name
				}
			default:
				operateName = "未知"
			}
		}

		inTimeAt := ""
		if item.Edges.Asset != nil {
			inTimeAt = item.Edges.Asset.CreatedAt.Format(carbon.DateTimeLayout)
		}

		res = &model.AssetScrapListRes{
			ID:          item.ID,
			SN:          sn,
			Model:       modelStr,
			Brand:       brandName,
			ScrapReason: model.ScrapReasonType(item.ScrapReasonType).String(),
			OperateName: operateName,
			Remark:      item.Remark,
			CreatedAt:   item.CreatedAt.Format(carbon.DateTimeLayout),
			ScrapAt:     item.ScrapAt.Format(carbon.DateTimeLayout),
			InTimeAt:    inTimeAt,
		}

		if item.Edges.Asset != nil {
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
		return res
	})
}

// ScrapListOther 其它物资报废列表
func (s *assetScrapService) ScrapListOther(ctx context.Context, req *model.AssetScrapListReq) *model.PaginationRes {
	// q := s.orm.Query().WithAsset(func(query *ent.AssetQuery) {
	// 	query.WithBrand().WithModel().WithMaterial()
	// }).WithAgent().WithEmployee().WithMaintainer().WithManager()
	q := ent.Database.Asset.QueryNotDeleted().WithMaterial().WithScrap(func(query *ent.AssetScrapQuery) {
		query.WithManager().WithMaintainer().WithEmployee().WithAgent()
	}).WithValues()
	q.Order(ent.Desc(assetscrap.FieldCreatedAt))

	// q.Modify(func(s *sql.Selector) {
	// 	// 关联物资表 根据material_id分组统计数量
	// 	s.Join("material", "material.id = asset.material_id")
	//
	//
	// })
	if req.OperateName != nil {
		q.Where(
			asset.HasScrapWith(
				assetscrap.Or(
					// 门店管理员
					assetscrap.HasEmployeeWith(employee.NameContains(*req.OperateName)),
					// 代理
					assetscrap.HasAgentWith(agent.NameContains(*req.OperateName)),
					// 运维
					assetscrap.HasMaintainerWith(maintainer.NameContains(*req.OperateName)),
					// 后台
					assetscrap.HasManagerWith(manager.NameContains(*req.OperateName)),
				),
			),
		)
	}
	if req.ScrapReasonType != nil {
		q.Where(asset.HasScrapWith(assetscrap.ScrapReasonType((*req.ScrapReasonType).Value())))
	}

	q.Modify(func(sel *sql.Selector) {
		t := sql.Table(assetscrap.Table)
		sel.LeftJoin(t).On(t.C(assetscrap.FieldAssetID), sel.C(asset.FieldID))
		// if req.ScrapReasonType != nil {
		// 	sel.Where(sql.EQ(t.C(assetscrap.FieldScrapReasonType), req.ScrapReasonType))
		// }

		sel.GroupBy(asset.MaterialColumn).
			AppendSelectExprAs(sql.Raw(fmt.Sprintf("COUNT(%s)", sel.C(asset.FieldID))), "num")
	})
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Asset) (res *model.AssetScrapListRes) {
		value, err := item.Value("num")
		if err != nil {
			return nil
		}
		fmt.Println(value)

		return nil
	})

}

// 公共筛选条件
func (s *assetScrapService) filter(ctx context.Context, q *ent.AssetScrapQuery, req *model.AssetScrapListReq) {
	if req.AssetType != nil {
		q.Where(assetscrap.HasAssetWith(asset.Type(req.AssetType.Value())))
	}
	if req.SN != nil {
		q.Where(assetscrap.HasAssetWith(asset.Sn(*req.SN)))
	}
	if req.ModelID != nil {
		q.Where(assetscrap.HasAssetWith(asset.ModelID(*req.ModelID)))
	}
	if req.ScrapReasonType != nil {
		q.Where(assetscrap.ScrapReasonType((*req.ScrapReasonType).Value()))
	}
	// 报废时间
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			assetscrap.ScrapAt(start),
			assetscrap.ScrapAt(end),
		)
	}
	// 操作人
	if req.OperateName != nil {
		q.Where(
			assetscrap.Or(
				// 门店管理员
				assetscrap.HasEmployeeWith(employee.NameContains(*req.OperateName)),
				// 代理
				assetscrap.HasAgentWith(agent.NameContains(*req.OperateName)),
				// 运维
				assetscrap.HasMaintainerWith(maintainer.NameContains(*req.OperateName)),
				// 后台
				assetscrap.HasManagerWith(manager.NameContains(*req.OperateName)),
			),
		)
	}
	// 属性查询
	if req.Attribute != nil {
		for _, v := range req.Attribute {
			q.Where(
				assetscrap.HasAssetWith(
					asset.HasValuesWith(
						assetattributevalues.AttributeID(v.AttributeID),
						assetattributevalues.ValueContains(v.AttributeValue),
					),
				),
			)
		}
	}
	if req.AssetName != nil && req.AssetType != nil {
		switch *req.AssetType {
		case model.AssetTypeEbike, model.AssetTypeSmartBattery:
			q.Where(assetscrap.HasAssetWith(asset.NameContains(*req.AssetName)))
		case model.AssetTypeNonSmartBattery, model.AssetTypeEbikeAccessory, model.AssetTypeCabinetAccessory, model.AssetTypeOtherAccessory:
			q.Where(assetscrap.HasAssetWith(asset.HasMaterialWith(material.NameContains(*req.AssetName))))
		}
	}
}

// ScrapBatchRestore 报废批量还原
func (s *assetScrapService) ScrapBatchRestore(ctx context.Context, req *model.AssetScrapBatchRestoreReq, modifier *model.Modifier) error {
	// 判定是否可还原 只有报废状态能还原
	items, _ := ent.Database.Asset.QueryNotDeleted().Where(asset.Status(model.AssetStatusScrap.Value()), asset.IDIn(req.IDs...)).All(ctx)
	if len(items) != len(req.IDs) {
		return fmt.Errorf("电池不存在或状态不正确")
	}
	for _, item := range items {
		err := item.Update().
			SetStatus(model.AssetStatusStock.Value()).
			SetEnable(true).
			SetLastModifier(modifier).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("电池还原失败: %w", err)
		}
		_, err = ent.Database.AssetScrap.Delete().Where(assetscrap.AssetID(item.ID)).Exec(ctx)
		if err != nil {
			return fmt.Errorf("电池还原失败: %w", err)
		}
	}
	return nil
}

// Scrap 报废资产
func (s *assetScrapService) Scrap(ctx context.Context, req *model.AssetScrapReq, modifier *model.Modifier) error {
	q := ent.Database.Asset.QueryNotDeleted().Where(
		// 资产状态在库存或故障可报废
		asset.StatusIn(
			model.AssetStatusStock.Value(),
			model.AssetStatusFault.Value(),
		),
		// 资产库存位置在骑手或运维不可报废
		asset.LocationsTypeNotIn(
			model.AssetLocationsTypeRider.Value(),
			model.AssetLocationsTypeOperation.Value(),
		),
	)
	var ids []uint64
	switch req.AssetType {
	case model.AssetTypeEbike, model.AssetTypeSmartBattery:
		assetId, err := s.transferAssetWithSN(ctx, q, req)
		if err != nil {
			return err
		}
		ids = append(ids, assetId...)
	case model.AssetTypeOtherAccessory, model.AssetTypeCabinetAccessory, model.AssetTypeNonSmartBattery, model.AssetTypeEbikeAccessory:
		assetId, err := s.transferAssetWithoutSN(ctx, q, req)
		if err != nil {
			return err
		}
		ids = append(ids, assetId...)
	default:
		return errors.New("资产类型错误")
	}

	scrapBluk := make([]*ent.AssetScrapCreate, 0)
	for _, v := range ids {
		// 更新资产状态
		err := ent.Database.Asset.Update().Where(asset.ID(v)).SetStatus(model.AssetStatusScrap.Value()).Exec(ctx)
		if err != nil {
			return err
		}

		// 是否已经报废
		if b, _ := s.orm.Query().Where(assetscrap.AssetID(v)).Exist(ctx); b {
			return fmt.Errorf("资产已经报废")
		}

		// 创建报废记录
		scrapBluk = append(scrapBluk, s.orm.Create().
			SetScrapBatch(tools.NewUnique().NewSN28()).
			SetAssetID(v).
			SetScrapReasonType(req.ScrapReasonType.Value()).
			SetLastModifier(modifier).
			SetScrapAt(time.Now()).
			SetScrapOperateID(modifier.ID).
			SetScrapOperateRoleType(model.AssetOperateRoleAdmin.Value()).
			SetNillableRemark(req.Remark).
			SetCreator(modifier).
			SetLastModifier(modifier),
		)
	}
	_, err := s.orm.CreateBulk(scrapBluk...).Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

// 有编号资产报废
func (s *assetScrapService) transferAssetWithSN(ctx context.Context, q *ent.AssetQuery, req *model.AssetScrapReq) ([]uint64, error) {
	ids := make([]uint64, 0)
	for _, v := range req.Detail {
		if v.AssetID == nil {
			return nil, fmt.Errorf("资产ID不能为空")
		}
		bat, _ := q.Where(asset.Type(req.AssetType.Value()), asset.ID(*v.AssetID)).All(ctx)
		if bat == nil {
			return nil, fmt.Errorf("资产不存在或状态不正确")
		}
		for _, vl := range bat {
			ids = append(ids, vl.ID)
		}
	}
	return ids, nil
}

// 无编号资产报废
func (s *assetScrapService) transferAssetWithoutSN(ctx context.Context, q *ent.AssetQuery, req *model.AssetScrapReq) ([]uint64, error) {
	ids := make([]uint64, 0)
	for _, v := range req.Detail {
		if v.MaterialID == nil {
			return nil, fmt.Errorf("物资分类ID不能为空")
		}
		if v.Num == nil || *v.Num == 0 {
			return nil, fmt.Errorf("报废数量不能为空")
		}
		all, _ := q.Where(asset.Type(req.AssetType.Value()), asset.MaterialID(*v.MaterialID)).Limit(int(*v.Num)).All(ctx)
		if len(all) < int(*v.Num) {
			return nil, fmt.Errorf("物资数量不足")
		}
		for _, vl := range all {
			ids = append(ids, vl.ID)
		}
	}
	return ids, nil
}

// ScrapReasonSelect 报废原因下拉
func (s *assetScrapService) ScrapReasonSelect(ctx context.Context) (res []*model.SelectOption) {
	for k, v := range model.ScrapReasonTypeMap {
		res = append(res, &model.SelectOption{
			Value: uint64(k),
			Label: v,
		})
	}
	return res
}

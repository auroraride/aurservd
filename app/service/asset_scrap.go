package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetscrap"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/manager"
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
	q := s.orm.Query().Where(
		assetscrap.ScrapReasonTypeNotNil(),
	).WithAsset(func(query *ent.AssetQuery) {
		query.WithBrand().WithModel()
	}).WithAgent().WithEmployee().WithMaintainer().WithManager()
	if req.AssetType != nil {
		q.Where(assetscrap.HasAssetWith(asset.Type((*req.AssetType).Value())))
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
				// todo 物资管理员
				// assetscrap.HasEmployeeWith(employee.NameContains(*req.OperateName)),
			),
		)
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetScrap) *model.AssetScrapListRes {

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
			case model.ScrapOperateRoleAdmin.Value():
				if item.Edges.Manager != nil {
					operateName = "[后台]-" + item.Edges.Manager.Name
				}
			case model.ScrapOperateRoleStore.Value():
				if item.Edges.Employee != nil {
					operateName = "[门店管理员]-" + item.Edges.Employee.Name
				}
			case model.ScrapOperateRoleOperation.Value():
				if item.Edges.Maintainer != nil {
					operateName = "[运维]-" + item.Edges.Maintainer.Name
				}
			case model.ScrapOperateRoleAgent.Value():
				if item.Edges.Agent != nil {
					operateName = "[代理管理员]-" + item.Edges.Agent.Name
				}
			case model.ScrapOperateRoleMaterial.Value():
				// todo
				operateName = "[物资管理员]"
			default:
				operateName = "未知"
			}
		}

		return &model.AssetScrapListRes{
			ID:          item.ID,
			SN:          sn,
			Model:       modelStr,
			Brand:       brandName,
			ScrapReason: model.ScrapReasonType(item.ScrapReasonType).String(),
			OperateName: operateName,
			Remark:      item.Remark,
			CreatedAt:   item.CreatedAt.Format(carbon.DateTimeLayout),
			ScrapAt:     item.ScrapAt.Format(carbon.DateTimeLayout),
		}
	})
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
	// 判定是否能报废
	bat, _ := ent.Database.Asset.QueryNotDeleted().Where(
		asset.ID(req.ID),
		// 资产状态在库存或故障可报废
		asset.StatusIn(
			model.AssetStatusStock.Value(),
			model.AssetStatusFault.Value(),
		),
		// 资产库存位置在骑手或运维不可报废
		asset.LocationsTypeNotIn(
			model.AseetLocationsTypeRider.Value(),
			model.AseetLocationsTypeOperation.Value(),
		),
	).First(ctx)
	if bat == nil {
		return fmt.Errorf("资产不存在或状态不正确")
	}

	// 更新资产状态
	err := bat.Update().SetStatus(model.AssetStatusScrap.Value()).Exec(ctx)
	if err != nil {
		return err
	}

	// 是否已经报废
	if b, _ := s.orm.Query().Where(assetscrap.AssetID(bat.ID)).Exist(ctx); b {
		return fmt.Errorf("资产已经报废")
	}

	// 创建报废记录
	err = s.orm.Create().
		SetAssetID(bat.ID).
		SetScrapReasonType(req.ScrapReasonType.Value()).
		SetLastModifier(modifier).
		SetScrapAt(time.Now()).
		SetScrapOperateID(modifier.ID).
		SetScrapOperateRoleType(model.ScrapOperateRoleAdmin.Value()).
		SetNillableRemark(req.Remark).
		SetCreator(modifier).
		SetLastModifier(modifier).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("资产报废失败: %w", err)
	}
	return nil
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

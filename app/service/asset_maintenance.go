package service

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/pkg/errors"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetmaintenance"
	"github.com/auroraride/aurservd/internal/ent/assetmaintenancedetails"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assetMaintenanceService struct {
	orm *ent.AssetMaintenanceClient
}

func NewAssetMaintenance() *assetMaintenanceService {
	return &assetMaintenanceService{
		orm: ent.Database.AssetMaintenance,
	}
}

// Create 创建资产维护
func (s *assetMaintenanceService) Create(ctx context.Context, req *model.AssetMaintenanceCreateReq, modifier *model.Modifier) error {
	err := ent.Database.AssetMaintenance.Create().
		SetMaintainerID(modifier.ID).
		SetLastModifier(modifier).
		SetCreator(modifier).
		SetCabinetID(req.CabinetID).
		SetStatus(model.AssetMaintenanceStatusUnder.Value()).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Modify 修改维护记录
func (s *assetMaintenanceService) Modify(ctx context.Context, req *model.AssetMaintenanceModifyReq, modifier *model.Modifier) error {
	bulk := make([]*ent.AssetMaintenanceDetailsCreate, 0)
	for _, v := range req.Details {
		item, _ := ent.Database.Material.QueryNotDeleted().Where(material.ID(v.MaterialID)).First(ctx)
		if item == nil {
			return errors.New("分类不存在")
		}
		// all, _ := ent.Database.Asset.QueryNotDeleted().
		// 	Where(
		// 		asset.StatusIn(model.AssetStatusStock.Value()),
		// 		asset.MaterialID(v.MaterialID),
		// 	).
		// 	Limit(int(v.Num)).
		// 	Order(ent.Asc(asset.FieldCreatedAt)).All(ctx)
		// if len(all) < int(v.Num) {
		// 	return fmt.Errorf(strconv.FormatUint(v.MaterialID, 10) + "物资数量不足")
		// }
		//
		// for _, vl := range all {
		// bulk = append(bulk, ent.Database.AssetMaintenanceDetails.Create().SetAssetID(vl.ID))
		// }
		for i := 0; i < int(v.Num); i++ {
			bulk = append(bulk, ent.Database.AssetMaintenanceDetails.
				Create().
				SetMaterialID(v.MaterialID).
				SetLastModifier(modifier).
				SetCreator(modifier),
			)
		}
	}
	md, err := ent.Database.AssetMaintenanceDetails.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return err
	}
	err = ent.Database.AssetMaintenance.Update().Where(assetmaintenance.ID(req.ID)).
		SetReason(req.Reason).
		SetStatus(req.Status.Value()).
		SetMaintainerID(modifier.ID).
		SetLastModifier(modifier).
		SetContent(req.Content).
		AddMaintenanceDetails(md...).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// List 列表
func (s *assetMaintenanceService) List(ctx context.Context, req *model.AssetMaintenanceListReq) (res *model.PaginationRes) {
	q := ent.Database.AssetMaintenance.
		QueryNotDeleted().
		Order(ent.Desc(asset.FieldCreatedAt)).
		WithMaintainer().
		WithCabinet().
		WithMaintenanceDetails()
	if req.Keyword != nil {
		q.Where(
			assetmaintenance.Or(
				assetmaintenance.HasMaintainerWith(maintainer.NameContains(*req.Keyword)),
				assetmaintenance.HasMaintainerWith(maintainer.PhoneContains(*req.Keyword)),
				assetmaintenance.HasCabinetWith(cabinet.SnContains(*req.Keyword)),
			),
		)
	}
	if req.Status != nil {
		q.Where(assetmaintenance.Status(*req.Status))
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			assetmaintenance.CreatedAtGTE(start),
			assetmaintenance.CreatedAtLTE(end),
		)
	}

	q.Modify(func(sel *sql.Selector) {
		if req.IsUseAccessory != nil {
			if *req.IsUseAccessory {
				sel.Where(
					sql.Exists(
						sql.Select(assetmaintenancedetails.FieldID).
							From(sql.Table("asset_maintenance_details")).
							Where(sql.ColumnsEQ(sql.Table("asset_maintenance_details").C("maintenance_id"), sql.Table("asset_maintenance").C("id"))),
					),
				)
			} else {
				sel.Where(
					sql.NotExists(
						sql.Select(assetmaintenancedetails.FieldID).
							From(sql.Table("asset_maintenance_details")).
							Where(sql.ColumnsEQ(sql.Table("asset_maintenance_details").C("maintenance_id"), sql.Table("asset_maintenance").C("id"))),
					),
				)
			}
		}
	})
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetMaintenance) (res *model.AssetMaintenanceListRes) {
		var result []struct {
			Count      float64 `json:"count"`
			MaterialID uint64  `json:"material_id"`
		}
		res = &model.AssetMaintenanceListRes{
			ID:        item.ID,
			Reason:    item.Reason,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			Status:    model.AssetMaintenanceStatus(item.Status).String(),
			Content:   item.Content,
			Details:   make([]model.AssetMaintenanceDetail, 0),
		}
		if item.Edges.Maintainer != nil {
			res.OpratorName = item.Edges.Maintainer.Name
			res.OpratorPhone = item.Edges.Maintainer.Phone
		}
		if item.Edges.Cabinet != nil {
			res.CabinetName = item.Edges.Cabinet.Name
			res.CabinetSn = item.Edges.Cabinet.Sn
		}
		err := ent.Database.AssetMaintenanceDetails.QueryNotDeleted().Select(assetmaintenancedetails.FieldMaterialID).Where(assetmaintenancedetails.MaintenanceID(item.ID)).GroupBy(assetmaintenancedetails.FieldMaterialID).Aggregate(ent.Count()).
			Scan(context.Background(), &result)
		if err != nil {
			return
		}
		for _, v := range result {
			m, _ := ent.Database.Material.Query().Where(material.ID(v.MaterialID)).First(ctx)
			if m == nil {
				continue
			}
			res.Details = append(res.Details, model.AssetMaintenanceDetail{
				AssetName: m.Name,
				Num:       uint8(v.Count),
				AssetType: model.AssetType(m.Type).String(),
			})
		}
		return res
	})
}

// QueryMaintenanceByCabinetID 通过电柜ID查询维护中的数据
func (s *assetMaintenanceService) QueryMaintenanceByCabinetID(cabId uint64) *ent.AssetMaintenance {
	res, _ := s.orm.QueryNotDeleted().
		Where(
			assetmaintenance.CabinetID(cabId),
			assetmaintenance.Status(model.AssetMaintenanceStatusUnder.Value()),
		).First(context.Background())
	return res
}

// QueryByID 通过电柜ID查询维保
func (s *assetMaintenanceService) QueryByID(cabId uint64) (res model.AssetMaintenanceRes) {
	mt, _ := s.orm.QueryNotDeleted().Where(assetmaintenance.CabinetID(cabId)).First(context.Background())
	if mt == nil {
		return
	}

	return model.AssetMaintenanceRes{
		ID:     mt.ID,
		Status: model.AssetMaintenanceStatus(mt.Status),
	}
}

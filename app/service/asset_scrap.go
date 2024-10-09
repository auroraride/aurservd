package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assetscrap"
	"github.com/auroraride/aurservd/internal/ent/assetscrapdetails"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
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
	q := s.orm.Query().WithScrapDetails(func(query *ent.AssetScrapDetailsQuery) {
		query.WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		})
	}).WithAgent().WithEmployee().WithMaintainer().WithManager()
	// 公共筛选条件
	s.filter(ctx, q, req)
	q.Order(ent.Desc(assetscrap.FieldCreatedAt))
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetScrap) (res *model.AssetScrapListRes) {
		operateName := ""
		if item.OperateID != nil && item.OperateRoleType != nil {
			switch *item.OperateRoleType {
			case model.OperatorTypeAssetManager.Value():
				if item.Edges.Manager != nil {
					// 查询角色
					var roleName string
					if role, _ := item.Edges.Manager.QueryRole().First(ctx); role != nil {
						roleName = role.Name
					}
					operateName = "[" + roleName + "]-" + item.Edges.Manager.Name
				}
			case model.OperatorTypeEmployee.Value():
				if item.Edges.Employee != nil {
					operateName = "[门店]-" + item.Edges.Employee.Name
				}
			case model.OperatorTypeMaintainer.Value():
				if item.Edges.Maintainer != nil {
					operateName = "[运维]-" + item.Edges.Maintainer.Name
				}
			case model.OperatorTypeAgent.Value():
				if item.Edges.Agent != nil {
					operateName = "[代理]-" + item.Edges.Agent.Name
				}
			default:
				operateName = "未知"
			}
		}
		res = &model.AssetScrapListRes{
			ID:          item.ID,
			ScrapReason: model.ScrapReasonType(item.ReasonType).String(),
			OperateName: operateName,
			Remark:      item.Remark,
			ScrapAt:     item.ScrapAt.Format(carbon.DateTimeLayout),
		}

		if len(item.Edges.ScrapDetails) != 0 {
			if item.Edges.ScrapDetails[0].Edges.Asset != nil {
				v := item.Edges.ScrapDetails[0]
				switch model.AssetType(item.Edges.ScrapDetails[0].Edges.Asset.Type) {
				case model.AssetTypeEbike, model.AssetTypeSmartBattery:
					res.Num = 1
				case model.AssetTypeNonSmartBattery, model.AssetTypeEbikeAccessory, model.AssetTypeCabinetAccessory, model.AssetTypeOtherAccessory:
					res.Num = uint(len(item.Edges.ScrapDetails))
				}
				attributeValue, _ := v.Edges.Asset.QueryValues().WithAttribute().All(ctx)
				assetAttributeMap := make(map[uint64]model.AssetAttribute)

				for _, vl := range attributeValue {
					var attributeName, attributeKey string
					if vl.Edges.Attribute != nil {
						attributeName = vl.Edges.Attribute.Name
						attributeKey = vl.Edges.Attribute.Key
					}
					assetAttributeMap[vl.AttributeID] = model.AssetAttribute{
						AttributeID:      vl.AttributeID,
						AttributeValue:   vl.Value,
						AttributeName:    attributeName,
						AttributeKey:     attributeKey,
						AttributeValueID: vl.ID,
					}
				}
				var modelStr string
				if v.Edges.Asset != nil && v.Edges.Asset.Edges.Model != nil {
					modelStr = v.Edges.Asset.Edges.Model.Model
				}
				var brandName, sn string
				if v.Edges.Asset != nil {
					if v.Edges.Asset.Edges.Brand != nil {
						brandName = v.Edges.Asset.Edges.Brand.Name
					}
					if v.Edges.Asset.Type == model.AssetTypeNonSmartBattery.Value() || v.Edges.Asset.Type == model.AssetTypeSmartBattery.Value() {
						brandName = v.Edges.Asset.BrandName
					}
					sn = v.Edges.Asset.Sn
				}

				inTimeAt := ""
				if v.Edges.Asset != nil {
					inTimeAt = v.Edges.Asset.CreatedAt.Format(carbon.DateTimeLayout)
				}

				// 物资名称
				var name string
				if v.Edges.Asset.Type == model.AssetTypeEbike.Value() || v.Edges.Asset.Type == model.AssetTypeSmartBattery.Value() {
					name = v.Edges.Asset.Name
				} else {
					if v.Edges.Asset.Edges.Material != nil {
						name = v.Edges.Asset.Edges.Material.Name
					}
				}
				res.AssetID = v.AssetID
				res.SN = sn
				res.Model = modelStr
				res.Brand = brandName
				res.InTimeAt = inTimeAt
				res.Attribute = assetAttributeMap
				res.Name = name
				res.AssetType = model.AssetType(v.Edges.Asset.Type).String()
			}
		}
		return res
	})
}

// 公共筛选条件
func (s *assetScrapService) filter(ctx context.Context, q *ent.AssetScrapQuery, req *model.AssetScrapListReq) {
	if req.ScrapType != nil {
		switch *req.ScrapType {
		case model.ScrapTypeEbike:
			q.Where(assetscrap.HasScrapDetailsWith(assetscrapdetails.HasAssetWith(asset.Type(model.AssetTypeEbike.Value()))))
		case model.ScrapTypeSmartBattery:
			q.Where(assetscrap.HasScrapDetailsWith(assetscrapdetails.HasAssetWith(asset.Type(model.AssetTypeSmartBattery.Value()))))
		case model.ScrapTypeOther:
			q.Where(assetscrap.HasScrapDetailsWith(assetscrapdetails.HasAssetWith(asset.TypeIn(
				model.AssetTypeNonSmartBattery.Value(),
				model.AssetTypeCabinetAccessory.Value(),
				model.AssetTypeEbikeAccessory.Value(),
				model.AssetTypeOtherAccessory.Value(),
			))))
		}
	}
	if req.SN != nil {
		q.Where(assetscrap.HasScrapDetailsWith(assetscrapdetails.HasAssetWith(asset.SnContains(*req.SN))))
	}
	if req.ModelID != nil {
		q.Where(assetscrap.HasScrapDetailsWith(assetscrapdetails.HasAssetWith(asset.ModelID(*req.ModelID))))
	}
	if req.BrandID != nil {
		q.Where(assetscrap.HasScrapDetailsWith(assetscrapdetails.HasAssetWith(asset.BrandID(*req.BrandID))))
	}
	if req.ScrapReasonType != nil {
		q.Where(assetscrap.ReasonType((*req.ScrapReasonType).Value()))
	}
	// 报废时间
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			assetscrap.ScrapAtGTE(start),
			assetscrap.ScrapAtLT(end),
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
				assetscrap.HasManagerWith(assetmanager.NameContains(*req.OperateName)),
			),
		)
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
			q.Where(
				assetscrap.HasScrapDetailsWith(
					assetscrapdetails.HasAssetWith(
						asset.HasValuesWith(
							assetattributevalues.AttributeID(attributeID), assetattributevalues.ValueContains(attributeValue),
						),
					),
				),
			)
		}
	}

	if req.AssetType != nil {
		q.Where(assetscrap.HasScrapDetailsWith(assetscrapdetails.HasAssetWith(asset.Type(req.AssetType.Value()))))
	}

	if req.AssetName != nil {
		q.Where(
			assetscrap.HasScrapDetailsWith(
				assetscrapdetails.HasAssetWith(
					asset.Or(
						asset.NameContains(*req.AssetName),
						asset.HasMaterialWith(material.NameContains(*req.AssetName)),
					),
				),
			),
		)
	}
}

// ScrapBatchRestore 报废批量还原
func (s *assetScrapService) ScrapBatchRestore(ctx context.Context, req *model.AssetScrapBatchRestoreReq, modifier *model.Modifier) (failed []string) {
	items, _ := ent.Database.AssetScrap.Query().Where(assetscrap.IDIn(req.IDs...)).WithScrapDetails(func(query *ent.AssetScrapDetailsQuery) {
		query.WithAsset()
	}).All(ctx)
	if len(items) == 0 {
		return []string{"报废记录不存在"}
	}
	for _, item := range items {
		for _, v := range item.Edges.ScrapDetails {
			if v.Edges.Asset == nil {
				failed = append(failed, "资产"+strconv.FormatUint(v.AssetID, 10)+"不存在")
				continue
			}
			if v.Edges.Asset.Status != model.AssetStatusScrap.Value() {
				failed = append(failed, "资产"+strconv.FormatUint(v.AssetID, 10)+"状态不正确")
				continue
			}
			// 更新资产状态
			err := v.Edges.Asset.Update().SetStatus(model.AssetStatusStock.Value()).Exec(ctx)
			if err != nil {
				failed = append(failed, "资产"+strconv.FormatUint(v.AssetID, 10)+"还原失败")
				continue
			}
			// 删除报废详情记录
			_, err = ent.Database.AssetScrapDetails.Delete().Where(assetscrapdetails.AssetID(v.AssetID)).Exec(ctx)
			if err != nil {
				failed = append(failed, "资产"+strconv.FormatUint(v.AssetID, 10)+"还原失败")
				continue
			}
		}
		// 删除报废记录
		_, err := ent.Database.AssetScrap.Delete().Where(assetscrap.ID(item.ID)).Exec(ctx)
		if err != nil {
			failed = append(failed, "报废记录"+strconv.FormatUint(item.ID, 10)+"删除失败")
			continue
		}
	}
	return failed
}

// Scrap 报废资产
func (s *assetScrapService) Scrap(ctx context.Context, req *model.AssetScrapReq, modifier *model.Modifier) error {
	for _, v := range req.Details {
		switch v.AssetType {
		case model.AssetTypeEbike, model.AssetTypeSmartBattery:
			assetId, err := s.scrapAssetWithSN(ctx, &v)
			if err != nil {
				return err
			}
			err = s.createScrap(ctx, req, modifier, assetId)
			if err != nil {
				return err
			}
		case model.AssetTypeOtherAccessory, model.AssetTypeCabinetAccessory, model.AssetTypeNonSmartBattery, model.AssetTypeEbikeAccessory:
			assetId, err := s.scrapAssetWithoutSN(ctx, &v, req)
			if err != nil {
				return err
			}
			err = s.createScrap(ctx, req, modifier, assetId)
			if err != nil {
				return err
			}
		default:
			return errors.New("资产类型错误")
		}
	}
	return nil
}

// 创建报废
func (s *assetScrapService) createScrap(ctx context.Context, req *model.AssetScrapReq, modifier *model.Modifier, ids []uint64) error {
	scrapBluk := make([]*ent.AssetScrapDetailsCreate, 0)
	scrapBatchSn := tools.NewUnique().NewSN28()
	scrapAt := time.Now()
	for _, v := range ids {
		// 是否已经报废
		if b, _ := ent.Database.AssetScrapDetails.Query().Where(assetscrapdetails.AssetID(v)).Exist(ctx); b {
			return fmt.Errorf("资产已经报废")
		}
		// 更新资产状态
		err := ent.Database.Asset.Update().Where(asset.ID(v)).SetStatus(model.AssetStatusScrap.Value()).Exec(ctx)
		if err != nil {
			return err
		}
		// 创建报废记录
		scrapBluk = append(scrapBluk, ent.Database.AssetScrapDetails.Create().SetAssetID(v))
	}
	c, err := ent.Database.AssetScrapDetails.CreateBulk(scrapBluk...).Save(ctx)
	if err != nil {
		return err
	}

	// 创建报废记录
	err = s.orm.Create().
		SetSn(scrapBatchSn).
		SetReasonType(req.ScrapReasonType.Value()).
		SetLastModifier(modifier).
		SetScrapAt(scrapAt).
		SetOperateID(modifier.ID).
		SetOperateRoleType(model.OperatorTypeAssetManager.Value()).
		SetNillableRemark(req.Remark).
		SetCreator(modifier).
		SetLastModifier(modifier).
		SetNum(uint(len(c))).
		AddScrapDetails(c...).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 有编号资产报废
func (s *assetScrapService) scrapAssetWithSN(ctx context.Context, req *model.AssetScrapDetails) ([]uint64, error) {
	ids := make([]uint64, 0)
	if req.Sn == nil {
		return nil, fmt.Errorf("资产ID不能为空")
	}
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
	bat, _ := q.Where(asset.Type(req.AssetType.Value()), asset.Sn(*req.Sn)).First(ctx)
	if bat == nil {
		return nil, fmt.Errorf("资产" + *req.Sn + "不存在或状态不正确")
	}
	ids = append(ids, bat.ID)
	return ids, nil
}

// 无编号资产报废
func (s *assetScrapService) scrapAssetWithoutSN(ctx context.Context, asd *model.AssetScrapDetails, req *model.AssetScrapReq) ([]uint64, error) {
	ids := make([]uint64, 0)

	if req.LocationType != nil && *req.LocationType == model.AssetLocationsTypeRider {
		return nil, errors.New("骑手物资不支持报废")
	}

	if asd.AssetType == model.AssetTypeNonSmartBattery {
		// 非智能电池调拨
		if asd.ModelID == nil || *asd.ModelID == 0 {
			return nil, errors.New(asd.AssetType.String() + "型号ID不能为空")
		}
		item, _ := ent.Database.BatteryModel.Query().Where(batterymodel.ID(*asd.ModelID)).First(ctx)
		if item == nil {
			return nil, errors.New(asd.AssetType.String() + "型号不存在")
		}
	} else {
		if asd.MaterialID == nil || *asd.MaterialID == 0 {
			return nil, errors.New(asd.AssetType.String() + "分类ID不能为空")
		}
		// 判定其它物资类型是否存在
		item, _ := ent.Database.Material.QueryNotDeleted().Where(
			material.ID(*asd.MaterialID),
			material.Type(asd.AssetType.Value()),
		).First(ctx)
		if item == nil {
			return nil, errors.New(asd.AssetType.String() + "分类不存在")
		}
	}
	if asd.Num == nil || *asd.Num == 0 {
		return nil, fmt.Errorf("报废数量不能为空")
	}
	if req.LocationType == nil || req.LocationID == nil {
		return nil, fmt.Errorf("报废位置有误")
	}
	q := ent.Database.Asset.QueryNotDeleted().Where(
		// 资产状态在库存或故障可报废
		asset.StatusIn(
			model.AssetStatusStock.Value(),
			model.AssetStatusFault.Value(),
		),
	)

	// // 资产库存位置在骑手或运维不可报废
	// asset.LocationsTypeNotIn(
	// 	model.AssetLocationsTypeRider.Value(),
	// 	model.AssetLocationsTypeOperation.Value(),
	// ),

	switch *req.LocationType {
	case model.AssetLocationsTypeWarehouse:
		q.Where(asset.HasWarehouseWith(warehouse.ID(*req.LocationID)))
	case model.AssetLocationsTypeStore:
		q.Where(asset.HasStoreWith(store.ID(*req.LocationID)))
	case model.AssetLocationsTypeStation:
		q.Where(asset.HasStationWith(enterprisestation.ID(*req.LocationID)))
	case model.AssetLocationsTypeOperation:
		q.Where(asset.HasOperatorWith(maintainer.ID(*req.LocationID)))
	default:
		return nil, fmt.Errorf("报废位置类型错误")
	}

	if asd.ModelID != nil {
		q.Where(asset.ModelID(*asd.ModelID))
	}
	if asd.MaterialID != nil {
		q.Where(asset.MaterialID(*asd.MaterialID))
	}
	all, _ := q.Where(asset.Type(asd.AssetType.Value())).Limit(int(*asd.Num)).Order(ent.Asc(asset.FieldCreatedAt)).All(ctx)
	if len(all) < int(*asd.Num) {
		return nil, fmt.Errorf("物资数量不足")
	}
	for _, vl := range all {
		ids = append(ids, vl.ID)
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

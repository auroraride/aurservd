// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-18
// Based on aurservd by liasica, magicrolan@qq.com.

package fix

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/auroraride/adapter"
	"github.com/spf13/cobra"

	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assetrole"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/pkg/silk"
)

type CabBatInfo struct {
	modelID uint64
	num     uint
}

type AssetCreateReq struct {
	AssetType     model.AssetType              `json:"assetType" validate:"required"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	SN            *string                      `json:"sn"`                            // 资产编号
	CityID        *uint64                      `json:"cityId"`                        // 城市ID(AssetType为 2:智能电池 需要填写)
	LocationsType model.AssetLocationsType     `json:"locationsType" enums:"1"`       // 资产位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	LocationsID   uint64                       `json:"locationsId"`                   // 资产位置ID
	Attribute     []model.AssetAttributeCreate `json:"attribute"`                     // 属性
	Enable        *bool                        `json:"enable"`                        // 是否启用
	BrandID       *uint64                      `json:"brandId"`                       // 品牌ID(AssetType为 1:电车 需要填写)
	SubscribeID   *uint64                      `json:"subscribeId"`                   // 订阅ID
	Ordinal       *int                         `json:"ordinal"`                       // 序号
}

func Asset() *cobra.Command {

	return &cobra.Command{
		Use:   "asset",
		Short: "修复资产数据",
		Run: func(cmd *cobra.Command, args []string) {
			AssetDo()
		},
	}
}

func AssetDo() {
	ctx := context.Background()
	// 查找超级管理员
	r, _ := ent.Database.AssetRole.Query().Where(assetrole.Super(true)).Only(ctx)
	if r == nil {
		return
	}
	// 给曹博文创建一个账号
	am, _ := ent.Database.AssetManager.Query().Where(assetmanager.PhoneEQ("17719646710")).Only(ctx)
	if am == nil {
		err := biz.NewAssetManager().Create(&definition.AssetManagerCreateReq{
			AssetManagerSigninReq: definition.AssetManagerSigninReq{
				Phone:    "17719646710",
				Password: "123456",
			},
			Name:   "曹博文",
			RoleID: r.ID,
		})
		if err != nil {
			return
		}
		am = ent.Database.AssetManager.Query().Where(assetmanager.PhoneEQ("17719646710")).OnlyX(ctx)
	}

	modifier := &model.Modifier{
		ID:    am.ID,
		Name:  am.Name,
		Phone: am.Phone,
	}
	// 电池
	BatteryDo(ctx, modifier)
	// 车
	// EbikeDo(ctx, modifier)
	// 骑手身上的非智能电池
	// RiderBatteryDo(ctx, modifier)
}

func BatteryDo(ctx context.Context, modifier *model.Modifier) {
	// 插入骑手、电柜身上绑定智能电池的数据
	bat, _ := ent.Database.Battery.QueryNotDeleted().Where(
		battery.Or(
			battery.RiderIDNotNil(),
			battery.CabinetIDNotNil(),
		),
	).All(ctx)
	// 查询 在电柜中非智能电池
	cab, _ := ent.Database.Cabinet.QueryNotDeleted().WithModels().Where(cabinet.Intelligent(false)).All(ctx)
	// 同步电柜信息
	// service.NewCabinet().SyncCabinets(cab)
	// 判定电柜中有多少电池
	// 电柜中电池数量
	var cabBat = make(map[uint64]*CabBatInfo)
	// 查询电池型号
	allModels, _ := ent.Database.BatteryModel.Query().All(ctx)
	allModelsMap := make(map[string]uint64)
	for _, m := range allModels {
		allModelsMap[m.Model] = m.ID
	}

	// 查询非智能电池电池型号
	for _, c := range cab {
		if c.Edges.Models != nil && len(c.Edges.Models) > 0 {
			if modelID, ok := allModelsMap[c.Edges.Models[0].Model]; ok {
				for _, bin := range c.Bin {
					if bin.Battery {
						// 同一个电柜中的电池型号是一样的
						if _, ok := cabBat[c.ID]; !ok {
							cabBat[c.ID] = &CabBatInfo{
								modelID: modelID,
								num:     1,
							}
						} else {
							// 存在
							cabBat[c.ID].num++
						}
					}
				}
			}
		}
	}

	// 创建资产电池资产
	for _, b := range bat {
		var locationsType model.AssetLocationsType
		var locationsID uint64
		if b.RiderID != nil {
			locationsType = model.AssetLocationsTypeRider
			locationsID = *b.RiderID
		}
		if b.CabinetID != nil {
			locationsType = model.AssetLocationsTypeCabinet
			locationsID = *b.CabinetID
		}
		if b.CityID == nil {
			// 如果是骑手绑定电池 根据骑手套餐中的城市ID
			if b.RiderID != nil && b.SubscribeID != nil {
				sub := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(*b.SubscribeID)).OnlyX(ctx)
				if sub != nil {
					b.CityID = &sub.CityID
				}
			}
			// 如果是电柜绑定电池 根据电柜所在城市
			if b.CabinetID != nil {
				ca, _ := ent.Database.Cabinet.Query().Where(cabinet.ID(*b.CabinetID)).Only(ctx)
				if ca != nil {
					b.CityID = ca.CityID
				}
			}
		}

		err := createAsset(ctx, b.ID, &AssetCreateReq{
			AssetType:     model.AssetTypeSmartBattery,
			SN:            &b.Sn,
			CityID:        b.CityID,
			LocationsType: locationsType,
			LocationsID:   locationsID,
			Enable:        silk.Bool(true),
			SubscribeID:   b.SubscribeID,
			Ordinal:       b.Ordinal,
		}, modifier)
		if err != nil {
			continue
		}
	}
	// 创建非智能电池
	for cabID, info := range cabBat {
		// 判定是否重复初始调拨
		b, _ := ent.Database.Asset.QueryNotDeleted().
			Where(
				asset.Type(model.AssetTypeNonSmartBattery.Value()),
				asset.LocationsID(cabID),
				asset.ModelID(info.modelID),
				asset.Status(model.AssetStatusStock.Value()),
			).Count(ctx)
		if uint(b) >= info.num {
			continue
		}

		_, failed, err := service.NewAssetTransfer().Transfer(ctx, &model.AssetTransferCreateReq{
			ToLocationType: model.AssetLocationsTypeCabinet,
			ToLocationID:   cabID,
			Details: []model.AssetTransferCreateDetail{
				{
					AssetType: model.AssetTypeNonSmartBattery,
					Num:       silk.UInt(info.num),
					ModelID:   silk.UInt64(info.modelID),
				},
			},
			Reason:            "系统初始入库",
			AssetTransferType: model.AssetTransferTypeInitial,
			OperatorID:        modifier.ID,
			OperatorType:      model.OperatorTypeAssetManager,
			AutoIn:            true,
		}, modifier)
		if err != nil {
			fmt.Printf("创建非智能电池失败: %v\n", err)
			continue
		}
		if len(failed) > 0 {
			fmt.Printf("创建非智能电池失败: %v\n", failed[0])
			continue
		}
	}

}

// 创建资产
func createAsset(ctx context.Context, id uint64, req *AssetCreateReq, modifier *model.Modifier) error {
	// 去重复
	if req.SN != nil {
		if b, _ := ent.Database.Asset.QueryNotDeleted().Where(asset.Sn(*req.SN), asset.Type(req.AssetType.Value())).Exist(ctx); b {
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

	q := ent.Database.Asset.Create()
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
		name = getAssetName(ctx, req.AssetType, modelInfo.ID)
		q.SetNillableModelID(&modelInfo.ID).
			SetBrandName(ab.Brand.String()).
			SetNillableCityID(req.CityID)
	case model.AssetTypeEbike:
		if req.BrandID == nil {
			return errors.New("品牌不能为空")
		}
		name = getAssetName(ctx, req.AssetType, *req.BrandID)
		q.SetBrandID(*req.BrandID).SetBrandName(name)
	default:
		return errors.New("未知类型")
	}

	q.SetID(id).
		SetType(req.AssetType.Value()).
		SetName(name).
		SetNillableSn(req.SN).
		SetEnable(enable).
		SetStatus(assetStatus).
		SetCreator(modifier).
		SetLastModifier(modifier).
		SetLocationsType((req.LocationsType).Value()).
		SetLocationsID(req.LocationsID).
		SetNillableSubscribeID(req.SubscribeID).
		SetNillableOrdinal(req.Ordinal).
		SetRemark("系统初始入库")

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
	_, failed, err := service.NewAssetTransfer().Transfer(ctx, &model.AssetTransferCreateReq{
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
func getAssetName(ctx context.Context, assetType model.AssetType, materialID uint64) string {
	var name string
	switch assetType {
	case model.AssetTypeSmartBattery, model.AssetTypeNonSmartBattery:
		only, _ := ent.Database.BatteryModel.Query().Where(batterymodel.ID(materialID)).Only(ctx)
		name = only.Model
	case model.AssetTypeEbike:
		only, _ := ent.Database.EbikeBrand.Query().Where(ebikebrand.ID(materialID)).Only(ctx)
		name = only.Name
	case model.AssetTypeCabinetAccessory, model.AssetTypeEbikeAccessory, model.AssetTypeOtherAccessory:
		only, _ := ent.Database.Material.QueryNotDeleted().Where(material.ID(materialID)).Only(ctx)
		if only == nil {
			return "未知"
		}
		name = only.Name
	default:
		name = "未知"
	}
	return name
}

func EbikeDo(ctx context.Context, modifier *model.Modifier) {
	eb, _ := ent.Database.Ebike.Query().Where(ebike.RiderIDNotNil()).All(context.Background())
	// 查询属性
	attributes, _ := ent.Database.AssetAttributes.Query().All(context.Background())
	for _, v := range eb {
		attribute := make([]model.AssetAttributeCreate, 0)
		for _, attr := range attributes {
			var value string
			if attr.Key == "exFactory" {
				value = v.ExFactory
			}
			if attr.Key == "plate" && v.Plate != nil {
				value = *v.Plate
			}
			if attr.Key == "machine" && v.Machine != nil {
				value = *v.Machine
			}
			if attr.Key == "sim" && v.Sim != nil {
				value = *v.Sim
			}
			if attr.Key == "color" {
				value = v.Color
			}
			attribute = append(attribute, model.AssetAttributeCreate{
				AttributeID:    attr.ID,
				AttributeValue: value,
			})
		}

		err := createAsset(ctx, v.ID, &AssetCreateReq{
			AssetType:     model.AssetTypeEbike,
			SN:            silk.String(v.Sn),
			LocationsType: model.AssetLocationsTypeRider,
			LocationsID:   *v.RiderID,
			Attribute:     attribute,
			Enable:        silk.Bool(true),
			BrandID:       silk.UInt64(v.BrandID),
		}, modifier)
		if err != nil {
			fmt.Printf("创建电车失败: %v\n", err)
			continue
		}
	}
}

// RiderBatteryDo 骑手身上的非智能电池创建资产
func RiderBatteryDo(ctx context.Context, modifier *model.Modifier) {
	r, _ := ent.Database.Rider.QueryNotDeleted().All(ctx)
	for _, v := range r {

		subd, _ := service.NewSubscribe().RecentDetail(v.ID)
		if subd == nil || subd.Intelligent {
			continue
		}
		// 计费和逾期的骑手 才创建资产
		if subd.Status == model.SubscribeStatusUsing || subd.Status == model.SubscribeStatusOverdue {
			// 查询非智能订阅电池型号
			m, _ := ent.Database.BatteryModel.Query().Where(batterymodel.Model(subd.Model)).Only(ctx)
			if m == nil {
				parts := strings.Split(subd.Model, "V")

				// 提取第一个数字部分
				vPart := parts[0]

				// 提取第二个数字部分
				ahParts := strings.Split(parts[1], "AH")
				ahPart := ahParts[0]

				// 将字符串转换为uint
				vNum, _ := strconv.ParseUint(vPart, 10, 64)
				ahNum, _ := strconv.ParseUint(ahPart, 10, 64)
				err := service.NewBatteryModel().Create(&model.BatteryModelCreateReq{
					Voltage:  uint(vNum),
					Capacity: uint(ahNum),
				})
				if err != nil {
					fmt.Printf("创建电池型号失败: %v\n", err)
					continue
				}
				m = ent.Database.BatteryModel.Query().Where(batterymodel.Model(subd.Model)).OnlyX(ctx)
			}

			// 判定是否重复初始调拨
			b, _ := ent.Database.Asset.QueryNotDeleted().
				Where(
					asset.Type(model.AssetTypeNonSmartBattery.Value()),
					asset.LocationsID(v.ID),
					asset.ModelID(m.ID),
					asset.Status(model.AssetStatusStock.Value()),
				).Count(ctx)
			if b > 0 {
				continue
			}

			sn, failed, err := service.NewAssetTransfer().Transfer(ctx, &model.AssetTransferCreateReq{
				ToLocationType: model.AssetLocationsTypeRider,
				ToLocationID:   v.ID,
				Details: []model.AssetTransferCreateDetail{
					{
						AssetType: model.AssetTypeNonSmartBattery,
						Num:       silk.UInt(1),
						ModelID:   silk.UInt64(m.ID),
					},
				},
				Reason:            "骑手绑定非智能系统初始入库",
				AssetTransferType: model.AssetTransferTypeInitial,
				OperatorID:        modifier.ID,
				OperatorType:      model.OperatorTypeAssetManager,
				AutoIn:            true,
			}, modifier)
			if err != nil {
				fmt.Printf("创建非智能电池失败: %v\n", err)
				continue
			}
			if len(failed) > 0 {
				fmt.Printf("创建非智能电池失败: %v\n", failed[0])
				continue
			}

			// 更新资产订阅
			d, _ := ent.Database.AssetTransferDetails.QueryNotDeleted().Where(assettransferdetails.HasTransferWith(assettransfer.Sn(sn))).First(ctx)
			if d == nil {
				continue
			}
			err = ent.Database.Asset.Update().Where(asset.ID(d.AssetID)).SetSubscribeID(subd.ID).Exec(ctx)
			if err != nil {
				continue
			}
		}
	}
}

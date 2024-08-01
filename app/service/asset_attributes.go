package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assetattributes"
)

type assetAttributesService struct {
	orm *ent.AssetAttributesClient
}

func NewAssetAttributes() *assetAttributesService {
	return &assetAttributesService{
		orm: ent.Database.AssetAttributes,
	}
}

// Initialize 初始化资产属性
func (s *assetAttributesService) Initialize() {
	for _, v := range model.InitAssetAttributes {
		// 如果存在则不创建
		if b, err := s.orm.Query().Where(assetattributes.AssetType(v.AssetType.Value())).Exist(context.Background()); err != nil || b {
			continue
		}
		for _, vl := range v.Attribute {
			if err := s.orm.Create().
				SetName(vl.AttributeName).
				SetAssetType(v.AssetType.Value()).
				SetKey(vl.AttributeKey).
				Exec(context.Background()); err != nil {
				zap.L().Error("初始化资产属性失败", zap.Error(err))
				continue
			}
		}
	}
}

// Create 创建资产属性
func (s *assetAttributesService) Create(ctx context.Context, req *model.AssetAttributesCreateReq) (err error) {
	for _, v := range req.Attribute {
		attr, _ := s.orm.Create().SetName(v.AttributeName).SetAssetType(req.AssetType.Value()).Save(ctx)
		if attr != nil {
			err = ent.Database.AssetAttributeValues.Create().SetValue(v.AttributeValue).SetAttributeID(attr.ID).Exec(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// List 资产属性列表
func (s *assetAttributesService) List(ctx context.Context, req *model.AssetAttributesListReq) (res []*model.AssetAttributesListRes) {
	attrs, _ := s.orm.Query().Where(assetattributes.AssetType(req.AssetType.Value())).All(ctx)
	for _, v := range attrs {
		res = append(res, &model.AssetAttributesListRes{
			AttributeID:   v.ID,
			AttributeName: v.Name,
			AttributeKey:  v.Key,
		})
	}
	return res
}

package service

import (
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionachievement"
	"github.com/auroraride/aurservd/pkg/snag"
)

type promotionAchievementService struct {
	ctx context.Context
	*BaseService
}

func NewPromotionAchievementService(params ...any) *promotionAchievementService {
	return &promotionAchievementService{
		ctx:         context.Background(),
		BaseService: newService(params...),
	}
}

// List 获取成就列表
func (s *promotionAchievementService) List() []promotion.Achievement {
	item, _ := ent.Database.PromotionAchievement.QueryNotDeleted().Order(ent.Desc(promotionachievement.FieldCreatedAt)).All(s.ctx)
	res := make([]promotion.Achievement, 0, len(item))
	for _, v := range item {
		res = append(res, promotion.Achievement{
			ID:        v.ID,
			Name:      v.Name,
			Icon:      v.Icon,
			Type:      promotion.AchievementType(v.Type),
			Condition: v.Condition,
		})
	}
	return res
}

// Create 新增成就
func (s *promotionAchievementService) Create(req *promotion.Achievement) {
	_, err := ent.Database.PromotionAchievement.Create().
		SetName(req.Name).
		SetIcon(req.Icon).
		SetType(req.Type.Value()).
		SetCondition(req.Condition).
		Save(s.ctx)
	if err != nil {
		snag.Panic("新增成就失败")
	}
}

// Update 修改成就
func (s *promotionAchievementService) Update(req *promotion.Achievement) {
	_, err := ent.Database.PromotionAchievement.
		UpdateOneID(req.ID).
		SetName(req.Name).
		SetIcon(req.Icon).
		SetType(req.Type.Value()).
		SetCondition(req.Condition).
		Save(s.ctx)
	if err != nil {
		snag.Panic("修改任务失败")
	}
}

// Delete 删除成就
func (s *promotionAchievementService) Delete(id uint64) {
	// 软删除
	_, err := ent.Database.PromotionAchievement.Update().
		Where(promotionachievement.ID(id)).
		SetDeletedAt(time.Now()).
		Save(s.ctx)
	if err != nil {
		snag.Panic("删除成就失败")
	}
}

// UploadIcon 上传icon
func (s *promotionAchievementService) UploadIcon(ctx echo.Context) promotion.UploadIcon {
	// 获取单张图片
	f, err := ctx.FormFile("icon")
	if err != nil {
		snag.Panic("上传失败")
	}
	src, err := f.Open()
	if err != nil {
		snag.Panic("上传图片失败")
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	// 确保只接受指定的图片格式
	ext := filepath.Ext(f.Filename)
	if !NewFeedback().IsValidImageExtension(ext) {
		snag.Panic("只支持jpg、jpeg、png格式的图片")
	}
	// 生成相对路径
	randomNum := rand.Intn(1000) // 生成一个随机数，用于防止同一秒钟上传多个文件时的冲突
	r := filepath.Join("promotion", "achievement", fmt.Sprintf("%s%d%s", time.Now().
		Format(carbon.ShortDateTimeLayout), randomNum, ext))

	err = ali.NewOss().Bucket.PutObject(r, src)
	if err != nil {
		zap.L().Error("上传图片失败", zap.Error(err))
		snag.Panic("上传图片失败")
	}
	return promotion.UploadIcon{Icon: r}
}

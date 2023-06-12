package service

import (
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/feedback"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type feedbackService struct {
	*BaseService
	orm *ent.FeedbackClient
}

func NewFeedback(params ...any) *feedbackService {
	return &feedbackService{
		BaseService: newService(params...),
		orm:         ent.Database.Feedback,
	}
}

func (s *feedbackService) Query(id uint64) (*ent.Feedback, error) {
	return s.orm.Query().Where(feedback.ID(id)).WithEnterprise().First(s.ctx)
}

func (s *feedbackService) QueryX(id uint64) *ent.Feedback {
	ag, _ := s.Query(id)
	if ag == nil {
		snag.Panic("未找到反馈信息")
	}
	return ag
}

// Create 创建反馈
func (s *feedbackService) Create(req *model.FeedbackReq, ag *ent.Agent) bool {
	_, err := s.orm.Create().SetEnterpriseID(ag.EnterpriseID).
		SetContent(req.Content).
		SetType(req.Type).
		SetURL(req.Url).
		SetName(ag.Name).
		SetPhone(ag.Phone).
		Save(s.ctx)
	if err != nil {
		snag.Panic("添加失败")
	}
	return true
}

// FeedbackList List 反馈列表
func (s *feedbackService) FeedbackList(req *model.FeedbackListReq) *model.PaginationRes {
	q := s.orm.Query().WithEnterprise().Order(ent.Desc(feedback.FieldCreatedAt))
	// 筛选条件
	if req.Keyword != "" {
		q.Where(feedback.ContentContains(req.Keyword))
	}
	if req.Type != nil {
		q.Where(feedback.TypeEQ(*req.Type))
	}
	if req.Start != nil {
		q.Where(feedback.CreatedAtGTE(tools.NewTime().ParseDateStringX(*req.Start)))
	}
	if req.End != nil {
		q.Where(feedback.CreatedAtLT(tools.NewTime().ParseNextDateStringX(*req.End)))
	}

	if req.EnterpriseID != nil {
		q.Where(feedback.EnterpriseID(*req.EnterpriseID))
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Feedback) model.FeedbackDetail {
		rsp := model.FeedbackDetail{
			ID:                     item.ID,
			Content:                item.Content,
			Url:                    item.URL,
			Type:                   item.Type,
			EnterpriseName:         item.Edges.Enterprise.Name,
			EnterpriseContactName:  item.Name,
			EnterpriseContactPhone: item.Phone,
			CreatedAt:              item.CreatedAt.Format(carbon.DateTimeLayout),
		}
		return rsp
	})

}

// UploadImage 上传照片
func (s *feedbackService) UploadImage(c echo.Context) []string {
	const maxUploadSize = 50 * 1024 * 1024
	if c.Request().ParseMultipartForm(maxUploadSize) != nil {
		snag.Panic("上传图片失败")
	}

	files := c.Request().MultipartForm.File["images"]
	if len(files) > 5 {
		snag.Panic("最多上传5张图片")
	}
	var paths []string
	for _, f := range files {
		r, err := s.doFile(f)
		if err != nil {
			snag.Panic("上传图片失败")
		}
		paths = append(paths, r)
	}
	return paths
}

func (s *feedbackService) doFile(f *multipart.FileHeader) (string, error) {
	// 限制单张图片大小为10M以下
	if f.Size > 10<<20 {
		snag.Panic("单张图片大小不能超过10M")
	}

	src, err := f.Open()
	if err != nil {
		log.Println(err)
		snag.Panic("上传图片失败")
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	// 确保只接受指定的图片格式
	ext := filepath.Ext(f.Filename)
	if !s.isValidImageExtension(ext) {
		snag.Panic("只接受jpg、jpeg、png、svg格式的图片")
	}

	// 生成相对路径
	randomNum := rand.Intn(1000) // 生成一个随机数，用于防止同一秒钟上传多个文件时的冲突
	r := filepath.Join("agent", "feedback", fmt.Sprintf("%s%d%s", time.Now().
		Format(carbon.ShortDateTimeLayout), randomNum, ext))

	return r, ali.NewOss().Bucket.PutObject(r, src)
}

func (*feedbackService) isValidImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".svg"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

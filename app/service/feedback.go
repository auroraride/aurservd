package service

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/feedback"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
	"github.com/auroraride/aurservd/pkg/utils"
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
	return s.orm.QueryNotDeleted().Where(feedback.ID(id)).WithEnterprise().First(s.ctx)
}

func (s *feedbackService) QueryX(id uint64) *ent.Feedback {
	ag, _ := s.Query(id)
	if ag == nil {
		snag.Panic("未找到反馈信息")
	}
	return ag
}

// Create 创建反馈
func (s *feedbackService) Create(req *model.FeedbackReq, enterprise *ent.Enterprise) bool {
	_, err := s.orm.Create().SetEnterpriseID(enterprise.ID).
		SetContent(req.Content).
		SetType(req.Type).
		SetURL(req.Url).
		SetName(enterprise.Name).
		SetPhone(enterprise.ContactPhone).
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
	if req.StartTime != nil && req.EndTime != nil {
		q.Where(feedback.CreatedAtGTE(tools.NewTime().ParseDateStringX(*req.StartTime)), feedback.CreatedAtLT(tools.NewTime().ParseNextDateStringX(*req.EndTime)))
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
			EnterpriseName:         item.Name,
			EnterpriseContactName:  item.Edges.Enterprise.Name,
			EnterpriseContactPhone: item.Edges.Enterprise.ContactPhone,
			CreatedAt:              item.CreatedAt.Format(carbon.DateTimeLayout),
		}
		return rsp
	})

}

// UploadImage 上传照片本地文件夹
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
	for _, file := range files {
		// 限制单张图片大小为10M以下
		if file.Size > 10<<20 {
			snag.Panic("单张图片大小不能超过10M")
		}

		src, err := file.Open()
		if err != nil {
			log.Println(err)
			snag.Panic("上传图片失败")
		}
		defer src.Close()

		// 确保只接受指定的图片格式
		ext := filepath.Ext(file.Filename)
		if !IsValidImageExtension(ext) {
			snag.Panic("只接受jpg、jpeg、png、svg格式的图片")
		}
		// 生成相对路径
		randomNum := rand.Intn(1000) // 生成一个随机数，用于防止同一秒钟上传多个文件时的冲突
		r := filepath.Join("runtime", "uploads", fmt.Sprintf("%s%d%s", time.Now().
			Format(carbon.ShortDateTimeLayout), randomNum, ext))
		if utils.NewFile(r).CreateDirectoryIfNotExist() != nil {
			snag.Panic("上传图片失败")
		}
		// 创建目标文件
		dst, err := os.Create(r)
		if err != nil {
			snag.Panic("上传图片失败")
		}
		defer dst.Close()

		// 将源文件内容复制到目标文件
		if _, err = io.Copy(dst, src); err != nil {
			snag.Panic("上传图片失败")
		}
		paths = append(paths, r)
	}
	return paths
}

func IsValidImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".svg"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
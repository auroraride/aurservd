package biz

import (
	"context"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/feedback"
)

type feedbackBiz struct {
	orm *ent.FeedbackClient
	ctx context.Context
}

func NewFeedback() *feedbackBiz {
	return &feedbackBiz{
		orm: ent.Database.Feedback,
		ctx: context.Background(),
	}
}

// RiderCreate 骑手创建反馈
func (s *feedbackBiz) RiderCreate(r *ent.Rider, req *model.FeedbackReq) error {
	_, err := s.orm.Create().
		SetContent(req.Content).
		SetSource(model.SourceRider).
		SetType(req.Type).
		SetURL(req.Url).
		SetName(r.Name).
		SetPhone(r.Phone).
		SetRider(r).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// List 反馈列表
func (s *feedbackBiz) List(r *ent.Rider, req *model.FeedbackListReq) (res *model.PaginationRes) {
	q := s.orm.Query().Where(feedback.RiderID(r.ID)).Order(ent.Desc(feedback.FieldID))
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Feedback) (res *model.FeedbackDetail) {
		res = &model.FeedbackDetail{
			ID:        item.ID,
			Content:   item.Content,
			Url:       item.URL,
			Type:      item.Type,
			CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
			Name:      item.Name,
			Phone:     item.Phone,
		}
		return
	})
}

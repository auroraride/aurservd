package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

type feedbackBiz struct {
	orm *ent.FeedbackClient
}

func NewFeedback() *feedbackBiz {
	return &feedbackBiz{
		orm: ent.Database.Feedback,
	}
}

// RiderCreate 骑手创建反馈
func (s *feedbackBiz) RiderCreate(req *model.FeedbackReq, ri *ent.Rider) bool {
	ctx := context.Background()
	// 保存反馈信息
	query := s.orm.Create().SetEnterpriseID(*ri.EnterpriseID).
		SetContent(req.Content).
		SetSource(model.SourceRider). // 反馈来源
		SetType(req.Type).
		SetURL(req.Url).
		SetName(ri.Name).
		SetPhone(ri.Phone)
	if query.Exec(ctx) != nil {
		snag.Panic("添加失败")
	}
	return false
}

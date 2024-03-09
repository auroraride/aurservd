package biz

import "github.com/auroraride/aurservd/internal/ent"

type plan struct {
	orm *ent.PlanClient
}

func NewPlan() *plan {
	return &plan{
		orm: ent.Database.Plan,
	}
}

// 查询骑士卡详情

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/contracttemplate"
)

type ContractTemplate struct {
	orm *ent.ContractTemplateClient
	ctx context.Context
}

func NewContractTemplate() *ContractTemplate {
	return &ContractTemplate{
		orm: ent.Database.ContractTemplate,
		ctx: context.Background(),
	}
}

// Create 新增合同模板
func (s *ContractTemplate) Create(req *definition.ContractTemplateCreateReq) error {
	// 判断当前传入类型是否存在
	temp, err := s.orm.QueryNotDeleted().Where(
		contracttemplate.Aimed(req.Aimed.Value()),
		contracttemplate.PlanType(req.PlanType.Value()),
		contracttemplate.Enable(true),
	).First(s.ctx)
	// 如果存在并且需要开启
	if temp != nil && req.Enable != nil && *req.Enable {
		_, err = temp.Update().SetEnable(false).Save(s.ctx)
		if err != nil {
			return err
		}
	}

	_, err = s.orm.Create().
		SetName(req.Name).
		SetAimed(req.Aimed.Value()).
		SetPlanType(req.PlanType.Value()).
		SetURL(req.Url).
		SetHash(req.Hash).
		SetNillableRemark(req.Remark).
		SetNillableEnable(req.Enable).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// Modify 编辑合同模板
func (s *ContractTemplate) Modify(req *definition.ContractTemplateModifyReq) error {
	// 判断当前传入类型是否存在
	temp, _ := s.orm.QueryNotDeleted().Where(
		contracttemplate.PlanType(req.PlanType.Value()),
		contracttemplate.Aimed(req.Aimed.Value()),
		contracttemplate.Enable(true),
	).First(s.ctx)

	// 如果存在并且需要开启
	if temp != nil && *req.Enable {
		_, err := temp.Update().SetEnable(false).Save(s.ctx)
		if err != nil {
			return err
		}
	}

	_, err := s.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetAimed(req.Aimed.Value()).
		SetPlanType(req.PlanType.Value()).
		SetNillableEnable(req.Enable).
		SetNillableRemark(req.Remark).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// List 列表
func (s *ContractTemplate) List() (res []*definition.ContractTemplateListRes) {
	res = make([]*definition.ContractTemplateListRes, 0)
	all := s.orm.Query().Order(ent.Desc(contracttemplate.FieldID)).AllX(s.ctx)
	for _, v := range all {
		res = append(res, &definition.ContractTemplateListRes{
			ContractTemplate: definition.ContractTemplate{
				ID:        v.ID,
				Name:      v.Name,
				Aimed:     definition.ContractTemplateAimed(v.Aimed),
				PlanType:  model.PlanType(v.PlanType),
				Url:       v.URL,
				Hash:      v.Hash,
				Enable:    v.Enable,
				CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
				Remark:    v.Remark,
			},
		})
	}
	return res
}

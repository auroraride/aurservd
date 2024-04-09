package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
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
		contracttemplate.SubType(req.SubType),
		contracttemplate.UserType(req.UserType),
		contracttemplate.Enable(true),
	).First(s.ctx)
	// 如果存在并且需要开启
	if temp != nil && req.Enable != nil && *req.Enable {
		_, err = temp.Update().SetEnable(false).Save(s.ctx)
		if err != nil {
			return err
		}
	}

	// var fields []string

	// fields = append(fields, definition.ContractTemplateFields...)
	// if f, ok := definition.FieldsUserMap[req.UserType]; ok {
	// 	fields = append(fields, f...)
	// }
	// if f, ok := definition.FieldsSubMap[req.SubType]; ok {
	// 	fields = append(fields, f...)
	// }

	// sn, err := rpc.AddContractTemplate(s.ctx, req.Url, fields)
	// if err != nil {
	// 	return err
	// }
	// if sn == "" {
	// 	return errors.New("请求失败")
	// }

	sn := "123456"

	_, err = s.orm.Create().
		SetName(req.Name).
		SetUserType(req.UserType).
		SetSubType(req.SubType).
		SetURL(req.Url).
		SetSn(sn).
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
		contracttemplate.SubType(req.SubType),
		contracttemplate.UserType(req.UserType),
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
		SetUserType(req.UserType).
		SetSubType(req.SubType).
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
				UserType:  v.UserType,
				SubType:   v.SubType,
				Url:       v.URL,
				Sn:        v.Sn,
				Enable:    v.Enable,
				CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
				Remark:    v.Remark,
			},
		})
	}
	return res
}

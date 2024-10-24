package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/instructions"
)

type instructionsBiz struct {
	orm      *ent.InstructionsClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewInstructions() *instructionsBiz {
	return &instructionsBiz{
		orm: ent.Database.Instructions,
		ctx: context.Background(),
	}
}

func NewInstructionsWithModifierBiz(m *model.Modifier) *instructionsBiz {
	s := NewInstructions()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// InitInstructions 初始化
func (s *instructionsBiz) InitInstructions() {
	for k, v := range definition.InstructionsColumns {
		if !s.orm.QueryNotDeleted().Where(instructions.Key(k)).ExistX(s.ctx) {
			var content interface{} = ""
			s.orm.Create().
				SetContent(&content).
				SetTitle(v).
				SetKey(k).
				SaveX(s.ctx)
		}
	}
}

func (s *instructionsBiz) Modify(req *definition.InstructionsCreateReq) error {
	_, err := s.orm.Update().
		Where(instructions.KeyEQ(req.Key)).
		SetContent(&req.Content).
		SetTitle(req.Title).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// Detail 详情
func (s *instructionsBiz) Detail(key string) (res *definition.InstructionsRes, err error) {
	item, err := s.orm.Query().Where(instructions.KeyEQ(key)).First(s.ctx)
	if err != nil {
		return nil, err
	}
	res = &definition.InstructionsRes{
		Content: item.Content,
		Title:   item.Title,
		Key:     item.Key,
	}
	return res, nil
}

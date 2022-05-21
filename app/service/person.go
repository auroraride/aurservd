// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/snag"
)

type personService struct {
    ctx context.Context
    orm *ent.PersonClient
}

func NewPerson() *personService {
    return &personService{
        ctx: context.Background(),
        orm: ar.Ent.Person,
    }
}

func (s *personService) Query(id uint64) *ent.Person {
    item, err := ar.Ent.Rider.QueryNotDeleted().WithPerson().Where(rider.ID(id)).Only(s.ctx)
    if err != nil || item == nil || item.Edges.Person == nil {
        snag.Panic("未找到骑手实名信息")
    }
    return item.Edges.Person
}

// Block 封禁或取消封禁
func (s *personService) Block(m *model.Modifier, req *model.PersonBlockReq) {
    item := s.Query(req.ID)
    if req.Block == item.Banned {
        snag.Panic("骑手已是封禁状态")
    }
    _, err := s.orm.UpdateOne(item).SetBanned(req.Block).SetLastModifier(m).Save(s.ctx)
    if err != nil {
        snag.Panic(err)
    }
    nb := "未封禁"
    bd := "已封禁"
    ol := logging.CreateOperateLog().SetRef(item).SetModifier(m)
    if req.Block {
        // 封禁
        ol.SetOperate(logging.OperatePersonBan).SetDiff(nb, bd)
    } else {
        ol.SetOperate(logging.OperatePersonUnBan).SetDiff(bd, nb)
    }
    ol.PutOperateLog()
}

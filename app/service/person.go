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
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/snag"
)

type personService struct {
    ctx      context.Context
    orm      *ent.PersonClient
    modifier *model.Modifier
}

func NewPerson() *personService {
    return &personService{
        ctx: context.Background(),
        orm: ar.Ent.Person,
    }
}

func NewPersonWithModifier(m *model.Modifier) *personService {
    s := NewPerson()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *personService) Query(id uint64) (*ent.Person, *ent.Rider) {
    item, err := ar.Ent.Rider.QueryNotDeleted().WithPerson().Where(rider.ID(id)).Only(s.ctx)
    if err != nil || item == nil || item.Edges.Person == nil {
        snag.Panic("未找到骑手实名信息")
    }
    return item.Edges.Person, item
}

// Ban 封禁或取消封禁
func (s *personService) Ban(req *model.PersonBanReq) {
    item, r := s.Query(req.ID)
    if req.Ban == item.Banned {
        snag.Panic("骑手已是封禁状态")
    }
    _, err := s.orm.UpdateOne(item).SetBanned(req.Ban).Save(s.ctx)
    if err != nil {
        snag.Panic(err)
    }
    nb := "未封禁"
    bd := "已封禁"
    ol := logging.NewOperateLog().SetRef(r).SetModifier(s.modifier)
    if req.Ban {
        // 封禁
        ol.SetOperate(model.OperatePersonBan).SetDiff(nb, bd)
    } else {
        ol.SetOperate(model.OperatePersonUnBan).SetDiff(bd, nb)
    }
    ol.Send()
}

// GetNormalAuthedPerson 获取正常骑手已实名认证的信息
func (s *personService) GetNormalAuthedPerson(u *ent.Rider) *ent.Person {
    if u.Blocked {
        snag.Panic("你已被封禁")
    }

    p := u.Edges.Person
    if p == nil && u.PersonID != nil {
        p, _ = ar.Ent.Person.QueryNotDeleted().Where(person.ID(*u.PersonID)).First(s.ctx)
    }

    if p == nil {
        snag.Panic("未找到实名认证信息")
    }

    if model.PersonAuthStatus(p.Status).RequireAuth() {
        snag.Panic("未实名认证")
    }

    if p.Banned {
        snag.Panic("你已被拉黑")
    }

    return p
}

// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/auroraride/aurservd/pkg/snag"
)

type contractService struct {
    esign *esign.Esign
}

func New() *contractService {
    return &contractService{
        esign: esign.New(),
    }
}

func (s *contractService) Sign(u *ent.Rider) {
    var (
        orm       = ar.Ent
        person    = u.Edges.Person
        accountId = person.EsignAccountID
    )

    // 创建 / 获取 签约个人账号
    if accountId == nil {
        accountId = s.esign.CreatePersonAccount(esign.CreatePersonAccountReq{
            ThirdPartyUserId: person.IDCardNumber,
            Name:             person.Name,
            IdType:           "CRED_PSN_CH_IDCARD",
            IdNumber:         person.IDCardNumber,
            Mobile:           u.Phone,
        })
        if accountId == nil {
            snag.Panic("签署账号生成失败")
        }
        // 保存个人账号
        err := orm.Person.UpdateOneID(person.ID).SetNillableEsignAccountID(accountId).Exec(context.Background())
        if err != nil {
            snag.Panic(err)
        }
    }
}

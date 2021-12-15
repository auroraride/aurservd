// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/auroraride/aurservd/pkg/snag"
)

type contractService struct {
    esign *esign.Esign
}

func NewContract() *contractService {
    return &contractService{
        esign: esign.New(),
    }
}

// Sign 签署合同
// TODO 保存合同
func (s *contractService) Sign(u *ent.Rider) string {
    var (
        cfg        = s.esign.Config
        orm        = ar.Ent
        person     = u.Edges.Person
        accountId  = u.EsignAccountID
        isGroup    = u.GroupID != nil
        templateId = cfg.Person.TemplateId
        req        = esign.CreateFlowReq{
            Scene: cfg.Person.Scene,
        }
    )

    if isGroup {
        templateId = cfg.Group.TemplateId
    }

    // 创建 / 获取 签约个人账号
    if accountId == "" {
        accountId = s.esign.CreatePersonAccount(esign.CreatePersonAccountReq{
            ThirdPartyUserId: u.Phone,
            Name:             person.Name,
            IdType:           "CRED_PSN_CH_IDCARD",
            IdNumber:         person.IDCardNumber,
            Mobile:           u.Phone,
        })
        if accountId == "" {
            snag.Panic("签署账号生成失败")
        }
        // 保存个人账号
        err := orm.Rider.UpdateOneID(u.ID).SetEsignAccountID(accountId).Exec(context.Background())
        if err != nil {
            snag.Panic(err)
        }
    }

    // 设置个人账户ID
    req.PersonAccountId = accountId

    // 获取模板控件
    tmpl := s.esign.DocTemplate(templateId)
    m := make(ar.Map)
    for _, com := range tmpl.StructComponents {
        switch com.Key {
        case cfg.EntSignKey:
            req.EntSignBean = esign.PosBean{
                PosPage: fmt.Sprintf("%v", com.Context.Pos.Page),
                PosX:    com.Context.Pos.X,
                PosY:    com.Context.Pos.Y,
            }
        case cfg.PsnSignKey:
            req.PsnSignBean = esign.PosBean{
                PosPage: fmt.Sprintf("%v", com.Context.Pos.Page),
                PosX:    com.Context.Pos.X,
                PosY:    com.Context.Pos.Y,
            }
        default:
            // todo 动态配置真实值
            m[com.Key] = com.Key
        }
    }
    // 填充内容生成PDF
    pdf := s.esign.CreateByTemplate(esign.CreateByTemplateReq{
        Name:             req.Scene + ".pdf", // todo 文件名
        SimpleFormFields: m,
        TemplateId:       templateId,
    })
    req.FileId = pdf.FileId

    // 发起签署，获取flowId
    flowId := s.esign.CreateFlowOneStep(req)

    // 获取签署链接
    return s.esign.ExecuteUrl(flowId, accountId)
}

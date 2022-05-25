// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/contract"
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
)

const (
    snKey      = "sn"      // 合同编号
    entSignKey = "entSign" // 企业签章控件名称
    psnSignKey = "psnSign" // 个人签章控件名称
)

type contractService struct {
    esign *esign.Esign
    ctx   context.Context
    orm   *ent.ContractClient
}

func NewContract() *contractService {
    return &contractService{
        esign: esign.New(),
        orm:   ar.Ent.Contract,
        ctx:   context.Background(),
    }
}

// Effective 查询骑手是否存在生效中的合同
// 当用户退租之后触发合同失效, 需要重新签订 // TODO 需要实现逻辑
func (s *contractService) Effective(u *ent.Rider) bool {
    return s.orm.QueryNotDeleted().Where(
        contract.RiderID(u.ID),
        contract.Status(model.ContractStatusSuccess.Raw()),
        contract.Effective(true),
    ).ExistX(s.ctx)
}

// generateSn 生成合同编号
func (s *contractService) generateSn() string {
    return fmt.Sprintf("%s%06d", carbon.Now().ToShortDateTimeString(), utils.RandomIntMaxMin(1000, 999999))
}

// Sign 签署合同
func (s *contractService) Sign(u *ent.Rider, params *model.ContractSignReq) model.ContractSignRes {
    if s.Effective(u) {
        return model.ContractSignRes{Effective: true}
    }
    var (
        sn           = s.generateSn()
        cfg          = s.esign.Config
        orm          = ar.Ent
        person       = u.Edges.Person
        accountId    = u.EsignAccountID
        isEnterprise = u.EnterpriseID != nil
        templateId   = cfg.Person.TemplateId
        req          = esign.CreateFlowReq{
            Scene: cfg.Person.Scene,
        }
    )

    if isEnterprise {
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

    // 设置合同编号
    s.esign.SetSn(sn)

    // 设置个人账户ID
    req.PersonAccountId = accountId

    // 获取模板控件
    tmpl := s.esign.DocTemplate(templateId)
    m := make(ar.Map)
    for _, com := range tmpl.StructComponents {
        switch com.Key {
        case snKey:
            m[snKey] = sn
        case entSignKey:
            req.EntSignBean = esign.PosBean{
                PosPage: fmt.Sprintf("%v", com.Context.Pos.Page),
                PosX:    com.Context.Pos.X,
                PosY:    com.Context.Pos.Y,
            }
        case psnSignKey:
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
        Name:             fmt.Sprintf("%s-%s.pdf", req.Scene, sn), // todo 文件名
        SimpleFormFields: m,
        TemplateId:       templateId,
    })
    req.FileId = pdf.FileId

    // 发起签署，获取flowId
    flowId := s.esign.CreateFlowOneStep(req)

    // 获取签署链接
    link := s.esign.ExecuteUrl(flowId, accountId)

    // 存储合同信息
    err := ar.Ent.Contract.Create().
        SetFlowID(flowId).
        SetRiderID(u.ID).
        SetStatus(model.ContractStatusPending.Raw()).
        SetSn(sn).
        Exec(context.Background())
    if err != nil {
        snag.Panic(err)
    }
    return model.ContractSignRes{
        Url: link,
        Sn:  sn,
    }
}

// Result 合同签署结果
func (s *contractService) Result(u *ent.Rider, sn string) model.StatusResponse {
    orm := ar.Ent.Contract
    // 查询合同是否存在
    c, err := orm.QueryNotDeleted().
        Where(contract.Sn(sn), contract.RiderID(u.ID)).
        Only(context.Background())
    if err != nil || c == nil {
        snag.Panic("合同查询失败")
    }
    success := s.esign.Result(c.FlowID)
    update := orm.UpdateOneID(c.ID)
    if success {
        // 获取合同并上传到阿里云
        p := u.Edges.Person
        update.SetStatus(model.ContractStatusSuccess.Raw()).
            SetFiles(s.esign.DownloadDocument(fmt.Sprintf("%s-%s/contracts/", p.Name, p.IDCardNumber), c.FlowID))
    }
    err = update.Exec(context.Background())
    if err != nil {
        snag.Panic(err)
    }
    return model.StatusResponse{Status: success}
}

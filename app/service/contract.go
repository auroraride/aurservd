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
    "github.com/auroraride/aurservd/internal/ent/enterprise"
    "github.com/auroraride/aurservd/internal/ent/enterprisestation"
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "math"
    "strconv"
    "strings"
    "time"
)

const (
    snKey      = "sn"         // 合同编号
    aurSeal    = "aurSeal"    // 平台签章控件名称
    riderSeal  = "riderSeal"  // 骑手签章控件名称
    pagingSeal = "pagingSeal" // 骑缝章
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
// 当用户退租之后触发合同失效, 需要重新签订
func (s *contractService) Effective(u *ent.Rider) bool {
    if u.Contractual {
        return true
    }
    exists, _ := s.orm.QueryNotDeleted().Where(
        contract.RiderID(u.ID),
        contract.Status(model.ContractStatusSuccess.Raw()),
        contract.Effective(true),
    ).Exist(s.ctx)
    return exists
}

// planData 个签合同数据
func (s *contractService) planData(planID uint64, m ar.Map) {
    p := NewPlan().QueryEffectiveWithID(planID)
    month := math.Round(float64(p.Days) / 31.0)
    price := p.Price
    if month > 1 {
        price = price / month
    }
    m["month"] = fmt.Sprintf("%.0f", month)
    m["payMonth"] = fmt.Sprintf("%.0f", month)
    m["price"] = fmt.Sprintf("%.2f", price)
    m["amount"] = fmt.Sprintf("%.2f", p.Price)
}

func (s *contractService) enterpriseData(u *ent.Rider, m ar.Map) {
    e := u.Edges.Enterprise
    if e == nil {
        e, _ = ar.Ent.Enterprise.QueryNotDeleted().Where(enterprise.ID(*u.EnterpriseID)).First(s.ctx)
    }
    if e == nil {
        snag.Panic("骑手企业查找失败")
    }

    sta := u.Edges.Station
    if sta == nil {
        sta, _ = ar.Ent.EnterpriseStation.QueryNotDeleted().Where(enterprisestation.ID(*u.StationID)).First(s.ctx)
    }
    if sta == nil {
        snag.Panic("骑手站点查找失败")
    }

    // entName entPhone station payerEnt
    m["entName"] = e.Name
    m["entPhone"] = e.ContactPhone
    m["station"] = sta.Name
    m["payMonth"] = 1
    m["payerEnt"] = true
}

// Sign 签署合同
func (s *contractService) Sign(u *ent.Rider, params *model.ContractSignReq) model.ContractSignRes {
    if s.Effective(u) {
        return model.ContractSignRes{Effective: true}
    }

    if u.Contact == nil {
        snag.Panic("未补充紧急联系人")
    }

    now := make([]int, 3)
    arr := strings.Split(time.Now().Format(carbon.DateLayout), "-")
    for i, a := range arr {
        now[i], _ = strconv.Atoi(a)
    }

    var (
        m            = make(ar.Map)
        sn           = tools.NewUnique().NewSN()
        cfg          = s.esign.Config
        p            = NewPerson().GetNormalAuthedPerson(u)
        accountId    = p.EsignAccountID
        isEnterprise = u.EnterpriseID != nil
        templateId   = cfg.Person.TemplateId
        scene        = cfg.Person.Scene
    )

    m["sn"] = sn
    m["name"] = p.Name
    m["riderName"] = p.Name
    m["signName"] = p.Name
    m["idcard"] = p.IDCardNumber
    m["address"] = p.AuthResult.Address
    m["phone"] = u.Phone
    m["startYear"] = now[0]
    m["startMonth"] = now[1]
    m["startDay"] = now[2]
    m["aurYear"] = now[0]
    m["aurMonth"] = now[1]
    m["aurDay"] = now[2]
    m["riderYear"] = now[0]
    m["riderMonth"] = now[1]
    m["riderDay"] = now[2]
    m["contactPhone"] = u.Contact.Phone

    // 电池型号
    if strings.HasPrefix(strings.ToUpper(params.Model), "72V") {
        m["v72"] = true
    } else {
        m["v60"] = true
    }

    if params.PlanID != nil {
        s.planData(*params.PlanID, m)
    }

    if isEnterprise {
        templateId = cfg.Group.TemplateId
        scene = cfg.Group.Scene
        s.enterpriseData(u, m)
    }

    // 创建 / 获取 签约个人账号
    if accountId == "" {
        accountId = s.esign.CreatePersonAccount(esign.CreatePersonAccountReq{
            ThirdPartyUserId: u.Phone,
            Name:             p.Name,
            IdType:           "CRED_PSN_CH_IDCARD",
            IdNumber:         p.IDCardNumber,
            Mobile:           u.Phone,
        })
        if accountId == "" {
            snag.Panic("签署账号生成失败")
        }
        // 保存个人账号
        err := p.Update().SetEsignAccountID(accountId).Exec(context.Background())
        if err != nil {
            snag.Panic(err)
        }
    }

    // 设置合同编号
    s.esign.SetSn(sn)

    // 设置个人账户ID
    req := esign.CreateFlowReq{
        Scene:           scene,
        PersonAccountId: accountId,
    }

    // 获取模板控件
    tmpl := s.esign.DocTemplate(templateId)
    for _, com := range tmpl.StructComponents {
        switch com.Key {
        case snKey:
            m[snKey] = sn
            break
        case aurSeal:
            req.EntSignBean = esign.PosBean{
                PosPage: fmt.Sprintf("%v", com.Context.Pos.Page),
                PosX:    com.Context.Pos.X,
                PosY:    com.Context.Pos.Y,
            }
            break
        case riderSeal:
            req.PsnSignBean = esign.PosBean{
                PosPage: fmt.Sprintf("%v", com.Context.Pos.Page),
                PosX:    com.Context.Pos.X,
                PosY:    com.Context.Pos.Y,
            }
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

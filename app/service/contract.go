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
    "math"
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
    now   struct {
        year  int
        month int
        day   int
    }
}

func NewContract() *contractService {
    s := &contractService{
        esign: esign.New(),
        orm:   ent.Database.Contract,
        ctx:   context.Background(),
    }
    now := time.Now()
    s.now.year = now.Year()
    s.now.month = int(now.Month())
    s.now.day = now.Day()
    return s
}

// Effective 查询骑手是否存在生效中的合同
// 当用户退租之后触发合同失效, 需要重新签订
func (s *contractService) Effective(u *ent.Rider) bool {
    if u.Contractual {
        return true
    }
    exists, _ := s.orm.QueryNotDeleted().Where(
        contract.RiderID(u.ID),
        contract.Status(model.ContractStatusSuccess.Value()),
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
    // 总月数
    m["month"] = fmt.Sprintf("%.0f", month)
    // 首次缴纳月数
    m["payMonth"] = fmt.Sprintf("%.0f", month)
    // 月租金
    m["price"] = fmt.Sprintf("%.2f", price)
    // 总租金
    m["amount"] = fmt.Sprintf("%.2f", p.Price)
    // 截止年月日
    end := tools.NewTime().WillEnd(time.Now(), int(p.Days))
    m["endYear"] = end.Year()
    m["endMonth"] = int(end.Month())
    m["endDay"] = end.Day()
}

func (s *contractService) enterpriseData(u *ent.Rider, m ar.Map, cityID uint64, bm string) {
    e := u.Edges.Enterprise
    if e == nil {
        e, _ = ent.Database.Enterprise.QueryNotDeleted().Where(enterprise.ID(*u.EnterpriseID)).First(s.ctx)
    }
    if e == nil {
        snag.Panic("骑手企业查找失败")
    }

    sta := u.Edges.Station
    if sta == nil {
        sta, _ = ent.Database.EnterpriseStation.QueryNotDeleted().Where(enterprisestation.ID(*u.StationID)).First(s.ctx)
    }
    if sta == nil {
        snag.Panic("骑手站点查找失败")
    }

    // 获取企业费用信息
    srv := NewEnterprise()
    prices := srv.GetPriceValues(e)
    pk := srv.PriceKey(cityID, bm)
    price, ok := prices[pk]
    if !ok {
        snag.Panic("企业费用查询失败")
    }

    // entName entPhone station payerEnt
    m["entName"] = e.CompanyName
    // 企业联系电话
    m["entPhone"] = e.ContactPhone
    // 站点
    m["station"] = sta.Name
    // 首次缴纳月数
    m["payMonth"] = 1
    // 企业付款
    m["payerEnt"] = true
    // 月租金
    m["price"] = fmt.Sprintf("%.2f", price)
    // 总月数
    m["month"] = 1
    // 总租金
    m["amount"] = fmt.Sprintf("%.2f", price*31.0)
    // 截止年月日
    end := tools.NewTime().WillEnd(time.Now(), 31)
    m["endYear"] = end.Year()
    m["endMonth"] = int(end.Month())
    m["endDay"] = end.Day()
}

// Sign 签署合同
func (s *contractService) Sign(u *ent.Rider, params *model.ContractSignReq) model.ContractSignRes {
    if s.Effective(u) {
        return model.ContractSignRes{Effective: true}
    }

    if u.Contact == nil {
        snag.Panic("未补充紧急联系人")
    }

    ci := NewCity().Query(params.CityID)

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
    m["startYear"] = s.now.year
    m["startMonth"] = s.now.month
    m["startDay"] = s.now.day
    m["aurYear"] = s.now.year
    m["aurMonth"] = s.now.month
    m["aurDay"] = s.now.day
    m["riderYear"] = s.now.year
    m["riderMonth"] = s.now.month
    m["riderDay"] = s.now.day
    // 紧急联系人
    m["contact"] = u.Contact.String()
    // 勾选租赁电池
    m["battery"] = true
    // 电池型号
    m["model"] = strings.ToUpper(params.Model)
    // 限制城市
    m["city"] = ci.Name

    if params.PlanID != 0 {
        s.planData(params.PlanID, m)
    }

    if isEnterprise {
        templateId = cfg.Group.TemplateId
        scene = cfg.Group.Scene
        s.enterpriseData(u, m, params.CityID, params.Model)
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
    err := ent.Database.Contract.Create().
        SetFlowID(flowId).
        SetRiderID(u.ID).
        SetStatus(model.ContractStatusPending.Value()).
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
    orm := ent.Database.Contract
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
        update.SetStatus(model.ContractStatusSuccess.Value()).
            SetFiles(s.esign.DownloadDocument(fmt.Sprintf("%s-%s/contracts/", p.Name, p.IDCardNumber), c.FlowID))
    }
    err = update.Exec(context.Background())
    if err != nil {
        snag.Panic(err)
    }
    return model.StatusResponse{Status: success}
}

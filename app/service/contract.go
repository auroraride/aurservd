// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/socket"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/contract"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "io"
    "math"
    "net/http"
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
    esign      *esign.Esign
    ctx        context.Context
    orm        *ent.ContractClient
    rider      *ent.Rider
    monthDays  int
    timeLayout string
}

func NewContract() *contractService {
    s := &contractService{
        esign:      esign.New(),
        orm:        ent.Database.Contract,
        ctx:        context.Background(),
        monthDays:  30,
        timeLayout: "2006年01月02日",
    }
    return s
}

func NewContractWithRider(u *ent.Rider) *contractService {
    s := NewContract()
    s.rider = u
    return s
}

// Effective 查询骑手是否存在生效中的合同
// 当用户退租之后触发合同失效, 需要重新签订
func (s *contractService) Effective(u *ent.Rider) bool {
    if u.Contractual {
        return true
    }
    exists, _ := s.orm.Query().Where(
        contract.RiderID(u.ID),
        contract.Status(model.ContractStatusSuccess.Value()),
        contract.Effective(true),
    ).Exist(s.ctx)
    return exists
}

// planData 个签合同数据
func (s *contractService) planData(sub *ent.Subscribe) *model.ContractSignUniversal {
    p, _ := sub.QueryPlan().First(s.ctx)
    if p == nil {
        snag.Panic("未找到骑士卡信息")
    }
    month := math.Round(float64(p.Days) / float64(s.monthDays))
    price := p.Price
    if month > 1 {
        price = price / month
    }

    return &model.ContractSignUniversal{
        Price: fmt.Sprintf("%.2f", price),
        Month: int(month),
        Total: fmt.Sprintf("%.2f", p.Price),
        Stop:  tools.NewTime().WillEnd(time.Now(), int(p.Days)).Format(s.timeLayout),
    }
}

// enterpriseData 团签合同数据
func (s *contractService) enterpriseData(m ar.Map, sub *ent.Subscribe) *model.ContractSignUniversal {
    if sub.BrandID != nil || sub.EbikeID != nil {
        snag.Panic("暂不支持团签")
    }

    // 查询团签
    ee, _ := sub.QueryEnterprise().First(s.ctx)
    if ee == nil {
        snag.Panic("团签信息查询失败")
    }

    // 查询站点
    es, _ := sub.QueryStation().First(s.ctx)
    if es == nil {
        snag.Panic("站点信息查询失败")
    }

    // 获取企业费用信息
    srv := NewEnterprise()
    prices := srv.GetPriceValues(ee)
    pk := srv.PriceKey(sub.CityID, sub.Model)
    price, ok := prices[pk]
    if !ok {
        snag.Panic("团签费用查询失败")
    }

    // 团签公司名称
    m["entName"] = ee.CompanyName
    // 团签负责人及电话
    m["entContact"] = fmt.Sprintf("%s，%s", ee.ContactName, ee.ContactPhone)
    // 团签隶属站点
    m["entStation"] = es.Name
    // 团签代缴
    m["payerEnt"] = true

    stop := tools.NewTime().WillEnd(time.Now(), s.monthDays).Format(s.timeLayout)
    month := 1
    days := float64(s.monthDays)

    // 如果是代理
    if ee.Agent {
        if sub.AgentEndAt == nil {
            snag.Panic("代理团签订阅日期错误")
        }
        days = float64(tools.NewTime().LastDaysToNow(*sub.AgentEndAt))
        month = int(math.Round(days / float64(s.monthDays)))
    }

    return &model.ContractSignUniversal{
        Price: fmt.Sprintf("%.2f", price),
        Month: month,
        Total: fmt.Sprintf("%.2f", price*days),
        Stop:  stop,
    }
}

// Sign 签署合同
// 月数按s.monthDays(30)天计算, 出现小数四舍五入
// TODO 电柜激活电池(需要注意判定库存是否充足)
func (s *contractService) Sign(req *model.ContractSignReq) model.ContractSignRes {
    u := s.rider
    // 是否免签
    if s.Effective(u) {
        return model.ContractSignRes{Effective: true}
    }

    // 是否有紧急联系人
    if u.Contact == nil {
        snag.Panic("未补充紧急联系人")
    }

    // 查找订阅
    sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(req.SubscribeID), subscribe.Status(model.SubscribeStatusInactive)).WithCity().First(s.ctx)
    if sub == nil {
        snag.Panic("未找到骑士卡")
    }

    if !sub.NeedContract {
        snag.Panic("当前订阅无需签约")
    }

    if sub.BrandID == nil && sub.EbikeID != nil {
        snag.Panic("当前订阅错误")
    }

    // 城市
    ec := sub.Edges.City
    if ec == nil {
        snag.Panic("未找到有效城市")
    }

    // 判定门店或电柜库存
    stockable := true
    if req.CabinetID != nil {
        // 判定电柜库存
        stockable = NewStock().CheckCabinet(*req.CabinetID, sub.Model, 1)
    }
    if req.StoreID != nil {
        // 判定门店库存
        stockable = NewStock().CheckStore(*req.StoreID, sub.Model, 1)
    }
    if !stockable {
        snag.Panic("电池库存不足")
    }

    // 定义变量
    var (
        m            = make(ar.Map)
        sn           = tools.NewUnique().NewSN()
        cfg          = s.esign.Config
        p            = NewPerson().GetNormalAuthedPerson(u)
        accountId    = p.EsignAccountID
        isEnterprise = u.EnterpriseID != nil
        templateId   = cfg.Person.TemplateId
        scene        = cfg.Person.Scene

        // 电池型号
        bm = strings.ToUpper(sub.Model)
        // 当前日期
        now = time.Now().Format(s.timeLayout)
    )

    // 填充公共变量
    // 合同编号
    m["sn"] = sn
    // 骑手姓名
    m["name"] = p.Name
    // 身份证号
    m["idcard"] = p.IDCardNumber
    // 户口地址
    m["address"] = p.AuthResult.Address
    // 骑手电话
    m["phone"] = u.Phone
    // 限制城市
    m["city"] = ec.Name
    // 骑手签字
    m["riderSign"] = p.Name
    // 紧急联系人
    m["riderContact"] = u.Contact.String()
    // 企业签署日期
    m["aurDate"] = now
    // 骑手签署日期
    m["riderDate"] = now

    var un *model.ContractSignUniversal

    if isEnterprise {
        // 团签
        templateId = cfg.Group.TemplateId
        scene = cfg.Group.Scene
        // 设置团签字段
        un = s.enterpriseData(m, sub)
        // 团签代缴
        m["payEnt"] = true
    } else {
        // 个签骑士卡
        un = s.planData(sub)
        // 骑手缴费
        m["payRider"] = true
    }

    if un == nil {
        snag.Panic("合同信息错误")
    }

    m["payMonth"] = un.Month

    // 电车
    var employeeID, storeID, allocateID *uint64
    if sub.BrandID != nil {
        // 查找电车分配
        ea := NewEbikeAllocate().QueryEffectiveSubscribeIDX(sub.ID)
        employeeID = silk.UInt64(ea.EmployeeID)
        storeID = silk.UInt64(ea.StoreID)
        allocateID = silk.UInt64(ea.ID)

        bike := ea.Info.Ebike
        // 车加电方案
        m["schemaEbike"] = true
        // 车加电方案一
        m["ebikeScheme1"] = true
        // 车辆品牌
        m["ebikeBrand"] = bike.Brand.Name
        // 车辆颜色
        m["ebikeColor"] = bike.Color
        // 车架号
        m["ebikeSN"] = bike.SN
        // 车牌号
        m["ebikePlate"] = bike.Plate
        // 电池类型
        m["ebikeBattery"] = "时光驹电池"
        // 电池规格
        m["ebikeModel"] = bm
        // 车电方案一开始日期
        m["ebikeScheme1Start"] = now
        // 车电方案一截止日
        m["ebikeScheme1Stop"] = un.Stop
        // 车电方案一月租金
        m["ebikeScheme1Price"] = un.Price
        // 车电方案一首次缴纳月数
        m["ebikeScheme1PayMonth"] = un.Month
        // 车电方案一首次缴纳租金
        m["ebikeScheme1PayTotal"] = un.Total
    } else {
        // 单电方案
        m["schemaBattery"] = true
        // 电池规格
        m["batteryModel"] = bm
        // 单电方案起租日
        m["batteryStart"] = now
        // 单电方案结束日
        m["batteryStop"] = un.Stop
        // 单电月租金
        m["batteryPrice"] = un.Price
        // 单电方案首次缴纳月数
        m["batteryPayMonth"] = un.Month
        // 单电方案首次缴纳租金
        m["batteryPayTotal"] = un.Total
    }

    // 个签选项
    if sub.PlanID != nil {
        s.planData(sub)
    }

    // 创建 / 获取 签约个人账号
    if accountId == "" {
        accountId = s.esign.CreatePersonAccount(esign.PersonAccountReq{
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
    flow := esign.CreateFlowReq{
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
            flow.EntSignBean = esign.PosBean{
                PosPage: fmt.Sprintf("%v", com.Context.Pos.Page),
                PosX:    com.Context.Pos.X,
                PosY:    com.Context.Pos.Y,
            }
            break
        case riderSeal:
            flow.PsnSignBean = esign.PosBean{
                PosPage: fmt.Sprintf("%v", com.Context.Pos.Page),
                PosX:    com.Context.Pos.X,
                PosY:    com.Context.Pos.Y,
            }
        }
    }
    // 填充内容生成PDF
    pdf := s.esign.CreateByTemplate(esign.CreateByTemplateReq{
        Name:             fmt.Sprintf("%s-%s.pdf", flow.Scene, sn),
        SimpleFormFields: m,
        TemplateId:       templateId,
    })
    flow.FileId = pdf.FileId

    // 发起签署，获取flowId
    flowId := s.esign.CreateFlowOneStep(flow)

    // 获取签署链接
    link := s.esign.ExecuteUrl(flowId, accountId)

    // 存储合同信息
    var c *ent.Contract
    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        c, err = tx.Contract.Create().
            SetFlowID(flowId).
            SetRiderID(u.ID).
            SetStatus(model.ContractStatusSigning.Value()).
            SetSn(sn).
            SetNillableStoreID(storeID).
            SetNillableEmployeeID(employeeID).
            SetNillableEbikeAllocateID(allocateID).
            SetSubscribe(sub).
            SetRiderInfo(&model.ContractRider{
                Phone:        u.Phone,
                Name:         u.Name,
                IDCardNumber: u.IDCardNumber,
            }).
            Save(context.Background())
        if err != nil {
            return
        }
        return sub.Update().SetNillableStoreID(storeID).SetNillableEmployeeID(employeeID).Exec(s.ctx)
    })

    // 监听合同签署结果
    go s.checkResult(c.FlowID)

    return model.ContractSignRes{
        Url: link,
        Sn:  sn,
    }
}

func (s *contractService) checkResult(flowID string) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    start := time.Now()
    for {
        select {
        case t := <-ticker.C:
            // 签署是否过期
            isExpired := t.Sub(start).Minutes() > model.ContractExpiration
            stop, _ := s.doResult(flowID, isExpired)
            if stop {
                ticker.Stop()
                return
            }
            if isExpired {
                ticker.Stop()
            }
        }
    }
}

func (s *contractService) doResult(flowID string, isExpired bool) (stop, success bool) {
    defer func() {
        if v := recover(); v != nil {
            stop = true
            log.Errorf("合同查询失败: %v", v)
            return
        }
    }()

    // 查询合同
    c, _ := s.orm.Query().Where(contract.FlowID(flowID)).WithRider().First(s.ctx)
    if c == nil {
        stop = true
        return
    }

    result := model.ContractStatus(c.Status)

    // 合同流程是否结束
    if result.IsFinished() {
        stop = true
        success = result.IsSuccessed()
        return
    }

    // 查询骑手信息
    r := c.RiderInfo
    if r == nil {
        stop = true
        return
    }

    // 查询合同流程状态
    result = s.esign.Result(c.FlowID)
    updater := s.orm.UpdateOneID(c.ID)

    // 是否过期
    if isExpired {
        result = model.ContractStatusExpired
        updater.SetStatus(result.Value())
    }

    // 是否成功
    success = result.IsSuccessed()

    if success {
        // 获取合同并上传到阿里云
        updater.SetStatus(model.ContractStatusSuccess.Value()).SetFiles(s.esign.DownloadDocument(fmt.Sprintf("%s-%s/contracts/", r.Name, r.IDCardNumber), c.FlowID))
    }

    // 流程是否终止
    if result.IsFinished() {
        stop = true
        err := updater.Exec(context.Background())
        if err != nil {
            log.Errorf("合同更新失败: %v", err)
            stop = true
            return
        }
    }

    // 成功签署合同
    if success {
        err := s.update(c)
        if err != nil {
            log.Errorf("合同更新失败 [id = %d] %v", c.ID, err)
        }

        // 如有必要, 通知店员合同签署完成
        if c.EmployeeID != nil && c.AllocateID != nil {
            socket.SendMessage(NewEmployeeSocket(), *c.EmployeeID, &model.EmployeeSocketMessage{
                Speech:          "骑手已签约",
                EbikeAllocateID: c.AllocateID,
            })
        }
    }

    return
}

// 关联更新
// 包含业务 [激活 / 业务 / 出入库]
func (s *contractService) update(c *ent.Contract) (err error) {
    defer func() {
        if v := recover(); v != nil {
            err = fmt.Errorf("%v", v)
            return
        }
    }()

    if c.SubscribeID == nil {
        return errors.New("合同未关联订阅")
    }

    info, sub := NewBusinessRider(c.Edges.Rider).Inactive(*c.SubscribeID)
    if sub == nil {
        return errors.New("需要更新订阅, 但是未找到订阅信息")
    }

    // 激活
    srv := NewBusinessRider(c.Edges.Rider).SetStoreID(c.StoreID).SetCabinetID(c.CabinetID)

    // 查询分配信息
    var ea *ent.EbikeAllocate
    if c.AllocateID != nil {
        ea, _ = c.QueryEbikeAllocate().First(s.ctx)
        if ea == nil {
            return errors.New("未找到分配信息")
        }
        // 设置门店和电车属性
        srv.SetEbike(&model.EbikeBusinessInfo{
            ID:        ea.EbikeID,
            BrandID:   ea.Info.Ebike.Brand.ID,
            BrandName: ea.Info.Ebike.Brand.Name,
        })
    }

    srv.Active(info, sub, func(tx *ent.Tx) (err error) {
        if c.AllocateID != nil {
            // 更新分配
            err = tx.EbikeAllocate.UpdateOne(ea).SetStatus(model.EbikeAllocateStatusSigned.Value()).Exec(s.ctx)
            if err != nil {
                return
            }
            // 更新电车
            err = tx.Ebike.UpdateOneID(ea.EbikeID).SetRiderID(sub.RiderID).SetStatus(model.EbikeStatusUsing).Exec(s.ctx)
        }
        return
    })

    return
}

// Result 合同签署结果
func (s *contractService) Result(r *ent.Rider, sn string) model.StatusResponse {
    // 查询合同是否存在
    c, err := s.orm.Query().
        Where(contract.Sn(sn), contract.RiderID(r.ID)).
        First(context.Background())
    if err != nil || c == nil {
        snag.Panic("合同查询失败")
    }

    return model.StatusResponse{Status: model.ContractStatus(c.Status).IsSuccessed()}
}

// Notice 签约回调
func (s *contractService) Notice(req *http.Request) {
    b, err := io.ReadAll(req.Body)
    if len(b) == 0 || err != nil {
        log.Errorf("签约回调读取失败: %v", err)
        return
    }

    // 解析回调信息
    var result esign.Notice
    err = jsoniter.Unmarshal(b, &result)
    if err != nil {
        log.Errorf("签约回调解析失败: %v", err)
        return
    }

    switch result.Action {
    case "SIGN_FLOW_FINISH":
        if result.FlowId != "" {
            s.doResult(result.FlowId, false)
        }
    }
}

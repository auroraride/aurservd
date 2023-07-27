// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/auroraride/adapter/log"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/socket"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/esign"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
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
	pk := srv.PriceKey(sub.CityID, sub.Model, sub.BrandID)
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
		if sub.InitialDays == 0 {
			snag.Panic("订阅天数错误错误")
		}
		days = float64(tools.NewTime().LastDaysToNow(tools.NewTime().WillEnd(time.Now(), sub.InitialDays)))
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
// 电柜激活电池(需要注意判定库存是否充足)
func (s *contractService) Sign(req *model.ContractSignReq) model.ContractSignRes {
	u := s.rider

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

	// 是否免签或已签约
	if !NewSubscribe().NeedContract(sub) {
		return model.ContractSignRes{Effective: true}
	}

	// 查找分配信息
	allo := NewAllocate().QueryEffectiveSubscribeIDX(sub.ID)
	if allo == nil {
		snag.Panic("未找到有效分配")
	}

	if sub.BrandID != nil && allo.StoreID == nil {
		snag.Panic("电车必须由门店分配")
	}

	// 城市
	ec := sub.Edges.City
	if ec == nil {
		snag.Panic("未找到有效城市")
	}

	// 判定非智能套餐门店库存
	if allo.StoreID != nil && allo.BatteryID == nil && !NewStock().CheckStore(*allo.StoreID, sub.Model, 1) {
		snag.Panic("电池库存不足")
	}

	var link, flowId, sn string
	skip := false
	co, _ := s.orm.QueryNotDeleted().Where(contract.SubscribeID(sub.ID), contract.LinkNotNil(), contract.Status(model.ContractStatusSigning.Value())).First(s.ctx)
	// 判定是否生成过合同
	if co != nil {
		// 合同处于有效期内跳过生成
		if co.ExpiresAt != nil && co.ExpiresAt.After(time.Now()) {
			skip = true
			link = *co.Link
			flowId = co.FlowID
			sn = co.Sn
		} else {
			// 否则删除原合同重新生成
			s.orm.DeleteOne(co).ExecX(s.ctx)
		}
	}

	// 生成合同
	if !skip {
		sn = tools.NewUnique().NewSN()

		// 定义变量
		var (
			m            = make(ar.Map)
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
		if sub.BrandID != nil {
			// 查找电车分配
			bike, _ := allo.QueryEbike().WithBrand().First(s.ctx)
			if bike == nil || bike.Edges.Brand == nil {
				snag.Panic("未找到电车信息")
			}

			brand := bike.Edges.Brand

			// 车加电方案
			m["schemaEbike"] = true
			// 车加电方案一
			m["ebikeScheme1"] = true
			// 车辆品牌
			m["ebikeBrand"] = brand.Name
			// 车辆颜色
			m["ebikeColor"] = bike.Color
			// 车架号
			m["ebikeSN"] = bike.Sn
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

		// 合同过期时间
		expiresAt := time.Now().Add(model.ContractExpiration * time.Minute)

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
			case aurSeal:
				flow.EntSignBean = esign.PosBean{
					PosPage: fmt.Sprintf("%v", com.Context.Pos.Page),
					PosX:    com.Context.Pos.X,
					PosY:    com.Context.Pos.Y,
				}
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
		flowId = s.esign.CreateFlowOneStep(flow, expiresAt)

		// 获取签署链接
		link = s.esign.ExecuteUrl(flowId, accountId)

		// 存储合同信息
		ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
			// 删除原有合同
			_, _ = tx.Contract.Delete().Where(contract.AllocateID(allo.ID)).Exec(s.ctx)
			err = tx.Contract.Create().
				SetFlowID(flowId).
				SetRiderID(u.ID).
				SetStatus(model.ContractStatusSigning.Value()).
				SetSn(sn).
				SetNillableEmployeeID(allo.EmployeeID).
				SetAllocateID(allo.ID).
				SetSubscribe(sub).
				SetRiderInfo(&model.ContractRider{
					Phone:        u.Phone,
					Name:         u.Name,
					IDCardNumber: u.IDCardNumber,
				}).
				SetLink(link).
				SetExpiresAt(expiresAt).
				Exec(context.Background())
			if err != nil {
				return
			}
			return sub.Update().UpdateTarget(allo.CabinetID, allo.StoreID, allo.EmployeeID).Exec(s.ctx)
		})

		// 监听合同签署结果
		go s.checkResult(flowId)
	}

	return model.ContractSignRes{
		Url: link,
		Sn:  sn,
	}
}

func (s *contractService) checkResult(flowID string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		t := <-ticker.C
		// 获取合同
		c, _ := s.orm.Query().Where(contract.FlowID(flowID)).WithRider().First(s.ctx)

		// 如果合同状态不是签署中, 直接结束
		if c == nil || c.Status != model.ContractStatusSigning.Value() {
			ticker.Stop()
			return
		}

		// 签署是否过期
		isExpired := t.Sub(c.CreatedAt).Minutes() > model.ContractExpiration+5
		// 如果已过期, 直接结束
		if isExpired {
			_ = c.Update().SetStatus(model.ContractStatusExpired.Value()).Exec(s.ctx)
			ticker.Stop()
			return
		}
	}
}

func (s *contractService) doResult(flowID string) {
	defer func() {
		if v := recover(); v != nil {
			zap.L().Error("合同查询失败", zap.Error(fmt.Errorf("%v", v)))
			return
		}
	}()

	// 查询合同
	c, _ := s.orm.Query().Where(contract.FlowID(flowID)).WithRider().First(s.ctx)
	if c == nil {
		return
	}

	result := model.ContractStatus(c.Status)

	// 合同流程是否结束
	if result.IsFinished() {
		return
	}

	// 查询骑手信息
	r := c.RiderInfo
	if r == nil {
		return
	}

	updater := s.orm.UpdateOneID(c.ID)

	var endAt time.Time

	// 查询合同流程状态
	var signResult *esign.SignResult
	result, signResult = s.esign.Result(c.FlowID)
	endAt = signResult.EndAt()

	// 是否成功
	if result.IsSuccessed() {
		// 获取合同并上传到阿里云
		updater.SetSignedAt(endAt).SetStatus(model.ContractStatusSuccess.Value()).SetFiles(s.esign.DownloadDocument(fmt.Sprintf("%s-%s/contracts/", r.Name, r.IDCardNumber), flowID))
	}

	idstr := strconv.FormatUint(c.ID, 10)

	// 流程是否终止
	if result.IsFinished() {
		err := updater.Exec(context.Background())
		if err != nil {
			zap.L().Error("合同更新失败: "+idstr, zap.Error(err))
			return
		}
	}

	// 成功签署合同
	if result.IsSuccessed() {
		err := s.update(c)
		if err != nil {
			zap.L().Error("已成功签署合同, 但合同更新失败: "+idstr, zap.Error(err))
		}

		// 如有必要, 通知店员合同签署完成
		if c.EmployeeID != nil && c.AllocateID != nil {
			socket.SendMessage(NewEmployeeSocket(), *c.EmployeeID, &model.EmployeeSocketMessage{
				Speech:          "骑手已签约",
				EbikeAllocateID: c.AllocateID,
			})
		}
	}
}

// 关联更新
// 包含业务 [激活 / 业务 / 出入库]
func (s *contractService) update(c *ent.Contract) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
			zap.L().Error("合同关联更新失败", zap.Error(err), zap.Uint64("id", c.ID))
			return
		}
	}()

	if c.SubscribeID == nil {
		return errors.New("合同未关联订阅")
	}

	srv := NewBusinessRider(c.Edges.Rider)

	_, sub := srv.Inactive(*c.SubscribeID)
	if sub == nil {
		return errors.New("需要更新订阅, 但是未找到订阅信息")
	}

	// 查询分配信息
	allo, _ := c.QueryAllocate().First(s.ctx)

	if allo == nil {
		return errors.New("未找到分配信息")
	}

	// 更新分配状态
	err = allo.Update().SetStatus(model.AllocateStatusSigned.Value()).Exec(s.ctx)

	// 激活
	srv.SetModifier(allo.LastModifier).
		SetStoreID(allo.StoreID).
		SetCabinetID(allo.CabinetID).
		SetEmployeeID(allo.EmployeeID).
		SetBatteryID(allo.BatteryID)

	// 设置电车属性
	if allo.EbikeID != nil {
		bike, _ := allo.QueryEbike().WithBrand().First(s.ctx)
		if bike == nil || bike.Edges.Brand == nil {
			return errors.New("未找到电车信息")
		}
		brand := bike.Edges.Brand
		srv.SetEbike(&model.EbikeBusinessInfo{
			ID:        bike.ID,
			BrandID:   brand.ID,
			BrandName: brand.Name,
		})
	}

	// 完成签约后
	// 若是门店或站点分配则自动并激活 (骑手扫码电柜无需激活)
	if allo.StoreID != nil || allo.StationID != nil {
		zap.L().Info("开始激活流程, allocateId = " + strconv.FormatUint(allo.ID, 10))
		srv.Active(sub, allo)
	}

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
func (s *contractService) Notice(b []byte) {
	// 解析回调信息
	var result esign.Notice
	err := jsoniter.Unmarshal(b, &result)
	if err != nil {
		zap.L().Error("签约回调解析失败", zap.Error(err), log.ResponseBody(b))
		return
	}

	switch result.Action {
	case "SIGN_FLOW_FINISH":
		if result.FlowId != "" {
			s.doResult(result.FlowId)
		}
	}
}

func (s *contractService) List(req *model.ContractListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Order(ent.Desc(contract.FieldCreatedAt))
	if req.Keyword != "" {
		q.Where(contract.HasRiderWith(rider.Or(rider.PhoneContainsFold(req.Keyword), rider.NameContainsFold(req.Keyword))))
	}
	if req.Status != nil {
		q.Where(contract.Status(*req.Status))
	}
	if req.Effective != nil {
		q.Where(contract.Effective(*req.Effective))
	}
	if req.Start != "" {
		q.Where(contract.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
	}
	if req.End != "" {
		q.Where(contract.CreatedAtLT(tools.NewTime().ParseNextDateStringX(req.End)))
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Contract) (res model.ContractListRes) {
		status := model.ContractStatus(item.Status)
		res = model.ContractListRes{
			Name:      item.RiderInfo.Name,
			Phone:     item.RiderInfo.Phone,
			Status:    status.String(),
			Effective: item.Effective,
			CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
		}
		// 文档链接
		if len(item.Files) > 0 {
			res.Link = item.Files[0]
		}

		if status == model.ContractStatusSuccess {
			// 如果签署成功, 显示签署时间
			if item.SignedAt != nil {
				res.SignAt = item.SignedAt.Format(carbon.DateTimeLayout)
			}
		} else {
			// 如果还未签署成功, 显示过期时间
			res.ExpiresAt = item.ExpiresAt.Format(carbon.DateTimeLayout)
		}

		return
	})
}

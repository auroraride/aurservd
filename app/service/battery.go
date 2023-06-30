// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-24
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/log"
	"github.com/auroraride/adapter/rpc/pb"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/rpc"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type batteryService struct {
	*BaseService
	orm *ent.BatteryClient
}

func NewBattery(params ...any) *batteryService {
	return &batteryService{
		BaseService: newService(params...),
		orm:         ent.Database.Battery,
	}
}

func (s *batteryService) QueryID(id uint64) (*ent.Battery, error) {
	return s.orm.Query().Where(battery.ID(id)).First(s.ctx)
}

func (s *batteryService) QueryIDX(id uint64) (b *ent.Battery) {
	b, _ = s.QueryID(id)
	if b == nil {
		snag.Panic("未找到电池")
	}
	return
}

func (s *batteryService) QueryRiderID(id uint64) (*ent.Battery, error) {
	return s.orm.Query().Where(battery.RiderID(id)).First(s.ctx)
}

func (s *batteryService) QueryRiderIDX(id uint64) (b *ent.Battery) {
	b, _ = s.QueryRiderID(id)
	if b == nil {
		snag.Panic("未找到电池")
	}
	return
}

func (s *batteryService) QuerySn(sn string) (bat *ent.Battery, err error) {
	return s.orm.Query().Where(battery.Sn(sn)).First(s.ctx)
}

func (s *batteryService) QuerySnX(sn string) (bat *ent.Battery) {
	bat, _ = s.QuerySn(sn)
	if bat == nil {
		snag.Panic("未找到电池")
	}
	return
}

// LoadOrCreate 加载电池, 若电池不存在则先创建电池, 若电池存在, 则不更新电池直接返回
func (s *batteryService) LoadOrCreate(sn string, params ...any) (bat *ent.Battery, err error) {
	bat, _ = s.QuerySn(sn)
	if bat != nil {
		return
	}

	var (
		cabID   *uint64
		ordinal *int
	)

	for _, param := range params {
		switch v := param.(type) {
		case *model.BatteryInCabinet:
			cabID = silk.UInt64(v.CabinetID)
			ordinal = silk.Int(v.Ordinal)
		}
	}

	// 解析电池型号
	ab, err := adapter.ParseBatterySN(sn)
	if err != nil || ab.Model == "" || ab.SN == "" {
		zap.L().Error("型号错误: "+sn, log.Payload(ab))
		return nil, adapter.ErrorBatterySN
	}

	return s.orm.Create().SetModel(ab.Model).SetSn(sn).SetBrand(ab.Brand).SetNillableCabinetID(cabID).SetNillableOrdinal(ordinal).Save(s.ctx)
}

// TODO 电池需要做库存管理

// Create 创建电池
func (s *batteryService) Create(req *model.BatteryCreateReq) {
	enable := true
	if req.Enable != nil {
		enable = *req.Enable
	}

	// 解析电池编号
	ab, err := adapter.ParseBatterySN(req.SN)
	if err != nil || ab.Model == "" {
		snag.Panic("电池编号解析失败, 请擦亮你的双眼")
	}
	_, err = s.orm.Create().
		SetSn(req.SN).
		SetBrand(ab.Brand).
		SetModel(ab.Model).
		SetEnable(enable).
		SetCityID(req.CityID).
		Save(s.ctx)
	if err != nil {
		snag.Panic("电池创建失败: " + err.Error())
	}
}

// BatchCreate 批量创建电池
// 0-城市:city 1-型号:model 2-编号:sn
func (s *batteryService) BatchCreate(c echo.Context) []string {
	rows, sns, failed := s.BaseService.GetXlsxRows(c, 2, 2, 2)
	// 查重
	items, _ := s.orm.Query().Where(battery.SnIn(sns...)).All(s.ctx)
	m := make(map[string]bool)
	for _, item := range items {
		m[item.Sn] = true
	}

	// 查询城市
	cs := make(map[string]struct{})
	for _, row := range rows {
		cs[row[0]] = struct{}{}
	}
	var (
		cids []string
		cm   = make(map[string]uint64)
	)
	for k := range cs {
		cids = append(cids, k)
	}
	cities, _ := ent.Database.City.Query().Where(city.NameIn(cids...)).All(s.ctx)
	for _, ci := range cities {
		cm[ci.Name] = ci.ID
	}

	for _, row := range rows {
		sn := row[1]
		if m[sn] {
			failed = append(failed, fmt.Sprintf("编号%s已存在", sn))
			continue
		}

		// 解析电池编号
		ab, err := adapter.ParseBatterySN(sn)
		if err != nil || ab.Model == "" {
			failed = append(failed, fmt.Sprintf("电池编号%s解析失败, 请擦亮你的双眼", sn))
			continue
		}

		creator := s.orm.Create()

		// 城市
		if cid, ok := cm[row[0]]; ok {
			creator.SetCityID(cid)
		} else {
			failed = append(failed, fmt.Sprintf("城市%s查询失败", row[0]))
			continue
		}

		_, err = creator.SetModel(ab.Model).SetBrand(ab.Brand).SetSn(sn).Save(s.ctx)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s保存失败: %v", sn, err))
		}

	}

	return failed
}

func (s *batteryService) Modify(req *model.BatteryModifyReq) {
	// 查找电池
	b := s.QueryIDX(req.ID)
	u := b.Update()
	if req.Enable != nil {
		u.SetEnable(*req.Enable)
	}
	if req.CityID != nil {
		u.SetCityID(*req.CityID)
	}
	_, _ = u.Save(s.ctx)
}

func (s *batteryService) listFilter(req model.BatteryFilter) (q *ent.BatteryQuery, info ar.Map) {
	q = s.orm.Query().WithRider().WithCity().WithCabinet().WithStation().WithEnterprise().Order(ent.Desc(battery.FieldCreatedAt))
	info = make(ar.Map)

	var (
		status     = 1
		statusText = map[int]string{
			0: "全部",
			1: "启用",
			2: "禁用",
		}
	)
	if req.Status != nil {
		status = *req.Status
	}
	info["状态"] = statusText[status]
	switch status {
	case 1:
		q.Where(battery.Enable(true))
	case 2:
		q.Where(battery.Enable(false))
	default:
		info["状态"] = "-"
	}

	if req.CityID != 0 {
		info["城市"] = ent.NewExportInfo(req.CityID, city.Table)
		q.Where(battery.CityID(req.CityID))
	}

	if req.SN != "" {
		info["编号"] = req.SN
		q.Where(battery.SnContainsFold(req.SN))
	}

	if req.Model != "" {
		m := strings.ToUpper(req.Model)
		info["型号"] = m
		q.Where(battery.Model(m))
	}
	if req.EnterpriseID != nil {
		info["团签"] = ent.NewExportInfo(*req.EnterpriseID, enterprise.Table)
		q.Where(battery.EnterpriseID(*req.EnterpriseID))
	}
	if req.StationID != nil {
		info["站点"] = ent.NewExportInfo(*req.StationID, enterprisestation.Table)
		q.Where(battery.StationID(*req.StationID))
	}
	if req.CabinetName != nil {
		info["电柜名称"] = *req.CabinetName
		q.Where(battery.HasCabinetWith(cabinet.Name(*req.CabinetName)))
	}
	if req.Keyword != nil {
		info["关键字"] = *req.Keyword
		q.Where(battery.HasRiderWith(rider.NameContainsFold(*req.Keyword)))
	}

	if req.OwnerType != nil {
		var name string
		switch *req.OwnerType {
		case 1: // 平台
			q.Where(battery.EnterpriseIDIsNil())
			name = "平台"
		case 2: // 代理商
			q.Where(battery.EnterpriseIDNotNil())
			name = "代理商"
		}
		info["归属"] = name
	}

	if req.CabinetID != nil {
		info["电柜"] = ent.NewExportInfo(*req.CabinetID, cabinet.Table)
		q.Where(battery.CabinetID(*req.CabinetID))
	}

	if req.RiderID != nil {
		info["骑手"] = ent.NewExportInfo(*req.RiderID, rider.Table)
		q.Where(battery.RiderID(*req.RiderID))
	}

	switch req.Goal {
	case model.BatteryStation:
		info["查询目标"] = "站点"
		q.Where(
			battery.StationIDNotNil(),
			battery.CabinetIDIsNil(),
			battery.RiderIDIsNil(),
		)
	case model.BatteryCabinet:
		info["查询目标"] = "电柜"
		q.Where(
			battery.CabinetIDNotNil(),
			battery.RiderIDIsNil(),
		)
	case model.BatteryRider:
		info["查询目标"] = "骑手"
		q.Where(
			battery.RiderIDNotNil(),
			battery.CabinetIDIsNil(),
		)
	}

	return
}

func (s *batteryService) List(req *model.BatteryListReq) (res *model.PaginationRes) {
	q, _ := s.listFilter(req.BatteryFilter)
	snmap := make(map[adapter.BatteryBrand][]string)
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Battery) (res *model.BatteryListRes) {
		snmap[item.Brand] = append(snmap[item.Brand], item.Sn)
		res = &model.BatteryListRes{
			ID:     item.ID,
			Brand:  item.Brand,
			Model:  item.Model,
			Enable: item.Enable,
			SN:     item.Sn,
		}

		c := item.Edges.City
		if c != nil {
			res.City = &model.City{
				ID:   c.ID,
				Name: c.Name,
			}
		}

		r := item.Edges.Rider
		if r != nil {
			res.Rider = &model.Rider{
				ID:    r.ID,
				Phone: r.Phone,
				Name:  r.Name,
			}
		}

		cab := item.Edges.Cabinet
		if cab != nil {
			res.Cabinet = &model.CabinetBasicInfo{
				ID:     cab.ID,
				Brand:  cab.Brand,
				Serial: cab.Serial,
				Name:   cab.Name,
			}
		}
		if item.Edges.Station != nil {
			res.StationName = item.Edges.Station.Name
		}
		if item.Edges.Enterprise != nil {
			res.EnterpriseName = item.Edges.Enterprise.Name
		}
		return
	})

	// 请求bms rpc
	result := make(map[string]*pb.BatteryItem)
	for br, list := range snmap {
		r := rpc.BmsBatch(br, &pb.BatteryBatchRequest{Sn: list})
		if r == nil {
			continue
		}

		for _, rb := range r.Items {
			result[rb.Sn] = rb
		}
	}

	for _, data := range res.Items.([]*model.BatteryListRes) {
		if rb, ok := result[data.SN]; ok {
			if len(rb.Heartbeats) > 0 {
				data.BmsBattery = model.NewBmsBattery(rb.Heartbeats[0])
			}
		}
	}

	return
}

// RiderBusiness 骑手业务操作电池
func (s *batteryService) RiderBusiness(putin bool, sn string, r *model.Rider, cab *ent.Cabinet, ordinal int) {
	var before, after string
	var op model.Operate

	target := "电池: " + sn + ", 电柜: " + cab.Serial + ", " + strconv.Itoa(ordinal) + " 号仓"

	if putin {
		op = model.OperateRiderPutin
		after = target
	} else {
		op = model.OperateRiderPutout
		before = target
	}

	go logging.NewOperateLog().
		SetOperate(op).
		SetCabinet(&model.CabinetBasicInfo{
			ID:     cab.ID,
			Brand:  cab.Brand,
			Serial: cab.Serial,
			Name:   cab.Name,
		}).
		SetRefManually(rider.Table, r.ID).
		SetDiff(before, after).
		Send()
}

// RiderDetail 获取电池详情
func (s *batteryService) RiderDetail(riderID uint64) (res model.BatteryDetail) {
	bat, _ := s.QueryRiderID(riderID)
	if bat != nil {
		// 获取电量信息
		var soc float64 = 0
		sr := rpc.BmsSample(bat.Brand, &pb.BatterySnRequest{Sn: bat.Sn})
		if sr != nil {
			soc = float64(sr.Soc)
		}
		res = model.BatteryDetail{
			ID:    bat.ID,
			Model: bat.Model,
			SN:    bat.Sn,
			Soc:   soc,
		}
	}
	return
}

// BindRequest 绑定骑手
func (s *batteryService) BindRequest(req *model.BatteryBind) {
	// 查找订阅
	sub := NewSubscribe().QueryEffectiveIntelligentX(req.RiderID, ent.SubscribeQueryWithBattery, ent.SubscribeQueryWithRider)

	// 查找电池
	bat := NewBattery().QueryIDX(req.BatteryID)
	// 查看是否冲突
	if (bat.RiderID != nil && *bat.RiderID != sub.RiderID) || (bat.SubscribeID != nil && *bat.SubscribeID != sub.ID) {
		snag.Panic("当前电池有绑定中的骑手, 无法重复绑定")
	}

	if bat.CabinetID != nil {
		snag.Panic("电柜中的电池无法手动绑定骑手")
	}

	s.Bind(bat, sub, sub.Edges.Rider)
}

func (s *batteryService) Bind(bat *ent.Battery, sub *ent.Subscribe, rd *ent.Rider) {
	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		return s.Allocate(tx, bat, sub, false)
	})

	go logging.NewOperateLog().
		SetOperate(model.OperateBindBattery).
		SetRef(rd).
		SetDiff("", "新电池: "+bat.Sn).
		SetModifier(s.modifier).
		Send()
}

// Unbind 解绑电池
func (s *batteryService) Unbind(req *model.BatteryUnbindRequest) {
	// 查找订阅
	sub := NewSubscribe().QueryEffectiveIntelligentX(req.RiderID, ent.SubscribeQueryWithBattery, ent.SubscribeQueryWithRider, ent.SubscribeQueryWithBattery)

	bat := sub.Edges.Battery
	if bat == nil {
		snag.Panic("未找到绑定的电池")
	}

	err := s.Unallocate(bat.Update())
	if err != nil {
		snag.Panic(err)
	}

	go logging.NewOperateLog().
		SetOperate(model.OperateUnbindBattery).
		SetDiff("旧电池: "+bat.Sn, "无电池").
		SetModifier(s.modifier).
		SetRef(sub.Edges.Rider).
		Send()
}

// Allocate 将电池分配给骑手
// @param createFlowIfError 即使遇到错误也创建电池流转信息
func (s *batteryService) Allocate(tx *ent.Tx, bat *ent.Battery, sub *ent.Subscribe, createFlowIfError bool) (err error) {
	// 删除原有骑手电池信息
	err = tx.Battery.Update().Where(battery.SubscribeID(sub.ID)).ClearRiderID().ClearSubscribeID().Exec(s.ctx)
	if err != nil {
		return
	}

	// 分配电池给骑手
	updater := tx.Battery.UpdateOne(bat).SetRiderID(sub.RiderID).SetSubscribeID(sub.ID)

	if sub.StationID != nil {
		// 当前骑手属于代理站点时, 设置新的站点信息
		updater.SetNillableStationID(sub.StationID).SetNillableEnterpriseID(sub.EnterpriseID)
	} else {
		// 当前骑手属于平台时, 清除原有站点信息
		updater.ClearStationID().ClearEnterpriseID()
	}

	err = updater.Exec(s.ctx)

	// 如果无错误或忽略错误，更新流转
	if createFlowIfError || err == nil {
		NewBatteryFlow().Create(tx, bat, model.BatteryFlowCreateReq{
			RiderID:     silk.Pointer(sub.RiderID),
			SubscribeID: silk.Pointer(sub.ID),
		})
	}

	return
}

// Unallocate 清除骑手信息
func (s *batteryService) Unallocate(updater *ent.BatteryUpdateOne) (err error) {
	return updater.ClearSubscribeID().ClearRiderID().Exec(s.ctx)
}

// StationBusinessTransfer 站点之间业务自动调拨
// 目前仅有换电业务
// 骑手从电柜中取出电池, 并将自己的电池放入电柜中, 因此:
// sub.StationID / sub.EnterpriseID 被用作放入的代理商信息
// cab.StationID / cab.EnterpriseID 被用作取出的代理商信息
// 需要记录流转信息
func (s *batteryService) StationBusinessTransfer(cabinetID, exchangeID uint64, putin, putout *model.BatteryEnterpriseTransfer) {
	// 进行站点对比, 放入 == 取出, 直接跳过
	if putin.StationID == putout.StationID {
		return
	}

	// 放入电池
	in, _ := NewBattery().QuerySn(putin.Sn)

	// 取出电池
	out, _ := NewBattery().QuerySn(putout.Sn)

	// 未找到电池跳过
	if in == nil || out == nil {
		return
	}

	// 若非站点骑手取出站点电池, 需要更新
	// 若放入是其他站点的电池, 其本质是两个站点(代理)的电池互换

	// 放入到该站点的电池
	s.updateStation(in, putout.StationID, putout.EnterpriseID)

	// 从该站点取出的电池
	s.updateStation(out, putin.StationID, putin.EnterpriseID)

	// 记录
	err := ent.Database.EnterpriseBatterySwap.Create().
		SetCabinetID(cabinetID).
		SetExchangeID(exchangeID).
		SetPutinID(in.ID).
		SetPutinSn(in.Sn).
		SetNillablePutinEnterpriseID(putin.EnterpriseID).
		SetNillablePutinStationID(putin.StationID).
		SetPutoutID(out.ID).
		SetPutoutSn(out.Sn).
		SetNillablePutoutEnterpriseID(putout.EnterpriseID).
		SetNillablePutoutStationID(putout.StationID).
		Exec(s.ctx)
	if err != nil {
		inb, _ := jsoniter.Marshal(putin)
		outb, _ := jsoniter.Marshal(putout)
		zap.L().Error("电池交换记录失败", zap.Error(err), zap.ByteString("putin", inb), zap.ByteString("putout", outb))
	}
}

// 更新电池站点信息
func (s *batteryService) updateStation(bat *ent.Battery, stationID, enterpriseID *uint64) {
	updater := s.orm.UpdateOne(bat)
	switch {
	default:
		// 非站点电池, 清除站点和团签信息
		updater.ClearStationID().ClearEnterpriseID()
	case stationID != nil:
		// 站点电池, 记录站点和团签信息
		updater.SetNillableStationID(stationID).SetNillableEnterpriseID(enterpriseID)
	}
	err := updater.Exec(s.ctx)
	if err != nil {
		zap.L().Error("电池流转更新失败", zap.Error(err))
	}
}

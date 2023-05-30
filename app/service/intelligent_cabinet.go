// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/defs/cabdef"
	"github.com/auroraride/adapter/log"
	"github.com/go-resty/resty/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type intelligentCabinetService struct {
	*BaseService
}

func NewIntelligentCabinet(params ...any) *intelligentCabinetService {
	return &intelligentCabinetService{
		BaseService: newService(params...),
	}
}

// ExchangeUsable 获取换电信息
func (s *intelligentCabinetService) ExchangeUsable(bm string, cab *ent.Cabinet) (uid string, info *model.RiderCabinetOperateProcess) {
	payload := &cabdef.ExchangeUsableRequest{
		Serial: cab.Serial,
		Minsoc: cache.Float64(model.SettingExchangeMinBatteryKey),
		Lock:   10,
		Model:  bm,
	}

	v, err := adapter.Post[cabdef.CabinetBinUsableResponse](s.GetCabinetAdapterUrlX(cab, "/exchange/usable"), s.GetAdapterUserX(), payload)

	if err != nil {
		snag.Panic(err)
	}

	info = &model.RiderCabinetOperateProcess{
		EmptyBin: &model.BinInfo{
			Index: v.Empty.Ordinal - 1,
		},
	}

	fully := &model.BinInfo{
		Index:       v.Fully.Ordinal - 1,
		Voltage:     v.Fully.Voltage,
		Electricity: model.BatterySoc(v.Fully.Soc),
	}

	if v.Fully.Soc >= model.IntelligentBatteryFullSoc {
		info.FullBin = fully
	} else {
		info.Alternative = fully
	}

	uid = v.UUID

	return
}

func (s *intelligentCabinetService) exchangeCacheKey(uid string) string {
	return "INTELLIGENT-CABINET-EXCHANGE-" + uid
}

// Exchange 请求换电
func (s *intelligentCabinetService) Exchange(uid string, ex *ent.Exchange, sub *ent.Subscribe, old *ent.Battery, cab *ent.Cabinet) {
	id, err := uuid.Parse(uid)
	if err != nil || (cab.Intelligent && old == nil) {
		snag.Panic("请求参数错误")
	}

	var (
		stopAt   *time.Time
		duration float64
		success  bool
		putout   string
		putin    string
		empty    *model.BinInfo
		bs       = NewBattery()
	)

	defer func() {
		updater := ex.Update()

		if err != nil {
			updater.SetRemark(err.Error())

			// 若换电失败, 标记任务失败
			_ = cache.Set(s.ctx, s.exchangeCacheKey(uid), &model.ExchangeStepResultCache{
				Index: 0,
				Results: []*cabdef.ExchangeStepMessage{
					{Step: 1, Message: err.Error(), Success: false},
				},
			}, 10*time.Minute).Err()
		}

		if stopAt == nil {
			stopAt = silk.Pointer(time.Now())
			duration = stopAt.Sub(ex.StartAt).Seconds()
		}

		if empty != nil {
			updater.SetEmpty(empty)
			// ex.Info.Empty = empty
		}

		// 保存数据库
		_ = updater.
			SetSuccess(success).
			SetFinishAt(*stopAt).
			SetDuration(int(duration)).
			SetPutoutBattery(putout).
			SetPutinBattery(putin).
			// SetInfo(ex.Info).
			Exec(s.ctx)
	}()

	payload := &cabdef.ExchangeRequest{
		UUID:    id,
		Serial:  cab.Serial,
		Expires: model.CabinetBusinessScanExpires,
		Timeout: model.CabinetBusinessStepTimeout,
		Minsoc:  cache.Float64(model.SettingExchangeMinBatteryKey),
	}

	if old != nil {
		payload.Battery = old.Sn
	}

	var v cabdef.ExchangeResponse
	v, err = adapter.Post[cabdef.ExchangeResponse](s.GetCabinetAdapterUrlX(cab, "/exchange/do"), s.GetAdapterUserX(), payload, func(r *resty.Response) {
		zap.L().Info("换电请求完成", log.ResponseBody(r.Body()))
	})

	if err != nil {
		zap.L().Error("换电请求失败", zap.Error(err))
		return
	}

	// 查询结果
	for _, result := range v.Results {
		duration += result.Duration
		stopAt = result.StopAt

		// 如果成功并且是智能柜, 记录电池编码
		if result.Success && cab.Intelligent {
			after := result.After
			before := result.Before

			// 记录用户放入的电池
			if model.ExchangeStepPutInto.EqualInt(result.Step) && after != nil {
				putin = after.BatterySN
				empty = &model.BinInfo{
					Index:       after.Ordinal - 1,
					Electricity: model.BatterySoc(after.Current),
					Voltage:     after.Voltage,
				}

				// 清除旧电池分配信息
				_ = NewBattery().Unallocate(old)

				go bs.RiderBusiness(true, putin, s.rider, cab, after.Ordinal)
			}

			// 记录用户取走的电池
			// 判定第三步是否成功, 只要柜门开启就把电池绑定到骑手 - BY: 曹博文
			if result.Step == model.ExchangeStepOpenFull.Int() && before != nil {
				putout = before.BatterySN

				go bs.RiderBusiness(false, putout, s.rider, cab, before.Ordinal)

				// 更新新电池信息
				bat, _ := bs.LoadOrCreate(putout)
				if bat != nil {
					_ = NewBattery().Allocate(bat.Update(), bat, sub, true)
				}
			}
		}
	}

	// 若换电成功 直接返回
	if v.Success {
		success = true
		return
	}

	if v.Error != "" {
		err = fmt.Errorf("%s", v.Error)
	}
}

// ExchangeStepSync 换电步骤同步
func (s *intelligentCabinetService) ExchangeStepSync(items []*cabdef.ExchangeStepMessage) {
	for _, req := range items {
		if req.Step == 0 {
			return
		}

		// TODO 检查电池是否存在???
		key := s.exchangeCacheKey(req.UUID)

		c := &model.ExchangeStepResultCache{}
		_ = cache.Get(s.ctx, key).Scan(c)
		c.Results = append(c.Results, req)

		// 排序
		slices.SortFunc(c.Results, func(a, b *cabdef.ExchangeStepMessage) bool {
			return a.Step <= b.Step
		})

		err := cache.Set(s.ctx, key, c, 10*time.Minute).Err()
		if err != nil {
			return
		}
	}
}

// ExchangeResult 查询换电结果
func (s *intelligentCabinetService) ExchangeResult(uid string) (res *model.RiderExchangeProcessRes) {
	key := s.exchangeCacheKey(uid)
	res = &model.RiderExchangeProcessRes{
		Step:   uint8(model.ExchangeStepOpenEmpty),
		Status: uint8(model.TaskStatusProcessing),
	}

	start := time.Now()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		if time.Since(start).Seconds() > 30 {
			return
		}

		c := &model.ExchangeStepResultCache{}
		err := cache.Get(s.ctx, key).Scan(c)
		if err != nil {
			continue
		}

		n := len(c.Results)
		index := c.Index

		// 当前的数据
		if index == n {
			s.exchangeStepResultFromCache(index-1, c, res)
			if res.Step < uint8(model.ExchangeStepPutOut) {
				res.Step += 1
			}
			res.Status = uint8(model.TaskStatusProcessing)
		}

		if index >= n {
			continue
		}

		s.exchangeStepResultFromCache(index, c, res)

		if !res.Stop {
			ttl, _ := ar.Redis.TTL(s.ctx, key).Result()
			c.Index += 1
			cache.Set(s.ctx, key, c, ttl)
		}
		return
	}

	return
}

// stepResult 获取换电步骤结果
func (s *intelligentCabinetService) exchangeStepResultFromCache(index int, c *model.ExchangeStepResultCache, res *model.RiderExchangeProcessRes) {
	data := c.Results[index]
	res.Step = uint8(data.Step)
	res.Message = data.Message
	if data.Success {
		res.Status = uint8(model.TaskStatusSuccess)
	}

	res.Stop = data.Step == 4 || !data.Success
}

// BusinessCensorX 校验用户是否可以使用智能柜办理业务
func (s *intelligentCabinetService) BusinessCensorX(bus adapter.Business, sub *ent.Subscribe, cab *ent.Cabinet) (bat *ent.Battery) {
	if !cab.Intelligent {
		return
	}

	// 判定电柜状态
	if cab.Status == model.CabinetStatusMaintenance.Value() {
		snag.Panic("电柜开小差了, 请联系客服")
	}

	// 判定是否智能电柜套餐
	if !sub.Intelligent {
		snag.Panic("套餐不匹配")
	}

	// 获取电池
	bat, _ = sub.QueryBattery().First(s.ctx)

	// 业务如果需要电池, 查找电池信息
	if bus.BatteryNeed() {
		// 未找到当前绑定的电池信息
		if bat == nil {
			snag.Panic(adapter.ErrorBatteryNotFound)
		}

		// 检查电池型号与电柜型号兼容
		if !NewCabinet().ModelInclude(cab, bat.Model) {
			snag.Panic("电池型号不兼容")
		}
	}

	return
}

// BusinessUsable 获取可用的业务仓位信息
func (s *intelligentCabinetService) BusinessUsable(cab *ent.Cabinet, bus adapter.Business, bm string) (uid string, index int, err error) {
	payload := &cabdef.BusinuessUsableRequest{
		Minsoc:   cache.Float64(model.SettingExchangeMinBatteryKey),
		Business: bus,
		Serial:   cab.Serial,
		Model:    bm,
	}

	var v cabdef.CabinetBinUsableResponse
	v, err = adapter.Post[cabdef.CabinetBinUsableResponse](s.GetCabinetAdapterUrlX(cab, "/business/usable"), s.GetAdapterUserX(), payload)
	if err != nil {
		return
	}

	uid = v.UUID
	index = v.BusinessBin.Ordinal - 1
	return
}

// DoBusiness 请求办理业务
func (s *intelligentCabinetService) DoBusiness(uidstr string, bus adapter.Business, sub *ent.Subscribe, riderBat *ent.Battery, cab *ent.Cabinet) (info *model.BinInfo, batinfo *model.Battery, err error) {
	defer func() {
		// 缓存任务返回
		data := &model.BusinessCabinetStatusRes{
			Success: err == nil,
			Stop:    true,
		}
		if err != nil {
			data.Message = err.Error()
		}
		cache.Set(s.ctx, uidstr, data, 10*time.Minute)
	}()

	var uid uuid.UUID
	uid, err = uuid.Parse(uidstr)
	if err != nil {
		return
	}

	var batterySN string
	if riderBat != nil {
		batterySN = riderBat.Sn
	}

	payload := &cabdef.BusinessRequest{
		UUID:     uid,
		Business: bus,
		Serial:   cab.Serial,
		Timeout:  model.CabinetBusinessStepTimeout,
		Battery:  batterySN,
		Model:    sub.Model,
	}

	var v cabdef.BusinessResponse
	v, err = adapter.Post[cabdef.BusinessResponse](s.GetCabinetAdapterUrlX(cab, "/business/do"), s.GetAdapterUserX(), payload)

	if err != nil {
		return
	}

	// TODO 失败后电池信息是否更新
	if v.Error != "" {
		err = errors.New(v.Error)
		return
	}

	var sn string
	var putin bool
	results := v.Results

	switch bus {
	case adapter.BusinessActive, adapter.BusinessContinue:
		sn = results[0].Before.BatterySN
	case adapter.BusinessPause, adapter.BusinessUnsubscribe:
		sn = results[1].After.BatterySN
		putin = true
	}

	b := results[1].After
	info = &model.BinInfo{
		Index:       b.Ordinal - 1,
		Electricity: model.BatterySoc(b.Soc),
		Voltage:     b.Voltage,
	}

	// 若智能电柜, 需记录电池信息
	if cab.Intelligent {
		// 获取电池
		var bat *ent.Battery
		bat, err = NewBattery().LoadOrCreate(sn)
		if err != nil {
			zap.L().Error("业务记录失败", zap.Error(err))
			return
		}

		batinfo = &model.Battery{
			ID:    bat.ID,
			SN:    sn,
			Model: bat.Model,
		}

		// 放入电池
		// TODO 是否有必要?
		// if putin {
		//     _, _ = bs.Unallocate(bat)
		// }

		// 取走电池
		if !putin {
			_ = NewBattery().Allocate(bat.Update(), bat, sub, true)
		}
	}

	return
}

func (s *intelligentCabinetService) Operate(cab *ent.Cabinet, op cabdef.Operate, req *model.CabinetDoorOperateReq) (success bool) {
	now := time.Now()
	br := cab.Brand
	ordinal := *req.Index + 1

	go func() {
		// 上传日志
		dlog := &logging.DoorOperateLog{
			ID:            shortuuid.New(),
			Brand:         br.String(),
			OperatorName:  s.modifier.Name,
			OperatorID:    s.modifier.ID,
			OperatorPhone: s.modifier.Phone,
			Serial:        cab.Serial,
			Name:          strconv.Itoa(ordinal) + "号仓",
			Operation:     req.Operation.String(),
			OperatorRole:  model.CabinetDoorOperatorRoleManager,
			Success:       success,
			Remark:        req.Remark,
			Time:          now.Format(carbon.DateTimeLayout),
		}
		dlog.Send()
	}()

	payload := &cabdef.OperateBinRequest{
		Operate: op,
		Ordinal: silk.Int(ordinal),
		Serial:  cab.Serial,
		Remark:  req.Remark,
	}

	_, err := adapter.Post[[]*cabdef.BinOperateResult](s.GetCabinetAdapterUrlX(cab, "/operate/bin"), s.GetAdapterUserX(), payload)

	success = err == nil
	return
}

func (s *intelligentCabinetService) Deactivate(cab *ent.Cabinet, payload *cabdef.BinDeactivateRequest) (success bool) {
	if s.modifier == nil {
		snag.Panic("权限校验失败")
	}
	now := time.Now()
	br := cab.Brand

	operation := "启用仓位"
	if *payload.Deactivate {
		operation = "禁用仓位"
	}

	res, _ := adapter.Post[model.StatusResponse](s.GetCabinetAdapterUrlX(cab, "/bin/deactivate"), s.GetAdapterUserX(), payload)

	success = res.Status

	if success {
		// 上传日志
		dlog := &logging.DoorOperateLog{
			ID:            shortuuid.New(),
			Brand:         br.String(),
			OperatorName:  s.modifier.Name,
			OperatorID:    s.modifier.ID,
			OperatorPhone: s.modifier.Phone,
			Serial:        cab.Serial,
			Name:          strconv.Itoa(payload.Ordinal) + "号仓",
			Operation:     operation,
			OperatorRole:  model.CabinetDoorOperatorRoleManager,
			Success:       success,
			Remark:        *payload.Reason,
			Time:          now.Format(carbon.DateTimeLayout),
		}
		go dlog.Send()
	}
	return
}

// OpenBind 开电池仓并绑定骑手
func (s *intelligentCabinetService) OpenBind(req *model.CabinetOpenBindReq) {
	bs := NewBattery(s.modifier)
	// 查找骑手
	rd := NewRider().QueryPhoneX(req.Phone)
	// 查找订阅
	sub := NewSubscribe().QueryEffectiveIntelligentX(rd.ID, ent.SubscribeQueryWithBattery, ent.SubscribeQueryWithRider)
	if !sub.Intelligent {
		snag.Panic("非智能电柜套餐, 无法操作")
	}
	// 查询电柜
	cab := NewCabinet().QueryOne(req.ID)
	if !cab.Intelligent {
		snag.Panic("非智能电柜, 无法操作")
	}
	// 查询电柜最新信息
	info, _ := s.Bininfo(cab, *req.Index+1)
	if info == nil {
		snag.Panic("获取最新仓位信息失败")
	}
	if info.BatterySN != req.BatterySN {
		snag.Panic("电池编码有变动, 请刷新后重试")
	}
	// 判定
	if exists, _ := sub.QueryBattery().Where().Exist(s.ctx); exists {
		snag.Panic("该骑手当前有绑定的电池")
	}
	// 查找电池
	bat := bs.QuerySnX(req.BatterySN)
	// 开门
	success := s.Operate(cab, cabdef.OperateDoorOpen, &model.CabinetDoorOperateReq{
		ID:        req.ID,
		Index:     req.Index,
		Remark:    req.Remark,
		Operation: silk.Pointer(model.CabinetDoorOperateOpen),
	})
	if !success {
		snag.Panic("仓门开启失败")
	}
	// 绑定
	bs.Bind(bat, sub, rd)
}

func (s *intelligentCabinetService) Bininfo(cab *ent.Cabinet, ordinal int) (*cabdef.BinInfo, error) {
	return adapter.Post[*cabdef.BinInfo](s.GetCabinetAdapterUrlX(cab, "/device/bininfo"), s.GetAdapterUserX(), &cabdef.BinInfoRequest{
		Serial:  cab.Serial,
		Ordinal: silk.Int(ordinal),
	})
}

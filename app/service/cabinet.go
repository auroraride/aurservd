// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "github.com/lithammer/shortuuid/v4"
    log "github.com/sirupsen/logrus"
    "sort"
    "time"
)

type cabinetService struct {
    ctx      context.Context
    orm      *ent.CabinetClient
    modifier *model.Modifier
}

func NewCabinet() *cabinetService {
    return &cabinetService{
        ctx: context.Background(),
        orm: ar.Ent.Cabinet,
    }
}

func NewCabinetWithModifier(m *model.Modifier) *cabinetService {
    s := NewCabinet()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// QueryOne 查询单个电柜
func (s *cabinetService) QueryOne(id uint64) *ent.Cabinet {
    c := s.orm.QueryNotDeleted().Where(cabinet.ID(id)).OnlyX(s.ctx)
    if c == nil {
        snag.Panic("未找到电柜")
    }
    return c
}

// CreateCabinet 创建电柜
func (s *cabinetService) CreateCabinet(req *model.CabinetCreateReq) (res *model.CabinetItem) {
    q := s.orm.Create().
        SetName(req.Name).
        SetSerial(req.Serial).
        SetSn(shortuuid.New()).
        SetStatus(req.Status).
        SetDoors(req.Doors).
        SetNillableRemark(req.Remark).
        SetBrand(req.Brand.Value()).
        SetHealth(model.CabinetHealthStatusOffline)
    if req.BranchID != nil {
        q.SetBranchID(*req.BranchID)
    }

    // 查询设置电池型号
    bms := make([]model.BatteryModel, len(req.Models))
    models := NewBattery().QueryIDs(req.Models)
    for i, bm := range models {
        bms[i] = model.BatteryModel{
            ID:       bm.ID,
            Voltage:  bm.Voltage,
            Capacity: bm.Capacity,
        }
    }
    q.AddBms(models...)
    q.SetModels(bms)

    item := q.SaveX(s.ctx)
    res = new(model.CabinetItem)
    _ = copier.Copy(res, item)
    res.Models = bms
    return
}

// List 查询电柜
func (s *cabinetService) List(req *model.CabinetQueryReq) (res *model.PaginationRes) {
    q := s.orm.QueryNotDeleted().WithBranch(
        func(bq *ent.BranchQuery) {
            bq.WithCity()
        },
    )
    if req.Serial != nil {
        q.Where(cabinet.SerialContainsFold(*req.Serial))
    }
    if req.Name != nil {
        q.Where(cabinet.NameContainsFold(*req.Name))
    }
    if req.CityId != nil {
        q.Where(cabinet.HasBranchWith(branch.CityID(*req.CityId)))
    }
    if req.Brand != nil {
        q.Where(cabinet.Brand(*req.Brand))
    }
    if req.Status != nil {
        q.Where(cabinet.Status(*req.Status))
    }

    return model.ParsePaginationResponse[model.CabinetItem](q.PaginationResult(req.PaginationReq), func() []model.CabinetItem {
        items := q.Pagination(req.PaginationReq).AllX(s.ctx)
        out := make([]model.CabinetItem, len(items))
        for i, item := range items {
            _ = copier.Copy(&out[i], item)
            if item.Edges.Branch != nil {
                city := item.Edges.Branch.Edges.City
                out[i].City = &model.City{
                    ID:   city.ID,
                    Name: city.Name,
                }
            }
        }
        return out
    })
}

// Modify 修改电柜
func (s *cabinetService) Modify(req *model.CabinetModifyReq) {
    c := s.QueryOne(req.ID)
    q := s.orm.UpdateOne(c)
    if req.Models != nil {
        q.ClearBms()
        // 查询设置电池型号
        bms := make([]model.BatteryModel, len(*req.Models))
        models := NewBattery().QueryIDs(*req.Models)
        for i, bm := range models {
            bms[i] = model.BatteryModel{
                ID:       bm.ID,
                Voltage:  bm.Voltage,
                Capacity: bm.Capacity,
            }
        }
        q.AddBms(models...)
        q.SetModels(bms)
    }
    if req.BranchID != nil {
        q.SetBranchID(*req.BranchID)
    }
    if req.Status != nil {
        q.SetStatus(*req.Status)
    }
    if req.Brand != nil {
        q.SetBrand(req.Brand.Value())
    }
    if req.Serial != nil {
        q.SetSerial(*req.Serial)
    }
    if req.Name != nil {
        q.SetName(*req.Name)
    }
    if req.Remark != nil {
        q.SetRemark(*req.Remark)
    }
    q.SaveX(s.ctx)
}

// Delete 删除电柜
func (s *cabinetService) Delete(req *model.CabinetDeleteReq) {
    s.orm.SoftDeleteOneID(req.ID).SaveX(s.ctx)
}

// UpdateStatus 立即更新电柜状态
func (s *cabinetService) UpdateStatus(item *ent.Cabinet) *ent.Cabinet {
    var prov provider.Provider
    if item.Brand == model.CabinetBrandKaixin.Value() {
        prov = provider.NewKaixin()
    } else {
        prov = provider.NewYundong()
    }
    up := s.orm.UpdateOne(item)
    prov.UpdateStatus(up, item)
    return up.SaveX(s.ctx)
}

// DoorOpenStatus 获取柜门状态
func (s *cabinetService) DoorOpenStatus(item *ent.Cabinet, index int) model.CabinetBinDoorStatus {
    item = s.UpdateStatus(item)
    if len(item.Bin) < index {
        return model.CabinetBinDoorStatusUnknown
    }
    bin := item.Bin[index]
    if !bin.DoorHealth {
        return model.CabinetBinDoorStatusFail
    }
    if bin.OpenStatus {
        return model.CabinetBinDoorStatusOpen
    }
    return model.CabinetBinDoorStatusClose
}

// Detail 电柜详细信息
func (s *cabinetService) Detail(id uint64) *model.CabinetDetailRes {
    item := s.orm.QueryNotDeleted().
        Where(cabinet.ID(id)).
        OnlyX(s.ctx)
    if item == nil {
        snag.Panic("未找到电柜")
    }
    item = s.UpdateStatus(item)
    res := new(model.CabinetDetailRes)
    _ = copier.Copy(res, item)
    return res
}

// DoorOperate 操作柜门
func (s *cabinetService) DoorOperate(req *model.CabinetDoorOperateReq, operator model.CabinetDoorOperator) (state bool, err error) {
    opId := shortuuid.New()
    now := time.Now()
    // 查找柜子和仓位
    item := s.QueryOne(*req.ID)
    if len(item.Bin) < *req.Index {
        err = errors.New("该柜门未找到")
        return
    }

    brand := model.CabinetBrand(item.Brand)
    op, ok := req.Operation.Value(brand)
    if !ok {
        err = errors.New("操作方式错误")
        return
    }
    if *req.Operation == model.CabinetDoorOperateLock {
        if req.Remark == "" {
            err = errors.New("该操作必须携带操作原因")
            return
        } else {
            req.Remark = ""
        }
    }
    var prov provider.Provider
    up := ar.Ent.Cabinet.UpdateOne(item).SetHealth(model.CabinetHealthStatusOnline)
    switch brand {
    case model.CabinetBrandYundong:
        prov = provider.NewYundong()
        break
    case model.CabinetBrandKaixin:
        prov = provider.NewKaixin()
        break
    }
    prov.PrepareRequest()
    state = prov.DoorOperate(operator.Name+"-"+opId, item.Serial, op, *req.Index)
    // 如果成功, 重新获取状态更新数据
    if state {
        // 更新仓位备注
        bins := item.Bin
        bins[*req.Index].Remark = req.Remark
        prov.UpdateStatus(up, item)
        up.SetBin(bins).SaveX(s.ctx)
    }
    go func() {
        // 上传日志
        slsCfg := ar.Config.Aliyun.Sls
        lg := &sls.LogGroup{
            Logs: []*sls.Log{{
                Time: tea.Uint32(uint32(now.Unix())),
                Contents: logging.GenerateLogContent(&logging.DoorOperateLog{
                    ID:            opId,
                    Brand:         brand.String(),
                    OperatorName:  operator.Name,
                    OperatorID:    operator.ID,
                    OperatorPhone: operator.Phone,
                    Serial:        item.Serial,
                    Name:          item.Bin[*req.Index].Name,
                    Operation:     req.Operation.String(),
                    Success:       state,
                    Remark:        req.Remark,
                    Time:          now.Format(carbon.DateTimeLayout),
                }),
            }},
        }
        err = ali.NewSls().PutLogs(slsCfg.Project, slsCfg.DoorLog, lg)
        if err != nil {
            log.Error(err)
        }
    }()
    return
}

// Reboot 重启电柜
func (s *cabinetService) Reboot(req *model.IDPostReq) bool {
    if s.modifier == nil {
        snag.Panic("请求不正确")
    }
    now := time.Now()
    opId := shortuuid.New()

    item := s.QueryOne(req.ID)
    if item.Brand == model.CabinetBrandKaixin.Value() {
        snag.Panic("凯信电柜不支持该操作")
    }
    var prov provider.Provider
    var state bool
    prov = provider.NewYundong()
    state = prov.Reboot(s.modifier.Name+"-"+opId, item.Serial)

    // 如果成功, 重新获取状态更新数据
    up := ar.Ent.Cabinet.UpdateOne(item).SetHealth(model.CabinetHealthStatusOnline)
    if state {
        // 更新仓位备注
        prov.UpdateStatus(up, item)
        up.SaveX(s.ctx)
    }

    brand := model.CabinetBrand(item.Brand)
    go func() {
        // 上传日志
        slsCfg := ar.Config.Aliyun.Sls
        lg := &sls.LogGroup{
            Logs: []*sls.Log{{
                Time: tea.Uint32(uint32(now.Unix())),
                Contents: logging.GenerateLogContent(&logging.DoorOperateLog{
                    ID:            opId,
                    Brand:         brand.String(),
                    OperatorName:  s.modifier.Name,
                    OperatorID:    s.modifier.ID,
                    OperatorPhone: s.modifier.Phone,
                    OperatorRole:  model.CabinetDoorOperatorRoleManager,
                    Serial:        item.Serial,
                    Operation:     "重启",
                    Success:       state,
                    Time:          now.Format(carbon.DateTimeLayout),
                }),
            }},
        }
        err := ali.NewSls().PutLogs(slsCfg.Project, slsCfg.DoorLog, lg)
        if err != nil {
            log.Error(err)
            return
        }
    }()

    return state
}

// QueryWithSerial 根据序列号查找电柜
func (s *cabinetService) QueryWithSerial(serial string) *ent.Cabinet {
    cab, _ := s.orm.QueryNotDeleted().Where(cabinet.Serial(serial)).Only(s.ctx)
    if cab == nil {
        snag.Panic("未找到电柜")
    }
    return cab
}

// Usable 获取换电可用仓位信息
func (s *cabinetService) Usable(cab *ent.Cabinet) (op model.RiderCabinetOperateProcess) {
    sort.Slice(cab.Bin, func(i, j int) bool {
        return cab.Bin[i].Electricity.Value() > cab.Bin[j].Electricity.Value()
    })
    // 查看电柜是否有满电
    for index, bin := range cab.Bin {
        if !bin.Battery && op.EmptyBin == nil {
            // 获取空仓
            op.EmptyBin = &model.CabinetBinBasicInfo{
                Index:       index,
                Electricity: bin.Electricity,
            }
        }
        if bin.Electricity.IsBatteryFull() && op.FullBin == nil {
            op.FullBin = &model.CabinetBinBasicInfo{
                Index:       index,
                Electricity: bin.Electricity,
            }
        }
    }
    if op.FullBin == nil {
        op.Alternative = &model.CabinetBinBasicInfo{
            Index:       len(cab.Bin) - 1,
            Electricity: cab.Bin[len(cab.Bin)-1].Electricity,
        }
    }
    return
}

// Health 判定电柜是否可用
// TODO 上次获取状态多久后标记为offline
func (s *cabinetService) Health(cab *ent.Cabinet) bool {
    return cab.Status == model.CabinetStatusNormal &&
        cab.Health == model.CabinetHealthStatusOnline &&
        len(cab.Bin) > 0
}

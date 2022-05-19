// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
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
    "time"
)

type cabinetService struct {
    ctx context.Context
    orm *ent.CabinetClient
}

func NewCabinet() *cabinetService {
    return &cabinetService{
        ctx: context.Background(),
        orm: ar.Ent.Cabinet,
    }
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
func (s *cabinetService) CreateCabinet(modifier *model.Modifier, req *model.CabinetCreateReq) (res *model.CabinetItem) {
    q := s.orm.Create().
        SetName(req.Name).
        SetSerial(req.Serial).
        SetSn(shortuuid.New()).
        SetStatus(uint(req.Status)).
        SetDoors(req.Doors).
        SetLastModifier(modifier).
        SetCreator(modifier).
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

// Query 查询电柜
func (s *cabinetService) Query(req *model.CabinetQueryReq) (res *model.PaginationRes) {
    res = new(model.PaginationRes)
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
        q.Where(cabinet.Status(uint(*req.Status)))
    }

    res.Pagination = q.PaginationResult(req.PaginationReq)

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
    res.Items = out
    return
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
        q.SetStatus(uint(*req.Status))
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
func (s *cabinetService) Delete(modifier *model.Modifier, req *model.CabinetDeleteReq) {
    s.orm.SoftDeleteOneID(req.ID).SetLastModifier(modifier).SaveX(s.ctx)
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
func (s *cabinetService) DoorOperate(modifier *model.Modifier, req *model.CabinetDoorOperateReq) (state bool) {
    opId := shortuuid.New()
    now := time.Now()
    // 查找柜子和仓位
    item := s.QueryOne(*req.ID)
    if len(item.Bin) < *req.Index {
        snag.Panic("该柜门未找到")
    }

    brand := model.CabinetBrand(item.Brand)
    op, ok := req.Operation.Value(brand)
    if !ok {
        snag.Panic("操作方式错误")
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
    state = prov.DoorOperate(modifier.Name+"-"+opId, item.Serial, op, *req.Index)
    // 如果成功, 重新获取状态更新数据
    if state {
        // 更新仓位备注
        bins := item.Bin
        bins[*req.Index].Remark = *req.Remark
        prov.UpdateStatus(up, item)
        up.SetBin(bins).SaveX(s.ctx)
    }
    go func() {
        // 上传日志
        slsCfg := ar.Config.Aliyun.Sls
        lg := &sls.LogGroup{
            Logs: []*sls.Log{{
                Time: tea.Uint32(uint32(now.Unix())),
                Contents: provider.ParseLogContent(&provider.OperationLog{
                    ID:        opId,
                    Brand:     brand.String(),
                    User:      modifier.Name,
                    UserID:    modifier.ID,
                    Phone:     modifier.Phone,
                    Serial:    item.Serial,
                    Name:      item.Bin[*req.Index].Name,
                    Operation: req.Operation.String(),
                    Success:   state,
                    Remark:    *req.Remark,
                    Time:      now.Format(carbon.DateTimeLayout),
                }),
            }},
        }
        err := ali.NewSls().PutLogs(slsCfg.Project, slsCfg.Operation, lg)
        if err != nil {
            log.Error(err)
            return
        }
    }()
    return
}

// Reboot 重启电柜
func (s *cabinetService) Reboot(modifier *model.Modifier, req *model.IDPostReq) bool {
    now := time.Now()
    opId := shortuuid.New()

    item := s.QueryOne(req.ID)
    if item.Brand == model.CabinetBrandKaixin.Value() {
        snag.Panic("凯信电柜不支持该操作")
    }
    var prov provider.Provider
    var state bool
    prov = provider.NewYundong()
    state = prov.Reboot(modifier.Name+"-"+opId, item.Serial)

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
                Contents: provider.ParseLogContent(&provider.OperationLog{
                    ID:        opId,
                    Brand:     brand.String(),
                    User:      modifier.Name,
                    UserID:    modifier.ID,
                    Phone:     modifier.Phone,
                    Serial:    item.Serial,
                    Operation: "重启",
                    Success:   state,
                    Time:      now.Format(carbon.DateTimeLayout),
                }),
            }},
        }
        err := ali.NewSls().PutLogs(slsCfg.Project, slsCfg.Operation, lg)
        if err != nil {
            log.Error(err)
            return
        }
    }()

    return state
}

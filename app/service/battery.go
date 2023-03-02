// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-24
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/adapter"
    "github.com/auroraride/adapter/log"
    "github.com/auroraride/adapter/rpc/pb"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/rpc"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/battery"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/jackc/pgconn"
    "github.com/labstack/echo/v4"
    "go.uber.org/zap"
    "strconv"
    "strings"
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

    return s.orm.Create().SetModel(ab.Model).SetSn(sn).SetNillableCabinetID(cabID).SetNillableOrdinal(ordinal).Save(s.ctx)
}

// SyncPutout 同步消息 - 从电柜中取出
func (s *batteryService) SyncPutout(cabinetID uint64, ordinal int) {
    _ = s.orm.Update().Where(battery.CabinetID(cabinetID), battery.Ordinal(ordinal)).ClearCabinetID().ClearOrdinal().Exec(s.ctx)
}

// SyncPutin 同步消息 - 放入电柜中
func (s *batteryService) SyncPutin(sn, serial string, cabinetID uint64, ordinal int, old model.CabinetBins) (bat *ent.Battery, err error) {
    bat, err = s.LoadOrCreate(sn, &model.BatteryInCabinet{
        CabinetID: cabinetID,
        Ordinal:   ordinal,
    })
    if err != nil {
        return
    }

    // 移除别的电池信息
    s.SyncPutout(cabinetID, ordinal)

    // 更新电池电柜信息
    bat, err = bat.Update().SetCabinetID(cabinetID).SetOrdinal(ordinal).ClearRiderID().ClearSubscribeID().Save(s.ctx)

    // TODO 更新电池流转
    // go NewBatteryFlow().Create(model.BatteryFlowCreateReq{
    //     SN:        bat.Sn,
    //     BatteryID: bat.ID,
    //     CabinetID: silk.Pointer(cabinetID),
    //     Ordinal:   silk.Pointer(ordinal),
    //     Serial:    silk.Pointer(serial),
    // })
    return
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

        _, err = creator.SetModel(ab.Model).SetSn(sn).Save(s.ctx)
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
    q = s.orm.Query().WithRider().WithCity().WithCabinet().Order(ent.Desc(battery.FieldCreatedAt))
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

    return
}

func (s *batteryService) List(req *model.BatteryListReq) (res *model.PaginationRes) {
    q, _ := s.listFilter(req.BatteryFilter)
    var sn []string
    res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Battery) (res *model.BatteryListRes) {
        sn = append(sn, item.Sn)
        res = &model.BatteryListRes{
            ID:     item.ID,
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
        return
    })

    // 请求xcbms rpc
    r, _ := rpc.XcBmsBatch(s.ctx, &pb.BatteryBatchRequest{Sn: sn})

    if r == nil {
        return
    }

    for _, data := range res.Items.([]*model.BatteryListRes) {
        if rb, ok := r.Items[data.SN]; ok {
            if len(rb.Heartbeats) > 0 {
                data.XcBmsBattery = model.NewXcBmsBattery(rb.Heartbeats[0])
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
        SetRefManually(rider.Table, r.ID).
        SetDiff(before, after).
        Send()
}

// RiderDetail 获取电池详情
func (s *batteryService) RiderDetail(riderID uint64) (res model.BatteryDetail) {
    bat, _ := s.QueryRiderID(riderID)
    if bat != nil {
        res = model.BatteryDetail{
            ID:    bat.ID,
            Model: bat.Model,
            SN:    bat.Sn,
            Soc:   0,
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
    err := s.Allocate(bat.Update(), bat, sub, false)
    if err != nil {
        snag.Panic(err)
    }

    go logging.NewOperateLog().
        SetOperate(model.OperateBindBattery).
        SetRef(rd).
        SetDiff("", "新电池: "+bat.Sn).
        SetModifier(s.modifier).
        Send()
}

func (s *batteryService) Unbind(bat *ent.Battery, rd *ent.Rider) {
    err := s.Unallocate(bat)
    if err != nil {
        snag.Panic(err)
    }

    go logging.NewOperateLog().
        SetOperate(model.OperateUnbindBattery).
        SetRef(rd).
        SetDiff("旧电池: "+bat.Sn, "无电池").
        SetModifier(s.modifier).
        Send()
}

func (s *batteryService) UnbindRequest(req *model.BatteryUnbindRequest) {
    // 查找订阅
    sub := NewSubscribe().QueryEffectiveIntelligentX(req.RiderID, ent.SubscribeQueryWithBattery, ent.SubscribeQueryWithRider, ent.SubscribeQueryWithBattery)

    bat := sub.Edges.Battery
    if bat == nil {
        snag.Panic("未找到绑定的电池")
    }

    s.Unbind(bat, sub.Edges.Rider)
}

// Allocate 将电池分配给骑手
func (s *batteryService) Allocate(buo *ent.BatteryUpdateOne, bat *ent.Battery, sub *ent.Subscribe, ignoreError bool) (err error) {
    err = buo.SetSubscribeID(sub.ID).SetRiderID(sub.RiderID).Exec(s.ctx)
    if err != nil && ent.IsConstraintError(err) {
        switch v := err.(*ent.ConstraintError).Unwrap().(type) {
        case *pgconn.PgError:
            if v.ConstraintName == "battery_subscribe_id_key" || v.ConstraintName == "battery_rider_id_key" {
                // 删除原有信息
                err = s.orm.Update().Where(battery.SubscribeID(sub.ID)).ClearRiderID().ClearSubscribeID().Exec(s.ctx)
                if err != nil {
                    return
                }
                err = buo.SetSubscribeID(sub.ID).SetRiderID(sub.RiderID).Exec(s.ctx)
            }
        }
    }

    // 更新流转
    if ignoreError || err == nil {
        go NewBatteryFlow().Create(model.BatteryFlowCreateReq{
            SN:          bat.Sn,
            BatteryID:   bat.ID,
            RiderID:     silk.Pointer(sub.RiderID),
            SubscribeID: silk.Pointer(sub.ID),
        })
    }
    return
}

// Unallocate 清除骑手信息
func (s *batteryService) Unallocate(bat *ent.Battery) (err error) {
    return bat.Update().ClearSubscribeID().ClearRiderID().Exec(s.ctx)
}

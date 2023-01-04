// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-24
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/adapter"
    "github.com/auroraride/adapter/defs/cabdef"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/battery"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
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

func (s *batteryService) LoadOrStore(sn string) (bat *ent.Battery, err error) {
    bat, _ = s.QuerySn(sn)
    if bat != nil {
        return
    }
    ab := adapter.ParseBatterySN(sn)
    return s.orm.Create().SetModel(ab.Model).SetSn(sn).Save(s.ctx)
}

// PutinCabinet 电池放入电柜
func (s *batteryService) PutinCabinet(sn string, cab *ent.Cabinet) (bat *ent.Battery, err error) {
    ab := adapter.ParseBatterySN(sn)
    if ab.Model == "" {
        log.Errorf("电池更新失败, 未找到型号: %#v", *ab)
    }

    var (
        cityID    *uint64
        cabinetID *uint64
    )

    if cab == nil {
        log.Errorf("电池更新失败, 未找到电柜: %#v", *ab)
    } else {
        cityID = cab.CityID
        cabinetID = silk.UInt64(cab.ID)
    }

    bat, _ = s.orm.Query().Where(battery.Sn(sn)).First(s.ctx)
    if bat != nil {
        return bat.Update().SetNillableCabinetID(cabinetID).SetModel(ab.Model).SetNillableCityID(cityID).ClearRiderID().ClearSubscribeID().Save(s.ctx)
    }
    return s.orm.Create().SetSn(sn).SetModel(ab.Model).SetNillableCityID(cityID).SetNillableCabinetID(cabinetID).Save(s.ctx)
}

// TODO 电池需要做库存管理

// Create 创建电池
func (s *batteryService) Create(req *model.BatteryCreateReq) {
    enable := true
    if req.Enable != nil {
        enable = *req.Enable
    }

    // 解析电池编号
    ab := adapter.ParseBatterySN(req.SN)
    if ab.Model == "" {
        snag.Panic("电池编号解析失败, 请擦亮你的双眼")
    }
    _, err := s.orm.Create().
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
        ab := adapter.ParseBatterySN(sn)
        if ab.Model == "" {
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

        _, err := creator.SetModel(ab.Model).SetSn(sn).Save(s.ctx)
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
    q = s.orm.Query().WithRider().WithCity().WithCabinet()
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

func (s *batteryService) List(req *model.BatteryListReq) *model.PaginationRes {
    q, _ := s.listFilter(req.BatteryFilter)
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Battery) (res model.BatteryListRes) {
        c := item.Edges.City
        res = model.BatteryListRes{
            ID: item.ID,
            City: model.City{
                ID:   c.ID,
                Name: c.Name,
            },
            Model:  item.Model,
            Enable: item.Enable,
            SN:     item.Sn,
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
                Brand:  model.CabinetBrand(cab.Brand),
                Serial: cab.Serial,
                Name:   cab.Name,
            }
        }
        return
    })
}

func (s *batteryService) Sync(data *cabdef.BatteryMessage) {
    if data.Cabinet == "" {
        log.Errorf("[SYNC] 电池:%s 缺少参数 cabinetSerial", data.SN)
        return
    }

    if data.SN == "" {
        log.Error("[SYNC] 电池缺少参数 sn")
        return
    }

    // 查找电柜
    cab := NewCabinet().QueryOneSerial(data.Cabinet)
    if cab == nil {
        log.Errorf("[SYNC] 电池:%s 未找到电柜", data.SN)
        return
    }

    // 更新或创建电池
    _, _ = s.PutinCabinet(data.SN, cab)
}

// RiderPutout 骑手取走电池
func (s *batteryService) RiderPutout(sn string, sub *ent.Subscribe) (bat *ent.Battery) {
    bat, _ = s.LoadOrStore(sn)

    if bat == nil {
        log.Error("电池订阅更新失败, 未找到电池信息")
        return
    }

    // 更新电池
    _ = bat.Update().ClearCabinetID().SetSubscribeID(sub.ID).SetRiderID(sub.RiderID).Exec(s.ctx)

    return
}

// RiderPutin 骑手放入电池
func (s *batteryService) RiderPutin(sn string, sub *ent.Subscribe, cab *ent.Cabinet) (bat *ent.Battery) {
    // TODO 是否记录骑手信息?
    bat, _ = s.PutinCabinet(sn, cab)

    return
}

// Bind 绑定骑手
func (s *batteryService) Bind(req *model.BatteryBind) {
    // 查找骑手
    r := NewRider().Query(req.RiderID)

    // 查找订阅
    sub := NewSubscribe().QueryEffectiveX(r.ID, ent.SubscribeQueryWithBattery, ent.SubscribeQueryWithRider)
    if sub.BatteryID != nil {
        snag.Panic("当前骑手有绑定中的电池, 无法重复绑定")
    }

    // 查找电池
    bat := NewBattery().QueryIDX(req.BatteryID)
    if bat.RiderID != nil || bat.SubscribeID != nil {
        snag.Panic("当前电池有绑定中的骑手, 无法重复绑定")
    }

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        // 绑定骑手
        err = sub.Update().SetBatteryID(bat.ID).SetBatterySn(bat.Sn).Exec(s.ctx)
        if err != nil {
            return
        }

        // 更新电池
        return bat.Update().SetSubscribeID(sub.ID).SetRiderID(sub.RiderID).Exec(s.ctx)
    })

    // diff
    var before, after string

    go logging.NewOperateLog().
        SetOperate(model.OperateBindBattery).
        SetRef(sub.Edges.Rider).
        SetDiff(before, after).
        SetModifier(s.modifier).
        Send()
}

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

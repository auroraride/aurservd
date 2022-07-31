// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "entgo.io/ent/dialect/sql"
    "entgo.io/ent/dialect/sql/sqljson"
    "errors"
    "fmt"
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "github.com/lithammer/shortuuid/v4"
    log "github.com/sirupsen/logrus"
    "golang.org/x/exp/slices"
    "regexp"
    "sort"
    "strconv"
    "strings"
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
        orm: ent.Database.Cabinet,
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
    c, _ := s.orm.QueryNotDeleted().Where(cabinet.ID(id)).First(s.ctx)
    if c == nil {
        snag.Panic("未找到电柜")
    }
    return c
}

// CreateCabinet 创建电柜
func (s *cabinetService) CreateCabinet(req *model.CabinetCreateReq) (res *model.CabinetItem) {
    if req.Status == model.CabinetStatusNormal && req.BranchID == nil {
        snag.Panic("电柜投产必须选择网点")
    }

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
        b := NewBranch().Query(*req.BranchID)
        q.SetBranchID(*req.BranchID).SetCityID(b.CityID)
    }
    if req.SimSn != "" && req.SimDate != "" {
        q.SetSimSn(req.SimSn).
            SetSimDate(tools.NewTime().ParseDateStringX(req.SimDate))
    }

    res = new(model.CabinetItem)

    // 查询设置电池型号
    models := NewBattery().QueryModelsX(req.Models)
    for _, bm := range models {
        res.Models = append(res.Models, bm.Model)
    }
    q.AddBms(models...)

    item := q.SaveX(s.ctx)
    res = new(model.CabinetItem)
    _ = copier.Copy(res, item)

    if item.Status == model.CabinetStatusNormal {
        go s.Deploy(item)
    }

    return
}

// List 查询电柜
func (s *cabinetService) List(req *model.CabinetQueryReq) (res *model.PaginationRes) {
    q := s.orm.QueryNotDeleted().WithCity().WithBms()
    if req.Serial != nil {
        q.Where(cabinet.SerialContainsFold(*req.Serial))
    }
    if req.Name != nil {
        q.Where(cabinet.NameContainsFold(*req.Name))
    }
    if req.CityId != nil {
        q.Where(cabinet.CityID(*req.CityId))
    }
    if req.Brand != nil {
        q.Where(cabinet.Brand(*req.Brand))
    }
    if req.Status != nil {
        q.Where(cabinet.Status(*req.Status))
    }
    if req.Model != nil {
        q.Where(cabinet.HasBmsWith(batterymodel.Model(*req.Model)))
    }
    if req.Online != 0 {
        switch req.Online {
        case 1:
            q.Where(cabinet.Health(model.CabinetHealthStatusOnline))
            break
        case 2:
            q.Where(cabinet.Health(model.CabinetHealthStatusOffline))
            break
        }
    }

    return model.ParsePaginationResponse[model.CabinetItem, ent.Cabinet](q, req.PaginationReq, func(item *ent.Cabinet) (res model.CabinetItem) {
        _ = copier.Copy(&res, item)

        if !item.SimDate.IsZero() {
            res.SimDate = item.SimDate.Format(carbon.DateLayout)
        }

        res.CreatedAt = item.CreatedAt.Format(carbon.DateTimeLayout)

        city := item.Edges.City
        if city != nil {
            res.City = &model.City{
                ID:   city.ID,
                Name: city.Name,
            }
        }
        bms := item.Edges.Bms
        for _, bm := range bms {
            res.Models = append(res.Models, bm.Model)
        }
        return res
    })
}

// Modify 修改电柜
func (s *cabinetService) Modify(req *model.CabinetModifyReq) {
    cab, _ := s.orm.QueryNotDeleted().Where(cabinet.ID(req.ID)).WithBms().First(s.ctx)
    if cab == nil {
        snag.Panic("未找到电柜")
    }
    willDeploy := cab.Status == model.CabinetStatusPending && req.Status != nil && *req.Status == model.CabinetStatusNormal
    err := ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
        q := tx.Cabinet.UpdateOne(cab)
        if req.Models != nil {
            var models []string
            for _, bm := range cab.Edges.Bms {
                models = append(models, bm.Model)
            }
            // 排序
            sort.Slice(models, func(i, j int) bool {
                return strings.Compare(models[i], models[j]) < 0
            })
            rms := *req.Models
            sort.Slice(rms, func(i, j int) bool {
                return strings.Compare(rms[i], rms[j]) < 0
            })

            if slices.Compare(rms, models) != 0 {
                q.ClearBms()
                // 查询设置电池型号
                q.AddBms(NewBattery().QueryModelsX(*req.Models)...)
            }
        }
        if req.BranchID != nil {
            b := NewBranch().Query(*req.BranchID)
            q.SetLng(b.Lng).
                SetLat(b.Lat).
                SetAddress(b.Address).
                SetBranchID(*req.BranchID).
                SetCityID(b.CityID)
        } else if cab.BranchID == nil {
            // 检查网点
            if cab.Status == model.CabinetStatusNormal || (req.Status != nil && *req.Status == model.CabinetStatusNormal) {
                return errors.New("电柜投产必须选择网点")
            }
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

        if req.SimSn != nil {
            q.SetSimSn(*req.SimSn)
        }

        if req.SimDate != nil {
            end := tools.NewTime().ParseDateStringX(*req.SimDate)
            if time.Now().After(end) {
                snag.Panic("sim卡到期日期不能早于现在")
            }
            q.SetSimDate(end)
        }

        cab, err = q.Save(s.ctx)
        if err != nil {
            return
        }

        return
    })

    if err != nil {
        snag.Panic(err)
    }

    if willDeploy {
        go s.Deploy(cab)
    }
}

func (s *cabinetService) Deploy(c *ent.Cabinet) {
    if model.CabinetBrand(c.Brand) != model.CabinetBrandYundong {
        return
    }

    // 查找电柜
    b, _ := ent.Database.Branch.QueryNotDeleted().Where(branch.ID(*c.BranchID)).WithCity().First(s.ctx)
    if b == nil || b.Edges.City == nil {
        snag.Panic("投产失败, 未找到网点信息, 请将电柜改为未投放并调整好网点重试")
    }

    bec := b.Edges.City

    // 云动部署投产
    provider.NewYundong().UpdateBasicInfo(model.YundongDeployInfo{
        SN:       c.Serial,
        AreaCode: strconv.Itoa(int(bec.ID)),
        Address:  b.Address,
        Lat:      b.Lat,
        Lng:      b.Lng,
        Name:     c.Name,
        Phone:    "4000290929",
        Contact:  "极光出行客服",
        City:     c.Name,
    })
}

// Delete 删除电柜
func (s *cabinetService) Delete(req *model.CabinetDeleteReq) {
    s.orm.SoftDeleteOneID(req.ID).SaveX(s.ctx)
}

// UpdateStatus 立即更新电柜状态
func (s *cabinetService) UpdateStatus(item *ent.Cabinet, params ...any) {
    var prov provider.Provider
    if item.Brand == model.CabinetBrandKaixin.Value() {
        prov = provider.NewKaixin()
    } else {
        prov = provider.NewYundong()
    }
    err := prov.UpdateStatus(item, params...)
    // 如果返回失败, 则延迟2秒后重新请求一次
    if err != nil {
        time.Sleep(2 * time.Second)
    }
    _ = prov.UpdateStatus(item, params...)
}

// DoorOpenStatus 获取柜门状态
func (s *cabinetService) DoorOpenStatus(item *ent.Cabinet, index int, params ...any) model.CabinetBinDoorStatus {
    s.UpdateStatus(item, params...)
    if len(item.Bin) == 0 || len(item.Bin) < index {
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
        WithBms().
        OnlyX(s.ctx)
    if item == nil {
        snag.Panic("未找到电柜")
    }
    bms := item.Edges.Bms

    s.UpdateStatus(item)
    res := new(model.CabinetDetailRes)
    _ = copier.Copy(res, item)
    for _, bm := range bms {
        res.Models = append(res.Models, bm.Model)
    }

    return res
}

// DoorOperate 操作柜门
func (s *cabinetService) DoorOperate(req *model.CabinetDoorOperateReq, operator model.CabinetDoorOperator, params ...any) (state bool, err error) {
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
        }
    }
    var prov provider.Provider
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
        log.Infof("%s操作成功[%s %s]", item.Serial, req.Operation.String(), req.Remark)
        // 如果是锁仓, 需要更新仓位备注
        if *req.Operation == model.CabinetDoorOperateLock {
            item.Bin[*req.Index].Remark = req.Remark
        }
        // 如果是解锁, 需要清除仓位备注
        if *req.Operation == model.CabinetDoorOperateUnlock {
            item.Bin[*req.Index].Remark = ""
        }
        _ = prov.UpdateStatus(item, params...)
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
                    OperatorRole:  operator.Role,
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
    if state {
        // 更新仓位备注
        _ = prov.UpdateStatus(item)
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
    cab, _ := s.orm.QueryNotDeleted().Where(cabinet.Serial(serial)).WithBms().Only(s.ctx)
    if cab == nil {
        snag.Panic("未找到电柜")
    }
    return cab
}

// ModelInclude 电柜是否可用指定型号电池
func (s *cabinetService) ModelInclude(item *ent.Cabinet, model string) bool {
    bms := item.Edges.Bms
    if bms == nil {
        return false
    }
    for _, bm := range bms {
        if bm.Model == model {
            return true
        }
    }
    return false
}

// Usable 获取换电可用仓位信息
func (s *cabinetService) Usable(cab *ent.Cabinet) (op model.RiderCabinetOperateProcess) {
    sort.Slice(cab.Bin, func(i, j int) bool {
        return cab.Bin[i].Electricity.Value()+cab.Bin[i].Voltage*0.1 > cab.Bin[j].Electricity.Value()+cab.Bin[j].Voltage*0.1
    })
    // 查看电柜是否有满电
    var max *model.CabinetBinBasicInfo
    var empty *model.CabinetBinBasicInfo
    for _, bin := range cab.Bin {
        // 若仓门不正常直接跳过
        if !bin.DoorHealth {
            continue
        }

        if !bin.Battery && empty == nil {
            // 获取空仓
            empty = &model.CabinetBinBasicInfo{
                Index:       bin.Index,
                Electricity: bin.Electricity,
            }
        }

        if bin.Battery && max == nil {
            max = &model.CabinetBinBasicInfo{
                Index:       bin.Index,
                Electricity: bin.Electricity,
                Voltage:     bin.Voltage,
            }
        }
        if max != nil && empty != nil {
            break
        }
    }

    if max == nil || empty == nil {
        snag.Panic("电柜异常, 无法换电")
    }

    op.EmptyBin = empty

    if max.Electricity.IsBatteryFull() {
        op.FullBin = max
    } else {
        op.Alternative = max
    }

    return
}

// Health 判定电柜是否可用
func (s *cabinetService) Health(cab *ent.Cabinet) bool {
    return cab.Status == model.CabinetStatusNormal &&
        cab.Health == model.CabinetHealthStatusOnline &&
        time.Now().Sub(cab.UpdatedAt).Minutes() < 5 &&
        len(cab.Bin) > 0
}

func (s *cabinetService) Data(req *model.CabinetDataReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().WithBms()
    switch req.Status {
    case 1:
        q.Where(cabinet.Health(model.CabinetHealthStatusOnline))
        break
    case 2:
        q.Where(cabinet.Health(model.CabinetHealthStatusOffline))
        break
    case 3:
        q.Modify(func(sel *sql.Selector) {
            sel.Where(sqljson.ValueContains(sel.C(cabinet.FieldBin), []ar.Map{{"doorHealth": false}}))
        })
        break
    }

    if req.Name != "" {
        q.Where(cabinet.NameContainsFold(req.Name))
    }

    if req.Serial != "" {
        q.Where(cabinet.SerialContainsFold(req.Serial))
    }

    if req.CityID != 0 {
        q.Where(cabinet.CityID(req.CityID))
    }

    if req.Brand != "" {
        q.Where(cabinet.Brand(req.Brand.Value()))
    }

    if req.Votage != 0 {
        bm := fmt.Sprintf("%.0fV", req.Votage)
        q.Where(cabinet.HasBmsWith(batterymodel.ModelHasPrefix(bm)))
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Cabinet) model.CabinetDataRes {
        res := model.CabinetDataRes{
            Name:       item.Name,
            Serial:     item.Serial,
            Model:      "",
            Brand:      model.CabinetBrand(item.Brand).String(),
            Online:     item.Health == model.CabinetHealthStatusOnline,
            BinNum:     int(item.Doors),
            BatteryNum: int(item.BatteryNum),
            EmptyNum:   0,
            LockNum:    0,
            FullNum:    0,
        }

        bms := item.Edges.Bms
        if len(bms) > 0 {
            res.Model = regexp.MustCompile(`(?m)(\d+)V\d+AH`).ReplaceAllString(bms[0].Model, "${1}V")
        }

        res.Bins = make([]model.CabinetDataBin, len(item.Bin))
        for i, bin := range item.Bin {

            if bin.Battery {
                if bin.Full {
                    res.Bins[i].Status = model.CabinetDataBinStatusFull
                    res.FullNum += 1
                } else {
                    res.Bins[i].Status = model.CabinetDataBinStatusCharging
                }
            } else {
                res.Bins[i].Status = model.CabinetDataBinStatusEmpty
                res.EmptyNum += 1
            }

            if !bin.DoorHealth {
                res.Bins[i].Status = model.CabinetDataBinStatusLock
                res.Bins[i].Remark = bin.Remark
                res.LockNum += 1
            }
        }

        return res
    })
}

// Busy TODO 是否需要两次换电间隔
func (s *cabinetService) Busy(cab *ent.Cabinet) bool {
    if model.CabinetBusying(cab.Serial) {
        return true
    }
    last, _ := ent.Database.Exchange.QueryNotDeleted().Where(exchange.CabinetID(cab.ID)).Order(ent.Desc(exchange.FieldCreatedAt)).First(s.ctx)
    if last != nil {

    }
    return false
}

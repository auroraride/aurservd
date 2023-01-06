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
    "github.com/auroraride/adapter/defs/cabdef"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/stock"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/silk"
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

func (s *cabinetService) QueryOneSerial(serial string) *ent.Cabinet {
    serial = strings.ReplaceAll(serial, "https://www.yunfuture.cn/qrcode/cabinet?cabinetSN=", "")
    cab, _ := s.orm.QueryNotDeleted().Where(cabinet.Serial(serial)).WithModels().First(s.ctx)
    return cab
}

func (s *cabinetService) QueryOneSerialX(serial string) *ent.Cabinet {
    cab := s.QueryOneSerial(serial)
    if cab == nil {
        snag.Panic("未找到电柜")
    }
    return cab
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
        SetStatus(req.Status.Value()).
        SetDoors(req.Doors).
        SetNillableRemark(req.Remark).
        SetIntelligent(req.Intelligent).
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
    models := NewBatteryModel().QueryModelsX(req.Models)
    for _, bm := range models {
        res.Models = append(res.Models, bm.Model)
    }
    q.AddModels(models...)

    item := q.SaveX(s.ctx)
    res = new(model.CabinetItem)
    _ = copier.Copy(res, item)

    if model.CabinetStatus(item.Status) == model.CabinetStatusNormal {
        go s.Deploy(item)
    }

    return
}

// List 查询电柜
func (s *cabinetService) List(req *model.CabinetQueryReq) (res *model.PaginationRes) {
    q := s.orm.QueryNotDeleted().WithCity().WithModels()

    if s.modifier != nil && s.modifier.Phone == "15537112255" {
        req.CityID = silk.UInt64(410100)
    }

    if req.Serial != nil {
        q.Where(cabinet.SerialContainsFold(*req.Serial))
    }
    if req.Name != nil {
        q.Where(cabinet.NameContainsFold(*req.Name))
    }
    if req.CityID != nil {
        q.Where(cabinet.CityID(*req.CityID))
    }
    if req.Brand != nil {
        q.Where(cabinet.Brand(*req.Brand))
    }
    if req.Status != nil {
        q.Where(cabinet.Status(*req.Status))
    }
    if req.Model != nil {
        q.Where(cabinet.HasModelsWith(batterymodel.Model(*req.Model)))
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
    if req.Intelligent != 0 {
        q.Where(cabinet.Intelligent(req.Intelligent == 1))
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
        bms := item.Edges.Models
        for _, bm := range bms {
            res.Models = append(res.Models, bm.Model)
        }
        return res
    })
}

// Modify 修改电柜
func (s *cabinetService) Modify(req *model.CabinetModifyReq) {
    cab, _ := s.orm.QueryNotDeleted().Where(cabinet.ID(req.ID)).WithModels().First(s.ctx)
    if cab == nil {
        snag.Panic("未找到电柜")
    }
    willDeploy := model.CabinetStatus(cab.Status) == model.CabinetStatusPending && req.Status != nil && *req.Status == model.CabinetStatusNormal
    err := ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
        q := tx.Cabinet.UpdateOne(cab)
        if req.Models != nil {
            var models []string
            for _, bm := range cab.Edges.Models {
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
                q.ClearModels()
                // 查询设置电池型号
                q.AddModels(NewBatteryModel().QueryModelsX(*req.Models)...)
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
            if model.CabinetStatus(cab.Status) == model.CabinetStatusNormal || (req.Status != nil && *req.Status == model.CabinetStatusNormal) {
                return errors.New("电柜投产必须选择网点")
            }
        }
        if req.Status != nil {
            q.SetStatus(req.Status.Value())
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

        if req.Intelligent != nil {
            q.SetIntelligent(*req.Intelligent)
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
func (s *cabinetService) UpdateStatus(item *ent.Cabinet) (err error) {
    err = provider.NewUpdater(item).DoUpdate()
    // 如果返回失败, 则延迟2秒后重新请求一次
    if err != nil {
        time.Sleep(2 * time.Second)
    }
    return provider.NewUpdater(item).DoUpdate()
}

// DoorOpenStatus 获取柜门状态
func (s *cabinetService) DoorOpenStatus(item *ent.Cabinet, index int) ec.DoorStatus {
    _ = s.UpdateStatus(item)
    if len(item.Bin) == 0 || len(item.Bin) < index {
        return ec.DoorStatusUnknown
    }
    bin := item.Bin[index]
    if !bin.DoorHealth {
        return ec.DoorStatusFail
    }
    if bin.OpenStatus {
        return ec.DoorStatusOpen
    }
    return ec.DoorStatusClose
}

// DetailFromID 电柜详细信息
func (s *cabinetService) DetailFromID(id uint64) *model.CabinetDetailRes {
    item := s.orm.QueryNotDeleted().
        Where(cabinet.ID(id)).
        WithModels().
        OnlyX(s.ctx)
    if item == nil {
        snag.Panic("未找到电柜")
    }

    return s.Detail(item)
}

func (s *cabinetService) Detail(item *ent.Cabinet) *model.CabinetDetailRes {
    if !item.Intelligent && time.Now().Sub(item.UpdatedAt).Seconds() > 2 {
        err := s.UpdateStatus(item)
        if err != nil {
            snag.Panic(err)
        }
    }

    bms := item.Edges.Models
    if bms == nil {
        bms, _ = item.QueryModels().All(s.ctx)
    }

    res := new(model.CabinetDetailRes)
    _ = copier.Copy(res, item)
    for _, bm := range bms {
        res.Models = append(res.Models, bm.Model)
    }
    res.Reserves = make([]model.ReserveCabinetItem, 0)

    res.StockNum = NewStock().CurrentBattery(item.ID, stock.FieldCabinetID)

    // 获取生效中的预约
    revs := NewReserve().CabinetUnfinished(item.ID)
    for _, rev := range revs {
        res.Reserves = append(res.Reserves, model.ReserveCabinetItem{
            Name:     rev.Edges.Rider.Name,
            Phone:    rev.Edges.Rider.Phone,
            Business: rev.Type,
            Time:     rev.CreatedAt.Format(carbon.TimeLayout),
        })
    }

    return res
}

// DoorOperate 操作柜门
func (s *cabinetService) DoorOperate(req *model.CabinetDoorOperateReq, operator model.CabinetDoorOperator) (state bool, err error) {
    opId := shortuuid.New()
    now := time.Now()
    // 查找柜子和仓位
    item := s.QueryOne(*req.ID)
    if len(item.Bin) < *req.Index {
        err = errors.New("柜门不存在")
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
    // 请求开启柜门
    prov.PrepareRequest()
    state = prov.DoorOperate(operator.Name+"-"+opId, item.Serial, op, *req.Index)
    // 如果成功, 重新获取状态更新数据
    if state {
        // 更新一次电柜状态
        err = provider.NewUpdater(item).DoUpdate()
        log.Infof("%s操作成功[%s %s], Update: %v", item.Serial, req.Operation.String(), req.Remark, err)
        // 如果是锁仓, 需要更新仓位备注
        if *req.Operation == model.CabinetDoorOperateLock {
            item.Bin[*req.Index].Remark = req.Remark
        }
        // 如果是解锁, 需要清除仓位备注
        if *req.Operation == model.CabinetDoorOperateUnlock {
            item.Bin[*req.Index].Remark = ""
        }
        _, _ = item.Update().SetBin(item.Bin).Save(s.ctx)
    } else {
        err = errors.New("柜门操作失败")
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

// ModelInclude 电柜是否可用指定型号电池
func (s *cabinetService) ModelInclude(item *ent.Cabinet, model string) bool {
    bms := item.Edges.Models
    if bms == nil {
        bms, _ = item.QueryModels().All(s.ctx)
    }
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
func (s *cabinetService) Usable(cab *ent.Cabinet) (op *model.RiderCabinetOperateProcess) {
    max, empty := cab.Bin.MaxEmpty()
    op = &model.RiderCabinetOperateProcess{}

    min := cache.Float64(model.SettingExchangeMinBattery)
    if max.Electricity.Value() < min {
        snag.Panic("当前无可用电池")
    }

    if max == nil || empty == nil {
        snag.Panic("电柜异常, 无法换电")
    }

    op.EmptyBin = &model.BinInfo{
        Index: empty.Index,
    }

    maxInfo := &model.BinInfo{
        Index:       max.Index,
        Voltage:     max.Voltage,
        Electricity: max.Electricity,
    }

    if max.Electricity.IsBatteryFull() {
        op.FullBin = maxInfo
    } else {
        op.Alternative = maxInfo
    }

    return
}

// Businessable 判定电柜是否可用
func (s *cabinetService) Businessable(cab *ent.Cabinet) (health bool, maintenance bool) {
    maintenance = model.CabinetStatus(cab.Status) == model.CabinetStatusMaintenance
    health = model.CabinetStatus(cab.Status) == model.CabinetStatusNormal &&
        cab.Health == model.CabinetHealthStatusOnline &&
        time.Now().Sub(cab.UpdatedAt).Minutes() < 5 &&
        len(cab.Bin) > 0
    return
}

func (s *cabinetService) BusinessableX(cab *ent.Cabinet) {
    health, maintenance := s.Businessable(cab)
    if maintenance {
        snag.Panic("电柜开小差了, 请联系客服")
    }
    if !health {
        snag.Panic("电柜离线, 暂无法使用")
    }
}

func (s *cabinetService) Data(req *model.CabinetDataReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().WithModels().Order(ent.Desc(cabinet.FieldCreatedAt))
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

    if s.modifier != nil && s.modifier.Phone == "15537112255" {
        req.CityID = 410100
    }

    if req.CityID != 0 {
        q.Where(cabinet.CityID(req.CityID))
    }

    if req.Brand != "" {
        q.Where(cabinet.Brand(req.Brand.Value()))
    }

    if req.Votage != 0 {
        bm := fmt.Sprintf("%.0fV", req.Votage)
        q.Where(cabinet.HasModelsWith(batterymodel.ModelHasPrefix(bm)))
    }

    return s.dataItems(model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Cabinet) model.CabinetDataRes {
        return s.dataDetail(item)
    }))
}

func (s *cabinetService) dataItems(res *model.PaginationRes) *model.PaginationRes {
    items := res.Items.([]model.CabinetDataRes)
    ids := make([]uint64, len(items))
    for i, item := range items {
        ids[i] = item.ID
    }
    m := NewStock().CurrentBatteryNum(ids, stock.FieldCabinetID)
    for i, item := range items {
        items[i].StockNum = m[item.ID]
    }
    return res
}

func (s *cabinetService) dataDetail(item *ent.Cabinet) model.CabinetDataRes {
    res := model.CabinetDataRes{
        ID:         item.ID,
        Name:       item.Name,
        Serial:     item.Serial,
        Brand:      model.CabinetBrand(item.Brand).String(),
        Online:     item.Health == model.CabinetHealthStatusOnline,
        BinNum:     item.Doors,
        BatteryNum: item.BatteryNum,
    }

    bms := item.Edges.Models
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
}

// Transfer 电柜初始化调拨
func (s *cabinetService) Transfer(req *model.CabinetTransferReq) {
    cab := s.QueryOne(req.CabinetID)
    if cab.Transferred {
        snag.Panic("电柜已初始化过")
    }
    if cab.Health != model.CabinetHealthStatusOnline {
        snag.Panic("电柜不在线")
    }
    if !s.ModelInclude(cab, req.Model) {
        snag.Panic("电池型号错误")
    }
    NewStockWithModifier(s.modifier).Transfer(&model.StockTransferReq{
        Model:         req.Model,
        Num:           req.Num,
        InboundID:     cab.ID,
        InboundTarget: model.StockTargetCabinet,
        Force:         true,
        Remark:        "电柜初始化",
    })
    _, _ = cab.Update().SetTransferred(true).Save(s.ctx)
    return

}

// Sync 电柜同步
func (s *cabinetService) Sync(data *cabdef.CabinetMessage) {
    if data.Serial == "" {
        log.Error("[SYNC] 缺少参数 serial")
        return
    }

    cab, _ := s.orm.QueryNotDeleted().Where(cabinet.Serial(data.Serial)).WithModels().First(s.ctx)
    if cab == nil {
        log.Errorf("[SYNC] 未找到电柜信息, 请先添加电柜: %s", data.Serial)
        return
    }

    updater := cab.Update()

    defer func() {
        _ = updater.Exec(s.ctx)
    }()

    c := data.Cabinet
    if c != nil {
        health := model.CabinetHealthStatusOffline
        if c.Online {
            health = model.CabinetHealthStatusOnline
        }
        if c.Status == cabdef.StatusAbnormal {
            health = model.CabinetHealthStatusFault
        }
        updater.SetHealth(health)
    }

    var bins model.CabinetBins

    if !data.Full {
        bins = cab.Bin
    }

    if len(data.Bins) > 0 {
        var (
            bn, bf, bc, be, bl int
        )
        for _, b := range data.Bins {
            hasBattery := b.BatteryExists && b.BatterySn != ""
            var (
                isFull bool
                remark string
            )
            if b.Remark != nil {
                remark = *b.Remark
            }
            // 电池数
            if hasBattery {
                // 如果该仓位有电池
                _, _ = NewBattery().SyncPutin(b.BatterySn, cab.ID, b.Ordinal)
                bn += 1
                if b.Soc >= model.IntelligentBatteryFullSoc {
                    // 满电
                    bf += 1
                    isFull = true
                } else {
                    // 充电
                    bc += 1
                }
            } else {
                // 如果该仓位无电池
                NewBattery().SyncPutout(cab.ID, b.Ordinal)
                // 空仓
                be += 1
            }
            // 锁仓
            if !b.Enable {
                bl += 1
            }

            cb := &model.CabinetBin{
                Index:       b.Ordinal - 1,
                Name:        b.Name,
                BatterySN:   b.BatterySn,
                Full:        isFull,
                Battery:     hasBattery,
                Electricity: model.BatteryElectricity(b.Soc),
                OpenStatus:  b.Open,
                DoorHealth:  b.Health && b.Enable,
                Current:     b.Current,
                Voltage:     b.Voltage,
                Remark:      remark,
            }

            if data.Full || len(cab.Bin) < len(data.Bins) {
                bins = append(bins, cb)
            }

            if !data.Full {
                for i, xb := range bins {
                    if xb.Index+1 == b.Ordinal {
                        bins[i] = cb
                    }
                }
            }
        }

        sort.Slice(bins, func(i, j int) bool {
            return bins[i].Index <= bins[j].Index
        })

        updater.SetDoors(len(bins)).
            SetBatteryNum(bn).
            SetBatteryFullNum(bf).
            SetBatteryChargingNum(bc).
            SetEmptyBinNum(be).
            SetLockedBinNum(bl).
            SetBin(bins)
    }
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    stdsql "database/sql"
    "entgo.io/ent/dialect/sql"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/assets"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/exception"
    "github.com/auroraride/aurservd/internal/ent/predicate"
    "github.com/auroraride/aurservd/internal/ent/stock"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "math"
    "sort"
    "strings"
)

type stockService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.StockClient
}

func NewStock() *stockService {
    return &stockService{
        ctx: context.Background(),
        orm: ent.Database.Stock,
    }
}

func NewStockWithRider(r *ent.Rider) *stockService {
    s := NewStock()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewStockWithModifier(m *model.Modifier) *stockService {
    s := NewStock()
    if m != nil {
        s.ctx = context.WithValue(s.ctx, "modifier", m)
        s.modifier = m
    }
    return s
}

func NewStockWithEmployee(e *ent.Employee) *stockService {
    s := NewStock()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *stockService) StoreList(req *model.StockListReq) *model.PaginationRes {
    q := ent.Database.Store.QueryNotDeleted().
        Where(store.HasStocks()).
        WithCity().
        WithStocks().
        WithExceptions(func(eq *ent.ExceptionQuery) {
            eq.Where(exception.Status(model.ExceptionStatusPending))
        })
    if req.Name != nil {
        q.Where(
            store.NameContainsFold(*req.Name),
        )
    }
    if req.CityID != nil {
        q.Where(
            store.HasBranchWith(branch.CityID(*req.CityID)),
        )
    }
    if req.StoreID != nil {
        q.Where(
            store.ID(*req.StoreID),
        )
    }

    if req.Start != nil {
        start := carbon.Parse(*req.Start).StartOfDay().Carbon2Time()
        if start.IsZero() {
            snag.Panic("开始时间错误")
        }
        q.Where(
            store.Or(
                store.HasStocksWith(stock.CreatedAtGTE(start)),
            ),
        )
    }
    if req.End != nil {
        end := carbon.Parse(*req.End).EndOfDay().Carbon2Time()
        if end.IsZero() {
            snag.Panic("结束时间错误")
        }
        q.Where(
            store.Or(
                store.HasStocksWith(stock.CreatedAtLTE(end)),
            ),
        )
    }
    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Store) model.StockListRes {
            res := model.StockListRes{
                Store: model.Store{
                    ID:   item.ID,
                    Name: item.Name,
                },
                City: model.City{
                    ID:   item.Edges.City.ID,
                    Name: item.Edges.City.Name,
                },
                BatteryTotal: 0,
                Batteries:    make([]*model.StockMaterial, 0),
                Materials:    make([]*model.StockMaterial, 0),
            }

            // 计算所有物资
            batteries := make(map[string]*model.StockMaterial)
            materials := make(map[string]*model.StockMaterial)

            // 出入库
            for _, st := range item.Edges.Stocks {
                if st.Model != nil {
                    // 电池
                    s.calculate(batteries, st)
                } else {
                    // 物资
                    s.calculate(materials, st)
                }
            }

            for _, ex := range item.Edges.Exceptions {
                if ex.Model != nil {
                    s.calculateException(batteries, ex)
                } else {
                    s.calculateException(materials, ex)
                }
            }

            for _, battery := range batteries {
                res.Batteries = append(res.Batteries, battery)
                res.BatteryTotal += battery.Surplus
            }

            for _, material := range materials {
                res.Materials = append(res.Materials, material)
            }

            // 排序
            sort.Slice(res.Batteries, func(i, j int) bool {
                return strings.Compare(res.Batteries[i].Name, res.Batteries[j].Name) < 0
            })
            sort.Slice(res.Materials, func(i, j int) bool {
                return strings.Compare(res.Materials[i].Name, res.Materials[j].Name) < 0
            })

            return res
        },
    )
}

func (s *stockService) calculateException(items map[string]*model.StockMaterial, ex *ent.Exception) {
    name := ex.Name
    if _, ok := items[name]; !ok {
        items[name] = &model.StockMaterial{
            Name:     name,
            Outbound: 0,
            Inbound:  0,
            Surplus:  0,
        }
    }
    items[name].Exception += ex.Num
}

func (s *stockService) calculate(items map[string]*model.StockMaterial, st *ent.Stock) {
    name := st.Name
    if _, ok := items[name]; !ok {
        items[name] = &model.StockMaterial{
            Name:     name,
            Outbound: 0,
            Inbound:  0,
            Surplus:  0,
        }
    }
    if st.Num > 0 {
        items[name].Inbound += st.Num
    } else {
        items[name].Outbound += int(math.Abs(float64(st.Num)))
    }
    items[name].Surplus += st.Num
}

// Fetch 获取对应物资库存
// TODO 车架调拨另外做
func (s *stockService) Fetch(target uint8, id uint64, name string) int {
    var result []struct {
        Sum       int    `json:"sum"`
        StoreID   uint64 `json:"store_id"`
        CabinetID uint64 `json:"cabinet_id"`
    }

    var idw predicate.Stock
    switch target {
    case model.StockTargetStore:
        idw = stock.StoreID(id)
        break
    case model.StockTargetCabinet:
        idw = stock.CabinetID(id)
        break
    }
    q := s.orm.QueryNotDeleted().
        Where(stock.Name(name), idw).
        GroupBy(stock.FieldStoreID, stock.FieldCabinetID).
        Aggregate(ent.Sum(stock.FieldNum))
    err := q.Scan(s.ctx, &result)
    if err != nil {
        log.Error(err)
        snag.Panic("物资数量获取失败")
    }

    if result == nil || len(result) < 0 {
        return 0
    }
    return result[0].Sum
}

// Transfer 调拨
func (s *stockService) Transfer(req *model.StockTransferReq) {
    if req.Name == "" && req.Model == "" {
        snag.Panic("电池型号和物资名称不能同时为空")
    }
    if req.Name != "" && req.Model != "" {
        snag.Panic("电池型号和物资名称不能同时存在")
    }
    if req.Num <= 0 {
        snag.Panic("调拨物资数量错误")
    }
    if req.InboundID == 0 && req.OutboundID == 0 {
        snag.Panic("平台之间无法调拨物资")
    }
    if req.InboundTarget == model.StockTargetCabinet && req.OutboundTarget == model.StockTargetCabinet {
        snag.Panic("电柜之间无法调拨")
    }
    if (req.InboundTarget == model.StockTargetStore && req.InboundID == 0) || (req.InboundTarget == model.StockTargetPlaform && req.InboundID != 0) {
        snag.Panic("调入参数错误")
    }
    if (req.OutboundTarget == model.StockTargetStore && req.OutboundID == 0) || (req.OutboundTarget == model.StockTargetPlaform && req.OutboundID != 0) {
        snag.Panic("调出参数错误")
    }

    name := req.Name
    mt := stock.MaterialOthers
    if req.Model != "" {
        name = req.Model
        mt = stock.MaterialBattery
    }

    var cityID uint64

    // 查询电柜
    var cab *ent.Cabinet
    var cabID uint64

    // 检查电柜是否初始化调拨过
    if req.InboundTarget == model.StockTargetCabinet {
        cabID = req.InboundID
    }
    if req.OutboundTarget == model.StockTargetCabinet {
        cabID = req.OutboundID
    }
    if cabID > 0 {
        cab = NewCabinet().QueryOne(cabID)
        if !cab.Transferred && !req.Force {
            snag.Panic("电柜未初始化调拨")
        }
        if cab.CityID != nil {
            cityID = *cab.CityID
        }
    }

    // 查询门店
    var st *ent.Store
    var stID uint64
    if req.InboundTarget == model.StockTargetStore {
        stID = req.InboundID
    }
    if req.OutboundTarget == model.StockTargetStore {
        stID = req.OutboundID
    }
    if stID > 0 {
        st = NewStore().Query(stID)
        cityID = st.CityID
    }
    if cab != nil && cab.CityID != nil && st != nil && st.CityID != *cab.CityID {
        snag.Panic("不同城市电柜和门店无法调拨")
    }

    // 调出检查
    if req.OutboundID > 0 && s.Fetch(req.OutboundTarget, req.OutboundID, name) < req.Num {
        snag.Panic("操作失败, 调出物资大于库存物资")
    }

    sn := tools.NewUnique().NewSN()

    in := &req.InboundID
    if req.InboundID == 0 {
        in = nil
    }

    out := &req.OutboundID
    if req.OutboundID == 0 {
        out = nil
    }

    v := &req.Model
    if req.Model == "" {
        v = nil
    }

    err := ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
        // 调出
        oq := tx.Stock.Create().
            SetName(name).
            SetNillableModel(v).
            SetNum(-req.Num).
            SetCityID(cityID).
            SetType(model.StockTypeTransfer).
            SetMaterial(mt).
            SetNillableRemark(req.Remark).
            SetSn(sn)

        if req.OutboundTarget == model.StockTargetStore {
            oq.SetNillableStoreID(out)
        } else {
            oq.SetNillableCabinetID(out)
        }

        // 调入
        iq := tx.Stock.Create().
            SetName(name).
            SetNillableModel(v).
            SetNum(req.Num).
            SetCityID(cityID).
            SetType(model.StockTypeTransfer).
            SetMaterial(mt).
            SetNillableRemark(req.Remark).
            SetSn(sn)
        if req.InboundTarget == model.StockTargetStore {
            iq.SetNillableStoreID(in)
        } else {
            iq.SetNillableCabinetID(in)
        }

        // 调出
        var spouse *ent.Stock
        spouse, err = oq.Save(s.ctx)
        if err != nil {
            return
        }

        // 调入
        _, err = iq.SetSpouse(spouse).Save(s.ctx)
        if err != nil {
            return
        }
        return
    })
    if err != nil {
        snag.Panic(err)
    }
}

func (s *stockService) BatteryOverview(req *model.StockOverviewReq) (items []model.StockBatteryOverviewRes) {
    var extends []string

    switch req.Goal {
    case 1:
        if req.StoreID != 0 {
            extends = append(extends, fmt.Sprintf("AND (%s = %d)", stock.FieldStoreID, req.StoreID))
        } else {
            extends = append(extends, fmt.Sprintf("AND (%s IS NOT NULL OR (%s IS NULL AND %s IS NULL AND %s > 0))", stock.FieldStoreID, stock.FieldStoreID, stock.FieldCabinetID, stock.FieldType))
        }
        break
    case 2:
        if req.CabinetID != 0 {
            extends = append(extends, fmt.Sprintf("AND (%s = %d)", stock.FieldCabinetID, req.StoreID))
        } else {
            extends = append(extends, fmt.Sprintf("AND (%s IS NOT NULL)", stock.FieldCabinetID))
        }
        break
    default:
        extends = append(extends, fmt.Sprintf("AND (%s IS NOT NULL OR %s IS NOT NULL OR %s IS NOT NULL)", stock.FieldStoreID, stock.FieldCabinetID, stock.FieldRiderID))
        break
    }

    if req.CityID != 0 {
        extends = append(extends, fmt.Sprintf("AND (%s = %d)", stock.FieldCityID, req.CityID))
    }

    if req.Start != "" && req.End != "" {
        start := tools.NewTime().ParseDateStringX(req.Start).Format(carbon.DateTimeLayout)
        end := tools.NewTime().ParseNextDateStringX(req.End).Format(carbon.DateTimeLayout)
        extends = append(extends, fmt.Sprintf("AND (%s >= '%s'::timestamp AND %s < '%s'::timestamp)", stock.FieldCreatedAt, start, stock.FieldCreatedAt, end))
    }

    extend := fmt.Sprintf("WHERE model IS NOT NULL %s", strings.Join(extends, " "))
    query := strings.Replace(assets.SQLStockOverview, "WHERE model IS NOT NULL", extend, 1)

    rows, err := ent.Database.QueryContext(s.ctx, query)

    if err != nil {
        log.Error(err)
        snag.Panic("请求失败")
    }

    defer func(rows *stdsql.Rows) {
        _ = rows.Close()
    }(rows)

    for rows.Next() {
        var b []byte
        err = rows.Scan(&b)
        if err != nil {
            log.Error(err)
            break
        }
        var item model.StockBatteryOverviewRes
        _ = jsoniter.Unmarshal(b, &item)
        items = append(items, item)
    }

    return
}

// BatteryWithRider 和骑手交互电池出入库
func (s *stockService) BatteryWithRider(cr *ent.StockCreate, req *model.StockBusinessReq) error {
    num := model.StockNumberOfRiderBusiness(req.StockType)

    if req.StoreID == nil && req.CabinetID == nil {
        return errors.New("参数校验错误")
    }

    // TODO 平台管理员可操作性时处理出入库逻辑
    if req.StoreID != nil {
        cr.SetStoreID(*req.StoreID)
        if num < 0 && s.Fetch(model.StockTargetStore, *req.StoreID, req.Model) < int(math.Abs(float64(num))) {
            return errors.New("电池库存不足")
        }
    }

    if req.CabinetID != nil {
        cr.SetCabinetID(*req.CabinetID)
        if num < 0 && s.Fetch(model.StockTargetCabinet, *req.CabinetID, req.Model) < int(math.Abs(float64(num))) {
            return errors.New("电池库存不足")
        }
    }

    cr.SetNillableEmployeeID(req.EmployeeID).
        SetName(req.Model).
        SetRiderID(req.RiderID).
        SetType(req.StockType).
        SetModel(req.Model).
        SetNum(num).
        SetCityID(req.CityID).
        SetNillableSubscribeID(req.SubscribeID).
        SetMaterial(stock.MaterialBattery).
        SetSn(tools.NewUnique().NewSN())

    _, err := cr.Save(s.ctx)

    if err != nil {
        log.Error(err)
    }
    return err
}

// EmployeeOverview 店员物资概览
func (s *stockService) EmployeeOverview() (res model.StockEmployeeOverview) {
    st := s.employee.Edges.Store
    if st == nil {
        snag.Panic("未上班")
    }

    start := carbon.Now().StartOfDay().Timestamp()

    res = model.StockEmployeeOverview{
        Batteries: make([]*model.StockEmployeeOverviewBattery, 0),
        Materials: make([]*model.StockEmployeeOverviewMaterial, 0),
    }

    // 计算所有物资
    batteries := make(map[string]*model.StockEmployeeOverviewBattery)
    materials := make(map[string]*model.StockEmployeeOverviewMaterial)

    items, _ := s.orm.QueryNotDeleted().Where(stock.StoreID(st.ID)).All(s.ctx)
    for _, item := range items {
        name := st.Name
        if item.Model != nil {
            if _, ok := batteries[name]; !ok {
                batteries[name] = &model.StockEmployeeOverviewBattery{
                    Outbound: 0,
                    Inbound:  0,
                    Surplus:  0,
                    Model:    *item.Model,
                }
            }
            // 判断是否今日
            if item.CreatedAt.Unix() > start {
                if item.Num > 0 {
                    batteries[name].Inbound += item.Num
                } else {
                    batteries[name].Outbound += int(math.Abs(float64(item.Num)))
                }
            }
            batteries[name].Surplus += item.Num
        } else {
            if _, ok := materials[name]; !ok {
                materials[name] = &model.StockEmployeeOverviewMaterial{
                    Name:    name,
                    Surplus: 0,
                }
            }
            materials[name].Surplus += item.Num
        }
    }

    for _, battery := range batteries {
        res.Batteries = append(res.Batteries, battery)
    }

    for _, material := range materials {
        res.Materials = append(res.Materials, material)
    }

    // 排序
    sort.Slice(res.Batteries, func(i, j int) bool {
        return strings.Compare(res.Batteries[i].Model, res.Batteries[j].Model) < 0
    })
    sort.Slice(res.Materials, func(i, j int) bool {
        return strings.Compare(res.Materials[i].Name, res.Materials[j].Name) < 0
    })

    return
}

// listBasicQuery 列表基础查询语句
func (s *stockService) listBasicQuery(req *model.StockEmployeeListReq) *ent.StockQuery {
    tt := tools.NewTime()

    q := s.orm.QueryNotDeleted().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        })

    if req.Outbound {
        q.Where(stock.NumLT(0))
    } else {
        q.Where(stock.NumGT(0))
    }

    if req.Start != nil {
        q.Where(stock.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }

    if req.End != nil {
        q.Where(stock.CreatedAtLTE(tt.ParseNextDateStringX(*req.End)))
    }

    return q
}

func (s *stockService) EmployeeList(req *model.StockEmployeeListReq) model.StockEmployeeListRes {
    st := s.employee.Edges.Store
    if st == nil {
        snag.Panic("未上班")
    }
    q := s.listBasicQuery(req).Where(stock.StoreID(st.ID))

    res := model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Stock) model.StockEmployeeListResItem {
            r := item.Edges.Rider
            res := model.StockEmployeeListResItem{
                ID:    item.ID,
                Type:  item.Type,
                Model: item.Model,
                Num:   item.Num,
                Time:  item.CreatedAt.Format(carbon.DateTimeLayout),
            }
            if r != nil {
                res.Phone = r.Phone
                if p := r.Edges.Person; p != nil {
                    res.Name = p.Name
                }
            }
            return res
        },
    )

    var today *int

    if res.Pagination.Current == 1 {
        today = new(int)
        // 获取今日数量
        var result []struct {
            ID  uint64 `json:"id"`
            Sum int    `json:"sum"`
        }
        cq := s.orm.QueryNotDeleted().Where(
            stock.CreatedAtGTE(carbon.Now().StartOfDay().Carbon2Time()),
            stock.StoreID(st.ID),
        )
        if req.Outbound {
            cq.Where(stock.NumLT(0))
        } else {
            cq.Where(stock.NumGT(0))
        }
        _ = cq.Modify().GroupBy(stock.FieldID).Aggregate(ent.Sum(stock.FieldNum)).Scan(s.ctx, &result)

        if len(result) > 0 {
            today = &result[0].Sum
        }
    }

    return model.StockEmployeeListRes{
        Today:         today,
        PaginationRes: res,
    }
}

// StoreCurrent 列出当前门店所有库存物资
func (s *stockService) StoreCurrent(id uint64) []model.InventoryItemWithNum {
    ins := make([]model.InventoryItemWithNum, 0)
    err := s.orm.QueryNotDeleted().
        Where(stock.StoreID(id)).
        Modify(func(sel *sql.Selector) {
            sel.GroupBy(stock.FieldName, stock.FieldModel).Select(stock.FieldName, stock.FieldModel).
                AppendSelectExprAs(sql.Raw(fmt.Sprintf("%s IS NOT NULL", stock.FieldModel)), "battery").
                AppendSelectExprAs(sql.Raw(fmt.Sprintf("SUM(%s)", stock.FieldNum)), "num")
        }).
        Scan(s.ctx, &ins)

    if err != nil {
        log.Error(err)
    }

    return ins
}

func (s *stockService) StoreCurrentMap(id uint64) map[string]model.InventoryItemWithNum {
    inm := make(map[string]model.InventoryItemWithNum)
    for _, in := range s.StoreCurrent(id) {
        inm[in.Name] = in
    }
    return inm
}

func (s *stockService) CabinetList(req *model.StockCabinetListReq) *model.PaginationRes {
    q := ent.Database.Cabinet.QueryNotDeleted().
        Where(cabinet.HasStocks()).
        WithCity().
        WithStocks()

    if req.CabinetID != 0 {
        q.Where(cabinet.ID(req.CabinetID))
    }
    if req.Serial != "" {
        q.Where(cabinet.Serial(req.Serial))
    }
    if req.CityID != 0 {
        q.Where(cabinet.CityID(req.CityID))
    }
    if req.Start != "" {
        q.Where(cabinet.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
    }
    if req.End != "" {
        q.Where(cabinet.CreatedAtLT(tools.NewTime().ParseNextDateStringX(req.End)))
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Cabinet) model.StockCabinetListRes {
        res := model.StockCabinetListRes{
            ID:        item.ID,
            Serial:    item.Serial,
            Name:      item.Name,
            Batteries: make([]*model.StockMaterial, 0),
        }
        c := item.Edges.City
        if c != nil {
            res.City = model.City{
                ID:   c.ID,
                Name: c.Name,
            }
        }
        batteries := make(map[string]*model.StockMaterial)

        // 出入库
        for _, st := range item.Edges.Stocks {
            s.calculate(batteries, st)
        }

        for _, battery := range batteries {
            res.Batteries = append(res.Batteries, battery)
        }

        return res
    })
}

// Detail 出入库明细
func (s *stockService) Detail(req *model.StockDetailReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().WithCabinet().WithStore().WithSpouse(func(sq *ent.StockQuery) {
        sq.WithStore().WithCabinet().WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        })
    }).WithRider(func(rq *ent.RiderQuery) {
        rq.WithPerson()
    }).WithEmployee().WithCity()
    // 排序
    if req.Positive {
        q.Order(ent.Asc(stock.FieldSn))
    } else {
        q.Order(ent.Desc(stock.FieldSn))
    }

    if req.Start != "" {
        q.Where(stock.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
    }

    if req.End != "" {
        q.Where(stock.CreatedAtLT(tools.NewTime().ParseNextDateStringX(req.End)))
    }

    if req.CityID != 0 {
        q.Where(stock.CityID(req.CityID))
    }

    if req.Serial != "" {
        q.Where(stock.HasCabinetWith(cabinet.Serial(req.Serial)))
    }

    switch req.Goal {
    case 1:
        // 门店物资
        if req.StoreID != 0 {
            q.Where(stock.StoreID(req.StoreID))
        } else {
            q.Where(
                stock.StoreIDNotNil(),
                stock.CabinetIDIsNil(),
            )
        }
        break
    case 2:
        // 电柜物资
        if req.CabinetID != 0 {
            q.Where(stock.CabinetID(req.CabinetID))
        } else {
            q.Where(
                stock.StoreIDIsNil(),
                stock.CabinetIDNotNil(),
            )
        }
        break
    default:
        // 门店或电柜物资
        q.Where(
            stock.Or(
                stock.StoreIDNotNil(),
                stock.CabinetIDNotNil(),
            ),
        )
        break
    }

    // 筛选物资类别
    if req.Materials == "" {
        req.Materials = string(stock.MaterialBattery)
    } else {
        req.Materials = strings.ReplaceAll(req.Materials, " ", "")
    }
    materials := strings.Split(req.Materials, ",")
    var predicates []predicate.Stock
    for _, material := range materials {
        switch stock.Material(material) {
        case stock.MaterialBattery:
            predicates = append(predicates, stock.ModelNotNil())
            break
        case stock.MaterialOthers:
            predicates = append(predicates, stock.ModelIsNil())
            break
        }
    }
    q.Where(stock.Or(predicates...))

    q.Modify(func(sel *sql.Selector) {
        sel.Select("ON (sn) *")
    })

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Stock) model.StockDetailRes {
        return s.detailInfo(item)
    })
}

// detailInfo 库存出入明细信息
func (s *stockService) detailInfo(item *ent.Stock) model.StockDetailRes {
    res := model.StockDetailRes{
        ID:     item.ID,
        Sn:     item.Sn,
        Name:   item.Name,
        Num:    int(math.Abs(float64(item.Num))),
        Time:   item.CreatedAt.Format(carbon.DateTimeLayout),
        Remark: item.Remark,
    }

    // 城市
    c := item.Edges.City
    if c != nil {
        res.City = c.Name
    }

    em := item.Creator
    er := item.Edges.Rider
    ee := item.Edges.Employee
    es := item.Edges.Store
    ec := item.Edges.Cabinet

    if item.Type == model.StockTypeTransfer {
        // 平台调拨记录
        res.Type = "平台调拨"
        res.Operator = fmt.Sprintf("后台 - %s", em.Name)

        var ses *ent.Store
        var sec *ent.Cabinet

        sp := item.Edges.Spouse
        if sp != nil {
            ses = sp.Edges.Store
            sec = sp.Edges.Cabinet
        }

        // 出入库对象判定
        if item.Num > 0 {
            res.Inbound = s.target(es, ec)
            res.Outbound = s.target(ses, sec)
        } else {
            res.Inbound = s.target(ses, sec)
            res.Outbound = s.target(es, ec)
        }
    } else {
        // 业务调拨记录
        var riderName string

        if er != nil {
            riderName = er.Edges.Person.Name
            res.Rider = fmt.Sprintf("%s - %s", riderName, er.Phone)
        }

        tm := map[uint8]string{
            model.StockTypeRiderObtain:      "新签",
            model.StockTypeRiderPause:       "寄存",
            model.StockTypeRiderContinue:    "取消寄存",
            model.StockTypeRiderUnSubscribe: "退租",
        }

        var tmr string
        if ec != nil {
            res.Operator = fmt.Sprintf("骑手 - %s", riderName)
            tmr = "电柜"
        } else {
            if ee != nil {
                tmr = "门店"
                res.Operator = fmt.Sprintf("店员 - %s", ee.Name)
            } else if item.Creator != nil {
                tmr = "后台"
                res.Operator = fmt.Sprintf("后台 - %s", item.Creator.Name)
            }
        }

        res.Type = tmr + tm[item.Type]

        // 出入库对象
        target := fmt.Sprintf("[骑手] %s - %s", er.Phone, er.Edges.Person.Name)
        switch item.Type {
        case model.StockTypeRiderObtain, model.StockTypeRiderContinue:
            res.Inbound = target
            res.Outbound = s.target(es, ec)
            break
        case model.StockTypeRiderPause, model.StockTypeRiderUnSubscribe:
            res.Inbound = s.target(es, ec)
            res.Outbound = target
            break
        }
    }

    return res
}

// target 出入库对象
func (s *stockService) target(es *ent.Store, ec *ent.Cabinet) (target string) {
    target = "平台"
    if es != nil {
        target = fmt.Sprintf("[门店] %s", es.Name)
    }
    if ec != nil {
        target = fmt.Sprintf("[电柜] %s - %s", ec.Name, ec.Serial)
    }
    return
}

// CurrentNum 获取当前库存
func (s *stockService) CurrentNum(ids []uint64, field string) map[uint64]int {
    var result []struct {
        TargetID uint64 `json:"target_id"`
        Sum      int    `json:"sum"`
    }
    v := make([]interface{}, len(ids))
    for i := range v {
        v[i] = ids[i]
    }
    _ = s.orm.QueryNotDeleted().
        Modify(func(sel *sql.Selector) {
            sel.Where(sql.In(sel.C(field), v...)).
                Select(
                    sql.As(sel.C(field), "target_id"),
                    sql.As(sql.Sum(stock.FieldNum), "sum"),
                ).
                GroupBy(field)
        }).
        Scan(s.ctx, &result)
    m := make(map[uint64]int)
    for _, r := range result {
        m[r.TargetID] = r.Sum
    }
    return m
}

func (s *stockService) Current(id uint64, field string) int {
    return s.CurrentNum([]uint64{id}, field)[id]
}

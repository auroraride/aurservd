// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    stdsql "database/sql"
    "errors"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/stock"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "math"
    "sort"
    "strconv"
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
        orm: ar.Ent.Stock,
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
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewStockWithEmployee(e *ent.Employee) *stockService {
    s := NewStock()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *stockService) List(req *model.StockListReq) *model.PaginationRes {
    q := ar.Ent.Store.QueryNotDeleted().
        Where(
            store.Or(
                store.HasStocks(),
            ),
        ).
        WithCity().
        WithStocks()
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

            // 入库
            for _, st := range item.Edges.Stocks {
                if st.Voltage != nil {
                    // 电池
                    s.calculate(batteries, st)
                } else {
                    // 物资
                    s.calculate(materials, st)
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

// Fetch 获取门店对应物资库存
func (s *stockService) Fetch(storeID uint64, name string) int {
    var result []struct {
        Sum     int    `json:"sum"`
        StoreID uint64 `json:"store_id"`
    }
    q := s.orm.QueryNotDeleted().
        Where(stock.Name(name), stock.StoreID(storeID)).
        GroupBy(stock.FieldStoreID).
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
    if req.Name == "" && req.Voltage == 0 {
        snag.Panic("电池型号和物资名称不能同时为空")
    }
    if req.Name != "" && req.Voltage != 0 {
        snag.Panic("电池型号和物资名称不能同时存在")
    }
    if req.Num <= 0 {
        snag.Panic("调拨物资数量错误")
    }
    if req.InboundID == 0 && req.OutboundID == 0 {
        snag.Panic("平台之间无法调拨物资")
    }

    tx, _ := ar.Ent.Tx(s.ctx)

    // 调出检查
    name := req.Name
    if req.Voltage != 0 {
        name = NewBattery().VoltageName(req.Voltage)
    }

    if req.OutboundID > 0 && s.Fetch(req.OutboundID, name) < req.Num {
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

    v := &req.Voltage
    if req.Voltage == 0 {
        v = nil
    }

    // 调出
    _, err := tx.Stock.Create().
        SetName(name).
        SetNillableVoltage(v).
        SetNum(-req.Num).
        SetNillableStoreID(out).
        SetType(model.StockTypeTransfer).
        SetSn(sn).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 调入
    _, err = tx.Stock.Create().
        SetName(name).
        SetNillableVoltage(v).
        SetNum(req.Num).
        SetNillableStoreID(in).
        SetType(model.StockTypeTransfer).
        SetSn(sn).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _ = tx.Commit()
}

func (s *stockService) Overview() (res model.StockOverview) {
    rows, err := ar.Ent.QueryContext(s.ctx, `SELECT DISTINCT ABS(SUM(num)) AS sum,
                NOT store_id IS NULL AND num < 0 AS outbound,
                NOT store_id IS NULL AND num > 0 AS inbound,
                store_id IS NULL                 AS plaform
FROM stock
WHERE voltage IS NOT NULL AND deleted_at IS NULL
GROUP BY outbound, inbound, plaform`)

    if err != nil {
        log.Error(err)
        snag.Panic("请求失败")
    }

    defer func(rows *stdsql.Rows) {
        _ = rows.Close()
    }(rows)

    for rows.Next() {
        var sum string
        var outbound, inbound, plaform bool
        err = rows.Scan(&sum, &outbound, &inbound, &plaform)
        if err != nil {
            log.Error(err)
            break
        }
        total, _ := strconv.Atoi(sum)
        if outbound {
            res.Outbound = total
        }
        if inbound {
            res.Inbound = total
        }
        if plaform {
            res.Total = total
        }
    }

    res.Surplus = res.Inbound - res.Outbound
    return
}

// BatteryWithRider 和骑手交互电池出入库
func (s *stockService) BatteryWithRider(cr *ent.StockCreate, req *model.StockWithRiderReq) error {
    name := NewBattery().VoltageName(req.Voltage)

    num := model.StockNumberOfRiderBusiness(req.StockType)

    // TODO 平台管理员可操作性时处理出入库逻辑
    if req.StoreID != 0 {
        cr.SetStoreID(req.StoreID)
        if num < 0 && s.Fetch(req.StoreID, name) < int(math.Abs(float64(num))) {
            return errors.New("电池库存不足")
        }
    }

    if req.EmployeeID != 0 {
        cr.SetEmployeeID(req.EmployeeID)
    }

    if req.ManagerID != 0 {
        cr.SetManagerID(req.ManagerID)
    }

    cr.SetName(NewBattery().VoltageName(req.Voltage)).
        SetRiderID(req.RiderID).
        SetType(req.StockType).
        SetVoltage(req.Voltage).
        SetNum(num).
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
        if item.Voltage != nil {
            if _, ok := batteries[name]; !ok {
                batteries[name] = &model.StockEmployeeOverviewBattery{
                    Outbound: 0,
                    Inbound:  0,
                    Surplus:  0,
                    Voltage:  *item.Voltage,
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
        return res.Batteries[i].Voltage < res.Batteries[j].Voltage
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
        q.Where(stock.CreatedAtLTE(tt.ParseDateStringX(*req.End)))
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
                ID:      item.ID,
                Type:    item.Type,
                Voltage: item.Voltage,
                Num:     item.Num,
                Time:    item.CreatedAt.Format(carbon.DateTimeLayout),
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

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    stdsql "database/sql"
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
        WithBranch(func(bq *ent.BranchQuery) {
            bq.WithCity()
        }).
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
                    ID:   item.Edges.Branch.Edges.City.ID,
                    Name: item.Edges.Branch.Edges.City.Name,
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
        Where(stock.Name(name), stock.StoreID(storeID), stock.NumGT(0)).
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
    if req.InboundID < 0 {
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

// BatteryOutboundWithRider 和骑手交互电池出入库
func (s *stockService) BatteryOutboundWithRider(cr *ent.StockCreate, riderID, storeID, employeeID uint64, voltage float64, stockType uint8) error {
    name := NewBattery().VoltageName(voltage)

    var num int
    switch stockType {
    case model.StockTypeRiderObtain:
        num = -1
        break
    case model.StockTypeRiderPause, model.StockTypeRiderUnSubscribe:
        num = 1
        break
    }

    if s.Fetch(storeID, name) < int(math.Abs(float64(num))) {
        snag.Panic("电池库存不足")
    }
    _, err := cr.SetName(NewBattery().VoltageName(voltage)).
        SetRiderID(riderID).
        SetEmployeeID(employeeID).
        SetStoreID(storeID).
        SetType(stockType).
        SetVoltage(voltage).
        SetNum(num).
        SetSn(tools.NewUnique().NewSN()).
        Save(s.ctx)
    if err != nil {
        log.Error(err)
    }
    return err
}

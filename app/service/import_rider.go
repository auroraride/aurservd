// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "encoding/csv"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "io"
    "os"
    "strconv"
    "strings"
    "time"
)

type importRiderService struct {
    modifier *model.Modifier
    ctx      context.Context
    plan     *ent.Plan
}

func NewImportRider() *importRiderService {
    return &importRiderService{
        ctx: context.Background(),
    }
}

func NewImportRiderWithModifier(m *model.Modifier) *importRiderService {
    s := NewImportRider()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *importRiderService) ParseCSV(path string) {
    csvfile, _ := os.Open(path)

    defer func(csvfile *os.File) {
        _ = csvfile.Close()
    }(csvfile)

    r := csv.NewReader(csvfile)

    i := 0

    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        i += 1
        if i == 1 {
            continue
        }

        if record[0] == "" {
            continue
        }

        t := strings.TrimSpace(record[7])
        var end string
        if strings.Contains(t, "/") {
            arr := strings.Split(t, "/")
            for i, str := range arr {
                if i == 0 {
                    continue
                }
                if len(str) != 2 {
                    str = "0" + str
                }
                arr[i] = str
            }
            end = strings.Join(arr, "-")
        }
        end = tools.NewTime().ParseDateStringX(end).Format(carbon.DateLayout)

        item := &model.ImportRiderFromCsv{
            Name:  strings.TrimSpace(record[0]),
            Phone: strings.TrimSpace(record[1]),
            Plan:  strings.TrimSpace(record[2]),
            Days:  strings.TrimSpace(record[3]),
            Model: strings.ToUpper(strings.TrimSpace(record[4])),
            City:  strings.TrimSpace(record[5]),
            Store: strings.TrimSpace(record[6]),
            End:   end,
        }

        // 查找骑行卡
        days, _ := strconv.Atoi(item.Days)
        s.plan = ent.Database.Plan.QueryNotDeleted().Where(plan.Name(item.Plan), plan.Days(uint(days))).FirstX(s.ctx)

        // 查找门店
        qs := ent.Database.Store.QueryNotDeleted().Where(store.Name(item.Store)).FirstX(s.ctx)
        if qs == nil {
            snag.Panic(fmt.Sprintf("未找到门店: %s", item.Store))
        }

        // 查找城市
        qc := ent.Database.City.QueryNotDeleted().Where(city.Name(item.City)).FirstX(s.ctx)
        err = s.Create(&model.ImportRiderCreateReq{
            Name:       item.Name,
            Phone:      item.Phone,
            PlanID:     s.plan.ID,
            CityID:     qc.ID,
            StoreID:    qs.ID,
            EmployeeID: 38654705685,
            End:        end,
            Model:      item.Model,
        })
        if err != nil {
            log.Errorf("[%s] 添加失败: %s", item, err)
        }
    }
}

// Create 手动添加骑手
func (s *importRiderService) Create(req *model.ImportRiderCreateReq) error {
    return ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
        var p *ent.Person
        var r *ent.Rider
        var o *ent.Order
        var sub *ent.Subscribe

        if r, _ = ent.Database.Rider.QueryNotDeleted().WithPerson().Where(rider.Phone(req.Phone)).First(s.ctx); r != nil {
            if exist, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.RiderID(r.ID)).Exist(s.ctx); exist {
                snag.Panic(fmt.Sprintf("%s:%s 已存在", req.Phone, req.Name))
            }
        }

        // 结束时间
        end := carbon.Parse(req.End)

        // 查询plan
        if s.plan == nil {
            s.plan = ent.Database.Plan.QueryNotDeleted().Where(plan.ID(req.PlanID)).FirstX(s.ctx)
        }

        // 计算开始时间
        start := end.SubDays(int(s.plan.Days)).Tomorrow().Carbon2Time()

        // 创建用户
        if r == nil || r.Edges.Person == nil {
            p, err = tx.Person.Create().SetName(req.Name).Save(s.ctx)
            if err != nil {
                return
            }
        }
        if r == nil {
            // 创建骑手并设置为不需要签约
            r, err = tx.Rider.Create().SetPhone(req.Phone).SetPerson(p).SetContractual(true).Save(s.ctx)
            if err != nil {
                return
            }
        } else if r.Edges.Person == nil {
            r, err = tx.Rider.UpdateOne(r).SetPerson(p).SetContractual(true).Save(s.ctx)
            if err != nil {
                return
            }
        }

        // 添加订阅
        sub, err = tx.Subscribe.Create().
            SetEmployeeID(req.EmployeeID).
            SetRider(r).
            SetInitialDays(int(s.plan.Days)).
            SetType(model.OrderTypeNewly).
            SetStatus(model.SubscribeStatusUsing).
            SetStartAt(start).
            SetStoreID(req.StoreID).
            SetPlanID(req.PlanID).
            SetCityID(req.CityID).
            SetModel(req.Model).
            SetRemaining(tools.NewTime().LastDays(end.Carbon2Time(), time.Now())).
            Save(s.ctx)
        if err != nil {
            return
        }

        // 创建订单
        o, err = tx.Order.
            Create().
            SetRiderID(r.ID).
            SetSubscribeID(sub.ID).
            SetStatus(model.OrderStatusPaid).
            SetRemark("导入数据").
            SetPayway(model.OrderPaywayManual).
            SetOutTradeNo(tools.NewUnique().NewSonyflakeID()).
            SetType(model.OrderTypeNewly).
            SetTradeNo(tools.NewUnique().NewSonyflakeID()).
            SetAmount(0).
            SetTotal(0).
            SetInitialDays(sub.InitialDays).
            SetCityID(sub.CityID).
            SetNillablePlanID(sub.PlanID).
            Save(s.ctx)
        if err != nil {
            return
        }

        sub, err = tx.Subscribe.UpdateOneID(sub.ID).SetInitialOrderID(o.ID).Save(s.ctx)
        if err != nil {
            return
        }
        return nil
    })
}

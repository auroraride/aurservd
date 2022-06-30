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
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
    "time"
)

type importRiderService struct {
    ctx context.Context
}

type riderCsvData struct {
    Name  string
    Phone string
    Plan  string
    Days  string
    Model string
    City  string
    Store string
    End   string
}

func NewImportRider() *importRiderService {
    return &importRiderService{
        ctx: context.Background(),
    }
}

func (s *importRiderService) ParseCSV(path string) {
    csvfile, _ := os.Open(path)

    defer func(csvfile *os.File) {
        _ = csvfile.Close()
    }(csvfile)

    r := csv.NewReader(csvfile)

    i := 0
    var items []riderCsvData

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

        arr := strings.Split(strings.TrimSpace(record[7]), "/")
        var end string
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

        items = append(items, riderCsvData{
            Name:  strings.TrimSpace(record[0]),
            Phone: strings.TrimSpace(record[1]),
            Plan:  strings.TrimSpace(record[2]),
            Days:  strings.TrimSpace(record[3]),
            Model: strings.ToUpper(strings.TrimSpace(record[4])),
            City:  strings.TrimSpace(record[5]),
            Store: strings.TrimSpace(record[6]),
            End:   end,
        })
    }

    s.insert(items)
}

func (s *importRiderService) insert(items []riderCsvData) {
    qe := ent.Database.Employee.QueryNotDeleted().Where(employee.Name("曹博文")).FirstX(s.ctx)

    var subs []*ent.Subscribe

    tx, _ := ent.Database.Tx(s.ctx)
    for _, item := range items {
        var p *ent.Person
        var r *ent.Rider
        var sub *ent.Subscribe
        var err error

        if exists, _ := ent.Database.Rider.QueryNotDeleted().Where(rider.Phone(item.Phone)).Exist(s.ctx); exists {
            snag.Panic(fmt.Sprintf("%s已存在", item.Phone))
        }

        // 查找骑行卡
        days, _ := strconv.Atoi(item.Days)
        qp := ent.Database.Plan.QueryNotDeleted().Where(plan.Name(item.Plan), plan.Days(uint(days))).FirstX(s.ctx)

        // 查找门店
        qs := ent.Database.Store.QueryNotDeleted().Where(store.Name(item.Store)).FirstX(s.ctx)

        // 查找城市
        qc := ent.Database.City.QueryNotDeleted().Where(city.Name(item.City)).FirstX(s.ctx)

        // 结束时间
        end := carbon.Parse(item.End)

        // 计算开始时间
        start := end.SubDays(days).Carbon2Time()

        // 创建用户
        p, err = tx.Person.Create().SetName(item.Name).Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)

        // 创建骑手
        r, err = tx.Rider.Create().SetPhone(item.Phone).SetPerson(p).Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)

        // 添加订阅
        sub, err = tx.Subscribe.Create().
            SetEmployeeID(qe.ID).
            SetRider(r).
            SetInitialDays(days).
            SetType(model.OrderTypeNewly).
            SetStatus(model.SubscribeStatusUsing).
            SetStartAt(start).
            SetStore(qs).
            SetPlan(qp).
            SetCity(qc).
            SetModel(item.Model).
            SetRemaining(tools.NewTime().DiffDays(end.Carbon2Time(), time.Now())).
            Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)

        log.Printf("[添加完成] %#v", item)

        subs = append(subs, sub)
    }

    _ = tx.Commit()
}

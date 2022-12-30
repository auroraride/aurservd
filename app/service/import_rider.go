// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/allocate"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/stock"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "github.com/xuri/excelize/v2"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

type importRiderService struct {
    modifier *model.Modifier
    ctx      context.Context
    plan     *ent.Plan
    epoch    time.Time
}

func NewImportRider() *importRiderService {
    return &importRiderService{
        ctx:   context.Background(),
        epoch: time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC),
    }
}

func NewImportRiderWithModifier(m *model.Modifier) *importRiderService {
    s := NewImportRider()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *importRiderService) BatchFile(path string) (err error) {
    var xlsx *excelize.File
    xlsx, err = excelize.OpenFile(path)
    if err != nil {
        return
    }
    defer func() {
        // 关闭工作簿
        _ = xlsx.Close()
    }()

    var styleID int
    styleID, err = xlsx.NewStyle(&excelize.Style{
        Fill: excelize.Fill{
            Type:    "pattern",
            Color:   []string{"#EA3323"},
            Pattern: 1,
        },
        Font: &excelize.Font{Color: "#FFFFFF"},
    })
    if err != nil {
        return
    }

    var rows [][]string

    // 获取 Sheet1 上所有单元格
    rows, err = xlsx.GetRows("Sheet1")
    if err != nil {
        return
    }

    err = xlsx.SetCellStyle("Sheet1", "H2", fmt.Sprintf("H%d", len(rows)), 0)
    if err != nil {
        return
    }

    rows, _ = xlsx.GetRows("Sheet1")

    for i, record := range rows {
        if i == 0 {
            continue
        }
        if record[0] == "" {
            continue
        }

        var item *model.ImportRiderFromExcel
        item, err = s.parseRow(record)

        _ = xlsx.SetCellValue("Sheet1", fmt.Sprintf("H%d", i+1), item.End)

        if err != nil {
            log.Errorf("[%s] 添加失败: %s", item, err)
            // 设置错误
            _ = xlsx.SetCellValue("Sheet1", fmt.Sprintf("I%d", i+1), err.Error())
            _ = xlsx.SetCellStyle("Sheet1", fmt.Sprintf("A%d", i+1), fmt.Sprintf("I%d", i+1), styleID)
            continue
        }
    }

    // 保存结果
    r := filepath.Join("runtime", "import", fmt.Sprintf("%s.xlsx", time.Now().Format(carbon.ShortDateTimeLayout)))
    err = utils.NewFile(r).CreateDirectoryIfNotExist()
    if err != nil {
        return
    }

    return xlsx.SaveAs(r)
}

// parseRow 解析行
func (s *importRiderService) parseRow(record []string) (item *model.ImportRiderFromExcel, err error) {
    x, _ := strconv.Atoi(record[7])
    end := s.epoch.Add(time.Second * time.Duration(x*86400)).Format(carbon.DateLayout)
    bm := strings.ToUpper(strings.TrimSpace(record[4]))
    item = &model.ImportRiderFromExcel{
        Name:  strings.TrimSpace(record[0]),
        Phone: strings.TrimSpace(record[1]),
        Plan:  strings.TrimSpace(record[2]),
        Days:  strings.TrimSpace(record[3]),
        Model: bm,
        City:  strings.TrimSpace(record[5]),
        Store: strings.TrimSpace(record[6]),
        End:   end,
    }

    // 查找城市
    qc := ent.Database.City.QueryNotDeleted().Where(city.Name(item.City)).FirstX(s.ctx)

    // 查找骑行卡
    days, _ := strconv.Atoi(item.Days)
    s.plan, _ = ent.Database.Plan.QueryNotDeleted().Where(
        plan.Name(item.Plan),
        plan.Days(uint(days)),
        plan.HasCitiesWith(city.ID(qc.ID)),
        plan.Model(item.Model),
    ).First(s.ctx)
    if s.plan == nil {
        err = errors.New("未找到骑行卡")
        return
    }

    // 查找门店
    qs, _ := ent.Database.Store.QueryNotDeleted().Where(store.Name(item.Store)).First(s.ctx)
    if qs == nil {
        err = errors.New("未找到门店")
        return
    }

    err = s.Create(&model.ImportRiderCreateReq{
        Name:       item.Name,
        Phone:      item.Phone,
        PlanID:     s.plan.ID,
        CityID:     qc.ID,
        StoreID:    qs.ID,
        EmployeeID: 38654705685,
        End:        end,
    })
    return
}

// Create 手动添加骑手
// TODO 导入车电套餐
func (s *importRiderService) Create(req *model.ImportRiderCreateReq) error {
    return ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
        var (
            p   *ent.Person
            r   *ent.Rider
            o   *ent.Order
            sub *ent.Subscribe
        )

        if r, _ = ent.Database.Rider.QueryNotDeleted().WithPerson().Where(rider.Phone(req.Phone)).First(s.ctx); r != nil {
            if existSub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.RiderID(r.ID)).First(s.ctx); existSub != nil {
                return fmt.Errorf("%s:%s 已存在 <%d>", req.Phone, req.Name, existSub.ID)
            }
        }

        // 结束时间
        end := carbon.Parse(req.End)

        // 查询plan
        if s.plan == nil {
            s.plan, err = ent.Database.Plan.QueryNotDeleted().Where(plan.ID(req.PlanID)).First(s.ctx)
            if err != nil {
                return
            }
        }

        // 定义变量
        var (
            bike     *ent.Ebike
            brand    *ent.EbikeBrand
            bikeID   *uint64
            brandID  *uint64
            alloType = allocate.TypeBattery
        )

        // 查询电车
        if s.plan.BrandID != nil {
            bike, err = NewEbike().AllocatableBaseFilter().
                Where(ebike.Sn(req.EbikeSN), ebike.StoreID(req.StoreID)).
                WithBrand().
                First(s.ctx)
            if bike == nil {
                return errors.New("电车未找到")
            }

            brand = bike.Edges.Brand
            if brand == nil {
                return errors.New("电车型号未找到")
            }

            bikeID = silk.UInt64(bike.ID)
            brandID = silk.UInt64(brand.ID)
            alloType = allocate.TypeEbike
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
            rc := tx.Rider.Create().SetPhone(req.Phone).SetName(req.Name).SetRemark("导入骑手")
            if p != nil {
                rc.SetPersonID(p.ID).SetIDCardNumber(p.IDCardNumber)
            }
            r, err = rc.Save(s.ctx)
            if err != nil {
                return
            }
        } else {
            ru := tx.Rider.UpdateOne(r).SetRemark("导入骑手 & 更新")
            if r.Edges.Person == nil {
                ru.SetPersonID(p.ID).SetIDCardNumber(p.IDCardNumber).SetName(p.Name)
            }
            r, err = ru.Save(s.ctx)
            if err != nil {
                return
            }
        }

        // 添加订阅
        sub, err = tx.Subscribe.Create().
            SetRemark("导入数据").
            SetEmployeeID(req.EmployeeID).
            SetRider(r).
            SetInitialDays(int(s.plan.Days)).
            SetType(model.OrderTypeNewly).
            SetStatus(model.SubscribeStatusUsing).
            SetStartAt(start).
            SetStoreID(req.StoreID).
            SetPlanID(req.PlanID).
            SetCityID(req.CityID).
            SetModel(s.plan.Model).
            SetIntelligent(s.plan.Intelligent).
            SetNeedContract(false).
            SetRemaining(tools.NewTime().LastDays(end.Carbon2Time(), time.Now())).
            SetNillableBrandID(brandID).
            SetNillableEbikeID(bikeID).
            Save(s.ctx)
        if err != nil {
            return
        }

        // 添加分配信息
        _, err = tx.Allocate.Create().
            SetType(alloType).
            SetEmployeeID(req.EmployeeID).
            SetStoreID(req.StoreID).
            SetNillableEbikeID(bikeID).
            SetNillableBrandID(brandID).
            SetSubscribe(sub).
            SetRider(r).
            SetStatus(model.AllocateStatusSigned.Value()).
            SetTime(time.Now()).
            SetModel(sub.Model).
            Save(s.ctx)
        if err != nil {
            return
        }

        // 创建订单
        o, err = tx.Order.Create().
            SetRemark("导入数据").
            SetRiderID(r.ID).
            SetSubscribeID(sub.ID).
            SetStatus(model.OrderStatusPaid).
            SetPayway(model.OrderPaywayManual).
            SetOutTradeNo(tools.NewUnique().NewSonyflakeID()).
            SetType(model.OrderTypeNewly).
            SetTradeNo(tools.NewUnique().NewSonyflakeID()).
            SetAmount(0).
            SetTotal(0).
            SetInitialDays(sub.InitialDays).
            SetCityID(sub.CityID).
            SetNillablePlanID(sub.PlanID).
            SetNillableBrandID(brandID).
            SetNillableEbikeID(bikeID).
            Save(s.ctx)
        if err != nil {
            return
        }

        sub, err = tx.Subscribe.UpdateOneID(sub.ID).SetInitialOrderID(o.ID).Save(s.ctx)
        if err != nil {
            return
        }

        // 创建 stock
        var stockParent *ent.Stock
        stockParent, err = tx.Stock.Create().
            SetRemark("导入数据").
            SetStoreID(req.StoreID).
            SetEmployeeID(req.EmployeeID).
            SetName(s.plan.Model).
            SetRiderID(sub.RiderID).
            SetType(model.StockTypeRiderActive).
            SetModel(s.plan.Model).
            SetNum(-1).
            SetCityID(req.CityID).
            SetSubscribeID(sub.ID).
            SetMaterial(stock.MaterialBattery).
            SetSn(tools.NewUnique().NewSN()).
            Save(s.ctx)
        if err != nil {
            return
        }

        // 更新电车
        if bike != nil {
            err = tx.Ebike.UpdateOneID(bike.ID).SetRiderID(r.ID).SetStatus(model.EbikeStatusUsing).Exec(s.ctx)
            if err != nil {
                return
            }
            err = tx.Stock.Create().
                SetRemark("导入数据").
                SetStoreID(req.StoreID).
                SetEmployeeID(req.EmployeeID).
                SetName(brand.Name).
                SetRiderID(sub.RiderID).
                SetType(model.StockTypeRiderActive).
                SetNum(-1).
                SetCityID(req.CityID).
                SetSubscribeID(sub.ID).
                SetMaterial(stock.MaterialEbike).
                SetSn(tools.NewUnique().NewSN()).
                SetNillableEbikeID(bikeID).
                SetNillableBrandID(brandID).
                SetParent(stockParent).
                Exec(s.ctx)
            if err != nil {
                return
            }
        }

        // 创建 business
        _, err = tx.Business.Create().
            SetRemark("导入数据").
            SetStoreID(req.StoreID).
            SetEmployeeID(req.EmployeeID).
            SetRiderID(sub.RiderID).
            SetSubscribeID(sub.ID).
            SetCityID(sub.CityID).
            SetNillableEnterpriseID(sub.EnterpriseID).
            SetNillableStationID(sub.StationID).
            SetNillablePlanID(sub.PlanID).
            SetType(business.TypeActive).
            SetStockSn(stockParent.Sn).
            Save(s.ctx)

        return
    })
}

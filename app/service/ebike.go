// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-01
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/allocate"
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/internal/ent/ebikebrand"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/labstack/echo/v4"
    "strings"
)

type ebikeService struct {
    *BaseService
    orm *ent.EbikeClient
}

func NewEbike(params ...any) *ebikeService {
    return &ebikeService{
        BaseService: newService(params...),
        orm:         ent.Database.Ebike,
    }
}

func (s *ebikeService) Query(id uint64) (*ent.Ebike, error) {
    return s.orm.Query().Where(ebike.ID(id)).First(s.ctx)
}

func (s *ebikeService) QueryX(id uint64) *ent.Ebike {
    e, _ := s.Query(id)
    if e == nil {
        snag.Panic("未找到电车")
    }
    return e
}

func (s *ebikeService) QueryKeyword(keyword string) (*ent.Ebike, error) {
    return s.orm.Query().Where(
        ebike.Or(
            ebike.Sn(keyword),
            ebike.Plate(keyword),
        ),
    ).First(s.ctx)
}

func (s *ebikeService) QueryKeywordX(keyword string) *ent.Ebike {
    result, _ := s.QueryKeyword(keyword)
    if result == nil {
        snag.Panic("未找到电车")
    }
    return result
}

// AllocatableBaseFilter 可分配车辆查询条件 (不包含门店筛选)
func (s *ebikeService) AllocatableBaseFilter() *ent.EbikeQuery {
    return s.orm.Query().Where(
        ebike.Enable(true),
        ebike.Status(model.EbikeStatusInStock),
        ebike.PlateNotNil(),
        ebike.MachineNotNil(),
        ebike.SimNotNil(),
        ebike.RiderIDIsNil(),
    )
}

// IsAllocated 电车是否已分配
func (s *ebikeService) IsAllocated(id uint64) bool {
    exists, _ := ent.Database.Allocate.Query().Where(
        allocate.EbikeID(id),
        allocate.Status(model.AllocateStatusPending.Value()),
        allocate.TimeGTE(carbon.Now().SubSeconds(model.AllocateExpiration).Carbon2Time()),
    ).Exist(s.ctx)
    return exists
}

func (s *ebikeService) IsAllocatedX(id uint64) {
    if s.IsAllocated(id) {
        snag.Panic("电车已被分配")
    }
}

func (s *ebikeService) QueryAllocatable(id, storeID uint64) (bike *ent.Ebike) {
    if s.IsAllocated(id) {
        return
    }
    q := s.AllocatableBaseFilter().WithBrand().Where(ebike.StoreIDNotNil(), ebike.ID(id))
    if storeID > 0 {
        q.Where(ebike.StoreID(storeID))
    }
    bike, _ = q.First(s.ctx)
    return bike
}

func (s *ebikeService) QueryAllocatableX(id, storeID uint64) *ent.Ebike {
    bike := s.QueryAllocatable(id, storeID)
    if bike == nil {
        snag.Panic("未找到可分配电车")
    }
    return bike
}

func (s *ebikeService) listFilter(req model.EbikeListFilter) (q *ent.EbikeQuery, info ar.Map) {
    info = make(ar.Map)

    q = s.orm.Query().Order(ent.Desc(ebike.FieldCreatedAt)).WithRider().WithStore().WithBrand()

    // 启用状态
    if req.Enable != nil {
        q.Where(ebike.Enable(*req.Enable))
        if *req.Enable {
            info["启用"] = "是"
        } else {
            info["启用"] = "否"
        }
    }

    // 状态
    if req.Status != nil {
        info["状态"] = req.Status.String()
        q.Where(ebike.Status(*req.Status))
    }

    // 骑手
    if req.RiderID != 0 {
        info["骑手"] = ent.NewExportInfo(req.RiderID, rider.Table)
        q.Where(ebike.RiderID(req.RiderID))
    }

    // 门店
    if req.StoreID != 0 {
        info["门店"] = ent.NewExportInfo(req.StoreID, store.Table)
        q.Where(ebike.StoreID(req.StoreID))
    }

    // 品牌
    if req.BrandID != 0 {
        info["品牌"] = ent.NewExportInfo(req.BrandID, ebikebrand.Table)
        q.Where(ebike.BrandID(req.BrandID))
    }

    // 关键词
    if req.Keyword != "" {
        info["关键词"] = req.Keyword
        q.Where(ebike.Or(
            ebike.Sn(req.Keyword),
            ebike.Plate(req.Keyword),
            ebike.Machine(req.Keyword),
            ebike.Sim(req.Keyword),
            ebike.HasRiderWith(rider.Or(
                rider.Name(req.Keyword),
                rider.Phone(req.Keyword),
            )),
        ))
    }

    // 生产批次
    if req.ExFactory != "" {
        info["生产批次"] = req.ExFactory
        q.Where(ebike.ExFactory(req.ExFactory))
    }

    return
}

func (s *ebikeService) List(req *model.EbikeListReq) *model.PaginationRes {
    q, _ := s.listFilter(req.EbikeListFilter)
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Ebike) model.EbikeListRes {
        eb := item.Edges.Brand
        er := item.Edges.Rider
        es := item.Edges.Store
        res := model.EbikeListRes{
            ID:        item.ID,
            SN:        item.Sn,
            BrandID:   item.BrandID,
            ExFactory: item.ExFactory,
            Status:    item.Status.String(),
            EbikeAttributes: model.EbikeAttributes{
                Enable:  silk.Pointer(item.Enable),
                Plate:   item.Plate,
                Machine: item.Machine,
                Sim:     item.Sim,
                Color:   silk.Pointer(item.Color),
            },
        }
        if eb != nil {
            res.Brand = eb.Name
        }
        if er != nil {
            res.Rider = fmt.Sprintf("%s-%s", er.Name, er.Phone)
        }
        if es != nil {
            res.Store = es.Name
        }
        return res
    })
}

func (s *ebikeService) Create(req *model.EbikeCreateReq) {
    b := NewEbikeBrand().QueryX(req.BrandID)
    s.orm.Create().
        SetBrandID(b.ID).
        SetExFactory(req.ExFactory).
        SetSn(req.SN).
        SetNillableMachine(req.Machine).
        SetNillablePlate(req.Plate).
        SetNillableSim(req.Sim).
        SetNillableEnable(req.Enable).
        SetNillableColor(req.Color).
        ExecX(s.ctx)
}

func (s *ebikeService) Modify(req *model.EbikeModifyReq) {
    updater := s.QueryX(req.ID).Update()

    if req.ExFactory != nil {
        updater.SetExFactory(*req.ExFactory)
    }

    updater.
        SetNillableMachine(req.Machine).
        SetNillablePlate(req.Plate).
        SetNillableSim(req.Sim).
        SetNillableEnable(req.Enable).
        SetNillableColor(req.Color).
        ExecX(s.ctx)
}

func (s *ebikeService) BatchCreate(c echo.Context) (failed []string) {
    xlsx := s.BaseService.GetXlsxDataX(c)
    if len(xlsx) < 2 {
        snag.Panic("至少有一条车辆信息")
    }
    // 获取所有型号
    brands := NewEbikeBrand().All()
    bm := make(map[string]uint64)
    for _, brand := range brands {
        bm[brand.Name] = brand.ID
    }

    s.orm.CreateBulk()

    failed = make([]string, 0)

    // 轮询一遍去重
    var sns []string
    var rows [][]string
    for _, columns := range xlsx[1:] {
        if len(columns) < 3 {
            failed = append(failed, fmt.Sprintf("格式错误:%s", strings.Join(columns, ",")))
            continue
        }
        sns = append(sns, columns[1])
        rows = append(rows, columns)
    }

    arr, _ := s.orm.Query().Where(ebike.SnIn(sns...)).All(s.ctx)
    exists := make(map[string]bool)
    for _, a := range arr {
        exists[a.Sn] = true
    }

    // 型号:brand(需查询) 车架号:sn 生产批次:exFactory 车牌号:plate 终端编号:machine SIM卡:sim 颜色:color
    for _, columns := range rows {

        bid, ok := bm[columns[0]]
        if !ok {
            failed = append(failed, fmt.Sprintf("型号未找到:%s", strings.Join(columns, ",")))
            continue
        }

        if _, ok = exists[columns[1]]; ok {
            failed = append(failed, fmt.Sprintf("车架号重复:%s", strings.Join(columns, ",")))
            continue
        }

        creator := s.orm.Create().SetBrandID(bid).SetSn(columns[1]).SetExFactory(columns[2]).SetRemark("批量导入")
        if len(columns) > 3 {
            creator.SetPlate(columns[3])
        }
        if len(columns) > 4 {
            creator.SetMachine(columns[4])
        }
        if len(columns) > 5 {
            creator.SetSim(columns[5])
        }
        color := model.EbikeColorDefault
        if len(columns) > 6 {
            color = strings.ReplaceAll(columns[6], "色", "")
        }
        creator.SetColor(color)

        err := creator.Exec(s.ctx)
        if err != nil {
            msg := "保存失败"
            if strings.Contains(err.Error(), "duplicate key value") {
                msg = "有重复项"
            }
            failed = append(failed, fmt.Sprintf("%s:%s", msg, strings.Join(columns, ",")))
        }
    }

    return
}

func (s *ebikeService) Detail(bike *ent.Ebike, brand *ent.EbikeBrand) *model.Ebike {
    if bike == nil && brand == nil {
        return nil
    }
    res := &model.Ebike{}
    if bike != nil {
        res.EbikeInfo = model.EbikeInfo{
            ID:        bike.ID,
            SN:        bike.Sn,
            ExFactory: bike.ExFactory,
            Plate:     bike.Plate,
            Color:     bike.Color,
        }
    }
    if brand != nil {
        res.Brand = &model.EbikeBrand{
            ID:   brand.ID,
            Name: brand.Name,
        }
    }
    return res
}

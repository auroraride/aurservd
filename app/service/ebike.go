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
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/internal/ent/ebikebrand"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
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
            ID: item.ID,
            EbikeAttributes: model.EbikeAttributes{
                Enable:  silk.Pointer(item.Enable),
                Plate:   item.Plate,
                Machine: item.Machine,
                Sim:     item.Sim,
                Color:   silk.Pointer(item.Color),
            },
            ExFactory: item.ExFactory,
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

    if req.SN != nil {
        updater.SetSn(*req.SN)
    }

    if req.BrandID != nil {
        updater.SetBrand(NewEbikeBrand().QueryX(*req.BrandID))
    }

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
    rows := s.BaseService.GetXlsxDataX(c)
    if len(rows) < 2 {
        snag.Panic("至少有一条车辆信息")
    }
    // 获取所有型号
    brands := NewEbikeBrand().All()
    bm := make(map[string]uint64)
    for _, brand := range brands {
        bm[brand.Name] = brand.ID
    }

    s.orm.CreateBulk()
    var bulk []*ent.EbikeCreate

    // 型号:brand(需查询) 车架号:sn 生产批次:exFactory 车牌号:plate 终端编号:machine SIM卡:sim 颜色:color
    for _, columns := range rows {
        if len(columns) != 7 {
            failed = append(failed, strings.Join(columns, ","))
        }

        bid, ok := bm[columns[0]]
        if !ok {
            failed = append(failed, strings.Join(columns, ","))
        }

        creater := s.orm.Create().SetBrandID(bid).SetSn(columns[1]).SetExFactory(columns[2]).SetRemark("批量导入")
        if columns[3] != "" {
            creater.SetPlate(columns[3])
        }
        if columns[4] != "" {
            creater.SetMachine(columns[4])
        }
        if columns[5] != "" {
            creater.SetSim(columns[5])
        }
        color := model.EbikeColorDefault
        if columns[6] != "" {
            color = columns[6]
        }
        creater.SetColor(color)

        bulk = append(bulk, creater)
    }

    if len(failed) > 0 {
        return
    }

    return nil
}

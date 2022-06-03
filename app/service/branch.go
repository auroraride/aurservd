// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "entgo.io/ent/dialect/sql"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/branchcontract"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/jinzhu/copier"
    "github.com/lithammer/shortuuid/v4"
    "time"
)

type branchService struct {
    orm *ent.BranchClient
    ctx context.Context
}

func NewBranch() *branchService {
    return &branchService{
        orm: ar.Ent.Branch,
        ctx: context.Background(),
    }
}

func NewBranchWithModifier(m *model.Modifier) *branchService {
    s := NewBranch()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    return s
}

// Query 根据ID查询网点
func (s *branchService) Query(id uint64) *ent.Branch {
    item, err := s.orm.QueryNotDeleted().Where(branch.ID(id)).Only(s.ctx)
    if err != nil {
        snag.Panic("未找到有效网点")
    }
    return item
}

// Create 新增网点
// TODO 从结构体新增
func (s *branchService) Create(req *model.BranchCreateReq) {
    tx, _ := ar.Ent.Tx(s.ctx)

    // TODO: 校验城市是否启用
    b, err := s.orm.Create().
        SetName(*req.Name).
        SetAddress(*req.Address).
        SetCityID(*req.CityID).
        SetLng(*req.Lng).
        SetLat(*req.Lat).
        SetGeom(&model.Geometry{
            Lng: *req.Lng,
            Lat: *req.Lat,
        }).
        SetPhotos(req.Photos).
        Save(s.ctx)
    if err != nil {
        _ = tx.Rollback()
        snag.Panic(err)
    }

    if len(req.Contracts) > 0 {
        for _, contract := range req.Contracts {
            s.AddContract(b.ID, contract)
        }
    }

    _ = tx.Commit()
}

// AddContract 新增合同
// TODO 从结构体新增
func (s *branchService) AddContract(id uint64, req *model.BranchContract) *ent.BranchContract {
    return ar.Ent.BranchContract.Create().
        SetBranchID(id).
        SetLandlordName(req.LandlordName).
        SetIDCardNumber(req.IDCardNumber).
        SetPhone(req.Phone).
        SetBankNumber(req.BankNumber).
        SetPledge(req.Pledge).
        SetRent(req.Rent).
        SetLease(req.Lease).
        SetElectricityPledge(req.ElectricityPledge).
        SetElectricity(req.Electricity).
        SetArea(req.Area).
        SetStartTime(req.StartTime).
        SetEndTime(req.EndTime).
        SetFile(req.File).
        SetSheets(req.Sheets).
        SaveX(s.ctx)
}

// List 网点列表
func (s *branchService) List(req *model.BranchListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        Order(ent.Desc(branch.FieldID))

    if req.CityID != nil {
        q.Where(branch.CityID(*req.CityID))
    }

    q.WithCity().
        WithStores().
        WithCabinets().
        WithContracts(func(query *ent.BranchContractQuery) {
            query.Order(ent.Desc(branchcontract.FieldID))
        })

    return model.ParsePaginationResponse[model.BranchItem, ent.Branch](
        q,
        req.PaginationReq,
        func(item *ent.Branch) model.BranchItem {
            var r model.BranchItem
            _ = copier.Copy(&r, item)

            cs := make([]model.BranchContract, len(item.Edges.Contracts))
            for n, contract := range item.Edges.Contracts {
                var c model.BranchContract
                _ = copier.Copy(&c, contract)
                cs[n] = c
            }

            r.Contracts = cs
            city := item.Edges.City
            r.City = model.City{
                ID:   city.ID,
                Name: city.Name,
            }
            r.StoreTotal = len(item.Edges.Stores)
            for _, c := range item.Edges.Cabinets {
                if c.Models[0].Voltage == 60 {
                    r.V60Total += 1
                } else {
                    r.V72Total += 1
                }
            }
            return r
        })
}

// Modify 修改网点
func (s *branchService) Modify(req *model.BranchModifyReq) {
    b := s.orm.QueryNotDeleted().Where(branch.ID(req.ID)).OnlyX(s.ctx)
    if b == nil {
        snag.Panic("网点不存在")
    }

    // 从结构体更新
    q := s.orm.ModifyOne(b, req)

    geom := b.Geom
    if req.Lng != nil {
        geom.Lng = *req.Lng
    }
    if req.Lat != nil {
        geom.Lat = *req.Lat
    }

    q.SetGeom(geom).SaveX(s.ctx)
}

// Selector 网点选择列表
func (s *branchService) Selector() *model.ItemListRes {
    items := make([]model.BranchSampleItem, 0)
    s.orm.QueryNotDeleted().Select(branch.FieldID, branch.FieldName).ScanX(s.ctx, &items)
    res := new(model.ItemListRes)
    model.SetItemListResItems[model.BranchSampleItem](res, items)
    return res
}

// ListByDistance 根据距离列出所有网点和电柜
func (s *branchService) ListByDistance(req *model.BranchWithDistanceReq) (items []*model.BranchWithDistanceRes) {
    var temps []struct {
        ID       uint64  `json:"id"`
        Distance float64 `json:"distance"`
        Name     string  `json:"name"`
        Lng      float64 `json:"lng"`
        Lat      float64 `json:"lat"`
        Image    string  `json:"image"`
        Address  string  `json:"address"`
    }
    if req.Distance == nil && req.CityID == nil {
        snag.Panic("距离和城市不能同时为空")
    }
    // rows, err := s.orm.QueryContext(s.ctx, fmt.Sprintf(`SELECT id, name, ST_Distance(%s, ST_GeogFromText('POINT(108.949969 34.333489)')) AS distance FROM %s WHERE ST_DWithin(%s, ST_GeogFromText('POINT(108.949969 34.333489)'), 10000000) ORDER BY distance;`, branch.Table, branch.FieldGeom, branch.FieldGeom))
    q := s.orm.QueryNotDeleted().
        WithCabinets(func(cq *ent.CabinetQuery) {
            cq.WithBms()
        }).
        WithStores().
        Modify(func(sel *sql.Selector) {
            // sq := sql.Select("*").From()
            bt := sql.Table(branch.Table)
            sel.Select(bt.C(branch.FieldID), bt.C(branch.FieldName), bt.C(branch.FieldAddress), bt.C(branch.FieldLat), bt.C(branch.FieldLng)).
                AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(%s, ST_GeogFromText('POINT(%f %f)'))`, branch.FieldGeom, *req.Lng, *req.Lat)), "distance").
                AppendSelectExprAs(sql.Raw(fmt.Sprintf(`TRIM('"' FROM %s[0]::TEXT)`, branch.FieldPhotos)), "image").
                GroupBy(bt.C(branch.FieldID)).
                OrderBy(sql.Asc("distance"))
            if req.Distance != nil {
                if *req.Distance > 100000 {
                    snag.Panic("请求距离太远")
                }
                sel.Where(sql.P(func(b *sql.Builder) {
                    b.WriteString(fmt.Sprintf(`ST_DWithin(%s, ST_GeogFromText('POINT(%f %f)'), %f)`, branch.FieldGeom, *req.Lng, *req.Lat, *req.Distance))
                }))
            }
        })
    if req.CityID != nil {
        q.Where(branch.CityID(*req.CityID))
    }
    err := q.Scan(s.ctx, &temps)
    items = make([]*model.BranchWithDistanceRes, 0)
    itemsMap := make(map[uint64]*model.BranchWithDistanceRes, len(temps))
    if err != nil || len(temps) == 0 {
        return
    }
    ids := make([]uint64, len(temps))
    for i, temp := range temps {
        ids[i] = temp.ID
        itemsMap[temp.ID] = &model.BranchWithDistanceRes{
            ID:       temp.ID,
            Distance: temp.Distance,
            Name:     temp.Name,
            Lng:      temp.Lng,
            Lat:      temp.Lat,
            Image:    temp.Image,
            Address:  temp.Address,
            Facility: make([]model.BranchFacility, 0),
        }
    }

    // 进行关联查询
    // 门店
    stores := ar.Ent.Store.QueryNotDeleted().Where(store.BranchIDIn(ids...)).AllX(s.ctx)
    for _, es := range stores {
        if es.Status == model.StoreStatusNormal {
            itemsMap[es.BranchID].Facility = append(itemsMap[es.BranchID].Facility, model.BranchFacility{
                ID:    es.ID,
                Type:  model.BranchFacilityTypeStore,
                Name:  es.Name,
                State: model.BranchFacilityStateOnline,
                Num:   0,
                Fid:   shortuuid.New(),
            })
        }
    }

    // 电柜
    cabinets := ar.Ent.Cabinet.QueryNotDeleted().Where(cabinet.BranchIDIn(ids...)).AllX(s.ctx)
    for _, c := range cabinets {
        if c.Status == model.CabinetStatusNormal {
            fa := model.BranchFacility{
                ID:    c.ID,
                Name:  c.Name,
                State: model.BranchFacilityStateOffline,
                Type:  model.BranchFacilityTypeV72,
                Fid:   shortuuid.New(),
            }
            // 获取健康状态
            // TODO 状态更新多久算离线 现在是5分钟
            if c.Health == model.CabinetHealthStatusOnline && time.Now().Sub(c.UpdatedAt).Minutes() < 5 {
                fa.State = model.BranchFacilityStateOnline
            }
            // 计算可用电池数量
            for _, bin := range c.Bin {
                fa.Total += 1
                if bin.Electricity.IsBatteryFull() {
                    fa.Num += 1
                }
            }
            // 判定电池型号
            // TODO 如果有多个电压怎么办
            if c.Models[0].Voltage == 60 {
                fa.Type = model.BranchFacilityTypeV60
            }
            itemsMap[c.BranchID].Facility = append(itemsMap[c.BranchID].Facility, fa)
        }
    }

    for _, m := range itemsMap {
        items = append(items, m)
    }
    return
}

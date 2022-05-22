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
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/jinzhu/copier"
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
func (s *branchService) Create(req *model.Branch, mod *model.Modifier) {
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
        SetLastModifier(mod).
        SetCreator(mod).
        Save(s.ctx)
    if err != nil {
        _ = tx.Rollback()
        snag.Panic(err)
    }

    if len(req.Contracts) > 0 {
        for _, contract := range req.Contracts {
            s.AddContract(b.ID, contract, mod)
        }
    }

    _ = tx.Commit()
}

// AddContract 新增合同
// TODO 从结构体新增
func (s *branchService) AddContract(id uint64, req *model.BranchContract, mod *model.Modifier) *ent.BranchContract {
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
        SetLastModifier(mod).
        SetCreator(mod).
        SaveX(s.ctx)
}

// List 网点列表
func (s *branchService) List(req *model.BranchListReq) (res model.PaginationRes) {
    q := s.orm.QueryNotDeleted().
        Order(ent.Desc(branch.FieldID))

    if req.CityID != nil {
        q.Where(branch.CityID(*req.CityID))
    }

    res.Pagination = q.PaginationResult(req.PaginationReq)

    items := q.
        WithContracts(func(query *ent.BranchContractQuery) {
            query.Order(ent.Desc(branchcontract.FieldID))
        }).
        Pagination(req.PaginationReq).
        AllX(s.ctx)

    rs := make([]*model.Branch, len(items))

    for m, item := range items {
        r := new(model.Branch)
        if err := copier.Copy(r, item); err != nil {
            snag.Panic(err)
        }

        cs := make([]*model.BranchContract, len(item.Edges.Contracts))
        for n, contract := range item.Edges.Contracts {
            c := new(model.BranchContract)
            if err := copier.Copy(c, contract); err != nil {
                snag.Panic(err)
            }
            cs[n] = c
        }

        r.Contracts = cs

        rs[m] = r
    }

    res.Items = rs
    return
}

// Modify 修改网点
// TODO 从结构体更新
func (s *branchService) Modify(req *model.Branch, mod *model.Modifier) {
    b := s.orm.QueryNotDeleted().Where(branch.ID(req.ID)).OnlyX(s.ctx)
    if b == nil {
        snag.Panic("网点不存在")
    }

    s.orm.UpdateOne(b).
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
        SetLastModifier(mod).
        SaveX(s.ctx)
}

// Selector 网点选择列表
func (s *branchService) Selector() *model.ItemListRes {
    items := make([]model.BranchSampleItem, 0)
    s.orm.QueryNotDeleted().Select(branch.FieldID, branch.FieldName).ScanX(s.ctx, &items)
    res := new(model.ItemListRes)
    model.SetItemListResItems[model.BranchSampleItem](res, items)
    return res
}

func (s *branchService) ListByDistance(req *model.BranchWithDistanceReq) (res []model.BranchWithDistanceRes) {
    // rows, err := s.orm.QueryContext(s.ctx, fmt.Sprintf(`SELECT id, name, ST_Distance(%s, ST_GeogFromText('POINT(108.949969 34.333489)')) AS distance FROM %s WHERE ST_DWithin(%s, ST_GeogFromText('POINT(108.949969 34.333489)'), 10000000) ORDER BY distance;`, branch.Table, branch.FieldGeom, branch.FieldGeom))
    var items []model.BranchWithDistanceRes
    q := s.orm.QueryNotDeleted()
    if req.Type > 0 {
        switch req.Type {
        case model.BranchBranchFacilityTypeStore:
            // q.Where(branch.HasCabinetsWith(cabinet.HasBmsWith(batterymodel.Voltage())))
        }
    }
    _ = q.
        WithCabinets().
        WithStores().
        Modify(func(sel *sql.Selector) {
            sel.Select(branch.FieldID, branch.FieldName).
                AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(%s, ST_GeogFromText('POINT(%f %f)'))`, branch.FieldGeom, *req.Lng, *req.Lat)), "distance").
                Where(sql.P(func(b *sql.Builder) {
                    b.WriteString(fmt.Sprintf(`ST_DWithin(%s, ST_GeogFromText('POINT(%f %f)'), 1000000000)`, branch.FieldGeom, *req.Lng, *req.Lat))
                })).
                OrderBy(sql.Asc("distance"))
        }).
        Scan(s.ctx, &items)
    res = make([]model.BranchWithDistanceRes, len(items))
    // for i, item := range items {
    //     res[i]
    // }
    return
}

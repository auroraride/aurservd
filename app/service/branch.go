// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
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

// Add 新增网点
func (s *branchService) Add(req *model.Branch) {
    tx, _ := ar.Ent.Tx(s.ctx)

    // TODO: 校验城市是否启用
    b, err := s.orm.Create().
        SetName(req.Name).
        SetAddress(req.Address).
        SetCityID(req.CityID).
        SetLng(req.Lng).
        SetLat(req.Lat).
        SetPhotos(req.Photos).
        Save(s.ctx)
    if err != nil {
        _ = tx.Rollback()
        snag.Panic(err)
    }

    if req.Contract != nil {
        s.AddContract(b.ID, req.Contract)
    }

    _ = tx.Commit()
}

// AddContract 新增合同
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
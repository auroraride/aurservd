// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/jinzhu/copier"
    "github.com/lithammer/shortuuid/v4"
    "time"
)

type storeService struct {
    ctx context.Context
    orm *ent.StoreClient
}

func NewStore() *storeService {
    return &storeService{
        ctx: context.Background(),
        orm: ar.Ent.Store,
    }
}

func (s *storeService) Query(id uint64) *ent.Store {
    item, err := s.orm.QueryNotDeleted().Where(store.ID(id)).Only(s.ctx)
    if err != nil {
        snag.Panic("未找到有效门店")
    }
    return item
}

// Create 创建门店
func (s *storeService) Create(m *model.Modifier, req *model.StoreCreateReq) model.StoreItem {
    b := NewBranch().Query(*req.BranchID)
    item := s.orm.Create().
        SetLastModifier(m).
        SetName(*req.Name).
        SetStatus(req.Status).
        SetBranch(b).
        SetSn(shortuuid.New()).
        SaveX(s.ctx)
    return s.Detail(item.ID)
}

// Modify 修改门店
func (s *storeService) Modify(m *model.Modifier, req *model.StoreModifyReq) model.StoreItem {
    item := s.Query(req.ID)
    q := s.orm.UpdateOne(item)
    if req.Status != nil {
        q.SetStatus(*req.Status)
    }
    if req.Name != nil {
        q.SetName(*req.Name)
    }
    if req.BranchID != nil {
        q.SetBranchID(*req.BranchID)
    }
    q.SetLastModifier(m).SaveX(s.ctx)
    return s.Detail(item.ID)
}

// Detail 获取门店详情
// TODO 店员
func (s *storeService) Detail(id uint64) model.StoreItem {
    item, err := s.orm.QueryNotDeleted().
        Where(store.ID(id)).
        WithBranch(func(bq *ent.BranchQuery) {
            bq.WithCity()
        }).
        Only(s.ctx)
    if err != nil {
        snag.Panic("未找到有效门店")
    }
    city := item.Edges.Branch.Edges.City
    return model.StoreItem{
        ID:     item.ID,
        Name:   item.Name,
        Status: item.Status,
        City: model.City{
            ID:   city.ID,
            Name: city.Name,
        },
    }
}

// Delete 删除门店
func (s *storeService) Delete(m *model.Modifier, req *model.IDParamReq) {
    item := s.Query(req.ID)
    s.orm.UpdateOne(item).SetDeletedAt(time.Now()).SetLastModifier(m).SaveX(s.ctx)
}

// List 列举门店
func (s *storeService) List(req *model.StoreListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().WithBranch(func(bq *ent.BranchQuery) {
        bq.WithCity()
    })
    if req.CityID != nil {
        q.Where(store.HasBranchWith(branch.CityID(*req.CityID)))
    }
    if req.Name != nil {
        q.Where(store.NameContainsFold(*req.Name))
    }
    if req.Status != nil {
        q.Where(store.Status(*req.Status))
    }

    return model.ParsePaginationResponse[model.StoreItem](q.PaginationResult(req.PaginationReq), func() []model.StoreItem {
        items := q.Pagination(req.PaginationReq).AllX(s.ctx)
        out := make([]model.StoreItem, len(items))
        for i, item := range items {
            _ = copier.Copy(&out[i], item)
            if item.Edges.Branch != nil {
                city := item.Edges.Branch.Edges.City
                out[i].City = model.City{
                    ID:   city.ID,
                    Name: city.Name,
                }
            }
        }
        return out
    })
}

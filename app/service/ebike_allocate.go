// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/mgo"
    "github.com/auroraride/aurservd/pkg/snag"
    log "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type ebikeAllocateService struct {
    *BaseService
    orm *ent.EbikeClient
}

func NewEbikeAllocate(params ...any) *ebikeAllocateService {
    return &ebikeAllocateService{
        BaseService: newService(params...),
        orm:         ent.Database.Ebike,
    }
}

// UnallocatedInfo 获取未分配车辆信息
func (s *ebikeAllocateService) UnallocatedInfo(keyword string) model.EbikeInfo {
    if s.entStore == nil {
        snag.Panic("店员未上班")
    }

    bike, _ := NewEbike().AllocatableFilter().Where(
        ebike.Or(
            ebike.Sn(keyword),
            ebike.Plate(keyword),
        ),
    ).First(s.ctx)
    if bike == nil {
        snag.Panic("未找到可分配电车")
    }

    // 查询是否已被分配
    NewEbike().IsAllocatedPendingX(bike.ID)

    return model.EbikeInfo{
        ID:        bike.ID,
        SN:        bike.Sn,
        ExFactory: bike.ExFactory,
        Plate:     bike.Plate,
    }
}

// Allocate 分配车辆
func (s *ebikeAllocateService) Allocate(req *model.EbikeAllocateReq) model.EbikeAllocateRes {
    if s.entStore == nil {
        snag.Panic("店员未上班")
    }
    // 查找骑手订阅
    sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
        subscribe.ID(req.RiderID),
        subscribe.Status(model.SubscribeStatusInactive),
    ).WithRider().First(s.ctx)
    if sub == nil {
        snag.Panic("未找到骑手订阅信息")
    }
    r := sub.Edges.Rider
    if r == nil {
        snag.Panic("骑手查询失败")
    }

    // 查找电车
    bike := NewEbike().QueryAllocatableX(req.EbikeID, s.entStore.ID)
    brand := bike.Edges.Brand
    if brand == nil {
        snag.Panic("电车型号查询失败")
    }

    data := model.EbikeAllocate{
        Rider: model.Rider{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  r.Name,
        },
        Ebike: model.Ebike{
            EbikeInfo: model.EbikeInfo{
                ID:        bike.ID,
                SN:        bike.Sn,
                ExFactory: bike.ExFactory,
                Plate:     bike.Plate,
            },
            Brand: model.EbikeBrand{
                ID:   brand.ID,
                Name: brand.Name,
            },
        },
        SubscribeID: sub.ID,
        Status:      model.EbikeAllocateStatussPending,
        Model:       sub.Model,
        EmployeeID:  s.employee.ID,
        StoreID:     s.entStore.ID,
    }
    // 缓存电车分配情况
    result, err := mgo.EbikeAllocate.InsertOne(s.ctx, data)
    if err != nil {
        log.Error(err)
        snag.Panic("电车分配失败")
    }
    return model.EbikeAllocateRes{
        AllocateID: result.InsertedID.(primitive.ObjectID).Hex(),
    }
}

// Info 电车分配信息
func (s *ebikeAllocateService) Info(req *model.EbikeAllocateIDQueryReq) model.EbikeAllocateInfo {
    var ea model.EbikeAllocate
    err := mgo.EbikeAllocate.Find(s.ctx, bson.M{"_id": req.AllocateID}).One(&ea)
    if err != nil {
        snag.Panic("电车分配信息查询失败")
    }
    return model.EbikeAllocateInfo{
        Status: ea.Status,
        Rider:  ea.Rider,
        Ebike:  ea.Ebike,
    }
}

// EmployeeList 电车分配店员列表
func (s *ebikeAllocateService) EmployeeList(req *model.EbikeAllocateEmployeeListReq) *model.PaginationRes {
    var items []model.EbikeAllocateInfo
    q := mgo.EbikeAllocate.
        Find(s.ctx, bson.M{"employeeId": s.employee.ID}).
        Skip(int64(req.GetOffset())).
        Limit(int64(req.GetLimit())).
        Sort("-createdAt")
    t, _ := q.Count()
    total := int(t)
    _ = q.All(&items)

    return &model.PaginationRes{
        Pagination: model.Pagination{
            Current: req.Current,
            Pages:   req.GetPages(total),
            Total:   total,
        },
        Items: items,
    }
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-13
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/LucaTheHacker/go-haversine"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/attendance"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "strings"
    "time"
)

type attendanceService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.AttendanceClient
}

func NewAttendance() *attendanceService {
    return &attendanceService{
        ctx: context.Background(),
        orm: ent.Database.Attendance,
    }
}

func NewAttendanceWithRider(r *ent.Rider) *attendanceService {
    s := NewAttendance()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewAttendanceWithModifier(m *model.Modifier) *attendanceService {
    s := NewAttendance()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewAttendanceWithEmployee(e *ent.Employee) *attendanceService {
    s := NewAttendance()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// check 检查是否满足打卡条件并获取需盘点物资清单
func (s *attendanceService) check(req *model.AttendancePrecheck) (*ent.Store, float64, []model.InventoryItemWithNum) {
    st := NewStore().QuerySn(*req.SN)
    b := NewBranch().Query(st.BranchID)
    if st.EmployeeID != nil && req.Duty {
        snag.Panic("当前已有员工上班")
    }
    if !req.Duty && (st.EmployeeID == nil || *st.EmployeeID != s.employee.ID) {
        snag.Panic("当前上班员工非本人")
    }
    // 判断距离
    if b == nil || b.Lat == 0 || b.Lng == 0 {
        snag.Panic("未找到门店地理信息")
    }
    distance := haversine.Distance(haversine.NewCoordinates(*req.Lat, *req.Lng), haversine.NewCoordinates(b.Lat, b.Lng))
    meters := distance.Kilometers() * 1000
    if meters > 1000 {
        snag.Panic("距离过远")
    }
    return st, meters, NewInventory().ListStockInventory(st.ID, model.InventoryListReq{Count: true})
}

func (s *attendanceService) Precheck(req *model.AttendancePrecheck) []model.InventoryItemWithNum {
    _, _, items := s.check(req)
    return items
}

// QueryDuty 获取员工最近的上班打卡信息
func (s *attendanceService) QueryDuty(storeID, employeeID uint64, panic ...bool) *ent.Attendance {
    at, _ := s.orm.QueryNotDeleted().Where(
        attendance.StoreID(storeID),
        attendance.EmployeeID(employeeID),
        attendance.Duty(true),
    ).First(s.ctx)
    if at == nil && (len(panic) == 0 || panic[0]) {
        snag.Panic("未找到员工上班信息")
    }
    return at
}

// dutyDate 获取上下班日期
func (s *attendanceService) dutyDate(duty bool, storeID, employeeID uint64) time.Time {
    now := carbon.Now().StartOfDay().Carbon2Time()
    // 下班需要计算日期
    if !duty {
        at := s.QueryDuty(storeID, employeeID)
        now = at.Date
    }
    return now
}

// Create 打卡
func (s *attendanceService) Create(req *model.AttendanceCreateReq) {
    if !strings.HasPrefix(*req.Photo, "employee/") {
        snag.Panic("照片错误")
    }

    inventory := make([]model.AttendanceInventory, 0)
    st, distance, items := s.check(req.AttendancePrecheck)
    for _, item := range items {
        n, ok := req.Inventory[item.Name]
        if !ok {
            snag.Panic("请提交全部的物资清单")
        }

        inventory = append(inventory, model.AttendanceInventory{
            Name:     item.Name,
            Num:      n,
            StockNum: item.Num,
            Model:    item.Model,
        })
    }

    tx, _ := ent.Database.Tx(s.ctx)
    _, err := tx.Attendance.Create().
        SetEmployee(s.employee).
        SetStore(st).
        SetDate(s.dutyDate(req.Duty, st.ID, s.employee.ID)).
        SetDuty(req.Duty).
        SetInventory(inventory).
        SetPhoto(*req.Photo).
        SetDuty(req.Duty).
        SetAddress(*req.Address).
        SetLng(*req.Lng).
        SetLat(*req.Lat).
        SetDistance(distance).
        Save(s.ctx)

    if err != nil {
        _ = tx.Rollback()
        snag.Panic("考勤打卡失败")
    }

    if req.Duty {
        _, err = tx.Store.UpdateOneID(st.ID).SetEmployee(s.employee).Save(s.ctx)
    } else {
        _, err = tx.Store.UpdateOneID(st.ID).ClearEmployeeID().Save(s.ctx)
    }
    if err != nil {
        _ = tx.Rollback()
        snag.Panic("考勤打卡失败")
    }

    _ = tx.Commit()
}

func (s *attendanceService) List(req *model.AttendanceListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().WithStore(func(sq *ent.StoreQuery) {
        sq.WithCity()
    }).WithEmployee()
    tt := tools.NewTime()
    if req.Start != nil {
        q.Where(attendance.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }
    if req.End != nil {
        q.Where(attendance.CreatedAtLT(carbon.Time2Carbon(tt.ParseDateStringX(*req.Start)).StartOfDay().Tomorrow().Carbon2Time()))
    }
    if req.Keyword != nil {
        q.Where(
            attendance.HasEmployeeWith(
                employee.Or(
                    employee.NameContainsFold(*req.Keyword),
                    employee.PhoneContainsFold(*req.Keyword),
                ),
            ),
        )
    }
    if req.Duty != 0 {
        q.Where(attendance.Duty(req.Duty == 1))
    }
    if req.CityID != nil {
        q.Where(attendance.HasStoreWith(store.CityID(*req.CityID)))
    }
    if req.StoreID != nil {
        q.Where(attendance.StoreID(*req.StoreID))
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Attendance) model.AttendanceListRes {
        es := item.Edges.Store
        esc := es.Edges.City
        ee := item.Edges.Employee
        return model.AttendanceListRes{
            ID: item.ID,
            City: model.City{
                ID:   esc.ID,
                Name: esc.Name,
            },
            Store: model.Store{
                ID:   es.ID,
                Name: es.Name,
            },
            Name:      ee.Name,
            Phone:     ee.Phone,
            Time:      item.CreatedAt.Format(carbon.DateTimeLayout),
            Photo:     *item.Photo,
            Inventory: item.Inventory,
            Duty:      item.Duty,
        }
    })
}

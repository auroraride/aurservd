// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-13
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/LucaTheHacker/go-haversine"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/attendance"
    "github.com/auroraride/aurservd/pkg/snag"
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
        orm: ar.Ent.Attendance,
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
func (s *attendanceService) check(req *model.AttendancePrecheck) (*ent.Store, float64, []model.InventoryItem) {
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
    miles := distance.Miles()
    if miles > 1000 {
        snag.Panic("距离过远")
    }
    return st, miles, NewInventory().ListInventory(model.InventoryListReq{Count: true})
}

func (s *attendanceService) Precheck(req *model.AttendancePrecheck) []model.InventoryItem {
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
    st, distance, items := s.check(req.AttendancePrecheck)
    for _, item := range items {
        if _, ok := req.Inventory[item.Name]; !ok {
            snag.Panic("请提交全部的物资清单")
        }
    }

    tx, _ := ar.Ent.Tx(s.ctx)
    _, err := tx.Attendance.Create().
        SetEmployee(s.employee).
        SetStore(st).
        SetDate(s.dutyDate(req.Duty, st.ID, s.employee.ID)).
        SetDuty(req.Duty).
        SetInventory(req.Inventory).
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

    _, err = tx.Store.UpdateOneID(st.ID).SetEmployee(s.employee).Save(s.ctx)
    if err != nil {
        _ = tx.Rollback()
        snag.Panic("考勤打卡失败")
    }

    _ = tx.Commit()
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/assistance"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/google/uuid"
    "github.com/rs/xid"
    log "github.com/sirupsen/logrus"
    "time"
)

type employeeService struct {
    cacheKeyPrefix string
    ctx            context.Context
    modifier       *model.Modifier
    rider          *ent.Rider
    employee       *ent.Employee
    orm            *ent.EmployeeClient
}

func NewEmployee() *employeeService {
    return &employeeService{
        cacheKeyPrefix: "EMPLOYEE_",
        ctx:            context.Background(),
        orm:            ar.Ent.Employee,
    }
}

func NewEmployeeWithRider(r *ent.Rider) *employeeService {
    s := NewEmployee()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEmployeeWithModifier(m *model.Modifier) *employeeService {
    s := NewEmployee()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEmployeeWithEmployee(e *ent.Employee) *employeeService {
    s := NewEmployee()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *employeeService) Query(id uint64) (*ent.Employee, error) {
    return s.orm.QueryNotDeleted().Where(employee.ID(id)).First(s.ctx)
}

func (s *employeeService) QueryX(id uint64) *ent.Employee {
    item, _ := s.Query(id)
    if item == nil {
        snag.Panic("未找到店员")
    }
    return item
}

func (s *employeeService) GetEmployeeByID(id uint64) (*ent.Employee, error) {
    return s.orm.QueryNotDeleted().Where(employee.ID(id)).WithStore().First(s.ctx)
}

// QueryByPhone 根据phone查找店员
func (s *employeeService) QueryByPhone(phone string) *ent.Employee {
    item, _ := s.orm.QueryNotDeleted().Where(employee.Phone(phone)).First(s.ctx)
    return item
}

// Create 添加店员
func (s *employeeService) Create(req *model.EmployeeCreateReq) *ent.Employee {
    // 判断重复
    em := s.QueryByPhone(req.Phone)
    if em != nil {
        snag.Panic("店员已存在")
    }
    var err error
    em, err = s.orm.Create().
        SetPhone(req.Phone).
        SetName(req.Name).
        SetCityID(req.CityID).
        Save(s.ctx)
    if em == nil {
        log.Error(err)
        snag.Panic("店员添加失败")
    }
    return em
}

// Modify 修改店员
func (s *employeeService) Modify(req *model.EmployeeModifyReq) {
    _, err := s.orm.ModifyOne(s.QueryX(*req.ID), req).Save(s.ctx)
    if err != nil {
        snag.Panic("保存失败")
        log.Error(err)
    }
}

func (s *employeeService) List(req *model.EmployeeListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        WithCity()

    if req.Keyword != nil {
        q.Where(
            employee.Or(
                employee.NameContainsFold(*req.Keyword),
                employee.PhoneContainsFold(*req.Keyword),
            ),
        )
    }

    if req.CityID != nil {
        q.Where(employee.CityID(*req.CityID))
    }

    if req.Status != 0 {
        q.Where(employee.Enable(req.Status == 1))
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Employee) model.EmployeeListRes {
        res := model.EmployeeListRes{
            ID:     item.ID,
            Name:   item.Name,
            Phone:  item.Phone,
            Enable: item.Enable,
        }

        ec := item.Edges.City
        if ec != nil {
            res.City = model.City{
                ID:   ec.ID,
                Name: ec.Name,
            }
        }
        return res
    })
}

// Activity 店员动态
func (s *employeeService) Activity(req *model.EmployeeActivityListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        WithStore().
        WithCity().
        WithExchanges().
        WithCommissions().
        WithAssistances(func(aq *ent.AssistanceQuery) {
            aq.Where(assistance.StatusIn(model.AssistanceStatusSuccess, model.AssistanceStatusUnpaid))
        })

    if req.Keyword != nil {
        q.Where(
            employee.Or(
                employee.NameContainsFold(*req.Keyword),
                employee.PhoneContainsFold(*req.Keyword),
            ),
        )
    }

    if req.StoreID != nil {
        q.Where(employee.HasStoreWith(store.ID(*req.StoreID)))
    }

    if req.CityID != nil {
        q.Where(employee.CityID(*req.CityID))
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Employee) model.EmployeeActivityListRes {
        res := model.EmployeeActivityListRes{
            ID:            item.ID,
            Name:          item.Name,
            Phone:         item.Phone,
            ExchangeTimes: len(item.Edges.Exchanges),
        }

        ec := item.Edges.City
        if ec != nil {
            res.City = model.City{
                ID:   ec.ID,
                Name: ec.Name,
            }
        }

        es := item.Edges.Store
        if es != nil {
            res.Store = &model.Store{
                ID:   es.ID,
                Name: es.Name,
            }
        }

        for _, cm := range item.Edges.Commissions {
            res.Amount = tools.NewDecimal().Sum(cm.Amount, res.Amount)
        }

        for _, ass := range item.Edges.Assistances {
            res.AssistanceTimes += 1
            res.AssistanceMeters += ass.Distance
        }
        return res
    })
}

// Delete 删除骑手
func (s *employeeService) Delete(req *model.EmployeeDeleteReq) {
    item := s.QueryX(req.ID)
    _, err := s.orm.SoftDeleteOne(item).Save(s.ctx)
    if err != nil {
        log.Error(err)
        snag.Panic("店员删除失败")
    }
}

func (s *employeeService) tokenKey(id uint64) string {
    return fmt.Sprintf("%s%d", s.cacheKeyPrefix, id)
}

// ExtendTokenTime 延长登录有效期
func (s *employeeService) ExtendTokenTime(id uint64, token string) {
    ctx := context.Background()
    cache.Set(ctx, s.tokenKey(id), token, 7*24*time.Hour)
    cache.Set(ctx, token, id, 7*24*time.Hour)
}

// Profile 骑手资料获取
func (s *employeeService) Profile(e *ent.Employee) model.EmployeeProfile {
    st := e.Edges.Store
    res := model.EmployeeProfile{
        ID:     e.ID,
        Qrcode: fmt.Sprintf("EMPLOYEE:%s", e.Sn),
        Phone:  e.Phone,
        Name:   e.Name,
    }
    if st != nil {
        res.Onduty = true
        res.Store = &model.StoreWithStatus{
            Status: st.Status,
            Store: model.Store{
                ID:   st.ID,
                Name: st.Name,
            },
        }
    }
    return res
}

// Signin 店员登录
func (s *employeeService) Signin(req *model.EmployeeSignReq) model.EmployeeProfile {
    // 校验短信
    debugPhones := ar.Config.App.Debug.Phone
    if !debugPhones[req.Phone] && !NewSms().VerifyCode(req.SmsId, req.SmsCode) {
        snag.Panic("短信验证码校验失败")
    }
    e, _ := s.orm.QueryNotDeleted().Where(employee.Phone(req.Phone)).WithStore().First(s.ctx)
    if e == nil {
        snag.Panic("未找到用户")
    }

    // 被禁用
    if !e.Enable {
        snag.Panic(snag.StatusForbidden, ar.BannedMessage)
    }

    // 生成token
    token := xid.New().String() + utils.RandTokenString()
    key := s.tokenKey(e.ID)

    // 删除旧的token
    if old := cache.Get(s.ctx, key).Val(); old != "" {
        cache.Del(s.ctx, key)
        cache.Del(s.ctx, old)
    }

    // 生成UUID
    sn := uuid.New()
    e.Update().SetSn(sn).SaveX(s.ctx)

    s.ExtendTokenTime(e.ID, token)

    e.Sn = sn
    res := s.Profile(e)
    res.Token = token

    return res
}

// RefreshQrcode 重新生成二维码
func (s *employeeService) RefreshQrcode() model.EmployeeQrcodeRes {
    e := s.employee.Update().SetSn(uuid.New()).SaveX(s.ctx)
    return model.EmployeeQrcodeRes{Qrcode: fmt.Sprintf("EMPLOYEE:%s", e.Sn)}
}

func (s *employeeService) Enable(req *model.EmployeeEnableReq) {
    e := s.QueryX(req.ID)
    e.Update().SetEnable(req.Enable).SaveX(s.ctx)
}

func (s *employeeService) Signout(e *ent.Employee) {
    ctx := context.Background()
    key := s.tokenKey(e.ID)
    token := cache.Get(ctx, key).Val()
    cache.Del(ctx, key)
    cache.Del(ctx, token)
}

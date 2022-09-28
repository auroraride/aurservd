// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "strconv"
    "strings"
    "sync/atomic"
    "time"
)

var couponTaskBusy = atomic.Bool{}

func init() {
    couponTaskBusy.Store(false)
}

type couponService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
    orm          *ent.CouponClient
    max          uint64
}

func NewCoupon() *couponService {
    return &couponService{
        ctx: context.Background(),
        orm: ent.Database.Coupon,
        max: 1679615,
    }
}

func NewCouponWithRider(r *ent.Rider) *couponService {
    s := NewCoupon()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewCouponWithModifier(m *model.Modifier) *couponService {
    s := NewCoupon()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewCouponWithEmployee(e *ent.Employee) *couponService {
    s := NewCoupon()
    if e != nil {
        s.employee = e
        s.employeeInfo = &model.Employee{
            ID:    e.ID,
            Name:  e.Name,
            Phone: e.Phone,
        }
        s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
    }
    return s
}

func (s *couponService) HexNumber(n uint64) string {
    return fmt.Sprintf("%02s", strings.ToUpper(strconv.FormatUint(n, 36)))
}

func (s *couponService) Generate(req *model.CouponGenerateReq) {
    if req.Amount < 0 {
        snag.Panic("金额至少为0.01")
    }

    var (
        n      = req.Number
        phones []string
        ids    []uint64
    )

    if len(req.Phones) > 0 {
        m := make(map[string]struct{})
        for _, phone := range req.Phones {
            m[phone] = struct{}{}
        }
        // 查询骑手
        riders, _ := ent.Database.Rider.QueryNotDeleted().Where(rider.PhoneIn(req.Phones...)).All(s.ctx)
        if len(riders) == 0 {
            snag.Panic("全部骑手查询失败")
        }
        // 找到的手机号
        for _, r := range riders {
            delete(m, r.Phone)
            ids = append(ids, r.ID)
        }
        // 未找到的手机号
        for k := range m {
            phones = append(phones, k)
        }
        n = len(ids)
    }

    if n == 0 {
        snag.Panic("生成券不能为0")
    }

    if couponTaskBusy.Load() {
        snag.Panic("其他券码任务执行中")
    }

    // 开始执行任务
    couponTaskBusy.Store(true)

    // 批量生成券码
    keys := s.GenerateCDKeys(n)
    // s.orm.CreateBulk()

    bulk := make([]*ent.CouponCreate, len(keys))
    for i, key := range keys {
        bulk[i] = s.orm.Create().SetCode(key)
    }

    couponTaskBusy.Store(false)
}

// GenerateCDKeys 生成一定数量的券码
func (s *couponService) GenerateCDKeys(n int) (keys []string) {
    for i := 0; i < n; i++ {
        str := s.HexNumber(uint64(time.Now().UnixMicro()) + uint64(i))
        sig := utils.DataSum([]byte(str))
        var r [10]rune
        for k, v := range model.CouponShuffle {
            r[k] = rune(str[v])
        }
        keys = append(keys, string(r[:])+sig)
    }
    return
}

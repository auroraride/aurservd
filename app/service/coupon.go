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
    "github.com/auroraride/aurservd/pkg/snag"
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

func (s *couponService) Generate(n int, amount float64) {
    if amount < 0 {
        snag.Panic("金额至少为0.01")
    }
    if couponTaskBusy.Load() {
        snag.Panic("其他券码任务执行中")
    }

    // 开始执行任务
    couponTaskBusy.Store(true)

    s.GenerateCDKeys(n)
}

// GenerateCDKeys 生成一定数量的券码
func (s *couponService) GenerateCDKeys(n int) (keys []string) {
    for i := 0; i < n; i++ {
        str := s.HexNumber(uint64(time.Now().UnixMicro()) + uint64(i))
        sig := s.DataSum([]byte(str))
        var r [10]rune
        for k, v := range model.CouponShuffle {
            r[k] = rune(str[v])
        }
        keys = append(keys, string(r[:])+sig)
    }
    return
}

// DataSum 校验和
func (s *couponService) DataSum(data []byte) string {
    var (
        sum    uint32
        length = len(data)
    )

    for i := 0; i < length; i++ {
        sum += uint32(data[i])
        if sum > 0xff {
            sum = ^sum
            sum += 1
        }
    }
    sum &= 0xff
    return fmt.Sprintf("%X", sum)
}

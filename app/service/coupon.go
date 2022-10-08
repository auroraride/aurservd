// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/coupon"
    "github.com/auroraride/aurservd/internal/ent/couponassembly"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
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

func (s *couponService) Query(id uint64) (*ent.Coupon, error) {
    return s.orm.Query().Where(coupon.ID(id)).First(s.ctx)
}

func (s *couponService) QueryX(id uint64) *ent.Coupon {
    c, _ := s.Query(id)
    if c == nil {
        snag.Panic("未找到优惠券")
    }
    return c
}

func (s *couponService) HexNumber(n uint64) string {
    return fmt.Sprintf("%02s", strings.ToUpper(strconv.FormatUint(n, 36)))
}

func (s *couponService) Generate(req *model.CouponGenerateReq) (phones []string) {
    if req.Amount < 0 {
        snag.Panic("金额至少为0.01")
    }

    ct := NewCouponTemplate().QueryEnableX(req.TemplateID)

    var (
        n   = req.Number
        ids []uint64
    )

    if len(req.Phones) > 0 {
        _, ids, phones = NewRider().QueryPhones(req.Phones)
        n = len(ids)
    }

    target := model.CouponTargetStock
    toRider := len(ids) > 0
    if toRider {
        target = model.CouponTargetRider
    }
    expiresAt := ct.Meta.ExpiresAt(toRider)

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
    if toRider && len(ids) != len(keys) {
        snag.Panic("优惠券生成失败")
    }

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        var as *ent.CouponAssembly
        as, err = tx.CouponAssembly.Create().
            SetNumber(n).
            SetAmount(req.Amount).
            SetRemark(req.Remark).
            SetTarget(target.Value()).
            SetTemplateID(req.TemplateID).
            SetName(ct.Name).
            SetMeta(ct.Meta).
            Save(s.ctx)
        if err != nil {
            return
        }

        bulk := make([]*ent.CouponCreate, len(keys))
        for i, key := range keys {
            bulk[i] = s.orm.Create().
                SetTemplateID(ct.ID).
                SetName(ct.Name).
                SetRule(ct.Meta.Rule.Value()).
                SetMultiple(ct.Meta.Multiple).
                SetAmount(req.Amount).
                SetCode(key).
                SetNillableExpiresAt(expiresAt).
                SetDuration(ct.Meta.CouponDuration).
                SetRemark(req.Remark).
                SetAssembly(as)
            if toRider {
                bulk[i].SetRiderID(ids[i])
            }
        }
        return tx.Coupon.CreateBulk(bulk...).Exec(s.ctx)
    })

    couponTaskBusy.Store(false)

    return
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

func (s *couponService) listFilter(req model.CouponListFilter) (q *ent.CouponQuery, info ar.Map) {
    info = make(ar.Map)
    q = s.orm.Query().Order(ent.Desc(coupon.FieldCreatedAt))
    if req.RiderID != 0 {
        q.Where(coupon.RiderID(req.RiderID))
        info["骑手"] = ent.NewExportInfo(req.RiderID, rider.Table)
    }
    if req.Keyword != "" {
        q.Where(
            coupon.HasRiderWith(rider.Or(
                rider.NameContainsFold(req.Keyword),
                rider.PhoneContainsFold(req.Keyword),
            )),
        )
        info["骑手关键词"] = req.Keyword
    }
    if req.Status != nil {
        info["状态"] = req.Status.String()
        switch *req.Status {
        case model.CouponStatusInStock:
            q.Where(coupon.RiderIDIsNil(), coupon.ExpiresAtGTE(time.Now()), coupon.UsedAtIsNil())
        case model.CouponStatusExpired:
            q.Where(coupon.UsedAtIsNil(), coupon.ExpiresAtLT(time.Now()))
        case model.CouponStatusUsed:
            q.Where(coupon.UsedAtNotNil())
        case model.CouponStatusUnused:
            q.Where(coupon.RiderIDNotNil(), coupon.UsedAtIsNil(), coupon.ExpiresAtGTE(time.Now()))
        }
    }
    if req.Target != nil {
        info["发送对象"] = req.Target.String()
        q.Where(coupon.HasAssemblyWith(couponassembly.Target(req.Target.Value())))
    }
    if req.Code != "" {
        if !req.Code.IsValid() {
            snag.Panic("券码错误")
        }
        q.Where(coupon.Code(req.Code.String()))
    }
    return
}

func (s *couponService) Status(c *ent.Coupon) model.CouponStatus {
    if c.UsedAt != nil {
        return model.CouponStatusUsed
    }

    if !c.ExpiresAt.IsZero() && c.ExpiresAt.Before(time.Now()) {
        return model.CouponStatusExpired
    }

    if c.RiderID == nil {
        return model.CouponStatusInStock
    }

    return model.CouponStatusUnused
}

func (s *couponService) List(req *model.CouponListReq) *model.PaginationRes {
    q, _ := s.listFilter(req.CouponListFilter)
    return model.ParsePaginationResponse(q.WithPlan().WithRider().WithOrder().WithAssembly(), req.PaginationReq, func(item *ent.Coupon) (res model.CouponListRes) {
        // 城市 / 骑士卡
        ea := item.Edges.Assembly
        var cities, plans []string
        for _, mc := range ea.Meta.Cities {
            cities = append(cities, mc.Name)
        }
        for _, mp := range ea.Meta.Plans {
            plans = append(plans, mp.Name)
        }
        // 骑手
        er := item.Edges.Rider
        var rn, rp string
        if er != nil {
            rn = er.Name
            rp = er.Phone
        }

        // 使用
        ep := item.Edges.Plan
        eo := item.Edges.Order
        var usedAt, expiredAt, uo, up string
        if item.UsedAt != nil {
            usedAt = item.UsedAt.Format(carbon.DateTimeLayout)
            if ep != nil {
                up = fmt.Sprintf("%s-%d天", ep.Name, ep.Days)
            }
            if eo != nil {
                uo = eo.TradeNo
            }
        }
        if item.ExpiresAt.IsZero() {
            expiredAt = item.ExpiresAt.Format(carbon.DateTimeLayout)
        } else {
            at := item.Duration.ExpiresAt(item.RiderID != nil)
            if at != nil {
                expiredAt = at.Format(carbon.DateTimeLayout)
            }
        }
        res = model.CouponListRes{
            ID:         item.ID,
            Amount:     item.Amount,
            Name:       item.Name,
            Code:       model.CouponCode(item.Code).Humanity(),
            TemplateID: item.TemplateID,
            AssemblyID: item.AssemblyID,
            Creator:    item.Creator.Name + " - " + item.Creator.Phone,
            Time:       item.CreatedAt.Format(carbon.DateTimeLayout),
            Rider:      rn,
            Phone:      rp,
            Status:     s.Status(item),
            Cities:     cities,
            Plans:      plans,
            UsedAt:     usedAt,
            ExpiredAt:  expiredAt,
            TradeNo:    uo,
            Plan:       up,
        }
        return
    })
}

// Allocate 分配优惠券
func (s *couponService) Allocate(req *model.CouponAllocateReq) {
    c := s.QueryX(req.ID)
    status := s.Status(c)
    if status != model.CouponStatusInStock {
        snag.Panic(fmt.Sprintf("优惠券状态错误: %s", status))
    }
    s.orm.UpdateOne(c).SetRiderID(req.RiderID).ExecX(s.ctx)
}

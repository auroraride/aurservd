// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "entgo.io/ent/dialect/sql"
    "fmt"
    "github.com/LucaTheHacker/go-haversine"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/assistance"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
)

type assistanceService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.AssistanceClient
}

func NewAssistance() *assistanceService {
    return &assistanceService{
        ctx: context.Background(),
        orm: ar.Ent.Assistance,
    }
}

func NewAssistanceWithRider(r *ent.Rider) *assistanceService {
    s := NewAssistance()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewAssistanceWithModifier(m *model.Modifier) *assistanceService {
    s := NewAssistance()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewAssistanceWithEmployee(e *ent.Employee) *assistanceService {
    s := NewAssistance()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *assistanceService) Breakdown() any {
    return NewSetting().GetSetting(model.SettingRescueReason)
}

// Unpaid 是否有未支付的救援订单
func (s *assistanceService) Unpaid(riderID uint64) *ent.Assistance {
    ass, _ := s.orm.QueryNotDeleted().
        Where(assistance.Status(model.AssistanceStatusUnpaid), assistance.RiderID(riderID)).
        First(s.ctx)
    return ass
}

func (s *assistanceService) List(req *model.AssistanceListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().WithRider(func(rq *ent.RiderQuery) {
        rq.WithPerson()
    }).WithCity().WithStore().WithEmployee()
    tt := tools.NewTime()
    if req.Start != "" {
        q.Where(assistance.CreatedAtGTE(tt.ParseDateStringX(req.Start)))
    }
    if req.End != "" {
        q.Where(assistance.CreatedAtLT(carbon.Time2Carbon(tt.ParseDateStringX(req.Start)).StartOfDay().Tomorrow().Carbon2Time()))
    }
    if req.CityID != 0 {
        q.Where(assistance.CityID(req.CityID))
    }
    if req.Keyword != "" {
        q.Where(
            assistance.HasRiderWith(rider.Or(
                rider.HasPersonWith(person.NameContainsFold(req.Keyword)),
                rider.PhoneContainsFold(req.Keyword),
            )),
        )
    }
    if req.Status != nil {
        q.Where(assistance.Status(*req.Status))
    }
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Assistance) model.AssistanceListRes {
        return s.BasicInfo(item)
    })
}

func (s *assistanceService) BasicInfo(item *ent.Assistance) model.AssistanceListRes {
    r := item.Edges.Rider
    p := r.Edges.Person
    c := item.Edges.City
    res := model.AssistanceListRes{
        ID:       item.ID,
        Status:   item.Status,
        Cost:     item.Cost,
        Distance: item.Distance,
        Time:     item.CreatedAt.Format(carbon.DateTimeLayout),
        Rider: model.RiderBasic{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  p.Name,
        },
        City: model.City{
            ID:   c.ID,
            Name: c.Name,
        },
    }

    e := item.Edges.Employee
    st := item.Edges.Store
    if e != nil {
        res.Employee = &model.Employee{
            ID:    e.ID,
            Name:  e.Name,
            Phone: e.Phone,
        }
    }

    if st != nil {
        res.Store = &model.Store{
            ID:   st.ID,
            Name: st.Name,
        }
    }

    return res
}

func (s *assistanceService) Query(id uint64) (*ent.Assistance, error) {
    return s.orm.QueryNotDeleted().
        Where(assistance.ID(id)).
        First(s.ctx)
}

func (s *assistanceService) QueryX(id uint64) *ent.Assistance {
    item, _ := s.Query(id)
    if item == nil {
        snag.Panic("未找到救援信息")
    }
    return item
}

func (s *assistanceService) QueryDetail(id uint64) (*ent.Assistance, error) {
    return s.orm.QueryNotDeleted().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithCity().
        WithStore().
        WithEmployee().
        Where(assistance.ID(id)).
        First(s.ctx)
}

func (s *assistanceService) QueryDetailX(id uint64) *ent.Assistance {
    item, _ := s.QueryDetail(id)
    if item == nil {
        snag.Panic("未找到救援信息")
    }
    return item
}

// Detail 救援详情
func (s *assistanceService) Detail(id uint64) model.AssistanceDetail {
    item := s.QueryDetailX(id)
    res := model.AssistanceDetail{
        AssistanceListRes: s.BasicInfo(item),
        Lng:               item.Lng,
        Lat:               item.Lat,
        Address:           item.Address,
        Breakdown:         item.Breakdown,
        BreakdownDesc:     item.BreakdownDesc,
        BreakdownPhotos:   item.BreakdownPhotos,
        Reason:            item.Reason,
        DetectPhoto:       item.DetectPhoto,
        JointPhoto:        item.JointPhoto,
        RefusedDesc:       item.RefusedDesc,
        FreeReason:        item.FreeReason,
    }
    if item.PayAt != nil {
        res.PayAt = tools.NewPointer().String(item.PayAt.Format(carbon.DateTimeLayout))
    }

    return res
}

// Create 发起救援订单
// 救援订单未支付的禁止办理所有业务
// TODO 救援订单支付状态可以直接在后台修改为不需要支付
func (s *assistanceService) Create(req *model.AssistanceCreateReq) model.AssistanceCreateRes {
    sub := NewSubscribe().Recent(s.rider.ID)
    if sub == nil || sub.Status != model.SubscribeStatusUsing {
        snag.Panic("无法发起救援")
    }

    // 检查是否可发起救援
    NewRiderPermissionWithRider(s.rider).BusinessX()

    // 检查是否已有救援订单
    if exists, _ := s.orm.QueryNotDeleted().Where(assistance.RiderID(s.rider.ID)).Exist(s.ctx); exists {
        snag.Panic("当前有进行中的救援订单")
    }

    as, _ := s.orm.Create().
        SetStatus(model.AssistanceStatusPending).
        SetLng(req.Lng).
        SetLat(req.Lat).
        // SetDistance(haversine.Distance(haversine.NewCoordinates(req.Lat, req.Lng), haversine.NewCoordinates(stb.Lat, stb.Lng)).Miles()).
        SetAddress(req.Address).
        SetBreakdown(req.Breakdown).
        SetBreakdownPhotos(req.BreakdownPhotos).
        SetBreakdownDesc(req.BreakdownDesc).
        SetOutTradeNo(tools.NewUnique().NewSN28()).
        SetRiderID(s.rider.ID).
        SetSubscribeID(sub.ID).
        SetCityID(sub.CityID).
        Save(s.ctx)

    if as == nil {
        snag.Panic("救援发起失败")
    }

    return model.AssistanceCreateRes{OutTradeNo: as.OutTradeNo}
}

// Nearby 救援订单附近门店
func (s *assistanceService) Nearby(req *model.IDQueryReq) (items []model.AssistanceNearbyRes) {
    ass, _ := s.orm.QueryNotDeleted().Where(assistance.ID(req.ID)).First(s.ctx)
    if ass == nil {
        snag.Panic("未找到救援订单")
    }

    var temps []struct {
        ID       uint64  `json:"id"`
        Name     string  `json:"name"`
        Lng      float64 `json:"lng"`
        Lat      float64 `json:"lat"`
        Distance float64 `json:"distance"`
        EID      uint64  `json:"eid"`
        EName    string  `json:"e_ame"`
        EPhone   string  `json:"e_hone"`
    }

    err := ar.Ent.Store.QueryNotDeleted().
        Where(store.EmployeeIDNotNil(), store.Status(model.StoreStatusOpen)).
        Modify(
            func(sel *sql.Selector) {
                sel.Select(sel.C(store.FieldID), sel.C(store.FieldLng), sel.C(store.FieldLat), sel.C(store.FieldName)).
                    Where(sql.P(func(b *sql.Builder) {
                        b.WriteString(fmt.Sprintf(`ST_DWithin(ST_GeogFromText('SRID=4326;POINT('||%s||' '||%s||')'), ST_GeogFromText('POINT(%f %f)'), %f)`, sel.C(store.FieldLng), sel.C(store.FieldLat), ass.Lng, ass.Lat, 50000.0))
                    })).
                    AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeogFromText('SRID=4326;POINT('||%s||' '||%s||')'), ST_GeogFromText('POINT(%f %f)'))`, sel.C(store.FieldLng), sel.C(store.FieldLat), ass.Lng, ass.Lat)), "distance").
                    OrderBy(sql.Asc("distance"))

                // 查找employee
                emt := sql.Table(employee.Table)
                sel.Join(emt).
                    On(sel.C(store.FieldEmployeeID), emt.C(employee.FieldID)).
                    AppendSelect(
                        sql.As(emt.C(employee.FieldID), "eid"),
                        sql.As(emt.C(employee.FieldName), "e_ame"),
                        sql.As(emt.C(employee.FieldPhone), "e_hone"),
                    )

                // // 查找city
                // ct := sql.Table(city.Table)
                // sel.Join(ct).
                //     On(sel.C(store.FieldCityID), ct.C(city.FieldID)).
                //     AppendSelect(
                //         sql.As(ct.C(city.FieldID), "cid"),
                //         sql.As(ct.C(city.FieldName), "c_name"),
                //     )
            },
        ).
        Scan(s.ctx, &temps)

    if err != nil {
        log.Error(err)
    }

    items = make([]model.AssistanceNearbyRes, len(temps))

    for i, temp := range temps {
        items[i] = model.AssistanceNearbyRes{
            ID:       temp.ID,
            Name:     temp.Name,
            Lng:      temp.Lng,
            Lat:      temp.Lat,
            Distance: temp.Distance,
            Employee: model.Employee{
                ID:    temp.EID,
                Name:  temp.EName,
                Phone: temp.EPhone,
            },
        }
    }

    return
}

// Allocate 分配救援任务
func (s *assistanceService) Allocate(req *model.AssistanceAllocateReq) {
    st, _ := ar.Ent.Store.QueryNotDeleted().Where(store.ID(req.StoreID), store.EmployeeIDNotNil(), store.Status(model.StoreStatusOpen)).Only(s.ctx)
    if st == nil {
        snag.Panic("未找到营业中的门店")
    }

    item := s.QueryX(req.ID)

    _, err := item.Update().
        SetDistance(haversine.Distance(haversine.NewCoordinates(item.Lat, item.Lng), haversine.NewCoordinates(st.Lat, st.Lng)).Miles()).
        SetStoreID(st.ID).
        SetEmployeeID(*st.EmployeeID).
        SetStatus(model.AssistanceStatusAllocated).
        Save(s.ctx)

    // TODO 处理接单响应
    if err != nil {
        snag.Panic("分配失败")
    }

    // 记录日志
    go logging.NewOperateLog().
        SetRef(item).
        SetModifier(s.modifier).
        SetOperate(model.OperateAssistanceAllocate).
        SetDiff(model.AssistanceStatus(item.Status), model.AssistanceStatus(model.AssistanceStatusAllocated)).
        Send()
}

// Free 救援免费
func (s *assistanceService) Free(req *model.AssistanceFreeReq) {
    item, _ := s.orm.QueryNotDeleted().Where(assistance.Status(model.AssistanceStatusUnpaid)).First(s.ctx)
    if item == nil {
        snag.Panic("未找到待支付订单")
    }

    _, err := item.Update().SetFreeReason(req.Reason).SetCost(0).SetStatus(model.AssistanceStatusSuccess).Save(s.ctx)
    if err != nil {
        snag.Panic("处理失败")
    }

    // TODO 处理免费响应

    // 记录日志
    go logging.NewOperateLog().
        SetRef(item).
        SetModifier(s.modifier).
        SetOperate(model.OperateAssistanceFree).
        SetDiff(fmt.Sprintf("%s (%.2f元)", model.AssistanceStatus(item.Status), item.Cost), model.AssistanceStatus(model.AssistanceStatusSuccess)).
        Send()
}

// Refuse 拒绝救援
func (s *assistanceService) Refuse(req *model.AssistanceRefuseReq) {
    item := s.QueryX(req.ID)
    if item.Status == model.AssistanceStatusSuccess || item.Status == model.AssistanceStatusUnpaid {
        snag.Panic("救援状态错误")
    }

    _, err := item.Update().
        ClearEmployeeID().
        ClearStoreID().
        ClearCost().
        SetStatus(model.AssistanceStatusRefused).
        SetRefusedDesc(req.Desc).
        Save(s.ctx)

    if err != nil {
        snag.Panic("操作失败")
    }

    // TODO 处理拒绝响应

    // 记录日志
    go logging.NewOperateLog().
        SetRef(item).
        SetModifier(s.modifier).
        SetOperate(model.OperateAssistanceRefuse).
        SetDiff(model.AssistanceStatus(item.Status), model.AssistanceStatus(model.AssistanceStatusRefused)).
        Send()
}

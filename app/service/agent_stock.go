package service

import (
	"fmt"
	"math"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/stock"
	"github.com/auroraride/aurservd/pkg/tools"
)

type agentStockService struct {
	*BaseService
	orm *ent.StockClient
}

func NewAgentStock(params ...any) *agentStockService {
	return &agentStockService{
		BaseService: newService(params...),
		orm:         ent.Database.Stock,
	}
}

func (s *agentStockService) Detail(req *model.AgentStockDetailReq) *model.PaginationRes {
	q := s.orm.Query().
		Where(stock.EnterpriseIDEQ(req.EnterpriseID)).
		Order(ent.Desc(stock.FieldCreatedAt)).
		Modify(func(sel *sql.Selector) {
			// 去重排除配偶
			sel.FromExpr(sql.Raw("(SELECT DISTINCT ON (id + COALESCE(stock_spouse, 0)) * FROM stock) stock"))
		}).
		WithCabinet().
		WithSpouse(func(sq *ent.StockQuery) {
			sq.WithStore().WithCabinet().WithRider().WithStation().WithBattery().WithEnterprise()
		}).
		WithRider().
		WithCity().
		WithEbike().
		WithStation().
		WithBattery().
		WithEnterprise().
		WithAgent()

	if req.Start != "" {
		q.Where(stock.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
	}

	if req.End != "" {
		q.Where(stock.CreatedAtLT(tools.NewTime().ParseNextDateStringX(req.End)))
	}

	// 筛选物资类别
	if req.Materials == "" {
		req.Materials = fmt.Sprintf("%s,%s", stock.MaterialBattery, stock.MaterialEbike)
	} else {
		req.Materials = strings.ReplaceAll(req.Materials, " ", "")
	}
	materials := strings.Split(req.Materials, ",")

	if len(materials) > 0 {
		var predicates []predicate.Stock
		for _, material := range materials {
			switch stock.Material(material) {
			case stock.MaterialBattery:
				predicates = append(predicates, stock.ModelNotNil())
			case stock.MaterialEbike:
				predicates = append(predicates, stock.EbikeIDNotNil())
			case stock.MaterialOthers:
				predicates = append(predicates, stock.ModelIsNil())
			}
		}
		q = q.Where(stock.Or(predicates...))
	}

	if req.Type != 0 {
		q.Where(stock.Type(req.Type))
	}

	if req.Model != "" {
		q.Where(stock.Model(req.Model))
	}

	if req.Keyword != "" {
		// 搜索关键字 查询电柜编号、车架号、电池编码
		q.Where(
			stock.Or(
				stock.HasCabinetWith(cabinet.SerialContains(req.Keyword)),
				stock.HasEbikeWith(ebike.SnContains(req.Keyword)),
				stock.HasBatteryWith(battery.SnContains(req.Keyword)),
			),
		)
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Stock) model.AgentStockDetailRes {
		return s.detailInfo(item)
	})
}

func (s *agentStockService) detailInfo(item *ent.Stock) model.AgentStockDetailRes {
	res := model.AgentStockDetailRes{
		ID:     item.ID,
		Sn:     item.Sn,
		Name:   item.Name,
		Num:    int(math.Abs(float64(item.Num))),
		Time:   item.CreatedAt.Format(carbon.DateTimeLayout),
		Remark: item.Remark,
	}

	// 城市
	c := item.Edges.City
	if c != nil {
		res.City = c.Name
	}

	// 电车
	bike := item.Edges.Ebike
	if bike != nil {
		res.Name = item.Name
		res.MaterialSn = bike.Sn
	}

	em := item.Creator
	er := item.Edges.Rider
	ec := item.Edges.Cabinet
	en := item.Edges.Enterprise
	st := item.Edges.Station
	ba := item.Edges.Battery
	ag := item.Edges.Agent

	// 站点调拨电池
	if ba != nil {
		res.Name = *item.Model
		res.MaterialSn = ba.Sn
	}
	if item.Type == model.StockTypeTransfer {
		// 平台调拨记录
		res.Type = "平台调拨"
		res.Operator = "后台"
		if em != nil {
			res.Operator = fmt.Sprintf("后台 - %s", em.Name)
		}
		var sec *ent.Cabinet
		var sst *ent.EnterpriseStation
		var sen *ent.Enterprise
		sp := item.Edges.Spouse
		if sp != nil {
			sec = sp.Edges.Cabinet
			sst = sp.Edges.Station
			sen = sp.Edges.Enterprise
		}

		// 出入库对象判定
		if item.Num > 0 {
			res.Inbound = s.target(ec, st, en)
			res.Outbound = s.target(sec, sst, sen)
		} else {
			res.Inbound = s.target(sec, sst, sen)
			res.Outbound = s.target(ec, st, en)
		}
	} else {
		// 业务调拨记录
		var riderName string
		var agentName string

		if er != nil {
			riderName = er.Name
			res.Rider = fmt.Sprintf("%s - %s", riderName, er.Phone)
		}

		if ag != nil {
			agentName = ag.Name
		}

		tm := map[uint8]string{
			model.StockTypeRiderActive:      "新签",
			model.StockTypeRiderPause:       "寄存",
			model.StockTypeRiderContinue:    "取消寄存",
			model.StockTypeRiderUnSubscribe: "退租",
		}

		var tmr string
		if ec != nil {
			res.Operator = fmt.Sprintf("骑手 - %s", riderName)
			tmr = "电柜"
		} else {
			if item.Creator != nil {
				tmr = "后台"
				res.Operator = fmt.Sprintf("后台 - %s", item.Creator.Name)
			} else if st != nil {
				tmr = "站点"
				res.Operator = fmt.Sprintf("代理 - %s", agentName)
			}
		}

		res.Type = tmr + tm[item.Type]

		// 骑手
		target := model.TransferInfo{
			RiderName:  er.Name,
			RiderPhone: er.Phone,
		}
		switch item.Type {
		case model.StockTypeRiderActive, model.StockTypeRiderContinue:
			res.Inbound = target
			res.Outbound = s.target(ec, st, en)
		case model.StockTypeRiderPause, model.StockTypeRiderUnSubscribe:
			res.Inbound = s.target(ec, st, en)
			res.Outbound = target
		}
	}

	return res
}

// target 出入库对象
func (s *agentStockService) target(ec *ent.Cabinet, st *ent.EnterpriseStation, en *ent.Enterprise) (target model.TransferInfo) {
	target = model.TransferInfo{}
	if ec == nil && st == nil && en == nil {
		target.PlatformName = "平台"
		return
	}

	if ec != nil {
		target.CabinetName = ec.Name
		target.CabinetSerial = ec.Serial
	}
	if st != nil {
		target.EnterpriseName = en.Name
		target.StationName = st.Name
	}
	return
}

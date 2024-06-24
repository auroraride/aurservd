package task

import (
	"context"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/cabinetec"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/tools"
)

type cabinetTask struct {
}

func NewCabinet() *cabinetTask {
	return &cabinetTask{}
}

func (t *cabinetTask) Start() {
	c := cron.New()
	_, err := c.AddFunc("@daily", func() {
		zap.L().Info("开始执行 @daily[cabinetTask] 电柜耗电定时任务")
		t.Do()
	})
	if err != nil {
		zap.L().Fatal("@daily[orderTradePay] 定时任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

func (t *cabinetTask) Do() {
	ctx := context.Background()
	now := carbon.Now()
	start := now.SubDay().StartOfDay().StdTime()
	end := now.SubDay().EndOfDay().StdTime()
	month := now.StartOfMonth().StdTime()
	lastMonth := now.StartOfMonth().SubMonth().StdTime()
	_, _, data := biz.NewCabinet().ListECData(definition.CabinetECDataSearchOptions{
		Start: silk.Time(start),
		End:   silk.Time(end),
	})
	if len(data) == 0 {
		return
	}

	// 按照Serial分组 分组中只保存最大读数 和最小读数
	groups := make(map[string]*definition.GroupCabinetECData)
	for _, item := range data {
		group, exists := groups[item.Serial]
		if !exists {
			group = &definition.GroupCabinetECData{
				Max:   item,
				Min:   item,
				Total: 0,
			}
			groups[item.Serial] = group
		}
		if item.Value > group.Max.Value {
			group.Max = item
		}

		if item.Value < 0 {
			item.Value = 0
		}

		if item.Value < group.Min.Value {
			group.Min = item
		}
	}
	for _, group := range groups {
		group.Total = tools.NewDecimal().Sub(group.Max.Value, group.Min.Value)
	}
	bulkCreate := make([]*ent.CabinetEcCreate, 0)
	for k, v := range groups {
		// 查询电柜id
		cab, _ := ent.Database.Cabinet.Query().Where(cabinet.Serial(k)).First(ctx)
		if cab == nil {
			zap.L().Error("电柜不存在", zap.String("serial", k))
			continue
		}
		// 查询当月数据 如果没有则创建 当月总数等于上个月endec - 本月startec
		ec, _ := ent.Database.CabinetEc.Query().Where(cabinetec.Date(month.Format("2006-01")), cabinetec.Serial(k)).First(ctx)
		if ec == nil {
			// 查询上个月的数据
			lastec, _ := ent.Database.CabinetEc.Query().Where(cabinetec.Date(lastMonth.Format("2006-01")), cabinetec.Serial(k)).First(ctx)
			var total float64
			startEC := v.Min.Value
			total = v.Total
			if lastec != nil && lastec.End > 0 && v.Max.Value-lastec.End > 0 {
				total = tools.NewDecimal().Sub(v.Max.Value, lastec.End)
				startEC = lastec.End
			}
			bulkCreate = append(bulkCreate, ent.Database.CabinetEc.Create().
				SetSerial(k).
				SetStart(decimal.NewFromFloat(startEC).Round(2).InexactFloat64()).
				SetEnd(decimal.NewFromFloat(v.Max.Value).Round(2).InexactFloat64()).
				SetDate(month.Format("2006-01")).
				SetCabinet(cab).
				SetTotal(total))
			continue
		}

		// 如果当天已经更新过不在更新
		if ec.UpdatedAt.Format(time.DateOnly) == now.StdTime().Format(time.DateOnly) {
			continue
		}
		// 更新当月数据
		err := ec.Update().SetEnd(decimal.NewFromFloat(v.Max.Value).Round(2).InexactFloat64()).SetTotal(tools.NewDecimal().Sub(v.Max.Value, ec.Start)).Exec(ctx)
		if err != nil {
			zap.L().Error("更新电柜耗电量失败", zap.Error(err))
			return
		}
	}

	err := ent.Database.CabinetEc.CreateBulk(bulkCreate...).Exec(ctx)
	if err != nil {
		zap.L().Error("创建电柜耗电量失败", zap.Error(err))
		return
	}

}

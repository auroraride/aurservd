// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/auroraride/aurservd/app/ec"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/pkg/snag"
	"go.uber.org/zap"
)

type LogCallback func(data any)

type Provider interface {
	Cabinets() ([]*ent.Cabinet, error)
	Brand() string
	Logger() *Logger
	FetchStatus(serial string) (bool, model.CabinetBins, error)
	DoorOperate(code, serial, operation string, door int) bool
	Reboot(code string, serial string) bool
}

func Run() {
	// if ar.Config.Cabinet.Provider {
	//     StartCabinetProvider(NewKaixin())
	// }
}

// cabinetCity 电柜城市获取
func cabinetCity(cab *ent.Cabinet) string {
	c := cab.Edges.City
	if c == nil && cab.CityID != nil {
		c, _ = ent.Database.City.Query().Where(city.ID(*cab.CityID)).First(context.Background())
	}
	return c.Name
}

// StartCabinetProvider 开始执行任务
func StartCabinetProvider(providers ...Provider) {
	for _, p := range providers {
		provider := p
		times := 0
		go func() {
			for {
				times += 1
				start := time.Now()

				dooLoop(times, start, provider)

				// 写入电柜日志
				provider.Logger().Write(fmt.Sprintf("完成第%d次%s电柜状态轮询, 耗时%.2fs\n\n", times, provider.Brand(), time.Now().Sub(start).Seconds()))

				time.Sleep(1 * time.Minute)
			}
		}()
	}
}

func dooLoop(times int, start time.Time, provider Provider) {
	snag.WithPanicStack(func() {
		slsCfg := ar.Config.Aliyun.Sls
		provider.Logger().Write(fmt.Sprintf("开始第%d次%s电柜状态轮询\n", times, provider.Brand()))

		items, err := provider.Cabinets()
		if err != nil {
			zap.L().Error(provider.Brand()+"电柜获取失败", zap.Error(err))
			return
		}

		for _, item := range items {
			// 换电过程不查询状态
			if ec.Busy(item.Serial) {
				continue
			}

			// 更新电柜信息
			err = NewUpdater(item).DoUpdate()

			if err == nil {
				// 提交日志
				if item.Health == model.CabinetHealthStatusOnline {
					go func() {
						// 保存历史仓位信息(转换后的)
						lg := GenerateSlsStatusLogGroup(item)
						if lg != nil {
							err = ali.NewSls().PutLogs(slsCfg.Project, slsCfg.CabinetLog, lg)
							if err != nil {
								zap.L().Error("阿里云SLS提交失败", zap.Error(err))
							}
						}
					}()
				}
			}

			time.Sleep(time.Duration((60000-int(time.Now().Sub(start).Milliseconds()))/len(items)) * time.Millisecond)
		}
	})
}

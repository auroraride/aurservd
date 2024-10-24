// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-17
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
)

// GenerateSlsStatusLogGroup 生成status log日志
func GenerateSlsStatusLogGroup(cab *ent.Cabinet) (lg *sls.LogGroup) {
	t := tea.Uint32(uint32(time.Now().Unix()))
	lg = &sls.LogGroup{}
	logs := make([]*sls.Log, len(cab.Bin))
	for i, bin := range cab.Bin {
		c := new(logging.CabinetLog)
		_ = copier.Copy(c, bin)
		c.Serial = cab.Serial
		c.Errors = strings.Join(bin.ChargerErrors, ",")
		c.Brand = cab.Brand.String()
		c.Time = time.Now().Format(carbon.DateTimeLayout)
		logs[i] = &sls.Log{
			Time:     t,
			Contents: logging.GenerateLogContent(c),
		}
	}
	lg.Logs = logs
	return
}

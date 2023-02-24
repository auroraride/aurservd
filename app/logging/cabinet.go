// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
)

// CabinetLog 电柜日志
type CabinetLog struct {
    Brand       string           `json:"brand" sls:"品牌" index:"doc"`
    Serial      string           `json:"serial" sls:"编码" index:"doc"`
    Name        string           `json:"name" sls:"仓位" index:"doc"`
    BatterySN   string           `json:"batterySN" sls:"电池序列号" index:"doc"`
    Full        bool             `json:"full" sls:"是否满电" index:"doc"`
    Battery     bool             `json:"battery" sls:"是否有电池" index:"doc"`
    Electricity model.BatterySoc `json:"electricity" sls:"当前电量" index:"doc"`
    OpenStatus  bool             `json:"openStatus" sls:"是否开门" index:"doc"`
    DoorHealth  bool             `json:"doorHealth" sls:"仓门是否正常" index:"doc"`
    Current     float64          `json:"current" sls:"充电电流(A)" index:"doc"`
    Voltage     float64          `json:"voltage" sls:"电压(V)" index:"doc"`
    Errors      string           `json:"errors" sls:"故障信息" index:"doc"`
    Remark      string           `json:"remark" sls:"备注"`
    Time        string           `json:"time" sls:"记录时间" index:"doc"`
}

func (c *CabinetLog) GetLogstoreName() string {
    return ar.Config.Aliyun.Sls.CabinetLog
}

func (c *CabinetLog) Send() {
    PutLog(c)
}

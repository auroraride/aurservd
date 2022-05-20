// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import "github.com/auroraride/aurservd/app/model"

// CabinetLog 电柜日志
type CabinetLog struct {
    Brand       string                   `json:"brand" sls:"品牌,index"`
    Serial      string                   `json:"serial" sls:"编码,index"`
    Name        string                   `json:"name" sls:"仓位,index"`
    BatterySN   string                   `json:"batterySN" sls:"电池序列号"`
    Full        bool                     `json:"full" sls:"是否满电"`
    Battery     bool                     `json:"battery" sls:"是否有电池"`
    Electricity model.BatteryElectricity `json:"electricity" sls:"当前电量"`
    OpenStatus  bool                     `json:"openStatus" sls:"是否开门"`
    DoorHealth  bool                     `json:"doorHealth" sls:"仓门是否正常"`
    Current     float64                  `json:"current" sls:"充电电流(A)"`
    Voltage     float64                  `json:"voltage" sls:"电压(V)"`
    Errors      string                   `json:"errors" sls:"故障信息"`
    Remark      string                   `json:"remark" sls:"备注"`
    Time        string                   `json:"time" sls:"记录时间"`
}

// DoorOperateLog 柜门操作日志
type DoorOperateLog struct {
    ID        string `json:"id" sls:"操作ID"`
    Brand     string `json:"brand" sls:"品牌"`
    Serial    string `json:"serial" sls:"编码"`
    Name      string `json:"name" sls:"仓位"`
    Operation string `json:"operation" sls:"操作"`
    Success   bool   `json:"success" sls:"是否成功"`
    Remark    string `json:"remark" sls:"备注"`
    UserID    uint64 `json:"userID" sls:"操作人ID"`
    Phone     string `json:"phone" sls:"操作人电话"`
    User      string `json:"user" sls:"操作人"`
    Time      string `json:"time" sls:"操作时间"`
}

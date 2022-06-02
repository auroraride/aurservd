// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import "github.com/auroraride/aurservd/app/model"

// CabinetLog 电柜日志
type CabinetLog struct {
    Brand       string                   `json:"brand" sls:"品牌" index:"doc"`
    Serial      string                   `json:"serial" sls:"编码" index:"doc"`
    Name        string                   `json:"name" sls:"仓位" index:"doc"`
    BatterySN   string                   `json:"batterySN" sls:"电池序列号" index:"doc"`
    Full        bool                     `json:"full" sls:"是否满电" index:"doc"`
    Battery     bool                     `json:"battery" sls:"是否有电池" index:"doc"`
    Electricity model.BatteryElectricity `json:"electricity" sls:"当前电量" index:"doc"`
    OpenStatus  bool                     `json:"openStatus" sls:"是否开门" index:"doc"`
    DoorHealth  bool                     `json:"doorHealth" sls:"仓门是否正常" index:"doc"`
    Current     float64                  `json:"current" sls:"充电电流(A)" index:"doc"`
    Voltage     float64                  `json:"voltage" sls:"电压(V)" index:"doc"`
    Errors      string                   `json:"errors" sls:"故障信息" index:"doc"`
    Remark      string                   `json:"remark" sls:"备注"`
    Time        string                   `json:"time" sls:"记录时间" index:"doc"`
}

// DoorOperateLog 柜门操作日志
type DoorOperateLog struct {
    ID            string                        `json:"id" sls:"操作ID" index:"doc"`
    Brand         string                        `json:"brand" sls:"品牌" index:"doc"`
    Serial        string                        `json:"serial" sls:"编码" index:"doc"`
    Name          string                        `json:"name" sls:"仓位" index:"doc"`
    Operation     string                        `json:"operation" sls:"操作" index:"doc"`
    Success       bool                          `json:"success" sls:"是否成功" index:"doc"`
    Remark        string                        `json:"remark" sls:"备注"`
    OperatorID    uint64                        `json:"operatorID" sls:"操作人ID" index:"doc"`
    OperatorPhone string                        `json:"operatorPhone" sls:"操作人电话" index:"doc"`
    OperatorName  string                        `json:"operatorName" sls:"操作者" index:"doc"`
    OperatorRole  model.CabinetDoorOperatorRole `json:"operatorRole" sls:"操作人角色" string:"true" index:"doc"`
    Time          string                        `json:"time" sls:"操作时间" index:"doc"`
}

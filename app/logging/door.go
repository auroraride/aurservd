// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
)

// DoorOperateLog 柜门操作日志
type DoorOperateLog struct {
	ID            string                        `json:"id" sls:"操作ID" index:"doc"`
	Brand         string                        `json:"brand" sls:"品牌" index:"doc"`
	Serial        string                        `json:"serial" sls:"编码" index:"doc"`
	Name          string                        `json:"name" sls:"仓位" index:"doc"`
	Operation     string                        `json:"operation" sls:"操作" index:"doc"`
	Success       bool                          `json:"success" sls:"是否成功" index:"doc"`
	Remark        string                        `json:"remark" sls:"备注"`
	OperatorID    uint64                        `json:"operatorId" sls:"操作人ID" index:"doc"`
	OperatorPhone string                        `json:"operatorPhone" sls:"操作人电话" index:"doc"`
	OperatorName  string                        `json:"operatorName" sls:"操作者" index:"doc"`
	OperatorRole  model.CabinetDoorOperatorRole `json:"operatorRole" sls:"操作人角色" string:"true" index:"doc"`
	Time          string                        `json:"time" sls:"操作时间" index:"doc"`
}

func (d *DoorOperateLog) GetLogstoreName() string {
	return ar.Config.Aliyun.Sls.DoorLog
}

func (d *DoorOperateLog) Send() {
	PutLog(d)
}

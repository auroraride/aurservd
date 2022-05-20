// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logger

type Operate uint

const (
    OperateBanPerson   = iota // 封禁用户
    OperateUnBanPerson        // 解封用户
)

type OperateRefTable string

type Operation struct {
    ID       string          `json:"id" sls:"操作ID"`
    RefID    uint64          `json:"refId" sls:"关联ID"`
    RefTable OperateRefTable `json:"refTable" sls:"关联表"`
    Operate  Operate         `json:"operate" sls:"操作类别"`
    Success  bool            `json:"success" sls:"是否成功"`
    Remark   string          `json:"remark" sls:"备注"`

    Before string `json:"before" sls:"操作前"`
    After  string `json:"after" sls:"操作后"`

    ManagerID   uint64 `json:"managerID" sls:"操作人ID"`
    Phone       string `json:"phone" sls:"操作人电话"`
    ManagerName string `json:"managerName" sls:"操作人"`
    Time        string `json:"time" sls:"操作时间"`
}

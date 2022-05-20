// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

type Operate uint

const (
    OperateBanPerson   = iota // 封禁用户
    OperateUnBanPerson        // 解封用户
)

type OperateRefTable string

// OperationLog 系统操作日志
type OperationLog struct {
    ID       string          `json:"id" sls:"操作ID"`
    RefID    uint64          `json:"refId" sls:"关联ID,index"`
    RefTable OperateRefTable `json:"refTable" sls:"关联表,index"`

    Operate Operate `json:"operate" sls:"类别,index"`
    Remark  string  `json:"remark" sls:"备注"`
    Before  string  `json:"before" sls:"操作前"`
    After   string  `json:"after" sls:"操作后"`

    ManagerID    uint64 `json:"managerID" sls:"操作人ID,index"`
    ManagerPhone string `json:"phone" sls:"操作人电话,index"`
    ManagerName  string `json:"managerName" sls:"操作人,index"`

    Time string `json:"time" sls:"时间,index"`
}

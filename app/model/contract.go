// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-17
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    ContractStatusDraft ContractStatus = iota
    ContractStatusPending
    ContractStatusSuccess
    ContractStatusUndo
    ContractStatusOverdue
    ContractStatusRefused
)

type ContractStatus uint8

func (s ContractStatus) Value() uint8 {
    return uint8(s)
}

func (s ContractStatus) String() string {
    switch s {
    case ContractStatusDraft:
        return "草稿"
    case ContractStatusPending:
        return "签署中"
    case ContractStatusSuccess:
        return "完成"
    case ContractStatusUndo:
        return "撤销"
    case ContractStatusOverdue:
        return "过期"
    case ContractStatusRefused:
        return "拒签"
    }
    return "未知"
}

// ContractSignReq 签约请求
type ContractSignReq struct {
    SubscribeID uint64 `json:"subscribeId" validate:"required" trans:"订阅ID"`
}

// ContractSignRes 合同签订返回
type ContractSignRes struct {
    Url       string `json:"url"`       // 签署URL
    Sn        string `json:"sn"`        // 签署识别码
    Effective bool   `json:"effective"` // 是否存在生效中的合同
}

// ContractSignResultReq 合同签署结果请求
type ContractSignResultReq struct {
    Sn string `json:"sn" param:"sn" validate:"required"`
}

type ContractSignUniversal struct {
    Price string // 月租金
    Month int    // 首次缴纳月数
    Total string // 首次缴纳总计
    Stop  string // 结束日期
}

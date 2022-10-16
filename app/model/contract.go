// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-17
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    ContractExpiration = 30 // 过期时间 (分钟)
)

const (
    ContractStatusDraft      ContractStatus = iota // 草稿
    ContractStatusSigning                          // 签署中
    ContractStatusSuccess                          // 已完成
    ContractStatusRevoked                          // 已撤销
    ContractStatusTerminated                       // 已终止
    ContractStatusExpired                          // 已过期
    ContractStatusDenied     ContractStatus = 7    // 已拒绝
)

type ContractStatus uint8

func (s ContractStatus) Value() uint8 {
    return uint8(s)
}

func (s ContractStatus) String() string {
    switch s {
    case ContractStatusDraft:
        return "草稿中"
    case ContractStatusSigning:
        return "签署中"
    case ContractStatusSuccess:
        return "已完成"
    case ContractStatusRevoked:
        return "已撤销"
    case ContractStatusExpired:
        return "已过期"
    case ContractStatusTerminated:
        return "已拒签"
    case ContractStatusDenied:
        return "已拒绝"
    }
    return "未知"
}

// IsSuccessed 合同签署流程是否成功
func (s ContractStatus) IsSuccessed() bool {
    return s == ContractStatusSuccess
}

// IsFailed 合同签署流程是否失败
func (s ContractStatus) IsFailed() bool {
    return s != ContractStatusDraft && s != ContractStatusSigning && s != ContractStatusSuccess
}

// IsFinished 合同签署流程是否已完成
func (s ContractStatus) IsFinished() bool {
    return s != ContractStatusDraft && s != ContractStatusSigning
}

type ContractRider struct {
    Phone        string `json:"phone"`
    Name         string `json:"name"`
    IDCardNumber string `json:"idCardNumber"`
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

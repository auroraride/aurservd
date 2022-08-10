// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    ExportStatusProcessing uint8 = iota // 生成中
    ExportStatusSuccess                 // 已生成
    ExportStatusFail                    // 已失败
)

type ExportRes struct {
    SN string `json:"sn"` // 导出编号
}

type ExportInfoNameCallback func() string

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type ExportStatus uint8

const (
    ExportStatusProcessing ExportStatus = iota // 生成中
    ExportStatusSuccess                        // 已生成
    ExportStatusFail                           // 已失败
)

func (s ExportStatus) String() string {
    return map[ExportStatus]string{
        ExportStatusProcessing: "生成中",
        ExportStatusSuccess:    "已生成",
        ExportStatusFail:       "已失败",
    }[s]
}

type ExportRes struct {
    SN string `json:"sn"` // 导出编号
}

type ExportInfoNameCallback func() string

type ExportListReq struct {
    PaginationReq
}

type ExportListRes struct {
    CreatedAt string                 `json:"createdAt"`          // 创建时间
    Operator  string                 `json:"operator"`           // 操作人
    Remark    string                 `json:"remark"`             // 备注
    Taxonomy  string                 `json:"taxonomy"`           // 类别
    SN        string                 `json:"sn"`                 // 编号
    Status    string                 `json:"status"`             // 状态
    Message   string                 `json:"message,omitempty"`  // 错误信息, 可能不存在
    FinishAt  string                 `json:"finishAt,omitempty"` // 完成时间, 可能不存在
    Info      map[string]interface{} `json:"info,omitempty"`     // 筛选条件, 可能不存在
}

type ExportDownloadReq struct {
    SN string `json:"sn" param:"sn" validate:"required"` // 下载文件
}

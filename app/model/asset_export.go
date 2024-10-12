package model

type AssetExportStatus uint8

const (
	AssetExportStatusProcessing AssetExportStatus = iota // 生成中
	AssetExportStatusSuccess                             // 已生成
	AssetExportStatusFail                                // 已失败
)

func (s AssetExportStatus) String() string {
	return map[AssetExportStatus]string{
		AssetExportStatusProcessing: "生成中",
		AssetExportStatusSuccess:    "已生成",
		AssetExportStatusFail:       "已失败",
	}[s]
}

type AssetExportRes struct {
	SN string `json:"sn"` // 导出编号
}

type AssetExportInfoNameCallback func() string

type AssetExportListReq struct {
	PaginationReq
	SN string `json:"sn" query:"sn"` // 导出编号
}

type AssetExportListRes struct {
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

type AssetExportDownloadReq struct {
	SN    string `json:"sn" param:"sn" validate:"required"`       // 下载文件
	Token string `json:"token" query:"token" validate:"required"` // 用户token
}

package definition

import (
	"time"

	"github.com/auroraride/aurservd/app/model"
)

// CabinetByRiderReq 根据获取电柜请求参数
type CabinetByRiderReq struct {
	Lng      *float64 `json:"lng" query:"lng" validate:"required" trans:"经度"`                      // 经度
	Lat      *float64 `json:"lat" query:"lat" validate:"required" trans:"纬度"`                      // 纬度
	Keyword  *string  `json:"keyword" query:"keyword"`                                             // 关键字
	Distance *float64 `json:"distance" query:"distance" trans:"distance"`                          // 距离(米)
	Model    *string  `json:"model" query:"model"`                                                 // 电池型号
	Business *string  `json:"business" query:"business" enums:"active,pause,continue,unsubscribe"` // 业务选项 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
}

// CabinetByRiderRes 电柜列表响应
type CabinetByRiderRes struct {
	model.CabinetDataRes
	Distance    float64                     `json:"distance"`           // 距离
	Lng         float64                     `json:"lng"`                // 经度
	Lat         float64                     `json:"lat"`                // 纬度
	StockNum    int                         `json:"stockNum,omitempty"` // 库存电池
	BranchID    uint64                      `json:"branchId"`           // 网点ID
	Fid         string                      `json:"fid"`                // 电柜FID
	Address     string                      `json:"address"`            // 详细地址
	Reserve     *model.ReserveUnfinishedRes `json:"reserve,omitempty"`  // 当前预约, 预约不存在时无此字段
	Businesses  []string                    `json:"businesses"`         // 可办理业务  active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
	ExchangNum  int                         `json:"exchangNum"`         // 可换电池数量
	EmptyBinNum int                         `json:"emptyBinNum"`        // 可用空仓数量
}

// CabinetDetailReq 电柜详情请求参数
type CabinetDetailReq struct {
	Serial string `json:"serial" param:"serial" validate:"required"` // 电柜编号
}

// CabinetECDataSearchOptions 电柜能耗查询选项
// 若开始和结束时间均为空，查询当日数据
// 若开始时间为空，结束时间不为空，查询结束时间当日
// 若开始时间不为空，结束时间为空，查询为开始时间当日
type CabinetECDataSearchOptions struct {
	Serial *string    // 电柜编号
	Start  *time.Time // 开始时间，包含
	End    *time.Time // 结束时间，包含
}

// CabinetECData 电柜能耗查询结果
type CabinetECData struct {
	Serial    string    `json:"serial"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"@timestamp"`
}

// CabinetECMonthExportReq 电柜耗电量请求参数
type CabinetECMonthExportReq struct {
	Serial *string `json:"serial" query:"serial"` // 电柜编号
	Date   *string `json:"date" query:"date"`     // 日期 yyyy-MM
	Remark string  `json:"remark" query:"remark"` // 备注
}

// CabinetECMonthReq 电柜耗电列表
type CabinetECMonthReq struct {
	model.PaginationReq
	Serial *string `json:"serial" query:"serial"` // 电柜编号
	Date   *string `json:"date" query:"date"`     // 日期 yyyy-MM
}

// CabinetECRes 电柜耗电列表响应
type CabinetECRes struct {
	Serial  string  `json:"serial"`  // 电柜编号
	StartAt string  `json:"startAt"` // 开始时间
	EndAt   string  `json:"endAt"`   // 结束时间
	StartEc float64 `json:"startEc"` // 开始电量
	EndEc   float64 `json:"endEc"`   // 结束电量
	Totoal  float64 `json:"total"`   // 总电量
}

// GroupCabinetECData 电柜能耗分组数据
type GroupCabinetECData struct {
	Max   *CabinetECData
	Min   *CabinetECData
	Total float64
}

// CabinetECReq 电柜耗电量明细
type CabinetECReq struct {
	model.PaginationReq
	Serial *string `json:"serial" query:"serial"` // 电柜编号
	Start  *string `json:"start" query:"start"`   // 开始时间
	End    *string `json:"end" query:"end"`       // 结束时间
}

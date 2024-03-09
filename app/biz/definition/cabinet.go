package definition

import "github.com/auroraride/aurservd/app/model"

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
	Distance   float64                     `json:"distance"`             // 距离
	Lng        float64                     `json:"lng"`                  // 经度
	Lat        float64                     `json:"lat"`                  // 纬度
	StockNum   int                         `json:"stockNum,omitempty"`   // 库存电池
	BranchID   uint64                      `json:"branchId"`             // 网点ID
	Fid        string                      `json:"fid"`                  // 电柜FID
	Address    string                      `json:"address"`              // 详细地址
	Reserve    *model.ReserveUnfinishedRes `json:"reserve,omitempty"`    // 当前预约, 预约不存在时无此字段
	Businesses []string                    `json:"businesses,omitempty"` // 可办理业务  active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
}

// CabinetDetailReq 电柜详情请求参数
type CabinetDetailReq struct {
	Serial string `json:"serial" param:"serial" validate:"required"` // 电柜编号
}

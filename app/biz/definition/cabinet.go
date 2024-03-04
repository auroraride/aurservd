package definition

import "github.com/auroraride/aurservd/app/model"

type CabinetByRiderReq struct {
	Keyword  *string  `json:"keyword" query:"keyword"`                                             // 关键字
	Lng      *float64 `json:"lng" query:"lng" validate:"required" trans:"经度"`                      // 经度
	Lat      *float64 `json:"lat" query:"lat" validate:"required" trans:"纬度"`                      // 纬度
	Distance *float64 `json:"distance" query:"distance" trans:"distance"`                          // 距离(米)
	Model    *string  `json:"model" query:"model"`                                                 // 电池型号
	Business *string  `json:"business" query:"business" enums:"active,pause,continue,unsubscribe"` // 业务选项 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
}

type CabinetByRiderRes struct {
	model.CabinetDataRes
	Distance float64 `json:"distance"`           // 距离
	Lng      float64 `json:"lng"`                // 经度
	Lat      float64 `json:"lat"`                // 纬度
	StockNum int     `json:"stockNum,omitempty"` // 库存电池
	BranchID uint64  `json:"branchId"`           // 网点ID
	Fid      string  `json:"fid"`                // 电柜FID
}

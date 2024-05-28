// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-05-28, by Jorjan

package definition

// BranchWithDistanceReq 根据距离获取网点请求
type BranchWithDistanceReq struct {
	Lng           *float64 `json:"lng" query:"lng" validate:"required" trans:"经度"`
	Lat           *float64 `json:"lat" query:"lat" validate:"required" trans:"纬度"`
	Distance      *float64 `json:"distance" query:"distance" trans:"距离"`
	CityID        *uint64  `json:"cityId" query:"cityId" trans:"城市ID"`
	Business      string   `json:"business" query:"business" enums:"active,pause,continue,unsubscribe"` // 业务选项 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
	Filter        string   `json:"filter" query:"filter"`                                               // 额外筛选参数
	Model         *string  `json:"model" query:"model"`                                                 // 电池型号
	StoreStatus   *uint8   `json:"storeStatus" query:"storeStatus"`                                     // 门店状态 1-营业 2-休息
	StoreBusiness *uint8   `json:"storeBusiness" query:"storeBusiness"`                                 // 门店业务 1-租车，2-修车，3-买车，4-驿站
}

// BranchWithDistanceRes 根据距离获取网点结果
type BranchWithDistanceRes struct {
	ID          uint64                     `json:"id"`       // 网点ID
	Distance    float64                    `json:"distance"` // 距离(前端处理: 超过1000米显示nKM)
	Name        string                     `json:"name"`     // 网点名称
	Lng         float64                    `json:"lng"`      // 经度
	Lat         float64                    `json:"lat"`      // 纬度
	Image       string                     `json:"image"`    // 网点图片
	Photos      []string                   `json:"photos"`   // 网点图片(V2)
	Address     string                     `json:"address"`  // 网点地址
	Facility    []*BranchFacility          `json:"facility"` // 网点设施
	FacilityMap map[string]*BranchFacility `json:"-" swaggerignore:"true"`
	Businesses  []string                   `json:"businesses"` // 可办理业务 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
}

// BranchFacility 网点设施
type BranchFacility struct {
	ID         uint64 `json:"id"`
	Fid        string `json:"fid"`             // 设施标识
	Type       string `json:"type"`            // 类别
	Name       string `json:"name"`            // 名称
	State      uint   `json:"state"`           // 状态 0不可用 1可用
	Num        int    `json:"num"`             // 满电数量
	Phone      string `json:"phone,omitempty"` // 联系电话
	Total      int    `json:"total"`           // 仓位数量
	CabinetNum int    `json:"cabinetNum"`      // 电柜数量
	BatteryNum int    `json:"batteryNum"`      // 电池数量
}

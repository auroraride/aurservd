// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type BranchFacilityType string

const (
	BranchFacilityTypeAll   BranchFacilityType = "all"
	BranchFacilityTypeStore BranchFacilityType = "store"
	BranchFacilityTypeV72   BranchFacilityType = "v72"
	BranchFacilityTypeV60   BranchFacilityType = "v60"
	BranchFacilityTypeRest  BranchFacilityType = "rest"
)

func (b BranchFacilityType) String() string {
	return string(b)
}

const (
	BranchFacilityStateOffline uint = iota // 不在线
	BranchFacilityStateOnline              // 在线
)

const (
	BranchFacilityFilterEbike       = "ebike"
	BranchFacilityFilterEbikeObtain = "ebikeObtain"
	BranchFacilityFilterEbikeRepair = "ebikeRepair"
)

type Branch struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"` // 网点名称
}

// BranchListReq 后台网点列表请求
type BranchListReq struct {
	PaginationReq

	Name   *string `json:"name" query:"name"`     // 网点名称
	CityID *uint64 `json:"cityId" query:"cityId"` // 城市ID
}

// BranchCreateReq 创建网点
type BranchCreateReq struct {
	CityID    *uint64           `json:"cityId" validate:"required" trans:"城市"`
	Name      *string           `json:"name" validate:"required" trans:"网点名称"`
	Lng       *float64          `json:"lng" validate:"required" trans:"经度"`
	Lat       *float64          `json:"lat" validate:"required" trans:"纬度"`
	Address   *string           `json:"address" validate:"required" trans:"详细地址"`
	Photos    []string          `json:"photos" validate:"required" trans:"网点照片"`
	Contracts []*BranchContract `json:"contracts,omitempty"`
}

// BranchModifyReq 编辑网点请求
type BranchModifyReq struct {
	ID uint64 `json:"id,omitempty" param:"id"`

	CityID    *uint64           `json:"cityId" trans:"城市"`
	Name      *string           `json:"name" trans:"网点名称"`
	Lng       *float64          `json:"lng" trans:"经度"`
	Lat       *float64          `json:"lat" trans:"纬度"`
	Address   *string           `json:"address" trans:"详细地址"`
	Photos    *[]string         `json:"photos" trans:"网点照片"`
	Contracts *[]BranchContract `json:"contracts,omitempty"`
}

// BranchItem 网点列表返回
type BranchItem struct {
	ID         uint64           `json:"id"`
	Name       string           `json:"name"`           // 名称
	Lng        float64          `json:"lng"`            // 经度
	Lat        float64          `json:"lat"`            // 纬度
	Address    string           `json:"address"`        // 地址
	Photos     []string         `json:"photos"`         // 照片
	City       City             `json:"city,omitempty"` // 城市
	Contracts  []BranchContract `json:"contracts"`      // 合同
	StoreTotal int              `json:"storeTotal"`     // 门店数量
	V72Total   int              `json:"v72Total"`       // V72电柜数量
	V60Total   int              `json:"v60Total"`       // V60电柜数量
}

type BranchSampleItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// BranchContract 网点合同请求体
type BranchContract struct {
	ID       uint64 `json:"id,omitempty" swaggerignore:"true"`
	BranchID uint64 `json:"branchId,omitempty" param:"id" swaggerignore:"true"`

	LandlordName      string   `json:"landlordName" validate:"required" trans:"房东姓名"`
	IDCardNumber      string   `json:"idCardNumber" validate:"required" trans:"房东身份证"`
	Phone             string   `json:"phone" validate:"required,phone" trans:"房东手机号"`
	BankNumber        string   `json:"bankNumber" trans:"房东银行卡号"`
	Pledge            float64  `json:"pledge" trans:"押金"`
	Rent              float64  `json:"rent" validate:"required" trans:"租金"`
	Lease             uint     `json:"lease" validate:"required" trans:"租期月数"`
	ElectricityPledge float64  `json:"electricityPledge" trans:"电费押金"`
	Electricity       string   `json:"electricity" validate:"required" trans:"电费单价"`
	Area              float64  `json:"area" trans:"网点面积"`
	StartTime         string   `json:"startTime" validate:"required" trans:"租期开始时间"`
	EndTime           string   `json:"endTime" validate:"required" trans:"租期结束时间"`
	File              string   `json:"file" validate:"required" trans:"合同文件"`
	Sheets            []string `json:"sheets" validate:"required" trans:"底单"`
}

type BranchContractSheetReq struct {
	ID     uint64   `json:"id" validate:"required" trans:"合同ID"`
	Sheets []string `json:"sheets"` // 底单, 需携带所有`未删除`的底单
}

// BranchWithDistanceReq 根据距离获取网点请求
type BranchWithDistanceReq struct {
	Lng           *float64           `json:"lng" query:"lng" validate:"required" trans:"经度"`
	Lat           *float64           `json:"lat" query:"lat" validate:"required" trans:"纬度"`
	Distance      *float64           `json:"distance" query:"distance" trans:"距离"`
	CityID        *uint64            `json:"cityId" query:"cityId" trans:"城市ID"`
	Business      string             `json:"business" query:"business" enums:"active,pause,continue,unsubscribe"` // 业务选项 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
	Filter        string             `json:"filter" query:"filter"`                                               // 额外筛选参数
	Model         *string            `json:"model" query:"model"`                                                 // 电池型号
	StoreStatus   *uint8             `json:"storeStatus" query:"storeStatus"`                                     // 门店状态 0-全部 1-营业 2-休息 （首页地图不要传入此值）
	StoreBusiness *StoreBusinessType `json:"storeBusiness" query:"storeBusiness"`                                 // 门店业务 0-全部 1-租车，2-修车，3-买车，4-驿站 （首页地图不要传入此值）
}

type BranchDistanceListReq struct {
	Lng      float64            `json:"lng" query:"lng"`           // 经度, 默认 `108.947713`
	Lat      float64            `json:"lat" query:"lat"`           // 纬度, 默认 `34.231657`
	Distance float64            `json:"distance" query:"distance"` // 请求距离(米), 默认 `50000`
	Type     BranchFacilityType `json:"type" query:"type"`         // 筛选类别 all:全部 store-门店 v72-v72电柜 v60-v60电柜 rest-驿站门店
	Name     string             `json:"name" query:"name"`         // 门店或电柜名称
}

type BranchDistanceListRes struct {
	ID       uint64                     `json:"id"`       // 网点ID
	Distance float64                    `json:"distance"` // 距离(前端处理: 超过1000米显示nKM)
	Name     string                     `json:"name"`     // 网点名称
	Lng      float64                    `json:"lng"`      // 经度
	Lat      float64                    `json:"lat"`      // 纬度
	Cabinets []CabinetListByDistanceRes `json:"cabinets"`
	Stores   []StoreWithStatus          `json:"stores"`
}

// BranchFacility 网点设施
type BranchFacility struct {
	ID          uint64             `json:"id"`
	Fid         string             `json:"fid"`             // 设施标识
	Type        BranchFacilityType `json:"type"`            // 类别  all:全部 store-门店 v72-v72电柜 v60-v60电柜 rest-驿站门店
	Name        string             `json:"name"`            // 名称
	State       uint               `json:"state"`           // 状态 0不可用 1可用
	Num         int                `json:"num"`             // 满电数量
	Phone       string             `json:"phone,omitempty"` // 联系电话
	Total       int                `json:"total"`           // 仓位数量
	CabinetNum  int                `json:"cabinetNum"`      // 电柜数量
	BatteryNum  int                `json:"batteryNum"`      // 电池数量
	ExchangNum  int                `json:"exchangNum"`      // 可换电池数量
	EmptyBinNum int                `json:"emptyBinNum"`     // 可用空仓数量
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

// BranchRidingReq 骑行规划时间请求
type BranchRidingReq struct {
	Origin      string `json:"origin" query:"origin" trans:"开始坐标"`
	Destination string `json:"destination" query:"destination" trans:"结束坐标"`
}

// BranchRidingRes 骑行规划时间返回
type BranchRidingRes struct {
	Minutes float64 `json:"minutes"` // 骑行规划时间(分钟)
}

type BranchExpriesNotice struct {
	City string `json:"city"`
	Name string `json:"name"`
	End  string `json:"end"`
}

type BranchFacilityReq struct {
	Fid string  `json:"fid" param:"fid"` // 设施标识
	Lng float64 `json:"lng" query:"lng"` // 经度
	Lat float64 `json:"lat" query:"lat"` // 纬度
}

type BranchFacilityRes struct {
	Name     string                  `json:"name"`                       // 网点名称
	Address  string                  `json:"address"`                    // 地址
	Lng      float64                 `json:"lng"`                        // 经度
	Lat      float64                 `json:"lat"`                        // 纬度
	Distance float64                 `json:"distance"`                   // 距离(前端处理: 超过1000米显示nKM)
	Type     string                  `json:"type" enums:"store,cabinet"` // 设施类型 store:门店 cabinet:电柜(此时cabinet字段为数组)
	Image    string                  `json:"image"`                      // 网点图片
	Photos   []string                `json:"photos"`                     // 网点图片(V2)
	Store    *BranchFacilityStore    `json:"store,omitempty"`            // 门店, 当type=store时存在
	Cabinet  []BranchFacilityCabinet `json:"cabinet,omitempty"`          // 电柜, 当type=cabinet时存在, 根据序号显示 1号柜/2号柜 等, 当仅有一个电柜时, 电柜切换tab隐藏
}

type BranchFacilityStore struct {
	Name   string   `json:"name"`   // 门店名称
	Models []string `json:"models"` // 可用电池型号
}

type BranchFacilityCabinet struct {
	ID                uint64                         `json:"id"`                // 电柜ID
	Status            uint8                          `json:"status"`            // 电柜状态 0:离线 1:在线 2:维护中
	Name              string                         `json:"name"`              // 电柜名称
	Serial            string                         `json:"serial"`            // 电柜编号
	Batteries         []BranchFacilityCabinetBattery `json:"batteries"`         // 电池情况
	Reserve           *ReserveUnfinishedRes          `json:"reserve,omitempty"` // 当前预约, 预约不存在时无此字段
	Bins              []BranchFacilityCabinetBin     `json:"bins"`              // 仓位详情
	Businesses        []string                       `json:"businesses"`        // 骑手可办理业务 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
	CabinetBusinesses []string                       `json:"cabinetBusinesses"` // 电柜可办理业务 active:激活, pause:寄存, continue:取消寄存, unsubscribe:退租
}

type BranchFacilityCabinetBattery struct {
	Voltage     string `json:"voltage"`     // 电压, 单位V
	Capacity    string `json:"capacity"`    // 容量, 单位AH
	Charging    int    `json:"charging"`    // 充电数量
	Fully       int    `json:"fully"`       // 满电数量
	ExchangNum  int    `json:"exchangNum"`  // 可换电池数量
	EmptyBinNum int    `json:"emptyBinNum"` // 可用空仓数量
}

type BranchFacilityCabinetBin struct {
	Status      uint8       `json:"status"`                // 状态 0:空仓 1:充电 2:可用 3:锁仓
	Electricity *BatterySoc `json:"electricity,omitempty"` // 当前电量 锁仓或空仓无此字段
	BatterySN   string      `json:"batterySN,omitempty"`   // 电池SN码 锁仓或空仓无此字段
}

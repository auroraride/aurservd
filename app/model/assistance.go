// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    AssistanceStatusPending   uint8 = iota // 待分配
    AssistanceStatusAllocated              // 已分配
    AssistanceStatusRefused                // 已拒绝
    AssistanceStatusFailed                 // 救援失败
    AssistanceStatusUnpaid                 // 救援成功 - 待支付
    AssistanceStatusSuccess                // 救援成功 - 已支付
)

// AssistanceStatusProcessing 获取进行中的状态
var AssistanceStatusProcessing = []uint8{
    AssistanceStatusPending,
    AssistanceStatusAllocated,
    AssistanceStatusUnpaid,
}

func AssistanceStatus(u uint8) string {
    switch u {
    case AssistanceStatusAllocated:
        return "已分配"
    case AssistanceStatusRefused:
        return "已拒绝"
    case AssistanceStatusFailed:
        return "已失败"
    case AssistanceStatusUnpaid:
        return "待支付"
    case AssistanceStatusSuccess:
        return "已支付"
    default:
        return "待分配"
    }
}

type AssistanceCreateReq struct {
    Lng             float64  `json:"lng" validate:"required" trans:"经度"`
    Lat             float64  `json:"lat" validate:"required" trans:"纬度"`
    Address         string   `json:"address" validate:"required" trans:"详细地址"`
    Breakdown       string   `json:"breakdown" validate:"required" trans:"故障原因"`
    BreakdownDesc   string   `json:"breakdownDesc" validate:"required" trans:"故障描述"`
    BreakdownPhotos []string `json:"breakdownPhotos" validate:"required,min=1,max=3" trans:"故障图片"`
}

type AssistanceCreateRes struct {
    ID uint64
}

type AssistanceNearbyRes struct {
    ID       uint64   `json:"id"`       // 门店ID
    Name     string   `json:"name"`     // 门店名称
    Lng      float64  `json:"lng"`      // 经度
    Lat      float64  `json:"lat"`      // 纬度
    Distance float64  `json:"distance"` // 距离
    Employee Employee `json:"employee"` // 当前员工
}

type AssistanceListReq struct {
    PaginationReq

    CityID  uint64 `json:"cityId"`                     // 城市ID
    Keyword string `json:"keyword"`                    // 骑手姓名或电话
    Start   string `json:"start"`                      // 救援订单发起时间开始
    End     string `json:"end"`                        // 救援订单发起时间结束
    Status  *uint8 `json:"status" enums:"0,1,2,3,4,5"` // 状态 0:待分配 1:已分配 2:已拒绝 3:已失败 4:待支付 5:已支付, 不携带此字段是获取全部
}

type AssistanceListRes struct {
    ID       uint64     `json:"id"`
    Status   uint8      `json:"status" enums:"0,1,2,3,4,5"` // 状态 0:待分配 1:已分配 2:已拒绝 3:已失败 4:待支付 5:已支付
    Cost     float64    `json:"cost"`                       // 费用, 未分配时为0
    Distance float64    `json:"distance"`                   // 救援距离(米), 未分配时为0
    Time     string     `json:"time"`                       // 救援发起时间
    Rider    RiderBasic `json:"rider"`                      // 骑手信息
    City     City       `json:"city"`                       // 城市
    Employee *Employee  `json:"employee,omitempty"`         // 店员信息, 未分配此字段不存在
    Store    *Store     `json:"store,omitempty"`            // 门店, 未分配此字段不存在
}

type AssistanceDetail struct {
    AssistanceListRes
    Lng             float64  `json:"lng"`                   // 经度
    Lat             float64  `json:"lat"`                   // 纬度
    Address         string   `json:"address"`               // 详细位置
    Breakdown       string   `json:"breakdown"`             // 故障
    BreakdownDesc   string   `json:"breakdownDesc"`         // 故障描述
    BreakdownPhotos []string `json:"breakdownPhotos"`       // 故障照片
    Reason          string   `json:"reason"`                // 成功救援 - 故障原因
    DetectPhoto     string   `json:"detectPhoto"`           // 检测照片
    JointPhoto      string   `json:"jointPhoto"`            // 合照
    AllocateAt      string   `json:"allocateAt"`            // 分配时间
    RefusedDesc     *string  `json:"refusedDesc,omitempty"` // 拒绝原因, 已拒绝会有此字段
    PayAt           *string  `json:"payAt,omitempty"`       // 支付时间, 已支付会有此字段
    FreeReason      *string  `json:"freeReason,omitempty"`  // 免费理由, 当订单被标记为免费的时候有此字段
    FailReason      *string  `json:"failReason,omitempty"`  // 失败原因, 救援失败的时候有此字段
}

type AssistanceAllocateReq struct {
    ID      uint64 `json:"id" validate:"required" trans:"救援ID"`
    StoreID uint64 `json:"storeId" validate:"required" trans:"门店ID"`
}

type AssistanceFreeReq struct {
    ID     uint64 `json:"id" validate:"required" trans:"救援ID"`
    Reason string `json:"reason" validate:"required" trans:"免费理由"`
}

type AssistanceRefuseReq struct {
    ID     uint64 `json:"id" validate:"required" trans:"救援ID"`
    Reason string `json:"reason" validate:"required" trans:"拒绝原因"`
}

type AssistanceCancelReq struct {
    ID     uint64 `json:"id" validate:"required" trans:"救援ID"`
    Reason string `json:"reason" validate:"required" trans:"取消原因"`
    Desc   string `json:"desc"` // 取消原因详细描述
}

type AssistanceEmployeeDetailRes struct {
    ID              uint64      `json:"id"`
    Status          uint8       `json:"status"`             // 状态 0:待分配 1:已分配 2:已拒绝 3:已失败 4:待支付 5:已支付
    Rider           RiderBasic  `json:"rider"`              // 骑手信息
    Lng             float64     `json:"lng"`                // 经度
    Lat             float64     `json:"lat"`                // 纬度
    Address         string      `json:"address"`            // 详细位置
    Breakdown       string      `json:"breakdown"`          // 故障
    BreakdownDesc   string      `json:"breakdownDesc"`      // 故障描述
    BreakdownPhotos []string    `json:"breakdownPhotos"`    // 故障照片
    Reason          string      `json:"reason"`             // 成功救援 - 故障原因
    DetectPhoto     string      `json:"detectPhoto"`        // 检测照片
    JointPhoto      string      `json:"jointPhoto"`         // 合照
    FailReason      string      `json:"failReason"`         // 失败原因
    Store           StoreLngLat `json:"store"`              // 门店详情
    Time            string      `json:"time"`               // 发起时间
    Distance        string      `json:"distance,omitempty"` // 救援距离
    Minutes         float64     `json:"minutes,omitempty"`  // 总用时 (分钟)
    Model           string      `json:"model"`              // 电池型号

    Polylines []string `json:"polylines,omitempty"` // 路径规划
    Configure struct {
        Breakdown []interface{} `json:"breakdown"` // 救援原因<选择>
    } `json:"configure,omitempty"`
}

type AssistanceProcessReq struct {
    ID uint64 `json:"id" validate:"required"`

    Success    bool   `json:"success"`    // 救援结果, TRUE成功 FALSE失败
    FailReason string `json:"failReason"` // 失败原因, 救援失败的时候必填

    Pay         bool   `json:"pay"`         // 是否需要付费, 救援成功需要店员判断是否需要付费
    Reason      string `json:"reason"`      // 救援原因, 救援成功必填
    DetectPhoto string `json:"detectPhoto"` // 检测照片, 救援成功必填
    JointPhoto  string `json:"jointPhoto" ` // 合照, 救援成功必填
}

type AssistanceProcessRes struct {
    Cost float64 `json:"cost"` // 待支付金额, 待支付为0则无需支付
}

type AssistancePayReq struct {
    ID     uint64 `json:"id" validate:"required" trans:"救援ID"`
    Payway uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2"` // 1支付宝 2微信
}

type AssistanceSimpleListRes struct {
    ID       uint64     `json:"id"`
    Status   uint8      `json:"status"`          // 状态 0:待分配 1:已分配 2:已拒绝 3:已失败 4:待支付 5:已支付
    Rider    RiderBasic `json:"rider"`           // 骑手信息
    Cost     float64    `json:"cost"`            // 费用
    Time     string     `json:"time"`            // 发起时间
    Reason   string     `json:"reason"`          // 救援原因
    Distance string     `json:"distance"`        // 救援距离
    Model    string     `json:"model"`           // 电池型号
    Store    *Store     `json:"store,omitempty"` // 门店
}

type AssistancePayRes struct {
    AssistanceSimpleListRes

    OutTradeNo string `json:"outTradeNo"` // 订单二维码 (用做查询支付结果)
    QR         string `json:"qr"`         // 支付码
}

type AssistanceEmployeeOverview struct {
    Times    int     `json:"times"`    // 总次数
    Success  int     `json:"success"`  // 成功次数
    Distance float64 `json:"distance"` // 总里程
}

type AssistanceNotice struct {
    Phone         string `json:"phone"`
    Reason        string `json:"reason"`
    Address       string `json:"address"`
    AddressDetail string `json:"addressDetail"`
}

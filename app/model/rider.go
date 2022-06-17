// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// RiderTokenPermission 骑手token权限, 以此判定登陆后动作
type RiderTokenPermission uint8

const (
    RiderTokenPermissionCommon    RiderTokenPermission = iota // 普通权限
    RiderTokenPermissionAuth                                  // 需要实名验证
    RiderTokenPermissionNewDevice                             // 更换设备需要人脸验证
)

const (
    RiderStatusNormal  uint8 = iota + 1 // 未认证
    RiderStatusBlocked                  // 禁用
    RiderStatusBanned                   // 黑名单
)

// RiderContext TODO 骑手上下文
type RiderContext struct {
}

type RiderBasic struct {
    ID    uint64 `json:"id"`
    Phone string `json:"phone"` // 电话
    Name  string `json:"name"`  // 姓名
}

// RiderSignupReq 骑手登录请求数据
type RiderSignupReq struct {
    Phone   string `json:"phone" validate:"required" trans:"电话"`
    SmsId   string `json:"smsId" validate:"required" trans:"短信ID"`
    SmsCode string `json:"smsCode" validate:"required" trans:"短信验证码"`
    CityID  uint64 `json:"cityId"`
}

// RiderSigninRes 骑手登录数据返回
type RiderSigninRes struct {
    ID              uint64           `json:"id"`
    Phone           string           `json:"phone"`                     // 电话
    Name            string           `json:"name"`                      // 姓名, 实名认证后才会有
    Token           string           `json:"token,omitempty"`           // 认证token
    IsNewDevice     bool             `json:"isNewDevice"`               // 是否新设备
    IsAuthed        bool             `json:"isAuthed"`                  // 是否已认证
    IsContactFilled bool             `json:"isContactFilled"`           // 联系人是否添加
    Contact         *RiderContact    `json:"contact,omitempty"`         // 联系人
    Qrcode          string           `json:"qrcode"`                    // 二维码
    Deposit         float64          `json:"deposit"`                   // 需缴押金
    OrderNotActived *bool            `json:"orderNotActived,omitempty"` // 是否存在未激活订单
    Subscribe       *Subscribe       `json:"subscribe,omitempty"`       // 骑士卡
    Enterprise      *EnterpriseBasic `json:"enterprise,omitempty"`      // 所属企业
}

// RiderContact 紧急联系人
type RiderContact struct {
    Name     string `json:"name" validate:"required" trans:"联系人姓名"`
    Phone    string `json:"phone" validate:"required,phone" trans:"联系人电话"`
    Relation string `json:"relation" validate:"required" trans:"关系"`
}

// RiderSampleInfo 骑手简单信息
type RiderSampleInfo struct {
    ID    uint64 `json:"id"`    // 骑手ID
    Name  string `json:"name"`  // 骑手姓名
    Phone string `json:"phone"` // 骑手电话
}

// RiderListReq 骑手列表请求
type RiderListReq struct {
    PaginationReq

    Keyword  *string `json:"keyword" query:"keyword"`   // 搜索关键词
    Modified *bool   `json:"modified" query:"modified"` // 是否被修改过
    Start    *string `json:"start" query:"start"`       // 注册开始时间, 格式为: 2022-01-01
    End      *string `json:"end" query:"end"`           // 注册结束时间, 格式为: 2022-01-01

    Status          *uint8            `json:"status" enums:"0,1,2,3,4" query:"status"`                           // 用户状态 1:正常 2:已禁用 3:黑名单
    SubscribeStatus *uint8            `json:"subscribeStatus" enums:"0,1,2,3,4,5,11,99" query:"subscribeStatus"` // 业务状态 0:未激活 1:计费中 2:寄存中 3:已逾期 4:已退订 5:已取消 11: 即将到期 99:未使用
    AuthStatus      *PersonAuthStatus `json:"authStatus" enums:"0,1,2,3" query:"authStatus"`                     // 认证状态 0:未认证 1:认证中 2:已认证 3:认证失败
    PlanID          *uint64           `json:"planId" query:"planId"`                                             // 骑士卡
}

// RiderItemSubscribe 骑手骑士卡简单信息
type RiderItemSubscribe struct {
    ID        uint64  `json:"id"`        // 订阅ID
    Status    uint8   `json:"status"`    // 订阅状态 0:未激活 1:计费中 2:寄存中 3:已逾期 4:已退订 5:已取消 11: 即将到期(计算状态) 当 status = 1 且 remaining <= 3 的时候是即将到期
    Remaining int     `json:"remaining"` // 剩余天数
    Voltage   float64 `json:"voltage"`   // 骑士卡可用电压型号
}

// RiderItem 骑手信息
type RiderItem struct {
    ID         uint64           `json:"id"`
    Name       string           `json:"name"`               // 用户姓名
    Phone      string           `json:"phone"`              // 手机号
    Status     uint8            `json:"status"`             // 用户状态, 优先显示状态值大的 1:正常 2:已禁用 3:黑名单
    AuthStatus PersonAuthStatus `json:"authStatus"`         // 认证状态 0:未认证 1:认证中 2:已认证 3:认证失败
    Deposit    float64          `json:"deposit"`            // 押金
    Address    string           `json:"address"`            // 户籍地址
    DeletedAt  string           `json:"deleteAt,omitempty"` // 账户删除时间
    Remark     string           `json:"remark"`             // 账户备注
    Contract   string           `json:"contract,omitempty"` // 合同(有可能不存在)
    // 团签企业信息, 若无此字段则为个签用户
    Enterprise *EnterpriseBasic `json:"enterprise,omitempty"`
    // 当前有效订阅信息, 若无此字段则代表当前无有效订阅 (订阅 = 骑手骑士卡)
    Subscribe *RiderItemSubscribe `json:"subscribe,omitempty"`
    // 认证信息, 有可能不存在, 内部字段也有可能不存在
    Person *Person `json:"person,omitempty"`
    // 紧急联系人, 有可能不存在
    Contact *RiderContact `json:"contact"`
}

// RiderBlockReq 封禁或解封骑手账号
type RiderBlockReq struct {
    ID    uint64 `json:"id" `   // 骑手ID
    Block bool   `json:"block"` // `true`封禁 `false`解封
}

const (
    RiderLogTypeAll           uint8 = iota // 全部
    RiderLogTypeStatus                     // 状态
    RiderLogTypeProfile                    // 资料
    RiderLogTypeSubscribeDate              // 时长
    RiderLogTypeDeposit                    // 押金
)

var (
    RiderLogTypes = map[uint8][]Operate{
        RiderLogTypeStatus: {
            OperatePersonBan,
            OperatePersonUnBan,
            OperateRiderBLock,
            OperateRiderUnBLock,
        },
        RiderLogTypeProfile: {
            OperateProfile,
        },
        RiderLogTypeSubscribeDate: {
            OperateSubscribeAlter,
            OperateSubscribePause,
            OperateSubscribeContinue,
        },
        RiderLogTypeDeposit: {
            OperateDeposit,
        },
    }
)

// RiderLogReq 骑手操作日志
type RiderLogReq struct {
    PaginationReq

    ID   uint64 `json:"id" query:"id" validate:"required" trans:"骑手ID"`
    Type uint8  `json:"type" enums:"0,1,2,3,4" query:"type"` // 操作类别 0:全部 1:状态 2:资料 3:时长 4:押金
}

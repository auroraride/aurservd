// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "errors"
    "fmt"
)

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
    ID                uint64             `json:"id"`
    Phone             string             `json:"phone"`                       // 电话
    Name              string             `json:"name"`                        // 姓名, 实名认证后才会有
    Token             string             `json:"token,omitempty"`             // 认证token
    IsNewDevice       bool               `json:"isNewDevice"`                 // 是否新设备
    IsAuthed          bool               `json:"isAuthed"`                    // 是否已认证
    IsContactFilled   bool               `json:"isContactFilled"`             // 联系人是否添加
    Contact           *RiderContact      `json:"contact,omitempty"`           // 联系人
    Qrcode            string             `json:"qrcode"`                      // 二维码
    Deposit           float64            `json:"deposit"`                     // 需缴押金
    OrderNotActived   *bool              `json:"orderNotActived,omitempty"`   // 是否存在未激活订单
    Subscribe         *Subscribe         `json:"subscribe,omitempty"`         // 骑士卡
    Enterprise        *Enterprise        `json:"enterprise,omitempty"`        // 所属企业
    UseStore          bool               `json:"useStore"`                    // 是否可使用门店办理业务
    EnterpriseContact *EnterpriseContact `json:"enterpriseContact,omitempty"` // 团签联系方式
}

// RiderContact 紧急联系人
type RiderContact struct {
    Name     string `json:"name" validate:"required" trans:"联系人姓名"`
    Phone    string `json:"phone" validate:"required,phone" trans:"联系人电话"`
    Relation string `json:"relation" validate:"required" trans:"关系"`
}

func (c *RiderContact) String() string {
    return fmt.Sprintf("[%s]%s - %s", c.Relation, c.Name, c.Phone)
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
    RiderListFilter
}

type RiderListExport struct {
    RiderListFilter
    Remark string `json:"remark" validate:"required" trans:"备注"`
}

type RiderListFilter struct {
    Keyword         *string           `json:"keyword,omitempty" query:"keyword"`                                           // 搜索关键词
    Modified        *bool             `json:"modified,omitempty" query:"modified"`                                         // 是否被修改过
    Start           *string           `json:"start,omitempty" query:"start"`                                               // 注册开始时间, 格式为: 2022-01-01
    End             *string           `json:"end,omitempty" query:"end"`                                                   // 注册结束时间, 格式为: 2022-01-01
    Enterprise      *uint8            `json:"enterprise,omitempty" query:"enterprise"`                                     // 是否团签, 0:全部 1:团签 2:个签
    EnterpriseID    *uint64           `json:"enterpriseId,omitempty" query:"enterpriseId"`                                 // 团签企业ID, `enterprise = 1`时才会生效
    CityID          *uint64           `json:"cityId,omitempty" query:"cityId"`                                             // 城市筛选
    Model           string            `json:"model,omitempty" query:"model"`                                               // 电池型号筛选
    Status          *uint8            `json:"status,omitempty" enums:"0,1,2,3,4" query:"status"`                           // 用户状态 1:正常 2:已禁用 3:黑名单
    SubscribeStatus *uint8            `json:"subscribeStatus,omitempty" enums:"0,1,2,3,4,5,11,99" query:"subscribeStatus"` // 业务状态 0:未激活 1:计费中 2:寄存中 3:已逾期 4:已退订 5:已取消 11: 即将到期 99:未使用
    AuthStatus      *PersonAuthStatus `json:"authStatus,omitempty" enums:"0,1,2,3" query:"authStatus"`                     // 认证状态 0:未认证 1:认证中 2:已认证 3:认证失败
    PlanID          *uint64           `json:"planId,omitempty" query:"planId"`                                             // 骑士卡
    Remaining       *string           `json:"remaining,omitempty" query:"remaining"`                                       // 订阅剩余天数区间, 逗号分隔, 例如 `0,7`
    Suspend         *bool             `json:"suspend,omitempty" query:"suspend"`                                           // 是否筛选暂停扣费中, 不携带此参数获取全部, 携带此参数`true`暂停中 `false`非暂停
}

// RiderItemSubscribe 骑手骑士卡简单信息
type RiderItemSubscribe struct {
    ID         uint64  `json:"id"`                // 订阅ID
    Status     uint8   `json:"status"`            // 订阅状态 0:未激活 1:计费中 2:寄存中 3:已逾期 4:已退订 5:已取消 11: 即将到期(计算状态) 当 status = 1 且 remaining <= 3 的时候是即将到期
    Remaining  int     `json:"remaining"`         // 剩余天数
    Model      string  `json:"model"`             // 骑士卡可用电池型号
    Suspend    bool    `json:"suspend"`           // 是否暂停中
    AgentEndAt string  `json:"agentEndAt"`        // 代理商处到期日期
    Formula    *string `json:"formula,omitempty"` // 订阅天数计算公式
}

// RiderItem 骑手信息
type RiderItem struct {
    ID         uint64           `json:"id"`
    Name       string           `json:"name"`               // 姓名
    Phone      string           `json:"phone"`              // 手机号
    Status     uint8            `json:"status"`             // 用户状态, 优先显示状态值大的 1:正常 2:已禁用 3:黑名单
    AuthStatus PersonAuthStatus `json:"authStatus"`         // 认证状态 0:未认证 1:认证中 2:已认证 3:认证失败
    Deposit    float64          `json:"deposit"`            // 押金
    Address    string           `json:"address"`            // 户籍地址
    DeletedAt  string           `json:"deleteAt,omitempty"` // 账户删除时间
    Remark     string           `json:"remark"`             // 账户备注
    Contract   string           `json:"contract,omitempty"` // 合同(有可能不存在)
    Points     int64            `json:"points"`             // 积分
    Balance    float64          `json:"balance"`            // 余额
    // 团签企业信息, 若无此字段则为个签用户
    Enterprise *Enterprise `json:"enterprise,omitempty"`
    // 当前有效订阅信息, 若无此字段则代表当前无有效订阅 (订阅 = 骑手骑士卡)
    Subscribe *RiderItemSubscribe `json:"subscribe,omitempty"`
    // 认证信息, 有可能不存在, 内部字段也有可能不存在
    Person *Person `json:"person,omitempty"`
    // 紧急联系人, 有可能不存在
    Contact *RiderContact `json:"contact,omitempty"`
    // 所在城市, 有可能不存在
    City *City `json:"city,omitempty"`
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

// RiderPhoneSearchReq 使用手机号查询骑手请求
type RiderPhoneSearchReq struct {
    Phone string `json:"phone" validate:"required" trans:"手机号" query:"phone"` // 精准匹配
}

type RiderExchangeReq struct {
    PaginationReq

    RiderID uint64 `json:"riderId" query:"RiderId" trans:"骑手ID"`
}

type RiderSelectionReq struct {
    Keyword *string `json:"keyword" query:"keyword"` // 筛选骑手姓名或电话
}

type RiderPermissionError struct {
    error
    ForceSignout bool // 是否强制退出
}

func NewRiderPermissionError(message string, params ...bool) *RiderPermissionError {
    err := &RiderPermissionError{
        error:        errors.New(message),
        ForceSignout: false,
    }
    if len(params) > 0 {
        err.ForceSignout = params[0]
    }
    return err
}

type RiderDepositRes struct {
    Deposit float64 `json:"deposit"`
}

type RiderFollowUpCreateReq struct {
    RiderID uint64 `json:"riderId" trans:"骑手ID" validate:"required"`
    Remark  string `json:"remark" trans:"跟进信息" validate:"required"`
}

type RiderFollowUpListReq struct {
    PaginationReq

    RiderID uint64 `json:"riderId" validate:"required" query:"riderId" trans:"骑手ID"`
}

type RiderFollowUpListRes struct {
    ID      uint64   `json:"id"`
    Manager Modifier `json:"manager"` // 管理员信息
    Remark  string   `json:"remark"`  // 跟进信息
    Time    string   `json:"time"`    // 跟进时间
}

type RiderAgentList struct {
    Keyword string `json:"keyword"` // 筛选关键词
    Status  uint8  `json:"status"`  // 状态 0:全部 1:未激活 2:计费中 3:已超期 4:已退租
    CityID  uint64 `json:"cityId"`  // 城市筛选

}

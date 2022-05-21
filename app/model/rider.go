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

// RiderContext TODO 骑手上下文
type RiderContext struct {
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
    Id              uint64        `json:"id"`
    Token           string        `json:"token"`
    IsNewDevice     bool          `json:"isNewDevice"`
    IsAuthed        bool          `json:"isAuthed"`        // 是否已认证
    IsContactFilled bool          `json:"isContactFilled"` // 联系人是否添加
    Contact         *RiderContact `json:"contact,omitempty"`
    Qrcode          string        `json:"qrcode"`
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
// TODO status 以下还未做
type RiderListReq struct {
    PaginationReq

    Keyword  *string `json:"keyword"`  // 搜索关键词
    Modified *bool   `json:"modified"` // 是否被修改过
    Start    *string `json:"start"`    // 注册开始时间, 格式为: 2022-01-01
    End      *string `json:"end"`      // 注册结束时间, 格式为: 2022-01-01

    Status *uint8  `json:"status"` // 用户状态 0:未使用 1:未开通 2:计费中 3:寄存中 4:已过期 5:暂停中 6:已退租 7:已欠费 8:未认证 9:未办理 10:即将到期 11:已禁用 12:黑名单
    PlanID *uint64 `json:"planId"` // 骑士卡
}

// RiderItem 骑手信息
type RiderItem struct {
    ID           uint64          `json:"id"`
    Name         string          `json:"name"`                 // 用户姓名
    Phone        string          `json:"phone"`                // 手机号
    Status       uint8           `json:"status"`               // 用户状态 0:未使用 1:未开通 2:计费中 3:寄存中 4:已过期 5:暂停中 6:已退租 7:已欠费 8:未认证 9:未办理 10:即将到期 11:已禁用 12:黑名单
    IDCardNumber string          `json:"idCardNumber"`         // 身份证
    Deposit      float64         `json:"deposit"`              // 押金
    Enterprise   *EnterpriseItem `json:"enterprise,omitempty"` // 团签企业信息, 若无此字段则为个签用户
    UserPlan     *UserPlanItem   `json:"userPlan,omitempty"`   // 当前有效骑士卡, 若无此字段则代表当前无有效骑士卡
}

// RiderBlockReq 封禁或解封骑手账号
type RiderBlockReq struct {
    ID    uint64 `json:"id" `   // 骑手ID
    Block bool   `json:"block"` // `true`封禁 `false`解封
}

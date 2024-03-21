// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-10
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type RiderMgrDepositReq struct {
	Amount float64 `json:"amount"` // 押金金额
	ID     uint64  `json:"id" validate:"required" trans:"骑手ID"`
}

type RiderMgrModifyReq struct {
	ID         uint64            `json:"id" validate:"required" trans:"骑手ID"`
	Phone      *string           `json:"phone"`                      // 手机号
	AuthStatus *PersonAuthStatus `json:"authStatus" enums:"0,1,2,3"` // 认证状态 0:未认证 1:认证中 2:已认证 3:认证失败
	Contact    *RiderContact     `json:"contact"`                    // 联系人

	IdCardNumber   *string `json:"idCardNumber"`   // 身份证号
	IdCardPortrait *string `json:"idCardPortrait"` // 身份证正面
	IdCardNational *string `json:"idCardNational"` // 身份证国徽面
}

type RiderEmployeeSearchRes struct {
	ID              uint64           `json:"id"`
	Name            string           `json:"name"`                 // 骑手姓名
	Status          uint8            `json:"status"`               // 用户状态, 优先显示状态值大的 1:正常 2:已禁用 3:黑名单
	AuthStatus      PersonAuthStatus `json:"authStatus"`           // 认证状态 0:未认证 1:认证中 2:已认证 3:认证失败
	Phone           string           `json:"phone"`                // 骑手电话
	Enterprise      *Enterprise      `json:"enterprise,omitempty"` // 团签企业, 团签骑手判定依据(非空是团签), 个签无此字段
	Plan            *Plan            `json:"plan,omitempty"`       // 骑行卡, 个签才有此字段, 团签无
	Overview        ExchangeOverview `json:"overview"`             // 换电预览
	SubscribeStatus uint8            `json:"subscribeStatus"`      // 骑手订阅状态
}

type RiderMgrBatchOperate struct {
	Operate uint8    `json:"operate" validate:"required" enums:"1,2,3,4,5,6" trans:"操作类型"` // 1:禁用 2:解除禁用 3:拉黑 3:解除拉黑 4:暂停计费 5:继续计费
	IDS     []uint64 `json:"ids" validate:"required,min=1"`                                // 批量处理骑手ID
}

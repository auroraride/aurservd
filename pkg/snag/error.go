// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package snag

type StatusCode int

const (
	StatusOK                  StatusCode = iota << 8 // 0x000 请求成功
	StatusBadRequest                                 // 0x100 请求失败
	StatusUnauthorized                               // 0x200 需要登录
	StatusForbidden                                  // 0x300 没有权限
	StatusNotFound                                   // 0x400 资源未获
	StatusInternalServerError                        // 0x500 未知错误
	StatusRequireAuth                                // 0x600 需要实名
	StatusLocked                                     // 0x700 需要验证 (更换设备需要人脸验证)
	StatusRequireContact                             // 0x800 需要联系人
	StatusRequestTimeout                             // 0x900 请求过期
	StatusRequireSign                                // 0x1000 需要签约
)

type Error struct {
	Code    StatusCode `json:"code"`
	Message string     `json:"message"`
	Data    any        `json:"data,omitempty"`
}

func (e Error) Error() (message string) {
	message = e.Message
	if message == "" {
		switch e.Code {
		case StatusBadRequest:
			message = "请求失败"
		case StatusUnauthorized:
			message = "需要登录"
		case StatusForbidden:
			message = "没有权限"
		case StatusNotFound:
			message = "未找到资源"
		case StatusInternalServerError:
			message = "未知错误"
		case StatusRequireAuth:
			message = "需要实名验证"
		case StatusLocked:
			message = "需要人脸验证"
		case StatusRequireContact:
			message = "需要补充紧急联系人"
		case StatusRequestTimeout:
			message = "请求过期"
		case StatusRequireSign:
			message = "需要签约"
		default:
			message = "请求失败"
		}
	}
	return message
}

func NewError(params ...any) *Error {
	out := &Error{
		Code: StatusBadRequest,
	}

	for _, param := range params {
		switch v := param.(type) {
		case string:
			out.Message = v
		case error:
			out.Message = v.Error()
		case StatusCode:
			out.Code = v
		case int:
			out.Code = StatusCode(v)
		default:
			out.Data = v
		}
	}

	return out
}

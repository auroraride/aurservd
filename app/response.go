// Copyright (C) liasica. 2022-present.
//
// Created at 2022-04-20
// Based on api by liasica, magicrolan@qq.com.

package app

import (
	"bytes"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/pkg/snag"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// processParam 处理参数
func (r *Response) processParam(param any) {
	switch v := param.(type) {
	case string:
		r.Message = v
	case int:
		r.Code = v
	case error:
		r.Message = v.Error()
		if r.Code == int(snag.StatusOK) {
			r.Code = int(snag.StatusBadRequest)
		}
	default:
		if v != nil {
			r.Data = param
		}
	}
}

// SendResponse 发送响应
func (c *BaseContext) SendResponse(params ...any) error {
	r := &Response{
		Message: "ok",
		Code:    int(snag.StatusOK),
	}

	for _, param := range params {
		r.processParam(param)
	}

	if r.Code == int(snag.StatusOK) && r.Data == nil {
		r.Data = model.StatusResponse{Status: true}
	}

	if reflect2.IsNil(r.Data) {
		r.Data = nil
	}

	buffer := &bytes.Buffer{}
	encoder := jsoniter.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(r)

	return c.JSONBlob(http.StatusOK, buffer.Bytes())
}

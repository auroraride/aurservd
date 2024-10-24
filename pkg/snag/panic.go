// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package snag

import (
	"fmt"

	"go.uber.org/zap"
)

func Panic(params ...any) {
	panic(NewError(params...))
}

func PanicCallback(cb func(), params ...any) {
	cb()
	panic(NewError(params...))
}

func PanicCallbackX(cb func() error, params ...any) {
	_ = cb()
	panic(NewError(params...))
}

func PanicIfError(err error, params ...any) {
	if err != nil {
		panic(NewError(err))
	}
}

func PanicIfErrorX(err error, cb func() error, params ...any) {
	if err != nil {
		_ = cb()
		panic(NewError(err))
	}
}

func WithPanic(cb func()) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
		}
	}()

	cb()

	return
}

func WithPanicStack(cb func()) {

	defer func() {
		if v := recover(); v != nil {
			zap.L().Error("[WithPanicStack] 捕获未处理崩溃: %v\n%s", zap.Error(fmt.Errorf("%v", v)), zap.Stack("stack"))
		}
	}()

	cb()
}

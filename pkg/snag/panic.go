// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package snag

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

func PanicIfError(err error, cb func(), params ...any) {
    if err != nil {
        cb()
        panic(NewError(err))
    }
}

func PanicIfErrorX(err error, cb func() error, params ...any) {
    if err != nil {
        _ = cb()
        panic(NewError(err))
    }
}

// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package snag

func Panic(e interface{}) {
    panic(NewError(e))
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-25
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "fmt"
    "os"
)

func FmtFatalln(err any) {
    fmt.Println(err)
    os.Exit(0)
}

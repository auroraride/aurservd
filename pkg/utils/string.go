// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/10
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import "unicode"

func StrToFirstUpper(str string) string {
    if len(str) == 0 {
        return ""
    }
    tmp := []rune(str)
    tmp[0] = unicode.ToUpper(tmp[0])
    return string(tmp)
}

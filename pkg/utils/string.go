// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/10
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z\\d])([A-Z])")
)

func StrToFirstUpper(str string) string {
	if len(str) == 0 {
		return ""
	}
	tmp := []rune(str)
	tmp[0] = unicode.ToUpper(tmp[0])
	return string(tmp)
}

func StrToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

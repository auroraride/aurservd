// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import "regexp"

const (
	// RegexPhone `^((13[0-9])|(14[5-9])|(15([0-3]|[5-9]))|(16[6-7])|(17[1-8])|(18[0-9])|(19[1|3])|(19[5|6])|(19[8|9]))\d{8}$`
	RegexPhone        = `(?m)^((13[0-9])|(14[5-9])|(15([0-3]|[5-9]))|(16[6-7])|(17[1-8])|(18[0-9])|(19[0-3])|(19[5|6])|(19[8|9]))\d{8}$`
	RegexIDCardNumber = `^([1-6][1-9]|50)\d{4}(18|19|20)\d{2}((0[1-9])|10|11|12)(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`
)

type regex struct {
}

func NewRegex() *regex {
	return &regex{}
}

func (*regex) MatchPhone(str string) bool {
	return regexp.MustCompile(RegexPhone).MatchString(str)
}

func (r *regex) MatchIDCardNumber(str string) bool {
	return regexp.MustCompile(RegexIDCardNumber).MatchString(str)
}

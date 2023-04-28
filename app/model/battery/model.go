// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package battery

import "regexp"

type Model struct {
	vaild    bool
	Voltage  string `json:"voltage"`  // 电压, 单位: 伏特(V)
	Capacity string `json:"capacity"` // 容量, 单位: 安时(AH)
}

func New(str string) (m *Model) {
	matchs := regexp.MustCompile(`(?m)(\d+)?[V|v](\d+)[A|a][H|h]`).FindStringSubmatch(str)
	m = &Model{}
	if len(matchs) == 3 {
		m.vaild = true
		m.Voltage = matchs[1]
		m.Capacity = matchs[2]
	}
	return
}

func (m *Model) IsVaild() bool {
	return m.vaild
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-29
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

type CabinetHealth struct {
    Brand  string `json:"brand" sls:"品牌" index:"doc"`
    Serial string `json:"serial" sls:"编码" index:"doc"`
}

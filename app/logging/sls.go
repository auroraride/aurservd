// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
	"fmt"
	"reflect"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/golang-module/carbon/v2"
)

// GenerateLogContent 转换为sls日志
func GenerateLogContent(pointer any) (contents []*sls.LogContent) {
	t := reflect.TypeOf(pointer).Elem()
	n := t.NumField()
	value := reflect.ValueOf(pointer).Elem()

	contents = make([]*sls.LogContent, n)
	for i := 0; i < n; i++ {
		// tag, _ := t.Field(i).Tag.Lookup("sls") // 使用中文字段名
		json, _ := t.Field(i).Tag.Lookup("json") // 使用英文字段名
		v := value.Field(i)
		cv := ""
		if v.Type().Kind() == reflect.Bool {
			cv = "否"
			if v.Bool() {
				cv = "是"
			}
		} else {
			cv = fmt.Sprintf("%v", v.Interface())
		}
		if json == "time" && v.Type().Kind() == reflect.String && v.String() == "" {
			cv = time.Now().Format(carbon.DateTimeLayout)
		}
		contents[i] = &sls.LogContent{
			Key:   tea.String(json),
			Value: tea.String(cv),
		}
	}

	return
}

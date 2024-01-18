// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on ard by liasica, magicrolan@qq.com.

package main

import (
	"github.com/auroraride/aurservd/cmd/aurservd/internal/script"
	"github.com/auroraride/aurservd/pkg/snag"
)

func main() {

	snag.WithPanicStack(func() {
		// 启动脚本
		script.Execute()
	})
}

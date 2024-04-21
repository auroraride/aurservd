// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-05
// Based on aurservd by liasica, magicrolan@qq.com.

package rpc

import "github.com/auroraride/aurservd/internal/ar"

func serverAddress(key string) string {
	return ar.Config.RpcServer[key]
}

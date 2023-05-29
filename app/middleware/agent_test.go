// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
	"strings"
	"testing"

	"github.com/auroraride/aurservd/app"
)

const testagenttoken = "Bearer AGENT/A6X5xSRx3JprjDPRe-3OSkm7T7v8jM5N_U7Cmnn5816B/MHcCAQEEIJQIzKgg_feSHkuVBxni-IDJvk0-g431RcoNBP5JgGt-oAoGCCqGSM49AwEHoUQDQgAEpfnFJHHcmmuMM9F77c5KSbtPu_yMzk39TsKaefnzXoF_47xavaity_ffbRmzXMBGJyGLMOUFldez1cDJVAnqzQ"

func BenchmarkAgentToken1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strings.Split(testagenttoken, app.AgentBearer)
	}
}

func BenchmarkAgentToken2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.TrimLeft(testagenttoken, app.AgentBearer)
	}
}

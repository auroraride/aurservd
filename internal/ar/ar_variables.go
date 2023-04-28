// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import (
	"time"

	"github.com/go-redis/redis/v9"
)

var (
	TimeLocation        *time.Location
	Quit                chan bool
	Redis               *redis.Client
	CabinetNameCacheKey string
	TaskCacheKey        string
)

func init() {
	Quit = make(chan bool)
}

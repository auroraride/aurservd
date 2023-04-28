// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-25
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Endpoint        string
	RegionId        string
	LogPath         string
}

var (
	cfg        *config
	configFile string
)

func ParseConfig() {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg = new(config)
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

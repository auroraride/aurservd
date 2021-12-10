// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package ar


import (
    _ "embed"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/spf13/viper"
    "log"
    "os"
)

const (
    configFile = "./config/config.yaml"
)

// defaultConfigStr 默认配置
//go:embed default_config.yml
var defaultConfigStr string

type config struct {
    App struct {
        Address   string
        BodyLimit string
        RateLimit float64
    }
    Nsq struct {
        Url string
    }
    Database struct {
        Postgres struct {
            Dsn string
        }
        Mongo struct {
            Uri string
            Db  string
        }
    }
    ThirdParty struct {
        Yunfuture struct {
            Token string
        }
    } `mapstructure:"third_party"`
    Logging struct {
        Color bool   // 是否启用日志颜色
        Level string // 日志等级
        Age   int    // 日志保存时间（小时）
        Json  bool   // 日志以json格式保存
    }
}

func LoadConfig() {
    // 判断配置是否存在
    f := utils.NewFile(configFile)
    if !f.IsExist() {
        err := f.CreateDirectoryIfNotExist()
        if err != nil {
            log.Fatalf("配置目录创建失败: %v", err)
            return
        }
        err = os.WriteFile(configFile, []byte(defaultConfigStr), 0644)
        if err != nil {
            log.Fatalf("默认配置保存失败: %v", err)
            return
        }
    }

    viper.SetConfigFile(configFile)
    viper.AutomaticEnv()
    // 读取配置
    err := viper.ReadInConfig()
    if err != nil {
        log.Fatalf("配置读取失败: %v", err)
    }
    err = viper.Unmarshal(Config)
    if err != nil {
        log.Fatalf("配置解析失败: %v", err)
    }
}
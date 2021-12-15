// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import (
    _ "embed"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/fsnotify/fsnotify"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
    "os"
)

const (
    configFile = "./config/config.yaml"
)

var Config *config

// defaultConfigStr 默认配置
//go:embed default_config.yml
var defaultConfigStr string

type EsignConfig struct {
    Appid    string
    BaseUrl  string
    Secret   string
    RSA      string `mapstructure:"rsa"`
    Log      bool
    Callback string
    Redirect string
    Group    struct {
        FlowId     string
        TemplateId string
    }
    Person struct {
        FlowId     string
        TemplateId string
    }
}

type config struct {
    App struct {
        Address   string
        BodyLimit string
        RateLimit float64
        Captcha   struct {
            Names map[string]string
        }
    }
    Nsq struct {
        Url string
    }
    Database struct {
        Postgres struct {
            Dsn string
        }
        Redis struct {
            Addr     string
            Password string
            DB       int `mapstructure:"db"`
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
    Aliyun struct {
        Oss struct {
            AccessKeyId     string
            AccessKeySecret string
            Bucket          string
            Endpoint        string
            Url             string
        }
        Sms struct {
            AccessId     string
            AccessSecret string
            Endpoint     string
            Sign         string
            Template     struct {
                General struct {
                    Code string
                }
            }
        }
    }
    Baidu struct {
        Face struct {
            ApiKey     string
            SecretKey  string
            Callback   string
            AuthPlanId string
            FacePlanId string
        }
    }
    Esign struct {
        Target  string
        Sandbox EsignConfig
        Online  EsignConfig
    }
    Trans map[string]string
}

func readConfig() error {
    viper.SetConfigFile(configFile)
    viper.AutomaticEnv()
    // 读取配置
    err := viper.ReadInConfig()
    if err != nil {
        return err
    }
    Config = new(config)
    err = viper.Unmarshal(Config)
    return err
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

    err := readConfig()
    if err != nil {
        log.Fatalf("配置读取失败: %v", err)
    }

    viper.OnConfigChange(func(e fsnotify.Event) {
        log.Infof("配置已改动: %s, 重载配置: %v", e.Name, readConfig())
    })
    viper.WatchConfig()
}

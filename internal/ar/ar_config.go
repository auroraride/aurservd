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

type Version struct {
    Version     string `json:"version"`
    Description string `json:"description"`
    Link        string `json:"link"`
    Enable      bool   `json:"enable"`
}

type AppVersion struct {
    Android Version `json:"android"`
    IOS     Version `json:"ios" mapstructure:"ios"`
}

type EsignConfig struct {
    Appid    string
    BaseUrl  string
    Secret   string
    RSA      string `mapstructure:"rsa"`
    Log      bool
    Callback string
    Redirect string
    Group    struct {
        Scene      string
        TemplateId string
    }
    Person struct {
        Scene      string
        TemplateId string
    }
}

type config struct {
    App struct {
        Address      string
        Host         string
        Mode         string
        SQL          bool
        BodyLimit    string
        RateLimit    float64
        CabinetDebug bool
        Task         bool
        Debug        struct {
            Phone map[string]bool
        }
        Captcha struct {
            Names map[string]string
        }
    }
    Cabinet struct {
        Debug    bool
        Provider bool
    }
    Nsq struct {
        Url string
    }
    Amap struct {
        Key string
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
        Yundong struct {
            Appid  string
            Appkey string
            Url    string
        }
        Kaixin struct {
            Url string
            Key string
        }
    } `mapstructure:"third_party"`
    Logging struct {
        Color    bool   // 是否启用日志颜色
        Level    string // 日志等级
        Age      int    // 日志保存时间（小时）
        Json     bool   // 日志以json格式保存
        RootPath string // 移除根目录
    }
    Aliyun struct {
        Oss struct {
            AccessKeyId     string
            AccessKeySecret string
            Bucket          string
            Endpoint        string
            Url             string
            Arn             string
            RamRole         string
            RegionId        string
        }
        // 语音通知
        Vms struct {
            AccessKeyId     string
            AccessKeySecret string
            Endpoint        string
            Overdue         struct {
                Template *string // tts 模板ID
                Tel      *string // 独立号码为空的时候使用公共号码进行发送
            }
        }
        Sms struct {
            AccessKeyId     string
            AccessKeySecret string
            Endpoint        string
            Sign            string
            Template        struct {
                General struct {
                    Code string
                }
            }
        }
        Sls struct {
            AccessKeyId     string
            AccessKeySecret string
            Endpoint        string
            Project         string
            CabinetLog      string // 电柜日志logstore
            DoorLog         string // 柜门操作日志logstore
            OperateLog      string // 管理端操作日志logstore
            ExchangeLog     string // 换电日志logstore
            HealthLog       string // 电柜在线变化日志logstore
            BatteryLog      string // 电柜电池变化日志logstore
        }
    }
    Baidu struct {
        Face struct {
            ApiKey     string
            SecretKey  string
            Redirect   string
            AuthPlanId string
            FacePlanId string
        }
    }
    Esign struct {
        Target  string
        Sandbox EsignConfig
        Online  EsignConfig
    }
    Mob struct {
        Push struct {
            Env       string
            AppKey    string
            AppSecret string
        }
    }
    Trans   map[string]string
    Payment struct {
        Wechat struct {
            PrivateKeyPath             string
            MchID                      string
            AppID                      string
            MchCertificateSerialNumber string
            MchAPIv3Key                string
            NotifyUrl                  string
            RefundUrl                  string
        }
        Alipay struct {
            Appid         string
            PrivateKey    string
            AppPublicCert string
            RootCert      string
            PublicCert    string
            NotifyUrl     string
        }
    }
    RiderApp    AppVersion
    EmployeeApp AppVersion
    WxWork      struct {
        AgentID    int64  `mapstructure:"agentID"`
        CorpID     string `mapstructure:"corpID"`
        CorpSecret string `mapstructure:"corpSecret"`
    }
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

// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import (
    _ "embed"
    "github.com/auroraride/adapter"
    "log"
)

const (
    configFile = "./config/config.yaml"
)

var (
    Config *config
)

type Version struct {
    Version     string `json:"version"`
    Description string `json:"description"`
    Link        string `json:"link"`
    Enable      bool   `json:"enable"`
}

type AppVersion struct {
    Android Version `json:"android"`
    IOS     Version `json:"ios" koanf:"ios"`
}

type EsignConfig struct {
    Appid    string
    BaseUrl  string
    Secret   string
    RSA      string `koanf:"rsa"`
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
    Debug bool
    adapter.Configure

    Tcp struct {
        Bind struct {
            Yundong string
        }
    }

    App struct {
        Host         string
        Mode         string
        SQL          bool
        CabinetDebug bool
        Debug        struct {
            Phone map[string]bool
        }
        Captcha struct {
            Names map[string]string
        }
    }

    Sync struct {
        Kxcab struct {
            Api string
        }
        Kxnicab struct {
            Api string
        }
        Ydcab struct {
            Api string
        }
        Tbcab struct {
            Api string
        }
    }

    RpcServer map[string]string `koanf:"rpc-server"`

    Task struct {
        Branch     bool // 网点合同到期提醒
        Enterprise bool // 团签账单
        Sim        bool // SIM卡到期提醒
        Subscribe  bool // 个签订阅日期计算
        Reserve    bool // 预约到期计算
        Reminder   bool // 个签到期提醒
        Cabinet    bool // 电柜任务失效
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
    }

    ThirdParty struct {
        Kaixin struct {
            Url string
            Key string
        }
    } `koanf:"third_party"`

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
            // 语音催费配置
            Reminder struct {
                Template *string // tts 模板ID
                Tel      *string // 独立号码为空的时候使用公共号码进行发送
            }
            // 逾期通知
            Overdue struct {
                Template *string
                Tel      *string
            }
        }
        Sms struct {
            AccessKeyId     string
            AccessKeySecret string
            Endpoint        string
            Sign            string
            Template        struct {
                // 验证码
                General string
                // 短信催费
                Reminder string
                // 逾期通知
                Overdue string
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
        Map struct {
            ApiKey string
            Ak     string
            Sk     string
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

    Trans map[string]string

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
        AgentID    int64  `koanf:"agentID"`
        CorpID     string `koanf:"corpID"`
        CorpSecret string `koanf:"corpSecret"`
    }
}

func LoadConfig() {
    var err error

    Config = new(config)
    err = adapter.LoadConfigure(Config, configFile, nil)
    if err != nil {
        log.Fatal(err)
    }

    Config.setKeys()
}

func (c *config) setKeys() {
    CabinetNameCacheKey = c.GetCacheKey("CABINET_NAME")
    TaskCacheKey = c.GetCacheKey("TASK")
}

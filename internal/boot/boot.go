// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package boot

import (
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/auroraride/adapter/log"
	"github.com/go-redis/redis/v9"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/logging"
	_ "github.com/auroraride/aurservd/app/permission"
	"github.com/auroraride/aurservd/assets"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/payment"
)

func Bootstrap() {
	// 设置全局时区
	tz := "Asia/Shanghai"
	_ = os.Setenv("TZ", tz)
	loc, _ := time.LoadLocation(tz)
	ar.TimeLocation = loc
	time.Local = loc
	carbon.SetTimezone(tz)

	// 载入配置
	ar.LoadConfig()

	// 创建redis客户端
	ar.Redis = redis.NewClient(&redis.Options{
		Addr:     ar.Config.Redis.Address,
		Password: ar.Config.Redis.Password,
		DB:       ar.Config.Redis.DB,
	})

	// 初始化日志
	logcfg := &log.Config{
		FormatJson: !ar.Config.LoggerDebug,
		Stdout:     ar.Config.LoggerDebug,
		LoggerName: ar.Config.LoggerName,
		NoCaller:   true,
	}

	if !ar.Config.LoggerDebug {
		logcfg.Writers = append(logcfg.Writers, log.NewRedisWriter(redis.NewClient(&redis.Options{})))
	}

	if len(ar.Config.Kafka.Brokers) > 0 {
		config := sarama.NewConfig()                     // 设置日志输入到Kafka的配置
		config.Producer.RequiredAcks = sarama.WaitForAll // 等待服务器所有副本都保存成功后的响应
		config.Producer.Return.Successes = true          // 是否等待成功后的响应,只有上面的RequiredAcks设置不是NoReponse这里才有用.
		config.Producer.Return.Errors = true
		producer, err := sarama.NewSyncProducer(ar.Config.Kafka.Brokers, config)
		if err != nil {
			return
		}

		kw := log.NewKafkaWriter(&log.KafkaWriter{
			Topic:    ar.Config.Kafka.Topic,
			Producer: producer,
		})

		logcfg.Writers = append(logcfg.Writers, log.NewKafkaWriter(kw))
	}

	log.New(logcfg)

	// 加载数据库
	entInit()

	// 加载其他数据

	// 初始化日志
	logging.Boot()

	// 初始化支付
	payment.Boot()

	// 加载模板
	assets.LoadTemplates()
}

// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/internal/app/response"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
    mw "github.com/labstack/echo/v4/middleware"
    log "github.com/sirupsen/logrus"
    "golang.org/x/time/rate"
)

type router struct {
    *echo.Echo
}

var r *router

func Run() {
    r = &router{echo.New()}
    cfg := ar.Config.App
    // 加载全局中间件
    corsConfig := mw.DefaultCORSConfig
    corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, []string{
        ar.HeaderCaptchaID,
    }...)
    r.Use(
        mw.Recover(),
        mw.BodyLimit(cfg.BodyLimit),
        mw.CORSWithConfig(corsConfig),
        mw.GzipWithConfig(mw.GzipConfig{
            Level: 5,
        }),
        mw.RequestID(),
        mw.RateLimiter(mw.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit))),
        mw.LoggerWithConfig(mw.LoggerConfig{
            Format: `{"time":"${time_custom}","id":"${id}","remote_ip":"${remote_ip}",` +
                `"host":"${host}","method":"${method}","uri":"${uri}",` +
                `"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
                `,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
            CustomTimeFormat: "2006-01-02 15:04:05.00000",
        }),
    )

    r.HTTPErrorHandler = func(err error, c echo.Context) {
        _ = response.New(c).Error(response.StatusError).SetMessage(err.Error()).Send()
    }

    // 载入路由
    // 公共API
    r.commonRoute()

    // 骑手api
    r.rideRoute()

    log.Fatal(r.Start(cfg.Address))
}

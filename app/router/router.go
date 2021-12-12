// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/request"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
    mw "github.com/labstack/echo/v4/middleware"
    log "github.com/sirupsen/logrus"
    "golang.org/x/time/rate"
)

var r *router

type router struct {
    *echo.Echo
}

func Run() {
    e := echo.New()
    // e.Logger.SetHeader(`[time] ${time_rfc3339_nano}` + "\n")
    r = &router{e}
    cfg := ar.Config.App
    corsConfig := mw.DefaultCORSConfig
    corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, []string{
        app.HeaderCaptchaID,
        app.HeaderDeviceSerial,
        app.HeaderDeviceType,
    }...)
    // 加载全局中间件
    r.Use(
        func(next echo.HandlerFunc) echo.HandlerFunc {
            return func(ctx echo.Context) error {
                c := &app.Context{Context: ctx}
                return next(c)
            }
        },
        mw.LoggerWithConfig(mw.LoggerConfig{
            Format: `{"time":"${time_custom}","id":"${id}","remote_ip":"${remote_ip}",` +
                `"host":"${host}","method":"${method}","uri":"${uri}",` +
                `"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
                `,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
            CustomTimeFormat: "2006-01-02 15:04:05.00000",
        }),
        mw.Recover(),
        mw.BodyLimit(cfg.BodyLimit),
        mw.CORSWithConfig(corsConfig),
        mw.GzipWithConfig(mw.GzipConfig{
            Level: 5,
        }),
        mw.RequestID(),
        mw.RateLimiter(mw.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit))),
    )

    r.Validator = request.NewGlobalValidator()

    r.HTTPErrorHandler = func(err error, c echo.Context) {
        code := response.StatusError
        errCode, ok := c.Get("errCode").(int)
        if ok {
            code = errCode
        }
        _ = response.New(c).Error(code).SetMessage(err.Error()).Send()
    }

    // 载入路由
    // 公共API
    r.commonRoute()

    // 骑手api
    r.rideRoute()

    log.Fatal(r.Start(cfg.Address))
}

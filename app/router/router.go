// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/middleware"
    "github.com/auroraride/aurservd/app/request"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
    mw "github.com/labstack/echo/v4/middleware"
    log "github.com/sirupsen/logrus"
    "golang.org/x/time/rate"
)

var (
    root *echo.Group
)

func Run() {
    e := echo.New()
    root = e.Group("/")

    // 错误处理
    e.Validator = request.NewGlobalValidator()
    e.HTTPErrorHandler = func(err error, c echo.Context) {
        res := app.NewResponse(c).Error(app.StatusError).SetMessage(err.Error())
        if e, ok := err.(*snag.Error); ok {
            if data, ok := e.Data.(app.Response); ok {
                res.Error(data.Code).SetMessage(data.Message).SetData(data.Data)
            }
        }
        _ = res.Send()
    }

    // 先载入文档路由
    e.Use(newRedoc())

    // e.Logger.SetHeader(`[time] ${time_rfc3339_nano}` + "\n")
    cfg := ar.Config.App
    corsConfig := mw.DefaultCORSConfig
    corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, []string{
        app.HeaderCaptchaID,
        app.HeaderDeviceSerial,
        app.HeaderDeviceType,
    }...)
    // 加载全局中间件
    root.Use(
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
        // mw.Recover(),
        middleware.Recover(),
        mw.BodyLimit(cfg.BodyLimit),
        mw.CORSWithConfig(corsConfig),
        mw.GzipWithConfig(mw.GzipConfig{
            Level: 5,
        }),
        mw.RequestID(),
        mw.RateLimiter(mw.NewRateLimiterMemoryStore(rate.Limit(cfg.RateLimit))),
    )

    // 载入路由
    loadCommonRoutes()  // 公共API
    loadRideRoutes()    // 骑手路由
    loadManagerRoutes() // 管理员路由

    log.Fatal(e.Start(cfg.Address))
}

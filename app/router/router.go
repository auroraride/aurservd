// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "fmt"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/controller"
    "github.com/auroraride/aurservd/app/middleware"
    "github.com/auroraride/aurservd/app/request"
    "github.com/auroraride/aurservd/assets"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/goccy/go-json"
    "github.com/labstack/echo/v4"
    mw "github.com/labstack/echo/v4/middleware"
    log "github.com/sirupsen/logrus"
    "golang.org/x/time/rate"
    "net/http"
)

var (
    root *echo.Group
    e    *echo.Echo
)

// DefaultJSONSerializer implements JSON encoding using encoding/json.
type DefaultJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d DefaultJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
    enc := json.NewEncoder(c.Response())
    if indent != "" {
        enc.SetIndent("", indent)
    }
    return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d DefaultJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
    err := json.NewDecoder(c.Request().Body).Decode(i)
    if ute, ok := err.(*json.UnmarshalTypeError); ok {
        return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
    } else if se, ok := err.(*json.SyntaxError); ok {
        return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
    }
    return err
}

func Run() {
    e = echo.New()
    e.JSONSerializer = &DefaultJSONSerializer{}

    e.Renderer = assets.Templates

    e.Static("/pages", "public/pages")

    root = e.Group("/")

    // 校验规则
    e.Validator = request.NewGlobalValidator()

    // 错误处理
    e.HTTPErrorHandler = func(err error, c echo.Context) {
        ctx := app.Context(c)
        if ctx == nil {
            ctx = app.NewContext(c)
        }
        message := err.Error()
        code := int(snag.StatusBadRequest)
        var data any
        switch err.(type) {
        case *snag.Error:
            target := err.(*snag.Error)
            code = int(target.Code)
            data = target.Data
        case *echo.HTTPError:
            target := err.(*echo.HTTPError)
            message = fmt.Sprintf("%v", target.Message)
            switch target.Code {
            case http.StatusNotFound:
                code = int(snag.StatusNotFound)
                break
            default:
                code = int(snag.StatusBadRequest)
                break
            }
            break
        }
        _ = ctx.SendResponse(code, message, data)
    }

    // e.Logger.SetHeader(`[time] ${time_rfc3339_nano}` + "\n")
    cfg := ar.Config.App
    corsConfig := mw.DefaultCORSConfig
    corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, []string{
        app.HeaderContentType,
        app.HeaderCaptchaID,
        app.HeaderDeviceSerial,
        app.HeaderDeviceType,
        app.HeaderPushId,
        app.HeaderRiderToken,
        app.HeaderManagerToken,
        app.HeaderEmployeeToken,
    }...)
    corsConfig.ExposeHeaders = append(corsConfig.ExposeHeaders, []string{
        app.HeaderCaptchaID,
        app.HeaderContentType,
        app.HeaderDispositionType,
    }...)
    // 加载全局中间件
    root.Use(
        // AppContext
        func(next echo.HandlerFunc) echo.HandlerFunc {
            return func(ctx echo.Context) error {
                return next(app.NewContext(ctx))
            }
        },
        // mw.LoggerWithConfig(mw.LoggerConfig{
        //     Format: `{"time":"${time_custom}","id":"${id}","remote_ip":"${remote_ip}",` +
        //         `"method":"${method}","uri":"${uri}",` +
        //         `"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
        //         `,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
        //     CustomTimeFormat: "2006-01-02 15:04:05.00000",
        // }),
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
    root.GET("app/version", controller.Version.Get)

    loadDocRoutes()      // 文档
    loadCabinetRoutes()  // 电柜回调
    loadCommonRoutes()   // 公共API
    loadRideRoutes()     // 骑手路由
    loadManagerRoutes()  // 管理员路由
    loadEmployeeRoutes() // 门店端路由
    loadAgentRoutes()    // 代理端路由
    loadToolRoutes()

    log.Fatal(e.Start(cfg.Address))
}

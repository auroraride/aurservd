// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-10
// Based on aurservd by liasica, magicrolan@qq.com.

package v2

import (
	"strconv"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/app"
	"github.com/labstack/echo/v4"

	inapp "github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/controller/v2/rapi"
	"github.com/auroraride/aurservd/app/middleware"
	"github.com/auroraride/aurservd/internal/ent"
)

func LoadRiderV2Routes(root *echo.Group) {
	g := root.Group("rider/v2")

	rawDump := app.NewDumpLoggerMiddleware().WithConfig(&app.DumpConfig{
		RequestHeader:  true,
		ResponseHeader: true,
	})

	// 记录请求日志
	dumpSkipPaths := map[string]bool{}
	dumpReqHeaders := map[string]struct{}{
		inapp.HeaderCaptchaID:    {},
		inapp.HeaderDeviceSerial: {},
		inapp.HeaderDeviceType:   {},
		inapp.HeaderPushId:       {},
	}
	dump := app.NewDumpLoggerMiddleware().WithConfig(&app.DumpConfig{
		ResponseBodySkipper: func(c echo.Context) bool {
			return dumpSkipPaths[c.Path()]
		},
		RequestHeader: true,
		RequestHeaderSkipper: func(s string) bool {
			_, ok := dumpReqHeaders[s]
			return !ok
		},
		Extra: func(c echo.Context) []byte {
			if r, ok := c.Get("rider").(*ent.Rider); ok && r != nil {
				buf := adapter.NewBuffer()
				defer adapter.ReleaseBuffer(buf)

				buf.WriteString(`{"id":`)
				buf.WriteString(strconv.FormatUint(r.ID, 10))
				buf.WriteString(`,"phone":"`)
				buf.WriteString(r.Phone)
				buf.WriteString(`","name":"`)
				buf.WriteString(r.Name)
				buf.WriteString(`"}`)

				return buf.Bytes()
			}
			return nil
		},
	})

	g.Use(
		middleware.DeviceMiddleware(),
		middleware.RiderMiddlewareV2(),
		dump,
	)

	// 实名认证
	g.POST("/certification", rapi.Person.Certification)

	// 实名认证回调
	g.GET("/certification/callback", rapi.Person.CertificationCallback, rawDump)

	// 实名认证结果
	g.GET("/certification/result", rapi.Person.CertificationResult)
}

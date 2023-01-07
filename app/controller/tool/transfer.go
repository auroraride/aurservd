// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-07
// Based on aurservd by liasica, magicrolan@qq.com.

package tool

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/labstack/echo/v4"
    "net/http"
)

type transfer struct{}

var Transfer = new(transfer)

func (*transfer) Subscribe(c echo.Context) (err error) {
    var (
        message string
        formula string
        state   bool
    )

    if c.Request().Method == "POST" {
        phone := c.FormValue("phone")
        intelligent := c.FormValue("intelligent") == "1"

        ctx := context.Background()
        sub, _ := ent.Database.Subscribe.Query().Where(
            subscribe.HasRiderWith(rider.Phone(phone)),
            subscribe.StatusIn(model.SubscribeNotUnSubscribed()...),
        ).First(ctx)

        if sub == nil {
            message = "未找到有效订阅信息"
            goto RENDER
        }

        if !intelligent && (sub.BatterySn != nil || sub.BatteryID != nil) {
            message = "当前骑手已绑定电池, 无法转为非智能套餐"
            goto RENDER
        }

        sub, err = sub.Update().SetIntelligent(intelligent).Save(ctx)
        if err != nil {
            message = err.Error()
            goto RENDER
        }

        message = "成功"
        state = true

        if sub.Formula != nil {
            formula = *sub.Formula
        }
    }

RENDER:
    return c.Render(http.StatusOK, "transfer/subscribe.html", ar.Map{
        "message": message,
        "formula": formula,
        "state":   state,
    })
}

// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-07
// Based on aurservd by liasica, magicrolan@qq.com.

package tool

import (
    "context"
    "github.com/auroraride/adapter"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/silk"
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
        bsn := c.FormValue("battery")

        ctx := context.Background()
        sub, _ := ent.Database.Subscribe.Query().Where(
            subscribe.HasRiderWith(rider.Phone(phone)),
            subscribe.StatusIn(model.SubscribeNotUnSubscribed()...),
        ).First(ctx)

        if sub == nil {
            message = "未找到有效订阅信息"
            goto RENDER
        }

        if exists, _ := sub.QueryBattery().Exist(ctx); exists {
            message = "当前骑手已绑定电池, 无法转化"
            goto RENDER
        }

        var (
            bat *ent.Battery
            bid *uint64

            bm = sub.Model
        )

        // 如果是智能电池, 解析并查找电池信息
        if intelligent {
            ab := adapter.ParseBatterySN(bsn)
            if ab.Model == "" {
                message = "电池编号解析错误"
                goto RENDER
            }

            // 查找电池
            bat, _ = service.NewBattery().QuerySn(ab.SN)
            if bat == nil {
                message = "未找到电池信息"
                goto RENDER
            }

            // 设置电池信息
            bm = ab.Model
            bid = silk.UInt64(bat.ID)
        }

        err = ent.WithTx(ctx, func(tx *ent.Tx) (err error) {
            err = tx.Subscribe.UpdateOneID(sub.ID).
                SetIntelligent(intelligent).
                SetNillableBatteryID(bid).
                SetModel(bm).
                Exec(ctx)
            if err != nil {
                return
            }

            // 更新电池信息
            if intelligent && bat != nil {
                err = service.NewBattery().Allocate(tx.Battery.UpdateOne(bat), bat, sub, false)
            }

            return
        })

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

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/gorilla/websocket"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
    "net/http"
)

type socket struct{}

var Socket = new(socket)

func (*socket) Employee(c echo.Context) error {
    var upgrader = websocket.Upgrader{}
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

    if err != nil {
        log.Error(err)
        return err
    }

    srv := service.NewEmployeeSocket()
    key := c.Request().Header.Get("Sec-WebSocket-Key")

    defer func(ws *websocket.Conn) {
        log.Infof("%s disconnect", key)
        srv.DisConnect(ws)
        _ = ws.Close()
    }(ws)

    token := c.QueryParam("token")
    register := srv.Connect(ws, token)
    _ = ws.WriteMessage(websocket.TextMessage, register.Bytes())

    for {
        _, _, err = ws.ReadMessage()
        if err != nil {
            break
        }

        _ = ws.WriteMessage(websocket.PongMessage, nil)
    }

    return nil
}

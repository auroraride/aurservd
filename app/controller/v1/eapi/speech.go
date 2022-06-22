// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-22
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/gorilla/websocket"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
    "net/http"
)

type speech struct{}

var Speech = new(speech)

var upgrader = websocket.Upgrader{}

func (*speech) Store(c echo.Context) error {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

    if err != nil {
        log.Error(err)
        return err
    }

    srv := service.NewSpeech()
    key := c.Request().Header.Get("Sec-WebSocket-Key")

    defer func(ws *websocket.Conn) {
        log.Infof("%s disconnect", key)
        srv.SpeecherDisConnect(ws)
        _ = ws.Close()
    }(ws)

    token := c.QueryParam("token")
    register := srv.SpeecherConnect(ws, token)
    _ = ws.WriteMessage(websocket.TextMessage, register.Bytes())

    for {
        _, _, err = ws.ReadMessage()
        if err != nil {
            break
        }

        ws.WriteMessage(websocket.PongMessage, nil)
    }

    return nil
}

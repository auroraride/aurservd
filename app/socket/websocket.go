// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package socket

import (
    "github.com/google/uuid"
    "github.com/gorilla/websocket"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
    "net/http"
    "net/url"
)

type Websocket interface {
    Connect(hub *WebsocketHub, values url.Values)
    DisConnect(hub *WebsocketHub)
}

type WebsocketHub struct {
    *websocket.Conn
    ID string
}

func Wrap(c echo.Context, ws Websocket) error {
    var upgrader = websocket.Upgrader{}
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

    if err != nil {
        log.Error(err)
        return err
    }

    id := uuid.New().String()
    hub := &WebsocketHub{Conn: conn, ID: id}

    defer func(conn *websocket.Conn) {
        log.Infof("%s disconnect", id)
        ws.DisConnect(hub)
        _ = conn.Close()
    }(conn)

    ws.Connect(hub, c.QueryParams())

    for {
        _, _, err = conn.ReadMessage()
        if err != nil {
            break
        }

        _ = conn.WriteMessage(websocket.PongMessage, nil)
    }

    return nil
}

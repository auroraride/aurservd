// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package socket

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/gorilla/websocket"
    "github.com/labstack/echo/v4"
    "github.com/lithammer/shortuuid/v4"
    log "github.com/sirupsen/logrus"
    "net/http"
    "sync"
)

type Websocket interface {
    Connect(token string) (uint64, error)
    Prefix() string
}

var (
    clients sync.Map
)

type WebsocketHub struct {
    *websocket.Conn
    ClientID string
}

// SendMessage 发送消息
func (hub *WebsocketHub) SendMessage(message model.SocketBinaryMessage) {
    if hub == nil {
        return
    }
    _ = hub.WriteMessage(websocket.TextMessage, message.Bytes())
}

// DisConnect 断开连接
func (hub *WebsocketHub) DisConnect() {
    clients.Range(func(key, value any) bool {
        if v, ok := value.(*WebsocketHub); ok && hub.ClientID == v.ClientID {
            clients.Delete(key)
            return false
        }
        return true
    })
}

func GetKey(ws Websocket, id uint64) string {
    return fmt.Sprintf("%s#%d", ws.Prefix(), id)
}

func GetClientID(ws Websocket, id uint64) *WebsocketHub {
    return GetClient(GetKey(ws, id))
}

func GetClient(key string) *WebsocketHub {
    client, ok := clients.Load(key)
    if !ok {
        return nil
    }
    return client.(*WebsocketHub)
}

// Wrap 封装socket
func Wrap(c echo.Context, ws Websocket) error {
    var upgrader = websocket.Upgrader{}
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

    if err != nil {
        log.Error(err)
        return err
    }

    clientID := fmt.Sprintf("%s#%s", ws.Prefix(), shortuuid.New())

    hub := &WebsocketHub{Conn: conn, ClientID: clientID}

    defer func(conn *websocket.Conn) {
        log.Infof("%s disconnect", clientID)

        hub.DisConnect()
        _ = conn.Close()
    }(conn)

    var id uint64
    id, err = ws.Connect(c.QueryParam("token"))
    if err != nil {
        hub.SendMessage(&model.SocketMessage{Error: err.Error()})
    } else {
        hub.SendMessage(&model.SocketMessage{Error: ""})
    }

    key := GetKey(ws, id)
    log.Infof("%s Socket connected", key)

    // 断开已有的
    client := GetClient(key)
    if client != nil {
        _ = client.Close()
    }

    // 存储客户端
    clients.Store(key, hub)

    for {
        _, _, err = conn.ReadMessage()
        if err != nil {
            break
        }

        _ = conn.WriteMessage(websocket.PongMessage, nil)
    }

    return nil
}

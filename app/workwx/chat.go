// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package workwx

import (
    "fmt"
    "github.com/golang-module/carbon/v2"
    "time"
)

// MaterialAbnormality # 物资异常
// BatteryFault # 电池故障
// CabinetHealth # 电柜离线
// BatteryAbnormality # 电池异常变动
// SimExpires # SIM卡到期
// BranchExpires # 场地到期

type ChatMessage struct {
    Chatid  string `json:"chatid"`
    Msgtype string `json:"msgtype"`
}

type ChatContent struct {
    Content string `json:"content"`
}

type ChatMarkdown struct {
    ChatMessage
    Markdown ChatContent `json:"markdown"`
}

func (w *Client) SendMarkdown(chatid string, content string) error {
    var res baseResponse
    return w.RequestPost("/appchat/send", ChatMarkdown{
        ChatMessage: ChatMessage{
            Chatid:  chatid,
            Msgtype: "markdown",
        },
        Markdown: ChatContent{Content: content},
    }, &res)
}

func (w *Client) SendCabinetOffline(name, serial string) error {
    content := fmt.Sprintf(`电柜离线警告
>状态: <font color="warning">离线</font>
>电柜: %s
>编号: %s
>时间: <font color="comment">%s</font>`, name, serial, time.Now().Format(carbon.DateTimeLayout))
    return w.SendMarkdown("CabinetHealth", content)
}

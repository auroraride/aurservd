// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package workwx

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/golang-module/carbon/v2"
    "strings"
    "time"
)

// InventoryAbnormality # 物资异常
// ExchangeBinFault # 换电故障
// CabinetHealth # 电柜离线
// BatteryNumberAbnormality # 电池异常变动
// SimExpires # SIM卡到期
// BranchExpires # 场地到期
// Assistance # 救援提醒
// CabinetFault # 电柜故障提醒

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

func (w *Client) SendMarkdown(chatid string, content string) {
    var res baseResponse
    _ = w.RequestPost("/appchat/send", ChatMarkdown{
        ChatMessage: ChatMessage{
            Chatid:  chatid,
            Msgtype: "markdown",
        },
        Markdown: ChatContent{Content: content},
    }, &res)
}

// SendCabinetOffline 电柜离线警告
func (w *Client) SendCabinetOffline(name, serial, city string) {
    content := fmt.Sprintf(`电柜离线警告
>状态: <font color="warning">离线</font>
>城市: %s
>电柜: %s
>编号: %s
>时间: <font color="comment">%s</font>`, city, name, serial, time.Now().Format(carbon.DateTimeLayout))
    w.SendMarkdown("CabinetHealth", content)
}

// ExchangeBinFault 换电仓位处理失败警告
func (w *Client) ExchangeBinFault(city, name, serial, bin, rider, phone string, times int) {
    state := "处理失败"
    if times >= 2 {
        state += ", 已锁仓"
    }
    content := fmt.Sprintf(`换电仓位处理失败警告
>状态: <font color="warning">%s</font>
>城市: %s
>电柜: %s
>编号: %s
>仓位: %s
>骑手: %s
>电话: %s
>时间: <font color="comment">%s</font>`,
        state,
        city,
        name,
        serial,
        bin,
        rider,
        phone,
        time.Now().Format(carbon.DateTimeLayout),
    )
    w.SendMarkdown("ExchangeBinFault", content)
}

// SendInventory 发送物资警告
func (w *Client) SendInventory(duty bool, city, store string, e model.Employee, items []model.AttendanceInventory) {
    ds := "下班"
    if duty {
        ds = "上班"
    }

    arr := make([]string, len(items))
    for i, item := range items {
        arr[i] = fmt.Sprintf("%s: 库存`%d` 盘点`%d`", item.Name, item.StockNum, item.Num)
    }
    content := fmt.Sprintf(`物资异常警告
>类别: <font color="info">%s</font>
>城市: %s
>门店: %s
>店员: %s
>电话: %s
>时间: <font color="comment">%s</font>
>
%s`, ds, city, store, e.Name, e.Phone, time.Now().Format(carbon.DateTimeLayout), strings.Join(arr, "\n>    "))
    w.SendMarkdown("InventoryAbnormality", content)
}

// SendBatteryAbnormality 电池异常变动警告
func (w *Client) SendBatteryAbnormality(city, serial, name string, from, to, diff int) {
    content := fmt.Sprintf(`电池异常变动警告
>差值: <font color="warning">%d</font>
>城市: %s
>电柜: %s
>编号: %s
>前值: %d
>后值: %d
>时间: <font color="comment">%s</font>`,
        diff,
        city,
        name,
        serial,
        from,
        to,
        time.Now().Format(carbon.DateTimeLayout),
    )
    w.SendMarkdown("BatteryNumberAbnormality", content)
}

// SendSimExpires SIM卡到期警告
func (w *Client) SendSimExpires(data model.CabinetSimNotice) {
    c := ""
    if data.City != "" {
        c = fmt.Sprintf(">城市: %s", data.City)
    }
    content := fmt.Sprintf(`SIM卡到期警告
%s
>电柜: %s
>编号: %s
>卡号: %s
>到期: <font color="comment">%s</font>`,
        c,
        data.Name,
        data.Serial,
        data.Sim,
        data.End)
    w.SendMarkdown("SimExpires", content)
}

// SendBranchExpires 场地到期警告
func (w *Client) SendBranchExpires(data model.BranchExpriesNotice) {
    c := ""
    if data.City != "" {
        c = fmt.Sprintf(">城市: %s", data.City)
    }
    content := fmt.Sprintf(`场地到期警告
%s
>场地: %s
>到期: <font color="comment">%s</font>`,
        c,
        data.Name,
        data.End)
    w.SendMarkdown("BranchExpires", content)
}

// SendAssistance 救援提醒
func (w *Client) SendAssistance(data model.AssistanceNotice) {
    w.SendMarkdown("Assistance", fmt.Sprintf(`救援提醒
>电话: %s
>原因: %s
>地址: %s
>详细: %s
>时间: <font color="comment">%s</font>`,
        data.Phone,
        data.Reason,
        data.Address,
        data.AddressDetail,
        time.Now().Format(carbon.DateTimeLayout),
    ))
}

// SendCabinetFault 电柜故障提醒
func (w *Client) SendCabinetFault(data model.CabinetFaultNotice) {
    w.SendMarkdown("CabinetFault", fmt.Sprintf(`电柜故障提醒
>电话: %s
>城市: %s
>网点: %s
>电柜: %s
>编号: %s
>故障: <font color="warning">%s</font>
>描述: %s
>时间: <font color="comment">%s</font>`,
        data.Phone,
        data.City,
        data.Branch,
        data.Name,
        data.Serial,
        data.Fault,
        data.Description,
        time.Now().Format(carbon.DateTimeLayout),
    ))
}

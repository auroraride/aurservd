// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-11
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import "context"

func (c *Cabinet) GetSLSLogInfo() string {
    return "[电柜] " + c.Serial
}

func (r *Rider) GetSLSLogInfo() string {
    s := "[骑手] " + r.Phone
    p := r.Edges.Person
    if p == nil {
        p, _ = r.QueryPerson().First(context.Background())
    }
    if p != nil {
        s += " - " + p.Name
    }
    return s
}

func (e *Enterprise) GetSLSLogInfo() string {
    return "[企业] " + e.Name
}

func (pe *Person) GetSLSLogInfo() string {
    return "[用户] " + pe.Name
}

func (o *Order) GetSLSLogInfo() string {
    return "[订单] " + o.OutTradeNo
}

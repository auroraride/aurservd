// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-08
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "github.com/auroraride/adapter"
    jsoniter "github.com/json-iterator/go"
    "golang.org/x/exp/slices"
    "strconv"
)

type ExchangeLimiter interface {
    Duplicate() bool
    Sort()
}

type SettingExchangeLimits map[uint64]RiderExchangeLimit

func (s *SettingExchangeLimits) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(s)
}

func (s *SettingExchangeLimits) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, s)
}

type ExchangeLimit struct {
    Hours int `json:"hours"` // 时长
    Times int `json:"times"` // 时长内允许次数
}

type RiderExchangeLimit []*ExchangeLimit

type RiderExchangeLimitReq struct {
    IDPostReq
    ExchangeLimit RiderExchangeLimit `json:"exchangeLimit"`
}

func (el RiderExchangeLimit) String() string {
    buf := adapter.NewBuffer()
    defer adapter.ReleaseBuffer(buf)

    for i, limit := range el {
        buf.WriteString(strconv.Itoa(limit.Hours) + "小时内: ")
        buf.WriteString(strconv.Itoa(limit.Times) + "次")
        if i < len(el)-1 {
            buf.WriteString("; ")
        }
    }

    return buf.String()
}

func (el RiderExchangeLimit) Duplicate() bool {
    m := make(map[int]bool)
    for _, limit := range el {
        if m[limit.Hours] {
            return true
        }
        m[limit.Hours] = true
    }
    return false
}

func (el RiderExchangeLimit) Sort() {
    slices.SortFunc(el, func(a, b *ExchangeLimit) bool {
        return a.Hours < b.Hours
    })
}

type SettingExchangeFrequencies map[uint64]RiderExchangeFrequency

func (s *SettingExchangeFrequencies) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(s)
}

func (s *SettingExchangeFrequencies) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, s)
}

type ExchangeFrequency struct {
    Hours   int `json:"hours"`   // 时长
    Times   int `json:"times"`   // 时长内次数
    Minutes int `json:"minutes"` // 限制时间(分钟)
}

type RiderExchangeFrequency []*ExchangeFrequency

type RiderExchangeFrequencyReq struct {
    IDPostReq
    ExchangeFrequency RiderExchangeFrequency `json:"ExchangeFrequency"`
}

func (el RiderExchangeFrequency) String() string {
    buf := adapter.NewBuffer()
    defer adapter.ReleaseBuffer(buf)

    for i, limit := range el {
        buf.WriteString(strconv.Itoa(limit.Hours) + "小时内: ")
        buf.WriteString(strconv.Itoa(limit.Times) + "次, 限制")
        buf.WriteString(strconv.Itoa(limit.Minutes) + "分钟")
        if i < len(el)-1 {
            buf.WriteString("; ")
        }
    }

    return buf.String()
}

func (el RiderExchangeFrequency) Duplicate() bool {
    m := make(map[int]bool)
    for _, limit := range el {
        if m[limit.Hours] {
            return true
        }
        m[limit.Hours] = true
    }
    return false
}

func (el RiderExchangeFrequency) Sort() {
    slices.SortFunc(el, func(a, b *ExchangeFrequency) bool {
        return a.Hours < b.Hours
    })
}

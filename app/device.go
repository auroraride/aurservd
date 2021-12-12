// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
    "errors"
    "strings"
)

type DeviceType uint8

const (
    DeviceiOS = iota + 1
    DeviceAndroid
)

var (
    deviceStrMap = map[string]uint8{
        "ios":     DeviceiOS,
        "android": DeviceAndroid,
    }
)

type Device struct {
    Sn   string
    Type DeviceType
}

func NewDevice(sn, dt string) (d *Device, err error) {
    t, ok := deviceStrMap[strings.ToLower(dt)]
    if !ok {
        err = errors.New("设备类型错误")
    }
    d = &Device{
        Sn:   sn,
        Type: DeviceType(t),
    }
    return
}

func (d DeviceType) String() string {
    switch d {
    case DeviceiOS:
        return "iOS"
    case DeviceAndroid:
        return "Android"
    }
    return "unknown"
}

func (d DeviceType) Raw() uint8 {
    return uint8(d)
}

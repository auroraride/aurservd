// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
	"errors"
	"strings"

	jsoniter "github.com/json-iterator/go"
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
	Serial string     `json:"serial"`
	Type   DeviceType `json:"type"`
}

func NewDevice(sn, dt string) (d *Device, err error) {
	t, ok := deviceStrMap[strings.ToLower(dt)]
	if !ok {
		err = errors.New("设备类型错误")
	}
	d = &Device{
		Serial: sn,
		Type:   DeviceType(t),
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

func (d DeviceType) Value() uint8 {
	return uint8(d)
}

func (d *Device) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(d)
}

func (d *Device) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, d)
}

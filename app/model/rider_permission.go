// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type RiderPermissionType uint8

const (
	RiderPermissionTypeExchange   RiderPermissionType = iota // 换电
	RiderPermissionTypeBusiness                              // 业务
	RiderPermissionTypeAssistance                            // 救援
)

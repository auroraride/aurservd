// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// RiderContext 骑手上下文
type RiderContext struct {
	*Rider
}

// EmployeeContext 店员上下文
type EmployeeContext struct {
	*Employee
}

// ManagerContext 管理上下文
type ManagerContext struct {
	*Modifier
}

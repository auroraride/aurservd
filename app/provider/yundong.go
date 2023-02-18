// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
)

type yundong struct{}

func (p *yundong) Logger() *Logger {
    return nil
}

func NewYundong() *yundong {
    return &yundong{}
}

func (p *yundong) Cabinets() ([]*ent.Cabinet, error) {
    return ent.Database.Cabinet.QueryNotDeleted().Where(cabinet.Brand(model.CabinetBrandYundong.Value()), cabinet.Status(model.CabinetStatusNormal.Value())).All(context.Background())
}

func (p *yundong) Brand() string {
    return "云动"
}

func (p *yundong) FetchStatus(serial string) (online bool, bins model.CabinetBins, err error) {
    return
}

// DoorOperate 云动柜门操作
// user携带操作ID，比对操作日志实时获取状态
func (p *yundong) DoorOperate(code, serial, operation string, door int) (state bool) {
    return
}

func (p *yundong) Reboot(code string, serial string) (state bool) {
    return
}

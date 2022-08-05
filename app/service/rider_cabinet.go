// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/pkg/snag"
    log "github.com/sirupsen/logrus"
    "time"
)

// TODO 服务器崩溃后自动启动继续换电进程
// TODO 电柜缓存优化

type riderCabinetService struct {
    ctx     context.Context
    rider   *ent.Rider
    maxTime time.Duration // 单步骤最大处理时长
    logger  *logging.ExchangeLog

    cabinet   *ent.Cabinet
    subscribe *ent.Subscribe
    task      *ec.Task
}

func NewRiderCabinet(rider *ent.Rider) *riderCabinetService {
    s := &riderCabinetService{
        ctx:     context.Background(),
        maxTime: 180 * time.Second,
    }
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func (s *riderCabinetService) preprocess(serial string, bt business.Type) {
    cab := NewCabinet().QueryOneSerialX(serial)
    err := NewCabinet().UpdateStatus(cab)
    if err != nil {
        log.Error(err)
        snag.Panic("电柜状态获取失败")
    }

    var bn, en int

    for _, bin := range cab.Bin {
        // 仓位锁仓跳过计算
        if !bin.DoorHealth {
            continue
        }
        if bin.Battery {
            // 有电池
            bn += 1
        } else {
            // 无电池
            en += 1
        }
    }

    switch bt {
    case business.TypePause:
    case business.TypeUnsubscribe:
        if en < 2 {
            snag.Panic("仓位不足, 无法处理当前业务")
        }
        break
    case business.TypeActive:
    case business.TypeContinue:
        if bn < 2 {
            snag.Panic("电池不足, 无法处理当前业务")
        }
        break
    }

    s.cabinet = cab
}

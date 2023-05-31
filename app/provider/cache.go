// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-12
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/workwx"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/cache"
)

// getOfflineTime 获取离线时间
func getOfflineTime(serial string) time.Time {
	t, _ := cache.Get(context.Background(), fmt.Sprintf("OFFLINE-%s", serial)).Time()
	return t
}

// setOfflineTime 设置离线时间
func setOfflineTime(serial string, online bool) {
	key := fmt.Sprintf("OFFLINE-%s", serial)
	if !online {
		t := getOfflineTime(serial)
		if t.IsZero() {
			cache.Set(context.Background(), key, time.Now(), -1)
		}
	} else {
		// 在线则删除
		cache.Del(context.Background(), key)
	}
}

// isOffline 判定电柜是否离线, 3分钟以上算作离线
func isOffline(serial string) bool {
	t := getOfflineTime(serial)
	return !t.IsZero() && time.Since(t).Minutes() >= 3
}

func binFaultKey(serial string, index int) string {
	return fmt.Sprintf("AUTO-BINFAULT-%s-%d", serial, index)
}

// AutoBinFault 自动处理换电仓门操作失败
// 每次都推送, 超过两次锁仓
func AutoBinFault(operator model.CabinetDoorOperator, cab *ent.Cabinet, index int, status bool, lock func()) {
	binKey := binFaultKey(cab.Serial, index)
	ctx := context.Background()
	if status {
		cache.Del(ctx, binKey)
	} else {
		// 查询并保存失败次数
		times, _ := cache.Get(ctx, binKey).Int()
		times += 1
		cache.Set(ctx, binKey, times, -1)
		// 推送消息
		go workwx.New().ExchangeBinFault(cabinetCity(cab), cab.Name, cab.Serial, fmt.Sprintf("%d号仓", index+1), operator.Name, operator.Phone, times)
		// 锁仓
		if times > 2 {
			lock()
		}
	}
}

// delBinFault 删除仓门自动故障
func delBinFault(serial string, index int) {
	cache.Del(context.Background(), binFaultKey(serial, index))
}

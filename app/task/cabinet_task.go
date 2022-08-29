// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-29
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "context"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/mgo"
    "github.com/qiniu/qmgo/operator"
    "go.mongodb.org/mongo-driver/bson"
    "time"
)

type cabinetTask struct {
    max time.Duration
}

func NewCabinetTask() *cabinetTask {
    return &cabinetTask{
        max: ec.MaxTime * time.Second,
    }
}

func (t *cabinetTask) Start() {
    if ar.Config.Task.Cabinet {
        ticker := time.NewTicker(t.max)
        for {
            select {
            case <-ticker.C:
                // 检查并标记任务失效
                ts := time.Now().Add(-t.max).In(time.UTC)
                _, _ = mgo.CabinetTask.UpdateAll(context.Background(), bson.M{
                    "createAt": bson.M{
                        operator.Lt: ts,
                    },
                    "status":      ec.TaskStatusNotStart,
                    "deactivated": false,
                }, bson.M{
                    operator.Set: bson.M{"deactivated": true},
                })
                break
            }
        }
    }

}

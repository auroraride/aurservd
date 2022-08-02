// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-31
// Based on aurservd by liasica, magicrolan@qq.com.

package actuator

import (
    "context"
    "github.com/auroraride/aurservd/internal/mgo"
    "github.com/auroraride/aurservd/pkg/snag"
    jsoniter "github.com/json-iterator/go"
    "github.com/qiniu/qmgo/field"
    "github.com/qiniu/qmgo/operator"
    log "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

const (
    maxTime = 10.0
)

// TaskJob 电柜任务
type TaskJob string

const (
    JobExchange         TaskJob = "RDR_EXCHANGE"    // 骑手-换电
    JobRiderActive              = "RDR_ACTIVE"      // 骑手-激活
    JobRiderUnSubscribe         = "RDR_UNSUBSCRIBE" // 骑手-退租
    JobPause                    = "RDR_PAUSE"       // 骑手-寄存
    JobContinue                 = "RDR_CONTINUE"    // 骑手-取消寄存
    JobManagerOpen              = "MGR_OPEN"        // 管理-开门
    JobManagerExchange          = "MGR_EXCHANGE"    // 管理-换电
)

// Task 电柜任务详情
// TODO 存储骑手信息
type Task struct {
    field.DefaultField `bson:",inline"`

    Task        TaskJob    `json:"task" bson:"task"`                   // 任务类别
    StartAt     *time.Time `json:"startAt,omitempty" bson:"startAt"`   // 开始时间
    FinishAt    *time.Time `json:"finishAt,omitempty" bson:"finishAt"` // 结束时间
    Deactivated bool       `json:"deactivated" bson:"deactivated"`     // 是否已失效

    Cabinet  Cabinet   `json:"cabinet" bson:"cabinet"`   // 电柜信息
    Exchange *Exchange `json:"exchange" bson:"exchange"` // 换电信息
}

func (t *Task) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(t)
}

func (t *Task) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, t)
}

// Cabinet 任务电柜设备信息
type Cabinet struct {
    Serial         string `json:"serial" bson:"serial"`                 // 电柜编号
    Health         uint8  `json:"health" bson:"health"`                 // 电柜健康状态 0离线 1正常 2故障
    Doors          uint   `json:"doors" bson:"doors"`                   // 总仓位
    BatteryNum     uint   `json:"batteryNum" bson:"batteryNum"`         // 总电池数
    BatteryFullNum uint   `json:"batteryFullNum" bson:"batteryFullNum"` // 总满电电池数
}

// Save 存储任务信息
func (t *Task) Save() (string, error) {
    r, err := mgo.CabinetTask.InsertOne(context.Background(), t)
    if err != nil {
        return "", err
    }
    return r.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (t *Task) SaveX() string {
    id, err := t.Save()
    if err != nil {
        log.Error(err)
        snag.Panic("任务存储失败")
    }
    return id
}

// Start 开始任务
func (t *Task) Start() {
    ctx := context.Background()

    // 更新任务开始时间
    t.StartAt = Pointer(time.Now())
    _ = mgo.CabinetTask.UpdateId(ctx, t.Id, bson.M{
        operator.Set: bson.M{
            "startAt": Pointer(time.Now()),
        },
    })

    // 更新非当前任务为失效
    _, _ = mgo.CabinetTask.UpdateAll(ctx, bson.M{
        operator.Not: bson.M{
            "_id": t.Id,
        },
    }, bson.M{"deactivated": true})
}

// Deactive 设为失效
func (t *Task) Deactive() {
    _ = mgo.CabinetTask.UpdateId(context.Background(), t.Id, bson.M{"deactivated": true})
}

type ObtainReq struct {
    Serial      string             `json:"serial" bson:"serial"`
    ID          primitive.ObjectID `json:"id" bson:"_id"`
    Deactivated bool               `json:"deactivated" bson:"deactivated"`
}

// Obtain 获取任务
func Obtain(req ObtainReq) (t *Task) {
    ctx := context.Background()
    _ = mgo.CabinetTask.Find(ctx, req).One(t)
    if t == nil {
        return
    }
    // 任务未开始且超过10秒设置为超时
    if t.StartAt == nil && time.Now().Sub(t.UpdateAt).Seconds() > maxTime {
        t.Deactive()
        return nil
    }
    return t
}

// Busy 查询电柜是否繁忙
func Busy(serial string) bool {
    return Obtain(ObtainReq{Serial: serial}) != nil
}

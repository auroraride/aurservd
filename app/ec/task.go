// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-31
// Based on aurservd by liasica, magicrolan@qq.com.

package ec

import (
    "context"
    "fmt"
    "github.com/auroraride/adapter"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/snag"
    jsoniter "github.com/json-iterator/go"
    "github.com/rs/xid"
    "time"
)

const (
    DeactivateTime = 10.0 // 失效时间
)

type Updater func(task *Task)

// Task 电柜任务详情
// TODO 存储开仓信息, 业务信息, 管理员信息
type Task struct {
    ID               string           `json:"id,omitempty"`               // 任务ID
    CabinetID        uint64           `json:"cabinetID,omitempty"`        // 电柜ID
    Serial           string           `json:"serial,omitempty"`           // 电柜编码
    Job              model.Job        `json:"job,omitempty"`              // 任务类别
    Status           model.TaskStatus `json:"status,omitempty"`           // 任务状态
    StartAt          *time.Time       `json:"startAt,omitempty"`          // 开始时间
    StopAt           *time.Time       `json:"stopAt,omitempty"`           // 结束时间
    Message          string           `json:"message,omitempty"`          // 失败消息
    Cabinet          *Cabinet         `json:"cabinet,omitempty"`          // 电柜信息
    Rider            *Rider           `json:"rider,omitempty"`            // 骑手信息
    Exchange         *Exchange        `json:"exchange,omitempty"`         // 换电信息
    BussinessBinInfo *model.BinInfo   `json:"bussinessBinInfo,omitempty"` // 业务仓位
}

func (t *Task) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(t)
}

func (t *Task) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, t)
}

func (t *Task) String() string {
    // TODO 开仓信息, 业务信息, 管理员信息
    info := ""
    if t.Job == model.JobExchange {
        info += fmt.Sprintf(
            "骑手电话: %s, 名字: %s\n步骤: %s, 空: %d仓, 满: %d仓",
            t.Rider.Phone,
            t.Rider.Name,
            t.Exchange.CurrentStep().Step,
            t.Exchange.Empty.Index+1,
            t.Exchange.Fully.Index+1,
        )
    }
    return info
}

type Rider struct {
    ID    uint64 `json:"id"`
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

func (t *Task) Save() {
    ar.Redis.HSet(context.Background(), ar.TaskCacheKey, t.ID, t)
}

func (t *Task) Delete() {
    ar.Redis.HDel(context.Background(), ar.TaskCacheKey, t.ID)
}

func (t *Task) Create() *Task {
    t.ID = xid.New().String()
    hub.setter <- t.ID
    t.Save()
    return t
}

// Start 开始任务
func (t *Task) Start(cbs ...Updater) {
    // 更新任务开始时间
    t.Update(func(t *Task) {
        for _, cb := range cbs {
            cb(t)
        }
        t.StartAt = Pointer(time.Now())
        t.Status = model.TaskStatusProcessing
        hub.updater <- t.ID
    })

    // 删除所有未开始的非当前任务
    DeleteRange(func(x *Task) bool {
        return x.ID != t.ID && x.Serial == t.Serial && x.Status == model.TaskStatusNotStart
    })
}

// Stop 结束任务
func (t *Task) Stop(status model.TaskStatus) {
    t.Update(func(t *Task) {
        if status != model.TaskStatusSuccess {
            status = model.TaskStatusFail
        }
        t.StopAt = Pointer(time.Now())
        t.Status = status
    })
}

// Update 更新任务
func (t *Task) Update(cb Updater) {
    cb(t)
    t.Save()
}

// QueryID 查询任务
func QueryID(id xid.ID) (t *Task) {
    t = new(Task)
    err := ar.Redis.HGet(context.Background(), ar.TaskCacheKey, id.String()).Scan(t)
    if err != nil || t == nil {
        if t.Job == model.JobExchange && t.Exchange == nil {
            return nil
        }
        return nil
    }
    return
}

type ObtainReq struct {
    Serial    string `json:"serial,omitempty" bson:"serial,omitempty"`
    CabinetID uint64 `json:"cabinetId,omitempty" bson:"cabinetId,omitempty"`
}

// Obtain 获取进行中的任务信息
func Obtain(req ObtainReq) *Task {
    m := ar.Redis.HGetAll(context.Background(), ar.TaskCacheKey).Val()
    for _, v := range m {
        t := new(Task)
        if jsoniter.Unmarshal(adapter.ConvertString2Bytes(v), t) == nil {
            if t.Status == model.TaskStatusProcessing && (t.Serial == req.Serial || t.CabinetID == req.CabinetID) {
                return t
            }
        }
    }
    return nil
}

// Busy 查询电柜是否繁忙
func Busy(serial string) bool {
    task := Obtain(ObtainReq{Serial: serial})
    return task != nil
}

func BusyX(serial string) {
    task := Obtain(ObtainReq{Serial: serial})
    if task != nil {
        snag.Panic("电柜忙")
    }
}

// BusyFromID 查询电柜是否繁忙
func BusyFromID(id uint64) bool {
    task := Obtain(ObtainReq{CabinetID: id})
    return task != nil
}

func BusyFromIDX(id uint64) {
    if BusyFromID(id) {
        snag.Panic("电柜忙")
    }
}

// DeleteRange 删除所有指定条件的任务, 返回为非指定条件的任务
func DeleteRange(delcon func(x *Task) bool) (tasks map[string]*Task) {
    ctx := context.Background()
    m := ar.Redis.HGetAll(ctx, ar.TaskCacheKey).Val()
    tasks = make(map[string]*Task)
    for _, v := range m {
        t := new(Task)
        if jsoniter.Unmarshal(adapter.ConvertString2Bytes(v), t) == nil {
            if delcon(t) {
                t.Delete()
            } else {
                tasks[t.ID] = t
            }
        }
    }
    return
}

func Clear() {
    ar.Redis.Del(context.Background(), ar.TaskCacheKey)
}

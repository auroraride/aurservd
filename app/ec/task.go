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

// Job 电柜任务
type Job string

const (
    JobExchange         Job = "RDR_EXCHANGE"    // 骑手-换电
    JobRiderActive      Job = "RDR_ACTIVE"      // 骑手-激活
    JobRiderUnSubscribe Job = "RDR_UNSUBSCRIBE" // 骑手-退租
    JobPause            Job = "RDR_PAUSE"       // 骑手-寄存
    JobContinue         Job = "RDR_CONTINUE"    // 骑手-继续
    JobManagerOpen      Job = "MGR_OPEN"        // 管理-开门
    JobManagerLock      Job = "MGR_LOCK"        // 管理-锁仓
    JobManagerUnLock    Job = "MGR_UNLOCK"      // 管理-解锁
    JobManagerReboot    Job = "MGR_REBOOT"      // 管理-重启
    JobManagerExchange  Job = "MGR_EXCHANGE"    // 管理-换电
)

func (j Job) Label() string {
    switch j {
    case JobExchange:
        return "骑手换电"
    case JobRiderActive:
        return "骑手激活"
    case JobRiderUnSubscribe:
        return "骑手退租"
    case JobPause:
        return "骑手寄存"
    case JobContinue:
        return "骑手继续"
    case JobManagerOpen:
        return "管理开门"
    case JobManagerLock:
        return "管理锁仓"
    case JobManagerUnLock:
        return "管理解锁"
    case JobManagerReboot:
        return "管理重启"
    case JobManagerExchange:
        return "管理换电"
    }
    return "未知任务"
}

// Task 电柜任务详情
// TODO 存储开仓信息, 业务信息, 管理员信息
type Task struct {
    ID               string           // 任务ID
    CabinetID        uint64           // 电柜ID
    Serial           string           // 电柜编码
    Job              Job              // 任务类别
    Status           model.TaskStatus // 任务状态
    StartAt          *time.Time       // 开始时间
    StopAt           *time.Time       // 结束时间
    Message          string           // 失败消息
    Cabinet          *Cabinet         // 电柜信息
    Rider            *Rider           // 骑手信息
    Exchange         *Exchange        // 换电信息
    BussinessBinInfo *model.BinInfo   // 业务仓位

    deactivate *time.Timer // 失效处理
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
    if t.Job == JobExchange {
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
    t.deactivate = time.AfterFunc(DeactivateTime*time.Second, func() {
        // TODO 标记任务失败
        t.Delete()
    })
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
        t.deactivate.Reset(10 * time.Minute)
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
    _ = ar.Redis.HGet(context.Background(), ar.TaskCacheKey, id.String()).Scan(t)
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

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-31
// Based on aurservd by liasica, magicrolan@qq.com.

package ec

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/mgo"
    "github.com/auroraride/aurservd/pkg/snag"
    jsoniter "github.com/json-iterator/go"
    "github.com/qiniu/qmgo/operator"
    log "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

const (
    MaxTime = 10.0
)

type Updater func(task *Task)

// Job 电柜任务
type Job string

const (
    JobExchange         Job = "RDR_EXCHANGE"    // 骑手-换电
    JobRiderActive          = "RDR_ACTIVE"      // 骑手-激活
    JobRiderUnSubscribe     = "RDR_UNSUBSCRIBE" // 骑手-退租
    JobPause                = "RDR_PAUSE"       // 骑手-寄存
    JobContinue             = "RDR_CONTINUE"    // 骑手-继续
    JobManagerOpen          = "MGR_OPEN"        // 管理-开门
    JobManagerLock          = "MGR_LOCK"        // 管理-锁仓
    JobManagerUnLock        = "MGR_UNLOCK"      // 管理-解锁
    JobManagerReboot        = "MGR_REBOOT"      // 管理-重启
    JobManagerExchange      = "MGR_EXCHANGE"    // 管理-换电
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

type TaskStatus uint8

const (
    TaskStatusNotStart   TaskStatus = iota // 未开始
    TaskStatusProcessing                   // 处理中
    TaskStatusSuccess                      // 成功
    TaskStatusFail                         // 失败
)

func (ts TaskStatus) String() string {
    switch ts {
    case TaskStatusNotStart:
        return "未开始"
    case TaskStatusSuccess:
        return "成功"
    case TaskStatusFail:
        return "失败"
    default:
        return "处理中"
    }
}

// IsSuccess 是否成功
func (ts TaskStatus) IsSuccess() bool {
    return ts == TaskStatusSuccess
}

// Task 电柜任务详情
// TODO 存储开仓信息, 业务信息, 管理员信息
type Task struct {
    ID       primitive.ObjectID `bson:"_id"`
    CreateAt time.Time          `bson:"createAt"`
    UpdateAt time.Time          `bson:"updateAt"`

    CabinetID   uint64     `json:"cabinetId" bson:"cabinetId"`                 // 电柜ID
    Serial      string     `json:"serial" bson:"serial"`                       // 电柜编码
    Deactivated bool       `json:"deactivated" bson:"deactivated"`             // 是否已失效
    Job         Job        `json:"job" bson:"job"`                             // 任务类别
    Status      TaskStatus `json:"status" bson:"status"`                       // 任务状态
    StartAt     *time.Time `json:"startAt,omitempty" bson:"startAt,omitempty"` // 开始时间
    StopAt      *time.Time `json:"stopAt,omitempty" bson:"stopAt,omitempty"`   // 结束时间
    Message     string     `json:"message,omitempty" bson:"message,omitempty"` // 失败消息

    Cabinet          *Cabinet       `json:"cabinet" bson:"cabinet"`                                       // 电柜信息
    Rider            *Rider         `json:"rider" bson:"rider,omitempty"`                                 // 骑手信息
    Exchange         *Exchange      `json:"exchange" bson:"exchange,omitempty"`                           // 换电信息
    BussinessBinInfo *model.BinInfo `json:"bussinessBinInfo,omitempty" bson:"bussinessBinInfo,omitempty"` // 业务仓位
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

// Create 创建任务并存储
func (t *Task) Create() (primitive.ObjectID, error) {
    t.ID = primitive.NewObjectID()
    t.UpdateAt = time.Now()
    t.CreateAt = time.Now()
    r, err := mgo.CabinetTask.InsertOne(context.Background(), t)
    if err != nil {
        return primitive.NilObjectID, err
    }
    return r.InsertedID.(primitive.ObjectID), nil
}

func (t *Task) CreateX() *Task {
    id, err := t.Create()
    if err != nil {
        log.Error(err)
        snag.Panic("任务存储失败")
    }
    t.ID = id
    return t
}

// Start 开始任务
func (t *Task) Start(cb ...Updater) {
    ctx := context.Background()

    // 更新任务开始时间
    t.Update(func(t *Task) {
        if len(cb) > 0 {
            cb[0](t)
        }
        t.StartAt = Pointer(time.Now())
        t.Status = TaskStatusProcessing
    })

    // 更新非当前任务为失效
    _, _ = mgo.CabinetTask.UpdateAll(ctx, bson.M{
        "_id": bson.M{
            operator.Ne: t.ID,
        },
        "status":    0,
        "serial":    t.Serial,
        "deactived": false,
    }, bson.M{
        operator.Set: bson.M{"deactivated": true},
    })
}

// Stop 结束任务
func (t *Task) Stop(status TaskStatus) {
    t.Update(func(t *Task) {
        if status != TaskStatusSuccess {
            status = TaskStatusFail
        }
        t.StopAt = Pointer(time.Now())
        t.Status = status
    })
}

// Update 更新任务
func (t *Task) Update(cb Updater) {
    cb(t)
    t.UpdateAt = time.Now()
    _ = mgo.CabinetTask.UpdateId(context.Background(), t.ID, bson.M{
        operator.Set: t,
    })
}

// Deactive 设为失效
func (t *Task) Deactive() {
    _ = mgo.CabinetTask.UpdateId(context.Background(), t.ID, bson.M{"deactivated": true})
}

// IsDeactived 是否失效
func (t *Task) IsDeactived() bool {
    if t.StartAt == nil && time.Now().Sub(t.UpdateAt).Seconds() > MaxTime {
        t.Deactive()
        return true
    }
    return t.Deactivated
}

// QueryID 查询任务
func QueryID(id primitive.ObjectID) (t *Task) {
    t = new(Task)
    ctx := context.Background()
    _ = mgo.CabinetTask.Find(ctx, bson.M{"_id": id}).One(t)
    return
}

type ObtainReq struct {
    Serial      string     `json:"serial,omitempty" bson:"serial,omitempty"`
    Deactivated bool       `json:"deactivated" bson:"deactivated"`
    CabinetID   uint64     `json:"cabinetId,omitempty" bson:"cabinetId,omitempty"`
    Status      TaskStatus `json:"status" bson:"status"` // 任务状态
}

// Obtain 获取进行中的任务信息
func Obtain(req ObtainReq) (t *Task) {
    t = new(Task)
    if req.Status == 0 {
        req.Status = TaskStatusProcessing
    }
    ctx := context.Background()
    _ = mgo.CabinetTask.Find(ctx, req).One(t)
    if t == nil {
        return
    }
    // 任务未开始且超过10秒设置为超时
    if t.IsDeactived() {
        return nil
    }
    return t
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

// GetAllProcessing 获取所有正在进行中的任务
func GetAllProcessing() (tasks []*Task) {
    _ = mgo.CabinetTask.Find(context.Background(), bson.M{"status": TaskStatusProcessing}).All(&tasks)
    return
}

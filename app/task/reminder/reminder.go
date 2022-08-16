// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-16
// Based on aurservd by liasica, magicrolan@qq.com.

package reminder

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/setting"
    "github.com/auroraride/aurservd/internal/ent/subscribereminder"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "strconv"
    "sync"
    "time"
)

var tasks sync.Map

type reminderTask struct {
    ctx      context.Context
    ticker   *time.Ticker
    vmstel   *string
    vmstempl *string
    vmsdays  int
    smstempl string
    smsdays  int
    running  bool
}

var runner *reminderTask

func Run() {
    newReminder()
    runner.run()
}

func newReminder() {
    vmscfg := ar.Config.Aliyun.Vms.Reminder
    runner = &reminderTask{
        ctx:      context.Background(),
        ticker:   time.NewTicker(5 * time.Second),
        vmstempl: vmscfg.Template,
        vmstel:   vmscfg.Tel,
        smstempl: ar.Config.Aliyun.Sms.Template.Reminder,
    }

    notice := new(model.SettingOverdueNotice)
    sm, _ := ent.Database.Setting.Query().Where(setting.Key(model.SettingOverdue)).First(context.Background())
    if sm != nil {
        err := jsoniter.Unmarshal([]byte(sm.Content), notice)
        if err == nil {
            var smserr, vmserr error
            runner.smsdays, smserr = strconv.Atoi(notice.Sms)
            runner.vmsdays, vmserr = strconv.Atoi(notice.Call)
            runner.running = smserr == nil && vmserr == nil
        }
    }

    return
}

func Subscribe(sub *ent.Subscribe) {
    ri := sub.Edges.Rider
    if ri == nil {
        return
    }
    pe := ri.Edges.Person
    pl := sub.Edges.Plan
    if pe == nil || pl == nil {
        return
    }

    task := &Task{
        Name:        pe.Name,
        Phone:       ri.Phone,
        Days:        sub.Remaining,
        SubscribeID: sub.ID,
        PlanName:    pl.Name,
    }

    switch true {
    case runner.vmsdays == sub.Remaining:
        task.Type = subscribereminder.TypeVms
        break
    case runner.smsdays == sub.Remaining:
        task.Type = subscribereminder.TypeSms
        break
    default:
        return
    }

    tasks.Store(task.Phone, task)
}

func (r *reminderTask) run() {
    if !r.running {
        log.Error("催费任务未启动")
    }
    for {
        select {
        case <-r.ticker.C:
            h := time.Now().Hour()
            // 下午4点发送
            if h == 16 {

                duplicateRemove()

                tasks.Range(func(_, v any) bool {
                    switch i := v.(type) {
                    case *Task:
                        switch i.Type {
                        case subscribereminder.TypeVms:
                            r.sendvms(i)
                            break
                        case subscribereminder.TypeSms:
                            r.sendsms(i)
                            break
                        }
                        break
                    }
                    return true
                })
            }
            break
        }
    }
}

// duplicateRemove 删除已执行成功的任务
func duplicateRemove() {
    // 查询今日推送
    items, _ := ent.Database.SubscribeReminder.Query().Where(subscribereminder.Date(time.Now().Format(carbon.DateLayout))).All(context.Background())
    // 判断任务是否完成
    for _, item := range items {
        if item.Success {
            // 如果任务完成, 直接删除
            tasks.Delete(item.Phone)
        } else {
            // 如果未完成, 携带以便更新
            v, ok := tasks.Load(item.Phone)
            if ok {
                v.(*Task).SubscribeReminder = item
                tasks.Store(item.Phone, v)
            }
        }
    }
}

func (r *reminderTask) sendvms(task *Task) {
    task.Success = ali.NewVms().SendVoiceMessageByTts(
        tools.NewPointerInterface(task.Phone),
        tools.NewPointerInterface(fmt.Sprintf(`{"name":"%s","product": "%s"}`, task.Name, task.PlanName)),
        r.vmstel,
        r.vmstempl,
    )
    r.updateOrSave(task)
}

func (r *reminderTask) sendsms(task *Task) {
    client, err := ali.NewSms()
    if err != nil {
        log.Error(err)
    }
    id, _ := client.SetTemplate(r.smstempl).
        SetParam(map[string]string{
            "name":    task.Name,
            "product": task.PlanName,
        }).
        SendCode(task.Phone)
    task.Success = id != ""

    r.updateOrSave(task)
}

func (r *reminderTask) updateOrSave(task *Task) {
    if task.SubscribeReminder != nil {
        _, _ = task.SubscribeReminder.Update().SetSuccess(task.Success).Save(r.ctx)
    } else {
        _, _ = ent.Database.SubscribeReminder.Create().
            SetPhone(task.Phone).
            SetSubscribeID(task.SubscribeID).
            SetPlanName(task.PlanName).
            SetName(task.Name).
            SetType(task.Type).
            SetDays(task.Days).
            SetSuccess(task.Success).
            SetDate(time.Now().Format(carbon.DateLayout)).
            Save(r.ctx)
    }
    tasks.Delete(task.Phone)
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-16
// Based on aurservd by liasica, magicrolan@qq.com.

package reminder

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/setting"
	"github.com/auroraride/aurservd/internal/ent/subscribereminder"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var tasks sync.Map

type Task struct {
	Name        string
	Phone       string
	Days        int
	SubscribeID uint64
	RiderID     uint64
	PlanName    string
	PlanID      uint64
	PlanDays    uint
	Success     bool
	Fee         *float64
	FeeFormula  *string

	vms *vmsconfig
	sms string
}

type vmsconfig struct {
	tel  *string
	tmpl *string
}

type reminderTask struct {
	ctx     context.Context
	ticker  *time.Ticker
	vms     map[int]*vmsconfig
	sms     map[int]string
	running bool
}

var runner *reminderTask

func Run() {
	newReminder()
	runner.run()
}

func Reset() {
	tasks.Range(func(key, _ any) bool {
		tasks.Delete(key)
		return true
	})
}

func newReminder() {
	vmscfg := ar.Config.Aliyun.Vms
	smscfg := ar.Config.Aliyun.Sms.Template
	runner = &reminderTask{
		ctx:    context.Background(),
		ticker: time.NewTicker(1 * time.Minute),
		vms:    make(map[int]*vmsconfig),
		sms:    make(map[int]string),
	}

	notice := new(model.SettingReminderNotice)
	sm, _ := ent.Database.Setting.Query().Where(setting.Key(model.SettingReminderKey)).First(context.Background())
	if sm != nil {
		err := jsoniter.Unmarshal([]byte(sm.Content), notice)
		if err == nil {
			for _, d := range notice.Sms {
				if d >= 0 {
					runner.sms[d] = smscfg.Reminder
				} else {
					runner.sms[d] = smscfg.Overdue
				}
			}
			for _, d := range notice.Vms {
				if d >= 0 {
					runner.vms[d] = &vmsconfig{
						tel:  vmscfg.Reminder.Tel,
						tmpl: vmscfg.Reminder.Template,
					}
				} else {
					runner.vms[d] = &vmsconfig{
						tel:  vmscfg.Overdue.Tel,
						tmpl: vmscfg.Overdue.Template,
					}
				}
			}
			runner.running = ar.Config.Task.Reminder
		}
	}

	return
}

func Subscribe(sub *ent.Subscribe) {
	if !runner.running {
		return
	}

	ri := sub.Edges.Rider
	if ri == nil {
		return
	}
	pl := sub.Edges.Plan
	if pl == nil {
		return
	}

	task := &Task{
		Name:        ri.Name,
		Phone:       ri.Phone,
		Days:        sub.Remaining,
		SubscribeID: sub.ID,
		RiderID:     sub.RiderID,
		PlanName:    pl.Name,
		PlanDays:    pl.Days,
		PlanID:      pl.ID,
	}

	if cfg, ok := runner.vms[sub.Remaining]; ok {
		task.vms = cfg
	}

	if tmpl, ok := runner.sms[sub.Remaining]; ok {
		task.sms = tmpl
	}

	if task.sms == "" && task.vms == nil {
		return
	}

	if sub.Remaining < 0 {
		f, fl := pl.OverdueFee(sub.Remaining)
		task.FeeFormula = silk.Pointer(fl)
		task.Fee = silk.Pointer(f)
	}

	tasks.Store(task.Phone, task)
}

func (r *reminderTask) run() {
	if !r.running {
		zap.L().Info("催费任务未启动")
	}
	for {
		select {
		case <-r.ticker.C:
			h := time.Now().Hour()
			// 下午4点发送
			if h == 16 {

				duplicateRemove()

				tasks.Range(func(_, v any) bool {
					switch t := v.(type) {
					case *Task:
						if t.sms != "" {
							r.sendsms(t)
						}
						if t.vms != nil {
							r.sendvms(t)
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
		}
	}
}

func (r *reminderTask) sendvms(task *Task) {

	type template struct {
		Name    string   `json:"name"`
		Product string   `json:"product"`
		Days    *int     `json:"days,omitempty"`
		Fee     *float64 `json:"fee,omitempty"`
	}

	data := template{
		Name:    task.Name,
		Product: task.PlanName,
		Fee:     task.Fee,
	}
	if task.Days < 0 {
		data.Days = silk.Pointer(task.Days)
	}
	b, _ := jsoniter.Marshal(data)

	vms := task.vms
	task.Success = ali.NewVms().SendVoiceMessageByTts(
		silk.Pointer(task.Phone),
		silk.Pointer(string(b)),
		vms.tel,
		vms.tmpl,
	)

	r.updateOrSave(task, subscribereminder.TypeVms)
}

func (r *reminderTask) sendsms(task *Task) {
	client, err := ali.NewSms()
	if err != nil {
		zap.L().Error("短信发送失败", zap.Error(err))
	}
	data := map[string]string{
		"name":    task.Name,
		"product": task.PlanName,
	}
	if task.Days < 0 {
		data["days"] = fmt.Sprintf("%d", task.Days)
		data["fee"] = fmt.Sprintf("%.2f", *task.Fee)
	}
	id, _ := client.SetTemplate(task.sms).
		SetParam(data).
		SendCode(task.Phone)
	task.Success = id != ""

	r.updateOrSave(task, subscribereminder.TypeVms)
}

func (r *reminderTask) updateOrSave(task *Task, typ subscribereminder.Type) {
	_, _ = ent.Database.SubscribeReminder.Create().
		SetRiderID(task.RiderID).
		SetPhone(task.Phone).
		SetSubscribeID(task.SubscribeID).
		SetPlanName(fmt.Sprintf("%s - %d天", task.PlanName, task.PlanDays)).
		SetName(task.Name).
		SetType(typ).
		SetDays(task.Days).
		SetPlanID(task.PlanID).
		SetSuccess(task.Success).
		SetDate(time.Now().Format(carbon.DateLayout)).
		SetNillableFeeFormula(task.FeeFormula).
		SetNillableFee(task.Fee).
		Save(r.ctx)
	tasks.Delete(task.Phone)
}

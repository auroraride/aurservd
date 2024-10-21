package task

import (
	"context"
	"strconv"

	"github.com/auroraride/adapter/log"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/purchasepayment"
	"github.com/auroraride/aurservd/pkg/silk"
)

type vmsconfig struct {
	tel  *string
	tmpl *string
}

type Task struct {
	Name    string
	Phone   string
	RiderID uint64
	Time    string
	Amount  string

	vms *vmsconfig
	sms string
}

type purchaseReminderTask struct {
}

func NewPurchaseReminder() *purchaseReminderTask {
	return &purchaseReminderTask{}
}

func (t *purchaseReminderTask) Start() {
	if !ar.Config.Task.PurchaseReminder {
		return
	}

	c := cron.New()
	_, err := c.AddFunc("0 9 * * *", func() {
		zap.L().Info("开始执行 @daily[purchaseReminderTask] 定时任务")
		t.Do()
	})
	if err != nil {
		zap.L().Fatal("@daily[purchaseReminderTask] 定时任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

// Do 提前提醒用户还款 提前三天短信提醒 提前一天语音提醒
func (t *purchaseReminderTask) Do() {
	// 查询即将逾期的账单
	now := carbon.Now().AddDay().StartOfDay()
	payments, _ := ent.Database.PurchasePayment.QueryNotDeleted().Where(
		purchasepayment.StatusEQ(purchasepayment.StatusObligation),
		purchasepayment.BillingDateLTE(now.AddDays(3).StdTime()),
	).WithRider().All(context.Background())
	tk := Task{
		vms: &vmsconfig{
			tel:  ar.Config.Aliyun.Vms.PurchaseOverdue.Tel,
			tmpl: ar.Config.Aliyun.Vms.PurchaseOverdue.Template,
		},
		sms: ar.Config.Aliyun.Sms.Template.PurchaseOverdue,
	}
	for _, v := range payments {
		// 发送短信
		tk.Name = v.Edges.Rider.Name
		tk.Phone = v.Edges.Rider.Phone
		tk.Amount = strconv.FormatFloat(v.Amount, 'f', 2, 64)
		tk.Time = v.BillingDate.Format("2006-01-02")
		t.sendsms(&tk)
	}
	payments, _ = ent.Database.PurchasePayment.QueryNotDeleted().Where(
		purchasepayment.StatusEQ(purchasepayment.StatusObligation),
		purchasepayment.BillingDateLTE(now.EndOfDay().AddDay().StdTime()),
	).WithRider().All(context.Background())
	for _, v := range payments {
		// 发送语音
		tk.Name = v.Edges.Rider.Name
		tk.Time = v.BillingDate.Format("2006-01-02")
		t.sendvms(&tk)
	}
}

// sendvms 发送语音
func (t *purchaseReminderTask) sendvms(task *Task) {
	type template struct {
		Name       string `json:"name"`
		DateFormat string `json:"date_format"`
	}

	data := template{
		Name:       task.Name,
		DateFormat: task.Time,
	}

	b, _ := jsoniter.Marshal(data)

	vms := task.vms
	ali.NewVms().SendVoiceMessageByTts(
		silk.Pointer(task.Phone),
		silk.Pointer(string(b)),
		vms.tel,
		vms.tmpl,
	)
	zap.L().Info("逾期语音发送成功", log.JsonData(data))
}

// sendsms 发送短信
func (t *purchaseReminderTask) sendsms(task *Task) {
	client, err := ali.NewSms()
	if err != nil {
		zap.L().Error("短信发送失败", zap.Error(err))
	}
	data := map[string]string{
		"name":   task.Name,
		"time":   task.Time,
		"amount": task.Amount,
	}

	id, _ := client.SetTemplate(task.sms).
		SetParam(data).
		SendCode(task.Phone)
	zap.L().Info("逾期短信发送成功", log.JsonData(data), zap.String("id", id))
}

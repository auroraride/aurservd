package task

import (
	"context"

	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/contract"
)

type contractTask struct {
}

func NewContractTask() *contractTask {
	return &contractTask{}
}

func (t *contractTask) Start() {
	c := cron.New()
	_, err := c.AddFunc("0 1 * * *", func() {
		zap.L().Info("开始执行合同过期任务")
		go t.Do()
	})
	if err != nil {
		zap.L().Fatal("执行合同过期任务失败", zap.Error(err))
		return
	}
	c.Start()
}

// Do 更新企业订单
func (*contractTask) Do() {
	expiresAt := carbon.Now().StdTime()
	list, _ := ent.Database.Contract.QueryNotDeleted().Where(contract.ExpiresAtLTE(expiresAt), contract.SignedAtIsNil()).All(context.Background())
	for _, v := range list {
		err := ent.Database.Contract.UpdateOne(v).SetStatus(model.ContractStatusExpired.Value()).SetEffective(false).Exec(context.Background())
		if err != nil {
			zap.L().Error("更新合同失败", zap.Error(err))
			continue
		}
	}
}

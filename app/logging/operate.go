// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/auroraride/adapter"
	"github.com/lithammer/shortuuid/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

type Operator struct {
	Type  model.OperatorType `json:"operatorType"`
	ID    uint64             `json:"operatorId"`
	Phone string             `json:"operatorPhone"`
	Name  string             `json:"operatorName"`
}

func (o *Operator) OperatorRole() model.CabinetDoorOperatorRole {
	switch o.Type {
	default:
		return model.CabinetDoorOperatorRoleManager
	case model.OperatorTypeRider:
		return model.CabinetDoorOperatorRoleRider
	case model.OperatorTypeMaintainer:
		return model.CabinetDoorOperatorRoleMaintainer
	}
}

func GetOperator(data any) (*Operator, error) {
	switch v := data.(type) {
	default:
		return nil, adapter.ErrorUserRequired
	case *model.Modifier:
		return &Operator{
			Type:  model.OperatorTypeManager,
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}, nil
	case *ent.Employee:
		return &Operator{
			Type:  model.OperatorTypeEmployee,
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}, nil
	case *model.Employee:
		return &Operator{
			Type:  model.OperatorTypeEmployee,
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}, nil
	case *ent.Cabinet:
		return &Operator{
			Type:  model.OperatorTypeCabinet,
			ID:    v.ID,
			Phone: v.Serial,
			Name:  v.Name,
		}, nil
	case *ent.Agent:
		return &Operator{
			Type:  model.OperatorTypeAgent,
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}, nil
	case *ent.Maintainer:
		return &Operator{
			Type:  model.OperatorTypeMaintainer,
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}, nil
	case *ent.Rider:
		return &Operator{
			Type:  model.OperatorTypeRider,
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}, nil
	case *ent.AssetManager:
		return &Operator{
			Type:  model.OperatorTypeAssetManager,
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}, nil
	}
}

func GetOperatorX(data any) *Operator {
	o, err := GetOperator(data)
	if err != nil {
		snag.Panic(err)
	}
	return o
}

func (o *Operator) GetAdapterUser() (*adapter.User, error) {
	switch o.Type {
	default:
		return nil, adapter.ErrorUserRequired
	case model.OperatorTypeRider:
		return &adapter.User{
			Type: adapter.UserTypeRider,
			ID:   o.Phone,
		}, nil
	case model.OperatorTypeEmployee:
		return &adapter.User{
			Type: adapter.UserTypeEmployee,
			ID:   o.Phone,
		}, nil
	case model.OperatorTypeManager:
		return &adapter.User{
			Type: adapter.UserTypeManager,
			ID:   o.Phone,
		}, nil
	case model.OperatorTypeAgent:
		return &adapter.User{
			Type: adapter.UserTypeAgent,
			ID:   o.Phone,
		}, nil
	case model.OperatorTypeMaintainer:
		return &adapter.User{
			Type: adapter.UserTypeMaintainer,
			ID:   o.Phone,
		}, nil
	}
}

func (o *Operator) GetAdapterUserX() *adapter.User {
	user, err := o.GetAdapterUser()
	if err != nil {
		snag.Panic(err)
	}
	return user
}

// OperateLog 系统操作日志
type OperateLog struct {
	ID string `json:"id" sls:"操作ID"`

	RefID    uint64 `json:"refId" sls:"关联ID" index:"doc"`
	RefTable string `json:"refTable" sls:"关联表" index:"doc"`

	Operate model.Operate `json:"operate" sls:"类别" string:"true" index:"doc"`
	Remark  string        `json:"remark" sls:"备注"`
	Before  string        `json:"before" sls:"操作前"`
	After   string        `json:"after" sls:"操作后"`

	OperatorType  model.OperatorType `json:"operatorType" sls:"操作人类别" index:"doc"` // 0管理员 1员工
	OperatorID    uint64             `json:"operatorId" sls:"操作人ID" index:"doc"`
	OperatorPhone string             `json:"operatorPhone" sls:"操作人电话" index:"doc"`
	OperatorName  string             `json:"operatorName" sls:"操作人" index:"doc"`

	Info string `json:"info" sls:"信息"`
	Time string `json:"time" sls:"时间"`
}

func (o *OperateLog) GetLogstoreName() string {
	return ar.Config.Aliyun.Sls.OperateLog
}

func NewOperateLog() *OperateLog {
	return &OperateLog{
		ID: shortuuid.New(),
	}
}

func (o *OperateLog) SetInfo(info string) *OperateLog {
	o.Info = info
	return o
}

func (o *OperateLog) SetRef(ref model.Table) *OperateLog {
	o.RefTable = ref.GetTableName()
	o.RefID = ref.GetID()
	switch v := ref.(type) {
	case model.TableSlsInfo:
		o.Info = v.GetSLSLogInfo()
	}
	return o
}

func (o *OperateLog) SetRefManually(table string, id uint64) *OperateLog {
	o.RefTable = table
	o.RefID = id
	return o
}

func (o *OperateLog) SetOperate(operate model.Operate) *OperateLog {
	o.Operate = operate
	return o
}

func (o *OperateLog) SetRemark(remark string) *OperateLog {
	o.Remark = remark
	return o
}

func (o *OperateLog) SetDiff(before, after string) *OperateLog {
	o.Before = before
	o.After = after
	return o
}

func (o *OperateLog) SetOperator(op *Operator) *OperateLog {
	if op != nil {
		o.OperatorID = op.ID
		o.OperatorPhone = op.Phone
		o.OperatorName = op.Name
		o.OperatorType = op.Type
	}
	return o
}

func (o *OperateLog) SetModifier(m *model.Modifier) *OperateLog {
	if m != nil {
		o.OperatorID = m.ID
		o.OperatorPhone = m.Phone
		o.OperatorName = m.Name
		o.OperatorType = model.OperatorTypeManager
	}
	return o
}

func (o *OperateLog) SetAgent(ag *ent.Agent) *OperateLog {
	if ag != nil {
		o.OperatorID = ag.ID
		o.OperatorPhone = ag.Phone
		o.OperatorName = ag.Name
		o.OperatorType = model.OperatorTypeAgent
	}
	return o
}

func (o *OperateLog) SetEmployee(e *model.Employee) *OperateLog {
	if e != nil {
		o.OperatorID = e.ID
		o.OperatorPhone = e.Phone
		o.OperatorName = e.Name
		o.OperatorType = model.OperatorTypeEmployee
	}
	return o
}

func (o *OperateLog) SetCabinet(c *model.CabinetBasicInfo) *OperateLog {
	if c != nil {
		o.OperatorID = c.ID
		o.OperatorPhone = c.Serial
		o.OperatorName = c.Name
		o.OperatorType = model.OperatorTypeCabinet
	}
	return o
}

func (o *OperateLog) SetMaintainer(c *ent.Maintainer) *OperateLog {
	if c != nil {
		o.OperatorID = c.ID
		o.OperatorPhone = c.Phone
		o.OperatorName = c.Name
		o.OperatorType = model.OperatorTypeMaintainer
	}
	return o
}

func (o *OperateLog) Send() {
	PutLog(o)
}

// GetLogs 获取日志
// 参考 https://help.aliyun.com/document_detail/29029.htm?spm=a2c4g.11186623.0.0.21752a73O1u0I4#t13238.html
func (o *OperateLog) GetLogs(from time.Time, query string, offset, limit int64) []map[string]string {
	cfg := ar.Config.Aliyun.Sls
	response, err := ali.NewSls().GetLogsV2(cfg.Project, cfg.OperateLog, &sls.GetLogRequest{
		From:    from.Unix(),
		To:      time.Now().Unix(),
		Reverse: true,
		Query:   query,
		Lines:   limit,
		Offset:  offset,
	})
	if err != nil {
		return make([]map[string]string, 0)
	}
	return response.Logs
}

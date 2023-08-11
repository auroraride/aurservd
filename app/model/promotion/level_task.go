package promotion

type LevelTaskType uint8

const (
	LevelTaskTypeSign  LevelTaskType = iota + 1 // 签约
	LevelTaskTypeRenew                          // 续费
)

func (t LevelTaskType) Value() uint8 {
	return uint8(t)
}

// tosting
func (t LevelTaskType) String() string {
	switch t {
	case LevelTaskTypeSign:
		return "签约"
	case LevelTaskTypeRenew:
		return "续费"
	}
	return ""
}

type LevelTask struct {
	ID          uint64        `json:"id" param:"id"`                    // id
	Name        string        `json:"name" `                            // 任务名称
	Description string        `json:"description" `                     // 任务描述
	Type        LevelTaskType `json:"type" `                            // 类型 1:签约 2:续费
	GrowthValue uint64        `json:"growthValue" validate:"required" ` // 成长值
}

type LevelTaskSelect struct {
	ID   uint64 `json:"id" param:"id"` // id
	Name string `json:"name" `         // 任务名称
}

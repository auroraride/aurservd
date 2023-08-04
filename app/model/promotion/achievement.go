package promotion

type AchievementType uint8

const (
	AchievementInvite   AchievementType = iota + 1 // 邀请成就
	AchievementEarnings                            // 收益成就
)

func (a AchievementType) Value() uint8 {
	return uint8(a)
}

func (a AchievementType) String() string {
	switch a {
	case AchievementInvite:
		return "邀请成就"
	case AchievementEarnings:
		return "收益成就"
	}
	return ""
}

// Achievement 会员成就
type Achievement struct {
	ID        uint64          `json:"id" param:"id"`                         // id
	Name      string          `json:"name" validate:"required"`              // 成就名称
	Type      AchievementType `json:"type"  validate:"required" enums:"1,2"` // 成就类型 1:邀请成就 2:收益成就
	Icon      string          `json:"icon" validate:"required"`              // 成就图标
	Condition uint64          `json:"condition" validate:"required" `        // 完成条件
}

type UploadIcon struct {
	Icon string `json:"icon" validate:"required" ` // 成就图标
}

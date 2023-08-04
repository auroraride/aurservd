package promotion

type PrivilegeType uint8

const (
	PrivilegeNone       PrivilegeType = iota // 无权益
	PrivilegeCommission                      // 佣金提高(%)
)

func (a PrivilegeType) Value() uint8 {
	return uint8(a)
}

func (a PrivilegeType) String() string {
	switch a {
	case PrivilegeNone:
		return "无权益"
	case PrivilegeCommission:
		return "佣金提高(%)"
	}
	return ""
}

type Privilege struct {
	ID          uint64 `json:"id" `          // id
	Type        uint64 `json:"type" `        // 权益类型 0:无权权益 1: 佣金提高(%)
	Name        string `json:"name" `        // 权益名称
	Description string `json:"description" ` // 权益描述
	Value       uint64 `json:"value" `       // 权益值
}

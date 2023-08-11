package promotion

const (
	PersonUnauthenticated      PersonAuthStatus = iota // 未认证
	PersonAuthenticated                                // 已认证
	PersonAuthenticationFailed                         // 认证失败
)

type PersonAuthStatus uint8

func (s PersonAuthStatus) String() string {
	switch s {
	case PersonUnauthenticated:
		return "未认证"
	case PersonAuthenticated:
		return "已认证"
	}
	return "认证失败"
}

func (s PersonAuthStatus) Value() uint8 {
	return uint8(s)
}

type RealNameAuthReq struct {
	Name   string `json:"realName" validate:"required" ` // 真实姓名
	IdCard string `json:"idCard" validate:"required"`    // 身份证号
}

type RealNameAuthRes struct {
	Success bool `json:"success" ` // 是否成功
}

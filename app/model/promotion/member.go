package promotion

import (
	"database/sql"

	"github.com/auroraride/aurservd/app/model"
)

// MemberReq  会员请求参数
type MemberReq struct {
	model.PaginationReq
	MemberFilter
	ID uint64 `json:"id" param:"id"` // id
}

type MemberUpdateReq struct {
	ID     uint64  `json:"id" param:"id" validate:"required"` // id
	Enable *bool   `json:"enable"`                            // 用户状态 false:禁用 true:启用
	Name   *string `json:"name"`                              // 姓名
}

type MemberFilter struct {
	Keyword      *string           `json:"keyword" query:"keyword"`                     // 关键词 姓名/手机号
	Enable       *bool             `json:"enable" query:"enable"  enums:"true,false"`   // 用户状态 false:禁用 true:启用
	CommissionID *uint64           `json:"commissionId" query:"commissionId" `          // 返佣方案id
	LevelID      *uint64           `json:"levelId" query:"levelId"`                     // 会员等级
	Start        *string           `json:"start" query:"start" `                        // 开始日期
	End          *string           `json:"end"  query:"end"`                            // 结束日期
	AuthStatus   *PersonAuthStatus `json:"authStatus" query:"authStatus" enums:"0,1,2"` // 是否实名认证
}

// MemberRes 会员返回参数
type MemberRes struct {
	ID                  uint64           `json:"id" `                  // id
	MemberInfo          MemberBaseInfo   `json:"memberInfo" `          // 会员信息
	Enable              bool             `json:"enable"`               // 用户状态 false:禁用 true:启用
	CommissionName      string           `json:"commissionName" `      // 返佣方案
	CommissionID        uint64           `json:"commissionId" `        // 返佣方案id
	Level               uint64           `json:"level" `               // 会员等级
	TotalBalance        float64          `json:"totalBalance" `        // 总余额
	ParentInfo          *MemberBaseInfo  `json:"parentId,omitempty"`   // 上级信息
	FirstLevel          uint64           `json:"firstLevel" `          // 第一级团员人数
	SecondLevel         uint64           `json:"secondLevel" `         // 第二级团员人数
	CreatedAt           string           `json:"createdAt" `           // 注册时间
	CurrentGrowthValue  uint64           `json:"currentGrowthValue" `  // 当前成长值
	PrivilegeCommission float64          `json:"privilegeCommission" ` // 等级权益 佣金提高(%)
	AuthStatusName      string           `json:"authStatusName"`       // 实名认证状态名称
	AuthStatus          PersonAuthStatus `json:"authStatus"`           // 实名认证状态 0:未认证 1:已认证 2:认证失败
	BankCard            []*BankCardRes   `json:"bankCard,omitempty"`   // 银行卡
}

// MemberBaseInfo 会员基本信息
type MemberBaseInfo struct {
	ID           uint64 `json:"id" `                     // id
	Phone        string `json:"phone" `                  // 手机号
	Name         string `json:"name,omitempty" `         // 姓名
	IDCardNumber string `json:"idCardNumber,omitempty" ` // 身份证号
}

// MemberTeamReq 会员团队请求参数
type MemberTeamReq struct {
	ID uint64 `json:"id" param:"id"` // id
	model.PaginationReq
	MemberTeamFilter // 会员团队筛选条件
}

type TemaLevel uint8

const (
	TemaLevel1 TemaLevel = iota + 1 // 一级团员
	TemaLevel2                      // 二级团员
)

func (l TemaLevel) Value() uint8 {
	return uint8(l)
}

// 转换为字符串
func (l TemaLevel) String() string {
	switch l {
	case TemaLevel1:
		return "一级团员"
	case TemaLevel2:
		return "二级团员"
	default:
		return "未知"
	}

}

// MemberTeamRes 会员团队列表信息
type MemberTeamRes struct {
	ID                  uint64 `json:"id" `                 // id
	Phone               string `json:"phone" `              // 手机号
	Name                string `json:"name" `               // 姓名
	Level               string `json:"level" `              // 团员层级
	RenewalCount        uint64 `json:"renewalCount" `       // 续费次数
	SubscribeStatus     uint8  `json:"subscribeStatus" `    // 订阅状态 0:未激活 1:计费中 4:已退订
	SubscribeStatusName string `json:"subscribeStatusName"` // 订阅状态名称
	SubscribeStartAt    string `json:"subscribeStartAt" `   // 订阅开始时间
}

type SubscribeStatus uint8

const (
	SubscribeStatusInactive SubscribeStatus = iota
	SubscribeStatusUsing
	SubscribeStatusUnSubscribed SubscribeStatus = 4
)

func (s SubscribeStatus) Value() uint8 {
	return uint8(s)
}
func (s SubscribeStatus) String() string {
	switch s {
	case SubscribeStatusInactive:
		return "未激活"
	case SubscribeStatusUsing:
		return "计费中"
	case SubscribeStatusUnSubscribed:
		return "已退订"
	default:
		return ""
	}
}

// MemberTeamFilter 会员团队筛选条件
type MemberTeamFilter struct {
	Keyword         *string `json:"keyword" query:"keyword"`                 // 关键词
	SubscribeStatus *uint8  `json:"subscribeStatus" query:"subscribeStatus"` // 订阅状态 0:未激活 1:计费中 4:已退订
	Level           *uint8  `json:"level" query:"level" enums:"1,2"`         // 团员层级
	Start           *string `json:"start" query:"start"`                     // 开始日期
	End             *string `json:"end" query:"end"`                         // 结束日期
}

type MemberTeamRows struct {
	ID               uint64         `json:"id" `               // id
	Phone            string         `json:"phone" `            // 手机号
	Name             sql.NullString `json:"name" `             // 姓名
	Level            TemaLevel      `json:"level" `            // 团员层级
	RenewalCount     uint64         `json:"renewalCount" `     // 续费次数
	SubscribeStatus  sql.NullInt64  `json:"subscribeStatus" `  // 订阅状态 0:未激活 1:计费中 4:已退订
	SubscribeStartAt sql.NullTime   `json:"subscribeStartAt" ` // 订阅开始时间
}

// MemberTeamListRes 会员团队列表
type MemberTeamListRes struct {
	model.PaginationRes
}

// MemberTeamStatisticsRes 会员团队概览
type MemberTeamStatisticsRes struct {
	Total            uint64 `json:"total" `            // 总人数
	FirstLevelCount  uint64 `json:"firstLevelCount" `  // 一级团员人数
	SecondLevelCount uint64 `json:"secondLevelCount" ` // 二级团员人数
}

const (
	MemberSigninTypeSms    uint64 = iota + 1 // 短信登录
	MemberSigninTypeWechat                   // 微信授权登录
)

// MemberSigninReq 注册会员请求参数
type MemberSigninReq struct {
	Phone             string  `json:"phone,omitempty" validate:"required_if=SigninType 1"`                         // 手机号
	SmsID             string  `json:"smsID,omitempty" validate:"required_if=SigninType 1"`                         // 短信ID
	Code              string  `json:"code,omitempty" validate:"required_if=SigninType 1,required_if=SigninType 2"` // 图形验证码或者授权登录code
	SigninType        uint64  `json:"signinType" validate:"required,oneof=1 2" enums:"1,2"`                        // 登录类型 1:短信登录 2:微信授权登录
	Name              *string `json:"name"`                                                                        // 姓名
	ReferringMemberID *uint64 `json:"referringMemberID"`                                                           // 推荐人
}

type MemberCreateReq struct {
	Phone             string  `json:"phone"`             // 手机号
	Name              *string `json:"name"`              // 姓名
	RiderID           *uint64 `json:"riderID"`           // 骑手ID
	ReferringMemberID *uint64 `json:"referringMemberID"` // 推荐人ID
	SubscribeID       *uint64 `json:"subscribeID"`       // 订阅ID
}

// MemberSigninRes 注册会员返回参数
type MemberSigninRes struct {
	// token
	Token string `json:"token"`
	// 会员信息
	Profile *MemberProfile `json:"profile"`
}

type MemberProfile struct {
	MemberBaseInfo
	Level          uint64           `json:"level"`          // 会员等级
	AvatarURL      string           `json:"avatarUrl"`      // 头像
	AuthStatusName string           `json:"authStatusName"` // 实名认证状态名称
	AuthStatus     PersonAuthStatus `json:"authStatus"`     // 实名认证状态 0:未认证 1:已认证 2:认证失败
}

// MemberCommissionReq 会员设置返佣
type MemberCommissionReq struct {
	ID     uint64          `json:"id" validate:"required"`
	Rule   *CommissionRule `json:"rule" validate:"required"` // 返佣规则
	Desc   *string         `json:"desc"`                     // 返佣说明
	PlanID []uint64        `json:"planId" validate:"unique"` // 骑士卡方案ID
}

type UploadAvatar struct {
	Avatar string `json:"avatar" validate:"required" ` // 头像
}

type InviteType uint8

const (
	MemberSignSuccess    InviteType = iota + 1 // 注册成功
	MemberBindSuccess                          // 绑定成功
	MemberInviteFail                           // 已被邀请
	MemberActivationFail                       // 已被激活
	MemberInviteSelfFail                       // 自己不能邀请自己
)

type MemberInviteRes struct {
	InviteType InviteType `json:"inviteType" ` // 邀请类型 1:注册成功 2:绑定成功 3:已被邀请 4:已被激活
}

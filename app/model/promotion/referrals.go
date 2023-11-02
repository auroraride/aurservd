package promotion

import (
	"github.com/auroraride/aurservd/app/model"
)

type ReferralsStatus uint8

const (
	ReferralsStatusInviting ReferralsStatus = iota // 邀请中
	ReferralsStatusSuccess                         // 邀请成功
	ReferralsStatusFail                            // 邀请失败
)

func (s ReferralsStatus) Value() uint8 {
	return uint8(s)
}

func (s ReferralsStatus) String() string {
	switch s {
	case ReferralsStatusInviting:
		return "邀请中"
	case ReferralsStatusSuccess:
		return "邀请成功"
	case ReferralsStatusFail:
		return "邀请失败"
	default:
		return ""
	}
}

type Referrals struct {
	ReferringMemberId *uint64         `json:"referringMemberId"` // 推荐人id
	ReferredMemberId  *uint64         `json:"referredMemberId"`  // 被推荐人id
	RiderID           *uint64         `json:"riderId"`           // 被推荐人骑手id
	Status            ReferralsStatus `json:"status"`            // 推荐关系状态 0: 邀请中 1:邀请成功 2:邀请失败
	Name              string          `json:"name"`              // 姓名
	Remark            string          `json:"remark"`            // 备注
}

type ReferralsProgressReq struct {
	model.PaginationReq
	MemberID *uint64 `json:"memberId" query:"memberId" param:"id"`
	ReferralsProgressFilter
}

type ReferralsProgressFilter struct {
	Status  *ReferralsStatus `json:"status" enums:"0,1,2" query:"status"` // 推荐关系状态 0: 邀请中 1:邀请成功 2:邀请失败
	Start   *string          `json:"start" query:"start"`                 // 开始时间
	End     *string          `json:"end" query:"end"`                     // 结束时间
	Keyword *string          `json:"keyword" query:"keyword"`             // 关键词
}

type ReferralsProgressRes struct {
	Name      string `json:"name"`      // 姓名
	Phone     string `json:"phone"`     // 手机号
	Status    uint8  `json:"status"`    // 状态
	CreatedAt string `json:"createdAt"` // 创建时间
	Remark    string `json:"remark"`    // 备注
}

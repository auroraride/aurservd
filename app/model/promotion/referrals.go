package promotion

type Referrals struct {
	ReferringMemberId *uint64 `json:"referringMemberId"` // 推荐人id
	ReferredMemberId  uint64  `json:"referredMemberId"`  // 被推荐人id
	RiderID           *uint64 `json:"riderId"`           // 被推荐人骑手id
}

package definition

// RefundReq 退款申请
type RefundReq struct {
	SubscribeID *uint64 `json:"subscribeId"` // 骑士卡ID
}

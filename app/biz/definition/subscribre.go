package definition

// SubscribeStoreModifyReq 车电套餐修改激活门店
type SubscribeStoreModifyReq struct {
	SubscribeID uint64 `json:"subscribeId" query:"subscribeId" validate:"required"` // 订阅ID
	StoreID     uint64 `json:"storeId" query:"storeId" validate:"required"`         // 门店ID
}

const (
	DepositTypeAlipayAuthFreeze uint8 = iota + 1 // 芝麻免押
	DepositTypeWechatDeposit                     // 微信支付分免押
	DepositTypeContract                          // 合同押金
	DepositTypePay                               // 支付押金
)

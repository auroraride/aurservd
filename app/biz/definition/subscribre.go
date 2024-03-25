package definition

// SubscribeStoreModifyReq 车电套餐修改激活门店
type SubscribeStoreModifyReq struct {
	SubscribeID uint64 `json:"subscribeId" query:"subscribeId" validate:"required"` // 订阅ID
	StoreID     uint64 `json:"storeId" query:"storeId" validate:"required"`         // 门店ID
}

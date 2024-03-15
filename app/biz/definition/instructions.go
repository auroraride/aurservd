package definition

var InstructionsColumns = map[string]string{
	"rent":    "租电买前必读",
	"rentCar": "租车买前必读",
	"buyCar":  "买车买前必读",
	"coupon":  "优惠券使用说明",
	"point":   "积分使用说明",
}

type InstructionsCreateReq struct {
	Content interface{} `json:"content" validate:"required"` // 内容
	Title   string      `json:"title" validate:"required"`   // 标题
	Key     string      `json:"key" validate:"required"`     // 键
}

type InstructionsRes struct {
	Content interface{} `json:"content"` // 内容
	Title   string      `json:"title"`   // 标题
	Key     string      `json:"key"`     // 键  rent 租电买前必读 rentCar 租车买前必读 buyCar 买车买前必读 coupon 优惠券使用说明 point 积分使用说明
}

type InstructionsDetailReq struct {
	Key string `json:"key" validate:"required" param:"key"` // 键 rent 租电买前必读 rentCar 租车买前必读 buyCar 买车买前必读 coupon 优惠券使用说明 point 积分使用说明
}

package model

// EbikeBrandAttribute 属性
type EbikeBrandAttribute struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

// EbikeBrandAttributeCreateReq 创建属性
type EbikeBrandAttributeCreateReq struct {
	BrandAttribute []EbikeBrandAttribute
	BrandID        uint64 `json:"brandId"`
}

// EbikeBrandAttributeUpdateReq 更新属性
type EbikeBrandAttributeUpdateReq struct {
	BrandAttribute []EbikeBrandAttribute
	BrandID        uint64 `json:"brandId"`
}

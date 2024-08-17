package model

// AssetAttributesCreateReq 创建属性请求
type AssetAttributesCreateReq struct {
	AssetType AssetType        `json:"assetType" validate:"required"`  // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	AssetID   uint64           `json:"assetId" validate:"required"`    // 资产ID
	Attribute []AssetAttribute `json:"attributes" validate:"required"` // 属性
}

// AssetAttributesListRes 属性列表返回
type AssetAttributesListRes struct {
	AttributeID   uint64 `json:"attributeId"`             // 属性ID
	AttributeName string `json:"attributeName,omitempty"` // 属性名称
	AttributeKey  string `json:"attributeKey,omitempty"`  // 属性键
}

// AssetAttributesListReq 资产属性列表请求
type AssetAttributesListReq struct {
	AssetType AssetType `json:"assetType" validate:"required" query:"assetType"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
}

// InitAssetAttributes 初始化资产属性
var InitAssetAttributes = []AssetAttributesCreateReq{
	{
		AssetType: AssetTypeEbike,
		Attribute: []AssetAttribute{
			{
				AttributeName: "生产批次",
				AttributeKey:  "exFactory",
			},
			{
				AttributeName: "车牌号",
				AttributeKey:  "plate",
			},
			{
				AttributeName: "终端编号",
				AttributeKey:  "machine",
			},
			{
				AttributeName: "SIM卡",
				AttributeKey:  "sim",
			},
			{
				AttributeName: "颜色",
				AttributeKey:  "color",
			},
		},
	},
}

// AssetAttributeCreate 属性创建
type AssetAttributeCreate struct {
	AttributeID    uint64 `json:"attributeId"`              // 属性ID
	AttributeValue string `json:"attributeValue,omitempty"` // 属性值
}

// AssetAttributeUpdate 属性更新
type AssetAttributeUpdate struct {
	AttributeValue string `json:"attributeValue,omitempty"` // 属性值
	AttributeID    uint64 `json:"attributeId"`              // 属性ID
}

// AssetAttribute 属性返回
type AssetAttribute struct {
	AttributeID      uint64 `json:"attributeId"`                // 属性ID
	AttributeName    string `json:"attributeName,omitempty"`    // 属性名称
	AttributeKey     string `json:"attributeKey,omitempty"`     // 属性键
	AttributeValueID uint64 `json:"attributeValueId,omitempty"` // 属性值ID
	AttributeValue   string `json:"attributeValue,omitempty"`   // 属性值
}

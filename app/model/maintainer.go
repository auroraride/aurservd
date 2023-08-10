package model

type Maintainer struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`   // 姓名
	Enable bool   `json:"enable"` // 是否启用
	Phone  string `json:"phone"`  // 电话
	Cities []City `json:"cities"` // 城市列表
}

type MaintainerListReq struct {
	PaginationReq
	CityID  uint64 `json:"cityId" query:"cityId"`   // 城市
	Keyword string `json:"keyword" query:"keyword"` // 关键词
}

type MaintainerCreateReq struct {
	Name     string   `json:"name" validate:"required" trans:"姓名"`
	Phone    string   `json:"phone" validate:"required" trans:"电话"`
	Password string   `json:"password" validate:"required" trans:"密码"`
	CityIDs  []uint64 `json:"cityIds" validate:"required" trans:"城市列表"`
	Enable   *bool    `json:"enable"` // 是否启用，默认`true`
}

type MaintainerModifyReq struct {
	ID       uint64   `json:"id" param:"id" validate:"required"`
	Name     *string  `json:"name"`     // 姓名
	Phone    *string  `json:"phone"`    // 电话
	Password *string  `json:"password"` // 密码
	CityIDs  []uint64 `json:"cityIDs"`  // 城市列表
	Enable   *bool    `json:"enable"`   // 是否启用
}

type MaintainerSigninReq struct {
	Phone    string `json:"phone" validate:"required" trans:"电话"`
	Password string `json:"password" validate:"required" trans:"密码"`
}

type MaintainerSigninRes struct {
	*Maintainer
	Token string `json:"token"`
}

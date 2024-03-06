package model

import (
	"time"
)

type GuideReq struct {
	Name   string `json:"name"`
	Sort   uint8  `json:"sort"`
	Answer string `json:"answer"`
	Remark string `json:"remark"`
}

type GuideDetail struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Sort      uint8     `json:"sort"`
	Answer    string    `json:"answer"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GuideModifyReq struct {
	IDParamReq
	Name      string    `json:"name"`
	Sort      uint8     `json:"sort"`
	Answer    string    `json:"answer"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GuideListReq struct {
	PaginationReq
}

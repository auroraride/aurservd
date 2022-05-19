// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "math"
)

type PaginationReq struct {
    Current  int `json:"current,omitempty" query:"current"`   // 当前页, 从1开始, 默认1
    PageSize int `json:"pageSize,omitempty" query:"pageSize"` // 每页数据, 默认20
}

type Pagination struct {
    Current int `json:"current"` // 当前页
    Pages   int `json:"pages"`   // 总页数
    Total   int `json:"total"`   // 总条数
}

type PaginationRes struct {
    Pagination Pagination `json:"pagination"` // 分页属性
    Items      any        `json:"items"`      // 返回数据
}

func (p PaginationReq) GetCurrent() int {
    c := p.Current
    if c < 1 {
        c = 1
    }
    return c
}

func (p PaginationReq) GetLimit() int {
    limit := p.PageSize
    if p.PageSize < 1 {
        limit = 20
    }
    return limit
}

func (p PaginationReq) GetOffset() int {
    return (p.GetCurrent() - 1) * p.GetLimit()
}

func (p PaginationReq) GetPages(total int) int {
    return int(math.Ceil(float64(total) / float64(p.GetLimit())))
}

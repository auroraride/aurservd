// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package workwx

type Department struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	NameEn           string   `json:"name_en"`
	DepartmentLeader []string `json:"department_leader"`
	Parentid         int      `json:"parentid"`
	Order            int      `json:"order"`
}

type DepartmentResponse struct {
	baseResponse
	Department []Department `json:"department"`
}

// DepartmentList 部门列表
func (w *Client) DepartmentList() (res DepartmentResponse, err error) {
	err = w.RequestGet("/department/list", &res)
	return res, err
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package workwx

import "fmt"

type UserSimple struct {
    Userid     string `json:"userid"`
    Name       string `json:"name"`
    Department []int  `json:"department"`
    OpenUserid string `json:"open_userid"`
}

type UserSimpleListResponse struct {
    baseResponse
    Userlist []UserSimple `json:"userlist"`
}

func (w *Client) UserSimpleList(departmentID int) (res UserSimpleListResponse, err error) {
    err = w.RequestGet(fmt.Sprintf("/user/simplelist?department_id=%d&fetch_child=1", departmentID), &res)
    return
}

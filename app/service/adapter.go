// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    am "github.com/auroraride/adapter/model"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/go-resty/resty/v2"
)

type adapterService struct {
    *BaseService
    baseUrl string
}

func NewAdapter(baseUrl string, params ...any) *adapterService {
    return &adapterService{
        BaseService: newService(params...),
        baseUrl:     baseUrl,
    }
}

func (s *adapterService) Request() *resty.Request {
    r := resty.New().R()

    switch true {
    default:
        snag.Panic("为找到用户信息")
    case s.rider != nil:
        r.SetHeader(am.HeaderUserID, s.rider.Phone).SetHeader(am.HeaderUserType, am.UserTypeRider.String())
    case s.employee != nil:
        r.SetHeader(am.HeaderUserID, s.employee.Phone).SetHeader(am.HeaderUserType, am.UserTypeEmployee.String())
    case s.modifier != nil:
        r.SetHeader(am.HeaderUserID, s.modifier.Phone).SetHeader(am.HeaderUserType, am.UserTypeManager.String())
    }

    return r
}

func (s *adapterService) Post(url string, playload any) (*resty.Response, error) {
    return s.Request().SetBody(playload).Post(url)
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/adapter"
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
    r := resty.New().SetBaseURL(s.baseUrl).R()

    switch true {
    default:
        snag.Panic("未找到用户信息")
    case s.rider != nil:
        r.SetHeader(adapter.HeaderUserID, s.rider.Phone).SetHeader(adapter.HeaderUserType, adapter.UserTypeRider.String())
    case s.employee != nil:
        r.SetHeader(adapter.HeaderUserID, s.employee.Phone).SetHeader(adapter.HeaderUserType, adapter.UserTypeEmployee.String())
    case s.modifier != nil:
        r.SetHeader(adapter.HeaderUserID, s.modifier.Phone).SetHeader(adapter.HeaderUserType, adapter.UserTypeManager.String())
    }

    return r
}

func (s *adapterService) Post(url string, playload any) (*resty.Response, error) {
    return s.Request().SetBody(playload).Post(url)
}

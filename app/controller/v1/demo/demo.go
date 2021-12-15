// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package demo

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
)

func Esign(c echo.Context) error {
    e := esign.New()
    var req esign.CreatePersonAccountReq
    c.(*app.GlobalContext).BindValidate(&req)
    ai := e.CreatePersonAccount(req)
    log.Info(ai)
    doc := e.DocTemplate()
    log.Println(doc)
    return nil
}
// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package demo

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
)

func Esign(c echo.Context) error {
    e := esign.New()
    var req esign.CreatePersonAccountReq
    c.(*app.Context).BindValidate(&req)
    ai := e.CreatePersonAccount(req)
    log.Info(ai) // f8c1df05edb047f4869a4d24749d76df
    tmpl := e.DocTemplate("8b04852aa1ee4963b10db41668e81f4d")
    m := make(ar.Map)
    for _, com := range tmpl.StructComponents {
        if com.Key != "entSign" && com.Key != "riderSign" {
            m[com.Key] = com.Key
        }
    }
    pdf := e.CreateByTemplate(esign.CreateByTemplateReq{
        Name:             "签署测试.pdf",
        SimpleFormFields: m,
        TemplateId:       tmpl.TemplateId,
    })
    return app.NewResponse(c).Success().SetData(pdf).Send()
}

func EsignDo(c echo.Context) error {
    snag.Panic("test")
    log.Info(esign.New().CreateFlowOneStep())
    return nil
}

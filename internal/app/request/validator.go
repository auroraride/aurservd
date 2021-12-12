// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package request

import (
    "errors"
    "github.com/auroraride/aurservd/internal/ar"
    zhLocales "github.com/go-playground/locales/zh"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    zhTranslations "github.com/go-playground/validator/v10/translations/zh"
    "github.com/iancoleman/strcase"
    "reflect"
)

type GlobalValidator struct {
    validator *validator.Validate
    trans     ut.Translator
}

func (v *GlobalValidator) Validate(i interface{}) error {
    err := v.validator.Struct(i)
    if err != nil {
        errs := err.(validator.ValidationErrors)
        msg := ""
        for _, e := range errs {
            msg += e.Translate(v.trans)
        }
        err = errors.New(msg)
    }
    return err
}

func NewGlobalValidator() *GlobalValidator {
    val := validator.New()
    zh := zhLocales.New()
    uni := ut.New(zh, zh)
    trans, _ := uni.GetTranslator("zh")
    _ = zhTranslations.RegisterDefaultTranslations(val, trans)
    val.RegisterTagNameFunc(func(field reflect.StructField) string {
        t := field.Tag.Get("trans")
        m := ar.Config.Trans
        if t != "" {
            return t
        }
        if str, ok := m[strcase.ToLowerCamel(field.Name)]; ok {
            return str
        }
        return field.Name
    })
    return &GlobalValidator{validator: val, trans: trans}
}

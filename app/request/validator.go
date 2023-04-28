// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package request

import (
	"errors"
	"reflect"
	"strings"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/utils"
	zhLocales "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/iancoleman/strcase"
)

type GlobalValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

type RegisterValidationFunc func(fn validator.Func)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func (v *GlobalValidator) Validate(i any) error {
	err := v.validator.Struct(i)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		var msg []string
		for _, e := range errs {
			msg = append(msg, e.Translate(v.trans))
		}
		err = errors.New(strings.Join(msg, "，"))
	}
	return err
}

// NewGlobalValidator 校验方法
func NewGlobalValidator() *GlobalValidator {
	validate = validator.New()
	zh := zhLocales.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)

	// 校验手机号
	customValidation("phone")(func(fl validator.FieldLevel) bool {
		return utils.NewRegex().MatchPhone(fl.Field().String())
	})

	// 校验枚举
	customValidation("enum")(validateEnum)

	// 从字段tag中获取字段翻译
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
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

	_ = validate.RegisterTranslation("excluded_with", trans, func(ut ut.Translator) error {
		return nil
	}, func(ut ut.Translator, fe validator.FieldError) string {
		return fe.Field() + "为排除字段"
	})

	_ = validate.RegisterTranslation("excluded_if", trans, func(ut ut.Translator) error {
		return nil
	}, func(ut ut.Translator, fe validator.FieldError) string {
		return fe.Field() + "为排除字段"
	})

	return &GlobalValidator{validator: validate, trans: trans}
}

// customValidation 自定义验证规则
func customValidation(tag string, message ...string) RegisterValidationFunc {
	return func(fn validator.Func) {
		_ = validate.RegisterValidation(tag, fn)
		_ = validate.RegisterTranslation(
			tag,
			trans,
			func(ut ut.Translator) error {
				text := "{0}验证失败"
				if len(message) > 0 {
					text = message[0]
				}
				return ut.Add(tag, text, true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(tag, fe.Field())
				return t
			},
		)
	}
}

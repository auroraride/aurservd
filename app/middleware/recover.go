// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					// buf := make([]byte, 1<<20)
					// stacklen := runtime.Stack(buf, true)
					// fmt.Printf("---------\n捕获错误: %v\n%s\n\n", r, buf[:stacklen])
					switch v := r.(type) {
					case *snag.Error:
						c.Error(r.(*snag.Error))
					case *ent.ValidationError:
						zap.L().Error("捕获错误: ent.ValidationError", zap.Error(v), zap.Stack("stack"))
						c.Error(r.(*ent.ValidationError).Unwrap())
					case error:
						zap.L().Error("捕获错误", zap.Error(v), zap.Stack("stack"))
						c.Error(r.(error))
					default:
						x := fmt.Errorf("%v", r)
						zap.L().Error("捕获错误: 其他", zap.Error(x), zap.Stack("stack"))
						// _ = mw.Recover()(next)(c)
						c.Error(x)
					}
				}
			}()
			return next(c)
		}
	}
}

// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/internal/ar"
    bc "github.com/mojocn/base64Captcha"
    "math/rand"
    "time"
)

type captcha struct {
    driver *bc.DriverDigit
    ctx    context.Context
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

func NewCaptcha() *captcha {
    d := bc.NewDriverDigit(92, 240, 6, 0.85, 120)
    return &captcha{driver: d, ctx: context.Background()}
}

// Set 存储验证码 10分钟有效期
func (c *captcha) Set(id string, value string) error {
    return ar.Cache.Set(c.ctx, id, value, 600*time.Second).Err()
}

// Get 从Redis中获取验证码
func (c *captcha) Get(id string) (code string) {
    code = ar.Cache.Get(c.ctx, id).Val()
    return
}

// Verify 校验Captcha
// clear 验证成功后是否删除缓存
func (c *captcha) Verify(id, answer string, clear bool) (ok bool) {
    ok = c.Get(id) == answer
    if ok && clear {
        ar.Cache.Del(c.ctx, id)
    }
    return
}

// DrawCaptcha 生成Captcha Item
func (c *captcha) DrawCaptcha() (id string, item bc.Item, err error) {
    id, content, answer := c.driver.GenerateIdQuestionAnswer()
    item, err = c.driver.DrawCaptcha(content)
    if err != nil {
        return
    }
    if err = c.Set(id, answer); err != nil {
        return
    }
    return
}

// Generate 生成Captcha base64图片
func (c *captcha) Generate() (id, b64s string, err error) {
    var item bc.Item
    id, item, err = c.DrawCaptcha()
    if err != nil {
        return
    }
    b64s = item.EncodeB64string()
    return
}

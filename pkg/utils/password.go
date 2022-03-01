// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import "golang.org/x/crypto/bcrypt"

// PasswordGenerate 生成密码
func PasswordGenerate(password string) (hash string, err error) {
    var b []byte
    b, err = bcrypt.GenerateFromPassword([]byte(password), 8)
    if err != nil {
        return
    }

    return string(b), nil
}

// PasswordCompare 对比密码
func PasswordCompare(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

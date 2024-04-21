// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package model

import "database/sql/driver"

type AppPlatform string

const (
	AppPlatformAndroid AppPlatform = "android"
	AppPlatformIOS     AppPlatform = "ios"
)

func (b *AppPlatform) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		return nil
	case string:
		*b = AppPlatform(v)
	}
	return nil
}

func (b AppPlatform) Value() (driver.Value, error) {
	return string(b), nil
}

// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package model

import "database/sql/driver"

type Platform string

const (
	PlatformUnknown Platform = ""
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
)

func (b *Platform) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		return nil
	case string:
		*b = Platform(v)
	}
	return nil
}

func (b Platform) Value() (driver.Value, error) {
	return b, nil
}

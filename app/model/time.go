// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-10
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/auroraride/aurservd/pkg/tools"
	"github.com/golang-module/carbon/v2"
)

type Date struct {
	time.Time
}

func DateNow() Date {
	return DateFromTime(time.Now())
}

func DateFromTime(t time.Time) Date {
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return Date{Time: d}
}

func DateFromStringX(str string) Date {
	return DateFromTime(tools.NewTime().ParseDateStringX(str))
}

func (t *Date) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateFromTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to date", v)
}

func (t *Date) String() string {
	return t.Format(carbon.DateLayout)
}

func (t Date) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t Date) Carbon() carbon.Carbon {
	return carbon.Time2Carbon(t.Time)
}

func (t Date) Tomorrow() Date {
	return Date{Time: t.AddDate(0, 0, 1)}
}

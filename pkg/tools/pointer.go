// Copyright (C) liasica. 2022-present.
//
// Created at 2022-04-30
// Based on api by liasica, magicrolan@qq.com.

package tools

import "time"

type pointer struct {
}

func NewPointer() *pointer {
    return &pointer{}
}

// Time 复制 time.Time 对象，并返回复制体的指针
func (*pointer) Time(t time.Time) *time.Time {
    return &t
}

// String 复制 string 对象，并返回复制体的指针
func (*pointer) String(s string) *string {
    return &s
}

// Bool 复制 bool 对象，并返回复制体的指针
func (*pointer) Bool(b bool) *bool {
    return &b
}

// Float64 复制 float64 对象，并返回复制体的指针
func (*pointer) Float64(f float64) *float64 {
    return &f
}

// Float32 复制 float32 对象，并返回复制体的指针
func (*pointer) Float32(f float32) *float32 {
    return &f
}

// UInt64 复制 uint64 对象，并返回复制体的指针
func (*pointer) UInt64(i uint64) *uint64 {
    return &i
}

// Int 复制 int 对象，并返回复制体的指针
func (*pointer) Int(i int) *int {
    return &i
}

// Int64 复制 int64 对象，并返回复制体的指针
func (*pointer) Int64(i int64) *int64 {
    return &i
}

// Int32 复制 int64 对象，并返回复制体的指针
func (*pointer) Int32(i int32) *int32 {
    return &i
}

func Pointer[T any](i T) *T {
    return &i
}

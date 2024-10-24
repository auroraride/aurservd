// Copyright (C) liasica. 2022-present.
//
// Created at 2022-04-30
// Based on api by liasica, magicrolan@qq.com.

package silk

import "time"

// Time 复制 time.Time 对象，并返回复制体的指针
func Time(t time.Time) *time.Time {
	return &t
}

// String 复制 string 对象，并返回复制体的指针
func String(s string) *string {
	return &s
}

// Bool 复制 bool 对象，并返回复制体的指针
func Bool(b bool) *bool {
	return &b
}

// Float64 复制 float64 对象，并返回复制体的指针
func Float64(f float64) *float64 {
	return &f
}

// Float32 复制 float32 对象，并返回复制体的指针
func Float32(f float32) *float32 {
	return &f
}

// UInt64 复制 uint64 对象，并返回复制体的指针
func UInt64(i uint64) *uint64 {
	return &i
}

// Int 复制 int 对象，并返回复制体的指针
func Int(i int) *int {
	return &i
}

// UInt 复制 uint 对象，并返回复制体的指针
func UInt(i uint) *uint {
	return &i
}

// Int64 复制 int64 对象，并返回复制体的指针
func Int64(i int64) *int64 {
	return &i
}

// Int32 复制 int64 对象，并返回复制体的指针
func Int32(i int32) *int32 {
	return &i
}

func Pointer[T any](i T) *T {
	return &i
}

func Compare[T comparable](a *T, b *T) bool {
	return (a != nil && b != nil) || (a == nil && b == nil)
}

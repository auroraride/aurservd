// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-01
// Based on aurservd by liasica, magicrolan@qq.com.

package types

import (
	"fmt"
	"strconv"
)

type Value interface {
	RawValue() any
}

func Uint8(t interface{}) (p uint8, err error) {
	switch x := t.(type) {
	case nil:
	case uint8:
		p = x
	case int64:
		p = uint8(x)
	case int, int8, int16, int32, uint, uint16, uint32, uint64, uintptr, float32, float64:
		var n uint64
		n, err = strconv.ParseUint(fmt.Sprintf("%x", x), 10, 64)
		if err != nil {
			return
		}
		p = uint8(n)
	default:
		err = fmt.Errorf("unexpected type %T", t)
	}
	return
}

// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-28
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import "fmt"

// DataSum 校验和
func DataSum(data []byte) string {
	var (
		sum    uint32
		length = len(data)
	)

	for i := 0; i < length; i++ {
		sum += uint32(data[i])
		if sum > 0xff {
			sum = ^sum
			sum += 1
		}
	}
	sum &= 0xff
	return fmt.Sprintf("%X", sum)
}

// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-18
// Based on aurservd by liasica, magicrolan@qq.com.

package definition

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPersonIdentity(t *testing.T) {
	p := NewPersonIdentity("11010119900101100X", "张三")
	str := p.Pack()
	t.Logf("PACKED: %s", str)

	up := new(PersonIdentity)
	err := up.UnPack(str)
	require.NoError(t, err)

	require.Equal(t, *p, *up)
}

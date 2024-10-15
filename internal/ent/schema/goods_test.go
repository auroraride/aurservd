// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-10-15, by liasica

package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortInstallment(t *testing.T) {
	tests := []struct {
		input [][]float64
		want  [][]float64
		error string
	}{
		{
			input: [][]float64{
				{1, 2, 3},
				{1, 2},
				{1},
				{1, 2, 3, 4},
			},
			want: [][]float64{
				{1},
				{1, 2},
				{1, 2, 3},
				{1, 2, 3, 4},
			},
			error: "",
		},
		{
			input: [][]float64{
				{1, 2, 3},
				{1, 2},
				{1, 2},
				{1},
			},
			want:  nil,
			error: "分期方案重复",
		},
	}

	for _, tt := range tests {
		err := SortInstallment(tt.input)
		if err != nil {
			require.EqualError(t, err, tt.error)
			continue
		}
		if !reflect.DeepEqual(tt.input, tt.want) {
			t.Errorf("SortInstallment(%v) = %v, want %v", tt.input, tt.input, tt.want)
		}
	}
}

// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-10-15, by liasica

package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	tests := []struct {
		input GoodsPaymentPlans
		want  GoodsPaymentPlans
		error error
	}{
		{
			input: GoodsPaymentPlans{
				{
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
				},
				{
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
				},
			},
			want: GoodsPaymentPlans{
				{
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
				},
				{
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
				},
			},
			error: nil,
		},
		{
			input: GoodsPaymentPlans{{}},
			want:  nil,
			error: ErrorGoodsPaymentEmpty,
		},
		{
			input: GoodsPaymentPlans{
				{
					{Period: GoodsPaymentPeriodOnce, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
				},
			},
			want:  nil,
			error: ErrorGoodsPaymentInvalid,
		},
		{
			input: GoodsPaymentPlans{
				{
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 0},
				},
			},
			want:  nil,
			error: ErrorGoodsPaymentAmount,
		},
		{
			input: GoodsPaymentPlans{
				{
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
				},
				{
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
					{Period: GoodsPaymentPeriodDaily, Unit: 1, Amount: 100},
				},
			},
			want:  nil,
			error: ErrorGoodsPaymentDuplicate,
		},
	}

	for _, tt := range tests {
		err := tt.input.Valid()
		if err != nil {
			require.ErrorAs(t, err, &tt.error)
			continue
		} else {
			if !reflect.DeepEqual(tt.input, tt.want) {
				t.Errorf("Valid(%v) = %v, want %v", tt.input, tt.input, tt.want)
			}
		}
	}
}

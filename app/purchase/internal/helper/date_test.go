// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-10-21, by liasica

package helper

import (
	"testing"
	"time"
)

func TestOverdueDays(t *testing.T) {
	tests := []struct {
		t1   time.Time
		t2   time.Time
		want int
	}{
		{
			t1:   time.Date(2024, 10, 21, 0, 8, 0, 0, time.Local),
			t2:   time.Date(2024, 10, 22, 0, 0, 0, 0, time.Local),
			want: 1,
		},
		{
			t1:   time.Date(2024, 10, 21, 23, 59, 58, 0, time.Local),
			t2:   time.Date(2024, 10, 28, 10, 0, 0, 0, time.Local),
			want: 7,
		},
	}

	for _, tt := range tests {
		got := OverdueDays(tt.t1, tt.t2)
		t.Logf("OverdueDays(%v, %v) = %v", tt.t1, tt.t2, got)
		if got != tt.want {
			t.Errorf("OverdueDays(%v, %v) = %v, want %v", tt.t1, tt.t2, got, tt.want)
		}
	}
}

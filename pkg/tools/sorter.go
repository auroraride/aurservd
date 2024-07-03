// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-02, by liasica

package tools

import (
	"reflect"
	"sort"
)

type Sorter struct{ keys []Key }

func NewSorter() *Sorter { return new(Sorter) }

func (l *Sorter) AddStr(key StringKey) *Sorter  { l.keys = append(l.keys, key); return l }
func (l *Sorter) AddInt(key IntKey) *Sorter     { l.keys = append(l.keys, key); return l }
func (l *Sorter) AddFloat(key FloatKey) *Sorter { l.keys = append(l.keys, key); return l }

func (l *Sorter) SortStable(slice any) {
	value := reflect.ValueOf(slice)
	sort.SliceStable(slice, func(i, j int) bool {
		si := value.Index(i).Interface()
		sj := value.Index(j).Interface()
		for _, key := range l.keys {
			if key.Less(si, sj) {
				return true
			}
			if key.Less(sj, si) {
				return false
			}
		}
		return false
	})
}

type Key interface {
	Less(a, b any) bool
}

type StringKey func(any) string

func (k StringKey) Less(a, b any) bool { return k(a) < k(b) }

type IntKey func(any) int

func (k IntKey) Less(a, b any) bool { return k(a) < k(b) }

type FloatKey func(any) float64

func (k FloatKey) Less(a, b any) bool { return k(a) < k(b) }

// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"reflect"
	"testing"

	"github.com/ctx42/testing/internal/affirm"
	"github.com/ctx42/testing/pkg/testcases"
)

func Test_Equal_pair_matrix_tabular(t *testing.T) {
	type hidden struct {
		items []int
		dict  map[string]int
	}
	type boxed struct{ Val any }

	vals := []struct {
		desc string
		make func() any
	}{
		{"nil slice", func() any { return []int(nil) }},
		{"empty slice", func() any { return make([]int, 0) }},
		{"slice 0", func() any { return []int{0} }},
		{"slice 1", func() any { return []int{1} }},
		{"slice 0 0", func() any { return []int{0, 0} }},
		{"nil map", func() any { return map[string]int(nil) }},
		{"empty map", func() any { return make(map[string]int) }},
		{"map a 0", func() any { return map[string]int{"a": 0} }},
		{"map a 1", func() any { return map[string]int{"a": 1} }},
		{"map b 1", func() any { return map[string]int{"b": 1} }},
		{"nil zero slice", func() any { return []struct{}(nil) }},
		{"empty zero slice", func() any { return make([]struct{}, 0) }},
		{"one zero", func() any { return make([]struct{}, 1) }},
		{"two zero", func() any { return make([]struct{}, 2) }},
		{"array len 0", func() any { return [0]int{} }},
		{"array 0", func() any { return [1]int{0} }},
		{"array 1", func() any { return [1]int{1} }},
		{"nil int ptr", func() any { return (*int)(nil) }},
		{"int ptr 1", func() any { val := 1; return &val }},
		{"int ptr 2", func() any { val := 2; return &val }},
		{"nil zero ptr", func() any { return (*struct{})(nil) }},
		{"zero ptr", func() any { return new(struct{}) }},
		{"hidden nil", func() any { return hidden{} }},
		{
			"hidden empty slice",
			func() any { return hidden{items: make([]int, 0)} },
		},
		{"hidden slice", func() any { return hidden{items: []int{1}} }},
		{
			"hidden empty map",
			func() any { return hidden{dict: make(map[string]int)} },
		},
		{
			"hidden map a",
			func() any { return hidden{dict: map[string]int{"a": 1}} },
		},
		{
			"hidden map b",
			func() any { return hidden{dict: map[string]int{"b": 1}} },
		},
		{"boxed nil", func() any { return boxed{} }},
		{"boxed nil slice", func() any { return boxed{[]int(nil)} }},
		{"boxed empty slice", func() any { return boxed{make([]int, 0)} }},
		{"boxed int", func() any { return boxed{1} }},
		{"boxed int64", func() any { return boxed{int64(1)} }},
	}

	for i, val0 := range vals {
		for j, val1 := range vals {
			t.Run(val0.desc+" vs "+val1.desc, func(t *testing.T) {
				// --- Given ---
				want, have := val0.make(), val1.make()

				// --- When ---
				err := Equal(want, have)

				// --- Then ---
				affirm.Equal(t, i == j, err == nil)
			})
		}
	}
}

func Test_Equal_ZENValues_matrix_tabular(t *testing.T) {
	vals := testcases.ZENValues()
	for _, val0 := range vals {
		for _, val1 := range vals {
			t.Run(val0.Desc+" vs "+val1.Desc, func(t *testing.T) {
				// --- When ---
				err := Equal(val0.Val, val1.Val)

				// --- Then ---
				want := reflect.DeepEqual(val0.Val, val1.Val)
				affirm.Equal(t, want, err == nil)
			})
		}
	}
}

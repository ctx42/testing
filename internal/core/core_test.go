// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package core

import (
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"github.com/ctx42/testing/internal/cases"
	"github.com/ctx42/testing/internal/types"
)

const expMsg = "expected same:\n  want: %v\n  have: %v"

func Test_IsNil_tabular_ZENValues(t *testing.T) {
	for _, tc := range cases.ZENValues() {
		t.Run("Nil "+tc.Desc, func(t *testing.T) {
			// --- When ---
			hNil, hWrapped := IsNil(tc.Val)

			// --- Then ---
			if tc.IsNil && !hNil {
				msg := "expected nil value:\n  have: %#v"
				t.Errorf(msg, hNil)
			}
			if !tc.IsNil && hNil {
				t.Errorf("expected not-nil value:\n  have: %#v", hNil)
			}
			if tc.IsWrappedNil != hWrapped {
				t.Errorf("expected wrapped nil value:\n  have: %#v", tc.Val)
			}
		})
	}
}

func Test_WillPanic(t *testing.T) {
	t.Run("panicked", func(t *testing.T) {
		// --- Given ---
		fn := func() { panic("panic") }

		// --- When ---
		val, stack := WillPanic(fn)

		// --- Then ---
		if val.(string) != "panic" {
			t.Error("expected WillPanic to return value 'panic'")
		}
		if stack == "" {
			t.Error("expected WillPanic to return stack trace")
		}
	})

	t.Run("panicked with nil", func(t *testing.T) {
		// --- Given ---
		fn := func() { panic(nil) } // nolint: govet

		// --- When ---
		val, stack := WillPanic(fn)

		// --- Then ---
		if _, ok := val.(*runtime.PanicNilError); !ok {
			t.Error("expected WillPanic to return value 'panic'")
		}
		if stack == "" {
			t.Error("expected WillPanic to return stack trace")
		}
	})

	t.Run("no panic", func(t *testing.T) {
		// --- Given ---
		fn := func() {}

		// --- When ---
		val, stack := WillPanic(fn)

		// --- Then ---
		if val != nil {
			t.Error("expected WillPanic to return an empty string")
		}
		if stack != "" {
			t.Error("expected WillPanic to return an empty string")
		}
	})
}

func Test_Same_tabular(t *testing.T) {
	// Pointers.
	ptr0 := &types.TPtr{Val: "0"}
	ptr1 := &types.TPtr{Val: "1"}

	// Interfaces.
	var itfPtr0, itfPtr1 types.TItf
	itfPtr0, itfPtr1 = &types.TPtr{Val: "0"}, &types.TPtr{Val: "1"}

	// Functions.
	fn0 := func() {}
	fn1 := func() {}
	type TFn0 func()
	type TFn1 func()
	var fnNil0 TFn0
	var fnNil1 TFn1
	var fnA, fnB TFn0

	// Slices.
	s0 := []int{1, 2, 3}
	s1 := []int{1, 2, 3}
	var sNil0 []int
	var sNil1 []float64
	type Ts []int
	var sA, sB Ts

	// Arrays.
	a0 := [2]int{1, 2}
	a1 := [2]int{1, 2}
	var aNil0 [2]int
	var aNil1 [2]float64
	type Ta []int
	var aA, aB Ta

	// Maps.
	m0 := map[string]int{"a": 1, "b": 2}
	m1 := map[string]int{"a": 1, "b": 2}
	var mNil0 map[string]int
	var mNil1 map[string]float64
	type Tm map[string]int
	var mA, mB Tm

	// Channels.
	c0 := make(chan int)
	c1 := make(chan int)
	var cNil0 chan int
	type Tc chan int
	var cNilA chan byte
	var cA, cB Tc

	tt := []struct {
		testN string

		want any
		have any
		same bool
	}{
		{"ptr same", ptr0, ptr0, true},
		{"ptr not same", ptr0, ptr1, false},
		{"ptr different types not same", &types.TPtr{}, &types.TVal{}, false},
		{"prt nil both", nil, nil, false},

		{"itf ptr same", itfPtr0, itfPtr0, true},
		{"itf ptr not same", itfPtr0, itfPtr1, false},
		{"obj val not same", types.TVal{}, types.TVal{}, false},

		{"func same", fn0, fn0, true},
		{"func not same", fn0, fn1, false},
		{"func nil want", nil, fn1, false},
		{"func nil type want", fnNil0, fn1, false},
		{"func nil have", fn0, nil, false},
		{"func nil type have", fn0, fnNil0, false},
		{"func nil type both", fnNil0, fnNil0, true},
		{"func nil different type both", fnNil0, fnNil1, false},
		{"func nil same type both", fnA, fnB, true},

		{"slice same", s0, s0, true},
		{"slice not same", s0, s1, false},
		{"slice nil want", nil, s1, false},
		{"slice nil var want", sNil0, s1, false},
		{"slice nil have", s0, nil, false},
		{"slice nil var have", s0, sNil0, false},
		{"slice nil var both", sNil0, sNil0, true},
		{"slice nil different type both", sNil0, sNil1, false},
		{"slice nil same type both", sA, sB, true},

		{"array same", a0, a0, false},
		{"array not same", a0, a1, false},
		{"array nil want", nil, a1, false},
		{"array nil var want", aNil0, a1, false},
		{"array nil have", a0, nil, false},
		{"array nil var have", a0, aNil0, false},
		{"array nil var both", aNil0, aNil0, false},
		{"array nil different type both", aNil0, aNil1, false},
		{"array nil same type both", aA, aB, true},

		{"map same", m0, m0, true},
		{"map not same", m0, m1, false},
		{"map nil want", nil, s1, false},
		{"map nil var want", mNil0, m1, false},
		{"map nil have", m0, nil, false},
		{"map nil var have", m0, mNil0, false},
		{"map nil both", nil, nil, false},
		{"map nil var both", mNil0, mNil0, true},
		{"map nil different type both", mNil0, mNil1, false},
		{"map nil same type both", mA, mB, true},

		{"chanel same", c0, c0, true},
		{"chanel not the same", c0, c1, false},
		{"chanel nil want", nil, c1, false},
		{"chanel nil var want", mNil0, c1, false},
		{"chanel nil have", c0, nil, false},
		{"chanel nil var have", c0, cNil0, false},
		{"chanel nil both", nil, nil, false},
		{"chanel nil var both", cNil0, cNil0, true},
		{"chanel nil different type both", cNil0, cNilA, false},
		{"chanel nil same type both", cA, cB, true},
	}

	wMsg := "expected same:\n  want: %t\n  have: %t"

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := Same(tc.want, tc.have)

			// --- Then ---
			if !reflect.DeepEqual(tc.same, have) {
				t.Errorf(wMsg, have, tc.same)
			}
		})
	}
}

func Test_Pointer_tabular(t *testing.T) {
	s2 := []*types.TA{{TAp: &types.TA{Int: 42}}}

	tt := []struct {
		testN string

		in  reflect.Value
		nil bool
	}{
		{"nil", reflect.ValueOf(nil), true},
		{"struct value", reflect.ValueOf(struct{}{}), true},
		{"nil struct pointer", reflect.ValueOf((*types.TIntStr)(nil)), true},
		{
			"struct pointer",
			reflect.ValueOf(&struct{ Int int }{Int: 123}),
			false,
		},
		{
			"struct pointer",
			reflect.ValueOf(s2).Index(0),
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := Pointer(tc.in)

			// --- Then ---
			if tc.nil && have != nil {
				t.Error("expected nil value")
			}
			if !tc.nil && have == nil {
				t.Error("expected not-nil value")
			}
		})
	}
}

func Test_Value(t *testing.T) {
	t.Run("struct by pointer - with exported fields", func(t *testing.T) {
		// --- Given ---
		ps := &types.TIntStr{Int: 42, Str: "abc"}
		val := reflect.ValueOf(ps)

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		have := haveVal.(*types.TIntStr)
		if !reflect.DeepEqual(have, &types.TIntStr{Int: 42, Str: "abc"}) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by pointer - exported field", func(t *testing.T) {
		// --- Given ---
		ps := &types.TIntStr{Int: 42, Str: "abc"}
		val := reflect.ValueOf(ps).Elem().FieldByName("Int")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		have := haveVal.(int)
		if have != 42 {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("function", func(t *testing.T) {
		// --- Given ---
		add := func(a, b int) int { return a + b }
		val := reflect.ValueOf(add)

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		if haveVal == nil {
			t.Error("expected non-nil value")
		}
		fn, ok := haveVal.(func(int, int) int)
		if !ok {
			t.Errorf("unexpected type")
		}
		have := fn(1, 2)
		if have != 3 {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("slice of slices of any interfaces", func(t *testing.T) {
		// --- Given ---
		s := [][]any{{"str"}}
		val := reflect.ValueOf(s).Index(0).Index(0)

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		if haveVal.(string) != "str" {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private int field", func(t *testing.T) {
		// --- Given ---
		prv := types.NewTPrv().SetInt(42)
		val := reflect.ValueOf(prv).FieldByName("vInt")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		if haveVal.(int) != 42 {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private pointer field", func(t *testing.T) {
		// --- Given ---
		prv := types.NewTPrv().SetPtr(&types.TVal{Val: "abc"})
		val := reflect.ValueOf(prv).FieldByName("ptr")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		have := haveVal.(*types.TVal)
		if have.Val != "abc" {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private slice field", func(t *testing.T) {
		// --- Given ---
		prv := types.NewTPrv().SetSInt([]int{1, 2, 3})
		val := reflect.ValueOf(prv).FieldByName("sInt")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		have := haveVal.([]int)
		if !reflect.DeepEqual(have, []int{1, 2, 3}) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private array field", func(t *testing.T) {
		// --- Given ---
		prv := types.NewTPrv().SetAInt([...]int{1, 2})
		val := reflect.ValueOf(prv).FieldByName("aInt")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		have := haveVal.([2]int)
		if !reflect.DeepEqual(have, [...]int{1, 2}) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private map field", func(t *testing.T) {
		// --- Given ---
		prv := types.NewTPrv().SetMapII(map[int]int{1: 11, 2: 22})
		val := reflect.ValueOf(prv).FieldByName("vMap")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		have := haveVal.(map[int]int)
		if !reflect.DeepEqual(have, map[int]int{1: 11, 2: 22}) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private time field", func(t *testing.T) {
		// --- Given ---
		tim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		prv := types.NewTPrv().SetTim(tim)
		val := reflect.ValueOf(prv).FieldByName("tim")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		if !haveVal.(time.Time).Equal(tim) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private function field", func(t *testing.T) {
		// --- Given ---
		fn := func() int { return 42 }
		prv := types.NewTPrv().SetFn(fn)
		val := reflect.ValueOf(prv).FieldByName("fn")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		if Same(haveVal, fn) != true {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private chan field", func(t *testing.T) {
		// --- Given ---
		ch := make(chan int)
		prv := types.NewTPrv().SetCh(ch)
		val := reflect.ValueOf(prv).FieldByName("ch")

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		if Same(haveVal, ch) != true {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("interface", func(t *testing.T) {
		// --- Given ---
		var itf types.TItf
		itf = &types.TPtr{}
		val := reflect.ValueOf(itf)

		// --- When ---
		haveVal, haveOK := Value(val)

		// --- Then ---
		if !haveOK {
			t.Error("expected success")
		}
		if Same(haveVal, itf) != true {
			t.Errorf("expected correct field value")
		}
	})
}

func Test_Value_tabular_success(t *testing.T) {
	chn := make(chan int)
	m := make(map[string]int)
	ptr := &types.TPtr{}

	tt := []struct {
		testN string

		in      any
		wantVal any
	}{
		{"nil", nil, nil},
		{"bool - true", true, true},
		{"bool - false", false, false},
		{"int", 42, 42},
		{"int8", int8(42), int8(42)},
		{"int16", int16(42), int16(42)},
		{"int32", int32(42), int32(42)},
		{"int64", int64(42), int64(42)},
		{"uint", uint(42), uint(42)},
		{"uint8", uint8(42), uint8(42)},
		{"uint16", uint16(42), uint16(42)},
		{"uint32", uint32(42), uint32(42)},
		{"uint64", uint64(42), uint64(42)},
		{"float32", float32(42), float32(42)},
		{"float64", float64(42), float64(42)},
		{"complex64", complex64(42), complex64(42)},
		{"complex128", complex128(42), complex128(42)},
		{"array", [...]int{1, 2, 3}, [...]int{1, 2, 3}},
		{"chan", chn, chn},
		{"map", m, m},
		{"pointer", ptr, ptr},
		{"slice", []int{1, 2, 3}, []int{1, 2, 3}},
		{"string", "abc", "abc"},
		{"struct", types.TPtr{}, types.TPtr{}},

		{"uintptr", uintptr(42), uintptr(42)},
		{"unsafe pointer", unsafe.Pointer(ptr), unsafe.Pointer(ptr)},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			haveVal, haveOK := Value(reflect.ValueOf(tc.in))

			// --- Then ---
			if !reflect.DeepEqual(tc.wantVal, haveVal) {
				t.Errorf(expMsg, tc.wantVal, haveVal)
			}
			if !haveOK {
				t.Error("expected true")
			}
		})
	}
}

func Test_value(t *testing.T) {
	t.Run("nil vHeader data pointer", func(t *testing.T) {
		// --- Given ---
		typ := reflect.TypeOf(123)
		val := reflect.ValueOf(nil)

		// --- When ---
		haveVal, haveOK := value(typ, val)

		// --- Then ---
		if haveOK {
			t.Error("expected false")
		}
		if haveVal != nil {
			t.Errorf("expected nil")
		}
	})
}

func Test_ValueSimple_tabular(t *testing.T) {
	chn := make(chan int)
	m := make(map[string]int)
	ptr := &types.TPtr{}

	tt := []struct {
		testN string

		in      any
		wantVal any
		wantOK  bool
	}{
		{"bool - false", false, false, true},
		{"bool - true", false, false, true},
		{"int", 42, 42, true},
		{"int8", int8(42), int8(42), true},
		{"int16", int16(42), int16(42), true},
		{"int32", int32(42), int32(42), true},
		{"int64", int64(42), int64(42), true},
		{"uint", uint(42), uint(42), true},
		{"uint8", uint8(42), uint8(42), true},
		{"uint16", uint16(42), uint16(42), true},
		{"uint32", uint32(42), uint32(42), true},
		{"uint64", uint64(42), uint64(42), true},
		{"uintptr", uintptr(42), nil, false},
		{"float32", float32(42), float32(42), true},
		{"float64", float64(42), float64(42), true},
		{"complex64", complex64(42), complex64(42), true},
		{"complex128", complex128(42), complex128(42), true},

		{"nil", nil, nil, false},
		{"array", [...]int{1, 2, 3}, nil, false},
		{"chan", chn, nil, false},
		{"map", m, nil, false},
		{"pointer", ptr, nil, false},
		{"slice", []int{1, 2, 3}, nil, false},
		{"string", "abc", "abc", true},
		{"struct", types.TPtr{}, nil, false},
		{"unsafe pointer", uintptr(unsafe.Pointer(ptr)), nil, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			haveVal, haveOK := ValueSimple(reflect.ValueOf(tc.in))

			// --- Then ---
			if tc.wantVal != haveVal {
				t.Errorf(expMsg, tc.wantVal, haveVal)
			}
			if tc.wantOK != haveOK {
				t.Errorf(expMsg, tc.wantOK, haveOK)
			}
		})
	}
}

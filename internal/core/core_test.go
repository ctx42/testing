// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package core

import (
	"io"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"github.com/ctx42/testing/pkg/testcases"
)

const expMsg = "expected same:\n  want: %v\n  have: %v"

func Test_IsNil_tabular(t *testing.T) {
	for _, tc := range testcases.ZENValues() {
		t.Run("Nil "+tc.Desc, func(t *testing.T) {
			// --- When ---
			have := IsNil(tc.Val)

			// --- Then ---
			if tc.IsNil != have {
				format := "expected IsNil to return %v for %#v"
				t.Errorf(format, tc.IsNil, tc.Val)
			}
		})
	}
}

func Test_IsNil(t *testing.T) {
	t.Run("nil unsafe pointer", func(t *testing.T) {
		// --- When ---
		have := IsNil(unsafe.Pointer(nil))

		// --- Then ---
		if !have {
			t.Error("expected nil value")
		}
	})

	t.Run("not nil unsafe pointer", func(t *testing.T) {
		// --- Given ---
		var i int

		// --- When ---
		have := IsNil(unsafe.Pointer(&i))

		// --- Then ---
		if have {
			t.Error("expected not-nil value")
		}
	})
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
		fn := func() { panic(nil) } //nolint:govet

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

func Test_Same(t *testing.T) {
	t.Run("method values of different receivers", func(t *testing.T) {
		// --- Given ---
		a, b := &testcases.TPtr{Val: "a"}, &testcases.TPtr{Val: "b"}

		// --- When ---
		have := Same(a.AAA, b.AAA)

		// --- Then ---
		if have {
			t.Error("expected not same")
		}
	})

	t.Run("closures with different captures", func(t *testing.T) {
		// --- Given ---
		mk := func(i int) func() int { return func() int { return i } }

		// --- When ---
		have := Same(mk(1), mk(2))

		// --- Then ---
		if have {
			t.Error("expected not same")
		}
	})

	t.Run("the same closure", func(t *testing.T) {
		// --- Given ---
		mk := func(i int) func() int { return func() int { return i } }
		fn := mk(1)

		// --- When ---
		have := Same(fn, fn)

		// --- Then ---
		if !have {
			t.Error("expected same")
		}
	})

	t.Run("sub-slice of the same array", func(t *testing.T) {
		// --- Given ---
		s := []int{1, 2, 3}

		// --- When ---
		have := Same(s[:1], s)

		// --- Then ---
		if have {
			t.Error("expected not same")
		}
	})

	t.Run("distinct empty slices", func(t *testing.T) {
		// --- When ---
		have := Same([]int{}, make([]int, 0))

		// --- Then ---
		if have {
			t.Error("expected not same")
		}
	})
}

func Test_Same_tabular(t *testing.T) {
	// Pointers.
	ptr0 := &testcases.TPtr{Val: "0"}
	ptr1 := &testcases.TPtr{Val: "1"}

	// Interfaces.
	var itfPtr0, itfPtr1 testcases.TItf
	itfPtr0, itfPtr1 = &testcases.TPtr{Val: "0"}, &testcases.TPtr{Val: "1"}

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
	type TS []int
	var sA, sB TS

	// Arrays.
	a0 := [2]int{1, 2}
	a1 := [2]int{1, 2}
	var aNil0 [2]int
	var aNil1 [2]float64
	type Ta [2]int
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
		{
			"ptr different types not same",
			&testcases.TPtr{},
			&testcases.TVal{},
			false,
		},
		{"ptr nil both", nil, nil, false},

		{"itf ptr same", itfPtr0, itfPtr0, true},
		{"itf ptr not same", itfPtr0, itfPtr1, false},
		{"obj val not same", testcases.TVal{}, testcases.TVal{}, false},

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
		{"array nil same type both", aA, aB, false},

		{"map same", m0, m0, true},
		{"map not same", m0, m1, false},
		{"map nil want", nil, m1, false},
		{"map nil var want", mNil0, m1, false},
		{"map nil have", m0, nil, false},
		{"map nil var have", m0, mNil0, false},
		{"map nil both", nil, nil, false},
		{"map nil var both", mNil0, mNil0, true},
		{"map nil different type both", mNil0, mNil1, false},
		{"map nil same type both", mA, mB, true},

		{"chan same", c0, c0, true},
		{"chan not the same", c0, c1, false},
		{"chan nil want", nil, c1, false},
		{"chan nil var want", cNil0, c1, false},
		{"chan nil have", c0, nil, false},
		{"chan nil var have", c0, cNil0, false},
		{"chan nil both", nil, nil, false},
		{"chan nil var both", cNil0, cNil0, true},
		{"chan nil different type both", cNil0, cNilA, false},
		{"chan nil same type both", cA, cB, true},
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
	ptr := &struct{ Int int }{Int: 123}
	s2 := []*testcases.TA{{TAp: &testcases.TA{Int: 42}}}

	tt := []struct {
		testN string

		in   reflect.Value
		want unsafe.Pointer
	}{
		{"nil", reflect.ValueOf(nil), nil},
		{"struct value", reflect.ValueOf(struct{}{}), nil},
		{"nil struct pointer", reflect.ValueOf((*testcases.TIntStr)(nil)), nil},
		{"struct pointer", reflect.ValueOf(ptr), unsafe.Pointer(ptr)},
		{
			"addressable struct pointer",
			reflect.ValueOf(s2).Index(0),
			unsafe.Pointer(s2[0]),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := Pointer(tc.in)

			// --- Then ---
			if tc.want != have {
				t.Errorf(expMsg, tc.want, have)
			}
		})
	}
}

//nolint:cyclop
func Test_Value(t *testing.T) {
	t.Run("zero reflect value", func(t *testing.T) {
		// --- When ---
		hVal, hOK := Value(reflect.Value{})

		// --- Then ---
		if hVal != nil {
			t.Error("expected nil value")
		}
		if !hOK {
			t.Error("expected true (zero reflect.Value is treated as nil)")
		}
	})

	t.Run("struct by pointer - with exported fields", func(t *testing.T) {
		// --- Given ---
		ps := &testcases.TIntStr{Int: 42, Str: "abc"}
		val := reflect.ValueOf(ps)

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		have := hVal.(*testcases.TIntStr)
		if !reflect.DeepEqual(have, &testcases.TIntStr{Int: 42, Str: "abc"}) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by pointer - exported field", func(t *testing.T) {
		// --- Given ---
		ps := &testcases.TIntStr{Int: 42, Str: "abc"}
		val := reflect.ValueOf(ps).Elem().FieldByName("Int")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		have := hVal.(int)
		if have != 42 {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("function", func(t *testing.T) {
		// --- Given ---
		add := func(a, b int) int { return a + b }
		val := reflect.ValueOf(add)

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		if hVal == nil {
			t.Error("expected non-nil value")
		}
		fn, ok := hVal.(func(int, int) int)
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
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		if hVal.(string) != "str" {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private int field", func(t *testing.T) {
		// --- Given ---
		prv := testcases.NewTPrv().SetInt(42)
		val := reflect.ValueOf(prv).FieldByName("vInt")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		if hVal.(int) != 42 {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private named simple field", func(t *testing.T) {
		// --- Given ---
		prv := struct{ dur time.Duration }{dur: 42}
		val := reflect.ValueOf(prv).Field(0)

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		if have, ok := hVal.(time.Duration); !ok || have != 42 {
			t.Errorf(expMsg, time.Duration(42), hVal)
		}
	})

	t.Run("struct by value - private pointer field", func(t *testing.T) {
		// --- Given ---
		prv := testcases.NewTPrv().SetPtr(&testcases.TVal{Val: "abc"})
		val := reflect.ValueOf(prv).FieldByName("ptr")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		have := hVal.(*testcases.TVal)
		if have.Val != "abc" {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private slice field", func(t *testing.T) {
		// --- Given ---
		prv := testcases.NewTPrv().SetSInt([]int{1, 2, 3})
		val := reflect.ValueOf(prv).FieldByName("sInt")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		have := hVal.([]int)
		if !reflect.DeepEqual(have, []int{1, 2, 3}) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private array field", func(t *testing.T) {
		// --- Given ---
		prv := testcases.NewTPrv().SetAInt([...]int{1, 2})
		val := reflect.ValueOf(prv).FieldByName("aInt")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		have := hVal.([2]int)
		if !reflect.DeepEqual(have, [...]int{1, 2}) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private map field", func(t *testing.T) {
		// --- Given ---
		prv := testcases.NewTPrv().SetMapII(map[int]int{1: 11, 2: 22})
		val := reflect.ValueOf(prv).FieldByName("vMap")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		have := hVal.(map[int]int)
		if !reflect.DeepEqual(have, map[int]int{1: 11, 2: 22}) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private time field", func(t *testing.T) {
		// --- Given ---
		tim := time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
		prv := testcases.NewTPrv().SetTim(tim)
		val := reflect.ValueOf(prv).FieldByName("tim")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		if !hVal.(time.Time).Equal(tim) {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private function field", func(t *testing.T) {
		// --- Given ---
		fn := func() int { return 42 }
		prv := testcases.NewTPrv().SetFn(fn)
		val := reflect.ValueOf(prv).FieldByName("fn")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		if Same(hVal, fn) != true {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("struct by value - private chan field", func(t *testing.T) {
		// --- Given ---
		ch := make(chan int)
		prv := testcases.NewTPrv().SetCh(ch)
		val := reflect.ValueOf(prv).FieldByName("ch")

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		if Same(hVal, ch) != true {
			t.Errorf("expected correct field value")
		}
	})

	t.Run("interface", func(t *testing.T) {
		// --- Given ---
		itf := &testcases.TPtr{}
		val := reflect.ValueOf(itf)

		// --- When ---
		hVal, hOK := Value(val)

		// --- Then ---
		if !hOK {
			t.Error("expected success")
		}
		if Same(hVal, itf) != true {
			t.Errorf("expected correct field value")
		}
	})
}

func Test_Value_success_tabular(t *testing.T) {
	chn := make(chan int)
	m := make(map[string]int)
	ptr := &testcases.TPtr{}

	tt := []struct {
		testN string

		in   any
		wVal any
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
		{"struct", testcases.TPtr{}, testcases.TPtr{}},
		{"named int", testcases.TIntType(42), testcases.TIntType(42)},
		{"named string", testcases.TStrType("abc"), testcases.TStrType("abc")},
		{"duration", time.Duration(42), time.Duration(42)},

		{"uintptr", uintptr(42), uintptr(42)},
		{"unsafe pointer", unsafe.Pointer(ptr), unsafe.Pointer(ptr)},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			hVal, hOK := Value(reflect.ValueOf(tc.in))

			// --- Then ---
			if !reflect.DeepEqual(tc.wVal, hVal) {
				t.Errorf(expMsg, tc.wVal, hVal)
			}
			if !hOK {
				t.Error("expected true")
			}
		})
	}
}

func Test_Value_unexported_tabular(t *testing.T) {
	i := 42
	chn := make(chan int)
	m := map[string]int{"a": 1}

	tt := []struct {
		testN string

		in   reflect.Value
		want any
	}{
		{
			"nil pointer",
			reflect.ValueOf(struct{ loc *time.Location }{}).Field(0),
			(*time.Location)(nil),
		},
		{
			"pointer",
			reflect.ValueOf(struct{ ptr *int }{ptr: &i}).Field(0),
			&i,
		},
		{
			"map",
			reflect.ValueOf(struct{ m map[string]int }{m: m}).Field(0),
			m,
		},
		{
			"chan",
			reflect.ValueOf(struct{ c chan int }{c: chn}).Field(0),
			chn,
		},
		{
			"interface",
			reflect.ValueOf(struct {
				err error
				i   int
			}{err: io.EOF}).Field(0),
			io.EOF,
		},
		{
			"nil interface",
			reflect.ValueOf(struct {
				err error
				i   int
			}{}).Field(0),
			nil,
		},
		{
			"uintptr",
			reflect.ValueOf(struct {
				u uintptr
				i int
			}{u: 42}).Field(0),
			uintptr(42),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have, hOK := Value(tc.in)

			// --- Then ---
			if !hOK {
				t.Error("expected true")
			}
			if !reflect.DeepEqual(tc.want, have) {
				t.Errorf(expMsg, tc.want, have)
			}
		})
	}
}

func Test_ValueSimple_tabular(t *testing.T) {
	chn := make(chan int)
	m := make(map[string]int)
	ptr := &testcases.TPtr{}

	tt := []struct {
		testN string

		in   any
		wVal any
		wOK  bool
	}{
		{"bool - false", false, false, true},
		{"bool - true", true, true, true},
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
		{"struct", testcases.TPtr{}, nil, false},
		{"unsafe pointer", unsafe.Pointer(ptr), nil, false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			hVal, hOK := ValueSimple(reflect.ValueOf(tc.in))

			// --- Then ---
			if tc.wVal != hVal {
				t.Errorf(expMsg, tc.wVal, hVal)
			}
			if tc.wOK != hOK {
				t.Errorf(expMsg, tc.wOK, hOK)
			}
		})
	}
}

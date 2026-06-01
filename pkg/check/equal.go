// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"unsafe"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/notice"
)

// Equal recursively checks that want and have are equal.
//
// See the package documentation for the customization model
// ([DefaultOptions], [WithTrail], [WithTypeChecker], etc.) and
// [dump] options.
func Equal(want, have any, opts ...any) error {
	ops := DefaultOptions(opts...)
	if _, ok := ops.Dumper.Dumpers[typByte]; !ok {
		ops.Dumper.Dumpers[typByte] = dumpByte
	}
	wVal := reflect.ValueOf(want)
	hVal := reflect.ValueOf(have)
	return deepEqual(wVal, hVal, make(map[visit]bool), WithOptions(ops))
}

// NotEqual checks that want and have are not equal.
//
// See the package documentation for the customization model.
func NotEqual(want, have any, opts ...any) error {
	if err := Equal(want, have, opts...); err == nil {
		return equalError(want, have, opts...).
			SetHeader("expected values not to be equal")
	}
	return nil
}

// visit tracks pointers during [deepEqual] to detect cycles and prevent
// infinite recursion.
type visit struct {
	want unsafe.Pointer
	have unsafe.Pointer
	typ  reflect.Type
}

// deepEqual is the internal recursive comparison engine used by [Equal].
// It handles all Go kinds, custom checkers (see [WithTypeChecker] /
// [WithTrailChecker]), trail skipping ([WithSkipTrail]), unexported field
// skipping ([WithSkipUnexported]), and cycle detection.
//
// nolint: gocognit, cyclop
func deepEqual(
	wVal, hVal reflect.Value,
	visited map[visit]bool,
	opts ...any,
) error {

	ops := DefaultOptions(opts...)

	// Return when the trail should be skipped (see [WithSkipTrail]).
	if i := slices.Index(ops.SkipTrails, ops.Trail); i >= 0 {
		ops.Trail += " <skipped>"
		ops.LogTrail()
		return nil
	}

	// Skip unexported fields if the option is turned on
	// (see [WithSkipUnexported]).
	if wVal.IsValid() && !wVal.CanInterface() && ops.SkipUnexported {
		trail := ops.Trail
		ops.Trail += " <skipped>"
		ops.LogTrail()
		ops.Trail = trail
		return nil
	}

	// Both are untyped nil value.
	if !wVal.IsValid() && !hVal.IsValid() {
		ops.LogTrail()
		return nil
	}

	// One of the values is untyped nil.
	if !wVal.IsValid() || !hVal.IsValid() {
		ops.LogTrail()
		if wVal.IsValid() {
			wStr := ops.Dumper.Value(wVal)
			msg := notice.New("expected values to be equal").
				Want("%s", wStr).
				Have("%s", dump.ValNil)
			return AddRows(ops, msg)
		}

		hStr := ops.Dumper.Value(hVal)
		msg := notice.New("expected values to be equal").
			Want("%s", dump.ValNil).
			Have("%s", hStr)
		return AddRows(ops, msg)
	}

	// Check both types are the same.
	wTyp := wVal.Type()
	hTyp := hVal.Type()
	if wTyp != hTyp {
		// Compare simple types if any of the types is an alias.
		if ops.CmpSimpleType {
			wSmp, wOK := core.ValueSimple(wVal)
			hSmp, hOK := core.ValueSimple(hVal)
			if wOK && hOK {
				wVal = reflect.ValueOf(wSmp)
				hVal = reflect.ValueOf(hSmp)
				return deepEqual(wVal, hVal, visited, WithOptions(ops))
			}
		}

		ops.LogTrail()
		msg := notice.New("expected values to be equal").
			Append("want type", "%s", wTyp).
			Append("have type", "%s", hTyp)
		return AddRows(ops, msg)
	}

	// Detect already compared pointers (cycle protection via the visited map).
	wPtr := core.Pointer(wVal)
	hPtr := core.Pointer(hVal)
	if wPtr != nil && hPtr != nil {
		v := visit{wPtr, hPtr, wTyp}
		if visited[v] {
			return nil
		}
		visited[v] = true
	}

	var chk Checker
	if chk = ops.TrailCheckers[ops.Trail]; chk == nil {
		chk = ops.TypeCheckers[wTyp]
	}

	if chk != nil {
		ops.LogTrail()
		wItf, wOk := core.Value(wVal)
		hItf, hOk := core.Value(hVal)
		if !wOk || !hOk {
			// core.Value could not extract the underlying value for a
			// registered custom checker (can happen with unexported fields).
			// Return a clear error instead of calling the checker with
			// incomplete data.
			msg := notice.New("not able to compare using a custom checker").
				Append("want type", "%s", wTyp).
				Append("have type", "%s", hTyp)
			return AddRows(ops, msg)
		}
		return chk(wItf, hItf, WithOptions(ops))
	}

	switch knd := wVal.Kind(); knd {
	case reflect.Pointer:
		// Recurse into the pointed-to values (nil == nil is already handled above).
		if wVal.IsNil() && hVal.IsNil() {
			ops.LogTrail()
			return nil
		}
		return deepEqual(wVal.Elem(), hVal.Elem(), visited, WithOptions(ops))

	case reflect.Struct:
		// Compare exported fields one by one, building per-field trails.
		var err error
		typeName := wVal.Type().Name()
		sOps := ops.StructTrail(typeName, "")
		for i := 0; i < wVal.NumField(); i++ {
			wfVal := wVal.Field(i)
			hfVal := hVal.Field(i)
			if !wfVal.IsValid() {
				continue
			}
			wSF := wVal.Type().Field(i)
			iOps := sOps.StructTrail("", wSF.Name)
			if e := deepEqual(wfVal, hfVal, visited, WithOptions(iOps)); e != nil {
				err = notice.Join(err, e)
			}
		}
		return err

	case reflect.Slice, reflect.Array:
		// Lengths must match first. Then compare element-wise with index trails.
		if wVal.Len() != hVal.Len() {
			ops.LogTrail()
			wStr, hStr, diff := ops.Dumper.DiffValue(wVal, hVal)
			msg := notice.New("expected values to be equal").
				Prepend("have len", "%d", hVal.Len()).
				Prepend("want len", "%d", wVal.Len()).
				Want("%s", wStr).
				Have("%s", hStr).
				Append("diff", "%s", diff)
			return AddRows(ops, msg)
		}
		if knd == reflect.Slice && wVal.Pointer() == hVal.Pointer() {
			ops.LogTrail()
			return nil
		}
		var err error
		for i := 0; i < wVal.Len(); i++ {
			wiVal := wVal.Index(i)
			hiVal := hVal.Index(i)
			iOps := ops.ArrTrail(knd.String(), i)
			if e := deepEqual(wiVal, hiVal, visited, WithOptions(iOps)); e != nil {
				err = notice.Join(err, e)
			}
		}
		return err

	case reflect.Map:
		// Lengths must match. Compare only keys present in "want" (extra keys in
		// "have" are allowed). Keys are sorted for deterministic ordering.
		if wVal.Len() != hVal.Len() {
			ops.LogTrail()
			wStr, hStr, diff := ops.Dumper.DiffValue(wVal, hVal)
			msg := notice.New("expected values to be equal").
				Prepend("have len", "%d", hVal.Len()).
				Prepend("want len", "%d", wVal.Len()).
				Want("%s", wStr).
				Have("%s", hStr).
				Append("diff", "%s", diff)
			return AddRows(ops, msg)
		}
		if wVal.Pointer() == hVal.Pointer() {
			ops.LogTrail()
			return nil
		}

		keys := wVal.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return valToString(keys[i]) < valToString(keys[j])
		})

		var err error
		for _, key := range keys {
			wkVal := wVal.MapIndex(key)
			hkVal := hVal.MapIndex(key)
			kOps := ops.MapTrail(valToString(key))
			if !hkVal.IsValid() {
				hItf := hVal.Interface()
				e := equalError(hItf, nil, WithOptions(kOps))
				err = notice.Join(err, e)
				continue
			}
			if e := deepEqual(wkVal, hkVal, visited, WithOptions(kOps)); e != nil {
				err = notice.Join(err, e)
			}
		}
		return err

	case reflect.Interface:
		// Compare the concrete values inside the interfaces.
		wElem := wVal.Elem()
		hElem := hVal.Elem()
		return deepEqual(wElem, hElem, visited, WithOptions(ops))

	case reflect.Bool:
		// For all primitive kinds we simply compare the underlying values
		// and use equalError for the failure message.
		ops.LogTrail()
		w, h := wVal.Bool(), hVal.Bool()
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Int:
		ops.LogTrail()
		w, h := int(wVal.Int()), int(hVal.Int())
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Int8:
		ops.LogTrail()
		w, h := int8(wVal.Int()), int8(hVal.Int()) // nolint: gosec
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Int16:
		ops.LogTrail()
		w, h := int16(wVal.Int()), int16(hVal.Int()) // nolint: gosec
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Int32:
		ops.LogTrail()
		w, h := int32(wVal.Int()), int32(hVal.Int()) // nolint: gosec
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Int64:
		ops.LogTrail()
		w, h := wVal.Int(), hVal.Int()
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Uint:
		ops.LogTrail()
		w, h := uint(wVal.Uint()), uint(hVal.Uint())
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Uint8:
		ops.LogTrail()
		w, h := uint8(wVal.Uint()), uint8(hVal.Uint()) // nolint: gosec
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Uint16:
		ops.LogTrail()
		w, h := uint16(wVal.Uint()), uint16(hVal.Uint()) // nolint: gosec
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Uint32:
		ops.LogTrail()
		w, h := uint32(wVal.Uint()), uint32(hVal.Uint()) // nolint: gosec
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Uint64:
		ops.LogTrail()
		w, h := wVal.Uint(), hVal.Uint()
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Float32:
		ops.LogTrail()
		w, h := float32(wVal.Float()), float32(hVal.Float()) // nolint: gosec
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Float64:
		ops.LogTrail()
		w, h := wVal.Float(), hVal.Float()
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Complex64:
		ops.LogTrail()
		w, h := complex64(wVal.Complex()), complex64(hVal.Complex())
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Complex128:
		ops.LogTrail()
		w, h := wVal.Complex(), hVal.Complex()
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.String:
		ops.LogTrail()
		w, h := wVal.String(), hVal.String()
		if w == h {
			return nil
		}
		return equalError(w, h, WithOptions(ops))

	case reflect.Chan:
		ops.LogTrail()
		w, h := wVal.Pointer(), hVal.Pointer()
		if w == h {
			return nil
		}
		msg := notice.New("expected values to be equal").
			Want("%s", dump.ChanDumper(ops.Dumper, 0, wVal)).
			Have("%s", dump.ChanDumper(ops.Dumper, 0, hVal))
		return AddRows(ops, msg)

	case reflect.Func:
		ops.LogTrail()
		w, h := wVal.Pointer(), hVal.Pointer()
		if w == h {
			return nil
		}
		msg := notice.New("expected values to be equal").
			Want("%s", dump.FuncDumper(ops.Dumper, 0, wVal)).
			Have("%s", dump.FuncDumper(ops.Dumper, 0, hVal))
		return AddRows(ops, msg)

	case reflect.Uintptr:
		ops.LogTrail()
		w, h := wVal.Uint(), hVal.Uint()
		if w == h {
			return nil
		}
		msg := notice.New("expected values to be equal").
			Want("%s", dump.HexPtrDumper(ops.Dumper, 0, wVal)).
			Have("%s", dump.HexPtrDumper(ops.Dumper, 0, hVal))
		return AddRows(ops, msg)

	case reflect.UnsafePointer:
		ops.LogTrail()
		w, h := wVal.Pointer(), hVal.Pointer()
		if w == h {
			return nil
		}
		msg := notice.New("expected values to be equal").
			Want("%s", dump.HexPtrDumper(ops.Dumper, 0, wVal)).
			Have("%s", dump.HexPtrDumper(ops.Dumper, 0, hVal))
		return AddRows(ops, msg)

	default:
		ops.LogTrail()
		msg := notice.New("cannot compare values").
			Append("cause", "%s", "value cannot be used without panicking").
			Append("hint", "%s", "use WithSkipTrail or WithSkipUnexported "+
				"option to skip this field")
		return AddRows(ops, msg)
	}
}

// equalError builds the standard "expected values to be equal" notice.
func equalError(want, have any, opts ...any) *notice.Notice {
	wTyp, hTyp := fmt.Sprintf("%T", want), fmt.Sprintf("%T", have)
	if wTyp == hTyp {
		wTyp, hTyp = "", ""
	}

	ops := DefaultOptions(opts...)
	if _, ok := ops.Dumper.Dumpers[typByte]; !ok {
		ops.Dumper.Dumpers[typByte] = dumpByte
	}

	msg := notice.New("expected values to be equal")
	if wTyp != "" {
		_ = msg.
			Append("want type", "%s", wTyp).
			Append("have type", "%s", hTyp)
	}

	wStr, hStr, diff := ops.Dumper.Diff(want, have)
	_ = msg.Want("%s", wStr).Have("%s", hStr)

	var assignable bool
	if want != nil && have != nil {
		assignable = reflect.TypeOf(want).AssignableTo(reflect.TypeOf(have))
	}
	if diff != "" && assignable {
		_ = msg.Append("diff", "%s", diff)
	}
	return AddRows(ops, msg)
}

// dumpByte is a custom bumper for bytes.
func dumpByte(dmp dump.Dump, lvl int, val reflect.Value) string {
	if !val.CanInterface() {
		return dump.ValCannotPrint
	}
	v := val.Interface().(byte) // nolint: forcetypeassert
	var str string
	if isPrintableChar(v) {
		str = fmt.Sprintf("0x%02x ('%s')", v, string(v))
	} else {
		str = fmt.Sprintf("0x%02x", v)
	}
	prn := dump.NewPrinter(dmp)
	return prn.Tab(dmp.Indent + lvl).Write(str).String()
}

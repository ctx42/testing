// SPDX-FileCopyrightText: (c) 2025 Rafal Zajac <rzajac@gmail.com>
// SPDX-License-Identifier: MIT

package check

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"time"

	"github.com/ctx42/testing/pkg/dump"
	"github.com/ctx42/testing/pkg/notice"
)

// Equal recursively checks both values are equal. Returns nil if they are,
// otherwise it returns an error with a message indicating the expected and
// actual values.
func Equal(want, have any, opts ...Option) error {
	wVal := reflect.ValueOf(want)
	hVal := reflect.ValueOf(have)
	return deepEqual(wVal, hVal, opts...)
}

// NotEqual checks both values are not equal using. Returns nil if they are not,
// otherwise it returns an error with a message indicating the expected and
// actual values.
func NotEqual(want, have any, opts ...Option) error {
	if err := Equal(want, have, opts...); err == nil {
		return equalError(want, have, opts...).
			SetHeader("expected values not to be equal")
	}
	return nil
}

// deepEqual is the internal comparison function which is called recursively.
//
// nolint: gocognit, cyclop
func deepEqual(wVal, hVal reflect.Value, opts ...Option) error {
	ops := DefaultOptions(opts...)

	if i := slices.Index(ops.SkipTrails, ops.Trail); i >= 0 {
		ops.Trail += " <skipped>"
		ops.LogTrail()
		return nil
	}

	if !wVal.IsValid() && !hVal.IsValid() {
		ops.LogTrail()
		return nil
	}

	if !wVal.IsValid() || !hVal.IsValid() {
		var wItf, hItf any
		if wVal.IsValid() {
			wItf = wVal.Interface()
		}
		if hVal.IsValid() {
			hItf = hVal.Interface()
		}
		ops.LogTrail()
		return equalError(wItf, hItf, WithOptions(ops))
	}

	if !wVal.CanInterface() {
		trail := ops.Trail
		ops.Trail += " <skipped>"
		ops.LogTrail()
		ops.Trail = trail
		if ops.SkipUnexported {
			return nil
		}
		return notice.New("cannot compare values").
			SetTrail(ops.Trail).
			Append("cause", "%s", "value cannot be used without panicking").
			Append("hint", "%s", "use WithSkipTrail or WithSkipUnexported "+
				"option to skip this field")
	}

	wType := wVal.Type()
	hType := hVal.Type()
	if wType != hType {
		ops.LogTrail()
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))
	}

	if chk, ok := ops.TrailCheckers[ops.Trail]; ok {
		ops.LogTrail()
		return chk(wVal.Interface(), hVal.Interface(), WithOptions(ops))
	}

	if chk, ok := ops.TypeCheckers[wType]; ok {
		// TODO(rz): Log we are using custom checker.
		ops.LogTrail()
		return chk(wVal.Interface(), hVal.Interface(), opts...)
	}

	switch knd := wVal.Kind(); knd {
	case reflect.Ptr:
		if wType == typTimeLocPtr && hType == typTimeLocPtr {
			ops.LogTrail()
			wZone := wVal.Interface().(*time.Location) // nolint: forcetypeassert
			hZone := hVal.Interface().(*time.Location) // nolint: forcetypeassert
			return Zone(wZone, hZone, WithOptions(ops))
		}

		if wVal.IsNil() && hVal.IsNil() {
			ops.LogTrail()
			return nil
		}
		if wVal.IsNil() || hVal.IsNil() {
			ops.LogTrail()
			wItf := wVal.Interface()
			hItf := hVal.Interface()
			return equalError(wItf, hItf, WithOptions(ops))
		}

		return deepEqual(wVal.Elem(), hVal.Elem(), WithOptions(ops))

	case reflect.Struct:
		wTyp := wVal.Type()
		hTyp := hVal.Type()
		if wTyp == typTime && hTyp == typTime {
			ops.LogTrail()
			return Time(wVal.Interface(), hVal.Interface(), opts...)
		}
		if wTyp == typTimeLoc && hTyp == typTimeLoc {
			ops.LogTrail()
			wZone := wVal.Interface().(time.Location) // nolint: forcetypeassert
			hZone := hVal.Interface().(time.Location) // nolint: forcetypeassert
			return Zone(&wZone, &hZone, opts...)
		}
		typeName := wVal.Type().Name()

		var err error
		for i := 0; i < wVal.NumField(); i++ {
			wfVal := wVal.Field(i)
			hfVal := hVal.Field(i)
			if !wfVal.IsValid() {
				continue
			}
			wSF := wVal.Type().Field(i)
			iOps := ops.StructTrail(typeName, wSF.Name)
			if e := deepEqual(wfVal, hfVal, WithOptions(iOps)); e != nil {
				err = notice.Join(err, e)
			}
		}
		return err

	case reflect.Slice, reflect.Array:
		if wVal.Len() != hVal.Len() {
			ops.LogTrail()
			wItf := wVal.Interface()
			hItf := hVal.Interface()
			return equalError(wItf, hItf, WithOptions(ops)).
				Prepend("have len", "%d", hVal.Len()).
				Prepend("want len", "%d", wVal.Len())
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
			if e := deepEqual(wiVal, hiVal, WithOptions(iOps)); e != nil {
				err = notice.Join(err, e)
			}
		}
		return err

	case reflect.Map:
		if wVal.Len() != hVal.Len() {
			ops.LogTrail()
			wItf := wVal.Interface()
			hItf := hVal.Interface()
			return equalError(wItf, hItf, WithOptions(ops)).
				Prepend("have len", "%d", hVal.Len()).
				Prepend("want len", "%d", wVal.Len())
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
			if e := deepEqual(wkVal, hkVal, WithOptions(kOps)); e != nil {
				err = notice.Join(err, e)
			}
		}
		return err

	case reflect.Interface:
		wElem := wVal.Elem()
		hElem := hVal.Elem()
		return deepEqual(wElem, hElem, WithOptions(ops))

	case reflect.Bool:
		ops.LogTrail()
		if wVal.Bool() == hVal.Bool() {
			return nil
		}
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ops.LogTrail()
		if wVal.Int() == hVal.Int() {
			return nil
		}
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		ops.LogTrail()
		if wVal.Uint() == hVal.Uint() {
			return nil
		}
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))

	case reflect.Float32, reflect.Float64:
		ops.LogTrail()
		if wVal.Float() == hVal.Float() {
			return nil
		}
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))

	case reflect.Complex64, reflect.Complex128:
		ops.LogTrail()
		if wVal.Complex() == hVal.Complex() {
			return nil
		}
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))

	case reflect.String:
		ops.LogTrail()
		if wVal.String() == hVal.String() {
			return nil
		}
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))

	case reflect.Chan, reflect.Func:
		ops.LogTrail()
		if wVal.Pointer() == hVal.Pointer() {
			return nil
		}
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))

	default:
		ops.LogTrail()
		// For types, we haven't explicitly handled, use DeepEqual.
		if reflect.DeepEqual(wVal.Interface(), hVal.Interface()) {
			return nil
		}
		return equalError(wVal.Interface(), hVal.Interface(), WithOptions(ops))
	}
}

// equalError returns error for not equal values.
func equalError(want, have any, opts ...Option) *notice.Notice {
	wTyp, hTyp := fmt.Sprintf("%T", want), fmt.Sprintf("%T", have)
	if wTyp == hTyp {
		wTyp, hTyp = "", ""
	}

	ops := DefaultOptions(opts...)
	if _, ok := ops.Dumper.Dumpers[typByte]; !ok {
		ops.Dumper.Dumpers[typByte] = dumpByte
	}

	msg := notice.New("expected values to be equal").SetTrail(ops.Trail)
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
	return msg
}

// dumpByte is a custom bumper for bytes.
func dumpByte(dmp dump.Dump, lvl int, val reflect.Value) string {
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

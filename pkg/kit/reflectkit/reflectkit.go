package reflectkit

import (
	"reflect"

	"github.com/ctx42/testing/pkg/tester"
)

// GetField returns the [reflect.StructField] for a named field in struct s.
// The struct s must be a struct pointer. On error, the test is marked as
// failed and an error message is logged to the test log.
func GetField(t tester.T, s any, name string) reflect.StructField {
	t.Helper()

	if name == "" {
		t.Error("the struct field name must not be empty")
		return reflect.StructField{}
	}
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Ptr {
		t.Error("pointer to struct is required")
		return reflect.StructField{}
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		t.Errorf("type %T is not struct", s)
		return reflect.StructField{}
	}
	fld, exist := typ.FieldByName(name)
	if !exist {
		t.Errorf("struct `%T` has no field %#q", s, name)
		return reflect.StructField{}
	}
	return fld
}

// GetValue returns the [reflect.Value] for a named field in struct s.
// The struct s must be a struct pointer. On error, the test is marked as
// failed and an error message is logged to the test log.
func GetValue(t tester.T, s any, name string) reflect.Value {
	t.Helper()

	if name == "" {
		t.Error("the struct field name must not be empty")
		return reflect.Value{}
	}

	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		t.Errorf("type `%T` is not struct", s)
		return reflect.Value{}
	}

	sv := reflect.ValueOf(s)
	fld := reflect.Indirect(sv).FieldByName(name)
	if !fld.IsValid() {
		t.Errorf("cannot get value for `%T.%s`", s, name)
		return reflect.Value{}
	}
	return fld
}

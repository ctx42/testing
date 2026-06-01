// Package reflectkit provides lightweight reflection utilities for testing.
//
// It is part of the [github.com/ctx42/testing/pkg/kit] curated collection
// and offers safe helpers for inspecting struct fields via [tester.T],
// ensuring proper error reporting and test lifecycle integration.
//
// See the [kit] package README for context on the full collection.
// All helpers are designed to work with [tester.T] and the assertion
// packages.
package reflectkit

import (
	"reflect"

	"github.com/ctx42/testing/pkg/tester"
)

// GetField returns the [reflect.StructField] for the named field in struct s.
// The argument s must be a pointer to struct. On any error it reports via
// t.Error (never panics) and returns the zero StructField.
func GetField(t tester.T, s any, name string) reflect.StructField {
	t.Helper()

	if name == "" {
		t.Error("the struct field name must not be empty")
		return reflect.StructField{}
	}
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Pointer {
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

// GetValue returns the [reflect.Value] for the named field in struct s.
// The argument s may be a struct or pointer to struct. On any error it
// reports via t.Error (never panics) and returns the zero Value.
func GetValue(t tester.T, s any, name string) reflect.Value {
	t.Helper()

	if name == "" {
		t.Error("the struct field name must not be empty")
		return reflect.Value{}
	}

	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Pointer {
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

// SPDX-FileCopyrightText: (c) 2026 Rafal Zajac
// SPDX-License-Identifier: MIT

package mock_test

import (
	"fmt"
	"testing"

	"github.com/ctx42/testing/pkg/mock"
	"github.com/ctx42/testing/pkg/testcases"
)

// ExampleMock_On demonstrates basic expectation setup.
func ExampleMock_On() {
	m := mock.NewMock(&testing.T{})

	m.On("GetName", 42).Return("Alice", nil)

	name, err := m.Call("GetName", 42).String(0), m.Call("GetName", 42).Error(1)
	fmt.Println(name, err)

	// Output:
	// Alice <nil>
}

// ExampleCall_Alter shows how to modify arguments before the real method
// is called using the fluent Call API. This is a powerful advanced pattern.
func ExampleCall_Alter() {
	m := mock.NewMock(&testing.T{})

	// We want to mutate the map argument before the method "sees" it.
	m.On("Process", mock.Any).
		Alter(func(args mock.Arguments) {
			m := args.Get(0).(*map[string]any)
			(*m)["injected"] = true
		}).
		Return(nil)

	data := map[string]any{"original": 1}
	_ = m.Call("Process", &data)

	fmt.Printf("%+v\n", data)

	// Output:
	// map[injected:true original:1]
}

// --- Advanced Examples ---

// ExampleMock_Proxy demonstrates using Proxy to forward calls to a real
// implementation. This is one of the most powerful advanced features.
func ExampleMock_Proxy() {
	m := mock.NewMock(&testing.T{})

	// We use a real object from testcases for demonstration.
	target := &testcases.TPtr{Val: "original"}

	// Proxy all calls to realTarget.Wrap through the mock.
	m.Proxy(target.Wrap)

	// Call the proxied method through the mock.
	result := m.Call("Wrap", "pre-", "-post").String(0)
	fmt.Println("returned:", result)
	fmt.Println("real value:", target.Val)

	// Output:
	// returned: pre-original-post
	// real value: original
}

// ExampleCall_Requires demonstrates using Requires to express ordering
// dependencies between expectations. This is very useful for complex
// test scenarios.
func ExampleCall_Requires() {
	m := mock.NewMock(&testing.T{})

	// Init must be called before DoWork can succeed.
	initCall := m.On("Init").Return(nil)
	m.On("DoWork").Return("done").Requires(initCall)

	// Correct order
	_ = m.Call("Init")
	result := m.Call("DoWork").String(0)
	fmt.Println(result)

	// Output:
	// done
}

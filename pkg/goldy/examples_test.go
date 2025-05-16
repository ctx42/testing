package goldy_test

import (
	"fmt"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/goldy"
)

func ExampleNew() {
	tspy := core.NewSpy()

	content := goldy.Open(tspy, "testdata/text_case1.gld")

	fmt.Println(content)
	// Output:
	// Content #1.
	// Content #2.
}

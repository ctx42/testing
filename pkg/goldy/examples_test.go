package goldy_test

import (
	"fmt"

	"github.com/ctx42/testing/internal/core"
	"github.com/ctx42/testing/pkg/goldy"
)

func ExampleText() {
	tspy := core.NewSpy()

	content := goldy.Text(tspy, "testdata/text_case1.txt")

	fmt.Println(content)
	// Output:
	// Content #1.
	// Content #2.
}

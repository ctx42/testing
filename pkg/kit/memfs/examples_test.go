package fskit_test

import (
	"fmt"
	"io"

	"github.com/ctx42/testing/pkg/kit/memfs"
)

func ExampleFile() {
	buf := &fskit.File{}

	_, _ = buf.Write([]byte{0, 1, 2, 3})
	_, _ = buf.Seek(-2, io.SeekEnd)
	_, _ = buf.Write([]byte{4, 5})
	_, _ = buf.Seek(0, io.SeekStart)

	data, _ := io.ReadAll(buf)
	fmt.Println(data)

	// Output: [0 1 4 5]
}

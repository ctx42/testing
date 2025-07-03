package tstkit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

func ExampleClockFixed() {
	clk := ClockFixed(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC))

	fmt.Println(clk())
	fmt.Println(clk())
	fmt.Println(clk())
	// Output:
	// 2022-01-01 00:00:00 +0000 UTC
	// 2022-01-01 00:00:00 +0000 UTC
	// 2022-01-01 00:00:00 +0000 UTC
}

func ExampleClockDeterministic() {
	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	clk := ClockDeterministic(start, time.Hour)

	fmt.Println(clk())
	fmt.Println(clk())
	fmt.Println(clk())
	// Output:
	// 2022-01-01 00:00:00 +0000 UTC
	// 2022-01-01 01:00:00 +0000 UTC
	// 2022-01-01 02:00:00 +0000 UTC
}

func ExampleTikTak() {
	start := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	clk := TikTak(start)

	fmt.Println(clk())
	fmt.Println(clk())
	fmt.Println(clk())
	// Output:
	// 2022-01-01 00:00:00 +0000 UTC
	// 2022-01-01 00:00:01 +0000 UTC
	// 2022-01-01 00:00:02 +0000 UTC
}

func ExampleErrReader() {
	rdr := strings.NewReader("some text")
	rcs := ErrReader(rdr, 3)

	data, err := io.ReadAll(rcs)

	fmt.Printf("error: %v\n", err)
	fmt.Printf(" data: %s\n", string(data))
	// Output:
	// error: read error
	//  data: som
}

func ExampleErrReader_custom_error() {
	mye := errors.New("my error")
	rdr := strings.NewReader("some text")
	rcs := ErrReader(rdr, 4, WithReadErr(mye))

	data, err := io.ReadAll(rcs)

	fmt.Printf("error: %v\n", err)
	fmt.Printf(" data: %s\n", string(data))
	// Output:
	// error: my error
	//  data: some
}

func ExampleErrWriter() {
	dst := &bytes.Buffer{}
	ce := errors.New("my error")
	ew := ErrWriter(dst, 3, WithWriteErr(ce))

	n, err := ew.Write([]byte{0, 1, 2, 3})

	fmt.Printf("    n: %d\n", n)
	fmt.Printf("error: %v\n", err)
	fmt.Printf("  dst: %v\n", dst.Bytes())
	// Output:
	//     n: 3
	// error: my error
	//   dst: [0 1 2]
}

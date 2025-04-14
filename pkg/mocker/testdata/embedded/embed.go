package embedded

import (
	"io"

	"github.com/ctx42/testing/internal/types"
)

type Itf0 interface {
	Method0()
	Method1() error
}

type Itf1 interface {
	Method0()
	Method2(a int)
}

type EmbedLocal interface {
	Itf0
	Itf1
}

type EmbedLocalAndStd interface {
	Itf0
	Itf1
	io.Closer
	types.TItf
}

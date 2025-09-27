package cases

import (
	"io"

	"github.com/ctx42/testing/pkg/mocker/testdata/pkga"
	"github.com/ctx42/testing/pkg/testcases"
)

type ItfA interface {
	Method0()
	Method1() error
}

type ItfB interface {
	Method0()
	Method2(a int)
}

type ItfAlias = pkga.CaseA0
type ItfAliasLocal = ItfB

// EmbedLocal embeds two local interfaces where one of the methods has the same
// signature. The resulting mock should have only one "Method0" method.
type EmbedLocal interface {
	ItfA
	ItfB
}

// Embedder embeds local, STD library and interface from some other package.
type Embedder interface {
	ItfA
	ItfB
	io.Closer
	testcases.TItf
}

// EmptyEmbed embeds empty interface.
type EmptyEmbed interface {
	Empty
	ItfA
}

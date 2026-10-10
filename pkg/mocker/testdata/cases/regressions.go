package cases

import (
	"crypto/sha256"
	"net/url"

	xa "github.com/ctx42/testing/pkg/mocker/testdata/pkga"

	"github.com/ctx42/testing/pkg/mocker/testdata/dup/pkga"
	"github.com/ctx42/testing/pkg/tester"
)

const ArrLen = 3

// Parameter named like the local the generated body declares.
type CaseOkParam interface{ SetOn(ok bool) error }

// Parameter named like an imported package.
type CaseImportParam interface {
	Parse(url string) (*url.URL, error)
}

// Parameter type from the tester package the mock imports itself.
type CaseTesterParam interface{ Run(t tester.T) }

// Array lengths given by constants.
type CaseArrayLen interface {
	Sum() [sha256.Size]byte
	Fixed() [ArrLen]int
}

// Package name shared by an aliased and a not aliased import.
type CaseAliasClash interface{ Dup(v pkga.D1) xa.A1 }

// Interface embedded through an aliased import, using its local types.
type CaseEmbedAliased interface{ xa.CaseA1 }

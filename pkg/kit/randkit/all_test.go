package randkit

import (
	"github.com/ctx42/testing/pkg/check"
	"github.com/ctx42/testing/pkg/notice"
)

func init() {
	check.RegisterTypeChecker(options{}, optionsCheck)
}

// optionsCheck is a custom checker, matching [check.Check] signature,
// comparing two instances of options.
func optionsCheck(want, have any, opts ...any) error {
	ops := check.DefaultOptions(opts...)
	stOpt := check.WithOptions(ops)
	if _, err := check.SameType(options{}, have, stOpt); err != nil {
		return err
	}
	w, h := want.(options), have.(options)

	fName := check.FieldName(ops, "gopkg")
	ers := []error{
		check.Equal(w.chars, h.chars, fName("chars")),
		check.Equal(w.n, h.n, fName("n")),
		check.Equal(w.prefix, h.prefix, fName("prefix")),
		check.Equal(w.suffix, h.suffix, fName("suffix")),
		check.Fields(4, w, fName("{field count}")),
	}
	return notice.Join(ers...)
}

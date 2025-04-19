package cases

import (
	"fmt"
	"io/fs"
	mt "time"

	"github.com/ctx42/testing/pkg/mocker/testdata/pkga"
	"github.com/ctx42/testing/pkg/mocker/testdata/pkgb"
	"github.com/ctx42/testing/pkg/mocker/testdata/pkgc"
	"github.com/ctx42/testing/pkg/mocker/testdata/pkgd"
	. "github.com/ctx42/testing/pkg/mocker/testdata/pkge"
)

// BuildIn represents build-in type.
var BuildIn int

// Concrete represents a non-interface type.
type Concrete struct{}

// Brutal is a struct type alias.
type Brutal = Concrete

// Empty represents an interface with no methods.
type Empty interface{}

// EmptyAny represents an interface with no methods.
type EmptyAny any

// Small represents interface with one simple method.
type Small interface {
	// Method is the only method in the interface.
	Method() error
}

// Medium represents an interface with two simple methods.
type Medium interface {
	// Method0 method documentation.
	Method0() error
	// Method1 method documentation.
	Method1(a *pkga.A1) error
}

// Cases represent an interface with methods covering all the cases mock
// generator must know how to handle.
type Cases interface {
	Method0()
	Method1(a int)
	Method2(a, b int)
	Method3(a int, b int)
	Method4(a, b int, c bool)
	Method5(_ int)
	Method6(a, _ int, b bool)
	Method7() error
	Method8() (err error)
	Method9() (err0, err1 error)
	Method10() (int, error)
	Method11(int, float64)
	Method12(a ...int)
	Method13(tim mt.Time) error
	Method14(func())
	Method15(func(a int))
	Method16(a func(a ...int))
	Method17() Concrete
	Method18() *Concrete
	Method19() pkga.A1
	Method20() *pkga.A1
	Method21(a fmt.Stringer) fs.File
	Method22(a Concrete)
	Method23(a *Concrete)
	Method24(a ...Concrete) int
	Method25(a ...*Concrete)
	Method26(a ...pkga.A1)
	Method27(a ...*pkga.A1)
	Method28(a *int)
	Method29(a pkga.A1)
	Method30(a *pkga.A1)
	Method31(a [2]int)
	Method32(a [2]*int)
	Method33(a [2]pkga.A1)
	Method34(a [2]*pkga.A1)
	Method35(a []int)
	Method36(a []*int)
	Method37(a []pkga.A1)
	Method38(a []*pkga.A1)
	Method39(a map[int]string)
	Method40(a map[int]*string)
	Method41(a map[*int]string)
	Method42(a map[pkga.A1]string)
	Method43(a map[*pkga.A1]string)
	Method44(a map[*pkga.A1]*pkgb.B1)
	Method45(a chan map[*pkga.A1]*pkgb.B1)
	Method46(a map[*pkga.A1]*pkgb.B1, b pkgc.C1) *pkgd.D1
	Method47(a func(func(mt.Time, *pkga.A1)))
	Method48(a map[*pkga.A1]func(b pkgb.B1) func(a pkga.A1) error)
	Method49(a <-chan *pkga.A1) chan<- int
	Method50(a int) (b, c int, d error)
	Method51(e E1) error
	Method52(a Concrete, b mt.Time, c int) (*Concrete, error)
	Method53(a int, b bool) (int, bool, string, error)
	Method54(a int, b Concrete, c pkga.A1, d E1)
	Method55() *Other
	Method56(a string, b float64, c ...int) error
}

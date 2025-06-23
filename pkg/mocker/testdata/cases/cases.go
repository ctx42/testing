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

type Case00 interface{ Method00() }
type Case01 interface{ Method01(a int) }
type Case02 interface{ Method02(a, b int) }
type Case03 interface{ Method03(a int, b int) }
type Case04 interface{ Method04(a, b int, c bool) }
type Case05 interface{ Method05(_ int) }
type Case06 interface{ Method06(a, _ int, b bool) }
type Case07 interface{ Method07() error }
type Case08 interface{ Method08() (err error) }
type Case09 interface{ Method09() (err0, err1 error) }
type Case10 interface{ Method10() (int, error) }
type Case11 interface{ Method11(int, float64) }
type Case12 interface{ Method12(a ...int) }
type Case13 interface{ Method13(tim mt.Time) error }
type Case14 interface{ Method14(func()) }
type Case15 interface{ Method15(func(a int)) }
type Case16 interface{ Method16(a func(a ...int)) }
type Case17 interface{ Method17() Concrete }
type Case18 interface{ Method18() *Concrete }
type Case19 interface{ Method19() pkga.A1 }
type Case20 interface{ Method20() *pkga.A1 }
type Case21 interface{ Method21(a fmt.Stringer) fs.File }
type Case22 interface{ Method22(a Concrete) }
type Case23 interface{ Method23(a *Concrete) }
type Case24 interface{ Method24(a ...Concrete) int }
type Case25 interface{ Method25(a ...*Concrete) }
type Case26 interface{ Method26(a ...pkga.A1) }
type Case27 interface{ Method27(a ...*pkga.A1) }
type Case28 interface{ Method28(a *int) }
type Case29 interface{ Method29(a pkga.A1) }
type Case30 interface{ Method30(a *pkga.A1) }
type Case31 interface{ Method31(a [2]int) }
type Case32 interface{ Method32(a [2]*int) }
type Case33 interface{ Method33(a [2]pkga.A1) }
type Case34 interface{ Method34(a [2]*pkga.A1) }
type Case35 interface{ Method35(a []int) }
type Case36 interface{ Method36(a []*int) }
type Case37 interface{ Method37(a []pkga.A1) }
type Case38 interface{ Method38(a []*pkga.A1) }
type Case39 interface{ Method39(a map[int]string) }
type Case40 interface{ Method40(a map[int]*string) }
type Case41 interface{ Method41(a map[*int]string) }
type Case42 interface{ Method42(a map[pkga.A1]string) }
type Case43 interface{ Method43(a map[*pkga.A1]string) }
type Case44 interface{ Method44(a map[*pkga.A1]*pkgb.B1) }

type Case45 interface {
	Method45(a chan map[*pkga.A1]*pkgb.B1)
}

type Case46 interface {
	Method46(a map[*pkga.A1]*pkgb.B1, b pkgc.C1) *pkgd.D1
}

type Case47 interface {
	Method47(a func(func(mt.Time, *pkga.A1)))
}

type Case48 interface {
	Method48(a map[Concrete]func(b pkgb.B1) func(a pkga.A1) error)
}

type Case49 interface {
	Method49(a <-chan *pkga.A1) chan<- int
}

type Case50 interface {
	Method50(a int) (b, c int, d error)
}

type Case51 interface{ Method51(e E1) error }

type Case52 interface {
	Method52(a Concrete, b mt.Time, c int) (*Concrete, error)
}

type Case53 interface {
	Method53(a int, b bool) (int, bool, string, error)
}

type Case54 interface {
	Method54(a int, b Concrete, c pkga.A1, d E1)
}

type Case55 interface{ Method55() *Other }

type Case56 interface {
	Method56(a string, b float64, c ...int) error
}

type Case57 interface{ Method57() ParamOne[int] }

type Case58 interface {
	Method58() ParamTwo[int, *Concrete]
}

type Case59 interface{ Method59(...int) }

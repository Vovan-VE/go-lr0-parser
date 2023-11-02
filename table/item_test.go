package table

import (
	"testing"

	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/internal/testutils"
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

func TestItem_IsEqual(t *testing.T) {
	const (
		tInt symbol.Id = iota + 1
		tPlus
		nSum
	)
	l := lexer.New(
		lexer.NewTerm(tInt, "int").Str("1"),
	)
	ntDef := grammar.NewNT(nSum, "Sum").Is(tInt, tPlus, tInt).Do(calcStub3)

	rule1orig := ntDef.GetRules(l)[0]
	rule1copy := ntDef.GetRules(l)[0]

	s1 := newItem(rule1orig)

	if s1 != newItem(rule1orig) {
		t.Fatal("same not equal")
	}
	if s1 == newItem(rule1copy) {
		t.Error("different equal")
	}

	s2 := s1.Shift()
	if s1 != newItem(rule1orig) {
		t.Error("s1 changed")
	}
	if s2 == s1 {
		t.Error("s2 is same")
	}
}

func TestItem_Navigate(t *testing.T) {
	const (
		tInt symbol.Id = iota + 1
		tPlus
		nSum
		nValue
	)
	l := lexer.New(
		lexer.NewTerm(tInt, "int").Str("1"),
		lexer.NewTerm(tPlus, "plus").Str("+"),
	)
	ntDef := grammar.NewNT(nSum, "Sum").Is(nValue, tPlus, tInt).Do(calcStub3)

	r := ntDef.GetRules(l)[0]
	i0 := newItem(r)
	if i0.Expected() != nValue {
		t.Error("i0 expect() wrong: ", i0.Expected())
	}
	if !i0.HasFurther() {
		t.Error("i0 has further wrong")
	}

	i1 := i0.Shift()
	if i1.Expected() != tPlus {
		t.Error("i1 expect() wrong: ", i1.Expected())
	}
	if !i1.HasFurther() {
		t.Error("i1 has further wrong")
	}

	i2 := i1.Shift()
	if i2.Expected() != tInt {
		t.Error("i2 expect() wrong: ", i2.Expected())
	}
	if !i2.HasFurther() {
		t.Error("i2 has further wrong")
	}

	i3 := i2.Shift()
	if i3.Expected() != symbol.InvalidId {
		t.Error("i3 expect() wrong: ", i3.Expected())
	}
	if i3.HasFurther() {
		t.Error("i3 has further wrong")
	}

	t.Run("shift in the end", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, nil, func(t *testing.T, err error) {
			if err.Error() != "bad usage: internal error" {
				t.Error("another error", err)
			}
		})
		i3.Shift()
	})
}

func calcStub3(any, any, any) (any, error) { return nil, nil }

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

	rule1orig := grammar.ToImplementation(grammar.NewRule(nSum, []symbol.Id{tInt, tPlus, tInt}, calcStub3), stubHidden)
	rule1copy := grammar.ToImplementation(grammar.NewRule(nSum, []symbol.Id{tInt, tPlus, tInt}, calcStub3), stubHidden)

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
	r := grammar.ToImplementation(grammar.NewRule(nSum, []symbol.Id{nValue, tPlus, tInt}, calcStub3), stubHidden)
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
			if err.Error() != "internal bad usage" {
				t.Error("another error", err)
			}
		})
		i3.Shift()
	})
}

func calcStub3(any, any, any) (any, error) { return nil, nil }

var stubHidden lexer.HiddenRegistry = &_stubHidden{}

type _stubHidden struct{}

func (_ *_stubHidden) IsHidden(symbol.Id) bool { return false }

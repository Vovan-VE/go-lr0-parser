package lr0

import (
	"testing"

	"github.com/vovan-ve/go-lr0-parser/internal/testutils"
)

func TestGrammar(t *testing.T) {
	var (
		terminals = []Terminal{
			NewTerm(tInt, "integer").FuncByte(isDigit, bytesToInt),
			NewTerm(tPlus, "plus").Hide().Str("+"),
			NewTerm(tMinus, "minus").Hide().Str("-"),
		}
	)

	t.Run("panic: multiple main rules", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine, func(t *testing.T, err error) {
			if err.Error() != "another rule Sum2 has Main flag too, previous was Sum1: invalid definition" {
				t.Fatal("wrong error message:", err)
			}
		})

		newGrammar(terminals, []NonTerminalDefinition{
			NewNT(nVal, "Value").Is(tInt),
			NewNT(nSum, "Sum1").Main().Is(nVal),
			NewNT(nSum, "Sum2").Main().Is(nSum, tPlus, nVal).Do(calc2IntSum),
			NewNT(nSum, "Sum3").Main().Is(nSum, tMinus, nVal).Do(calc2IntSub),
		})
	})

	t.Run("panic: rule subject is Terminal", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine, func(t *testing.T, err error) {
			if err.Error() != "Non-Terminal Plus is Terminal: invalid definition" {
				t.Fatal("wrong error message:", err)
			}
		})

		newGrammar(terminals, []NonTerminalDefinition{
			NewNT(nVal, "Value").Is(tInt),
			NewNT(tPlus, "Plus").Is(nVal),
			NewNT(nSum, "Sum").
				Is(nSum, tPlus, nVal).Do(calc2IntSum).
				Is(nSum, tMinus, nVal).Do(calc2IntSub),
		})
	})

	t.Run("panic: non-terminal without rule", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine, func(t *testing.T, err error) {
			const (
				expected = `undefined non-terminals without rules:
- #10 in NT Sum rules[0] (Sum : #10) definitions[0]
: invalid definition`
			)
			if err.Error() != expected {
				t.Fatal("wrong error message:", err)
			}
		})

		newGrammar(terminals, []NonTerminalDefinition{
			NewNT(nSum, "Sum").
				Is(nVal).
				Is(nSum, tPlus, nVal).Do(calc2IntSum).
				Is(nSum, tMinus, nVal).Do(calc2IntSub),
		})
	})

	t.Run("panic: no main rule", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine, func(t *testing.T, err error) {
			if err.Error() != "no main rule with EOF flag: invalid definition" {
				t.Fatal("wrong error message:", err)
			}
		})

		newGrammar(terminals, []NonTerminalDefinition{
			NewNT(nVal, "Value").Is(tInt),
			NewNT(nSum, "Sum").
				Is(nVal).
				Is(nSum, tPlus, nVal).Do(calc2IntSum).
				Is(nSum, tMinus, nVal).Do(calc2IntSub),
			NewNT(nGoal, "Goal").Is(nSum),
		})
	})

	t.Run("panic: unused terminals", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine, func(t *testing.T, err error) {
			const (
				expected = `following Terminals are not used in any Rule:
- div
- MUL
: invalid definition`
			)
			if err.Error() != expected {
				t.Fatal("wrong error message:", err)
			}
		})

		terminals2 := append(
			terminals,
			NewTerm(tDiv, "div").Str("/"),
			NewTerm(tMul, "MUL").Str("*"),
		)

		newGrammar(terminals2, []NonTerminalDefinition{
			NewNT(nVal, "Value").Is(tInt),
			NewNT(nSum, "Sum").
				Is(nVal).
				Is(nSum, tPlus, nVal).Do(calc2IntSum).
				Is(nSum, tMinus, nVal).Do(calc2IntSub),
			NewNT(nGoal, "Goal").Main().Is(nSum),
		})
	})

	g := newGrammar(terminals, []NonTerminalDefinition{
		NewNT(nVal, "Value").Is(tInt),
		NewNT(nSum, "Sum").
			Is(nVal).
			Is(nSum, tPlus, nVal).Do(calc2IntSum).
			Is(nSum, tMinus, nVal).Do(calc2IntSub),
		NewNT(nGoal, "Goal").Main().Is(nSum),
	})

	t.Run("SymbolRegistry", func(t *testing.T) {
		var reg SymbolRegistry = g
		if s := reg.SymbolName(tInt); s != "integer" {
			t.Error("name of tInt: ", s)
		}
		if s := reg.SymbolName(nSum); s != "Sum" {
			t.Error("name of nSum: ", s)
		}

		if s := dumpId(tInt, g); s != "integer" {
			t.Error("name of tInt: ", s)
		}
		if s := dumpId(nSum, g); s != "Sum" {
			t.Error("name of nSum: ", s)
		}
	})

	t.Run("main rule", func(t *testing.T) {
		mr := g.MainRule()
		if mr == nil {
			t.Fatal("no main rule")
		}
		if mr.Subject() != nGoal || !mr.HasEOF() {
			t.Errorf("incorrect main rule: %#v", mr)
		}
		if mrd := mr.Definition(); len(mrd) != 1 || mrd[0] != nSum {
			t.Errorf("incorrect main rule definition: %#v", mrd)
		}
	})

	t.Run("RulesFor", func(t *testing.T) {
		if rv := g.RulesFor(nVal); len(rv) != 1 || rv[0].Definition()[0] != tInt {
			t.Fatalf("incorrect rules for nValue: %#v", rv)
		}

		rs := g.RulesFor(nSum)
		if len(rs) != 3 {
			t.Fatalf("incorrect rules for nSum: %#v", rs)
		}

		if rg := g.RulesFor(nGoal); len(rg) != 1 || rg[0].Definition()[0] != nSum {
			t.Fatalf("incorrect rules for nGoal: %#v", rg)
		}

		t.Run("panic: for terminal", func(t *testing.T) {
			defer testutils.ExpectPanicError(t, ErrDefine, func(t *testing.T, err error) {
				if err.Error() != "no rule - integer is Terminal: invalid definition" {
					t.Fatal("wrong error message:", err)
				}
			})

			g.RulesFor(tInt)
		})

		t.Run("panic: unknown symbol", func(t *testing.T) {
			defer testutils.ExpectPanicError(t, ErrDefine, func(t *testing.T, err error) {
				if err.Error() != "no rule for #6: invalid definition" {
					t.Fatal("wrong error message:", err)
				}
			})

			g.RulesFor(tMul)
		})
	})
}

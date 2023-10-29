package grammar

import (
	"testing"

	"github.com/vovan-ve/go-lr0-parser/internal/testutils"
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

func TestNew(t *testing.T) {
	const (
		tInt symbol.Id = iota + 1
		tPlus
		tMinus
		tDiv
		tMul
	)
	const (
		nValue symbol.Id = iota + 100
		nSum
		nGoal
	)
	var (
		terminals = []lexer.Terminal{
			lexer.NewFunc(tInt, matchDigits),
			lexer.NewFixedStr(tPlus, "+"),
			lexer.NewFixedStr(tMinus, "-"),
		}
	)

	t.Run("panic: multiple main rules", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, symbol.ErrDefine, func(t *testing.T, err error) {
			if err.Error() != "another rule [2] has EOF flag too, previous was [1]: invalid definition" {
				t.Fatal("wrong error message", err)
			}
		})

		New(terminals, []Rule{
			NewRule(nValue, []symbol.Id{tInt}),
			NewRuleMain(nSum, []symbol.Id{nValue}),
			NewRuleMain(nSum, []symbol.Id{nSum, tPlus, nValue}),
			NewRuleMain(nSum, []symbol.Id{nSum, tMinus, nValue}),
		})
	})

	t.Run("panic: rule subject is Terminal", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, symbol.ErrDefine, func(t *testing.T, err error) {
			if err.Error() != "rules[1] subject is Terminal: invalid definition" {
				t.Fatal("wrong error message", err)
			}
		})

		New(terminals, []Rule{
			NewRule(nValue, []symbol.Id{tInt}),
			NewRule(tPlus, []symbol.Id{nValue}),
			NewRule(nSum, []symbol.Id{nSum, tPlus, nValue}),
			NewRule(nSum, []symbol.Id{nSum, tMinus, nValue}),
		})
	})

	t.Run("panic: non-terminal without rule", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, symbol.ErrDefine, func(t *testing.T, err error) {
			const (
				expected = `undefined non-terminals without rules:
- #100 in rules[0] definitions[0]
: invalid definition`
			)
			if err.Error() != expected {
				t.Fatal("wrong error message:", err)
			}
		})

		New(terminals, []Rule{
			NewRule(nSum, []symbol.Id{nValue}),
			NewRule(nSum, []symbol.Id{nSum, tPlus, nValue}),
			NewRule(nSum, []symbol.Id{nSum, tMinus, nValue}),
		})
	})

	const (
		tagSumValue symbol.Tag = iota + 1
		tagSumPlus
		tagSumMinus
	)

	t.Run("panic: no main rule", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, symbol.ErrDefine, func(t *testing.T, err error) {
			if err.Error() != "no main rule with EOF flag: invalid definition" {
				t.Fatal("wrong error message:", err)
			}
		})

		New(terminals, []Rule{
			NewRule(nValue, []symbol.Id{tInt}),
			NewRuleTag(nSum, tagSumValue, []symbol.Id{nValue}),
			NewRuleTag(nSum, tagSumPlus, []symbol.Id{nSum, tPlus, nValue}),
			NewRuleTag(nSum, tagSumMinus, []symbol.Id{nSum, tMinus, nValue}),
			NewRule(nGoal, []symbol.Id{nSum}),
		})
	})

	t.Run("panic: unused terminals", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, symbol.ErrDefine, func(t *testing.T, err error) {
			const (
				expected = `following Terminals are not used in any Rule:
- #4
- MUL
: invalid definition`
			)
			if err.Error() != expected {
				t.Fatal("wrong error message:", err)
			}
		})

		terminals2 := append(
			terminals,
			lexer.NewFixedStr(tDiv, "/"),
			lexer.Name("MUL", lexer.NewFixedStr(tMul, "*")),
		)

		New(terminals2, []Rule{
			NewRule(nValue, []symbol.Id{tInt}),
			NewRuleTag(nSum, tagSumValue, []symbol.Id{nValue}),
			NewRuleTag(nSum, tagSumPlus, []symbol.Id{nSum, tPlus, nValue}),
			NewRuleTag(nSum, tagSumMinus, []symbol.Id{nSum, tMinus, nValue}),
			NewRuleMain(nGoal, []symbol.Id{nSum}),
		})
	})

	g := New(terminals, []Rule{
		NewRule(nValue, []symbol.Id{tInt}),
		NewRuleTag(nSum, tagSumValue, []symbol.Id{nValue}),
		NewRuleTag(nSum, tagSumPlus, []symbol.Id{nSum, tPlus, nValue}),
		NewRuleTag(nSum, tagSumMinus, []symbol.Id{nSum, tMinus, nValue}),
		NewRuleMain(nGoal, []symbol.Id{nSum}),
	})

	t.Run("main rule", func(t *testing.T) {
		mr := g.MainRule()
		if mr == nil {
			t.Fatal("no main rule")
		}
		if mr.Subject() != nGoal || mr.Tag() != 0 || !mr.HasEOF() {
			t.Errorf("incorrect main rule: %#v", mr)
		}
		if mrd := mr.Definition(); len(mrd) != 1 || mrd[0] != nSum {
			t.Errorf("incorrect main rule definition: %#v", mrd)
		}
	})

	t.Run("RulesFor", func(t *testing.T) {
		if rv := g.RulesFor(nValue); len(rv) != 1 || rv[0].Definition()[0] != tInt {
			t.Fatalf("incorrect rules for nValue: %#v", rv)
		}

		rs := g.RulesFor(nSum)
		if len(rs) != 3 || rs[0].Tag() != tagSumValue || rs[1].Tag() != tagSumPlus || rs[2].Tag() != tagSumMinus {
			t.Fatalf("incorrect rules for nSum: %#v", rs)
		}

		if rg := g.RulesFor(nGoal); len(rg) != 1 || rg[0].Definition()[0] != nSum {
			t.Fatalf("incorrect rules for nGoal: %#v", rg)
		}

		t.Run("panic: for terminal", func(t *testing.T) {
			defer testutils.ExpectPanicError(t, symbol.ErrDefine, func(t *testing.T, err error) {
				if err.Error() != "no rule - #1 is Terminal: invalid definition" {
					t.Fatal("wrong error message:", err)
				}
			})

			g.RulesFor(tInt)
		})

		t.Run("panic: unknown symbol", func(t *testing.T) {
			defer testutils.ExpectPanicError(t, symbol.ErrDefine, func(t *testing.T, err error) {
				if err.Error() != "no rule for #5: invalid definition" {
					t.Fatal("wrong error message:", err)
				}
			})

			g.RulesFor(tMul)
		})
	})
}

func matchDigits(state *lexer.State) (next *lexer.State, value any) {
	if state.IsEOF() {
		return
	}
	st, b := state.TakeBytesFunc(isDigit)
	if b == nil {
		return
	}
	next = st
	value = string(state.BytesTo(next))
	return
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

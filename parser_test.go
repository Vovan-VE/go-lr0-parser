package lr0

import (
	"strings"
	"testing"
	"unicode"

	"github.com/pkg/errors"
)

func TestParser(t *testing.T) {
	p := newParser(newGrammar(
		[]Terminal{
			NewTerm(tInt, "int").FuncByte(isDigit, bytesToInt),
			NewTerm(tPlus, `"+"`).Hide().Str("+"),
			NewTerm(tMinus, `"-"`).Hide().Str("-"),
			NewTerm(tMul, `"*"`).Hide().Str("*"),
			NewTerm(tDiv, `"/"`).Hide().Str("/"),

			NewWhitespace().FuncRune(unicode.IsSpace),
		},
		[]NonTerminalDefinition{
			NewNT(nGoal, "Goal").Main().Is(nSum),
			NewNT(nSum, "Sum").
				Is(nSum, tPlus, nProd).Do(func(a, b int) int { return a + b }).
				Is(nSum, tMinus, nProd).Do(func(a, b int) int { return a - b }).
				Is(nProd),
			NewNT(nProd, "Prod").
				Is(nProd, tMul, nVal).Do(func(a, b int) int { return a * b }).
				Is(nProd, tDiv, nVal).Do(
				func(a, b int) (int, error) {
					if b == 0 {
						return 0, errDivZero
					}
					return a / b, nil
				}).
				Is(nVal),
			NewNT(nVal, "Val").Is(tInt),
		},
	))

	t.Run("success", func(t *testing.T) {
		const result = 42*23/3 + 90/15 - 17*19
		v, err := p.Parse(NewState([]byte("42*23/3 + 90/15 - 17*19 ")))
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if v != result {
			t.Fatalf("result is %#v", v)
		}
	})
	t.Run("reduce value error", func(t *testing.T) {
		_, err := p.Parse(NewState([]byte("42/0-5")))
		if !errors.Is(err, errDivZero) {
			t.Error("no expected error:", err)
		}
	})
	// REFACT: never happen: cannot trigger this error in this grammar without whitespace skipping
	//t.Run("no reduce rule, unexpected token", func(t *testing.T) {
	//	_, err := p.Parse(NewState([]byte("42/3**7")))
	//	if !errors.Is(err, ErrParse) {
	//		t.Fatal("wrong error type:", err)
	//	}
	//	if !strings.Contains(err.Error(), "unexpected input 1") {
	//		t.Fatal("wrong error message:", err)
	//	}
	//})
	t.Run("unexpected input 2, unexpected char", func(t *testing.T) {
		_, err := p.Parse(NewState([]byte("42/3*?0")))
		if !errors.Is(err, ErrParse) {
			t.Fatal("wrong error type:", err)
		}
		if err.Error() != "unexpected input: expected int: parse error near ⟪42/3*⟫⏵⟪?0⟫" {
			t.Fatal("wrong error message:", err)
		}
	})
	t.Run("no reduce rule 3, unexpected token `A op > FAIL`", func(t *testing.T) {
		_, err := p.Parse(NewState([]byte("42/3**7")))
		if !errors.Is(err, ErrParse) {
			t.Fatal("wrong error type:", err)
		}
		if err.Error() != "unexpected input: expected int: parse error near ⟪42/3*⟫⏵⟪*7⟫" {
			t.Fatal("wrong error message:", err)
		}
	})
	t.Run("no reduce rule, unexpected token", func(t *testing.T) {
		_, err := p.Parse(NewState([]byte("42/3 7")))
		if !errors.Is(err, ErrParse) {
			t.Fatal("wrong error type:", err)
		}
		if !strings.Contains(err.Error(), "unexpected input instead of EOF: parse error near ⟪42/3⟫⏵⟪␠7⟫") {
			t.Fatal("wrong error message:", err)
		}
	})
}

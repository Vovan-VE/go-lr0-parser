package parser

import (
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

const (
	tInt symbol.Id = iota + 1
	tPlus
	tMinus
	tMul
	tDiv

	nVal
	nProd
	nSum
	nGoal
)

var errDivZero = errors.New("division by zero")

var g = grammar.New(
	[]lexer.Terminal{
		lexer.NewTerm(tInt, "int").Func(matchDigits),
		lexer.NewTerm(tPlus, `"+"`).Hide().Str("+"),
		lexer.NewTerm(tMinus, `"-"`).Hide().Str("-"),
		lexer.NewTerm(tMul, `"*"`).Hide().Str("*"),
		lexer.NewTerm(tDiv, `"/"`).Hide().Str("/"),

		lexer.NewWhitespace().Func(matchWS),
	},
	[]grammar.NonTerminalDefinition{
		grammar.NewNT(nGoal, "Goal").Main().Is(nSum),
		grammar.NewNT(nSum, "Sum").
			Is(nSum, tPlus, nProd).Do(func(a, b int) int { return a + b }).
			Is(nSum, tMinus, nProd).Do(func(a, b int) int { return a - b }).
			Is(nProd),
		grammar.NewNT(nProd, "Prod").
			Is(nProd, tMul, nVal).Do(func(a, b int) int { return a * b }).
			Is(nProd, tDiv, nVal).Do(
			func(a, b int) (int, error) {
				if b == 0 {
					return 0, errDivZero
				}
				return a / b, nil
			}).
			Is(nVal),
		grammar.NewNT(nVal, "Val").Is(tInt),
	},
)

func TestParser(t *testing.T) {
	p := New(g)

	t.Run("success", func(t *testing.T) {
		const result = 42*23/3 + 90/15 - 17*19
		v, err := p.Parse(lexer.NewState([]byte("42*23/3 + 90/15 - 17*19 ")))
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if v != result {
			t.Fatalf("result is %#v", v)
		}
	})
	t.Run("reduce value error", func(t *testing.T) {
		_, err := p.Parse(lexer.NewState([]byte("42/0-5")))
		if !errors.Is(err, errDivZero) {
			t.Error("no expected error:", err)
		}
	})
	// REFACT: never happen: cannot trigger this error in this grammar without whitespace skipping
	//t.Run("no reduce rule, unexpected token", func(t *testing.T) {
	//	_, err := p.Parse(lexer.NewState([]byte("42/3**7")))
	//	if !errors.Is(err, lexer.ErrParse) {
	//		t.Fatal("wrong error type:", err)
	//	}
	//	if !strings.Contains(err.Error(), "unexpected input 1") {
	//		t.Fatal("wrong error message:", err)
	//	}
	//})
	t.Run("unexpected input 2, unexpected char", func(t *testing.T) {
		_, err := p.Parse(lexer.NewState([]byte("42/3*?0")))
		if !errors.Is(err, lexer.ErrParse) {
			t.Fatal("wrong error type:", err)
		}
		if err.Error() != "unexpected input: expected int: parse error near ⟪42/3*⟫⏵⟪?0⟫" {
			t.Fatal("wrong error message:", err)
		}
	})
	t.Run("no reduce rule 3, unexpected token `A op > FAIL`", func(t *testing.T) {
		_, err := p.Parse(lexer.NewState([]byte("42/3**7")))
		if !errors.Is(err, lexer.ErrParse) {
			t.Fatal("wrong error type:", err)
		}
		if err.Error() != "unexpected input: expected int: parse error near ⟪42/3*⟫⏵⟪*7⟫" {
			t.Fatal("wrong error message:", err)
		}
	})
	t.Run("no reduce rule, unexpected token", func(t *testing.T) {
		_, err := p.Parse(lexer.NewState([]byte("42/3 7")))
		if !errors.Is(err, lexer.ErrParse) {
			t.Fatal("wrong error type:", err)
		}
		if !strings.Contains(err.Error(), "unexpected input instead of EOF: parse error near ⟪42/3⟫⏵⟪␠7⟫") {
			t.Fatal("wrong error message:", err)
		}
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

	value, err := strconv.Atoi(string(state.BytesTo(next)))
	if err != nil {
		value = err
	}
	return
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

func matchWS(st *lexer.State) (next *lexer.State, v any) {
	to, _ := st.TakeRunesFunc(unicode.IsSpace)
	if to.Offset() == st.Offset() {
		return nil, nil
	}
	return to, nil
}

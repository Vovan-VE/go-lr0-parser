package lr0parser

import (
	"strconv"
	"strings"
	"testing"

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
		lexer.NewFunc(tInt, matchDigits),
		lexer.Hide(lexer.NewFixedStr(tPlus, "+")),
		lexer.Hide(lexer.NewFixedStr(tMinus, "-")),
		lexer.Hide(lexer.NewFixedStr(tMul, "*")),
		lexer.Hide(lexer.NewFixedStr(tDiv, "/")),
	},
	[]grammar.RuleDefinition{
		grammar.NewRuleMain(nGoal, []symbol.Id{nSum}, nil),
		grammar.NewRule(nSum, []symbol.Id{nSum, tPlus, nProd}, func(a, b int) int { return a + b }),
		grammar.NewRule(nSum, []symbol.Id{nSum, tMinus, nProd}, func(a, b int) int { return a - b }),
		grammar.NewRule(nSum, []symbol.Id{nProd}, nil),
		grammar.NewRule(nProd, []symbol.Id{nProd, tMul, nVal}, func(a, b int) int { return a * b }),
		grammar.NewRule(nProd, []symbol.Id{nProd, tDiv, nVal}, func(a, b int) (int, error) {
			if b == 0 {
				return 0, errDivZero
			}
			return a / b, nil
		}),
		grammar.NewRule(nProd, []symbol.Id{nVal}, nil),
		grammar.NewRule(nVal, []symbol.Id{tInt}, nil),
	},
)

func TestParser(t *testing.T) {
	p := New(g)

	t.Run("success", func(t *testing.T) {
		const result = 42*23/3 + 90/15 - 17*19
		v, err := p.Parse(lexer.NewState([]byte("42*23/3+90/15-17*19")))
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
			t.Error("no expected error", err)
		}
	})
	// TODO: cannot trigger this error in this grammar without whitespace skipping
	//t.Run("no reduce rule, unexpected token", func(t *testing.T) {
	//	_, err := p.Parse(lexer.NewState([]byte("42/3**7")))
	//	if !errors.Is(err, lexer.ErrParse) {
	//		t.Fatal("wrong error type", err)
	//	}
	//	if !strings.Contains(err.Error(), "unexpected input 1") {
	//		t.Fatal("wrong error message", err)
	//	}
	//})
	t.Run("unexpected input 2, unexpected char", func(t *testing.T) {
		_, err := p.Parse(lexer.NewState([]byte("42/3*?0")))
		if !errors.Is(err, lexer.ErrParse) {
			t.Fatal("wrong error type", err)
		}
		if !strings.Contains(err.Error(), "unexpected input 2") {
			t.Fatal("wrong error message", err)
		}
	})
	t.Run("no reduce rule 3, unexpected token `A op > FAIL`", func(t *testing.T) {
		_, err := p.Parse(lexer.NewState([]byte("42/3**7")))
		if !errors.Is(err, lexer.ErrParse) {
			t.Fatal("wrong error type", err)
		}
		if !strings.Contains(err.Error(), "unexpected input 3") {
			t.Fatal("wrong error message", err)
		}
	})
	// TODO: cannot trigger "unexpected input instead of EOF" in this grammar
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

package lr0_test

import (
	"fmt"
	"strconv"
	"testing"
	"unicode"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser"
)

const (
	tInt lr0.Id = iota + 1
	tPlus
	tMinus
	tMul
	tDiv
	tParensOpen
	tParensClose

	nVal
	nProd
	nSum
	nGoal
)

var errDivZero = errors.New("division by zero")

var parser = lr0.New(
	[]lr0.Terminal{
		lr0.NewTerm(tInt, "int").Func(matchDigits),
		lr0.NewTerm(tPlus, `"+"`).Hide().Str("+"),
		lr0.NewTerm(tMinus, `"-"`).Hide().Str("-"),
		lr0.NewTerm(tMul, `"*"`).Hide().Str("*"),
		lr0.NewTerm(tDiv, `"/"`).Hide().Str("/"),
		lr0.NewTerm(tParensOpen, `"("`).Hide().Str("("),
		lr0.NewTerm(tParensClose, `")"`).Hide().Str(")"),

		lr0.NewWhitespace().Func(matchWS),
	},
	[]lr0.NonTerminalDefinition{
		lr0.NewNT(nGoal, "Goal").Main().Is(nSum),
		lr0.NewNT(nSum, "Sum").
			Is(nSum, tPlus, nProd).Do(func(a, b int) int { return a + b }).
			Is(nSum, tMinus, nProd).Do(func(a, b int) int { return a - b }).
			Is(nProd),
		lr0.NewNT(nProd, "Prod").
			Is(nProd, tMul, nVal).Do(func(a, b int) int { return a * b }).
			Is(nProd, tDiv, nVal).Do(
			func(a, b int) (int, error) {
				if b == 0 {
					return 0, errDivZero
				}
				return a / b, nil
			}).
			Is(nVal),
		lr0.NewNT(nVal, "Val").
			Is(tInt).
			Is(tParensOpen, nSum, tParensClose),
	},
)

func TestParser(t *testing.T) {
	type testCase struct {
		input  string
		result int
		err    error
	}
	for i, c := range []testCase{
		{input: "42*23/3 + 90/15 - 17*19 ", result: 42*23/3 + 90/15 - 17*19},
		{input: "42*23/3 + 90/(15-17)*( (19 )) ", result: 42*23/3 + 90/(15-17)*19},

		{input: "(2 + 3) / (3 + 8 - 12 + 1)", err: errDivZero},
		{input: "1+2(3", err: lr0.ErrParse},
	} {
		t.Run(fmt.Sprintf("case %d: %s", i, c.input), func(t *testing.T) {
			v, err := parser.Parse(lr0.NewState([]byte(c.input)))
			if !errors.Is(err, c.err) {
				t.Fatalf("wrong error %v ; expected %v", err, c.err)
			}
			if err == nil && v != c.result {
				t.Fatalf("result is %#v ; expected %v", v, c.result)
			}
		})
	}
}

func matchDigits(state *lr0.State) (next *lr0.State, value any) {
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

func matchWS(st *lr0.State) (next *lr0.State, v any) {
	to, _ := st.TakeRunesFunc(unicode.IsSpace)
	if to.Offset() == st.Offset() {
		return nil, nil
	}
	return to, nil
}

func TestCommentExample1(t *testing.T) {
	p := lr0.New(
		[]lr0.Terminal{
			lr0.NewTerm(tInt, "int").Func(matchDigits),
			lr0.NewTerm(tPlus, `"+"`).Hide().Str("+"),
			lr0.NewWhitespace().Func(matchWS),
		},
		[]lr0.NonTerminalDefinition{
			lr0.NewNT(nGoal, "Goal").Main().Is(nSum),
			lr0.NewNT(nSum, "Sum").
				Is(nSum, tPlus, nVal).Do(func(a, b int) int { return a + b }).
				Is(nVal),
			lr0.NewNT(nVal, "Val").Is(tInt),
		},
	)
	result, err := p.Parse(lr0.NewState([]byte("42 + 37")))
	if err != nil {
		t.Error("error", err)
	}
	if result != 42+37 {
		t.Error("result", result)
	}
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"

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

func main() {
	if len(os.Args) <= 1 {
		log.Println("no args to calc")
		return
	}
	for i, input := range os.Args[1:] {
		fmt.Printf("%d> %s", i, input)
		result, err := parser.Parse(lr0.NewState([]byte(input)))
		if err != nil {
			fmt.Println("\t=> Error:", err)
		} else {
			fmt.Println("\t=>", result)
		}
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

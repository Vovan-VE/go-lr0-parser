package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/vovan-ve/go-lr0-parser"
)

const (
	tInt lr0.Id = iota + 1
	tPlus
	tMinus

	nVal
	nSum
	nGoal
)

var parser = lr0.New(
	[]lr0.Terminal{
		lr0.NewTerm(tInt, "int").FuncByte(isDigit, bytesToInt),
		lr0.NewTerm(tPlus, `"+"`).Hide().Str("+"),
		lr0.NewTerm(tMinus, `"-"`).Hide().Str("-"),
	},
	[]lr0.NonTerminalDefinition{
		lr0.NewNT(nGoal, "Goal").
			Main().
			Is(nSum),
		lr0.NewNT(nSum, "Sum").
			Is(nSum, tPlus, nVal).Do(func(a, b int) int { return a + b }).
			Is(nSum, tMinus, nVal).Do(func(a, b int) int { return a - b }).
			Is(nVal),
		lr0.NewNT(nVal, "Val").
			Is(tInt),
	},
)

func main() {
	if len(os.Args) <= 1 {
		log.Println("no args to calc")
		return
	}
	for i, input := range os.Args[1:] {
		fmt.Printf("%d> %s\t=> ", i, input)
		result, err := parser.Parse(lr0.NewState([]byte(input)))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(result)
		}
	}
}

func isDigit(b byte) bool              { return b >= '0' && b <= '9' }
func bytesToInt(b []byte) (int, error) { return strconv.Atoi(string(b)) }

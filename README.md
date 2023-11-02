# LR(0) Parser

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vovan-ve/go-lr0-parser)
![GitHub release](https://img.shields.io/github/v/release/vovan-ve/go-lr0-parser)
[![License](https://img.shields.io/github/license/vovan-ve/go-lr0-parser)](./LICENSE)

This package contains [LR(0) parser][lr-parser.wiki] to parse text according
to defined LR(0) grammar.

It's based on my previous [PHP library](https://github.com/Vovan-VE/parser),
but with some package API enhancements.

## Example

```go
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
		lr0.NewTerm(tInt, "int").Func(matchDigits),
		lr0.NewTerm(tPlus, `"+"`).Hide().Str("+"),
		lr0.NewTerm(tMinus, `"-"`).Hide().Str("-"),
	},
	[]lr0.NonTerminalDefinition{
		lr0.NewNT(nGoal, "Goal").Main().Is(nSum),
		lr0.NewNT(nSum, "Sum").
			Is(nSum, tPlus, nVal).Do(func(a, b int) int { return a + b }).
			Is(nSum, tMinus, nVal).Do(func(a, b int) int { return a - b }).
			Is(nVal),
		lr0.NewNT(nVal, "Val").Is(tInt),
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
```

```sh
$ go build -o calc examples/01-calc-tiny/main.go
...
$ ./calc "3+8-5" "3+ 8-5" "3+8*5"
0> 3+8-5        => 6
1> 3+ 8-5       => Error: unexpected input: expected int: parse error near ⟪3+⟫⏵⟪␠8-5⟫
2> 3+8*5        => Error: unexpected input: expected "+" or "-": parse error near ⟪3+8⟫⏵⟪*5⟫
```

See examples in [examples/](./examples/) and [tests](./lr0_test.go).

Theory
------

[LR parser][lr-parser.wiki].

License
-------

[MIT][mit]

[lr-parser.wiki]: https://en.wikipedia.org/wiki/LR_parser
[mit]: https://opensource.org/licenses/MIT

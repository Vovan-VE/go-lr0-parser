package grammar

import (
	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

type NonTerminalDefinition interface {
	symbol.Symbol
	GetRules(l lexer.HiddenRegistry) []Rule
}

type NonTerminal struct {
	id          symbol.Id
	name        string
	main        bool
	definitions []nonTerminalDefinition
}

// NewNT created new non-terminal definition
//
//	NewNT(nGoal, "Goal").Main().
//		Is(nSum)
//
//	NewNT(nSum, "Sum").
//		Is(nSum, tPlus, nVal).Do(func (a, b int) int { return a+b }).
//		// ^^^^  ^^^^^  ^^^^           ^  ^
//		//  a    hidden   b   ---------+--'
//		Is(nSum, tMinus, nVal).Do(func (a, b int) int { return a-b }).
//		Is(nVal)
//		// -----^^^^ here is no `Do(func (v int) int { return v })`
//		//           `Do(nil)` will do the same in this case
func NewNT(id symbol.Id, name string) *NonTerminal {
	return &NonTerminal{id: id, name: name}
}

func (n *NonTerminal) Id() symbol.Id { return n.id }

func (n *NonTerminal) Name() string { return n.name }

// Main marks this non-terminal as main
//
// Main non-terminal must have exactly one definition. Exactly one main rule
// must be defined in grammar.
func (n *NonTerminal) Main() *NonTerminal {
	if l := len(n.definitions); l > 1 {
		panic(errors.Wrapf(symbol.ErrDefine, "main non-terminal must have the only definition, here are %d", l))
	}
	n.main = true
	return n
}

// Is adds one more alternative definition for the non-terminal
//
// Is can be followed by `Do()` to define evaluation for this definition.
func (n *NonTerminal) Is(id symbol.Id, ids ...symbol.Id) *NonTerminal {
	if n.main && len(n.definitions) > 0 {
		panic(errors.Wrap(symbol.ErrDefine, "main non-terminal must have the only definition"))
	}
	n.definitions = append(n.definitions, nonTerminalDefinition{
		items: append([]symbol.Id{id}, ids...),
	})
	return n
}

// Do sets a func how to evaluate return value of this non-terminal for the
// latest `Is()` case.
//
// If this definition with respect to hidden terminals evaluates exactly one
// value
//
//	[]any{ /* one value here */ }
//
// then `Do()` subsequent can be omitted. In this case the value of children
// node will be returned from this rule:
//
//	NewNT(nSum, "Sum").
//		Is(nSum, tPlus, nVal).Do(func (a, b int) int { return a+b }).
//		// ^^^^  ^^^^^  ^^^^           ^  ^
//		//  a    hidden   b   ---------+--'
//		Is(nSum, tMinus, nVal).Do(func (a, b int) int { return a-b }).
//		Is(nVal)
//		// -----^^^^ here is no `Do(func (v int) int { return v })`
//		//           `Do(nil)` will do the same in this case
func (n *NonTerminal) Do(calcHandler any) *NonTerminal {
	l := len(n.definitions)
	if l == 0 {
		panic(errors.Wrap(symbol.ErrDefine, "using Do() without Is()"))
	}
	to := &n.definitions[l-1]
	if to.calcHandler != nil {
		panic(errors.Wrap(symbol.ErrDefine, "using Do() again without Is()"))
	}
	to.calcHandler = calcHandler
	return n
}

// GetRules return actual rules built for this non-terminal
func (n *NonTerminal) GetRules(l lexer.HiddenRegistry) []Rule {
	c := len(n.definitions)
	if c == 0 {
		panic(errors.Wrap(symbol.ErrDefine, "no definitions by Is()"))
	}
	res := make([]Rule, 0, c)
	for _, def := range n.definitions {
		res = append(res, newRule(n, n.main, def, l))
	}
	return res
}

type nonTerminalDefinition struct {
	items       []symbol.Id
	calcHandler any
}

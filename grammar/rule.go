package grammar

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

// Rule is one of possible definition for a non-Terminal
type Rule interface {
	fmt.Stringer
	// Subject of the rule
	Subject() symbol.Id
	// HasEOF tells whether EOF must be found in the end of input
	HasEOF() bool
	// Definition of what it consists of
	Definition() []symbol.Id
	Value([]any) (any, error)
	IsHidden(index int) bool
}

func newRule(s symbol.Symbol, main bool, d nonTerminalDefinition, l lexer.HiddenRegistry) *rule {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		err, ok := e.(error)
		if !ok || !errors.Is(err, symbol.ErrDefine) {
			panic(e)
		}
		panic(errors.Wrapf(err, "rule for %s", symbol.Dump(s)))
	}()

	hidden := whichHidden(l, d.items)
	return &rule{
		subject:    s.Id(),
		eof:        main,
		definition: d.items,
		calc:       prepareHandler(d.calcHandler, len(d.items)-len(hidden)),
		hidden:     hidden,
	}
}

type rule struct {
	subject    symbol.Id
	eof        bool
	definition []symbol.Id
	calc       calcFunc
	hidden     map[int]struct{}
}

func (r *rule) Subject() symbol.Id      { return r.subject }
func (r *rule) HasEOF() bool            { return r.eof }
func (r *rule) Definition() []symbol.Id { return r.definition }

// REFACT: names from grammar
func (r *rule) String() string {
	s := fmt.Sprintf("#%d :", r.subject)
	for _, id := range r.definition {
		s += fmt.Sprintf(" #%d", id)
	}
	if r.eof {
		s += " $"
	}
	return s
}

func (r *rule) Value(v []any) (any, error) { return r.calc(v) }

func (r *rule) IsHidden(index int) bool {
	_, ok := r.hidden[index]
	return ok
}

func whichHidden(l lexer.HiddenRegistry, ids []symbol.Id) map[int]struct{} {
	hidden := make(map[int]struct{})
	for i, id := range ids {
		if l.IsHidden(id) {
			hidden[i] = struct{}{}
		}
	}
	return hidden
}

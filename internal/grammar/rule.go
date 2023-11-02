package grammar

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/lexer"
	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
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

func newRule(s symbol.Symbol, main bool, d nonTerminalDefinition, l lexer.NamedHiddenRegistry) *rule {
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
		nameReg:    l,
	}
}

type rule struct {
	subject    symbol.Id
	eof        bool
	definition []symbol.Id
	calc       calcFunc
	hidden     map[int]struct{}
	nameReg    symbol.Registry
}

func (r *rule) Subject() symbol.Id      { return r.subject }
func (r *rule) HasEOF() bool            { return r.eof }
func (r *rule) Definition() []symbol.Id { return r.definition }

func (r *rule) String() string {
	s := symbol.DumpId(r.subject, r.nameReg) + " :"
	for _, id := range r.definition {
		s += " " + symbol.DumpId(id, r.nameReg)
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

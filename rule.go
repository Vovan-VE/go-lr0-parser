package lr0

import (
	"github.com/pkg/errors"
)

func newRule(s Symbol, main bool, d nonTerminalDefinition, l NamedHiddenRegistry) *rule {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		err, ok := e.(error)
		if !ok || !errors.Is(err, ErrDefine) {
			panic(e)
		}
		panic(errors.Wrapf(err, "rule for %s", dumpSymbol(s)))
	}()

	hidden := make(map[int]struct{})
	for i, id := range d.items {
		if l.IsHidden(id) {
			hidden[i] = struct{}{}
		}
	}

	return &rule{
		subject:    s.Id(),
		eof:        main,
		definition: d.items,
		calc:       newCalcFunc(d.calcHandler, len(d.items)-len(hidden)),
		hidden:     hidden,
		nameReg:    l,
	}
}

type rule struct {
	subject    Id
	eof        bool
	definition []Id
	calc       calcFunc
	hidden     map[int]struct{}
	nameReg    SymbolRegistry
}

func (r *rule) Subject() Id      { return r.subject }
func (r *rule) HasEOF() bool     { return r.eof }
func (r *rule) Definition() []Id { return r.definition }

func (r *rule) String() string {
	s := dumpId(r.subject, r.nameReg) + " :"
	for _, id := range r.definition {
		s += " " + dumpId(id, r.nameReg)
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

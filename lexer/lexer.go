package lexer

import (
	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

type HiddenRegistry interface {
	// IsHidden returns true if the given Id is a hidden symbol
	IsHidden(id symbol.Id) bool
}

// Lexer does search predefined Terminals in an input stream.
type Lexer interface {
	HiddenRegistry
	// IsTerminal returns true if the given Id is one of defined Terminal
	IsTerminal(id symbol.Id) bool
	// GetTerminalIdsSet returns new set of all defined Id
	GetTerminalIdsSet() symbol.Set
	// Match tries to parse one of expected Terminals in the given State.
	//
	// Order in `expected` does not matter. Only definition order ot Terminals
	// does matter.
	//
	// If none of expected Terminals matched, it will try to match first of the
	// rest unexpected Terminals.
	//
	// At EOF or when nothing matched, a ErrParse wrapped error will be returned.
	Match(state *State, expected symbol.ReadonlySet) (next *State, m *Match, err error)
}

type termMap = map[symbol.Id]Terminal

type lexer struct {
	list      []Terminal
	terminals termMap
	// TODO: whitespaces - ignore terminals on every match
}

// New creates a new empty Configurable
func New(t ...Terminal) Lexer {
	l := &lexer{
		list:      make([]Terminal, 0, len(t)),
		terminals: make(termMap),
	}
	for _, ti := range t {
		id := ti.Id()
		if id == symbol.InvalidId {
			panic(errors.Wrap(symbol.ErrDefine, "zero id"))
		}
		if prev, exists := l.terminals[id]; exists {
			if prev == ti {
				continue
			}
			panic(errors.Wrapf(symbol.ErrDefine, "redefine terminal %v with %v", symbol.Dump(prev), symbol.Dump(ti)))
		}
		l.list = append(l.list, ti)
		l.terminals[id] = ti
	}
	return l
}

func (l *lexer) IsTerminal(id symbol.Id) bool {
	_, ok := l.terminals[id]
	return ok
}

func (l *lexer) IsHidden(id symbol.Id) bool {
	t, ok := l.terminals[id]
	return ok && t.IsHidden()
}

func (l *lexer) GetTerminalIdsSet() symbol.Set {
	m := symbol.NewSetOfId()
	for id := range l.terminals {
		m.Add(id)
	}
	return m
}

func (l *lexer) Match(state *State, expected symbol.ReadonlySet) (*State, *Match, error) {
	if !state.IsEOF() {
		var (
			altNext *State
			altM    *Match
		)
		for _, t := range l.list {
			nextS, v := t.Match(state)
			if v != nil {
				if err, ok := v.(error); ok {
					return nil, nil, err
				}
			}
			if nextS != nil {
				m2 := &Match{
					Term:  t.Id(),
					Value: v,
				}
				if expected.Has(t.Id()) {
					return nextS, m2, nil
				} else if altNext == nil {
					altNext = nextS
					altM = m2
				}
			}
		}

		if altNext != nil {
			return altNext, altM, nil
		}
	}
	return nil, nil, WithSource(expectationError(expected, l.list), state)
}

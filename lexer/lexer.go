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
	GetTerminalIdsSet() symbol.SetOfId
	// Match tries to parse one of expected Terminals in the given State.
	//
	// Order in `expected` does matter. First matched Terminal will be returned in
	// Match with next State for further parsing.
	//
	// At EOF or when none of `expected` Terminals matched, a ErrParse wrapped error
	// will be returned.
	Match(state *State, expected []symbol.Id) (next *State, m *Match, err error)
}

type termMap = map[symbol.Id]Terminal

type lexer struct {
	terminals termMap
}

// New creates a new empty Configurable
func New(t ...Terminal) Lexer {
	l := &lexer{
		terminals: make(termMap),
	}
	for _, ti := range t {
		l.add(ti)
	}
	return l
}

func (l *lexer) add(t Terminal) {
	id := t.Id()
	if id == symbol.InvalidId {
		panic(errors.Wrap(symbol.ErrDefine, "zero id"))
	}
	if prev, exists := l.terminals[id]; exists {
		if prev == t {
			return
		}
		panic(errors.Wrapf(symbol.ErrDefine, "redefine terminal %v with %v", symbol.Dump(prev), symbol.Dump(t)))
	}
	l.terminals[id] = t
}

func (l *lexer) IsTerminal(id symbol.Id) bool {
	_, ok := l.terminals[id]
	return ok
}

func (l *lexer) IsHidden(id symbol.Id) bool {
	t, ok := l.terminals[id]
	return ok && t.IsHidden()
}

func (l *lexer) GetTerminalIdsSet() symbol.SetOfId {
	m := symbol.NewSetOfId()
	for id := range l.terminals {
		m.Add(id)
	}
	return m
}

func (l *lexer) Match(state *State, expected []symbol.Id) (next *State, m *Match, err error) {
	if !state.IsEOF() {
		var (
			t     Terminal
			value any
			ok    bool
		)
		// FIXME: Order may be important, so enhance to return the longest match.
		//   Maybe Terminal can optionally have MatchLength() to return cont length if available.
		for _, expect := range expected {
			t = l.terminals[expect]
			next, value = t.Match(state)
			if err, ok = value.(error); ok {
				return
			}
			if next != nil {
				m = &Match{
					Term:  expect,
					Value: value,
				}
				return
			}
		}
	}
	err = WithSource(expectationError(expected, l.terminals), state)
	return
}

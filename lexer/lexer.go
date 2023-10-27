package lexer

import (
	"github.com/pkg/errors"
)

type termMap = map[Term]Terminal

// Lexer describes a set of possible Terminal to parse any appropriate input
// into defined Terminal tokens
type Lexer struct {
	terminals termMap
}

// New creates a new empty Lexer
func New() *Lexer {
	return &Lexer{}
}

// NewFromMap creates new Lexer with all the given Terminals already defined
func NewFromMap(m termMap) *Lexer {
	return &Lexer{terminals: m}
}

func (l *Lexer) init() {
	if l.terminals == nil {
		l.terminals = make(termMap)
	}
}

// Add adds a Terminal definition and returns the Lexer itself for chaining
func (l *Lexer) Add(id Term, t Terminal) *Lexer {
	l.init()
	l.add(id, t)
	return l
}

// AddMap adds multiple Terminals and returns the Lexer itself for chaining
func (l *Lexer) AddMap(m termMap) *Lexer {
	l.init()
	for id, t := range m {
		l.add(id, t)
	}
	return l
}

func (l *Lexer) add(id Term, t Terminal) {
	if prev, exists := l.terminals[id]; exists {
		if prev == t {
			return
		}
		panic(errors.Wrapf(ErrDefine, "redefine terminal %v with %v", dumpT(id, prev), dumpT(id, t)))
	}
	l.terminals[id] = t
}

// Match tries to parse one of expected Terminals in the given State.
//
// Order in `expected` does matter. First matched Terminal will be returned in
// Match with next State for further parsing.
//
// At EOF or when none of `expected` Terminals matched, a ErrParse wrapped error
// will be returned.
func (l *Lexer) Match(state *State, expected []Term) (next *State, m *Match, err error) {
	if !state.IsEOF() {
		var (
			t     Terminal
			value any
			ok    bool
		)
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

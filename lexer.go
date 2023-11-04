package lr0

import (
	"io"

	"github.com/pkg/errors"
)

//// Lexer does search predefined Terminals in an input stream.
//type Lexer interface {
//	NamedHiddenRegistry
//	// IsTerminal returns true if the given Id is one of defined Terminal
//	IsTerminal(id Id) bool
//	// GetTerminalIdsSet returns new set of all defined Id
//	GetTerminalIdsSet() idSet
//	// Match tries to parse one of expected Terminals in the given State.
//	//
//	// Order in `expected` does not matter. Only definition order ot Terminals
//	// does matter.
//	//
//	// If none of expected Terminals matched, it will try to match first of the
//	// rest unexpected Terminals.
//	//
//	// At EOF the final eof State and io.EOF will be returned.
//	//
//	// When nothing matched, a ErrParse wrapped error will be returned.
//	Match(state *State, expected readonlyIdSet) (next *State, m *Match, err error)
//	ExpectationError(expected readonlyIdSet, pre string) error
//}

type termMap = map[Id]Terminal

type lexer struct {
	list          []Terminal
	terminals     termMap
	internalTerms map[Id][]Terminal
}

// newLexer creates a new empty Configurable
func newLexer(t ...Terminal) *lexer {
	l := &lexer{
		list:          make([]Terminal, 0, len(t)),
		terminals:     make(termMap),
		internalTerms: make(map[Id][]Terminal),
	}
	for _, ti := range t {
		id := ti.Id()
		if id == InvalidId {
			panic(errors.Wrap(ErrDefine, "zero id"))
		}
		if id < 0 {
			prev, _ := l.internalTerms[id]
			l.internalTerms[id] = append(prev, ti)
			continue
		}
		if prev, exists := l.terminals[id]; exists {
			if prev == ti {
				continue
			}
			panic(errors.Wrapf(ErrDefine, "redefine terminal %v with %v", dumpSymbol(prev), dumpSymbol(ti)))
		}
		l.list = append(l.list, ti)
		l.terminals[id] = ti
	}
	return l
}

func (l *lexer) SymbolName(id Id) string {
	if s, ok := l.terminals[id]; ok {
		return s.Name()
	}
	return ""
}

func (l *lexer) IsTerminal(id Id) bool {
	_, ok := l.terminals[id]
	return ok
}

func (l *lexer) IsHidden(id Id) bool {
	t, ok := l.terminals[id]
	return ok && t.IsHidden()
}

func (l *lexer) GetTerminalIdsSet() idSet {
	m := newIdSet()
	for id := range l.terminals {
		m.Add(id)
	}
	return m
}

func (l *lexer) Match(state *State, expected readonlyIdSet) (*State, *Match, error) {
	state = l.skipWhitespaces(state)
	if state.IsEOF() {
		return state, nil, io.EOF
	}
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
			}
			if altNext == nil {
				altNext = nextS
				altM = m2
			}
		}
	}
	if altNext != nil {
		return altNext, altM, nil
	}

	return nil, nil, WithSource(l.ExpectationError(expected, ""), state)
}

func (l *lexer) ExpectationError(expected readonlyIdSet, pre string) error {
	s := pre
	if s != "" {
		s += ": "
	}
	s += "expected "
	i, last := 0, expected.Count()-1
	for _, t := range l.list {
		if !expected.Has(t.Id()) {
			continue
		}
		if i > 0 {
			if i < last {
				s += ", "
			} else {
				s += " or "
			}
		}
		i++
		s += dumpSymbol(t)
	}
	return NewParseError(s)
}

func (l *lexer) skipWhitespaces(state *State) (next *State) {
	next = state
	wsList, ok := l.internalTerms[tWhitespace]
	if !ok {
		return
	}
WsType:
	for !next.IsEOF() {
		for _, ws := range wsList {
			to, _ := ws.Match(next)
			// is this ws type matched, move further and retry if we will
			// find another ws type
			if to != nil {
				next = to
				continue WsType
			}
		}
		// no ws here, done
		break
	}
	return
}

package lr0

import (
	"github.com/pkg/errors"
)

const (
	tWhitespace Id = -iota - 1
)

var (
	metaWS = term{id: tWhitespace, name: "whitespace"}
)

type TerminalFactory struct {
	term
	//Hide() TerminalFactory
	//Reveal() TerminalFactory
	//Terminal() Terminal
}

// NewTerm starts new Terminal creation.
//
//	NewTerm(tInt, "integer").Func(matchDigits)
//	NewTerm(tPlus, "plus").Hide().Byte('+')
func NewTerm(id Id, name string) *TerminalFactory {
	return &TerminalFactory{
		term: term{id: id, name: name},
	}
}

// NewWhitespace can be used to define internal terminals to skip whitespaces
//
// Can be used multiple times to define different kinds of whitespaces.
//
// Whitespace tokens will be silently skipped before every terminal match
func NewWhitespace() *TerminalFactory {
	return &TerminalFactory{
		term: metaWS,
	}
}

// Hide sets "is hidden" flag for further Terminal `IsHidden()` result.
func (t *TerminalFactory) Hide() *TerminalFactory {
	t.hide = true
	return t
}

//func (t *TerminalFactory) Reveal() *TerminalFactory {
//	t.hide = false
//	return t
//}

// Byte creates a Terminal to match the given sequence of bytes exactly.
//
// On match the returned value is matched bytes
//
//	NewTerm(tInc, "increment").Byte('+', '+')
//
//	NewTerm(tPlus, "plus").Byte('+')
func (t *TerminalFactory) Byte(b byte, more ...byte) Terminal {
	return &termFixed{
		term: t.term,
		b:    append([]byte{b}, more...),
	}
}

// Bytes creates a Terminal to match the given sequence of bytes exactly.
//
// On match the returned value is matched bytes
func (t *TerminalFactory) Bytes(b []byte) Terminal {
	if len(b) == 0 {
		panic(errors.Wrap(ErrDefine, "empty bytes slice"))
	}
	return &termFixed{
		term: t.term,
		b:    b,
	}
}

// Str creates a Terminal to match the given substring exactly.
//
// On match the returned value is matched substring
//
//	NewTerm(tInc, "increment").Str("++")
func (t *TerminalFactory) Str(s string) Terminal {
	if s == "" {
		panic(errors.Wrap(ErrDefine, "empty string"))
	}
	return &termFixed{
		term: t.term,
		b:    []byte(s),
		v:    toString,
	}
}

// Func wraps a MatchFunc to Terminal
func (t *TerminalFactory) Func(fn MatchFunc) Terminal {
	return &termCallback{
		term: t.term,
		fn:   fn,
	}
}

type termFixed struct {
	term
	b []byte
	v func([]byte) any
}

func (f *termFixed) Match(state *State) (next *State, value any) {
	if state.IsEOF() {
		return
	}
	next, ok := state.ExpectByteOk(f.b...)
	if !ok {
		return nil, nil
	}

	b := state.BytesTo(next)
	if f.v != nil {
		value = f.v(b)
	} else {
		value = b
	}
	return
}

type termCallback struct {
	term
	fn MatchFunc
}

func (c *termCallback) Match(state *State) (*State, any) {
	return c.fn(state)
}

type term struct {
	id   Id
	name string
	hide bool
}

func (m *term) Id() Id         { return m.id }
func (m *term) Name() string   { return m.name }
func (m *term) IsHidden() bool { return m.hide }

func toString(b []byte) any { return string(b) }

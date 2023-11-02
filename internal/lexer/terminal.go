package lexer

import (
	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
)

// ToValue returns either a value for parsed token, or `error`. Valid value
// cannot be `error`. May return input slice too.
type ToValue = func([]byte) any

func toString(b []byte) any { return string(b) }

// MatchFunc is a signature for common function to match an underlying token.
//
// It accepts current State to parse from.
//
// If the token parsed, the function returns next State to continue from and
// the token value from ToValue.
//
// If the token was not parsed, the function returns `nil, nil`.
//
// Must not return the same State as input State.
type MatchFunc = func(*State) (next *State, value any)

// Terminal is interface to parse specific type of token from input State
type Terminal interface {
	symbol.Symbol
	// IsHidden returns whether the terminal is hidden
	//
	// Hidden terminal does not produce a value to calc non-terminal value.
	// For example if in the following rule:
	//	Sum : Sum plus Val
	// a `plus` terminal is hidden, then only two values will be passed to calc
	// function - value of `Sum` and value of `Val`:
	//	func(sum any, val any) any
	IsHidden() bool
	// Match is MatchFunc
	Match(*State) (next *State, value any)
}

type TerminalFactory struct {
	meta
	//Hide() TerminalFactory
	//Reveal() TerminalFactory
	//Terminal() Terminal
}

// NewTerm starts new Terminal creation.
//
//	NewTerm(tInt, "integer").Func(matchDigits)
//	NewTerm(tPlus, "plus").Hide().Byte('+')
func NewTerm(id symbol.Id, name string) *TerminalFactory {
	return &TerminalFactory{
		meta: meta{id: id, name: name},
	}
}

// NewWhitespace can be used to define internal terminals to skip whitespaces
//
// Can be used multiple times to define different kinds of whitespaces.
//
// Whitespace tokens will be silently skipped before every terminal match
func NewWhitespace() *TerminalFactory {
	return &TerminalFactory{
		meta: metaWS,
	}
}

const (
	tWhitespace symbol.Id = -iota - 1
)

var (
	metaWS = meta{id: tWhitespace, name: "whitespace"}
)

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
	return &fixed{
		meta: t.meta,
		b:    append([]byte{b}, more...),
	}
}

// Bytes creates a Terminal to match the given sequence of bytes exactly.
//
// On match the returned value is matched bytes
func (t *TerminalFactory) Bytes(b []byte) Terminal {
	if len(b) == 0 {
		panic(errors.Wrap(symbol.ErrDefine, "empty bytes slice"))
	}
	return &fixed{
		meta: t.meta,
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
		panic(errors.Wrap(symbol.ErrDefine, "empty string"))
	}
	return &fixed{
		meta: t.meta,
		b:    []byte(s),
		v:    toString,
	}
}

// Func wraps a MatchFunc to Terminal
func (t *TerminalFactory) Func(fn MatchFunc) Terminal {
	return &callback{
		meta: t.meta,
		fn:   fn,
	}
}

type fixed struct {
	meta
	b []byte
	v ToValue
}

func (f *fixed) Match(state *State) (next *State, value any) {
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

func (f *fixed) copy() Terminal {
	ret := *f
	return &ret
}

type callback struct {
	meta
	fn MatchFunc
}

func (c *callback) Match(state *State) (*State, any) {
	return c.fn(state)
}

func (c *callback) copy() Terminal {
	ret := *c
	return &ret
}

type meta struct {
	id   symbol.Id
	name string
	hide bool
}

func (m *meta) Id() symbol.Id  { return m.id }
func (m *meta) Name() string   { return m.name }
func (m *meta) IsHidden() bool { return m.hide }

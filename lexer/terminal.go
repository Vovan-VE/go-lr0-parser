package lexer

import (
	"github.com/vovan-ve/go-lr0-parser/symbol"
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

// Terminal is base interface to parse specific type of token from input State
type Terminal interface {
	symbol.Meta
	// Match is MatchFunc
	Match(*State) (next *State, value any)
}

// NewFixed creates a Terminal to match the given sequence of bytes exactly.
//
// On match the returned value is matched bytes
func NewFixed(id symbol.Id, b []byte) Terminal {
	return &fixed{b: b, meta: meta{id: id}}
}

// NewFixedStr creates a Terminal to match the given substring exactly.
//
// On match the returned value is matched substring
func NewFixedStr(id symbol.Id, s string) Terminal {
	return &fixed{b: []byte(s), v: toString, meta: meta{id: id}}
}

type fixed struct {
	b []byte
	v ToValue
	meta
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

// NewFunc wraps a MatchFunc to Terminal
func NewFunc(id symbol.Id, fn MatchFunc) Terminal {
	return &callback{fn: fn, meta: meta{id: id}}
}

type callback struct {
	fn MatchFunc
	meta
}

func (c *callback) Match(state *State) (*State, any) {
	return c.fn(state)
}

func (c *callback) copy() Terminal {
	ret := *c
	return &ret
}

// Name creates a copy of the given Terminal with Name() changed to the given one.
//
// Works only with Terminal created by this module.
func Name(name string, t Terminal) Terminal {
	if t.Name() == name {
		return t
	}
	cpy := t.(terminalCopier).copy()
	cpy.(terminalCustomizer).setName(name)
	return cpy
}

// Hide creates a copy of the given Terminal with IsHidden() changed to `true`.
//
// Works only with Terminal created by this module.
func Hide(t Terminal) Terminal {
	if t.IsHidden() {
		return t
	}
	cpy := t.(terminalCopier).copy()
	cpy.(terminalCustomizer).setHidden(true)
	return cpy
}

type terminalCopier interface {
	copy() Terminal
}
type terminalCustomizer interface {
	setName(name string)
	setHidden(hidden bool)
}

type meta struct {
	id   symbol.Id
	name string
	hide bool
}

func (m *meta) Id() symbol.Id  { return m.id }
func (m *meta) Name() string   { return m.name }
func (m *meta) IsHidden() bool { return m.hide }

func (m *meta) setName(name string)   { m.name = name }
func (m *meta) setHidden(hidden bool) { m.hide = hidden }

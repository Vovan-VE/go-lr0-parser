package lexer

import (
	"testing"

	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
)

const tTemp symbol.Id = 17

func TestFixedBytes_Match(t *testing.T) {
	spaceship := NewTerm(tTemp, "spaceship").Byte('<', '=', '>')

	state := NewState([]byte("42 <=> 37"))

	if next, v := spaceship.Match(state); next != nil {
		t.Errorf("unexpected match to next %v with value %v", next, v)
	}
	if next, v := spaceship.Match(state.to(9000)); next != nil {
		t.Errorf("unexpected match to next %v with value %v", next, v)
	}

	next, v := spaceship.Match(state.FF(3))
	if next == nil {
		t.Error("failed match")
	}
	b, ok := v.([]byte)
	if !ok {
		t.Errorf("value is %#v", v)
	}
	if string(b) != "<=>" {
		t.Errorf("value is %q", b)
	}

	if next.Offset() != 6 {
		t.Errorf("next is %v", next)
	}
}

func TestFixedStr_Match(t *testing.T) {
	spaceship := NewTerm(tTemp, "spaceship").Str("<=>")

	state := NewState([]byte("42 <=> 37"))

	if next, v := spaceship.Match(state); next != nil {
		t.Errorf("unexpected match to next %v with value %v", next, v)
	}
	if next, v := spaceship.Match(state.to(9000)); next != nil {
		t.Errorf("unexpected match to next %v with value %v", next, v)
	}

	next, v := spaceship.Match(state.FF(3))
	if next == nil {
		t.Error("failed match")
	}
	b, ok := v.(string)
	if !ok {
		t.Errorf("value is %#v", v)
	}
	if b != "<=>" {
		t.Errorf("value is %q", b)
	}

	if next.Offset() != 6 {
		t.Errorf("next is %v", next)
	}
}

func TestFunc_Match(t *testing.T) {
	tIdent := NewTerm(tTemp, "ident").Func(matchIdentifier)

	a := NewState(source)

	if b, v := tIdent.Match(a.to(9000)); b != nil {
		t.Errorf("b: unexpected match %v, %q", b, v)
	}

	c, v := tIdent.Match(a)
	if c == nil {
		t.Error("c: no match")
	}
	if c.Offset() != 5 {
		t.Errorf("c: offset is %v: %v", c.Offset(), c)
	}
	cs, ok := v.(string)
	if !ok {
		t.Errorf("c: v is %#v", v)
	}
	if cs != "Lorem" {
		t.Errorf("c: v string is %q", cs)
	}
}

func matchIdentifier(state *State) (next *State, value any) {
	if state.IsEOF() {
		return
	}
	if next, _ = state.TakeByteFunc(isAlpha); next == nil {
		return
	}
	next, _ = next.TakeBytesFunc(isAlphaNum)
	value = string(state.BytesTo(next))
	return
}

func matchDigits(state *State) (next *State, value any) {
	if state.IsEOF() {
		return
	}
	st, b := state.TakeBytesFunc(isDigit)
	if b == nil {
		return
	}
	next = st
	value = string(state.BytesTo(next))
	return
}

func isAlphaNum(b byte) bool { return isAlpha(b) || isDigit(b) }

func isAlpha(b byte) bool {
	switch {
	case b >= 'A' && b <= 'Z', b >= 'a' && b <= 'z', b == '_':
		return true
	default:
		return false
	}
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

func TestName(t *testing.T) {
	const name = "T_IDENT"
	a := NewTerm(tTemp, name).Func(matchIdentifier)
	if a.Name() != name {
		t.Errorf("a name was set to %q", a.Name())
	}
}

func TestHide(t *testing.T) {
	a := NewTerm(tTemp, "ident").Hide().Func(matchIdentifier)
	if !a.IsHidden() {
		t.Error("a is not hidden")
	}
}
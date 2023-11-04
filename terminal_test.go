package lr0

import (
	"testing"
)

const tTemp Id = 17

func TestTermFixedBytes_Match(t *testing.T) {
	const name = "spaceship"
	spaceship := NewTerm(tTemp, name).Byte('<', '=', '>')
	if spaceship.Name() != name {
		t.Errorf("a name was set to %q", spaceship.Name())
	}

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

func TestTermFixedStr_Match(t *testing.T) {
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

func TestTermFunc_Match(t *testing.T) {
	ident := NewTerm(tTemp, "ident").Func(matchIdentifier)

	a := NewState(testStateSource)

	if b, v := ident.Match(a.to(9000)); b != nil {
		t.Errorf("b: unexpected match %v, %q", b, v)
	}

	c, v := ident.Match(a)
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

func TestTerminalHide(t *testing.T) {
	a := NewTerm(tTemp, "ident").Hide().Func(matchIdentifier)
	if !a.IsHidden() {
		t.Error("a is not hidden")
	}
}

package lexer

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/testutils"
)

func TestLexer_Add(t *testing.T) {
	const (
		tIdent Term = iota + 1
		tPlus
		tMinus
	)

	l := New()
	l.Add(tIdent, Name("Identifier", NewFunc(matchIdentifier)))
	if len(l.terminals) != 1 {
		t.Errorf("where is it? %#v", l.terminals)
	}
	if x := l.terminals[tIdent]; x.Name() != "Identifier" {
		t.Errorf("name %q don't match - %#v", x.Name(), x)
	}

	m := l.
		Add(tPlus, Name("Plus", NewFixedStr("+"))).
		Add(tMinus, Name("Minus", NewFixedStr("-")))
	if m != l {
		t.Error("it's something else")
	}

	if len(l.terminals) != 3 {
		t.Errorf("what just happened? %#v", l.terminals)
	}

	t.Run("panic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		l.Add(tPlus, NewFixedStr("+="))
	})
}

func TestLexer_AddMap(t *testing.T) {
	const (
		tIdent Term = iota + 1
		tPlusAssign
		tPlus
		tMinusAssign
		tMinus
		tDiv
	)

	l := NewFromMap(termMap{
		tIdent: Name("Identifier", NewFunc(matchIdentifier)),
		tPlus:  Name("Plus", NewFixedStr("+")),
		tMinus: Name("Minus", NewFixedStr("-")),
	})
	if len(l.terminals) != 3 {
		t.Errorf("where are they? %#v", l.terminals)
	}

	l.AddMap(map[Term]Terminal{
		tPlusAssign:  Name("PlusAssign", NewFixedStr("+=")),
		tMinusAssign: Name("MinusAssign", NewFixedStr("-=")),
	})
	if len(l.terminals) != 5 {
		t.Errorf("what now? %#v", l.terminals)
	}

	t.Run("panic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		l.AddMap(map[Term]Terminal{
			tDiv:        Name("Div", NewFixedStr("/")),
			tPlusAssign: Name("PlusAssign dupe", NewFixedStr("+=")),
		})
	})
}

func TestLexer_Match(t *testing.T) {
	const (
		tInt Term = iota + 1
		tPlus
		tMinus
		tInc
	)

	l := New().
		Add(tInt, Name("Int", NewFunc(matchDigits))).
		Add(tPlus, Name("Plus", NewFixedStr("+"))).
		Add(tMinus, Name("Minus", NewFixedStr("-"))).
		Add(tInc, Name("Increment", NewFixedStr("++")))

	start := NewState([]byte("38+23-19++"))

	a, m, err := l.Match(start, []Term{tInt})
	if err != nil {
		t.Fatalf("a: match failed: %+v", err)
	}
	if m.Term != tInt {
		t.Fatal("a: match term:", m.Term)
	}
	if v, ok := m.Value.(string); !ok || v != "38" {
		t.Fatal("a: match value:", m.Value)
	}
	if a.Offset() != 2 {
		t.Fatal("a: offset", a.Offset())
	}

	b, m, err := l.Match(a, []Term{tPlus, tMinus})
	if err != nil {
		t.Fatalf("b: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tPlus || !ok || v != "+" {
		t.Fatalf("b: match wrong: %+v", m)
	}
	if b.Offset() != 3 {
		t.Fatal("b: offset", b.Offset())
	}

	c, m, err := l.Match(b, []Term{tInt})
	if err != nil {
		t.Fatalf("c: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tInt || !ok || v != "23" {
		t.Fatalf("c: match wrong: %+v", m)
	}
	if c.Offset() != 5 {
		t.Fatal("c: offset", c.Offset())
	}

	d, m, err := l.Match(c, []Term{tPlus, tMinus})
	if err != nil {
		t.Fatalf("d: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tMinus || !ok || v != "-" {
		t.Fatalf("d: match wrong: %+v", m)
	}
	if d.Offset() != 6 {
		t.Fatal("d: offset", d.Offset())
	}

	e, m, err := l.Match(d, []Term{tInt})
	if err != nil {
		t.Fatalf("e: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tInt || !ok || v != "19" {
		t.Fatalf("e: match wrong: %+v", m)
	}
	if e.Offset() != 8 {
		t.Fatal("e: offset", e.Offset())
	}

	// ++ first, then +
	f1, m, err := l.Match(e, []Term{tInc, tPlus})
	if err != nil {
		t.Fatalf("f1: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tInc || !ok || v != "++" {
		t.Fatalf("f1: match wrong: %+v", m)
	}
	if f1.Offset() != 10 {
		t.Fatal("f1: offset", f1.Offset())
	}

	// + first, then ++
	f2, m, err := l.Match(e, []Term{tPlus, tInc})
	if err != nil {
		t.Fatalf("f2: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tPlus || !ok || v != "+" {
		t.Fatalf("f2: match wrong: %+v", m)
	}
	if f2.Offset() != 9 {
		t.Fatal("f2: offset", f2.Offset())
	}

	// no match
	_, _, err = l.Match(e, []Term{tMinus, tInt})
	if !errors.Is(err, ErrParse) {
		t.Fatal("no expected: wrong error", err)
	}
	const (
		expectStr     = `expected Minus or Int: parse error near ⟪38+23-19⟫⏵⟪++⟫`
		expectStrPlus = `expected Minus or Int: parse error near:
38+23-19++
--------^
`
	)
	if fmt.Sprintf("%v", err) != expectStr {
		t.Errorf("err: <<<<%s>>>>", err)
	}
	if fmt.Sprintf("%+v", err) != expectStrPlus {
		t.Errorf("err+ wrong: <<<<%+v>>>>", err)
	}

	// EOF
	_, _, err = l.Match(start.to(9000), []Term{tInt, tPlus, tMinus, tInc})
	if !errors.Is(err, ErrParse) {
		t.Fatal("eof: wrong error", err)
	}
	const (
		expectEofStr     = `expected Int, Plus, Minus or Increment: parse error near ⟪38+23-19++⟫⏵<EOF>`
		expectEofStrPlus = `expected Int, Plus, Minus or Increment: parse error near:
38+23-19++<EOF>
----------^
`
	)
	if fmt.Sprintf("%v", err) != expectEofStr {
		t.Errorf("err: <<<<%s>>>>", err)
	}
	if fmt.Sprintf("%+v", err) != expectEofStrPlus {
		t.Errorf("err+ wrong: <<<<%+v>>>>", err)
	}
}
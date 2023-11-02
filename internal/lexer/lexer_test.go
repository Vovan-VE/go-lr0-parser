package lexer

import (
	"io"
	"testing"
	"unicode"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
	"github.com/vovan-ve/go-lr0-parser/internal/testutils"
)

func TestLexer_New(t *testing.T) {
	const (
		tIdent symbol.Id = iota + 1
		tPlus
		tMinus
	)

	l := New(
		NewTerm(tIdent, "Identifier").Func(matchIdentifier),
	).(*lexer)
	if len(l.terminals) != 1 {
		t.Errorf("where is it? %#v", l.terminals)
	}
	if x := l.terminals[tIdent]; x.Name() != "Identifier" {
		t.Errorf("name %q don't match - %#v", x.Name(), x)
	}

	l = New(
		NewTerm(tIdent, "Identifier").Func(matchIdentifier),
		NewTerm(tPlus, "Plus").Str("+"),
		NewTerm(tMinus, "Minus").Str("-"),

		NewWhitespace().Func(matchWS),
	).(*lexer)
	if len(l.terminals) != 3 {
		t.Errorf("what just happened? %#v", l.terminals)
	}
	if len(l.internalTerms) != 1 {
		t.Errorf("what just happened? %#v", l.internalTerms)
	}

	t.Run("panic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, symbol.ErrDefine)
		New(
			NewTerm(tPlus, "Plus").Str("+"),
			NewTerm(tPlus, "PlusAssign").Str("+="),
		)
	})
}

func TestLexer_Match(t *testing.T) {
	const (
		tInt symbol.Id = iota + 1
		tPlus
		tMinus
		tInc
	)

	l := New(
		NewTerm(tInt, "Int").Func(matchDigits),
		// ++ first, + after
		NewTerm(tInc, "Increment").Str("++"),
		NewTerm(tPlus, "Plus").Str("+"),
		NewTerm(tMinus, "Minus").Str("-"),

		NewWhitespace().Func(matchWS),
	)

	start := NewState([]byte("38+23 - \n 19++"))

	a, m, err := l.Match(start, symbol.NewSetOfId(tInt))
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

	b, m, err := l.Match(a, symbol.NewSetOfId(tPlus, tMinus))
	if err != nil {
		t.Fatalf("b: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tPlus || !ok || v != "+" {
		t.Fatalf("b: match wrong: %+v", m)
	}
	if b.Offset() != 3 {
		t.Fatal("b: offset", b.Offset())
	}

	c, m, err := l.Match(b, symbol.NewSetOfId(tInt))
	if err != nil {
		t.Fatalf("c: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tInt || !ok || v != "23" {
		t.Fatalf("c: match wrong: %+v", m)
	}
	if c.Offset() != 5 {
		t.Fatal("c: offset", c.Offset())
	}

	d, m, err := l.Match(c, symbol.NewSetOfId(tPlus, tMinus))
	if err != nil {
		t.Fatalf("d: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tMinus || !ok || v != "-" {
		t.Fatalf("d: match wrong: %+v", m)
	}
	if d.Offset() != 7 {
		t.Fatal("d: offset", d.Offset())
	}

	e, m, err := l.Match(d, symbol.NewSetOfId(tInt))
	if err != nil {
		t.Fatalf("e: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tInt || !ok || v != "19" {
		t.Fatalf("e: match wrong: %+v", m)
	}
	if e.Offset() != 12 {
		t.Fatal("e: offset", e.Offset())
	}

	// ++ first, then +
	f1, m, err := l.Match(e, symbol.NewSetOfId(tInc, tPlus))
	if err != nil {
		t.Fatalf("f1: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tInc || !ok || v != "++" {
		t.Fatalf("f1: match wrong: %+v", m)
	}
	if f1.Offset() != 14 {
		t.Fatal("f1: offset", f1.Offset())
	}

	// + first, then ++ - anyway ++ match first due to definition order
	f2, m, err := l.Match(e, symbol.NewSetOfId(tPlus, tInc))
	if err != nil {
		t.Fatalf("f2: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tInc || !ok || v != "++" {
		t.Fatalf("f2: match wrong: %+v", m)
	}
	if f2.Offset() != 14 {
		t.Fatal("f2: offset", f2.Offset())
	}

	// no match
	f3, m, err := l.Match(e, symbol.NewSetOfId(tMinus, tInt))
	if err != nil {
		t.Fatalf("f3: match failed: %+v", err)
	}
	if v, ok := m.Value.(string); m.Term != tInc || !ok || v != "++" {
		t.Fatalf("f3: match wrong: %+v", m)
	}
	if f3.Offset() != 14 {
		t.Fatal("f3: offset", f3.Offset())
	}

	// EOF
	_, _, err = l.Match(start.to(9000), symbol.NewSetOfId(tInt, tPlus, tMinus, tInc))
	if !errors.Is(err, io.EOF) {
		t.Fatal("eof: wrong error", err)
	}
	//	const (
	//		expectEofStr     = `expected Int, Increment, Plus or Minus: parse error near ⟪38+23␠-␠␊␠19++⟫⏵<EOF>`
	//		expectEofStrPlus = `expected Int, Increment, Plus or Minus: parse error near:
	//38+23␠-␠␊␠19++<EOF>
	//--------------^
	//`
	//	)
	//	if fmt.Sprintf("%v", err) != expectEofStr {
	//		t.Errorf("err: <<<<%s>>>>", err)
	//	}
	//	if fmt.Sprintf("%+v", err) != expectEofStrPlus {
	//		t.Errorf("err+ wrong: <<<<%+v>>>>", err)
	//	}
}

func matchWS(st *State) (next *State, v any) {
	to, _ := st.TakeRunesFunc(unicode.IsSpace)
	if to.Offset() == st.Offset() {
		return nil, nil
	}
	return to, nil
}
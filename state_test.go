package lr0

import (
	"fmt"
	"io"
	"testing"
	"unicode/utf8"

	"github.com/vovan-ve/go-lr0-parser/internal/testutils"
)

const testStateSrcStr = "Lorem ipsum dolor sit amet, consectepture"
const testStateSrcLen = len(testStateSrcStr)

var testStateSource = []byte(testStateSrcStr)

func TestState_to(t *testing.T) {
	a := NewState(testStateSource)
	b := a.to(10)
	if a.at != 0 {
		t.Errorf("a.at changed to %v", a.at)
	}
	if b.at != 10 {
		t.Errorf("b.at is %v", b.at)
	}
	if string(b.source) != testStateSrcStr {
		t.Errorf("b.source is %q", string(b.source))
	}

	t.Run("panic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrNegativeOffset)
		a.to(-7)
	})

	c := b.to(15)
	if c.at != 15 {
		t.Errorf("c.at is %v", c.at)
	}

	d := b.to(testStateSrcLen + 7)
	if d.at != testStateSrcLen {
		t.Errorf("d.at is %v", d.at)
	}

	e := c.to(4)
	if e.at != 4 {
		t.Errorf("e.at is %v", e.at)
	}

	a2 := a.to(0)
	if a2 != a {
		t.Error("a2 is not a")
	}
}

func TestState_IsEOF(t *testing.T) {
	a := NewState(testStateSource)
	if a.IsEOF() {
		t.Error("a: it's not")
	}
	if a.to(15).IsEOF() {
		t.Error("b: it's not")
	}
	if !a.to(testStateSrcLen).IsEOF() {
		t.Error("c: it is")
	}
}

func TestState_Len(t *testing.T) {
	a := NewState(testStateSource)
	if a.Len() != testStateSrcLen {
		t.Errorf("a.Len() is %v", a.Len())
	}
	if b := a.to(15); b.Len() != testStateSrcLen {
		t.Errorf("b.Len() is %v", b.Len())
	}
	if c := a.to(testStateSrcLen); c.Len() != testStateSrcLen {
		t.Errorf("c.Len() is %v", c.Len())
	}
}

func TestState_Offset(t *testing.T) {
	a := NewState(testStateSource)
	if a.Offset() != 0 {
		t.Errorf("a.Offset() is %v", a.Offset())
	}
	if b := a.to(15); b.Offset() != 15 {
		t.Errorf("b.Offset() is %v", b.Offset())
	}
	if c := a.to(testStateSrcLen + 3); c.Offset() != testStateSrcLen {
		t.Errorf("c.Offset() is %v", c.Offset())
	}
}

func TestState_RestLen(t *testing.T) {
	a := NewState(testStateSource)
	if a.RestLen() != testStateSrcLen {
		t.Errorf("a.RestLen() is %v", a.RestLen())
	}
	if b := a.to(15); b.RestLen() != testStateSrcLen-15 {
		t.Errorf("b.RestLen() is %v", b.RestLen())
	}
	if c := a.to(testStateSrcLen); c.RestLen() != 0 {
		t.Errorf("c.RestLen() is %v", c.RestLen())
	}
}

func TestState_RestBytes(t *testing.T) {
	a := NewState(testStateSource)
	if string(a.RestBytes()) != testStateSrcStr {
		t.Errorf("a.RestLen() is %q", a.RestBytes())
	}
	if b := a.to(15); string(b.RestBytes()) != testStateSrcStr[15:] {
		t.Errorf("b.RestLen() is %q", b.RestBytes())
	}
	if c := a.to(testStateSrcLen); string(c.RestBytes()) != "" {
		t.Errorf("c.RestLen() is %q", c.RestBytes())
	}
}

func TestState_BytesToOffset(t *testing.T) {
	a := NewState(testStateSource)
	b := a.to(10)
	c := a.to(27)
	if s := a.BytesToOffset(15); string(s) != testStateSrcStr[:15] {
		t.Errorf("[1] is %q", s)
	}
	if s := b.BytesToOffset(31); string(s) != testStateSrcStr[10:31] {
		t.Errorf("[2] is %q", s)
	}
	if s := c.BytesToOffset(testStateSrcLen + 5); string(s) != testStateSrcStr[27:] {
		t.Errorf("[3] is %q", s)
	}

	t.Run("panic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrNegativeOffset)
		b.BytesToOffset(3)
	})
}

func TestState_Byte(t *testing.T) {
	a := NewState(testStateSource)
	if a.Byte() != 'L' {
		t.Errorf("a.Byte() is %q", a.Byte())
	}
	if b := a.to(9); b.Byte() != 'u' {
		t.Errorf("b.Byte() is %q", b.Byte())
	}
	if c := a.to(26); c.Byte() != ',' {
		t.Errorf("c.Byte() is %q", c.Byte())
	}
	if d := a.to(testStateSrcLen - 1); d.Byte() != 'e' {
		t.Errorf("d.Byte() is %q", d.Byte())
	}

	t.Run("panic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, io.EOF)
		a.to(testStateSrcLen).Byte()
	})
}

func TestState_Rune(t *testing.T) {
	a := NewState([]byte("#№±"))
	if r, n := a.Rune(); r != '#' || n != 1 {
		t.Errorf("a.Rune() is %q, %v", r, n)
	}
	if r, n := a.to(1).Rune(); r != '№' || n != 3 {
		t.Errorf("b.Rune() is %q, %v", r, n)
	}
	if r, n := a.to(4).Rune(); r != '±' || n != 2 {
		t.Errorf("c.Rune() is %q, %v", r, n)
	}

	// wrong octet offset
	if r, n := a.to(2).Rune(); r != utf8.RuneError || n != 1 {
		t.Errorf("a.to(2).Rune() is %q, %v", r, n)
	}
	if r, n := a.to(3).Rune(); r != utf8.RuneError || n != 1 {
		t.Errorf("a.to(3).Rune() is %q, %v", r, n)
	}

	t.Run("panic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, io.EOF)
		a.to(6).Rune()
	})
}

func TestState_FF(t *testing.T) {
	a := NewState(testStateSource)
	b := a.FF(10)
	if b.at != 10 {
		t.Errorf("b.at is %v", b.at)
	}

	t.Run("panic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrNegativeOffset)
		a.FF(-7)
	})

	c := b.FF(5)
	if c.at != 15 {
		t.Errorf("c.at is %v", c.at)
	}

	d := b.FF(9000)
	if d.at != testStateSrcLen {
		t.Errorf("d.at is %v", d.at)
	}

	e := c.FF(-11)
	if e.at != 4 {
		t.Errorf("e.at is %v", e.at)
	}

	a2 := a.FF(0)
	if a2 != a {
		t.Error("a2 is not a")
	}
}

func TestState_TakeByte(t *testing.T) {
	a := NewState(testStateSource)
	b, x := a.TakeByte()
	if x != 'L' {
		t.Errorf("b x is %q", x)
	}
	if b.Offset() != 1 {
		t.Errorf("b offset is %v", b.Offset())
	}

	c := a
	for rest := testStateSource[:]; len(rest) != 0; rest = rest[1:] {
		d, b := c.TakeByte()
		if b != rest[0] {
			t.Fatalf("got %q at %s; expect from %q", b, c, rest)
		}
		c = d
	}
	t.Run("eof", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, io.EOF)
		c.TakeByte()
	})
}

func TestState_TakeByteFunc(t *testing.T) {
	a := NewState(testStateSource)

	b, x := a.TakeByteFunc(byteIsNotSpace)
	if b == nil {
		t.Error("b: no match")
	}
	if x != 'L' {
		t.Errorf("b: x is %q", x)
	}
	if b.Offset() != 1 {
		t.Errorf("b: %v", b)
	}

	c, x := a.to(5).TakeByteFunc(byteIsNotSpace)
	if c != nil {
		t.Errorf("c: %v, %q", c, x)
	}
}

func TestState_TakeBytes(t *testing.T) {
	a := NewState(testStateSource)

	b, v := a.TakeBytes(4)
	if string(v) != "Lore" {
		t.Errorf("b v is %q", v)
	}
	if b.Offset() != 4 {
		t.Errorf("b offset is %v", b.Offset())
	}

	c, v := b.TakeBytes(10)
	if string(v) != "m ipsum do" {
		t.Errorf("c v is %q", v)
	}
	if c.Offset() != 14 {
		t.Errorf("c offset is %v", c.Offset())
	}

	z, v := a.FF(testStateSrcLen - 5).TakeBytes(7)
	if string(v) != "pture" {
		t.Errorf("z v is %q", v)
	}
	if !z.IsEOF() {
		t.Error("it is eof")
	}
}

func TestState_TakeBytesFunc(t *testing.T) {
	a := NewState(testStateSource)

	b, v := a.TakeBytesFunc(byteIsNotSpace)
	if string(v) != "Lorem" {
		t.Errorf("b v is %q", v)
	}
	if b.Offset() != 5 {
		t.Errorf("b offset is %v", b.Offset())
	}

	b, v = a.to(22).TakeBytesFunc(byteIsNotSpace)
	if string(v) != "amet," {
		t.Errorf("c v is %q", v)
	}
	if b.Offset() != 27 {
		t.Errorf("c offset is %v", b.Offset())
	}

	b, v = a.to(28).TakeBytesFunc(byteIsNotSpace)
	if string(v) != "consectepture" {
		t.Errorf("d v is %q", v)
	}
	if !b.IsEOF() {
		t.Error("d is eof")
	}
}

func TestState_TakeRune(t *testing.T) {
	a := NewState([]byte("#№±"))

	b, r := a.TakeRune()
	if r != '#' {
		t.Errorf("b r is %q", r)
	}
	if b.Offset() != 1 {
		t.Errorf("b offset is %v", b.Offset())
	}

	c, r := b.TakeRune()
	if r != '№' {
		t.Errorf("c r is %q", r)
	}
	if c.Offset() != 4 {
		t.Errorf("c offset is %v", c.Offset())
	}

	d, r := c.TakeRune()
	if r != '±' {
		t.Errorf("d r is %q", r)
	}
	if d.Offset() != 6 {
		t.Errorf("d offset is %v", d.Offset())
	}
	if !d.IsEOF() {
		t.Error("it is EOF")
	}
}

func TestState_TakeRunes(t *testing.T) {
	a := NewState([]byte("#№±°ᴁ¼⇒©😎"))

	b, v := a.TakeRunes(4)
	if string(v) != "#№±°" {
		t.Errorf("b v is %q", v)
	}
	if b.Offset() != 8 {
		t.Errorf("b offset is %v", b.Offset())
	}

	c, v := b.TakeRunes(3)
	if string(v) != "ᴁ¼⇒" {
		t.Errorf("c v is %q", v)
	}
	if c.Offset() != 16 {
		t.Errorf("c offset is %v", c.Offset())
	}

	d, v := c.TakeRunes(2)
	if string(v) != "©😎" {
		t.Errorf("d v is %q", v)
	}
	if !d.IsEOF() {
		t.Error("it is eof")
	}
}

func TestState_TakeRunesFunc(t *testing.T) {
	a := NewState([]byte("#№±°ᴁ¼ ⇒©😎"))

	b, v := a.TakeRunesFunc(runeIsNotSpace)
	if string(v) != "#№±°ᴁ¼" {
		t.Errorf("b v is %q", v)
	}
	if b.Offset() != 13 {
		t.Errorf("b offset is %v", b.Offset())
	}

	c, v := b.FF(1).TakeRunesFunc(runeIsNotSpace)
	if string(v) != "⇒©😎" {
		t.Errorf("b v is %q", v)
	}
	if c.Offset() != 23 {
		t.Errorf("c offset is %v", c.Offset())
	}
	if !c.IsEOF() {
		t.Error("c is EOF")
	}
}

func TestState_ExpectByte(t *testing.T) {
	a := NewState(testStateSource)

	b := a.ExpectByte('L')
	if b == nil {
		t.Errorf("got: %q", a.Byte())
	}

	c := b.ExpectByte('o', 'r', 'e', 'm')
	if c == nil {
		t.Errorf("got: %q", b.Byte())
	}

	if cBad := b.ExpectByte('o', 'r', 'X', 'Y'); cBad != nil {
		t.Errorf("got: %+v", cBad)
	}

	d := c.ExpectByte([]byte(" ipsum dolor")...)
	if d == nil {
		t.Errorf("got: %q", c.Byte())
	}

	if bad := a.ExpectByte('-'); bad != nil {
		t.Errorf("bad got %+v", bad)
	}
}

func TestState_ExpectByteOk(t *testing.T) {
	a := NewState(testStateSource)

	b, ok := a.ExpectByteOk('L', 'o', 'r', 'e', 'm')
	if !ok {
		t.Error("it's OK")
	}
	if b.Offset() != 5 {
		t.Errorf("b offset is %v", b.Offset())
	}

	c, ok := b.ExpectByteOk(' ', 'i', 'p', 'X', 'Y')
	if ok {
		t.Error("c isn't OK")
	}
	if c.Offset() != 8 {
		t.Errorf("c offset is %v", c.Offset())
	}

	d, ok := b.ExpectByteOk('X', 'Y')
	if ok {
		t.Error("d isn't OK")
	}
	if d != b {
		t.Error("d is not b")
	}
}

func TestState_Fmt(t *testing.T) {
	a := NewState([]byte(testStateSrcStr + ". " + testStateSrcStr))
	aStr := "⏵⟪Lorem␠ipsum␠dolor␠sit␠amet,␠co⟫"
	if a.String() != aStr {
		t.Errorf("a str got %q", a)
	}
	if fmt.Sprintf("%q", a) != aStr {
		t.Errorf("a fmt got %q", a)
	}

	aStr = "" +
		"Lorem␠ipsum␠dolor␠sit␠amet,␠co\n" +
		"^\n"
	if fmt.Sprintf("%+v", a) != aStr {
		t.Errorf("a fmt+ got<\n%+v>", a)
	}

	b := a.to(35)
	bStr := "⟪␠ipsum␠dolor␠sit␠amet,␠consect⟫⏵⟪epture.␠Lorem␠ipsum␠dolor␠sit␠⟫"
	if b.String() != bStr {
		t.Errorf("b str got %q", b)
	}
	if fmt.Sprintf("%q", b) != bStr {
		t.Errorf("b fmt got %q", b)
	}

	bStr = "" +
		"␠ipsum␠dolor␠sit␠amet,␠consectepture.␠Lorem␠ipsum␠dolor␠sit␠\n" +
		"------------------------------^\n"
	if fmt.Sprintf("%+v", b) != bStr {
		t.Errorf("b fmt+ got<\n%+v>", b)
	}

	b = a.to(9000)
	bStr = "⟪␠dolor␠sit␠amet,␠consectepture⟫⏵<EOF>"
	if b.String() != bStr {
		t.Errorf("b2 got %q", b)
	}
	if fmt.Sprintf("%q", b) != bStr {
		t.Errorf("b2 fmt got %q", b)
	}

	bStr = "" +
		"␠dolor␠sit␠amet,␠consectepture<EOF>\n" +
		"------------------------------^\n"
	if fmt.Sprintf("%+v", b) != bStr {
		t.Errorf("b2 fmt+ got<\n%+v>", b)
	}
}

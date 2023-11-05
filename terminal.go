package lr0

import (
	"reflect"

	"github.com/pkg/errors"
)

const (
	tWhitespace Id = -iota - 1
)

var (
	metaWS = term{id: tWhitespace, name: "whitespace"}
)

// TerminalFactory is a helper API to define a Terminal
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
//
//	NewWhitespace().FuncRune(unicode.IsSpace)
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
//
//	NewTerm(tCRLF, "CRLF").Byte('\r', '\n')
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

// FuncByte lets use alternative (more native and probably predefined in packages)
// match and calc functions instead of writing MatchFunc for every trivial case.
//
// `ok` checks if the subsequent byte is valid for this Terminal. The Terminal
// will match input bytes until `ok` returns `false` or EOF happens. If at least
// one byte matched, MatchFunc succeeds.
//
// `calc` is optional `func(b []byte) V` or `func(b []byte) (V, error)` to define
// how to evaluate the value of this Terminal. If `calc` is omitted or `nil`, the
// `[]byte` itself will be the value of this Terminal.
//
//	NewTerm(tInt, "int").FuncByte(isDigit, bytesToInt)
//	NewWhitespace().FuncByte(func(b byte) bool { return b == ' ' || b == '\t' })
//	...
//	func isDigit(b byte) bool              { return b >= '0' && b <= '9' }
//	func bytesToInt(b []byte) (int, error) { return strconv.Atoi(string(b)) }
func (t *TerminalFactory) FuncByte(ok func(byte) bool, calc ...any) Terminal {
	return &termCallback{
		term: t.term,
		fn:   newMatchFunc((*State).TakeBytesFunc, ok, calc...),
	}
}

// FuncRune lets use alternative (more native and probably predefined in packages)
// match and calc functions instead of writing MatchFunc for every trivial case.
//
// `ok` checks if the subsequent rune is valid for this Terminal. The Terminal
// will match input rune until `ok` returns `false` or EOF happens. If at least
// one rune matched, MatchFunc succeeds.
//
// `calc` is optional `func(b []rune) V` or `func(b []rune) (V, error)` to define
// how to evaluate the value of this Terminal. If `calc` is omitted or `nil`, the
// `[]rune` itself will be the value of this Terminal.
//
//	NewWhitespace().FuncRune(unicode.IsSpace)
func (t *TerminalFactory) FuncRune(ok func(rune) bool, calc ...any) Terminal {
	return &termCallback{
		term: t.term,
		fn:   newMatchFunc((*State).TakeRunesFunc, ok, calc...),
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

// newMatchFunc is generic wrapper to craft a MatchFunc from `ok` and optional
// `calc`
//
// `take` is a method of State
//
// `ok` returns whether the given T character from State is acceptable
//
// `calc` is optional `func([]T)V | func([]T)(V,error)` to evaluate value of a
// Terminal from the matched `[]T`. It can be nil or omitted to let the
// MatchFunc to return `[]T`.
func newMatchFunc[T any](
	take func(st *State, ok func(T) bool) (next *State, value []T),
	ok func(T) bool,
	calc ...any,
) MatchFunc {
	if len(calc) == 0 {
		return func(st *State) (next *State, value any) {
			next, value = take(st, ok)
			if next.Offset() == st.Offset() {
				return nil, nil
			}
			return
		}
	}

	valFunc := newValueFunc[T](calc[0])
	return func(st *State) (next *State, value any) {
		next, b := take(st, ok)
		if next.Offset() == st.Offset() {
			return nil, nil
		}
		value = valFunc(b)
		return
	}
}

// newValueFunc wraps the given `func([]T)V | func([]T)(V,error)` to `func([]T) any`
func newValueFunc[T any](fn any) func(v []T) any {
	funcV := reflect.ValueOf(fn)
	if funcV.Kind() != reflect.Func {
		panic(errors.Wrapf(ErrDefine, "fn contains not a func: %s", funcV.Kind()))
	}
	if funcV.IsNil() {
		panic(errors.Wrap(ErrDefine, "fn func is nil"))
	}

	funcT := funcV.Type()
	if funcT.NumIn() != 1 {
		panic(errors.Wrapf(ErrDefine, "fn arguments count is %d when wanted 1", funcT.NumIn()))
	}
	if funcT.IsVariadic() {
		panic(errors.Wrap(ErrDefine, "fn func is variadic"))
	}

	sliceOfTTyp := reflect.TypeOf([]T(nil))
	in1 := funcT.In(0)
	if in1.Kind() != reflect.Slice || !sliceOfTTyp.AssignableTo(in1) {
		panic(errors.Wrapf(ErrDefine, "fn argument of type `%s` cannot be assigned with value of type %s", in1, sliceOfTTyp))
	}

	switch funcT.NumOut() {
	case 1:
		return func(v []T) any {
			res := funcV.Call([]reflect.Value{reflect.ValueOf(v)})
			return res[0].Interface()
		}
	case 2:
		if t1 := funcT.Out(1); t1.Kind() != reflect.Interface || !t1.Implements(typeOfError) || !typeOfError.AssignableTo(t1) {
			panic(errors.Wrapf(ErrDefine, "fn func 2nd result must be `error`, given %v", t1))
		}
		return func(v []T) any {
			res := funcV.Call([]reflect.Value{reflect.ValueOf(v)})
			v0 := res[0].Interface()
			v1 := res[1].Interface()
			if v1 == nil {
				return v0
			}
			return v1.(error)
		}
	default:
		panic(errors.Wrapf(
			ErrDefine,
			"unexpected calc type `%T`, expected `func(%s)V` or `func(%s)(V, error)`",
			fn,
			sliceOfTTyp,
			sliceOfTTyp,
		))
	}
}

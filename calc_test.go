package lr0

import (
	"io"
	"strings"
	"testing"

	"github.com/vovan-ve/go-lr0-parser/internal/testutils"
)

func TestNewCalcFunc(t *testing.T) {
	t.Run("strings.Repeat", func(t *testing.T) {
		h1 := newCalcFunc(strings.Repeat, 2)
		v, err := h1([]any{"ab", 3})
		if err != nil {
			t.Errorf("h1 err is %#v", err)
		}
		if s, ok := v.(string); !ok {
			t.Errorf("h1 v is %#v", v)
		} else if s != "ababab" {
			t.Errorf("h1 v string is %q", s)
		}
	})
	t.Run("someTestFunc", func(t *testing.T) {
		type someTestType struct {
			i int
			s string
			b byte
		}

		someTestFunc := func(i int, s string, b byte) (*someTestType, error) {
			if i < 0 {
				return nil, io.EOF
			}
			return &someTestType{i: i, s: s, b: b}, nil
		}

		h2 := newCalcFunc(someTestFunc, 3)
		v, err := h2([]any{-5, "", byte(7)})
		if err != io.EOF {
			t.Errorf("h2.1 err is %#v", err)
		}
		v, err = h2([]any{5, "foo", byte(19)})
		if err != nil {
			t.Errorf("h2.2 err is %#v", err)
		}
		expect := someTestType{i: 5, s: "foo", b: 19}
		if s, ok := v.(*someTestType); !ok {
			t.Errorf("h2.2 v is %#v", v)
		} else if s == nil || *s != expect {
			t.Errorf("h2.2 v data is %#v", s)
		}
	})

	t.Run("panic: not func", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		newCalcFunc(42, 0)
	})
	t.Run("panic: null func", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		var fn func(any)
		newCalcFunc(fn, 1)
	})
	t.Run("panic: args count", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		newCalcFunc(func(any, any) {}, 4)
	})
	t.Run("panic: args count+variadic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		newCalcFunc(func(any, any, ...any) {}, 4)
	})
	t.Run("panic: variadic", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		newCalcFunc(func(any, any, ...any) {}, 3)
	})
	t.Run("panic: 2nd result not `error`, got any", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		newCalcFunc(func(any) (any, any) { return nil, nil }, 1)
	})
	t.Run("panic: 2nd result not `error`, got moreThenError", func(t *testing.T) {
		type moreThenError interface {
			error
			Foo() int
		}
		defer testutils.ExpectPanicError(t, ErrDefine)
		newCalcFunc(func(any) (any, moreThenError) { return nil, nil }, 1)
	})
	t.Run("panic: results count", func(t *testing.T) {
		defer testutils.ExpectPanicError(t, ErrDefine)
		newCalcFunc(func(any) {}, 1)
	})
}

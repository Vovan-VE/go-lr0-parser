package grammar

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
)

type calcFunc func([]any) (any, error)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

func prepareHandler(fn any, expectArgsCount int) calcFunc {
	if fn == nil && expectArgsCount == 1 {
		return bubble
	}

	funcV := reflect.ValueOf(fn)
	if funcV.Kind() != reflect.Func {
		panic(errors.Wrapf(symbol.ErrDefine, "fn contains not a func: %s", funcV.Kind()))
	}
	if funcV.IsNil() {
		panic(errors.Wrap(symbol.ErrDefine, "fn func is nil"))
	}

	funcT := funcV.Type()
	if funcT.NumIn() != expectArgsCount {
		panic(errors.Wrapf(symbol.ErrDefine, "fn arguments count is %d when wanted %d", funcT.NumIn(), expectArgsCount))
	}
	if funcT.IsVariadic() {
		panic(errors.Wrap(symbol.ErrDefine, "fn func is variadic"))
	}

	switch funcT.NumOut() {
	case 1:
		return func(v []any) (any, error) {
			res := funcV.Call(prepareArgs(v))
			return res[0].Interface(), nil
		}
	case 2:
		if t1 := funcT.Out(1); t1.Kind() != reflect.Interface || !t1.Implements(typeOfError) || !typeOfError.AssignableTo(t1) {
			panic(errors.Wrapf(symbol.ErrDefine, "fn func 2nd result must be `error`, given %v", t1))
		}
		return func(v []any) (any, error) {
			res := funcV.Call(prepareArgs(v))
			v0 := res[0].Interface()
			v1 := res[1].Interface()
			if v1 == nil {
				return v0, nil
			}
			return v0, v1.(error)
		}
	default:
		panic(errors.Wrapf(symbol.ErrDefine, "fn results count is %d when wanted 1 (value any) or 2 (value any, err error)", funcT.NumOut()))
	}
}

func prepareArgs(vs []any) []reflect.Value {
	res := make([]reflect.Value, 0, len(vs))
	for _, v := range vs {
		res = append(res, reflect.ValueOf(v))
	}
	return res
}

func bubble(v []any) (any, error) { return v[0], nil }

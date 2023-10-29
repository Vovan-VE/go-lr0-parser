package testutils

import (
	"testing"

	"github.com/pkg/errors"
)

func ExpectPanicError(t *testing.T, err error, fn ...func(*testing.T, error)) {
	e := recover()
	gotErr, ok := e.(error)
	if !ok {
		t.Fatalf("expected `error`, got %#v", e)
	}
	if err != nil && !errors.Is(gotErr, err) {
		t.Fatalf("unexpected %+v", e)
	}
	for _, f := range fn {
		f(t, gotErr)
	}
}

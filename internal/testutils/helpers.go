package testutils

import (
	"testing"

	"github.com/pkg/errors"
)

func ExpectPanicError(t *testing.T, err error) {
	e := recover()
	gotErr, ok := e.(error)
	if !ok {
		t.Fatalf("expected `error`, got %#v", e)
	}
	if !errors.Is(gotErr, err) {
		t.Fatalf("unexpected %+v", e)
	}
}

package lexer

import (
	"testing"

	"github.com/pkg/errors"
)

func TestParseError(t *testing.T) {
	err := WithSource(
		NewParseError("foo bar"),
		NewState(append(source, source...)).to(35),
	)
	if !errors.Is(err, ErrParse) {
		t.Errorf("unexpected false: %v", err)
	}
}

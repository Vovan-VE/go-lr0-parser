package lr0

import (
	"testing"

	"github.com/pkg/errors"
)

func TestParseError(t *testing.T) {
	err := WithSource(
		NewParseError("foo bar"),
		NewState(append(testStateSource, testStateSource...)).to(35),
	)
	if !errors.Is(err, ErrParse) {
		t.Errorf("unexpected false: %v", err)
	}
}

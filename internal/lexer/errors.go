package lexer

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
)

var (
	// ErrNegativeOffset will be raised by panic if some operation with State
	// cause negative offset
	ErrNegativeOffset = errors.New("negative position")
	// ErrParse is base error for run-time errors about parsing.
	//
	//	errors.Is(err, lexer.ErrParse)
	//
	//	errors.Wrap(lexer.ErrParse, "unexpected thing found")
	ErrParse = errors.New("parse error")
)

// NewParseError creates new Error by wrapping ErrParse
func NewParseError(msg string) error {
	return &parseError{
		error: errors.Wrap(ErrParse, msg),
	}
}

//// NewParseErrorf creates new Error by wrapping ErrParse
//func NewParseErrorf(format string, args ...any) error {
//	return &parseError{
//		error: errors.Wrapf(ErrParse, format, args...),
//	}
//}

// WithSource wraps the given error to append State info to error message
func WithSource(err error, s *State) error {
	return &withSource{
		error: err,
		src:   s,
	}
}

func ExpectationError(expected symbol.ReadonlySet, terminals []Terminal) error {
	s := "expected "
	i, last := 0, expected.Count()-1
	for _, t := range terminals {
		if !expected.Has(t.Id()) {
			continue
		}
		if i > 0 {
			if i < last {
				s += ", "
			} else {
				s += " or "
			}
		}
		i++
		s += symbol.Dump(t)
	}
	return NewParseError(s)
}

type parseError struct {
	error
}

func (p *parseError) Unwrap() error { return p.error }

func (p *parseError) Format(s fmt.State, verb rune) {
	if x, ok := p.error.(fmt.Formatter); ok {
		x.Format(s, verb)
		return
	}

	switch verb {
	case 'v', 's', 'q':
		io.WriteString(s, p.Error())
	}
}

type withSource struct {
	error
	src *State
}

func (w *withSource) Error() string {
	return w.error.Error() + " near " + fmt.Sprintf("%s", w.src)
}

func (w *withSource) Unwrap() error {
	return w.error
}

func (w *withSource) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, w.error.Error())
			io.WriteString(s, " near:\n")
			w.src.Format(s, verb)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

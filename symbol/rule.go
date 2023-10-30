package symbol

import (
	"github.com/pkg/errors"
)

// Rule is a common rule data
type Rule interface {
	// Subject of the rule
	Subject() Id
	// HasEOF tells whether EOF must be found in the end of input
	HasEOF() bool
}

// NewRule creates new Rule with only Subject Id
func NewRule(subject Id) Rule {
	return rule{subject: subject}
}

// WithEOF creates new copy of our local Rule implementation with HasEOF flag
// set to true
//
// Will return input rule as is if it already has HasEOF flag set to true.
//
// Will panic with wrapped ErrDefine if the given Rule implementation is not
// from this package.
func WithEOF(r Rule) Rule {
	if r.HasEOF() {
		return r
	}
	orig, ok := r.(rule)
	if !ok {
		panic(errors.Wrap(ErrDefine, "not my Rule implementation"))
	}
	result := orig
	result.eof = true
	return result
}

type rule struct {
	subject Id
	eof     bool
}

func (r rule) Subject() Id  { return r.subject }
func (r rule) HasEOF() bool { return r.eof }

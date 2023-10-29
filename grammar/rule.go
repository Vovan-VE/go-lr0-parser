package grammar

import (
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

// Rule is one of possible definition for a non-Terminal
type Rule interface {
	symbol.Rule
	// Definition of what it consists of
	Definition() []symbol.Id
}

// NewRule creates Rule from symbol.Rule
func NewRule(r symbol.Rule, definition []symbol.Id) Rule {
	return &rule{
		Rule:       r,
		definition: definition,
	}
}

// NewRuleId creates Rule without Tag and EOF flag
func NewRuleId(subject symbol.Id, definition []symbol.Id) Rule {
	return &rule{
		Rule:       symbol.NewRule(subject),
		definition: definition,
	}
}

type rule struct {
	symbol.Rule
	definition []symbol.Id
}

func (r *rule) Definition() []symbol.Id { return r.definition }

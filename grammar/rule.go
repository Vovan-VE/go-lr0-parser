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

// NewRule creates Rule without Tag and EOF flag
func NewRule(subject symbol.Id, definition []symbol.Id) Rule {
	return &rule{
		Rule:       symbol.NewRule(subject),
		definition: definition,
	}
}

// NewRuleTag creates Rule with Tag, but without EOF flag
func NewRuleTag(subject symbol.Id, tag symbol.Tag, definition []symbol.Id) Rule {
	return &rule{
		Rule:       symbol.NewRuleTag(subject, tag),
		definition: definition,
	}
}

// NewRuleMain creates Rule with EOF flag, but without Tag. Since a grammar must
// have exactly one main rule, a Tag is useless in main rule
func NewRuleMain(subject symbol.Id, definition []symbol.Id) Rule {
	return &rule{
		Rule:       symbol.WithEOF(symbol.NewRule(subject)),
		definition: definition,
	}
}

type rule struct {
	symbol.Rule
	definition []symbol.Id
}

func (r *rule) Definition() []symbol.Id { return r.definition }

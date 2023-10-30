package grammar

import (
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

// Rule is one of possible definition for a non-Terminal
type Rule interface {
	symbol.Rule
	// Definition of what it consists of
	Definition() []symbol.Id
}
type RuleDefinition interface {
	Rule
	CalcHandler() any
	// TODO: Hide(index... int)
}
type RuleImplementation interface {
	Rule
	Value([]any) (any, error)
	IsHidden(index int) bool
}

// TODO: chainable RuleDefinition creation
//   NewRule(nSum).Main().Is(nSum, tPlus, nVal).Do(calcSum)

// NewRule creates Rule without Tag and EOF flag
func NewRule(subject symbol.Id, definition []symbol.Id, calcHandler any) RuleDefinition {
	return &ruleDefinition{
		Rule: &rule{
			Rule:       symbol.NewRule(subject),
			definition: definition,
		},
		calcHandler: calcHandler,
	}
}

// NewRuleMain creates Rule with EOF flag, but without Tag. Since a grammar must
// have exactly one main rule, a Tag is useless in main rule
func NewRuleMain(subject symbol.Id, definition []symbol.Id, calcHandler any) RuleDefinition {
	return &ruleDefinition{
		Rule: &rule{
			Rule:       symbol.WithEOF(symbol.NewRule(subject)),
			definition: definition,
		},
		calcHandler: calcHandler,
	}
}

func ToImplementation(rd RuleDefinition, l lexer.HiddenRegistry) RuleImplementation {
	hidden := whichHidden(l, rd)
	return &ruleImplementation{
		Rule:   rd,
		calc:   prepareHandler(rd.CalcHandler(), len(rd.Definition())-len(hidden)),
		hidden: hidden,
	}
}

type rule struct {
	symbol.Rule
	definition []symbol.Id
}

func (r *rule) Definition() []symbol.Id { return r.definition }

type ruleDefinition struct {
	Rule
	calcHandler any
}

func (r *ruleDefinition) CalcHandler() any { return r.calcHandler }

type ruleImplementation struct {
	Rule
	calc   CalcFunc
	hidden map[int]struct{}
}

func (r *ruleImplementation) Value(v []any) (any, error) { return r.calc(v) }

func (r *ruleImplementation) IsHidden(index int) bool {
	_, ok := r.hidden[index]
	return ok
}

// How many output values the rule will provide. A hidden terminals designed to
// skip unnecessary useless values.
func whichHidden(l lexer.HiddenRegistry, r Rule) map[int]struct{} {
	hidden := make(map[int]struct{})
	for i, id := range r.Definition() {
		if l.IsHidden(id) {
			hidden[i] = struct{}{}
		}
	}
	return hidden
}

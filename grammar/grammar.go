package grammar

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

// Grammar defines full grammar how to parse an input stream
type Grammar interface {
	lexer.Lexer
	RuleImpl(index int) RuleImplementation
	// MainRule returns main rule, the one with EOF flag
	MainRule() RuleImplementation
	// RulesFor returns set of rules for the given subject
	RulesFor(id symbol.Id) []RuleImplementation
}

// New creates new Grammar
//
// Violation of the following statements will cause panic symbol.ErrDefine since
// it is definition mistake.
//
// - Every Terminal must be used in Rules
//
// - Every symbol.Id in every rule definition must exist either in Terminals or
// in Rules Subject
//
// - Exactly one Rule must have EOF flag - this is Main Rule
func New(terminals []lexer.Terminal, rules []RuleDefinition) Grammar {
	var (
		l           = lexer.New(terminals...)
		rulesImpl   = make([]RuleImplementation, 0, len(rules))
		mainI       = -1
		si          = make(map[symbol.Id][]int)
		furtherNTAt = make(map[symbol.Id]string)
		usedT       = make(map[symbol.Id]struct{})
	)

	for ri, r := range rules {
		if r.HasEOF() {
			if mainI >= 0 {
				panic(errors.Wrapf(symbol.ErrDefine, "another rule [%d] has EOF flag too, previous was [%d]", ri, mainI))
			}
			mainI = ri
		}

		subjId := r.Subject()
		if subjId == symbol.InvalidId {
			panic(errors.Wrapf(symbol.ErrDefine, "rules[%d] subject id is zero", ri))
		}
		if l.IsTerminal(subjId) {
			panic(errors.Wrapf(symbol.ErrDefine, "rules[%d] subject is Terminal", ri))
		}
		si[subjId] = append(si[subjId], ri)
		// now this non-terminal is defined
		delete(furtherNTAt, subjId)

		for i, id := range r.Definition() {
			// defined Terminal - ok
			if l.IsTerminal(id) {
				usedT[id] = struct{}{}
				continue
			}
			// already defined non-terminal - ok
			if _, ok := si[id]; ok {
				continue
			}
			// already know where it was seen first time - ok
			if _, ok := furtherNTAt[id]; ok {
				continue
			}
			// seeing this non-terminal first time
			// TODO: Rule to string
			furtherNTAt[id] = fmt.Sprintf("#%d in rules[%d] definitions[%d]", id, ri, i)
		}

		rulesImpl = append(rulesImpl, ToImplementation(r, l))
	}

	if len(furtherNTAt) != 0 {
		msg := "undefined non-terminals without rules:\n"
		for _, at := range furtherNTAt {
			msg += "- " + at + "\n"
		}
		panic(errors.Wrap(symbol.ErrDefine, msg))
	}
	if mainI == -1 {
		panic(errors.Wrap(symbol.ErrDefine, "no main rule with EOF flag"))
	}
	if len(usedT) != len(terminals) {
		msg := "following Terminals are not used in any Rule:\n"
		for _, t := range terminals {
			if _, ok := usedT[t.Id()]; ok {
				continue
			}
			msg += "- " + symbol.Dump(t) + "\n"
		}
		panic(errors.Wrap(symbol.ErrDefine, msg))
	}

	return &grammar{
		Lexer:           l,
		rules:           rulesImpl,
		mainIndex:       mainI,
		subjectsIndices: si,
	}
}

type grammar struct {
	lexer.Lexer
	rules           []RuleImplementation
	mainIndex       int
	subjectsIndices map[symbol.Id][]int
}

func (g *grammar) RuleImpl(index int) RuleImplementation { return g.rules[index] }

func (g *grammar) MainRule() RuleImplementation {
	return g.rules[g.mainIndex]
}

func (g *grammar) RulesFor(id symbol.Id) []RuleImplementation {
	indices, ok := g.subjectsIndices[id]
	if !ok {
		if g.IsTerminal(id) {
			panic(errors.Wrapf(symbol.ErrDefine, "no rule - #%d is Terminal", id))
		}
		panic(errors.Wrapf(symbol.ErrDefine, "no rule for #%d", id))
	}
	ret := make([]RuleImplementation, 0, len(indices))
	for _, idx := range indices {
		ret = append(ret, g.rules[idx])
	}
	return ret
}

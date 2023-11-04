package lr0

import (
	"fmt"

	"github.com/pkg/errors"
)

//// Grammar defines full grammar how to parse an input stream
//type Grammar interface {
//	Lexer
//	RulesCount() int
//	Rule(index int) Rule
//	// MainRule returns main rule, the one with EOF flag
//	MainRule() Rule
//	// RulesFor returns set of rules for the given subject
//	RulesFor(id Id) []Rule
//}

// newGrammar creates new Grammar
//
// Violation of the following statements will cause panic ErrDefine since
// it is definition mistake.
//
// - All Terminal must be used in Rules. No unused Terminals.
//
// - Every Id in every rule definition must exist either in Terminals or
// in Rules Subject
//
// - Exactly one Rule must have EOF flag - this is Main Rule
func newGrammar(terminals []Terminal, nonTerminals []NonTerminalDefinition) *grammar {
	var (
		l           = newLexer(terminals...)
		nonTerm     = make(map[Id]Symbol)
		mainS       Symbol
		ruleIndex   int
		si          = make(map[Id][]int)
		furtherNTAt = make(map[Id]string)
		usedT       = make(map[Id]struct{})
	)

	gr := &grammar{
		lexer:           l,
		nonTerm:         nonTerm,
		rules:           make([]Rule, 0),
		mainIndex:       -1,
		subjectsIndices: si,
	}

	for _, ntDef := range nonTerminals {
		subjId := ntDef.Id()
		nonTerm[subjId] = ntDef
		delete(furtherNTAt, subjId)

		if subjId == InvalidId {
			panic(errors.Wrapf(ErrDefine, "Non-Terminal %s id is zero", dumpSymbol(ntDef)))
		}
		if l.IsTerminal(subjId) {
			panic(errors.Wrapf(ErrDefine, "Non-Terminal %s is Terminal", dumpSymbol(ntDef)))
		}

		for ri, r := range ntDef.GetRules(gr) {
			if r.HasEOF() {
				if mainS != nil {
					panic(errors.Wrapf(ErrDefine, "another rule %s has Main flag too, previous was %s", dumpSymbol(ntDef), dumpSymbol(mainS)))
				}
				mainS = ntDef
				gr.mainIndex = ruleIndex
			}

			si[subjId] = append(si[subjId], ruleIndex)
			// now this non-terminal is defined

			gr.rules = append(gr.rules, r)
			ruleIndex++

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
				furtherNTAt[id] = fmt.Sprintf("#%d in NT %s rules[%d] (%s) definitions[%d]", id, dumpSymbol(ntDef), ri, r, i)
			}
		}
	}

	if len(furtherNTAt) != 0 {
		msg := "undefined non-terminals without rules:\n"
		for _, at := range furtherNTAt {
			msg += "- " + at + "\n"
		}
		panic(errors.Wrap(ErrDefine, msg))
	}
	if gr.mainIndex == -1 {
		panic(errors.Wrap(ErrDefine, "no main rule with EOF flag"))
	}
	if len(usedT) != len(terminals) {
		msg := "following Terminals are not used in any Rule:\n"
		bad := false
		for _, t := range terminals {
			if t.Id() < 0 {
				continue
			}
			if _, ok := usedT[t.Id()]; ok {
				continue
			}
			msg += "- " + dumpSymbol(t) + "\n"
			bad = true
		}
		if bad {
			panic(errors.Wrap(ErrDefine, msg))
		}
	}

	return gr
}

type grammar struct {
	*lexer
	nonTerm         map[Id]Symbol
	rules           []Rule
	mainIndex       int
	subjectsIndices map[Id][]int
}

func (g *grammar) SymbolName(id Id) string {
	if n := g.lexer.SymbolName(id); n != "" {
		return n
	}
	if s, ok := g.nonTerm[id]; ok {
		return s.Name()
	}
	return ""
}

func (g *grammar) RulesCount() int     { return len(g.rules) }
func (g *grammar) Rule(index int) Rule { return g.rules[index] }

func (g *grammar) MainRule() Rule {
	return g.rules[g.mainIndex]
}

func (g *grammar) RulesFor(id Id) []Rule {
	indices, ok := g.subjectsIndices[id]
	if !ok {
		if g.IsTerminal(id) {
			panic(errors.Wrapf(ErrDefine, "no rule - %s is Terminal", dumpId(id, g)))
		}
		panic(errors.Wrapf(ErrDefine, "no rule for #%d", id))
	}
	ret := make([]Rule, 0, len(indices))
	for _, idx := range indices {
		ret = append(ret, g.rules[idx])
	}
	return ret
}

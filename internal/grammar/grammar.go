package grammar

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/lexer"
	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
)

// Grammar defines full grammar how to parse an input stream
type Grammar interface {
	lexer.Lexer
	RulesCount() int
	Rule(index int) Rule
	// MainRule returns main rule, the one with EOF flag
	MainRule() Rule
	// RulesFor returns set of rules for the given subject
	RulesFor(id symbol.Id) []Rule
}

// New creates new Grammar
//
// Violation of the following statements will cause panic symbol.ErrDefine since
// it is definition mistake.
//
// - All Terminal must be used in Rules. No unused Terminals.
//
// - Every symbol.Id in every rule definition must exist either in Terminals or
// in Rules Subject
//
// - Exactly one Rule must have EOF flag - this is Main Rule
func New(terminals []lexer.Terminal, nonTerminals []NonTerminalDefinition) Grammar {
	var (
		l           = lexer.New(terminals...)
		nonTerm     = make(map[symbol.Id]symbol.Symbol)
		mainS       symbol.Symbol
		ruleIndex   int
		si          = make(map[symbol.Id][]int)
		furtherNTAt = make(map[symbol.Id]string)
		usedT       = make(map[symbol.Id]struct{})
	)

	gr := &grammar{
		Lexer:           l,
		nonTerm:         nonTerm,
		rules:           make([]Rule, 0),
		mainIndex:       -1,
		subjectsIndices: si,
	}

	for _, ntDef := range nonTerminals {
		subjId := ntDef.Id()
		nonTerm[subjId] = ntDef
		delete(furtherNTAt, subjId)

		if subjId == symbol.InvalidId {
			panic(errors.Wrapf(symbol.ErrDefine, "Non-Terminal %s id is zero", symbol.Dump(ntDef)))
		}
		if l.IsTerminal(subjId) {
			panic(errors.Wrapf(symbol.ErrDefine, "Non-Terminal %s is Terminal", symbol.Dump(ntDef)))
		}

		for ri, r := range ntDef.GetRules(gr) {
			if r.HasEOF() {
				if mainS != nil {
					panic(errors.Wrapf(symbol.ErrDefine, "another rule %s has Main flag too, previous was %s", symbol.Dump(ntDef), symbol.Dump(mainS)))
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
				furtherNTAt[id] = fmt.Sprintf("#%d in NT %s rules[%d] (%s) definitions[%d]", id, symbol.Dump(ntDef), ri, r, i)
			}
		}
	}

	if len(furtherNTAt) != 0 {
		msg := "undefined non-terminals without rules:\n"
		for _, at := range furtherNTAt {
			msg += "- " + at + "\n"
		}
		panic(errors.Wrap(symbol.ErrDefine, msg))
	}
	if gr.mainIndex == -1 {
		panic(errors.Wrap(symbol.ErrDefine, "no main rule with EOF flag"))
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
			msg += "- " + symbol.Dump(t) + "\n"
			bad = true
		}
		if bad {
			panic(errors.Wrap(symbol.ErrDefine, msg))
		}
	}

	return gr
}

type grammar struct {
	lexer.Lexer
	nonTerm         map[symbol.Id]symbol.Symbol
	rules           []Rule
	mainIndex       int
	subjectsIndices map[symbol.Id][]int
}

func (g *grammar) SymbolName(id symbol.Id) string {
	if n := g.Lexer.SymbolName(id); n != "" {
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

func (g *grammar) RulesFor(id symbol.Id) []Rule {
	indices, ok := g.subjectsIndices[id]
	if !ok {
		if g.IsTerminal(id) {
			panic(errors.Wrapf(symbol.ErrDefine, "no rule - %s is Terminal", symbol.DumpId(id, g)))
		}
		panic(errors.Wrapf(symbol.ErrDefine, "no rule for #%d", id))
	}
	ret := make([]Rule, 0, len(indices))
	for _, idx := range indices {
		ret = append(ret, g.rules[idx])
	}
	return ret
}

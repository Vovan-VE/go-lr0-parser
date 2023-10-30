package table

import (
	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

// Row is a single row of a Table
//
// https://en.wikipedia.org/wiki/LR_parser#Table_construction
type Row interface {
	// AcceptEof returns true if this state accepts EOF
	AcceptEof() bool
	// Terminals returns all terminals possible in this state
	Terminals() []symbol.Id
	// TerminalAction returns next state index for the given terminal
	TerminalAction(id symbol.Id) (StateIndex, bool)
	// GotoAction returns next state index for the given non-terminal
	GotoAction(id symbol.Id) (StateIndex, bool)
	// ReduceRule returns a reduce rule if available or nil otherwise
	ReduceRule() grammar.RuleImplementation
	// IsReduceOnly returns true if this state can only be used for reduce
	IsReduceOnly() bool
}

func newRow() *row {
	return &row{
		terminals: make(stateActions),
		gotos:     make(stateActions),
	}
}

type stateActions map[symbol.Id]StateIndex

type row struct {
	acceptEof bool
	terminals stateActions
	gotos     stateActions
	// REFACT: this both can be merged together, add `terminals []Id` for `Terminals()`

	reduceRule grammar.RuleImplementation
}

func (r *row) AcceptEof() bool { return r.acceptEof }
func (r *row) SetAcceptEof()   { r.acceptEof = true }

func (r *row) ReduceRule() grammar.RuleImplementation     { return r.reduceRule }
func (r *row) SetReduceRule(v grammar.RuleImplementation) { r.reduceRule = v }

func (r *row) Terminals() []symbol.Id {
	ret := make([]symbol.Id, 0, len(r.terminals))
	for id := range r.terminals {
		ret = append(ret, id)
	}
	return ret
}

func (r *row) TerminalAction(id symbol.Id) (StateIndex, bool) {
	idx, ok := r.terminals[id]
	return idx, ok
}

func (r *row) SetTerminalAction(id symbol.Id, idx StateIndex) {
	// impossible to predict or check order of overlapping terminals here
	// example is plus `+` and increment `++`
	// a `+` can incorrectly match a part of increment `++` which is incorrect
	if v, ok := r.terminals[id]; ok && v != idx {
		panic(errors.New("already was set to different index"))
	}
	r.terminals[id] = idx
}

func (r *row) GotoAction(id symbol.Id) (StateIndex, bool) {
	idx, ok := r.gotos[id]
	return idx, ok
}

func (r *row) SetGoto(id symbol.Id, idx StateIndex) {
	if v, ok := r.gotos[id]; ok && v != idx {
		panic(errors.New("already was set to different index"))
	}
	r.gotos[id] = idx
}

func (r *row) IsReduceOnly() bool {
	return !r.acceptEof &&
		len(r.terminals) == 0 &&
		len(r.gotos) == 0 &&
		r.reduceRule != nil
}

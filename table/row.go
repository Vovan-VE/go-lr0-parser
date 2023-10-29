package table

import (
	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

type Row interface {
	AcceptEof() bool
	Terminals() []symbol.Id
	TerminalAction(id symbol.Id) (StateIndex, bool)
	GotoAction(id symbol.Id) (StateIndex, bool)
	ReduceRule() grammar.Rule
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
	acceptEof  bool
	terminals  stateActions
	gotos      stateActions
	reduceRule grammar.Rule
}

func (r *row) AcceptEof() bool { return r.acceptEof }
func (r *row) SetAcceptEof()   { r.acceptEof = true }

func (r *row) ReduceRule() grammar.Rule     { return r.reduceRule }
func (r *row) SetReduceRule(v grammar.Rule) { r.reduceRule = v }

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

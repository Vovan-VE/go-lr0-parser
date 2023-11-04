package lr0

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/internal/helpers"
)

//// Row is a single row of a Table
////
//// https://en.wikipedia.org/wiki/LR_parser#Table_construction
//type Row interface {
//	// AcceptEof returns true if this state accepts EOF
//	AcceptEof() bool
//	// TerminalsSet returns all terminals possible in this state
//	TerminalsSet() readonlyIdSet
//	// TerminalAction returns next state index for the given terminal
//	TerminalAction(id Id) (tableStateIndex, bool)
//	// GotoAction returns next state index for the given non-terminal
//	GotoAction(id Id) (tableStateIndex, bool)
//	// ReduceRule returns a reduce rule if available or nil otherwise
//	ReduceRule() Rule
//	// IsReduceOnly returns true if this state can only be used for reduce
//	IsReduceOnly() bool
//}

func newTableRow() *tableRow {
	return &tableRow{
		terminalsSet: newIdSet(),
		terminals:    make(stateActions),
		gotos:        make(stateActions),
	}
}

type tableRow struct {
	acceptEof    bool
	terminalsSet idSet
	terminals    stateActions
	gotos        stateActions

	reduceRule Rule
}

func (r *tableRow) AcceptEof() bool { return r.acceptEof }
func (r *tableRow) SetAcceptEof()   { r.acceptEof = true }

func (r *tableRow) ReduceRule() Rule     { return r.reduceRule }
func (r *tableRow) SetReduceRule(v Rule) { r.reduceRule = v }

func (r *tableRow) TerminalsSet() readonlyIdSet { return r.terminalsSet }

func (r *tableRow) TerminalAction(id Id) (tableStateIndex, bool) {
	idx, ok := r.terminals[id]
	return idx, ok
}

func (r *tableRow) SetTerminalAction(id Id, idx tableStateIndex) {
	// impossible to predict or check order of overlapping terminals here
	// example is plus `+` and increment `++`
	// a `+` can incorrectly match a part of increment `++` which is incorrect
	if v, ok := r.terminals[id]; ok && v != idx {
		panic(errors.Wrap(ErrInternal, "already was set to different index"))
	}
	r.terminalsSet.Add(id)
	r.terminals[id] = idx
}

func (r *tableRow) GotoAction(id Id) (tableStateIndex, bool) {
	idx, ok := r.gotos[id]
	return idx, ok
}

func (r *tableRow) SetGoto(id Id, idx tableStateIndex) {
	if v, ok := r.gotos[id]; ok && v != idx {
		panic(errors.Wrap(ErrInternal, "already was set to different index"))
	}
	r.gotos[id] = idx
}

func (r *tableRow) IsReduceOnly() bool {
	return !r.acceptEof &&
		len(r.terminals) == 0 &&
		len(r.gotos) == 0 &&
		r.reduceRule != nil
}

func (r *tableRow) dump(indent string, reg SymbolRegistry) string {
	res := indent + "EOF: "
	if r.acceptEof {
		res += "-"
	} else {
		res += "ACCEPT"
	}

	res += "\n" + indent + "terminals:"
	if len(r.terminals) != 0 {
		res += "\n" + r.terminals.dump(indent+"\t", reg)
	} else {
		res += " -\n"
	}

	res += indent + "goto:"
	if len(r.gotos) != 0 {
		res += "\n" + r.gotos.dump(indent+"\t", reg)
	} else {
		res += " -\n"
	}

	res += indent + "rule:"
	if r.reduceRule != nil {
		res += "\n" + indent + "\t" + r.ReduceRule().String() + "\n"
	} else {
		res += " -\n"
	}
	return res
}

type stateActions map[Id]tableStateIndex

func (s stateActions) dump(indent string, r SymbolRegistry) string {
	res := ""
	for _, p := range helpers.MapSortedInt(s) {
		res += indent + fmt.Sprintf("%s -> %v\n", dumpId(p.K, r), p.V)
	}
	return res
}

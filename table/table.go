package table

import (
	"fmt"

	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/internal/helpers"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

type StateIndex = int

// Table is states table controlling how parses state will behave
//
// https://en.wikipedia.org/wiki/LR_parser#Table_construction
type Table interface {
	//RowsCount() int

	Row(idx StateIndex) Row
}

// New creates new Table from the given Grammar
func New(g grammar.Grammar) Table {
	var (
		rows   []*row
		states []itemset
	)

	addStates := newAddStatesMap()
	addStates[0] = map[symbol.Id]itemset{
		symbol.InvalidId: newItemset([]item{newItem(g.MainRule())}, g),
	}
	for len(addStates) != 0 {
		nextStates := newAddStatesMap()
		// maps are sorted only for stable iterations order - it's better for
		// debug purpose, so table is stable for same input grammar
		for _, fromSets := range helpers.MapSortedInt(addStates) {
			fromSI := fromSets.K
		NewSets:
			for _, from := range helpers.MapSortedInt(fromSets.V) {
				fromId, fromState := from.K, from.V
				var fromIsT, fromIsNT bool
				if fromId != symbol.InvalidId {
					fromIsT = g.IsTerminal(fromId)
					fromIsNT = !fromIsT
				}
				for si, st := range states {
					if st.IsEqual(fromState) {
						thatRow := rows[fromSI]
						if fromIsT {
							thatRow.SetTerminalAction(fromId, si)
						} else if fromIsNT {
							thatRow.SetGoto(fromId, si)
						}
						continue NewSets
					}
				}

				newSI := len(states)
				newR := newRow()
				if fromState.HasFinalItem() {
					newR.SetAcceptEof()
				}
				states = append(states, fromState)
				rows = append(rows, newR)

				nextStates[newSI] = fromState.GetNextItemsets(g)

				fromRow := rows[fromSI]
				if fromIsT {
					fromRow.SetTerminalAction(fromId, newSI)
				} else if fromIsNT {
					fromRow.SetGoto(fromId, newSI)
				}
			}
		}
		addStates = nextStates
	}

	for si, st := range states {
		if r := st.ReduceRule(); r != nil {
			rows[si].SetReduceRule(r)
		}
	}

	return &table{rows: rows}
}

func newAddStatesMap() map[StateIndex]map[symbol.Id]itemset {
	return make(map[StateIndex]map[symbol.Id]itemset)
}

type table struct {
	rows []*row
}

//func (t *table) RowsCount() int { return len(t.rows) }

func (t *table) Row(idx StateIndex) Row { return t.rows[idx] }

func (t *table) dump() string {
	res := "====[ table ]====\n"
	for i, r := range t.rows {
		res += fmt.Sprintf("row %v ---------\n", i)
		res += r.dump("\t")
	}
	res += "=================\n"
	return res
}

package table

import (
	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

type StateIndex = int

type Table interface {
	Row(idx StateIndex) Row
}

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
		for fromSI, fromSets := range addStates {
		NewSets:
			for fromId, fromState := range fromSets {
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

func (t *table) Row(idx StateIndex) Row { return t.rows[idx] }

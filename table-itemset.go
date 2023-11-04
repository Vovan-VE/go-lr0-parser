package lr0

import (
	"github.com/pkg/errors"
)

func newTableItemset(items []tableItem, g *grammar) tableItemset {
	allItems := expandAllPossibleTableItems(items, g)
	validateTableItemsetDeterministic(allItems, g)
	return tableItemset{items: allItems}
}

// expandAllPossibleTableItems returns possible items by expansion all the given items
// by Grammar
func expandAllPossibleTableItems(items []tableItem, g *grammar) []tableItem {
	var final []tableItem
	knownNext := newIdSet()

	newItems := items
	for len(newItems) != 0 {
		var nextNew []Id
	NewItems:
		for _, newIt := range newItems {
			for _, have := range final {
				if newIt == have {
					continue NewItems
				}
			}

			final = append(final, newIt)

			nextId := newIt.Expected()
			if nextId == InvalidId || g.IsTerminal(nextId) {
				continue
			}
			if knownNext.Has(nextId) {
				continue
			}
			nextNew = append(nextNew, nextId)
		}

		knownNext.Add(nextNew...)

		newItems = nil
		for _, id := range nextNew {
			for _, r := range g.RulesFor(id) {
				newItems = append(newItems, newTableItem(r))
			}
		}
	}

	return final
}

// validateTableItemsetDeterministic checks for bad state in this tableItemset.
// A problem is reported by panic since it's grammar definition problem
//
// - ErrConflictReduceReduce
//
// - ErrConflictShiftReduce
//
// https://en.wikipedia.org/wiki/LR_parser#Conflicts_in_the_constructed_tables
func validateTableItemsetDeterministic(items []tableItem, g *grammar) {
	var finite, terminals, nonTerminals []tableItem
	for _, it := range items {
		if !it.HasFurther() {
			finite = append(finite, it)
			continue
		}
		if g.IsTerminal(it.Expected()) {
			terminals = append(terminals, it)
			continue
		}
		nonTerminals = append(nonTerminals, it)
	}

	// Check for Shift-Reduce conflicts
	if len(finite) != 0 && len(terminals) != 0 {
		restTerminals := g.GetTerminalIdsSet()
		for _, it := range terminals {
			restTerminals.Remove(it.Expected())
		}
		if restTerminals.Count() == 0 {
			panic(errors.WithStack(ErrConflictShiftReduce))
		}
	}

	// Check for Reduce-Reduce conflicts
	if len(finite) > 1 {
		panic(errors.WithStack(ErrConflictReduceReduce))
	}
}

type tableItemset struct {
	items []tableItem
}

// HasFinalItem checks if this set has tableItem where the next token must be EOF
func (s tableItemset) HasFinalItem() bool {
	for _, it := range s.items {
		if it.HasEOF() && !it.HasFurther() {
			return true
		}
	}
	return false
}

// ReduceRule returns reduction rule of this set if any, nil otherwise
func (s tableItemset) ReduceRule() Rule {
	for _, it := range s.items {
		if !it.HasFurther() {
			// this is the only due to Reduce-Reduce conflicts validation
			return it.Rule
		}
	}
	return nil
}

// IsEqual checks equality with another tableItemset
func (s tableItemset) IsEqual(to tableItemset) bool {
	if len(s.items) != len(to.items) {
		return false
	}
	my := make(map[tableItem]struct{})
	for _, it := range s.items {
		my[it] = struct{}{}
	}
	for _, it := range to.items {
		if _, ok := my[it]; !ok {
			return false
		}
	}
	return true
}

// GetNextItemsets creates next sets for next state by shifting current position
// in items
func (s tableItemset) GetNextItemsets(g *grammar) map[Id]tableItemset {
	ret := make(map[Id]tableItemset)
	for id, items := range s.getNextSetsMap() {
		ret[id] = newTableItemset(items, g)
	}
	return ret
}

// getNextSetsMap creates input data to create next sets for next states by
// shifting current position in items
func (s tableItemset) getNextSetsMap() map[Id][]tableItem {
	nextM := make(map[Id][]tableItem)
SOURCE:
	for _, it := range s.items {
		if !it.HasFurther() {
			continue
		}
		nextId := it.Expected()
		nextItem := it.Shift()
		has := nextM[nextId]
		for _, hasItem := range has {
			if hasItem == nextItem {
				continue SOURCE
			}
		}
		nextM[nextId] = append(has, nextItem)
	}
	return nextM
}

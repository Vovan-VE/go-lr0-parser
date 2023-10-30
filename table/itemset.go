package table

import (
	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

func newItemset(items []item, g grammar.Grammar) itemset {
	allItems := getAllPossibleItems(items, g)
	validateDeterministic(allItems, g)
	return itemset{items: allItems}
}

// getAllPossibleItems returns possible items by expansion all the given items
// by Grammar
func getAllPossibleItems(items []item, g grammar.Grammar) []item {
	var final []item
	knownNext := symbol.NewSetOfId()

	newItems := items
	for len(newItems) != 0 {
		var nextNew []symbol.Id
	NewItems:
		for _, newIt := range newItems {
			for _, have := range final {
				if newIt == have {
					continue NewItems
				}
			}

			final = append(final, newIt)

			nextId := newIt.Expected()
			if nextId == symbol.InvalidId || g.IsTerminal(nextId) {
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
				newItems = append(newItems, newItem(r))
			}
		}
	}

	return final
}

// validateDeterministic checks for bad state in this itemset. A problem is
// reported by panic since it's grammar definition problem
//
// - ErrConflictReduceReduce
//
// - ErrConflictShiftReduce
//
// https://en.wikipedia.org/wiki/LR_parser#Conflicts_in_the_constructed_tables
func validateDeterministic(items []item, g grammar.Grammar) {
	var finite, terminals, nonTerminals []item
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
		if restTerminals.IsEmpty() {
			panic(errors.WithStack(ErrConflictShiftReduce))
		}
	}

	// Check for Reduce-Reduce conflicts
	if len(finite) > 1 {
		panic(errors.WithStack(ErrConflictReduceReduce))
	}
}

type itemset struct {
	items []item
}

// HasFinalItem checks if this set has item where the next token must be EOF
func (s itemset) HasFinalItem() bool {
	for _, it := range s.items {
		if it.HasEOF() && !it.HasFurther() {
			return true
		}
	}
	return false
}

// ReduceRule returns reduction rule of this set if any, nil otherwise
func (s itemset) ReduceRule() grammar.RuleImplementation {
	for _, it := range s.items {
		if !it.HasFurther() {
			// this is the only due to Reduce-Reduce conflicts validation
			return it.RuleImplementation
		}
	}
	return nil
}

// IsEqual checks equality with another itemset
func (s itemset) IsEqual(to itemset) bool {
	if len(s.items) != len(to.items) {
		return false
	}
	my := make(map[item]struct{})
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
func (s itemset) GetNextItemsets(g grammar.Grammar) map[symbol.Id]itemset {
	ret := make(map[symbol.Id]itemset)
	for id, items := range s.getNextSetsMap() {
		ret[id] = newItemset(items, g)
	}
	return ret
}

// getNextSetsMap creates input data to create next sets for next states by
// shifting current position in items
func (s itemset) getNextSetsMap() map[symbol.Id][]item {
	nextM := make(map[symbol.Id][]item)
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

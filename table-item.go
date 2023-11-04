package lr0

import (
	"github.com/pkg/errors"
)

// Create tableItem from a Rule with position pointing to first symbol
//
//	input rule:
//		Sum : Sum "+" Product
//	output tableItem:
//		Sum : > Sum "+" Product
func newTableItem(r Rule) tableItem {
	return tableItem{Rule: r}
}

// Item of parser state items set
//
// Item is a rule-like object, where definition body is spited by Current
// Position into passed and further symbols. So each Rule can produce exactly
// N+1 Items where N is count of symbols in Rule definition.
//
//	rule:
//		Sum : Sum "+" Product
//	items possible:
//		Sum : > Sum   "+"   Product
//		Sum :   Sum > "+"   Product
//		Sum :   Sum   "+" > Product
//		Sum :   Sum   "+"   Product >
type tableItem struct {
	Rule
	nextIndex int
}

// Expected returns next expected Id or InvalidId in the last
// position
func (i tableItem) Expected() Id {
	if !i.HasFurther() {
		return InvalidId
	}
	return i.Definition()[i.nextIndex]
}

// HasFurther returns whether the tableItem has further expected symbols
func (i tableItem) HasFurther() bool {
	return i.nextIndex < len(i.Definition())
}

// Shift creates a new tableItem by shifting current position to the next symbol
//
//	in
//		Sum : Sum > "+"   Product
//	out
//		Sum : Sum   "+" > Product
func (i tableItem) Shift() tableItem {
	if !i.HasFurther() {
		panic(errors.Wrap(ErrInternal, "bad usage"))
	}
	next := i
	next.nextIndex++
	return next
}

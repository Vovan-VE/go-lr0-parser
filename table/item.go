package table

import (
	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

// Create item from a Rule with position pointing to first symbol
//
//	input rule:
//		Sum : Sum "+" Product
//	output item:
//		Sum : > Sum "+" Product
func newItem(r grammar.Rule) item {
	return item{Rule: r}
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
type item struct {
	grammar.Rule
	nextIndex int
}

// Expected returns next expected symbol.Id or symbol.InvalidId in the last
// position
func (i item) Expected() symbol.Id {
	if !i.HasFurther() {
		return symbol.InvalidId
	}
	return i.Definition()[i.nextIndex]
}

// HasFurther returns whether the item has further expected symbols
func (i item) HasFurther() bool {
	return i.nextIndex < len(i.Definition())
}

// Shift creates a new item by shifting current position to the next symbol
//
//	in
//		Sum : Sum > "+"   Product
//	out
//		Sum : Sum   "+" > Product
func (i item) Shift() item {
	if !i.HasFurther() {
		panic(errors.Wrap(symbol.ErrInternal, "bad usage"))
	}
	next := i
	next.nextIndex++
	return next
}

package stack

import (
	"github.com/vovan-ve/go-lr0-parser/symbol"
	"github.com/vovan-ve/go-lr0-parser/table"
)

type item struct {
	state table.StateIndex
	node  symbol.Id
	value any
	// at *lexer.State - can be useful or not - then complicated calc api
}

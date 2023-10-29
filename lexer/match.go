package lexer

import (
	"github.com/vovan-ve/go-lr0-parser/symbol"
)

// Match is a found token representation
type Match struct {
	// Which Terminal found
	Term symbol.Id
	// What value it returned
	Value any
}

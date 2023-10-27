package lexer

// Match is a found token representation
type Match struct {
	// Which Term found
	Term Term
	// What value it returned
	Value any
}

package lr0

// Match is a found token representation
type Match struct {
	// Which Terminal found
	Term Id
	// What value it returned
	Value any
}

package lr0

import (
	"github.com/vovan-ve/go-lr0-parser/internal/grammar"
	"github.com/vovan-ve/go-lr0-parser/internal/lexer"
	"github.com/vovan-ve/go-lr0-parser/internal/parser"
	"github.com/vovan-ve/go-lr0-parser/internal/symbol"
	"github.com/vovan-ve/go-lr0-parser/internal/table"
)

var (
	// ErrDefine will be raised by panic in case of invalid definition
	ErrDefine = symbol.ErrDefine
	// ErrInternal is an internal error which actually should not be raised, but
	// coded in some panic for debug purpose just in case
	ErrInternal = symbol.ErrInternal
	// ErrNegativeOffset will be raised by panic if some operation with State
	// cause negative offset
	ErrNegativeOffset = lexer.ErrNegativeOffset
	// ErrParse is base error for run-time errors about parsing.
	//
	//	errors.Is(err, ErrParse)
	//
	//	errors.Wrap(ErrParse, "unexpected thing found")
	ErrParse = lexer.ErrParse
	// ErrState is base wrap error for parsing state
	ErrState = table.ErrState
	// ErrConflictReduceReduce means that there are a number of rules which
	// applicable to reduce in the current state. Wraps ErrState.
	ErrConflictReduceReduce = table.ErrConflictReduceReduce
	// ErrConflictShiftReduce means that Shift and Reduce are both applicable
	// in the current state. Wraps ErrState.
	ErrConflictShiftReduce = table.ErrConflictShiftReduce
)

// Id is an identifier for terminals and non-terminals
//
// Zero value is InvalidId and must not be used:
//
//	const (
//		TInt Id = iota + 1
//		TPlus
//		TMinus
//
//		NValue
//		NSum
//		NGoal
//	)
type Id = symbol.Id

type Symbol = symbol.Symbol

// State describes an immutable state of reading the underlying buffer at the
// given position
type State = lexer.State

// MatchFunc is a signature for common function to match an underlying token.
//
// It accepts current State to parse from.
//
// If the token parsed, the function returns next State to continue from and
// the token value from ToValue.
//
// If the token was not parsed, the function returns `nil, nil`.
//
// Must not return the same State as input State.
type MatchFunc = func(*State) (next *State, value any)

type TerminalFactory = lexer.TerminalFactory

// Terminal is interface to parse specific type of token from input State
type Terminal = lexer.Terminal

// NonTerminal is non-terminal definition describing how to parse it from
// underlying parts.
type NonTerminal = grammar.NonTerminal

// NonTerminalDefinition is an interface to make Rules for non-terminal
type NonTerminalDefinition = grammar.NonTerminalDefinition

// Parser is object preconfigured for a specific grammar, ready to parse an
// input to evaluate the result.
type Parser = parser.Parser

const (
	// InvalidId id zero value for Id. It's used internally, and it's not
	// allowed to use in definition.
	InvalidId = symbol.InvalidId
)

// New created new Parser
//
// terminals can be defined with NewTerm or NewWhitespace
//
// rules can be defined by NewNT
//
//	parser := New(
//		[]Terminal{
//			NewTerm(tInt, "int").Func(matchInt),
//			NewTerm(tPlus, `"+"`).Hide().Str("+"),
//			NewWhitespace().Func(matchSpaces),
//		},
//		[]NonTerminalDefinition{
//			NewNT(nGoal, "Goal").Main().Is(nSum),
//			NewNT(nSum, "Sum").
//				Is(nSum, tPlus, nVal).Do(func (a, b int) int { return a+b }).
//				Is(nVal),
//			NewNT(nVal, "Val").Is(tInt),
//		},
//	)
//	result, err := parser.Parse(NewState([]byte("42 + 37")))
//	if err != nil {
//		fmt.Println("error", err)
//	} else {
//		fmt.Println("result", result)
//	}
func New(terminals []Terminal, rules []NonTerminalDefinition) Parser {
	return parser.New(grammar.New(terminals, rules))
}

// NewTerm starts new Terminal creation.
//
//	NewTerm(tInt, "integer").Func(matchDigits)
//	NewTerm(tPlus, "plus").Hide().Byte('+')
func NewTerm(id Id, name string) *TerminalFactory { return lexer.NewTerm(id, name) }

// NewWhitespace can be used to define internal terminals to skip whitespaces
//
// Can be used multiple times to define different kinds of whitespaces.
//
// Whitespace tokens will be silently skipped before every terminal match
//
//	NewWhitespace().Func(matchSpaces),
//	NewWhitespace().Func(matchComment),
func NewWhitespace() *TerminalFactory { return lexer.NewWhitespace() }

// NewNT created new non-terminal definition
//
//	NewNT(nGoal, "Goal").Main().
//		Is(nSum)
//
//	NewNT(nSum, "Sum").
//		Is(nSum, tPlus, nVal).Do(func (a, b int) int { return a+b }).
//		// ^^^^  ^^^^^  ^^^^           ^  ^
//		//  a    hidden   b   ---------+--'
//		Is(nSum, tMinus, nVal).Do(func (a, b int) int { return a-b }).
//		Is(nVal)
//		// -----^^^^ here is no `Do(func (v int) int { return v })`
//		//           `Do(nil)` will do the same in this case
func NewNT(id Id, name string) *NonTerminal { return grammar.NewNT(id, name) }

// NewState creates new State for the given buffer `input` pointing to its start
func NewState(input []byte) *State { return lexer.NewState(input) }

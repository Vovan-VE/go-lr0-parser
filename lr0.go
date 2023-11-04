package lr0

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	// ErrDefine will be raised by panic in case of invalid definition
	ErrDefine = errors.New("invalid definition")
	// ErrInternal is an internal error which actually should not be raised, but
	// coded in some panic for debug purpose just in case
	ErrInternal = errors.New("internal error")
	// ErrNegativeOffset will be raised by panic if some operation with State
	// cause negative offset
	// ErrNegativeOffset will be raised by panic if some operation with State
	// cause negative offset
	ErrNegativeOffset = errors.New("negative position")
	// ErrParse is base error for run-time errors about parsing.
	//
	//	errors.Is(err, ErrParse)
	//
	//	errors.Wrap(ErrParse, "unexpected thing found")
	ErrParse = errors.New("parse error")
	// ErrState is base wrap error for parsing state
	ErrState = errors.Wrap(ErrDefine, "bad state for table")
	// ErrConflictReduceReduce means that there are a number of rules which
	// applicable to reduce in the current state.
	ErrConflictReduceReduce = errors.Wrap(ErrState, "reduce-reduce conflict")
	// ErrConflictShiftReduce means that Shift and Reduce are both applicable
	// in the current state.
	ErrConflictShiftReduce = errors.Wrap(ErrState, "shift-reduce conflict")
)

// Id is an identifier for terminals and non-terminals
//
// Only positive values must be used. Zero value is `InvalidId` and negative
// values are reserved.
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
type Id int

const (
	// InvalidId id zero value for Id. It's used internally, and it's not
	// allowed to use in definition.
	InvalidId Id = 0
)

// Symbol is common interface to describe Symbol meta data
type Symbol interface {
	Id() Id
	// Name returns a human-recognizable name to not mess up with numeric Term
	Name() string
}

type SymbolRegistry interface {
	SymbolName(id Id) string
}

// MatchFunc is a signature for common function to match an underlying token.
//
// It accepts current State to parse from.
//
// If the token parsed, the function returns next State to continue from and
// evaluated token value.
//
// If the token was not parsed, the function returns `nil, nil`.
//
// Must not return the same State as input State.
type MatchFunc = func(*State) (next *State, value any)

// Terminal is interface to parse specific type of token from input State
type Terminal interface {
	Symbol
	// IsHidden returns whether the terminal is hidden
	//
	// Hidden terminal does not produce a value to calc non-terminal value.
	// For example if in the following rule:
	//	Sum : Sum plus Val
	// a `plus` terminal is hidden, then only two values will be passed to calc
	// function - value of `Sum` and value of `Val`:
	//	func(sum any, val any) any
	IsHidden() bool
	// Match is MatchFunc
	Match(*State) (next *State, value any)
}

// NonTerminalDefinition is an interface to make Rules for non-terminal
type NonTerminalDefinition interface {
	Symbol
	GetRules(l NamedHiddenRegistry) []Rule
}

type NamedHiddenRegistry interface {
	SymbolRegistry
	IsHidden(id Id) bool
}

// Rule is one of possible definition for a non-Terminal
type Rule interface {
	fmt.Stringer
	// Subject of the rule
	Subject() Id
	// HasEOF tells whether EOF must be found in the end of input
	HasEOF() bool
	// Definition of what it consists of
	Definition() []Id
	Value([]any) (any, error)
	IsHidden(index int) bool
}

// Parser is object preconfigured for a specific grammar, ready to parse an
// input to evaluate the result.
type Parser interface {
	// Parse parses the whole input stream State.
	//
	// Returns either evaluated result or error.
	Parse(input *State) (result any, err error)
}

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
	return newParser(newGrammar(terminals, rules))
}

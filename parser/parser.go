package parser

import (
	"io"

	"github.com/pkg/errors"
	"github.com/vovan-ve/go-lr0-parser/grammar"
	"github.com/vovan-ve/go-lr0-parser/lexer"
	"github.com/vovan-ve/go-lr0-parser/stack"
	"github.com/vovan-ve/go-lr0-parser/table"
)

type Parser interface {
	// Parse parses the whole input stream State.
	//
	// Returns either evaluated result or error.
	Parse(input *lexer.State) (result any, err error)
}

func New(g grammar.Grammar) Parser {
	return &parser{
		g: g,
		t: table.New(g),
	}
}

type parser struct {
	g grammar.Grammar
	t table.Table
}

func (p *parser) Parse(input *lexer.State) (result any, err error) {
	st := stack.New(p.t)

	var (
		next = input
		ok   bool
		to   table.StateIndex
	)
Goal:
	for {
		at := next
		for st.Current().IsReduceOnly() {
			ok, err = st.Reduce()
			if err != nil {
				return nil, lexer.WithSource(err, at)
			}
			if !ok {
				// if this happens ever?
				// REFACT: looks like this will never happen now
				// no reduce rule - unexpected input
				//st.Current().TerminalsSet()
				return nil, lexer.WithSource(lexer.NewParseError("unexpected input 1"), at)
			}
		}

		var m *lexer.Match
		if !next.IsEOF() {
			next, m, err = p.g.Match(next, st.Current().TerminalsSet())
			if err != nil && err != io.EOF {
				return nil, errors.Wrap(err, "unexpected input")
			}
		}

		for {
			if m != nil {
				if to, ok = st.Current().TerminalAction(m.Term); ok {
					st.Shift(to, m.Term, m.Value)
					break
				}
			}
			if st.Current().AcceptEof() {
				if m == nil {
					break Goal
				}
				return nil, lexer.WithSource(lexer.NewParseError("unexpected input instead of EOF"), at)
			}

			ok, err = st.Reduce()
			if err != nil {
				return nil, lexer.WithSource(err, at)
			}
			if !ok {
				return nil, lexer.WithSource(p.g.ExpectationError(st.Current().TerminalsSet(), "unexpected input"), at)
			}
		}
	}
	return st.Done(), nil
}

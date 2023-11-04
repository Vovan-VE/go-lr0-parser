package lr0

import (
	"io"

	"github.com/pkg/errors"
)

func newParser(g *grammar) Parser {
	return &parser{
		g: g,
		t: newTable(g),
	}
}

type parser struct {
	g *grammar
	t *table
}

func (p *parser) Parse(input *State) (result any, err error) {
	st := newStack(p.t)

	var (
		next = input
		ok   bool
		to   tableStateIndex
	)
Goal:
	for {
		at := next
		for st.Current().IsReduceOnly() {
			ok, err = st.Reduce()
			if err != nil {
				return nil, WithSource(err, at)
			}
			if !ok {
				// if this happens ever?
				// REFACT: looks like this will never happen now
				// no reduce rule - unexpected input
				//st.Current().TerminalsSet()
				return nil, WithSource(NewParseError("unexpected input 1"), at)
			}
		}

		var m *Match
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
				return nil, WithSource(NewParseError("unexpected input instead of EOF"), at)
			}

			ok, err = st.Reduce()
			if err != nil {
				return nil, WithSource(err, at)
			}
			if !ok {
				return nil, WithSource(p.g.ExpectationError(st.Current().TerminalsSet(), "unexpected input"), at)
			}
		}
	}
	return st.Done(), nil
}

package lr0

import (
	"github.com/pkg/errors"
)

//// Stack stores a parser state while last is processing an input stream
//type Stack interface {
//	// Current state Row from Table
//	Current() Row
//	// Shift does shift - push the given item into Stack
//	Shift(si tableStateIndex, id Id, value any)
//	// Reduce tries to perform reduce in current state. If no ReduceRule
//	// available, returns `false, nil`. If an error occurred while calculating
//	// a value, `false, error` will be returned. Of success `true, nil` will be
//	// returned.
//	Reduce() (bool, error)
//	// Done ends work, returns the final value from Stack and resets Stack to
//	// initial state
//	Done() any
//}

// newStack creates new stack for the given Table
func newStack(t *table) *stack {
	st := &stack{t: t}
	st.set(0)
	return st
}

type stack struct {
	t     *table
	items []stackItem
	si    tableStateIndex
	// cached `.t.Row(.si)`
	row *tableRow
}

func (s *stack) Current() *tableRow { return s.row }

func (s *stack) Shift(si tableStateIndex, id Id, value any) {
	s.set(si)
	s.items = append(s.items, stackItem{
		state: si,
		node:  id,
		value: value,
	})
}

func (s *stack) Reduce() (bool, error) {
	rule := s.row.ReduceRule()
	if rule == nil {
		return false, nil
	}

	reduceCount := len(rule.Definition())
	totalCount := len(s.items)
	if totalCount < reduceCount {
		panic(errors.Wrap(ErrInternal, "not enough items in stack"))
	}
	nextCount := totalCount - reduceCount

	values := make([]any, 0, reduceCount)
	def := rule.Definition()
	for i, it := range s.items[nextCount:] {
		if it.node != def[i] {
			panic(errors.Wrap(ErrInternal, "unexpected stack content"))
		}
		if !rule.IsHidden(i) {
			values = append(values, it.value)
		}
	}
	newValue, err := rule.Value(values)
	if err != nil {
		return false, err
	}

	var baseSI tableStateIndex
	if totalCount > reduceCount {
		baseSI = s.items[nextCount-1].state
	}
	baseRow := s.t.Row(baseSI)

	newId := rule.Subject()
	newSI, ok := baseRow.GotoAction(newId)
	if !ok {
		panic(errors.Wrap(ErrInternal, "unexpected state in gotos"))
	}

	s.items = s.items[:nextCount]
	s.Shift(newSI, newId, newValue)
	return true, nil
}

func (s *stack) Done() any {
	if len(s.items) != 1 {
		panic(errors.Wrap(ErrInternal, "unexpected stack content"))
	}
	v := s.items[0].value
	s.items = nil
	s.set(0)
	return v
}

func (s *stack) set(si tableStateIndex) {
	s.si, s.row = si, s.t.Row(si)
}

type stackItem struct {
	state tableStateIndex
	node  Id
	value any
	// at *State - can be useful or not - then complicated calc api
}

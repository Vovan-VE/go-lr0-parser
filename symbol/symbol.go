package symbol

import (
	"fmt"
)

// Id is a base identifier for a Symbol
//
// Zero value must not be user as a valid id:
//
//	const (
//		NValue symbol.Id = iota + 1
//		NProduct
//		NSum
//		NGoal
//	)
type Id int

const InvalidId Id = 0

// Meta is common interface to describe Symbol meta data
type Meta interface {
	Id() Id
	// Name returns a human-recognizable name to not mess up with numeric Term
	Name() string
	// IsHidden returns whether the token is hidden
	//
	// Hidden token does not produce a value to calc non-terminal value.
	// For example if in the following rule:
	//	Sum : Sum plus Val
	// a `plus` token is hidden, then only two values - value of `Sum` and value
	// of `Val` - will be passed to calc function: `func(any, any) any`
	IsHidden() bool
}

func Dump(m Meta) string {
	if s := m.Name(); s != "" {
		return s
	}
	return fmt.Sprintf("#%v", m.Id())
}

type ReadonlySetOfId interface {
	IsEmpty() bool
	Count() int
	Has(id Id) bool
	ForEach(fn func(Id))
}

type SetOfId map[Id]struct{}

func NewSetOfId(id ...Id) SetOfId {
	return make(SetOfId).Add(id...)
}

func (s SetOfId) Add(id ...Id) SetOfId {
	for _, v := range id {
		s[v] = struct{}{}
	}
	return s
}

func (s SetOfId) Remove(id Id)  { delete(s, id) }
func (s SetOfId) IsEmpty() bool { return len(s) == 0 }
func (s SetOfId) Count() int    { return len(s) }

func (s SetOfId) ForEach(fn func(Id)) {
	for id := range s {
		fn(id)
	}
}

func (s SetOfId) Has(id Id) bool {
	_, ok := s[id]
	return ok
}

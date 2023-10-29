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

// Tag is optional for extra identity when one non-Terminal has multiple
// alternative Rules
type Tag byte

// Meta is common interface to describe Symbol meta data
type Meta interface {
	Id() Id
	// Name returns a human-recognizable name to not mess up with numeric Term
	Name() string
	// IsHidden returns whether the token is hidden
	// TODO: docs
	IsHidden() bool
}

func Dump(m Meta) string {
	if s := m.Name(); s != "" {
		return s
	}
	return fmt.Sprintf("#%v", m.Id())
}

type SetOfId map[Id]struct{}

func NewSetOfId() SetOfId { return make(SetOfId) }

func (s SetOfId) Add(id ...Id) {
	for _, v := range id {
		s[v] = struct{}{}
	}
}

func (s SetOfId) Remove(id Id)  { delete(s, id) }
func (s SetOfId) IsEmpty() bool { return len(s) == 0 }

func (s SetOfId) ForEach(fn func(Id)) {
	for id := range s {
		fn(id)
	}
}

func (s SetOfId) Has(id Id) bool {
	_, ok := s[id]
	return ok
}

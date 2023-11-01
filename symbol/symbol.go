package symbol

import (
	"fmt"
)

// Id is an identifier for terminals and non-terminals
//
// Zero value must not be used:
//
//	const (
//		TInt symbol.Id = iota + 1
//		TPlus
//		TMinus
//
//		NValue
//		NSum
//		NGoal
//	)
type Id int

const InvalidId Id = 0

// Symbol is common interface to describe Symbol meta data
type Symbol interface {
	Id() Id
	// Name returns a human-recognizable name to not mess up with numeric Term
	Name() string
}

func Dump(m Symbol) string {
	if s := m.Name(); s != "" {
		return s
	}
	return fmt.Sprintf("#%v", m.Id())
}

type ReadonlySet interface {
	Count() int
	Has(id Id) bool

	//IsEmpty() bool
	//ForEach(fn func(Id))
}

type Set map[Id]struct{}

func NewSetOfId(id ...Id) Set {
	return make(Set).Add(id...)
}

func (s Set) Add(id ...Id) Set {
	for _, v := range id {
		s[v] = struct{}{}
	}
	return s
}

func (s Set) Remove(id Id) { delete(s, id) }
func (s Set) Count() int   { return len(s) }

//func (s Set) IsEmpty() bool { return len(s) == 0 }
//func (s Set) ForEach(fn func(Id)) {
//	for id := range s {
//		fn(id)
//	}
//}

func (s Set) Has(id Id) bool {
	_, ok := s[id]
	return ok
}

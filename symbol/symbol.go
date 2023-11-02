package symbol

import (
	"fmt"
)

// Id is an identifier for terminals and non-terminals
//
// Only positive values must be used. Zero value is `InvalidId` and negative
// values are reserved.
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

type Registry interface {
	SymbolName(id Id) string
}

func Dump(s Symbol) string {
	if n := s.Name(); n != "" {
		return n
	}
	return fmt.Sprintf("#%v", s.Id())
}

func DumpId(id Id, r Registry) string {
	if s := r.SymbolName(id); s != "" {
		return s
	}
	return fmt.Sprintf("#%v", id)
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

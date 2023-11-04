package lr0

import (
	"fmt"
)

func dumpSymbol(s Symbol) string {
	if n := s.Name(); n != "" {
		return n
	}
	return fmt.Sprintf("#%v", s.Id())
}

func dumpId(id Id, r SymbolRegistry) string {
	if s := r.SymbolName(id); s != "" {
		return s
	}
	return fmt.Sprintf("#%v", id)
}

type readonlyIdSet interface {
	Count() int
	Has(id Id) bool

	//IsEmpty() bool
	//ForEach(fn func(Id))
}

type idSet map[Id]struct{}

func newIdSet(id ...Id) idSet {
	return make(idSet).Add(id...)
}

func (s idSet) Add(id ...Id) idSet {
	for _, v := range id {
		s[v] = struct{}{}
	}
	return s
}

func (s idSet) Remove(id Id) { delete(s, id) }
func (s idSet) Count() int   { return len(s) }

//func (s idSet) IsEmpty() bool { return len(s) == 0 }
//func (s idSet) ForEach(fn func(Id)) {
//	for id := range s {
//		fn(id)
//	}
//}

func (s idSet) Has(id Id) bool {
	_, ok := s[id]
	return ok
}

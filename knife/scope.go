package knife

import (
	"fmt"
	"go/token"
	"go/types"
)

type Scope struct {
	TypesScope *types.Scope
	Parent     *Scope
	Children   []*Scope
	Pos        token.Pos
	End        token.Pos
	Objects    map[string]Object
	Names      []string
}

var _ fmt.Stringer = (*Scope)(nil)

func NewScope(s *types.Scope) *Scope {
	if s == nil {
		return nil
	}

	v, _ := cache.Load(s)
	cached, _ := v.(*Scope)
	if cached != nil {
		return cached
	}

	var ns Scope
	cache.Store(s, &ns)

	ns.TypesScope = s
	ns.Parent = NewScope(s.Parent())
	ns.Children = make([]*Scope, s.NumChildren())
	for i := range ns.Children {
		ns.Children[i] = NewScope(s.Child(i))
	}
	ns.Pos = s.Pos()
	ns.End = s.End()
	ns.Objects = make(map[string]Object, s.Len())
	ns.Names = make([]string, s.Len())
	for i, name := range s.Names() {
		ns.Names[i] = name
		o := s.Lookup(name)
		ns.Objects[name] = NewObject(o)
	}

	return &ns
}

func (s *Scope) String() string {
	return s.TypesScope.String()
}

package knife

import (
	"fmt"
	"go/types"
	"iter"
	"sync"
)

var (
	cache sync.Map
)

type Package struct {
	TypesPackage *types.Package
	Name         string
	Path         string
	Imports      []*Package
	Funcs        map[string]*Func
	FuncNames    []string
	Vars         map[string]*Var
	VarNames     []string
	Consts       map[string]*Const
	ConstNames   []string
	Types        map[string]*TypeName
	TypeNames    []string
}

var _ fmt.Stringer = (*Package)(nil)

func (pkg *Package) String() string {
	return pkg.TypesPackage.String()
}

func NewPackage(pkg *types.Package) *Package {
	if pkg == nil {
		return nil
	}

	v, _ := cache.Load(pkg)
	cached, _ := v.(*Package)
	if cached != nil {
		return cached
	}

	var np Package
	cache.Store(pkg, &np)

	np.TypesPackage = pkg
	np.Name = pkg.Name()
	np.Path = pkg.Path()
	np.Imports = make([]*Package, len(pkg.Imports()))
	np.Funcs = map[string]*Func{}
	np.Vars = map[string]*Var{}
	np.Consts = map[string]*Const{}
	np.Types = map[string]*TypeName{}

	for i, p := range pkg.Imports() {
		np.Imports[i] = NewPackage(p)
	}

	for _, n := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(n)
		switch obj := obj.(type) {
		case *types.Func:
			np.Funcs[n] = NewFunc(obj)
			np.FuncNames = append(np.FuncNames, n)
		case *types.Var:
			np.Vars[n] = NewVar(obj)
			np.VarNames = append(np.VarNames, n)
		case *types.Const:
			np.Consts[n] = NewConst(obj)
			np.ConstNames = append(np.ConstNames, n)
		case *types.TypeName:
			np.Types[n] = NewTypeName(obj)
			np.TypeNames = append(np.TypeNames, n)
		}
	}

	return &np

}

func (pkg *Package) Objects() iter.Seq2[string, Object] {
	return func(yield func(string, Object) bool) {
		for name, f := range pkg.Funcs {
			if !yield(name, f) {
				return
			}
		}

		for name, v := range pkg.Vars {
			if !yield(name, v) {
				return
			}
		}

		for name, c := range pkg.Consts {
			if !yield(name, c) {
				return
			}
		}
	}
}

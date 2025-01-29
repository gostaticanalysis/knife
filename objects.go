package knife

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"reflect"
)

type Object interface {
	TypesObject() types.Object
}

func NewObject(o types.Object) Object {
	switch o := o.(type) {
	case *types.Var:
		if !o.IsField() {
			return NewVar(o)
		}
	case *types.Const:
		return NewConst(o)
	case *types.Func:
		return NewFunc(o)
	case *types.TypeName:
		return NewTypeName(o)
	}
	return nil
}

type Field struct {
	TypesVar  *types.Var
	Struct    *Struct
	Tag       string
	Anonymous bool
	Exported  bool
	Name      string
	Type      *Type
}

var _ fmt.Stringer = (*Field)(nil)
var _ Object = (*Field)(nil)

func NewField(s *Struct, v *types.Var, tag string) *Field {
	if s == nil || v == nil {
		return nil
	}

	_v, _ := cache.Load(v)
	cached, _ := _v.(*Field)
	if cached != nil {
		return cached
	}

	var nf Field
	cache.Store(v, &nf)
	nf.TypesVar = v
	nf.Struct = s
	nf.Tag = tag
	nf.Anonymous = v.Anonymous()
	nf.Exported = v.Exported()
	nf.Name = v.Name()
	nf.Type = NewType(v.Type())

	return &nf
}

func (f *Field) Pos() token.Pos {
	return f.TypesVar.Pos()
}

func (f *Field) String() string {
	return f.TypesVar.String()
}

func (f *Field) TypesObject() types.Object {
	return f.TypesVar
}

type Var struct {
	TypesVar *types.Var
	Exported bool
	Name     string
	Type     *Type
	Package  *Package
}

var _ fmt.Stringer = (*Var)(nil)
var _ Object = (*Var)(nil)

func NewVar(v *types.Var) *Var {
	if v == nil {
		return nil
	}

	_v, _ := cache.Load(v)
	cached, _ := _v.(*Var)
	if cached != nil {
		return cached
	}

	var nv Var
	cache.Store(v, &nv)
	nv.TypesVar = v
	nv.Exported = v.Exported()
	nv.Name = v.Name()
	nv.Type = NewType(v.Type())
	nv.Package = NewPackage(v.Pkg())

	return &nv
}

func (v *Var) Pos() token.Pos {
	return v.TypesVar.Pos()
}

func (v *Var) String() string {
	return v.TypesVar.String()
}

func (v *Var) TypesObject() types.Object {
	return v.TypesVar
}

type Func struct {
	TypesFunc *types.Func
	Name      string
	Exported  bool
	Package   *Package
	Signature *Signature
}

var _ fmt.Stringer = (*Func)(nil)
var _ Object = (*Func)(nil)

func (f *Func) Pos() token.Pos {
	return f.TypesFunc.Pos()
}

func (f *Func) String() string {
	return f.TypesFunc.String()
}

func (f *Func) TypesObject() types.Object {
	return f.TypesFunc
}

func NewFunc(f *types.Func) *Func {
	if f == nil {
		return nil
	}

	v, _ := cache.Load(f)
	cached, _ := v.(*Func)
	if cached != nil {
		return cached
	}

	var nf Func
	cache.Store(f, &nf)
	nf.TypesFunc = f
	nf.Name = f.Name()
	nf.Exported = f.Exported()
	nf.Package = NewPackage(f.Pkg())
	nf.Signature = NewSignature(f.Type().(*types.Signature))

	return &nf
}

type TypeName struct {
	TypesTypeName *types.TypeName
	Exported      bool
	IsAlias       bool
	Name          string
	Package       *Package
	Type          *Type
}

var _ fmt.Stringer = (*TypeName)(nil)
var _ Object = (*TypeName)(nil)

func (tn *TypeName) Pos() token.Pos {
	return tn.TypesTypeName.Pos()
}

func (tn *TypeName) String() string {
	return tn.TypesTypeName.String()
}

func (tn *TypeName) TypesObject() types.Object {
	return tn.TypesTypeName
}

func NewTypeName(tn *types.TypeName) *TypeName {
	if tn == nil {
		return nil
	}

	v, _ := cache.Load(tn)
	cached, _ := v.(*TypeName)
	if cached != nil {
		return cached
	}

	var ntn TypeName
	cache.Store(tn, &ntn)
	ntn.TypesTypeName = tn
	ntn.Exported = tn.Exported()
	ntn.IsAlias = tn.IsAlias()
	ntn.Name = tn.Name()
	ntn.Package = NewPackage(tn.Pkg())
	ntn.Type = NewType(tn.Type())

	return &ntn
}

type Const struct {
	TypesConst *types.Const
	Exported   bool
	Name       string
	Package    *Package
	Type       *Type
	Value      constant.Value
}

var _ fmt.Stringer = (*Const)(nil)
var _ Object = (*Const)(nil)

func (c *Const) Pos() token.Pos {
	return c.TypesConst.Pos()
}

func (c *Const) String() string {
	return c.TypesConst.String()
}

func (c *Const) TypesObject() types.Object {
	return c.TypesConst
}

func NewConst(c *types.Const) *Const {
	if c == nil {
		return nil
	}

	v, _ := cache.Load(c)
	cached, _ := v.(*Const)
	if cached != nil {
		return cached
	}

	var nc Const
	cache.Store(c, &nc)
	nc.TypesConst = c
	nc.Exported = c.Exported()
	nc.Name = c.Name()
	nc.Package = NewPackage(c.Pkg())
	nc.Type = NewType(c.Type())
	nc.Value = c.Val()

	return &nc
}

func (c *Const) BoolVal() bool {
	return constant.BoolVal(c.Value)
}

func (c *Const) StringVal() string {
	return constant.StringVal(c.Value)
}

func (c *Const) Float32Val() float32 {
	v, ok := constant.Float32Val(c.Value)
	if !ok {
		panic("unkown kind")
	}
	return v
}

func (c *Const) Float64Val() float64 {
	v, ok := constant.Float64Val(c.Value)
	if !ok {
		panic("unkown kind")
	}
	return v
}

func (c *Const) Int64Val() int64 {
	v, ok := constant.Int64Val(c.Value)
	if !ok {
		panic("unkown kind")
	}
	return v

}

func (c *Const) Uint64Val() uint64 {
	v, ok := constant.Uint64Val(c.Value)
	if !ok {
		panic("unkown kind")
	}
	return v

}

func (c *Const) Val() any {
	return constant.Val(c.Value)
}

func Position(fset *token.FileSet, v any) token.Position {
	n, ok := v.(interface{ Pos() token.Pos })
	if ok && fset != nil {
		return fset.Position(n.Pos())
	}
	return token.Position{}
}

func Exported(list any) any {
	v := reflect.ValueOf(list)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return exportedSlice(v)
	case reflect.Map:
		return exportedMap(v)
	}
	panic("unexpected kind")
}

func exportedSlice(v reflect.Value) any {
	result := reflect.MakeSlice(v.Type(), 0, 0)
	for i := 0; i < v.Len(); i++ {
		elm := v.Index(i)
		if isExported(elm) {
			result = reflect.Append(result, elm)
		}
	}
	return result.Interface()
}

func exportedMap(v reflect.Value) any {
	result := reflect.MakeMap(v.Type())
	for _, key := range v.MapKeys() {
		elm := v.MapIndex(key)
		if isExported(elm) {
			result.SetMapIndex(key, elm)
		}
	}
	return result.Interface()
}

func isExported(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr:
		return isExported(v.Elem())
	case reflect.Struct:
		return v.FieldByName("Exported").Bool()
	}
	panic("unexpected kind")
}

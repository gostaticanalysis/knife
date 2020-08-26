package knife

import (
	"fmt"
	"go/types"
)

type Type struct {
	TypesType types.Type
}

var _ fmt.Stringer = (*Type)(nil)

func NewType(t types.Type) *Type {
	if t == nil {
		return nil
	}

	v, _ := cache.Load(t)
	cached, _ := v.(*Type)
	if cached != nil {
		return cached
	}
	nt := &Type{TypesType: t}
	cache.Store(t, nt)
	return nt
}

func (t *Type) Underlying() *Type {
	return NewType(t.TypesType.Underlying())
}

func (t *Type) String() string {
	return t.TypesType.String()
}

func (t *Type) Array() *Array {
	a, _ := under(t.TypesType).(*types.Array)
	return NewArray(a)
}

func (t *Type) Slice() *Slice {
	s, _ := under(t.TypesType).(*types.Slice)
	return NewSlice(s)
}

func (t *Type) Struct() *Struct {
	s, _ := under(t.TypesType).(*types.Struct)
	return NewStruct(s)
}

func (t *Type) Map() *Map {
	m, _ := under(t.TypesType).(*types.Map)
	return NewMap(m)
}

func (t *Type) Pointer() *Pointer {
	p, _ := under(t.TypesType).(*types.Pointer)
	return NewPointer(p)
}

func (t *Type) Chan() *Chan {
	c, _ := under(t.TypesType).(*types.Chan)
	return NewChan(c)
}

func (t *Type) Basic() *Basic {
	b, _ := under(t.TypesType).(*types.Basic)
	return NewBasic(b)
}

func (t *Type) Interface() *Interface {
	i, _ := under(t.TypesType).(*types.Interface)
	return NewInterface(i)
}

func (t *Type) Signature() *Signature {
	s, _ := under(t.TypesType).(*types.Signature)
	return NewSignature(s)
}

func (t *Type) Named() *Named {
	n, _ := under(t.TypesType).(*types.Named)
	return NewNamed(n)
}

type Array struct {
	TypesArray *types.Array
	Elem       *Type
	Len        int64
}

var _ fmt.Stringer = (*Array)(nil)

func NewArray(a *types.Array) *Array {
	if a == nil {
		return nil
	}

	v, _ := cache.Load(a)
	cached, _ := v.(*Array)
	if cached != nil {
		return cached
	}

	var na Array
	cache.Store(a, &na)
	na.TypesArray = a
	na.Elem = NewType(a.Elem())
	na.Len = a.Len()
	return &na
}

func ToArray(t interface{}) *Array {
	switch t := t.(type) {
	case *Type:
		return t.Array()
	case types.Type:
		return NewType(t).Array()
	}
	return nil
}

func (a *Array) String() string {
	return a.TypesArray.String()
}

type Slice struct {
	TypesSlice *types.Slice
	Elem       *Type
}

var _ fmt.Stringer = (*Slice)(nil)

func NewSlice(s *types.Slice) *Slice {
	if s == nil {
		return nil
	}

	v, _ := cache.Load(s)
	cached, _ := v.(*Slice)
	if cached != nil {
		return cached
	}

	var ns Slice
	cache.Store(s, &ns)
	ns.TypesSlice = s
	ns.Elem = NewType(s.Elem())

	return &ns
}

func ToSlice(t interface{}) *Slice {
	switch t := t.(type) {
	case *Type:
		return t.Slice()
	case types.Type:
		return NewType(t).Slice()
	}
	return nil
}

func (s *Slice) String() string {
	return s.TypesSlice.String()
}

type Struct struct {
	TypesStruct *types.Struct
	Fields      map[string]*Field
	FieldNames  []string
}

var _ fmt.Stringer = (*Struct)(nil)

func NewStruct(s *types.Struct) *Struct {
	if s == nil {
		return nil
	}

	v, _ := cache.Load(s)
	cached, _ := v.(*Struct)
	if cached != nil {
		return cached
	}

	var ns Struct
	cache.Store(s, &ns)
	ns.TypesStruct = s
	ns.Fields = make(map[string]*Field, s.NumFields())
	ns.FieldNames = make([]string, s.NumFields())

	for i := 0; i < s.NumFields(); i++ {
		v := s.Field(i)
		ns.Fields[v.Name()] = NewField(&ns, v, s.Tag(i))
		ns.FieldNames[i] = v.Name()
	}

	return &ns
}

func ToStruct(t interface{}) *Struct {
	switch t := t.(type) {
	case *Type:
		return t.Struct()
	case types.Type:
		return NewType(t).Struct()
	}
	return nil
}

func (s Struct) String() string {
	return s.TypesStruct.String()
}

type Map struct {
	TypesMap *types.Map
	Elem     *Type
	Key      *Type
}

var _ fmt.Stringer = (*Map)(nil)

func NewMap(m *types.Map) *Map {
	if m == nil {
		return nil
	}

	v, _ := cache.Load(m)
	cached, _ := v.(*Map)
	if cached != nil {
		return cached
	}

	var nm Map
	cache.Store(m, &nm)
	nm.TypesMap = m
	nm.Elem = NewType(m.Elem())
	nm.Key = NewType(m.Key())

	return &nm
}

func ToMap(t interface{}) *Map {
	switch t := t.(type) {
	case *Type:
		return t.Map()
	case types.Type:
		return NewType(t).Map()
	}
	return nil
}

func (m *Map) String() string {
	return m.TypesMap.String()
}

type Pointer struct {
	TypesPointer *types.Pointer
	Elem         *Type
}

var _ fmt.Stringer = (*Pointer)(nil)

func NewPointer(p *types.Pointer) *Pointer {
	if p == nil {
		return nil
	}

	v, _ := cache.Load(p)
	cached, _ := v.(*Pointer)
	if cached != nil {
		return cached
	}

	var np Pointer
	cache.Store(p, &np)
	np.TypesPointer = p
	np.Elem = NewType(p.Elem())

	return &np
}

func ToPointer(t interface{}) *Pointer {
	switch t := t.(type) {
	case *Type:
		return t.Pointer()
	case types.Type:
		return NewType(t).Pointer()
	}
	return nil
}

func (p *Pointer) String() string {
	return p.TypesPointer.String()
}

type Chan struct {
	TypesChan *types.Chan
	Dir       types.ChanDir
	Elem      *Type
}

var _ fmt.Stringer = (*Chan)(nil)

func NewChan(c *types.Chan) *Chan {
	if c == nil {
		return nil
	}

	v, _ := cache.Load(c)
	cached, _ := v.(*Chan)
	if cached != nil {
		return cached
	}

	var nc Chan
	cache.Store(c, &nc)
	nc.TypesChan = c
	nc.Dir = c.Dir()
	nc.Elem = NewType(c.Elem())

	return &nc
}

func ToChan(t interface{}) *Chan {
	switch t := t.(type) {
	case *Type:
		return t.Chan()
	case types.Type:
		return NewType(t).Chan()
	}
	return nil
}

func (c *Chan) String() string {
	return c.TypesChan.String()
}

type Basic struct {
	TypesBasic *types.Basic
	Info       types.BasicInfo
	Kind       types.BasicKind
	Name       string
}

var _ fmt.Stringer = (*Basic)(nil)

func NewBasic(b *types.Basic) *Basic {
	if b == nil {
		return nil
	}

	v, _ := cache.Load(b)
	cached, _ := v.(*Basic)
	if cached != nil {
		return cached
	}

	nb := &Basic{
		TypesBasic: b,
		Info:       b.Info(),
		Kind:       b.Kind(),
		Name:       b.Name(),
	}

	cache.Store(b, nb)
	return nb
}

func ToBasic(t interface{}) *Basic {
	switch t := t.(type) {
	case *Type:
		return t.Basic()
	case types.Type:
		return NewType(t).Basic()
	}
	return nil
}

func (b *Basic) String() string {
	return b.TypesBasic.String()
}

type Interface struct {
	TypesInterface      *types.Interface
	Empty               bool
	Embeddeds           []*Type
	Methods             map[string]*Func
	MethodNames         []string
	ExplicitMethods     map[string]*Func
	ExplicitMethodNames []string
}

var _ fmt.Stringer = (*Interface)(nil)

func NewInterface(iface *types.Interface) *Interface {
	if iface == nil {
		return nil
	}

	v, _ := cache.Load(iface)
	cached, _ := v.(*Interface)
	if cached != nil {
		return cached
	}

	var ni Interface
	cache.Store(iface, &ni)
	ni.TypesInterface = iface
	ni.Empty = iface.Empty()
	ni.Embeddeds = make([]*Type, iface.NumEmbeddeds())
	ni.Methods = make(map[string]*Func, iface.NumMethods())
	ni.MethodNames = make([]string, iface.NumMethods())
	ni.ExplicitMethods = make(map[string]*Func, iface.NumExplicitMethods())
	ni.ExplicitMethodNames = make([]string, iface.NumExplicitMethods())

	for i := 0; i < iface.NumEmbeddeds(); i++ {
		ni.Embeddeds[i] = NewType(iface.EmbeddedType(i))
	}

	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)
		ni.Methods[m.Name()] = NewFunc(m)
		ni.MethodNames[i] = m.Name()
	}

	for i := 0; i < iface.NumExplicitMethods(); i++ {
		m := iface.ExplicitMethod(i)
		ni.ExplicitMethods[m.Name()] = NewFunc(m)
		ni.ExplicitMethodNames[i] = m.Name()
	}

	return &ni
}

func ToInterface(t interface{}) *Interface {
	switch t := t.(type) {
	case *Type:
		return t.Interface()
	case types.Type:
		return NewType(t).Interface()
	}
	return nil
}

func (i *Interface) String() string {
	return i.TypesInterface.String()
}

type Signature struct {
	TypesSignature *types.Signature
	Recv           *Var
	Params         []*Var
	Results        []*Var
	Variadic       bool
}

var _ fmt.Stringer = (*Signature)(nil)

func NewSignature(s *types.Signature) *Signature {
	if s == nil {
		return nil
	}

	v, _ := cache.Load(s)
	cached, _ := v.(*Signature)
	if cached != nil {
		return cached
	}

	var ns Signature
	cache.Store(s, ns)
	ns.TypesSignature = s
	ns.Recv = NewVar(s.Recv())
	ns.Params = make([]*Var, s.Params().Len())
	ns.Results = make([]*Var, s.Results().Len())
	ns.Variadic = s.Variadic()

	for i := 0; i < s.Params().Len(); i++ {
		ns.Params[i] = NewVar(s.Params().At(i))
	}

	for i := 0; i < s.Results().Len(); i++ {
		ns.Results[i] = NewVar(s.Results().At(i))
	}

	return &ns
}

func ToSignature(t interface{}) *Signature {
	switch t := t.(type) {
	case *Type:
		return t.Signature()
	case types.Type:
		return NewType(t).Signature()
	}
	return nil
}

func (s *Signature) String() string {
	return s.TypesSignature.String()
}

type Named struct {
	TypesNamed  *types.Named
	Methods     map[string]*Func
	MethodNames []string
	Object      *TypeName
}

var _ fmt.Stringer = (*Named)(nil)

func NewNamed(n *types.Named) *Named {
	if n == nil {
		return nil
	}

	v, _ := cache.Load(n)
	cached, _ := v.(*Named)
	if cached != nil {
		return cached
	}

	var nn Named
	cache.Store(n, &nn)
	nn.TypesNamed = n
	nn.Methods = make(map[string]*Func, n.NumMethods())
	nn.MethodNames = make([]string, n.NumMethods())
	nn.Object = NewTypeName(n.Obj())

	for i := 0; i < n.NumMethods(); i++ {
		m := n.Method(i)
		nn.Methods[m.Name()] = NewFunc(m)
		nn.MethodNames[i] = m.Name()
	}

	return &nn
}

func ToNamed(t interface{}) *Named {
	switch t := t.(type) {
	case *Type:
		return t.Named()
	case types.Type:
		return NewType(t).Named()
	}
	return nil
}

func (c *Named) String() string {
	return c.TypesNamed.String()
}

func Methods(v interface{}) map[string]*Func {
	methods := map[string]*Func{}
	switch t := v.(type) {
	case *Type:
		return Methods(t.TypesType)
	case *TypeName:
		return Methods(t.TypesTypeName.Type())
	case *types.TypeName:
		return Methods(t.Type())
	case types.Type:
		ms := types.NewMethodSet(t)
		for i := 0; i < ms.Len(); i++ {
			m, _ := ms.At(i).Obj().(*types.Func)
			if m != nil {
				methods[m.Name()] = NewFunc(m)
			}
		}
		if _, isPtr := t.(*types.Pointer); !isPtr {
			ptrMethods := Methods(types.NewPointer(t))
			for n, m := range ptrMethods {
				if _, ok := methods[n]; !ok {
					methods[n] = m
				}
			}
		}
	}
	return methods
}

func under(t types.Type) types.Type {
	if named, _ := t.(*types.Named); named != nil {
		return under(named.Underlying())
	}
	return t
}


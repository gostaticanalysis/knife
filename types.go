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
	n, _ := t.TypesType.(*types.Named)
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
	na.TypesArray = a
	na.Elem = NewType(a.Elem())
	na.Len = a.Len()
	cache.Store(a, &na)
	return &na
}

func ToArray(t any) *Array {
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
	ns.TypesSlice = s
	ns.Elem = NewType(s.Elem())

	cache.Store(s, &ns)
	return &ns
}

func ToSlice(t any) *Slice {
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

func ToStruct(t any) *Struct {
	switch t := t.(type) {
	case *Type:
		return t.Struct()
	case *TypeName:
		return t.Type.Struct()
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
	nm.TypesMap = m
	nm.Elem = NewType(m.Elem())
	nm.Key = NewType(m.Key())

	cache.Store(m, &nm)
	return &nm
}

func ToMap(t any) *Map {
	switch t := t.(type) {
	case *Type:
		return t.Map()
	case *TypeName:
		return t.Type.Map()
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
	np.TypesPointer = p
	np.Elem = NewType(p.Elem())

	cache.Store(p, &np)
	return &np
}

func ToPointer(t any) *Pointer {
	switch t := t.(type) {
	case *Type:
		return t.Pointer()
	case *TypeName:
		return t.Type.Pointer()
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
	nc.TypesChan = c
	nc.Dir = c.Dir()
	nc.Elem = NewType(c.Elem())

	cache.Store(c, &nc)
	return &nc
}

func ToChan(t any) *Chan {
	switch t := t.(type) {
	case *Type:
		return t.Chan()
	case *TypeName:
		return t.Type.Chan()
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

func ToBasic(t any) *Basic {
	switch t := t.(type) {
	case *Type:
		return t.Basic()
	case *TypeName:
		return t.Type.Basic()
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

func ToInterface(t any) *Interface {
	switch t := t.(type) {
	case *Type:
		return t.Interface()
	case *TypeName:
		return t.Type.Interface()
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

func ToSignature(t any) *Signature {
	switch t := t.(type) {
	case *Type:
		return t.Signature()
	case *TypeName:
		return t.Type.Signature()
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

func ToNamed(t any) *Named {
	switch t := t.(type) {
	case *Type:
		return t.Named()
	case *TypeName:
		return t.Type.Named()
	case types.Type:
		return NewType(t).Named()
	}
	return nil
}

func (c *Named) String() string {
	return c.TypesNamed.String()
}

func Methods(v any) map[string]*Func {
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
	if t == nil {
		return nil
	}
	return t.Underlying()
}

func implements(t any, iface any) bool {
	if t == nil || iface == nil {
		return false
	}

	var (
		_t     types.Type
		_iface *types.Interface
	)

	switch t := t.(type) {
	case types.Type:
		if t == nil {
			return false
		}
		_t = t
	case *Type:
		if t == nil {
			return false
		}
		_t = t.TypesType
	case *TypeName:
		if t == nil || t.Type == nil {
			return false
		}
		_t = t.Type.TypesType
	case Object:
		if t == nil {
			return false
		}
		_t = t.TypesObject().Type()
	case types.Object:
		if t == nil {
			return false
		}
		_t = t.Type()
	default:
		return false
	}

	switch iface := iface.(type) {
	case *types.Interface:
		_iface = iface
	case types.Type:
		_iface, _ = under(iface).(*types.Interface)
	case *Type:
		_iface, _ = under(iface.TypesType.Underlying()).(*types.Interface)
	case Object:
		_iface, _ = under(iface.TypesObject().Type().Underlying()).(*types.Interface)
	case types.Object:
		_iface, _ = under(iface.Type()).(*types.Interface)
	case *Interface:
		_iface = iface.TypesInterface
	default:
		return false
	}

	if _t == nil || _iface == nil {
		return false
	}

	return types.Implements(_t, _iface) || types.Implements(types.NewPointer(_t), _iface)
}

func identical(t1, t2 any) bool {
	if t1 == nil || t2 == nil {
		return false
	}

	var _t1, _t2 types.Type

	switch t1 := t1.(type) {
	case types.Type:
		if t1 == nil {
			return false
		}
		_t1 = t1
	case *Type:
		if t1 == nil {
			return false
		}
		_t1 = t1.TypesType
	case Object:
		if t1 == nil {
			return false
		}
		_t1 = t1.TypesObject().Type()
	case types.Object:
		if t1 == nil {
			return false
		}
		_t1 = t1.Type()
	}

	switch t2 := t2.(type) {
	case types.Type:
		if t2 == nil {
			return false
		}
		_t2 = t2
	case *Type:
		if t2 == nil {
			return false
		}
		_t2 = t2.TypesType
	case Object:
		if t2 == nil {
			return false
		}
		_t2 = t2.TypesObject().Type()
	case types.Object:
		if t2 == nil {
			return false
		}
		_t2 = t2.Type()
	}

	return _t1 != nil && _t2 != nil && types.Identical(_t1, _t2)
}

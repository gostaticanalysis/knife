package knife

import (
	"bytes"
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"strings"
	"text/template"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/comment"
	"golang.org/x/tools/go/packages"
)

// newTemplate creates new a template with funcmap.
func newTemplate(pkg *packages.Package, extraData map[string]interface{}) *template.Template {
	prefix := pkg.Name
	if prefix == "" {
		prefix = pkg.ID
	}
	return template.New(prefix + "_format").Funcs(newFuncMap(pkg, extraData))
}

func newFuncMap(pkg *packages.Package, extraData map[string]interface{}) template.FuncMap {
	var cmaps comment.Maps
	return template.FuncMap{
		"pkg":        func() *Package { return NewPackage(pkg.Types) },
		"br":         fmt.Sprintln,
		"array":      ToArray,
		"basic":      ToBasic,
		"chan":       ToChan,
		"interface":  ToInterface,
		"map":        ToMap,
		"named":      ToNamed,
		"pointer":    ToPointer,
		"ptr":        ToPointer,
		"signature":  ToSignature,
		"slice":      ToSlice,
		"struct":     ToStruct,
		"len":        func(v interface{}) int { return reflect.ValueOf(v).Len() },
		"cap":        func(v interface{}) int { return reflect.ValueOf(v).Cap() },
		"exported":   Exported,
		"methods":    Methods,
		"names":      names,
		"implements": implements,
		"identical":  identical,
		"under":      under,
		"pos":        func(v interface{}) token.Position { return Position(pkg.Fset, v) },
		"objectof":   func(s string) Object { return objectOf(pkg.Types, s) },
		"typeof":     func(s string) *Type { return typeOf(pkg.Types, s) },
		"doc":        func(v interface{}) string { return doc(pkg, cmaps, v) },
		"data":       func(k string) interface{} { return extraData[k] },
		"at":         at,
	}
}

func names(slice interface{}) string {
	vs := reflect.ValueOf(slice)
	switch vs.Kind() {
	case reflect.Slice, reflect.Array:
		return nameSlice(vs)
	case reflect.Map:
		return nameMap(vs)
	}

	return ""
}

func nameSlice(vs reflect.Value) string {
	var buf bytes.Buffer
	for i := 0; i < vs.Len(); i++ {
		s := name(vs.Index(i))
		if s != "" {
			fmt.Fprintln(&buf, s)
		}
	}
	return strings.TrimRight(buf.String(), "\n")
}

func nameMap(vs reflect.Value) string {
	var buf bytes.Buffer
	for _, k := range vs.MapKeys() {
		s := name(vs.MapIndex(k))
		if s != "" {
			fmt.Fprintln(&buf, s)
		}
	}
	return strings.TrimRight(buf.String(), "\n")
}

func name(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Ptr:
		return name(v.Elem())
	case reflect.Struct:
		fv := v.FieldByName("Name")
		if !fv.IsZero() {
			return fv.String()
		}
	}
	return ""
}

func objectOf(typesPkg *types.Package, s string) Object {
	ss := strings.Split(s, ".")

	switch len(ss) {
	case 1:
		obj := types.Universe.Lookup(s)
		return NewObject(obj)
	case 2:
		pkg, name := ss[0], ss[1]
		obj := analysisutil.LookupFromImports(typesPkg.Imports(), pkg, name)
		if obj != nil {
			return NewObject(obj)
		}
		if analysisutil.RemoveVendor(typesPkg.Name()) != analysisutil.RemoveVendor(pkg) {
			return nil
		}
		return NewObject(typesPkg.Scope().Lookup(name))
	}
	return nil
}

func typeOf(typesPkg *types.Package, s string) *Type {
	if s == "" {
		return nil
	}

	if s[0] == '*' {
		typ := typeOf(typesPkg, s[1:])
		if typ == nil {
			return nil
		}
		return NewType(types.NewPointer(typ.TypesType))
	}

	obj := objectOf(typesPkg, s)
	if obj == nil {
		return nil
	}
	return NewType(obj.TypesObject().Type())
}

func doc(pkg *packages.Package, cmaps comment.Maps, v interface{}) string {
	node, ok := v.(interface{ Pos() token.Pos })
	if !ok {
		return ""
	}

	if cmaps == nil {
		cmaps = comment.New(pkg.Fset, pkg.Syntax)
	}

	pos := node.Pos()
	cgs := cmaps.CommentsByPosLine(pkg.Fset, pos)
	if len(cgs) > 0 {
		return strings.TrimSpace(cgs[len(cgs)-1].Text())
	}

	return ""
}

func at(v interface{}, expr string) interface{} {
	ret, err := At(v, expr)
	if err != nil {
		panic(err)
	}
	return ret
}

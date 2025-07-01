package knife

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os/exec"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/comment"
)

type TempalteData struct {
	Fset      *token.FileSet
	Files     []*ast.File
	TypesInfo *types.Info
	Pkg       *types.Package
	Extra     map[string]any
}

// NewTemplate creates new a template with funcmap.
func NewTemplate(td *TempalteData) *template.Template {
	prefix := td.Pkg.Name()
	return template.New(prefix + "_format").Funcs(newFuncMap(td))
}

func newFuncMap(td *TempalteData) template.FuncMap {
	var cmaps comment.Maps
	return template.FuncMap{
		"pkg":        func() *Package { return NewPackage(td.Pkg) },
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
		"len":        lenFunc,
		"cap":        capFunc,
		"last":       lastFunc,
		"exported":   Exported,
		"methods":    Methods,
		"names":      td.names,
		"implements": implements,
		"identical":  identical,
		"under":      func(v any) *Type { return NewType(under(v)) },
		"pos":        func(v any) token.Position { return Position(td.Fset, v) },
		"objectof":   func(s string) Object { return td.objectOf(s) },
		"typeof":     func(s string) *Type { return td.typeOf(s) },
		"doc":        func(v any) string { return td.doc(cmaps, v) },
		"data":       func(k string) any { return td.Extra[k] },
		"regexp":     regexpMatch,
		"godoc":      godoc,
	}
}

func (td *TempalteData) names(slice any) string {
	vs := reflect.ValueOf(slice)
	switch vs.Kind() {
	case reflect.Slice, reflect.Array:
		return td.nameSlice(vs)
	case reflect.Map:
		return td.nameMap(vs)
	}

	return ""
}

func (td *TempalteData) nameSlice(vs reflect.Value) string {
	var buf bytes.Buffer
	for i := 0; i < vs.Len(); i++ {
		s := td.name(vs.Index(i))
		if s != "" {
			fmt.Fprintln(&buf, s)
		}
	}
	return strings.TrimRight(buf.String(), "\n")
}

func (td *TempalteData) nameMap(vs reflect.Value) string {
	var buf bytes.Buffer
	for _, k := range vs.MapKeys() {
		s := td.name(vs.MapIndex(k))
		if s != "" {
			fmt.Fprintln(&buf, s)
		}
	}
	return strings.TrimRight(buf.String(), "\n")
}

func (td *TempalteData) name(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return ""
		}
		return td.name(v.Elem())
	case reflect.Struct:
		fv := v.FieldByName("Name")
		if !fv.IsValid() || fv.IsZero() {
			return ""
		}
		return fv.String()
	}
	return ""
}

func (td *TempalteData) objectOf(s string) Object {
	dotPos := strings.LastIndex(s, ".")

	if dotPos == -1 {
		obj := types.Universe.Lookup(s)
		return NewObject(obj)
	}

	pkg, name := s[:dotPos], s[dotPos+1:]
	obj := analysisutil.LookupFromImports(td.Pkg.Imports(), pkg, name)
	if obj != nil {
		return NewObject(obj)
	}

	if analysisutil.RemoveVendor(td.Pkg.Name()) != analysisutil.RemoveVendor(pkg) {
		return nil
	}

	return NewObject(td.Pkg.Scope().Lookup(name))
}

func (td *TempalteData) typeOf(s string) *Type {
	if s == "" {
		return nil
	}

	if s[0] == '*' {
		typ := td.typeOf(s[1:])
		if typ == nil {
			return nil
		}
		return NewType(types.NewPointer(typ.TypesType))
	}

	obj := td.objectOf(s)
	if obj == nil {
		return nil
	}
	return NewType(obj.TypesObject().Type())
}

func (td *TempalteData) doc(cmaps comment.Maps, v any) string {
	if v == nil {
		return ""
	}
	node, ok := v.(interface{ Pos() token.Pos })
	if !ok {
		return ""
	}

	if cmaps == nil {
		cmaps = comment.New(td.Fset, td.Files)
	}

	pos := node.Pos()
	cgs := cmaps.CommentsByPosLine(td.Fset, pos)
	if len(cgs) > 0 {
		return strings.TrimSpace(cgs[len(cgs)-1].Text())
	}

	return ""
}

// regexpMatch performs regular expression matching.
// Usage: regexp "pattern" "text" - returns true if pattern matches text
func regexpMatch(pattern, text string) (bool, error) {
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false, fmt.Errorf("regexp error: %v", err)
	}
	return matched, nil
}

// godoc executes go doc command with arguments and returns the output.
// Usage: godoc "fmt.Println" or godoc "fmt"
func godoc(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("go doc: no arguments provided")
	}

	cmdArgs := slices.Concat([]string{"doc"}, args)
	cmd := exec.Command("go", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("go doc command failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("go doc command failed: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func lenFunc(v any) (int, error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len(), nil
	default:
		return 0, fmt.Errorf("len: invalid type %T", v)
	}
}

func capFunc(v any) (int, error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		return rv.Cap(), nil
	default:
		return 0, fmt.Errorf("cap: invalid type %T", v)
	}
}

func lastFunc(v any) (any, error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		if rv.Len() == 0 {
			return nil, fmt.Errorf("last: empty collection")
		}
		return rv.Index(rv.Len() - 1).Interface(), nil
	case reflect.String:
		if rv.Len() == 0 {
			return "", fmt.Errorf("last: empty string")
		}
		return string(rv.String()[rv.Len()-1]), nil
	default:
		return nil, fmt.Errorf("last: invalid type %T", v)
	}
}

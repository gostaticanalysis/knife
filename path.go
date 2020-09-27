package knife

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
)

// At access a field of v by a path.
// v must be struct or pointer of struct.
// A path is represented by Go's expression which can be parsed by go/parser.ParseExpr.
// You can use selectors and indexes in a path.
// Slice and arrays index allow only expressions of int.
// Maps key allow only expressions of string, int and float64.
// If a map key is string, you can use Map.Key instead of Map["Key"].
// If a map key does not exist, At returns nil.
func At(v interface{}, expr string) (interface{}, error) {
	p, err := NewPath(expr)
	if err != nil {
		return nil, err
	}

	var dest interface{}
	if err := p.Eval(v, &dest); err != nil {
		return nil, err
	}
	
	return dest, nil
}

// Path represents a path to a field or an element of slice, key and map.
type Path struct {
	expr ast.Expr
}

// NewPath creates new Path.
func NewPath(expr string) (*Path, error) {

	_expr, err := parser.ParseExpr("v." + expr)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	return &Path{expr: _expr}, nil
}

// Eval evaluates and traces the path from root.
// Eval sets the value into dest.
func (p *Path) Eval(root, dest interface{}) error {
	rootV := reflect.ValueOf(root)
	destV := reflect.ValueOf(dest)
	if destV.Kind() != reflect.Ptr {
		return errors.New("dest must be pointer")
	}

	v, err := p.evalExpr(rootV, p.expr)
	if err != nil {
		return err
	}

	d := destV.Elem()
	if d.CanSet() && v.IsValid() && v.Type().AssignableTo(d.Type()) {
		d.Set(v)
	}

	return nil
}

func (p *Path) direct(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.Elem()
	default:
		return v
	}
}

func (p *Path) evalExpr(v reflect.Value, expr ast.Expr) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return p.evalExpr(v.Elem(), expr)
	}

	switch expr := expr.(type) {
	case *ast.Ident:
		return v, nil
	case *ast.SelectorExpr:
		return p.evalSelectorExpr(v, expr)
	case *ast.IndexExpr:
		return p.evalIndexExpr(v, expr)
	default:
		var buf bytes.Buffer
		types.WriteExpr(&buf, expr)
		return reflect.Value{}, fmt.Errorf("does not support expr: %s", &buf)
	}
}

func (p *Path) evalSelectorExpr(v reflect.Value, expr *ast.SelectorExpr) (reflect.Value, error) {
	ev, err := p.evalExpr(v, expr.X)
	if err != nil {
		return reflect.Value{}, err
	}

	ev = p.direct(ev)
	switch ev.Kind() {
	case reflect.Struct:
		name := expr.Sel.Name
		fv := ev.FieldByName(name)
		if !fv.IsValid() {
			return reflect.Value{}, fmt.Errorf("does find field: %s", name)
		}
		return fv, nil
	case reflect.Map:
		keyKind := ev.Type().Key().Kind()
		switch keyKind {
		case reflect.String:
			key := reflect.ValueOf(expr.Sel.Name)
			return ev.MapIndex(key), nil
		default:
			return reflect.Value{}, errors.New("does not support selector type")
		}
	default:
		return reflect.Value{}, errors.New("does not support selector type")
	}
}

func (p *Path) evalIndexExpr(v reflect.Value, expr *ast.IndexExpr) (reflect.Value, error) {
	ev, err := p.evalExpr(v, expr.X)
	if err != nil {
		return reflect.Value{}, err
	}
	ev = p.direct(ev)

	idx, err := p.evalIndex(expr.Index)
	if idx == nil {
		var buf bytes.Buffer
		types.WriteExpr(&buf, expr.Index)
		return reflect.Value{}, fmt.Errorf("does not support indexer: %s", &buf)
	}

	switch ev.Kind() {
	case reflect.Slice, reflect.Array:
		i, ok := constant.Int64Val(idx)
		if !ok {
			return reflect.Value{}, errors.New("an index of slice or array must be integer")
		}
		return ev.Index(int(i)), nil
	case reflect.Map:
		keyKind := ev.Type().Key().Kind()
		switch keyKind {
		case reflect.Int:
			key, ok := constant.Int64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(int(key))), nil
		case reflect.Int8:
			key, ok := constant.Int64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(int8(key))), nil
		case reflect.Int16:
			key, ok := constant.Int64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(int16(key))), nil
		case reflect.Int32:
			key, ok := constant.Int64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(int32(key))), nil
		case reflect.Int64:
			key, ok := constant.Int64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(key)), nil
		case reflect.Uint:
			key, ok := constant.Uint64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(uint(key))), nil
		case reflect.Uint8:
			key, ok := constant.Uint64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(uint8(key))), nil
		case reflect.Uint16:
			key, ok := constant.Uint64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(uint16(key))), nil
		case reflect.Uint32:
			key, ok := constant.Uint64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(uint32(key))), nil
		case reflect.Uint64:
			key, ok := constant.Uint64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(key)), nil
		case reflect.Float32:
			key, ok := constant.Float64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(float32(key))), nil
		case reflect.Float64:
			key, ok := constant.Float64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not match a key of map and expr")
			}
			return ev.MapIndex(reflect.ValueOf(key)), nil
		case reflect.String:
			key := constant.StringVal(idx)
			return ev.MapIndex(reflect.ValueOf(key)), nil
		case reflect.Bool:
			key := constant.BoolVal(idx)
			return ev.MapIndex(reflect.ValueOf(key)), nil
		default:
			return reflect.Value{}, fmt.Errorf("does not support key type: %s", keyKind)
		}
	default:
		return reflect.Value{}, errors.New("does not support expr type")
	}
}

func (p *Path) evalIndex(expr ast.Expr) (constant.Value, error) {
	info := &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{},
	}

	fset := token.NewFileSet()
	if err := types.CheckExpr(fset, nil, token.NoPos, expr, info); err != nil {
		return nil, fmt.Errorf("type check error: %w", err)
	}

	return info.Types[expr].Value, nil
}

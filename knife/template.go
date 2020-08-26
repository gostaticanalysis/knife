package knife

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

var Template = template.New("knife_format").Funcs(template.FuncMap{
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
})

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

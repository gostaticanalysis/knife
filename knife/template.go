package knife

import (
	"fmt"
	"reflect"
	"text/template"
)

var Template = template.New("knife_format").Funcs(template.FuncMap{
	"br":        fmt.Sprintln,
	"array":     ToArray,
	"basic":     ToBasic,
	"chan":      ToChan,
	"interface": ToInterface,
	"map":       ToMap,
	"named":     ToNamed,
	"pointer":   ToPointer,
	"ptr":       ToPointer,
	"signature": ToSignature,
	"slice":     ToSlice,
	"struct":    ToStruct,
	"len":       func(v interface{}) int { return reflect.ValueOf(v).Len() },
	"cap":       func(v interface{}) int { return reflect.ValueOf(v).Cap() },
	"exported":  Exported,
	"methods":   Methods,
})

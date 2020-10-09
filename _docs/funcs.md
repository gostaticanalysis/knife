# Functions for a template

knife provides functions which can be use in a template.

| Function | Example | Description |
| - | - | - |
| `pkg` | `{{pkg}}` | target package |
| `br` | `{{br}}` | new line |
| `array` | `{{(array .).Len}}` | convert type to [`knife.Array`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Array)<br>see: [knife.ToArray](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToArray) |
| `basic` | `{{(basic .).Kind}}` | convert type to [`knife.Basic`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Basic)<br>see: [knife.ToBasic](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToBasic) |
| `chan` | `{{(chan .).Dir}}` | convert type to [`knife.Chan`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Chan)<br>see: [knife.ToChan](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToChan) |
| `interface` | `{{(interface .).Methods}}` | convert type to [`knife.Interface`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Interface)<br>see: [knife.ToInterface](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToInterface) |
| `map` | `{{(map .).Key}}` | convert type to [`knife.Map`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Map)<br>see: [knife.ToMap](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToMap) |
| `named` | `{{(named .).Methods}}` | convert type to [`knife.Named`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Named)<br>see: [knife.ToNamed](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToNamed) |
| `pointer` | `{{(pointer .).Elem}}` | convert type to [`knife.Pointer`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Pointer)<br>see: [knife.ToPointer](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToPointer) |
| `ptr` | `{{(ptr .).Elem}}` | same as `pointer` |
| `signature` | `{{(signature .).Recv}}` | convert type to [`knife.Signature`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Signature)<br>see: [knife.ToSignature](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToSignature) |
| `slice` | `{{(slice .).Elem}}` | convert type to [`knife.Slice`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Slice)<br>see: [knife.ToSlice](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToSlice) |
| `struct` | `{{(struct .).Fields}}` | convert type to [`knife.Struct`](https://pkg.go.dev/github.com/gostaticanalysis/knife/Struct)<br>see: [knife.ToStruct](https://pkg.go.dev/github.com/gostaticanalysis/knife/ToStruct) |
| `len` | `{{len .}}` | `len(x)` calls `reflect.ValueOf(x).Len()` |
| `cap` | `{{cap .}}` | `cap(x)` calls `reflect.ValueOf(x).Cap()` |
| `last` | `{{last .}}` | `last(x)` returns last element of a slice, array or string |
| `exported` | `{{exported .Types}}` | `exported` filters out unexported objects |
| `methods` | `{{methods .Types.T}}` | `methods` returns methods of the type |
| `names` | `{{range names .Types}}{{.}}{{end}}` | slice, array or map of `Name` field |
| `implements` | `{{if implements . (typeof "error")}}{{.}}{{end}}` | `implements` reports whether the type implements the interface |
| `identical` | `{{if identical . (typeof "error")}}{{.}}{{end}}` | `identical` reports whether the two types are identical types |
| `under` | `{{under .Types.T}}` | `under` returns an underlying type of the type recursively |
| `pos` | `{{pos .}}` | `pos` returns `token.Position` by calling `Pos` methods |
| `objectof` | `{{objectof "panic"}}` | `objectof` returns `knife.Object` which is specified by name |
| `typeof` | `{{typeof "error"}}` | `typeof` returns `*knife.Type` which is specified by name |
| `doc` | `{{doc .Types.T}}` | `doc` returns corresponding document to the object |
| `data` | `{{data "key"}}` | `data` returns extra data which given via `knife.Option` |

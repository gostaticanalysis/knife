# Template Format Reference

This document describes the template syntax and functions available for formatting Go package information.

## Overview

Templates use Go's `text/template` package with additional custom functions for inspecting Go packages. The template system allows you to extract and format detailed information about types, functions, variables, and constants from Go packages.

## Template Context

The root context (`.`) is a `Package` object that provides access to:

- `.Types` - Map of type names to TypeName objects (`map[string]*TypeName`)
- `.Funcs` - Map of function names to Func objects (`map[string]*Func`)
- `.Vars` - Map of variable names to Var objects (`map[string]*Var`)
- `.Consts` - Map of constant names to Const objects (`map[string]*Const`)
- `.Name` - Package name (string)
- `.Path` - Package path (string)
- `.Imports` - Imported packages (`[]*Package`)

## Available Types and Properties

### Object Types

#### Function (`*Func`)
```go
.Name         // Function name (string)
.Exported     // Whether exported (bool)  
.Package      // Containing package (*Package)
.Signature    // Function signature (*Signature)
.Pos()        // Position (token.Position)
```

#### Variable (`*Var`) 
```go
.Name         // Variable name (string)
.Exported     // Whether exported (bool)
.Type         // Variable type (*Type)
.Package      // Containing package (*Package)
.Pos()        // Position (token.Position)
```

#### Constant (`*Const`)
```go
.Name         // Constant name (string)
.Exported     // Whether exported (bool)
.Type         // Constant type (*Type)
.Value        // Constant value (constant.Value)
.Package      // Containing package (*Package)
.BoolVal()    // Extract as bool
.StringVal()  // Extract as string
.Int64Val()   // Extract as int64
.Float64Val() // Extract as float64
.Pos()        // Position (token.Position)
```

#### Type Name (`*TypeName`)
```go
.Name         // Type name (string)
.Exported     // Whether exported (bool)
.IsAlias      // Whether type alias (bool)
.Type         // Actual type (*Type)
.Package      // Containing package (*Package)
.Pos()        // Position (token.Position)
```

### Type System

#### Core Type (`*Type`)
```go
.Underlying() // Underlying type (*Type)
.Array()      // Convert to *Array (if applicable)
.Slice()      // Convert to *Slice (if applicable)
.Struct()     // Convert to *Struct (if applicable)
.Map()        // Convert to *Map (if applicable)
.Pointer()    // Convert to *Pointer (if applicable)
.Chan()       // Convert to *Chan (if applicable)
.Basic()      // Convert to *Basic (if applicable)
.Interface()  // Convert to *Interface (if applicable)
.Signature()  // Convert to *Signature (if applicable)
.Named()      // Convert to *Named (if applicable)
```

#### Composite Types

**Array (`*Array`)**
```go
.Elem         // Element type (*Type)
.Len          // Array length (int64)
```

**Slice (`*Slice`)**
```go
.Elem         // Element type (*Type)
```

**Struct (`*Struct`)**
```go
.Fields       // Map of fields (map[string]*Field)
.FieldNames   // Field names ([]string)
```

**Map (`*Map`)**
```go
.Key          // Key type (*Type)
.Elem         // Value type (*Type)
```

**Pointer (`*Pointer`)**
```go
.Elem         // Pointed-to type (*Type)
```

**Channel (`*Chan`)**
```go
.Dir          // Channel direction (types.ChanDir)
.Elem         // Element type (*Type)
```

**Interface (`*Interface`)**
```go
.Empty        // Whether empty interface (bool)
.Methods      // All methods (map[string]*Func)
.MethodNames  // Method names ([]string)
.ExplicitMethods // Declared methods (map[string]*Func)
.Embeddeds    // Embedded types ([]*Type)
```

**Named Type (`*Named`)**
```go
.Methods      // Type methods (map[string]*Func)
.MethodNames  // Method names ([]string)
.Object       // Type name object (*TypeName)
```

**Function Signature (`*Signature`)**
```go
.Recv         // Receiver (*Var)
.Params       // Parameters ([]*Var)
.Results      // Return values ([]*Var)
.Variadic     // Whether variadic (bool)
```

**Basic Type (`*Basic`)**
```go
.Kind         // Basic type kind (types.BasicKind)
.Info         // Type info (types.BasicInfo)
.Name         // Type name (string)
```

**Struct Field (`*Field`)**
```go
.Name         // Field name (string)
.Type         // Field type (*Type)
.Tag          // Struct tag (string)
.Anonymous    // Whether anonymous (bool)
.Exported     // Whether exported (bool)
.Struct       // Containing struct (*Struct)
.Pos()        // Position (token.Position)
```

## Basic Template Examples

### Simple Output

```go
{{.}}  // Print everything (default)
```

### List Names

```go
{{range .Types}}{{.Name}}{{br}}{{end}}
```

## Template Functions

### Core Functions

| Function | Example | Description |
|----------|---------|-------------|
| `pkg` | `{{pkg}}` | Current package |
| `br` | `{{br}}` | Line break |
| `len` | `{{len .Types}}` | Length of slice/array/map |
| `cap` | `{{cap .}}` | Capacity of slice/array |
| `last` | `{{last .Types}}` | Last element of slice/array/string |

### Type Conversion Functions

| Function | Example | Description |
|----------|---------|-------------|
| `array` | `{{(array .).Len}}` | Convert to Array type |
| `basic` | `{{(basic .).Kind}}` | Convert to Basic type |
| `chan` | `{{(chan .).Dir}}` | Convert to Chan type |
| `interface` | `{{(interface .).Methods}}` | Convert to Interface type |
| `map` | `{{(map .).Key}}` | Convert to Map type |
| `named` | `{{(named .).Methods}}` | Convert to Named type |
| `pointer`/`ptr` | `{{(pointer .).Elem}}` | Convert to Pointer type |
| `signature` | `{{(signature .).Recv}}` | Convert to Signature type |
| `slice` | `{{(slice .).Elem}}` | Convert to Slice type |
| `struct` | `{{(struct .).Fields}}` | Convert to Struct type |

### Filtering and Analysis Functions

| Function | Example | Description |
|----------|---------|-------------|
| `exported` | `{{exported .Types}}` | Filter exported objects only |
| `methods` | `{{methods .Types.T}}` | Get methods of a type |
| `names` | `{{range names .Types}}{{.}}{{end}}` | Extract Name fields from slice/array/map |
| `implements` | `{{if implements . (typeof "error")}}{{.}}{{end}}` | Check if type implements interface |
| `identical` | `{{if identical . (typeof "error")}}{{.}}{{end}}` | Check if two types are identical |
| `under` | `{{under .Types.T}}` | Get underlying type recursively |

### Object and Type Lookup Functions

| Function | Example | Description |
|----------|---------|-------------|
| `objectof` | `{{objectof "panic"}}` | Get Object by name |
| `typeof` | `{{typeof "error"}}` | Get Type by name |
| `pos` | `{{pos .}}` | Get position (file:line) |
| `doc` | `{{doc .Types.T}}` | Get documentation comment |
| `data` | `{{data "key"}}` | Access extra data from `-data` flag |
| `regexp` | `{{regexp "^Get" .Name}}` | Check if text matches regex pattern |
| `godoc` | `{{godoc "fmt.Println"}}` | Execute go doc command and return output for the specified symbol or package |

### Documentation Functions

#### `godoc` Function

The `godoc` function executes the `go doc` command and returns the documentation for the specified package, type, function, or variable.

**Syntax:**
```go
{{godoc "package"}}           // Package documentation
{{godoc "package.Symbol"}}    // Symbol documentation
{{godoc "package.Type.Method"}} // Method documentation
{{godoc "-src" "package.Symbol"}} // Show source code instead of documentation
```

**Examples:**

**Package Documentation:**
```go
{{godoc "fmt"}}               // Documentation for fmt package
{{godoc "net/http"}}          // Documentation for net/http package
{{godoc .Path}}               // Documentation for current package
```

**Function Documentation:**
```go
{{godoc "fmt.Printf"}}        // Documentation for fmt.Printf
{{godoc "os.Open"}}           // Documentation for os.Open
{{godoc "-src" "fmt.Printf"}} // Source code for fmt.Printf
{{godoc "-src" "os.Open"}}    // Source code for os.Open
```

**Type Documentation:**
```go
{{godoc "http.Server"}}       // Documentation for http.Server type
{{godoc "io.Reader"}}         // Documentation for io.Reader interface
{{godoc "-src" "http.Server"}} // Source code for http.Server type
{{godoc "-src" "io.Reader"}}  // Source code for io.Reader interface
```

**Method Documentation:**
```go
{{godoc "http.Server.ListenAndServe"}}     // Documentation for Server.ListenAndServe method
{{godoc "-src" "http.Server.ListenAndServe"}} // Source code for Server.ListenAndServe method
```

**Dynamic Symbol Documentation:**
```go
{{range .Types}}{{godoc (printf "%s.%s" .Package.Path .Name)}}{{end}}
{{range .Funcs}}{{godoc (printf "%s.%s" .Package.Path .Name)}}{{end}}
```

**Using go doc Flags:**
```go
{{godoc "-src" "fmt.Printf"}}     // Show source code
{{godoc "-short" "fmt"}}          // Show package summary only
{{godoc "-u" "fmt"}}              // Show unexported symbols
{{godoc "-c" "fmt.Printf"}}       // Show examples
{{godoc "-all" "fmt"}}            // Show all symbols and methods
```

**Error Handling:**
If the `godoc` command fails or the symbol is not found, an empty string is returned.

**Performance Note:**
The `godoc` function executes external commands, so it may be slower than other template functions. Use it judiciously in loops.

## Template Patterns

### List All Exported Functions

```go
{{range exported .Funcs}}{{.Name}}{{br}}{{end}}
```

### Find Functions with Specific Signature

```go
{{range exported .Funcs}}{{if eq .Signature.Params.Len 1}}{{.Name}}{{br}}{{end}}{{end}}
```

### List Types Implementing an Interface

```go
{{range exported .Types}}{{if implements . (typeof "error")}}{{.Name}}{{br}}{{end}}{{end}}
```

### Show Type Positions

```go
{{range .Types}}{{.Name}} - {{pos .}}{{br}}{{end}}
```

### Access Struct Fields

```go
{{range .Types}}{{$t := .}}{{with struct .}}{{range .Fields}}{{$t.Name}}.{{.Name}}{{br}}{{end}}{{end}}{{end}}
```

### Filter Functions with Name Patterns

```go
{{range exported .Funcs}}{{if regexp "^Get" .Name}}{{.Name}}{{br}}{{end}}{{end}}
```

### Find Types with Specific Naming Convention

```go
{{range exported .Types}}{{if regexp "Interface$" .Name}}{{.Name}}{{br}}{{end}}{{end}}
```

### Filter Test Functions

```go
{{range .Funcs}}{{if regexp "^Test" .Name}}{{.Name}}{{br}}{{end}}{{end}}
```

### Find Error Variables

```go
{{range exported .Vars}}{{if implements .Type (typeof "error")}}{{.Name}}{{br}}{{end}}{{end}}
```

### List Interface Methods

```go
{{range exported .Types}}{{$t := .}}{{with interface .Type}}{{range .Methods}}{{$t.Name}}.{{.Name}}{{br}}{{end}}{{end}}{{end}}
```

### Find HTTP Handler Functions

```go
{{range exported .Funcs}}{{with .Signature}}{{if and (eq (len .Params) 2) (identical (index .Params 0).Type (typeof "net/http.ResponseWriter")) (identical (index .Params 1).Type (typeof "*net/http.Request"))}}{{.Recv.Name}}{{br}}{{end}}{{end}}{{end}}
```

### List Struct Fields with JSON Tags

```go
{{range exported .Types}}{{$t := .}}{{with struct .Type}}{{range .Fields}}{{if .Tag}}{{$t.Name}}.{{.Name}} `{{.Tag}}`{{br}}{{end}}{{end}}{{end}}{{end}}
```

### Find Context-Accepting Functions

```go
{{range exported .Funcs}}{{with .Signature}}{{if and (gt (len .Params) 0) (identical (index .Params 0).Type (typeof "context.Context"))}}{{.Recv.Name}}{{br}}{{end}}{{end}}{{end}}
```

### List Constants by Type

```go
{{range exported .Consts}}{{if identical .Type (typeof "string")}}{{.Name}} = {{.StringVal}}{{br}}{{end}}{{end}}
```

### Find Embedded Interfaces

```go
{{range exported .Types}}{{$t := .}}{{with interface .Type}}{{range .Embeddeds}}{{$t.Name}} embeds {{.}}{{br}}{{end}}{{end}}{{end}}
```

### List Method Receivers

```go
{{range exported .Types}}{{$t := .}}{{with named .Type}}{{range .Methods}}{{.Signature.Recv.Type}} {{.Name}}{{br}}{{end}}{{end}}{{end}}
```

### Find Generic Types

```go
{{range exported .Types}}{{with named .Type}}{{if .Object.IsAlias}}{{.Object.Name}} (alias){{br}}{{end}}{{end}}{{end}}
```

### List Variadic Functions

```go
{{range exported .Funcs}}{{with .Signature}}{{if .Variadic}}{{.Recv.Name}}{{br}}{{end}}{{end}}{{end}}
```

### Find Channel Types

```go
{{range exported .Types}}{{$t := .}}{{with chan .Type}}{{$t.Name}} chan {{.Elem}}{{br}}{{end}}{{end}}
```

### List Package Dependencies

```go
{{range .Imports}}{{.Path}}{{br}}{{end}}
```

### Find Types Implementing Multiple Interfaces

```go
{{range exported .Types}}{{if and (implements .Type (typeof "io.Reader")) (implements .Type (typeof "io.Writer"))}}{{.Name}}{{br}}{{end}}{{end}}
```

### Show Function Signatures

```go
{{range exported .Funcs}}{{.Name}}({{range $i, $p := .Signature.Params}}{{if $i}}, {{end}}{{$p.Name}} {{$p.Type}}{{end}}){{if .Signature.Results}} ({{range $i, $r := .Signature.Results}}{{if $i}}, {{end}}{{$r.Type}}{{end}}){{end}}{{br}}{{end}}
```

### Find Slice and Array Types

```go
{{range exported .Types}}{{$t := .}}{{with slice .Type}}{{$t.Name}} []{{.Elem}}{{br}}{{end}}{{with array .Type}}{{$t.Name}} [{{.Len}}]{{.Elem}}{{br}}{{end}}{{end}}
```

### List Map Types with Key-Value Info

```go
{{range exported .Types}}{{$t := .}}{{with map .Type}}{{$t.Name}} map[{{.Key}}]{{.Elem}}{{br}}{{end}}{{end}}
```

### Find Pointer Types

```go
{{range exported .Types}}{{$t := .}}{{with pointer .Type}}{{$t.Name}} *{{.Elem}}{{br}}{{end}}{{end}}
```

### Show Constructor Functions (New* pattern)

```go
{{range exported .Funcs}}{{if and (regexp "^New" .Name) (gt (len .Signature.Results) 0)}}{{.Name}} -> {{index .Signature.Results 0}.Type}}{{br}}{{end}}{{end}}
```

### Find Functions Returning Errors

```go
{{range exported .Funcs}}{{with .Signature}}{{if and (gt (len .Results) 0) (identical (index .Results -1).Type (typeof "error"))}}{{.Recv.Name}}{{br}}{{end}}{{end}}{{end}}
```

### Get Documentation for Types and Functions

```go
{{range exported .Types}}{{.Name}}: {{godoc .Name}}{{br}}{{end}}
```

### Show Package Documentation

```go
Package {{.Name}}: {{godoc .Path}}
```

### Get Function Signatures with Documentation

```go
{{range exported .Funcs}}{{.Name}}: {{godoc (printf "%s.%s" .Package.Path .Name)}}{{br}}{{end}}
```

### Get Full Documentation for Specific Function

```go
{{with (index .Funcs "FunctionName")}}{{godoc (printf "%s.%s" .Package.Path .Name)}}{{end}}
```

### Show Type Documentation with Methods

```go
{{range exported .Types}}{{.Name}}:
{{godoc (printf "%s.%s" .Package.Path .Name)}}
{{with named .Type}}{{range .Methods}}  {{.Name}}: {{godoc (printf "%s.%s.%s" .Package.Path .Signature.Recv.Type .Name)}}{{br}}{{end}}{{end}}
{{br}}{{end}}
```

### Get Standard Library Documentation

```go
{{godoc "fmt"}}  // Package documentation
{{godoc "fmt.Printf"}}  // Function documentation
{{godoc "io.Reader"}}  // Interface documentation
```

### Documentation for Current Package Elements

```go
{{range exported .Types}}
// Type: {{.Name}}
{{godoc (printf "%s.%s" .Package.Path .Name)}}
{{br}}{{end}}

{{range exported .Funcs}}
// Function: {{.Name}}
{{godoc (printf "%s.%s" .Package.Path .Name)}}
{{br}}{{end}}
```

## Advanced Features

### XPath Filtering (knife only)

Filter AST nodes before applying templates:

```go
{{range .}}{{.Name}}{{br}}{{end}}
```

### Extra Data Access

Access additional data passed via command line:

```go
{{.Name}} (version: {{data "version"}})
```

## Notes

- All standard Go template functions are available
- Type conversion functions return `nil` if the conversion is not possible
- Position information requires the object to implement `Pos() token.Pos`
- Documentation extraction looks for comments preceding declarations
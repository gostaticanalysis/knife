# knife [![PkgGoDev](https://pkg.go.dev/badge/github.com/gostaticanalysis/knife)](https://pkg.go.dev/github.com/gostaticanalysis/knife)

**knife** is a CLI tool for inspecting Go packages, focusing on listing type and object information.  
It provides a `-f` option similar to `go list`, allowing you to customize output via Go templates.

Additionally, this repository offers several **separate** CLI tools that work similarly:

- **knife**: The main CLI (comprehensive tool for listing and inspecting Go package objects)
- **cutter**: A lightweight tool focusing on listing type information
- **typels**: Lists types in a package
- **objls**: Lists objects in a package
- **hagane**: A template-based code generator

---

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
   1. [Common Options](#common-options)
   2. [Functions for Templates](#functions-for-templates)
   3. [Example Commands](#example-commands)
3. [Related Tools](#related-tools)
   1. [cutter](#cutter)
   2. [typels](#typels)
   3. [objls](#objls)
   4. [hagane](#hagane)
4. [License](#license)
5. [Author](#author)

---

## Installation

### knife

```sh
go install github.com/gostaticanalysis/knife/cmd/knife@latest
```

### cutter

```sh
go install github.com/gostaticanalysis/knife/cmd/cutter@latest
```

### typels

```sh
go install github.com/gostaticanalysis/knife/cmd/typels@latest
```

### objls

```sh
go install github.com/gostaticanalysis/knife/cmd/objls@latest
```

### hagane

```sh
go install github.com/gostaticanalysis/knife/cmd/hagane@latest
```

---

## Usage

### Common Options

`knife` works similarly to `go list`: You can specify a Go template using the `-f` option to customize the output.  
For example:

```sh
knife -f "{{.}}" fmt
```

For more details, see [Options](./_docs/options.md).

### Functions for Templates

Within the template specified by the `-f` flag, you can use various helper functions alongside the standard Go `text/template` package. For instance:

- `exported`: Filters for exported items only
- `identical`: Checks if two types are identical
- `implements`: Checks if a type implements a given interface
- `typeof`: Returns the type by name
- `pos`: Retrieves the position (file and line) of a definition
- `br`: Inserts a line break in the output

See the full list of available functions in [Functions for a template](./_docs/funcs.md).

### Example Commands

Below are some common examples:

1. **List functions in the `fmt` package whose names begin with `Print`:**

   ```sh
   knife -f "{{range exported .Funcs}}{{.Name}}{{br}}{{end}}" fmt | grep Print
   Print
   Printf
   Println
   ```

2. **List exported types in the `fmt` package:**

   ```sh
   knife -f "{{range exported .Types}}{{.Name}}{{br}}{{end}}" fmt
   Formatter
   GoStringer
   ScanState
   Scanner
   State
   Stringer
   ```

3. **List functions in `net/http` whose first parameter is `context.Context`:**

   ```sh
   knife -f '{{range exported .Funcs}}{{.Name}} {{with .Signature.Params}}{{index . 0}}{{end}}{{br}}{{end}}' net/http | grep context.Context
   NewRequestWithContext var ctx context.Context
   ```

4. **List variables in `net/http` that implement the `error` interface:**

   ```sh
   knife -f '{{range exported .Vars}}{{if implements . (typeof "error")}}{{.Name}}{{br}}{{end}}{{end}}' net/http
   ErrAbortHandler
   ErrBodyNotAllowed
   ErrBodyReadAfterClose
   ...
   ErrWriteAfterFlush
   ```

5. **List the position of fields whose type is `context.Context`:**

   ```sh
   knife -f '{{range .Types}}{{$t := .}}{{with struct .}}{{range .Fields}}{{if identical . (typeof "context.Context")}}{{$t.Name}} - {{pos .}}{{br}}{{end}}{{end}}{{end}}{{end}}' net/http
   Request - /usr/local/go/src/net/http/request.go:319:2
   http2ServeConnOpts - /usr/local/go/src/net/http/h2_bundle.go:3878:2
   ...
   ```

6. **Use an XPath expression to list AST node types (e.g., `FuncDecl` names starting with `Print`):**

   ```sh
   knife -f '{{range .}}{{.Name}}:{{with .Scope}}{{.Names}}{{end}}{{br}}{{end}}' -xpath '//*[@type="FuncDecl"]/Name[starts-with(@Name, "Print")]' fmt
   Printf:[a err format n]
   Print:[a err n]
   Println:[a err n]
   ```

---

## Related Tools

### cutter

`cutter` is a simplified version of `knife` that focuses on listing types.  
Usage is almost identical to `knife`:

```sh
cutter -f "{{range exported .Funcs}}{{.Name}}{{br}}{{end}}" fmt | grep Print
Print
Printf
Println
```

### typels

`typels` lists types in a package:

```sh
typels -f interface io | grep Writer
io.ByteWriter
io.ReadWriter
io.WriterAt
io.WriterTo
io.StringWriter
io.Writer
```

### objls

`objls` lists objects in a package:

```sh
objls -f const net/http | grep Status | head -5
net/http.StatusBadGateway
net/http.StatusMovedPermanently
net/http.StatusNotFound
net/http.StatusCreated
net/http.StatusForbidden
```

### hagane

`hagane` is a template-based code generator that can produce Go code based on a specified template and source file(s):

```sh
hagane -template template.go.tmpl -o sample_mock.go -data '{"type":"DB"}' sample.go
```

- `-o`: Output file path (defaults to stdout)
- `-f`: Template format (defaults to `{{.}}`)
- `-template`: Template file (used if `-f` is not set)
- `-data`: Extra data (JSON) passed into the template

For a complete example, see [this hagane sample](./_examples/hagane/).

---

## License

This project is licensed under the [MIT License](./LICENSE).

Contributions are always welcome! Feel free to open issues or PRs for bugs and enhancements.

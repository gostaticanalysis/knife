# tlist [![PkgGoDev](https://pkg.go.dev/badge/github.com/gostaticanalysis/tlist)](https://pkg.go.dev/github.com/gostaticanalysis/tlist)

`tlist` lists type information of the named packages, one per line.
`tlist` has `-f` option likes `go list` command.

## Install

```sh
$ go get -u github.com/tenntenn/gostaticanalysis/tlist
```

## How to use

### List `fmt` package's functions which name begins `Print`

```sh
$ tlist -f "{{range exported .Funcs}}{{.Name}}{{br}}{{end}}" fmt | grep Print
Print
Printf
Println
```

### List `fmt` package's exported types

```sh
$ tlist -f "{{range exported .Types}}{{.Name}}{{br}}{{end}}" fmt
Formatter
GoStringer
ScanState
Scanner
State
Stringer
```

### List `net/http` package's functions which first parameter is context.Context

```sh
$ tlist -f '{{range exported .Funcs}}{{.Name}} \
{{with .Signature.Params}}{{index . 0}}{{end}}{{br}}{{end}}' "net/http" | grep context.Context
NewRequestWithContext var ctx context.Context
```

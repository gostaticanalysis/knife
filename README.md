# knife [![PkgGoDev](https://pkg.go.dev/badge/github.com/gostaticanalysis/knife)](https://pkg.go.dev/github.com/gostaticanalysis/knife)

`knife` lists type information of the named packages, one per line.
`knife` has `-f` option likes `go list` command.

## Install

```sh
$ go get -u github.com/gostaticanalysis/knife/cmd/knife
```

## How to use

### List `fmt` package's functions which name begins `Print`

```sh
$ knife -f "{{range exported .Funcs}}{{.Name}}{{br}}{{end}}" fmt | grep Print
Print
Printf
Println
```

### List `fmt` package's exported types

```sh
$ knife -f "{{range exported .Types}}{{.Name}}{{br}}{{end}}" fmt
Formatter
GoStringer
ScanState
Scanner
State
Stringer
```

### List `net/http` package's functions which first parameter is context.Context

```sh
$ knife -f '{{range exported .Funcs}}{{.Name}} \
{{with .Signature.Params}}{{index . 0}}{{end}}{{br}}{{end}}' "net/http" | grep context.Context
NewRequestWithContext var ctx context.Context
```

### List net/http types which implements error interface

```sh
$ knife -f '{{range exported .Vars}}{{if implements . (typeof "error")}}{{.Name}}{{br}}{{end}}{{end}}' "net/http"
ErrAbortHandler
ErrBodyNotAllowed
ErrBodyReadAfterClose
ErrContentLength
ErrHandlerTimeout
ErrHeaderTooLong
ErrHijacked
ErrLineTooLong
ErrMissingBoundary
ErrMissingContentLength
ErrMissingFile
ErrNoCookie
ErrNoLocation
ErrNotMultipart
ErrNotSupported
ErrServerClosed
ErrShortBody
ErrSkipAltProtocol
ErrUnexpectedTrailer
ErrUseLastResponse
ErrWriteAfterFlush
```

### List position of fields which type is context.Context

```sh
$ knife -f '{{range .Types}}{{$t := .}}{{with struct .}}{{range .Fields}}{{if identical . (typeof "context.Context")}}{{$t.Name}} - {{pos .}}{{br}}{{end}}{{end}}{{end}}{{end}}' "net/http"
Request - /usr/local/go/src/net/http/request.go:319:2
http2ServeConnOpts - /usr/local/go/src/net/http/h2_bundle.go:3878:2
http2serverConn - /usr/local/go/src/net/http/h2_bundle.go:4065:2
http2stream - /usr/local/go/src/net/http/h2_bundle.go:4146:2
initALPNRequest - /usr/local/go/src/net/http/server.go:3393:2
timeoutHandler - /usr/local/go/src/net/http/server.go:3241:2
wantConn - /usr/local/go/src/net/http/transport.go:1162:2
```

### List type information of an AST node which is selected by a XPath expression

```sh
$ knife -f '{{range .}}{{.Name}}:{{with .Scope}}{{.Names}}{{end}}{{br}}{{end}}' -xpath '//*[@type="FuncDecl"]/Name[starts-with(@Name, "Print")]' fmt
Printf:[a err format n]
Print:[a err n]
Println:[a err n]
```

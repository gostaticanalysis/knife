package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"
	"strings"
	"text/template"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/knife/knife"
	"golang.org/x/tools/go/packages"
)

var (
	flagFormat string
)

func init() {
	flag.StringVar(&flagFormat, "f", "{{.}}", "output format")
	flag.Parse()
}

func main() {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps}
	pkgs, err := packages.Load(cfg, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(1)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		pkg := pkg
		var cmaps comment.Maps
		tmpl, err := knife.Template.Funcs(template.FuncMap{
			"pos": func(v interface{}) token.Position {
				return knife.Position(pkg.Fset, v)
			},
			"doc": func(v interface{}) string {
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
			},
		}).Parse(flagFormat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "template parse error: %v\n", err)
			os.Exit(1)
		}
		p := knife.NewPackage(pkg.Types)
		if err := tmpl.Execute(os.Stdout, p); err != nil {
			fmt.Fprintf(os.Stderr, "template execute: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
	}
}

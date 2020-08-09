package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"
	"text/template"

	"github.com/gostaticanalysis/tlist/tlist"
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
		tmpl, err := tlist.Template.Funcs(template.FuncMap{
			"pos": func(v interface{}) token.Position {
				return tlist.Position(pkg.Fset, v)
			},
		}).Parse(flagFormat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "template parse error: %v\n", err)
			os.Exit(1)
		}
		p := tlist.NewPackage(pkg.Types)
		if err := tmpl.Execute(os.Stdout, p); err != nil {
			fmt.Fprintf(os.Stderr, "template execute: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
	}
}

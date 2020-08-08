package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/gostaticanalysis/tlist/tlist"
	"golang.org/x/tools/go/packages"
)

func init() {
	var flagFormat string
	flag.StringVar(&flagFormat, "f", "{{.}}", "output format")
	flag.Parse()
	template.Must(tlist.Template.Parse(flagFormat))
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
		p := tlist.NewPackage(pkg.Types)
		if err := tlist.Template.Execute(os.Stdout, p); err != nil {
			fmt.Fprintf(os.Stderr, "template execute: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
	}
}

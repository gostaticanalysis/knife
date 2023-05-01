package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gostaticanalysis/knife"
)

var (
	flagFilter   string
	flagExported bool
)

func init() {
	flag.StringVar(&flagFilter, "f", "all", "object filter(all|interface|func|struct|chan|array|slice|map)")
	flag.BoolVar(&flagExported, "exported", false, "filter only exported types")
	flag.Parse()
}

func main() {
	if err := run(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	k, err := knife.New(args...)
	if err != nil {
		return err
	}

	var w io.Writer = os.Stdout

	pkgs := k.Packages()
	for i, pkg := range pkgs {

		kpkg := knife.NewPackage(pkg.Types)

		if len(pkgs) > 1 {
			fmt.Fprintf(w, "# %s\n", kpkg.Path)
		}

		for name, typ := range kpkg.Types {
			if flagExported && !typ.Exported {
				continue
			}

			if typ.Exported && match(flagFilter, typ) {
				fmt.Fprintln(w, name)
			}
		}

		if i != len(pkgs)-1 {
			fmt.Fprintln(w)
		}
	}

	return nil
}

func match(filter string, typ *knife.TypeName) bool {
	switch filter {
	case "interface":
		return knife.ToInterface(typ) != nil
	case "func":
		return knife.ToSignature(typ) != nil
	case "struct":
		return knife.ToStruct(typ) != nil
	case "chan":
		return knife.ToChan(typ) != nil
	case "array":
		return knife.ToArray(typ) != nil
	case "slice":
		return knife.ToSlice(typ) != nil
	case "map":
		return knife.ToMap(typ) != nil
	default:
		return true
	}
}

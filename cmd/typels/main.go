package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/newmo-oss/gogroup"

	"github.com/gostaticanalysis/analysisutil"

	"github.com/gostaticanalysis/knife"
	"github.com/gostaticanalysis/knife/cutter"
)

var (
	flagKind       string
	flagImplements string
	flagExported   bool
	flagPos        bool
)

func init() {
	flag.StringVar(&flagKind, "kind", "all", "all|interface|func|struct|chan|array|slice|map")
	flag.StringVar(&flagImplements, "implements", "", "implements interface")
	flag.BoolVar(&flagExported, "exported", true, "filter only exported types")
	flag.BoolVar(&flagPos, "pos", false, "print position")
	flag.Parse()
}

func main() {
	if err := run(context.Background(), flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	cutterOpt := &cutter.CutterOption{Tests: true}
	c, err := cutter.New(cutterOpt, args...)
	if err != nil {
		return err
	}

	pkgs := c.KnifePackages()
	readers := make([]io.Reader, len(pkgs))
	var g gogroup.Group
	for i, pkg := range pkgs {
		var buf bytes.Buffer
		readers[i] = &buf

		if len(pkg.Types) == 0 {
			continue
		}

		g.Add(func(ctx context.Context) error {

			for name, typ := range pkg.Types {
				if !match(typ) {
					continue
				}

				if flagPos {
					pos := c.Position(typ)
					if _, err := fmt.Fprintf(&buf, "%s.%s(%s:%d)\n", pkg.Path, name, filepath.Base(pos.Filename), pos.Line); err != nil {
						return err
					}
				} else {
					if _, err := fmt.Fprintf(&buf, "%s.%s\n", pkg.Path, name); err != nil {
						return err
					}
				}
			}

			return nil
		})
	}

	if err := g.Run(ctx); err != nil {
		return err
	}

	if _, err := io.Copy(os.Stdout, io.MultiReader(readers...)); err != nil {
		return err
	}

	return nil
}

func match(typ *knife.TypeName) bool {
	if flagExported && !typ.Exported {
		return false
	}

	if !matchKind(typ) {
		return false
	}

	if !checkImplements(typ) {
		return false
	}

	return true
}

func matchKind(typ *knife.TypeName) bool {
	switch flagKind {
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

func checkImplements(typ *knife.TypeName) bool {
	if flagImplements == "" {
		return false
	}

	pkg := typ.Package.TypesPackage
	if pkg == nil {
		return false
	}

	if implements(pkg, typ.Type.TypesType) {
		return true
	}

	if _, isInterface := typ.Type.TypesType.Underlying().(*types.Interface); isInterface {
		return false
	}

	return implements(pkg, types.NewPointer(typ.Type.TypesType))
}

func implements(pkg *types.Package, typ types.Type) bool {

	iface, _ := typeOf(pkg, flagImplements).(*types.Interface)
	if iface == nil {
		return false
	}

	return types.Implements(typ, iface)
}

func typeOf(pkg *types.Package, s string) types.Type {
	if s == "" {
		return nil
	}

	obj := objectOf(pkg, s)
	if obj == nil {
		return nil
	}

	return obj.Type().Underlying()
}

func objectOf(pkg *types.Package, s string) types.Object {
	dotPos := strings.LastIndex(s, ".")

	if dotPos == -1 {
		return types.Universe.Lookup(s)
	}

	pkgpath, name := s[:dotPos], s[dotPos+1:]
	obj := analysisutil.LookupFromImports(pkg.Imports(), pkgpath, name)
	if obj != nil {
		return obj
	}

	if analysisutil.RemoveVendor(pkg.Path()) != analysisutil.RemoveVendor(pkgpath) {
		return nil
	}

	return pkg.Scope().Lookup(name)
}

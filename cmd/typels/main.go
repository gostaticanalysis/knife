package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/newmo-oss/gogroup"

	"github.com/gostaticanalysis/knife"
	"github.com/gostaticanalysis/knife/cutter"
)

var (
	flagFilter   string
	flagExported bool
	flagPos      bool
)

func init() {
	flag.StringVar(&flagFilter, "f", "all", "object filter(all|interface|func|struct|chan|array|slice|map)")
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
	c, err := cutter.New(args...)
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
				if flagExported && !typ.Exported {
					continue
				}

				if typ.Exported && match(flagFilter, typ) {
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

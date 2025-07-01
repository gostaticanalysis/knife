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
	flag.StringVar(&flagFilter, "f", "all", "object filter(all|const|func|var)")
	flag.BoolVar(&flagExported, "exported", true, "filter only exported object")
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
			for name, obj := range pkg.Objects() {
				if flagExported && !obj.TypesObject().Exported() {
					continue
				}

				if match(flagFilter, obj) {
					if flagPos {
						pos := c.Position(obj)
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

func match(filter string, obj knife.Object) bool {
	if filter == "all" || filter == "" {
		return true
	}

	switch obj.(type) {
	case *knife.Const:
		return filter == "const"
	case *knife.Func:
		return filter == "func"
	case *knife.Var:
		return filter == "var"
	}
	return true
}

package cutter

import (
	"context"
	"fmt"
	"go/token"
	"io"

	"golang.org/x/tools/go/packages"

	"github.com/newmo-oss/gogroup"

	"github.com/gostaticanalysis/knife"
)

// Cutter is a lightweight version of Knife which is resterected for type information.
type Cutter struct {
	fset      *token.FileSet
	pkgs      []*packages.Package
	knifePkgs []*knife.Package
}

// New creates a [Cutter].
func New(patterns ...string) (*Cutter, error) {
	mode := packages.NeedName | packages.NeedTypes
	cfg := &packages.Config{
		Fset: token.NewFileSet(),
		Mode: mode,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	knifePkgs := make([]*knife.Package, len(pkgs))
	var g gogroup.Group
	for i := range pkgs {
		g.Add(func(ctx context.Context) error {
			knifePkgs[i] = knife.NewPackage(pkgs[i].Types)
			return nil
		})
	}

	if err := g.Run(context.Background()); err != nil {
		return nil, err
	}

	return &Cutter{
		fset:      cfg.Fset,
		pkgs:      pkgs,
		knifePkgs: knifePkgs,
	}, nil
}

// Packages returns packages.
func (c *Cutter) Packages() []*packages.Package {
	return c.pkgs
}

// KnifePackages returns knife packages.
func (c *Cutter) KnifePackages() []*knife.Package {
	return c.knifePkgs
}

// Option is a option of Execute.
type Option struct {
	ExtraData map[string]any
}

// Execute outputs the pkg with the format.
func (c *Cutter) Execute(w io.Writer, pkg *knife.Package, tmpl any, opt *Option) error {

	var tmplStr string
	switch tmpl := tmpl.(type) {
	case string:
		tmplStr = tmpl
	case []byte:
		tmplStr = string(tmpl)
	case io.Reader:
		b, err := io.ReadAll(tmpl)
		if err != nil {
			return fmt.Errorf("cannnot read template: %w", err)
		}
		tmplStr = string(b)
	default:
		return fmt.Errorf("template must be string, []byte or io.Reader: %T", tmpl)
	}

	td := &knife.TempalteData{
		Fset:  c.fset,
		Pkg:   pkg.TypesPackage,
		Extra: opt.ExtraData,
	}
	t, err := knife.NewTemplate(td).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	if err := t.Execute(w, pkg); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

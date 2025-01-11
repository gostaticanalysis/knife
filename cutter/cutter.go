package cutter

import (
	"fmt"
	"io"

	"golang.org/x/tools/go/packages"

	"github.com/gostaticanalysis/knife"
)

// Cutter is a lightweight version of Knife which is resterected for type information.
type Cutter struct {
	pkgs []*packages.Package
}

// New creates a [Cutter].
func New(patterns ...string) (*Cutter, error) {
	mode := packages.NeedName | packages.NeedTypes
	cfg := &packages.Config{Mode: mode}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	return &Cutter{
		pkgs: pkgs,
	}, nil
}

// Packages returns packages.
func (c *Cutter) Packages() []*packages.Package {
	return c.pkgs
}

// Option is a option of Execute.
type Option struct {
	ExtraData map[string]any
}

// Execute outputs the pkg with the format.
func (c *Cutter) Execute(w io.Writer, pkg *packages.Package, tmpl any, opt *Option) error {

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
		Fset:  pkg.Fset,
		Pkg:   pkg.Types,
		Extra: opt.ExtraData,
	}
	t, err := knife.NewTemplate(td).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	data := knife.NewPackage(pkg.Types)
	if err := t.Execute(w, data); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

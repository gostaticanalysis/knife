package knife

import (
	"fmt"
	"go/ast"
	"io"

	"github.com/gostaticanalysis/astquery"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

type Knife struct {
	pkgs []*packages.Package
	ins  map[*packages.Package]*inspector.Inspector
}

func New(patterns ...string) (*Knife, error) {
	mode := packages.NeedFiles | packages.NeedSyntax |
		packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo
	cfg := &packages.Config{Mode: mode}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	ins := make(map[*packages.Package]*inspector.Inspector, len(pkgs))
	for _, pkg := range pkgs {
		ins[pkg] = inspector.New(pkg.Syntax)
	}

	return &Knife{pkgs: pkgs, ins: ins}, nil
}

func (k *Knife) Format(w io.Writer, format string) error {
	for _, pkg := range k.pkgs {
		if err := k.formatPkg(w, format, pkg); err != nil {
			return err
		}
	}
	return nil
}

func (k *Knife) formatPkg(w io.Writer, format string, pkg *packages.Package) error {
	tmpl, err := NewTemplate(pkg, format)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	if err := tmpl.Execute(w, NewPackage(pkg.Types)); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

func (k *Knife) FormatWithXPath(w io.Writer, format, xpath string) error {
	for _, pkg := range k.pkgs {
		if err := k.formatPkgWithXPath(w, format, xpath, pkg); err != nil {
			return err
		}
	}
	return nil
}

func (k *Knife) formatPkgWithXPath(w io.Writer, format, xpath string, pkg *packages.Package) error {
	tmpl, err := NewTemplate(pkg, format)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	e := astquery.New(pkg.Fset, pkg.Syntax, k.ins[pkg])
	v, err := e.Eval(xpath)
	if err != nil {
		return fmt.Errorf("XPath parse error: %w", err)
	}

	var data interface{}
	switch v := v.(type) {
	case []ast.Node:
		ns := make([]*ASTNode, len(v))
		for i := range ns {
			ns[i] = NewASTNode(pkg.TypesInfo, v[i])
		}
		data = ns
	default:
		data = v
	}

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

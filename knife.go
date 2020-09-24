package knife

import (
	"fmt"
	"go/ast"
	"io"
	"io/ioutil"

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

// Packages returns packages.
func (k *Knife) Packages() []*packages.Package {
	return k.pkgs
}

// Option is a option of Execute.
type Option struct {
	XPath     string
	ExtraData map[string]interface{}
}

// Execute outputs the pkg with the format.
func (k *Knife) Execute(w io.Writer, pkg *packages.Package, format string, opt *Option) error {
	tmpl, err := newTemplate(pkg, opt.ExtraData).Parse(format)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	var data interface{}

	switch {
	case opt != nil && opt.XPath != "":
		data, err = k.evalXPath(pkg, opt.XPath)
		if err != nil {
			return err
		}
	default:
		data = NewPackage(pkg.Types)
	}

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

// ExecuteWithTemplate outputs the pkg with the template.
func (k *Knife) ExecuteWithTemplate(w io.Writer, pkg *packages.Package, tmplPath string, opt *Option) error {
	format, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("cannot read template file: %w", err)
	}

	tmpl, err := newTemplate(pkg, opt.ExtraData).Parse(string(format))
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	var data interface{}

	switch {
	case opt != nil && opt.XPath != "":
		data, err = k.evalXPath(pkg, opt.XPath)
		if err != nil {
			return err
		}
	default:
		data = NewPackage(pkg.Types)
	}

	if err := tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

func (k *Knife) evalXPath(pkg *packages.Package, xpath string) (interface{}, error) {
	e := astquery.New(pkg.Fset, pkg.Syntax, k.ins[pkg])
	v, err := e.Eval(xpath)
	if err != nil {
		return nil, fmt.Errorf("XPath parse error: %w", err)
	}

	switch v := v.(type) {
	case []ast.Node:
		ns := make([]*ASTNode, len(v))
		for i := range ns {
			ns[i] = NewASTNode(pkg.TypesInfo, v[i])
		}
		return ns, nil
	}

	return v, nil
}

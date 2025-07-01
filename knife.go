package knife

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"strings"

	"github.com/gostaticanalysis/astquery"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

//go:embed version.txt
var version string

// Version returns the version of knife.
func Version() string {
	return strings.TrimSpace(version)
}

type Knife struct {
	fset *token.FileSet
	pkgs []*packages.Package
	ins  map[*packages.Package]*inspector.Inspector
}

func New(patterns ...string) (*Knife, error) {
	mode := packages.NeedFiles | packages.NeedSyntax |
		packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo
	cfg := &packages.Config{
		Fset:  token.NewFileSet(),
		Mode:  mode,
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	ins := make(map[*packages.Package]*inspector.Inspector, len(pkgs))
	for _, pkg := range pkgs {
		ins[pkg] = inspector.New(pkg.Syntax)
	}

	return &Knife{
		fset: cfg.Fset,
		pkgs: pkgs,
		ins:  ins,
	}, nil
}

// Packages returns packages.
func (k *Knife) Packages() []*packages.Package {
	return k.pkgs
}

// Position returns position of v.
func (k *Knife) Position(v any) token.Position {
	n, ok := v.(interface{ Pos() token.Pos })
	if ok && k.fset != nil {
		return k.fset.Position(n.Pos())
	}
	return token.Position{}
}

// Option is a option of Execute.
type Option struct {
	XPath     string
	ExtraData map[string]any
}

// Execute outputs the pkg with the format.
func (k *Knife) Execute(w io.Writer, pkg *packages.Package, tmpl any, opt *Option) error {

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

	td := &TempalteData{
		Fset:      pkg.Fset,
		Files:     pkg.Syntax,
		TypesInfo: pkg.TypesInfo,
		Pkg:       pkg.Types,
		Extra:     opt.ExtraData,
	}
	t, err := NewTemplate(td).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template parse: %w", err)
	}

	var data any

	switch {
	case opt != nil && opt.XPath != "":
		data, err = k.evalXPath(pkg, opt.XPath)
		if err != nil {
			return err
		}
	default:
		data = NewPackage(pkg.Types)
	}

	if err := t.Execute(w, data); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}

func (k *Knife) evalXPath(pkg *packages.Package, xpath string) (any, error) {
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

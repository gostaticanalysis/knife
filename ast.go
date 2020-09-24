package knife

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
)

type ASTNode struct {
	Node   ast.Node
	Scope  *Scope
	Type   *Type
	Name   string
	Object Object
	Value  constant.Value
}

var _ fmt.Stringer = (*ASTNode)(nil)

func (n *ASTNode) Pos() token.Pos {
	return n.Node.Pos()
}

func (n *ASTNode) String() string {
	return fmt.Sprintf("%T", n.Node)
}

func (n *ASTNode) BoolVal() bool {
	return constant.BoolVal(n.Value)
}

func (n *ASTNode) StringVal() string {
	return constant.StringVal(n.Value)
}

func (n *ASTNode) Float32Val() float32 {
	v, ok := constant.Float32Val(n.Value)
	if !ok {
		panic("unkown kind")
	}
	return v
}

func (n *ASTNode) Float64Val() float64 {
	v, ok := constant.Float64Val(n.Value)
	if !ok {
		panic("unkown kind")
	}
	return v
}

func (n *ASTNode) Int64Val() int64 {
	v, ok := constant.Int64Val(n.Value)
	if !ok {
		panic("unkown kind")
	}
	return v

}

func (n *ASTNode) Uint64Val() uint64 {
	v, ok := constant.Uint64Val(n.Value)
	if !ok {
		panic("unkown kind")
	}
	return v

}

func (n *ASTNode) Val() interface{} {
	return constant.Val(n.Value)
}

func NewASTNode(typesInfo *types.Info, n ast.Node) *ASTNode {
	if n == nil {
		return nil
	}

	v, _ := cache.Load(n)
	cached, _ := v.(*ASTNode)
	if cached != nil {
		return cached
	}

	var nn ASTNode
	cache.Store(n, &nn)

	nn.Node = n
	nn.Scope = NewScope(typesInfo.Scopes[n])
	if id, ok := n.(*ast.Ident); ok {
		obj := typesInfo.ObjectOf(id)
		if obj != nil {
			nn.Object = NewObject(obj)
			nn.Name = obj.Name()
			if scopeHolder, ok := obj.(interface{ Scope() *types.Scope }); ok {
				nn.Scope = NewScope(scopeHolder.Scope())
			}
		}
	}
	if expr, ok := n.(ast.Expr); ok {
		nn.Type = NewType(typesInfo.TypeOf(expr))
		if tv, ok := typesInfo.Types[expr]; ok {
			nn.Value = tv.Value
		}
	}
	return &nn
}

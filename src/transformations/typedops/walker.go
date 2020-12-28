package typedops

import (
	"errors"
	"semi/src/ast"
	"semi/src/ast/childgroup"
	"semi/src/ast/nodetype"
)

type typedOpsWalker struct {
	*ast.MapWalker
}

func NewTypedOpsWalker() *typedOpsWalker {
	walker := &typedOpsWalker{
		MapWalker: ast.NewMapWalker(nil),
	}

	walker.AddVisitor(nodetype.Addition, walker.additionVisitor)

	return walker
}

func (w *typedOpsWalker) additionVisitor(node ast.Node, _ interface{}) (interface{}, error) {
	var typedNode ast.Node

	rhs, err := w.VisitNode(node.GetChildForGroup(childgroup.Rhs), nil)
	if err != nil {
		return nil, err
	}

	lhs, err := w.VisitNode(node.GetChildForGroup(childgroup.Lhs), nil)
	if err != nil {
		return nil, err
	}

	rhsNode, ok := rhs.(ast.Node)
	if !ok {
		return nil, errors.New("expected a node to be returned for RHS expression")
	}

	lhsNode, ok := lhs.(ast.Node)
	if !ok {
		return nil, errors.New("expected a node to be returned for LHS expression")
	}

	// TODO: Why are we doing function resolution here? Should this be moved to codegen, where we can do lookups?
	// Replace lhs/rhs nodes with any type coercion nodes needed, clone node with those children and a type-specific Addition type

	return typedNode, nil
}

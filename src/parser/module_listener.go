package parser

import (
	"semi/src/ast"
	"semi/src/ast/childgroup"
	"semi/src/ast/nodetype"
	"semi/src/parser/gen"
)

type ModuleListener struct {
	*NodeStackListener
}

func NewModuleListener() *ModuleListener {
	return &ModuleListener{
		NodeStackListener: NewNodeStackListener(),
	}
}

func (l *ModuleListener) ExitModule(ctx *gen.ModuleContext) {
	node := ast.NewNode(nodetype.Module, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	functionDeclarationListener := NewFunctionDeclarationListener()

	l.Walk(functionDeclarationListener, ctx)

	node.AddChildren(childgroup.FunctionDeclarations, functionDeclarationListener.Nodes())

	l.PushNode(node)
}

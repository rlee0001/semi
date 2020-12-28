package parser

import (
	"semi/src/ast"
	"semi/src/ast/childgroup"
	"semi/src/ast/nodetype"
	"semi/src/parser/gen"
)

type FunctionDeclarationListener struct {
	*NodeArrayAstListener
}

func NewFunctionDeclarationListener() *FunctionDeclarationListener {
	return &FunctionDeclarationListener{
		NodeArrayAstListener: NewNodeArrayListener(),
	}
}

func (l *FunctionDeclarationListener) EnterFunctionDeclaration(ctx *gen.FunctionDeclarationContext) {
	node := ast.NewNode(nodetype.FunctionDeclaration, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	parameterListItemListener := NewParameterListItemAstListener()
	typeListener := NewTypeExpressionListener()
	localListener := NewLocalDeclarationListener()
	statementListener := NewStatementListener()

	l.Walk(parameterListItemListener, ctx.ParameterDeclarationList())
	l.Walk(typeListener, ctx.TypeExpression())
	l.Walk(localListener, ctx.Block())
	l.Walk(statementListener, ctx.Block())

	node.SetName(ctx.GetIdent().GetText())
	node.AddChildren(childgroup.ParameterDeclarations, parameterListItemListener.Nodes())
	node.AddChild(childgroup.Type, typeListener.PopNode())
	node.AddChildren(childgroup.LocalDeclarations, localListener.Nodes())
	node.AddChildren(childgroup.Statements, statementListener.Nodes())

	l.AddNode(node)
}

type ParameterListItemAstListener struct {
	*NodeArrayAstListener
}

func NewParameterListItemAstListener() *ParameterListItemAstListener {
	return &ParameterListItemAstListener{
		NodeArrayAstListener: NewNodeArrayListener(),
	}
}

func (l *ParameterListItemAstListener) EnterParameterDeclaration(ctx *gen.ParameterDeclarationContext) {
	node := ast.NewNode(nodetype.ParameterDeclaration, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	typeExpressionListener := NewTypeExpressionListener()

	l.Walk(typeExpressionListener, ctx.TypeExpression())

	node.SetName(ctx.GetIdent().GetText())
	node.AddChild(childgroup.Type, typeExpressionListener.PopNode())

	l.AddNode(node)
}

type TypeExpressionListener struct {
	*NodeStackListener
}

func NewTypeExpressionListener() *TypeExpressionListener {
	return &TypeExpressionListener{
		NodeStackListener: NewNodeStackListener(),
	}
}

func (l *TypeExpressionListener) EnterTypeExpression(ctx *gen.TypeExpressionContext) {
	node := ast.NewNode(nodetype.TypeExpression, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	node.SetName(ctx.GetIdent().GetText())

	l.PushNode(node)
}

type LocalDeclarationListener struct {
	*NodeArrayAstListener
}

func NewLocalDeclarationListener() *LocalDeclarationListener {
	return &LocalDeclarationListener{
		NodeArrayAstListener: NewNodeArrayListener(),
	}
}

func (l *LocalDeclarationListener) EnterLocalDeclaration(ctx *gen.LocalDeclarationContext) {
	node := ast.NewNode(nodetype.LocalDeclaration, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	expressionListener := NewExpressionListener()

	l.Walk(expressionListener, ctx)

	node.SetName(ctx.GetIdent().GetText())
	node.AddChild(childgroup.Initializer, expressionListener.PopNode())

	l.AddNode(node)
}

type StatementListener struct {
	*NodeArrayAstListener
}

func NewStatementListener() *StatementListener {
	return &StatementListener{
		NodeArrayAstListener: NewNodeArrayListener(),
	}
}

func (l *StatementListener) EnterStatement(ctx *gen.StatementContext) {
	node := ast.NewNode(nodetype.Statement, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	expressionListener := NewExpressionListener()

	l.Walk(expressionListener, ctx)

	node.AddChild(childgroup.Arguments, expressionListener.PopNode())

	l.AddNode(node)
}

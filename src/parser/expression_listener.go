package parser

import (
	"strconv"

	"semi/src/ast"
	"semi/src/ast/childgroup"
	"semi/src/ast/nodetype"
	"semi/src/parser/gen"
)

type ExpressionListener struct {
	*NodeStackListener
}

func NewExpressionListener() *ExpressionListener {
	return &ExpressionListener{
		NodeStackListener: NewNodeStackListener(),
	}
}

func (l *ExpressionListener) EnterIdentifierFactor(ctx *gen.IdentifierFactorContext) {
	node := ast.NewNode(nodetype.Identifier, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	node.SetName(ctx.GetIdent().GetText())

	l.PushNode(node)
}

func (l *ExpressionListener) EnterIntegerFactor(ctx *gen.IntegerFactorContext) {
	node := ast.NewNode(nodetype.Integer, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	// TODO: Handle integer parsing error
	integerValue, _ := strconv.ParseInt(ctx.GetInteger().GetText(), 10, 64)

	node.SetValue(integerValue)

	l.PushNode(node)
}

func (l *ExpressionListener) EnterFloatFactor(ctx *gen.FloatFactorContext) {
	node := ast.NewNode(nodetype.Float, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	// TODO: Handle float parsing error
	floatValue, _ := strconv.ParseFloat(ctx.GetFloat().GetText(), 64)

	node.SetValue(floatValue)

	l.PushNode(node)
}

func (l *ExpressionListener) EnterTrueFactor(ctx *gen.TrueFactorContext) {
	node := ast.NewNode(nodetype.Boolean, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	node.SetValue(true)

	l.PushNode(node)
}

func (l *ExpressionListener) EnterFalseFactor(ctx *gen.FalseFactorContext) {
	node := ast.NewNode(nodetype.Boolean, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	node.SetValue(false)

	l.PushNode(node)
}

func (l *ExpressionListener) ExitMulDivTerm(ctx *gen.MulDivTermContext) {
	var nodeType nodetype.NodeType

	if ctx.GetOp().GetText() == "*" {
		nodeType = nodetype.Multiplication
	} else {
		nodeType = nodetype.Division
	}

	node := ast.NewNode(nodeType, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	node.AddChild(childgroup.Rhs, l.PopNode())
	node.AddChild(childgroup.Lhs, l.PopNode())

	l.PushNode(node)
}

func (l *ExpressionListener) ExitAddSubTerm(ctx *gen.AddSubTermContext) {
	var nodeType nodetype.NodeType

	if ctx.GetOp().GetText() == "+" {
		nodeType = nodetype.Addition
	} else {
		nodeType = nodetype.Subtraction
	}

	node := ast.NewNode(nodeType, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	node.AddChild(childgroup.Rhs, l.PopNode())
	node.AddChild(childgroup.Lhs, l.PopNode())

	l.PushNode(node)
}

func (l *ExpressionListener) ExitIsEqualToTerm(ctx *gen.IsEqualToTermContext) {
	node := ast.NewNode(nodetype.IsEqualTo, ctx.GetStart().GetLine(), ctx.GetStart().GetColumn())

	node.AddChild(childgroup.Rhs, l.PopNode())
	node.AddChild(childgroup.Lhs, l.PopNode())

	l.PushNode(node)
}


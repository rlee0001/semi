package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/golang-collections/collections/stack"

	"semi/src/ast"
	"semi/src/parser/gen"
)

type Listener interface {
	antlr.ParseTreeListener

	Walk(listener Listener, tree antlr.Tree)
}

type NodeStackListener struct {
	*gen.BaseSemiListener

	stack *stack.Stack
}

func NewNodeStackListener() *NodeStackListener {
	return &NodeStackListener{
		stack: stack.New(),
	}
}

func (l *NodeStackListener) Walk(listener Listener, tree antlr.Tree) {
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)
}

func (l *NodeStackListener) PushNode(node ast.Node) {
	l.stack.Push(node)
}

func (l *NodeStackListener) PopNode() ast.Node {
	node, _ := l.stack.Pop().(ast.Node)

	return node
}

type NodeArrayAstListener struct {
	*gen.BaseSemiListener

	nodes []ast.Node
}

func NewNodeArrayListener() *NodeArrayAstListener {
	return &NodeArrayAstListener{
		nodes: []ast.Node{},
	}
}

func (l *NodeArrayAstListener) Walk(listener Listener, tree antlr.Tree) {
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)
}

func (l *NodeArrayAstListener) AddNode(node ast.Node) {
	l.nodes = append(l.nodes, node)
}

func (l *NodeArrayAstListener) Nodes() []ast.Node {
	return l.nodes
}

package ast

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"semi/src/ast/childgroup"
	"semi/src/ast/nodetype"
)

type Walker interface {
	VisitNode(node Node, memo interface{}) interface{}
	VisitNodes(nodes []Node, memo interface{}) []interface{}
}

type Visitor func(node Node, memo interface{}) (interface{}, error)

type MapWalker struct {
	visitorMap     map[nodetype.NodeType]Visitor
	defaultVisitor Visitor
}

func NewMapWalker(defaultVisitor Visitor) *MapWalker {
	mapWalker := &MapWalker{
		visitorMap:     map[nodetype.NodeType]Visitor{},
		defaultVisitor: defaultVisitor,
	}

	if defaultVisitor == nil {
		defaultVisitor = mapWalker.DefaultDefaultVisitor
	}

	return mapWalker
}

func (w *MapWalker) lookupVisitorForNodeType(nodeType nodetype.NodeType) Visitor {
	if visitor, ok := w.visitorMap[nodeType]; ok {
		return visitor
	}

	return w.defaultVisitor
}

func (w *MapWalker) AddVisitor(nodeType nodetype.NodeType, visitor Visitor) {
	w.visitorMap[nodeType] = visitor
}

func (w *MapWalker) VisitNode(node Node, memo interface{}) (interface{}, error) {
	visitor := w.lookupVisitorForNodeType(node.GetNodeType())

	return visitor(node, memo)
}

func (w *MapWalker) VisitNodes(nodes []Node, memo interface{}) ([]interface{}, error) {
	var returns []interface{}
	errors := &multierror.Error{}

	for _, node := range nodes {
		ret, err := w.VisitNode(node, memo)
		if err != nil {
			// TODO: Stop here immediately or keep going?
			errors = multierror.Append(errors, err)
		}

		returns = append(returns, ret)
	}

	return returns, errors.ErrorOrNil()
}

func (w *MapWalker) DefaultDefaultVisitor(node Node, memo interface{}) (interface{}, error) {
	childGroups := map[childgroup.ChildGroup][]Node{}
	hasChangedChildren := false

	for _, childGroup := range node.GetChildGroups() {
		childGroups[childGroup] = []Node{}
		groupChildren := node.GetChildrenForGroup(childGroup)

		for _, childNode := range groupChildren {
			retValue, err := w.VisitNode(childNode, memo)
			if err != nil {
				return nil, err
			}

			newChildNode, ok := retValue.(Node)
			if !ok {
				return nil, errors.New("expected node from visitor")
			}

			hasChangedChildren = hasChangedChildren || newChildNode != childNode
			childGroups[childGroup] = append(childGroups[childGroup], newChildNode)
		}
	}

	if !hasChangedChildren {
		return node, nil
	}

	newNode := NewNode(node.GetNodeType(), node.GetSourceLine(), node.GetSourceColumn())

	newNode.SetName(node.GetName())
	newNode.SetDataType(*node.GetDataType())
	newNode.SetValue(node.GetValue())

	for _, childGroup := range node.GetChildGroups() {
		newNode.AddChildren(childGroup, childGroups[childGroup])
	}

	return newNode, nil
}

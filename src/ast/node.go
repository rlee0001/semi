package ast

import (
	"sort"

	"semi/src/ast/childgroup"
	"semi/src/ast/datatype"
	"semi/src/ast/nodetype"
)

type Node interface {
	// The type of the node itself (function declaration, integer literal, function call, addition operator, etc...)
	GetNodeType() nodetype.NodeType

	GetSourceLine() int
	GetSourceColumn() int

	// For nodes that have names (function and local declarations and references)
	GetName() string

	// For nodes that evaluate to a data type (can be used in expressions)
	GetDataType() *datatype.DataType // TODO: Remove this.

	// For nodes that have immediate values (literals, constants)
	GetValue() interface{}
	GetBooleanValue() bool
	GetIntegerValue() int64
	GetFloatValue() float64
	GetStringValue() string

	// For nodes that have children
	GetChildGroups() []childgroup.ChildGroup
	GetChildrenForGroup(childGroup childgroup.ChildGroup) []Node
	GetChildForGroup(childGroup childgroup.ChildGroup) Node
}

type BasicNode struct {
	nodeType     nodetype.NodeType
	sourceLine   int
	sourceColumn int
	childNodes   map[childgroup.ChildGroup][]Node
	name         string
	dataType     *datatype.DataType
	value        interface{}
}

func NewNode(nodeType nodetype.NodeType, sourceLine int, sourceColumn int) *BasicNode {
	return &BasicNode{
		nodeType:     nodeType,
		sourceLine:   sourceLine,
		sourceColumn: sourceColumn,
		childNodes:   map[childgroup.ChildGroup][]Node{},
	}
}

func (n *BasicNode) GetNodeType() nodetype.NodeType {
	return n.nodeType
}

func (n *BasicNode) GetSourceLine() int {
	return n.sourceLine
}

func (n *BasicNode) GetSourceColumn() int {
	return n.sourceColumn
}

func (n *BasicNode) GetName() string {
	return n.name
}

func (n *BasicNode) SetName(name string) {
	n.name = name
}

func (n *BasicNode) GetDataType() *datatype.DataType {
	return n.dataType
}

func (n *BasicNode) SetDataType(dataType datatype.DataType) {
	n.dataType = &dataType
}

func (n *BasicNode) GetValue() interface{} {
	return n.value
}

func (n *BasicNode) GetBooleanValue() bool {
	booleanValue, _ := n.value.(bool)

	return booleanValue
}

func (n *BasicNode) GetIntegerValue() int64 {
	integerValue, _ := n.value.(int64)

	return integerValue
}

func (n *BasicNode) GetFloatValue() float64 {
	floatValue, _ := n.value.(float64)

	return floatValue
}

func (n *BasicNode) GetStringValue() string {
	stringValue, _ := n.value.(string)

	return stringValue
}

func (n *BasicNode) SetValue(value interface{}) {
	n.value = value
}

func (n *BasicNode) GetChildGroups() []childgroup.ChildGroup {
	childGroups := make([]childgroup.ChildGroup, 0, len(n.childNodes))

	for childGroup := range n.childNodes {
		childGroups = append(childGroups, childGroup)
	}

	sort.Slice(childGroups, func(i, j int) bool {
		return childGroups[i] < childGroups[j]
	})

	return childGroups
}

func (n *BasicNode) GetChildrenForGroup(childGroup childgroup.ChildGroup) []Node {
	if children, ok := n.childNodes[childGroup]; ok {
		return children
	}

	return nil
}

func (n *BasicNode) GetChildForGroup(childGroup childgroup.ChildGroup) Node {
	if children, ok := n.childNodes[childGroup]; ok {
		return children[0]
	}

	return nil
}

func (n *BasicNode) AddChild(childGroup childgroup.ChildGroup, childNode Node) {
	if nodes, ok := n.childNodes[childGroup]; ok {
		n.childNodes[childGroup] = append(nodes, childNode)
	} else {
		n.childNodes[childGroup] = []Node{childNode}
	}
}

func (n *BasicNode) AddChildren(childGroup childgroup.ChildGroup, children []Node) {
	if nodes, ok := n.childNodes[childGroup]; ok {
		n.childNodes[childGroup] = append(nodes, children...)
	} else {
		n.childNodes[childGroup] = children
	}
}

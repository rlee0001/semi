package ast

import (
	"semi/src/ast/datatype"
	"testing"

	"github.com/stretchr/testify/assert"

	"semi/src/ast/nodetype"
)

func TestNodeTypeCorrect(t *testing.T) {
	node := NewNode(nodetype.Addition, 1, 1)
	assert.Equal(t, nodetype.Addition, node.GetNodeType())
	node2 := NewNode(nodetype.Module, 1, 1)
	assert.Equal(t, nodetype.Module, node2.GetNodeType())
	node3 := NewNode(nodetype.Integer, 1, 1)
	assert.Equal(t, nodetype.Integer, node3.GetNodeType())
}

func TestNodeSourceLineCorrect(t *testing.T) {
	node := NewNode(nodetype.Module, 12, 1)
	assert.Equal(t, 12, node.GetSourceLine())
	node2 := NewNode(nodetype.Module, 99, 1)
	assert.Equal(t, 99, node2.GetSourceLine())
	node3 := NewNode(nodetype.Module, -4, 1)
	assert.Equal(t, -4, node3.GetSourceLine())
}

func TestNodeSourceColumnCorrect(t *testing.T) {
	node := NewNode(nodetype.Module, 1, 12)
	assert.Equal(t, 12, node.GetSourceColumn())
	node2 := NewNode(nodetype.Module, 1, 99)
	assert.Equal(t, 99, node2.GetSourceColumn())
	node3 := NewNode(nodetype.Module, 1, -4)
	assert.Equal(t, -4, node3.GetSourceColumn())
}

func TestNodeGetSetNameWorks(t *testing.T) {
	node := NewNode(nodetype.Module, 1, 1)
	assert.Equal(t, "", node.GetName())

	node.SetName("FOO")
	assert.Equal(t, "FOO", node.GetName())

	node.SetName("BAR")
	assert.Equal(t, "BAR", node.GetName())
}

func TestNodeGetSetDataTypeWorks(t *testing.T) {
	node := NewNode(nodetype.Module, 1, 1)
	assert.Nil(t, node.GetDataType())

	node.SetDataType(datatype.Boolean)
	assert.Equal(t, datatype.Boolean, *node.GetDataType())

	node.SetDataType(datatype.String)
	assert.Equal(t, datatype.String, *node.GetDataType())
}

func TestNodeGetSetValueWorks(t *testing.T) {
	node := NewNode(nodetype.Module, 1, 1)
	assert.Nil(t, node.GetValue())

	node.SetValue(true)
	assert.Equal(t, true, node.GetValue())
	assert.Equal(t, true, node.GetBooleanValue())

	node.SetValue(false)
	assert.Equal(t, false, node.GetValue())
	assert.Equal(t, false, node.GetBooleanValue())

	node.SetValue(int64(42))
	assert.Equal(t, int64(42), node.GetValue())
	assert.Equal(t, int64(42), node.GetIntegerValue())

	node.SetValue(-12.5)
	assert.Equal(t, -12.5, node.GetValue())
	assert.Equal(t, -12.5, node.GetFloatValue())

	node.SetValue("BUZZ")
	assert.Equal(t, "BUZZ", node.GetValue())
	assert.Equal(t, "BUZZ", node.GetStringValue())
}

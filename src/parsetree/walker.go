package parsetree

import (
	"fmt"
	"strings"

	"semi/src/ast"
)

type parseTreeWalker struct {
	*ast.MapWalker
}

type consoleDumpMemo struct {
	level int
}

func NewParseTreeWalker() *parseTreeWalker {
	walker := &parseTreeWalker{}

	walker.MapWalker = ast.NewMapWalker(walker.defaultVisitor)

	return walker
}

func (w *parseTreeWalker) defaultVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	args, _ := memo.(consoleDumpMemo)
	indentation := strings.Repeat("    ", args.level)
	var attrs []string

	attrs = append(attrs, fmt.Sprintf("location=\"%d:%d\"", node.GetSourceLine(), node.GetSourceColumn() + 1))

	if node.GetDataType() != nil {
		attrs = append(attrs, fmt.Sprintf("type=\"%v\"", node.GetDataType()))
	}

	if node.GetName() != "" {
		attrs = append(attrs, fmt.Sprintf("name=\"%s\"", node.GetName()))
	}

	if node.GetValue() != nil {
		attrs = append(attrs, fmt.Sprintf("value=\"%v\"", node.GetValue()))
	}

	// 	var children []Node
	//
	//	for _, nodes := range n.childNodes {
	//		children = append(children, nodes...)
	//	}

	if len(node.GetChildGroups()) > 0 {
		fmt.Printf("%s<%s%s>\n", indentation, node.GetNodeType(), fmt.Sprintf(" %s", strings.Join(attrs, " ")))

		for _, childGroup := range node.GetChildGroups() {
			groupChildren := node.GetChildrenForGroup(childGroup)

			if len(groupChildren) > 0 {
				fmt.Printf("%s    <!-- %s -->\n", indentation, childGroup)

				_, err := w.VisitNodes(groupChildren, consoleDumpMemo{level: args.level + 1})
				if err != nil {
					return nil, err
				}

				fmt.Printf("%s\n", indentation)
			}
		}

		fmt.Printf("%s</%s>\n", indentation, node.GetNodeType())
	} else {
		fmt.Printf("%s<%s%s />\n", indentation, node.GetNodeType(), fmt.Sprintf(" %s", strings.Join(attrs, " ")))
	}

	return nil, nil
}

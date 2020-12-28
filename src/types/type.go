package types

import (
	"errors"
	"fmt"

	"semi/src/ast"
	"semi/src/ast/childgroup"
	"semi/src/ast/nodetype"
	"semi/src/scopes"
)

// Types should be in the same "names" namespace.
// A type itself should be an implementation of the "Type" class.
// Built in types should be instantiated and "inserted" into the appropriate namespace.
// User-defined types should be "inserted" into the appropriate namespace. (code-gen time? transform/analysis time?)

// What about unnamed ephemeral "qualified" generic types like "Map<string, float>" used in type expressions/generic function calls?

type Type interface {
	// TODO: maybe also Name() string
	Is(other Type) bool
	AssignableTo(other Type) bool
	Parameters() []TypeParameter
	Constrain(types... Type) (Type, error)
}

type ConstraintRelationship uint8

const (
	ConstraintRelationshipIs ConstraintRelationship = iota
	ConstraintRelationshipExtends
	ConstraintRelationshipSuper
)

type TypeParameter struct {
	name string
	constraintRelationship ConstraintRelationship
	constraintType Type
}

func TypeFromNode(node ast.Node, scope scopes.Scope) (Type, error) {
	// TODO: Given a "type-expression" node and a scope (which should have named Type instances in it),
	//       lookup or construct, and return a type

	switch node.GetNodeType() {
	case nodetype.NamedTypeExpression:
		////// map<string, float>		includes: array, map, boolean, integer, float, string, etc...
		// named-type-expression
		//     name
		//       identifier
		//     arguments
		//       type-expression...

		name := node.GetChildForGroup(childgroup.Name).GetStringValue()
		argumentNodes := node.GetChildrenForGroup(childgroup.Arguments)
		arguments := make([]Type, len(argumentNodes))

		for index, argumentNode := range argumentNodes {
			argument, err := TypeFromNode(argumentNode, scope)
			if err != nil {
				return nil, err
			}

			arguments[index] = argument
		}

		found, err := scope.Lookup(name) // TODO: Pass arguments to lookup
		if err != nil {
			return nil, errors.New(fmt.Sprintf("no type with name %s found in scope", name))
		}

		if type_, ok := found.(Type); ok {
			return type_, nil
		} else {
			return nil, errors.New(fmt.Sprintf("name %s cannot be used as a type", name))
		}
	default:
		return nil, errors.New(fmt.Sprintf("expected type expression: %+v", node))
	}

	// Node structures:
	////// struct { x T; y T; }
	// struct-type-expression
	//     members
	//       struct-member-decl...
	//         name
	//           identifier
	//         type
	//           type-expression
	//
	////// func(str1 string, str2 string) (int, err)
	// func-type-expression
	//     parameters
	//       parameter...
	//         name
	//           identifier
	//         type
	//           type-expression
	//     returns
	//       type-expression...
	//
	////// interface { foo(x int); bar(); }
	// interface-type-expression
	//     members
	//       interface-member-decl...
	//         name
	//           identifier
	//         signature
	//           func-type-expression

	return nil
}

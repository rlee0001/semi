package types

import "errors"

type ArrayType struct {
	of Type
}

func (t *ArrayType) Is(other Type) bool {
	if other, ok := other.(*ArrayType); ok {
		return other.of.Is(t.of)
	}

	return false
}

func (t *ArrayType) AssignableTo(other Type) bool {
	return t.Is(other)
}

func (t *ArrayType) Parameters() []TypeParameter {
	return []TypeParameter{
		{
			name: "T",
			constraintRelationship: ConstraintRelationshipExtends,
			constraintType: nil, // TODO: Should be the root type (interface{})
		},
	}
}

func (t *ArrayType) Constrain(types... Type) (Type, error) {
	if len(types) != 1 {
		return nil, errors.New("array expects one type argument")
	}

	return &ArrayType{of: types[0]}, nil
}

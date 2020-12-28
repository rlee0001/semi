package types

import "errors"

type MapType struct {
	key Type
	value Type
}

func (t *MapType) Is(other Type) bool {
	if other, ok := other.(*MapType); ok {
		return other.key.Is(t.key) && other.value.Is(t.value)
	}

	return false
}

func (t *MapType) AssignableTo(other Type) bool {
	return t.Is(other)
}

func (t *MapType) Parameters() []TypeParameter {
	return []TypeParameter{
		{
			name: "K",
			constraintRelationship: ConstraintRelationshipExtends,
			constraintType: nil, // TODO: Should be the hashable interface type: interface{HashCode; Equals}
		},
		{
			name: "V",
			constraintRelationship: ConstraintRelationshipExtends,
			constraintType: nil, // TODO: Should be the root type: interface{}
		},
	}
}

func (t *MapType) Constrain(types... Type) (Type, error) {
	if len(types) != 2 {
		return nil, errors.New("map expects two type arguments")
	}

	return &MapType{key: types[0], value: types[1]}, nil
}

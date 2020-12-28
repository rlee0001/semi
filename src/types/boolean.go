package types

import "errors"

type BooleanType struct {
}

func (t *BooleanType) Is(other Type) bool {
	_, ok := other.(*BooleanType)

	return ok
}

func (t *BooleanType) AssignableTo(other Type) bool {
	return t.Is(other)
}

func (t *BooleanType) Parameters() []TypeParameter {
	return []TypeParameter{}
}

func (t *BooleanType) Constrain(types... Type) (Type, error) {
	if len(types) != 0 {
		return nil, errors.New("boolean doesn't expect type arguments")
	}

	return t, nil
}

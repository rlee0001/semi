package types

import "errors"

type IntegerType struct {
	// TODO: bits?
}

func (t *IntegerType) Is(other Type) bool {
	_, ok := other.(*IntegerType)

	return ok
}

func (t *IntegerType) AssignableTo(other Type) bool {
	return t.Is(other)
}

func (t *IntegerType) Parameters() []TypeParameter {
	return []TypeParameter{}
}

func (t *IntegerType) Constrain(types... Type) (Type, error) {
	if len(types) != 0 {
		return nil, errors.New("integer doesn't expect type arguments")
	}

	return t, nil
}

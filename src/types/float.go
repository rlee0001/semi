package types

import "errors"

type FloatType struct {
	// TODO: bits?
}

func (t *FloatType) Is(other Type) bool {
	_, ok := other.(*FloatType)

	return ok
}

func (t *FloatType) AssignableTo(other Type) bool {
	return t.Is(other)
}

func (t *FloatType) Parameters() []TypeParameter {
	return []TypeParameter{}
}

func (t *FloatType) Constrain(types... Type) (Type, error) {
	if len(types) != 0 {
		return nil, errors.New("float doesn't expect type arguments")
	}

	return t, nil
}

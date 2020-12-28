package types

import "errors"

type StringType struct {
	// TODO: Encoding?
}

func (t *StringType) Is(other Type) bool {
	_, ok := other.(*StringType)

	return ok
}

func (t *StringType) AssignableTo(other Type) bool {
	return t.Is(other)
}

func (t *StringType) Parameters() []TypeParameter {
	return []TypeParameter{}
}

func (t *StringType) Constrain(types... Type) (Type, error) {
	if len(types) != 0 {
		return nil, errors.New("string doesn't expect type arguments")
	}

	return t, nil
}

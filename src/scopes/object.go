package scopes

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"semi/src/types"
)

type basicFunction struct {
	irFunc *ir.Func
}

func (f *basicFunction) IrFunc() *ir.Func {
	return f.irFunc
}

func NewFunction(f *ir.Func) (*basicFunction, error) {
	return &basicFunction{
		irFunc: f,
	}, nil
}

type basicValue struct {
	type_   types.Type
	irValue value.Value
}

func (v *basicValue) Type() types.Type {
	return v.type_
}

func (v *basicValue) IrValue() value.Value {
	return v.irValue
}

func NewBooleanValue(v value.Value) (*basicValue, error) {
	// TODO: Ensure that v is of the correct llvm data type
	return &basicValue{
		type_:   &types.BooleanType{},
		irValue: v,
	}, nil
}

func NewInt64Value(v value.Value) (*basicValue, error) {
	// TODO: Ensure that v is of the correct llvm data type
	return &basicValue{
		type_:   &types.IntegerType{},
		irValue: v,
	}, nil
}

func NewFloat64Value(v value.Value) (*basicValue, error) {
	// TODO: Ensure that v is of the correct llvm data type
	return &basicValue{
		type_:   &types.FloatType{},
		irValue: v,
	}, nil
}

func NewStringValue(v value.Value) (*basicValue, error) {
	// TODO: Ensure that v is of the correct llvm data type
	return &basicValue{
		type_:   &types.StringType{},
		irValue: v,
	}, nil
}

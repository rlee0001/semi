package codegen

import (
	"errors"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	scopes2 "semi/src/scopes"

	"semi/src/ast"
	"semi/src/ast/childgroup"
	"semi/src/ast/nodetype"
)

type codeGenWalker struct {
	*ast.MapWalker

	scopes   scopes2.ScopeStack
	irModule *ir.Module
}

type functionDeclarationMemo struct {
	irFunc            *ir.Func
	irBlock           *ir.Block
}

func (m *functionDeclarationMemo) GetIrFunc() *ir.Func {
	return m.irFunc
}

func (m *functionDeclarationMemo) GetIrBlock() *ir.Block {
	return m.irBlock
}

func NewCodeGenWalker() *codeGenWalker {
	walker := &codeGenWalker{
		MapWalker: ast.NewMapWalker(nil),
		scopes:    scopes2.NewScopeStack(),
		irModule:  ir.NewModule(),
	}

	walker.scopes.Push()


	printfFormat := ir.NewParam("format", types.I8Ptr)
	printfFormat.Attrs = []ir.ParamAttribute{enum.ParamAttrNoCapture, enum.ParamAttrReadOnly}
	printf := walker.irModule.NewFunc("printf", types.I64, printfFormat)
	printf.FuncAttrs = []ir.FuncAttribute{enum.FuncAttrNoUnwind}
	printf.Sig.Variadic = true

	printfFunc, _ := scopes2.NewFunction(printf)
	_ = walker.scopes.DeclareFunction("printf", printfFunc)

	intFormat, _ := scopes2.NewStringValue(walker.irModule.NewGlobalDef("intFormat", constant.NewCharArrayFromString("%d\n\u0000")))
	floatFormat, _ := scopes2.NewStringValue(walker.irModule.NewGlobalDef("floatFormat", constant.NewCharArrayFromString("%f\n\u0000")))
	stringFormat, _ := scopes2.NewStringValue(walker.irModule.NewGlobalDef("stringFormat", constant.NewCharArrayFromString("%s\n\u0000")))
	trueString, _ := scopes2.NewStringValue(walker.irModule.NewGlobalDef("trueString", constant.NewCharArrayFromString("true\u0000")))
	falseString, _ := scopes2.NewStringValue(walker.irModule.NewGlobalDef("falseString", constant.NewCharArrayFromString("false\u0000")))

	_ = walker.scopes.DeclareValue("intFormat", intFormat)
	_ = walker.scopes.DeclareValue("floatFormat", floatFormat)
	_ = walker.scopes.DeclareValue("stringFormat", stringFormat)
	_ = walker.scopes.DeclareValue("trueString", trueString)
	_ = walker.scopes.DeclareValue("falseString", falseString)

	walker.AddVisitor(nodetype.Module, walker.moduleVisitor)
	walker.AddVisitor(nodetype.FunctionDeclaration, walker.functionDeclarationVisitor)
	walker.AddVisitor(nodetype.ParameterDeclaration, walker.parameterDeclarationVisitor)
	walker.AddVisitor(nodetype.TypeExpression, walker.typeExpressionVisitor)
	walker.AddVisitor(nodetype.LocalDeclaration, walker.localDeclarationVisitor)
	walker.AddVisitor(nodetype.Statement, walker.statementVisitor)
	walker.AddVisitor(nodetype.Identifier, walker.identifierVisitor)
	walker.AddVisitor(nodetype.Integer, walker.integerVisitor)
	walker.AddVisitor(nodetype.Float, walker.floatVisitor)
	walker.AddVisitor(nodetype.Boolean, walker.booleanVisitor)
	walker.AddVisitor(nodetype.Addition, walker.additionVisitor)
	walker.AddVisitor(nodetype.Subtraction, walker.subtractionVisitor)
	walker.AddVisitor(nodetype.Multiplication, walker.multiplicationVisitor)
	walker.AddVisitor(nodetype.Division, walker.divisionVisitor)
	walker.AddVisitor(nodetype.IsEqualTo, walker.isEqualToVisitor)

	return walker
}

func (w *codeGenWalker) getIrModule() *ir.Module {
	return w.irModule
}

func (w *codeGenWalker) moduleVisitor(node ast.Node, _ interface{}) (interface{}, error) {
	_, err := w.VisitNodes(node.GetChildrenForGroup(childgroup.FunctionDeclarations), nil)

	return nil, err
}

func (w *codeGenWalker) functionDeclarationVisitor(node ast.Node, _ interface{}) (interface{}, error) {
	// TODO: Add parameters and return type to function from attributes and children
	name := node.GetName()
	irFunc := w.getIrModule().NewFunc(name, types.I64)
	irBlock := irFunc.NewBlock("")

	fnMemo := &functionDeclarationMemo{
		irFunc:  irFunc,
		irBlock: irBlock,
	}

	_, err := w.VisitNodes(node.GetChildrenForGroup(childgroup.LocalDeclarations), fnMemo)
	if err != nil {
		return nil, err
	}

	_, err = w.VisitNodes(node.GetChildrenForGroup(childgroup.Statements), fnMemo)
	if err != nil {
		return nil, err
	}

	fnMemo.irBlock.NewRet(constant.NewInt(types.I64, 0))

	return nil, nil
}

func (w *codeGenWalker) parameterDeclarationVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	// TODO: NYI
	return nil, nil
}

func (w *codeGenWalker) typeExpressionVisitor(node ast.Node, _ interface{}) (interface{}, error) {
	// TODO: Get datatype.DataType from Object?
	return node.GetName(), nil
}

func (w *codeGenWalker) localDeclarationVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	arg, _ := memo.(*functionDeclarationMemo)

	value, err := w.VisitNode(node.GetChildForGroup(childgroup.Initializer), arg)
	if err != nil {
		return nil, err
	}

	name := node.GetName()
	initializerValue, _ := value.(*scopes2.basicValue)
	irBlock := arg.GetIrBlock()
	irAllocation := irBlock.NewAlloca(initializerValue.irValue.Type())
	irAllocation.SetName(name)
	irBlock.NewStore(initializerValue.irValue, irAllocation)

	// TODO: This switch should be a helper: NewValue(value.Value, valueType)
	if initializerValue.type_ == scopes2.valueTypeInt32 {
		typedValue, _ := scopes2.NewInt64Value(irAllocation)

		_ = w.scopes.DeclareValue(name, typedValue)
	} else if initializerValue.type_ == scopes2.valueTypeFloat64 {
		typedValue, _ := scopes2.NewFloat64Value(irAllocation)

		_ = w.scopes.DeclareValue(name, typedValue)
	} else if initializerValue.type_ == scopes2.valueTypeBoolean {
		typedValue, _ := scopes2.NewInt64Value(irAllocation)

		_ = w.scopes.DeclareValue(name, typedValue)
	} else if initializerValue.type_ == scopes2.valueTypeString {
		typedValue, _ := scopes2.NewInt64Value(irAllocation)

		_ = w.scopes.DeclareValue(name, typedValue)
	}

	return nil, nil
}

func (w *codeGenWalker) statementVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	arg, _ := memo.(*functionDeclarationMemo)

	value, err := w.VisitNode(node.GetChildForGroup(childgroup.Arguments), arg)
	if err != nil {
		return nil, err
	}

	typedValue, _ := value.(*scopes2.basicValue)

	printf, _ := w.scopes.Lookup("printf")
	irFloatFormat, _ := w.scopes.Lookup("floatFormat")
	irIntFormat, _ := w.scopes.Lookup("intFormat")
	irStringFormat, _ := w.scopes.Lookup("stringFormat")
	irTrueString, _ := w.scopes.Lookup("trueString")
	irFalseString, _ := w.scopes.Lookup("falseString")
	irFunc := arg.GetIrFunc()
	irBlock := arg.GetIrBlock()
	irZero := constant.NewInt(types.I64, 0)

	if typedValue.type_ == scopes2.valueTypeFloat64 {
		floatingValue := irBlock.NewFPExt(typedValue.irValue, types.Double)
		formatPtr := irBlock.NewGetElementPtr(irFloatFormat.(scopes2.Value).IrValue().(*ir.Global).ContentType, irFloatFormat.(scopes2.Value).IrValue(), irZero, irZero)
		irBlock.NewCall(printf.(scopes2.Function).IrFunc(), formatPtr, floatingValue)
	} else if typedValue.type_ == scopes2.valueTypeInt32 {
		formatPtr := irBlock.NewGetElementPtr(irIntFormat.(scopes2.Value).IrValue().(*ir.Global).ContentType, irIntFormat.(scopes2.Value).IrValue(), irZero, irZero)
		irBlock.NewCall(printf.(scopes2.Function).IrFunc(), formatPtr, typedValue.irValue)
	} else if typedValue.type_ == scopes2.valueTypeBoolean {
		trueBlock := irFunc.NewBlock("")
		falseBlock := irFunc.NewBlock("")
		contBlock := irFunc.NewBlock("")

		irBlock.NewCondBr(typedValue.irValue, trueBlock, falseBlock)

		trueFormatPtr := trueBlock.NewGetElementPtr(irStringFormat.(scopes2.Value).IrValue().(*ir.Global).ContentType, irStringFormat.(scopes2.Value).IrValue(), irZero, irZero)
		truePtr := trueBlock.NewGetElementPtr(irTrueString.(scopes2.Value).IrValue().(*ir.Global).ContentType, irTrueString.(scopes2.Value).IrValue(), irZero, irZero)
		trueBlock.NewCall(printf.(scopes2.Function).IrFunc(), trueFormatPtr, truePtr)
		trueBlock.NewBr(contBlock)

		falseFormatPtr := falseBlock.NewGetElementPtr(irStringFormat.(scopes2.Value).IrValue().(*ir.Global).ContentType, irStringFormat.(scopes2.Value).IrValue(), irZero, irZero)
		falsePtr := falseBlock.NewGetElementPtr(irFalseString.(scopes2.Value).IrValue().(*ir.Global).ContentType, irFalseString.(scopes2.Value).IrValue(), irZero, irZero)
		falseBlock.NewCall(printf.(scopes2.Function).IrFunc(), falseFormatPtr, falsePtr)
		falseBlock.NewBr(contBlock)

		arg.irBlock = contBlock
	}

	return nil, nil
}

func (w *codeGenWalker) identifierVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	arg, _ := memo.(*functionDeclarationMemo)

	typedValue, err := w.scopes.Lookup(node.GetName())
	if err != nil {
		return nil, err
	}

	irBaseType := typedValue.(scopes2.Value).IrValue().Type()

	if irPtrType, ok := irBaseType.(*types.PointerType); ok {
		irBaseType = irPtrType.ElemType
	}

	irLoad := arg.GetIrBlock().NewLoad(irBaseType, typedValue.(scopes2.Value).IrValue())

	if typedValue.(scopes2.Value).Type() == scopes2.valueTypeInt32 {
		return scopes2.NewInt64Value(irLoad)
	} else if typedValue.(scopes2.Value).Type() == scopes2.valueTypeFloat64 {
		return scopes2.NewFloat64Value(irLoad)
	} else if typedValue.(scopes2.Value).Type() == scopes2.valueTypeBoolean {
		return scopes2.NewBooleanValue(irLoad)
	} else if typedValue.(scopes2.Value).Type() == scopes2.valueTypeString {
		return scopes2.NewStringValue(irLoad)
	}

	return nil, errors.New("identifier with unrecognized type")
}

func (w *codeGenWalker) integerVisitor(node ast.Node, _ interface{}) (interface{}, error) {
	return scopes2.NewInt64Value(constant.NewInt(types.I64, node.GetIntegerValue()))
}

func (w *codeGenWalker) floatVisitor(node ast.Node, _ interface{}) (interface{}, error) {
	return scopes2.NewFloat64Value(constant.NewFloat(types.Float, node.GetFloatValue()))
}

func (w *codeGenWalker) booleanVisitor(node ast.Node, _ interface{}) (interface{}, error) {
	var integerValue int64

	if node.GetBooleanValue() {
		integerValue = 1
	}

	return scopes2.NewBooleanValue(constant.NewInt(types.I1, integerValue))
}

func (w *codeGenWalker) additionVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	arg, _ := memo.(*functionDeclarationMemo)

	rhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Rhs), arg)
	if err != nil {
		return nil, err
	}

	lhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Lhs), arg)
	if err != nil {
		return nil, err
	}

	rhs, _ := rhsValue.(*scopes2.basicValue)
	lhs, _ := lhsValue.(*scopes2.basicValue)

	irBlock := arg.GetIrBlock()
	operatorBasicType := scopes2.valueTypeInt32

	// If either lhs or rhs is a float, cast the other to a float and emit an FAdd, otherwise emit an Add.
	if rhs.type_ == scopes2.valueTypeFloat64 {
		if lhs.type_ == scopes2.valueTypeInt32 {
			lhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(lhs.irValue, types.Float))
		}
		operatorBasicType = scopes2.valueTypeFloat64
	} else if lhs.type_ == scopes2.valueTypeFloat64 {
		rhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(rhs.irValue, types.Float))
		operatorBasicType = scopes2.valueTypeFloat64
	}

	if operatorBasicType == scopes2.valueTypeFloat64 {
		return scopes2.NewFloat64Value(irBlock.NewFAdd(lhs.irValue, rhs.irValue))
	} else {
		return scopes2.NewInt64Value(irBlock.NewAdd(lhs.irValue, rhs.irValue))
	}
}

func (w *codeGenWalker) subtractionVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	arg, _ := memo.(*functionDeclarationMemo)

	rhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Rhs), arg)
	if err != nil {
		return nil, err
	}

	lhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Lhs), arg)
	if err != nil {
		return nil, err
	}

	rhs, _ := rhsValue.(*scopes2.basicValue)
	lhs, _ := lhsValue.(*scopes2.basicValue)

	irBlock := arg.GetIrBlock()
	operatorBasicType := scopes2.valueTypeInt32

	// If either lhs or rhs is a float, cast the other to a float and emit an FSub, otherwise emit a Sub.
	if rhs.type_ == scopes2.valueTypeFloat64 {
		if lhs.type_ == scopes2.valueTypeInt32 {
			lhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(lhs.irValue, types.Float))
		}
		operatorBasicType = scopes2.valueTypeFloat64
	} else if lhs.type_ == scopes2.valueTypeFloat64 {
		rhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(rhs.irValue, types.Float))
		operatorBasicType = scopes2.valueTypeFloat64
	}

	if operatorBasicType == scopes2.valueTypeFloat64 {
		return scopes2.NewFloat64Value(irBlock.NewFSub(lhs.irValue, rhs.irValue))
	} else {
		return scopes2.NewInt64Value(irBlock.NewSub(lhs.irValue, rhs.irValue))
	}
}

func (w *codeGenWalker) multiplicationVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	arg, _ := memo.(*functionDeclarationMemo)

	rhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Rhs), arg)
	if err != nil {
		return nil, err
	}

	lhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Lhs), arg)
	if err != nil {
		return nil, err
	}

	rhs, _ := rhsValue.(*scopes2.basicValue)
	lhs, _ := lhsValue.(*scopes2.basicValue)

	irBlock := arg.GetIrBlock()
	operatorBasicType := scopes2.valueTypeInt32

	// If either lhs or rhs is a float, cast the other to a float and emit an FMul, otherwise emit a Mul.
	if rhs.type_ == scopes2.valueTypeFloat64 {
		if lhs.type_ == scopes2.valueTypeInt32 {
			lhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(lhs.irValue, types.Float))
		}
		operatorBasicType = scopes2.valueTypeFloat64
	} else if lhs.type_ == scopes2.valueTypeFloat64 {
		rhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(rhs.irValue, types.Float))
		operatorBasicType = scopes2.valueTypeFloat64
	}

	if operatorBasicType == scopes2.valueTypeFloat64 {
		return scopes2.NewFloat64Value(irBlock.NewFMul(lhs.irValue, rhs.irValue))
	} else {
		return scopes2.NewInt64Value(irBlock.NewMul(lhs.irValue, rhs.irValue))
	}
}

func (w *codeGenWalker) divisionVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	arg, _ := memo.(*functionDeclarationMemo)

	rhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Rhs), arg)
	if err != nil {
		return nil, err
	}

	lhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Lhs), arg)
	if err != nil {
		return nil, err
	}

	rhs, _ := rhsValue.(*scopes2.basicValue)
	lhs, _ := lhsValue.(*scopes2.basicValue)

	irBlock := arg.GetIrBlock()
	operatorBasicType := scopes2.valueTypeInt32

	// If either lhs or rhs is a float, cast the other to a float and emit an FDiv, otherwise emit a SDiv.
	if rhs.type_ == scopes2.valueTypeFloat64 {
		if lhs.type_ == scopes2.valueTypeInt32 {
			lhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(lhs.irValue, types.Float))
		}
		operatorBasicType = scopes2.valueTypeFloat64
	} else if lhs.type_ == scopes2.valueTypeFloat64 {
		rhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(rhs.irValue, types.Float))
		operatorBasicType = scopes2.valueTypeFloat64
	}

	if operatorBasicType == scopes2.valueTypeFloat64 {
		return scopes2.NewFloat64Value(irBlock.NewFDiv(lhs.irValue, rhs.irValue))
	} else {
		return scopes2.NewInt64Value(irBlock.NewSDiv(lhs.irValue, rhs.irValue))
	}
}

func (w *codeGenWalker) isEqualToVisitor(node ast.Node, memo interface{}) (interface{}, error) {
	arg, _ := memo.(*functionDeclarationMemo)

	rhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Rhs), arg)
	if err != nil {
		return nil, err
	}

	lhsValue, err := w.VisitNode(node.GetChildForGroup(childgroup.Lhs), arg)
	if err != nil {
		return nil, err
	}

	rhs, _ := rhsValue.(*scopes2.basicValue)
	lhs, _ := lhsValue.(*scopes2.basicValue)

	irBlock := arg.GetIrBlock()
	operatorBasicType := scopes2.valueTypeInt32

	if rhs.type_ == scopes2.valueTypeFloat64 {
		if lhs.type_ == scopes2.valueTypeInt32 {
			lhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(lhs.irValue, types.Float))
		}
		operatorBasicType = scopes2.valueTypeFloat64
	} else if lhs.type_ == scopes2.valueTypeFloat64 {
		rhs, _ = scopes2.NewFloat64Value(irBlock.NewSIToFP(rhs.irValue, types.Float))
		operatorBasicType = scopes2.valueTypeFloat64
	}

	if operatorBasicType == scopes2.valueTypeFloat64 {
		return scopes2.NewBooleanValue(irBlock.NewFCmp(enum.FPredOEQ,lhs.irValue, rhs.irValue))
	} else {
		return scopes2.NewBooleanValue(irBlock.NewICmp(enum.IPredEQ, lhs.irValue, rhs.irValue))
	}
}

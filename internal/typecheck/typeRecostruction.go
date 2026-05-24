package typecheck

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/typecheck/constraint"
)

func reconstruct(ctx *Context, node nodes.Node) (type_ nodes.StellaType, err *TypecheckError) {
	switch v := node.(type) {
	case *nodes.FunctionDeclaration:
		// If successful add function to upper level scope
		defer func() {
			if err != nil {
				err.AddIfEmptyFunctionName(v.Name)
			}
		}()

		// function top-level scope
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		// // Add current function for recursion
		// freshVar := ctx.constraintHandler.GetFreshVar()
		// ctx.AddVar(v.Name, &freshVar)
		// defer ctx.RemoveLastScope()

		// function parameters scope
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		params, err := addParametersToContext(ctx, v.Params)
		if err != nil {
			err.AddIfEmptyExpr(v)
			return nil, err
		}

		// subdeclarations scope
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		// Collect all names
		for _, declaration := range v.Declarations {

			switch decl := declaration.(type) {
			case *nodes.FunctionDeclaration:
				type_, err := constructTypeFromDeclaration(ctx, &declaration)

				if err != nil {
					return nil, err
				}

				success := ctx.AddVar(decl.Name, type_)

				if !success {
					err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_FUNCTION_DECLARATION)
					err.AddIfEmptyFunctionName(v.Name)
					// err.AddIfEmptyExpr(decl)
					return nil, &err
				}
			case *nodes.ExceptionTypeDeclaration:
				err := NewTypeCheckErrorErrorType(ERROR_ILLEGAL_LOCAL_EXCEPTION_TYPE)
				err.AddAdditionalInfo("Exception declarations permitted only in global scope")
				err.AddIfEmptyFunctionName(v.Name)
				err.AddIfEmptyExpr(decl)
				err.Freeze()
				return nil, &err
			case *nodes.ExceptionVariantDeclaration:
				err := NewTypeCheckErrorErrorType(ERROR_ILLEGAL_LOCAL_OPEN_VARIANT_EXCEPTION)
				err.AddAdditionalInfo("Exception declarations permitted only in global scope")
				err.AddIfEmptyFunctionName(v.Name)
				err.AddIfEmptyExpr(decl)
				err.Freeze()
				return nil, &err
			default:
				err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
				err.AddIfEmptyExpr(decl)
				err.AddAdditionalInfo(fmt.Sprintf("Not implemented check type function declaration switch for %s", decl))
				return nil, &err
			}
		}

		// Check types
		for _, decl := range v.Declarations {
			switch decl := decl.(type) {
			case *nodes.FunctionDeclaration:
				type_ := ctx.GetVarType(decl.Name).Require()

				inferredType, err := reconstruct(ctx, decl)

				if err != nil {
					return nil, err
				}

				ctx.constraintHandler.AddEquation(
					constraint.NewEquation(inferredType, type_, decl),
				)
			}
		}

		// body scope
		ctx.AddNewScope()

		retType, err := transformTypeAuto(ctx, v.ReturnType.OrElse(&nodes.TypeUnit{}))
		if err != nil {
			return nil, err
		}

		inferredType, err := reconstruct(ctx, v.Expr)

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredType, retType, v),
		)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeFun{ParamTypes: getTypeFromParameters(params), ReturnType: inferredType}, nil
	case *nodes.ConstUnit:
		return &nodes.TypeUnit{}, nil
	case *nodes.ConstBool:
		return &nodes.TypeBool{}, nil
	case *nodes.ConstInt:
		return &nodes.TypeNat{}, nil
	case *nodes.Var:
		type_ := ctx.GetVarType(v.Name)

		if type_.IsEmpty() {
			var err = NewTypeCheckErrorErrorType(ERROR_UNDEFINED_VARIABLE)
			err.AddIfEmptyExpr(&v.Name)
			return nil, &err
		}
		return type_.Require(), nil
	case *nodes.If:
		condType, err := reconstruct(ctx, v.Condition)

		if err != nil {
			return nil, err
		}

		inferredType1, err := reconstruct(ctx, v.ThenExpr)

		if err != nil {
			return nil, err
		}

		inferredType2, err := reconstruct(ctx, v.ElseExpr)

		if err != nil {
			return nil, err
		}

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(condType, &nodes.TypeBool{}, v.Condition),
		)

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredType1, inferredType2, v),
		)

		return inferredType1, nil
	case *nodes.Succ:
		natType := nodes.TypeNat{}
		inferredType, err := reconstruct(ctx, v.N)

		if err != nil {
			err.AddIfEmptyExpr(node)
			return nil, err
		}

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredType, &natType, v.N),
		)

		return &natType, nil
	case *nodes.IsZero:
		natType := nodes.TypeNat{}
		inferredType, err := reconstruct(ctx, v.N)

		if err != nil {
			err.AddIfEmptyExpr(node)
			return nil, err
		}

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredType, &natType, v.N),
		)

		return &nodes.TypeBool{}, nil
	case *nodes.NatRec:
		natType := nodes.TypeNat{}
		nType, err := reconstruct(ctx, v.N)

		if err != nil {
			err.AddIfEmptyExpr(node)
			return nil, err
		}

		inferredType, err := reconstruct(ctx, v.Initial)

		if err != nil {
			return nil, err
		}

		stepSubType := nodes.TypeFun{ParamTypes: []nodes.StellaType{inferredType}, ReturnType: inferredType}
		expectedStepType := nodes.TypeFun{ParamTypes: []nodes.StellaType{&natType}, ReturnType: &stepSubType}

		inferredStepType, err := reconstruct(ctx, v.Step)

		if err != nil {
			return nil, err
		}

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(nType, &natType, v.N),
		)

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredStepType, &expectedStepType, v.Step),
		)

		return inferredType, nil
	case *nodes.Abstraction:
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		params, err := addParametersToContext(ctx, v.Params)

		if err != nil {
			err.AddIfEmptyExpr(v)
			return nil, err
		}

		paramTypes := getTypeFromParameters(params)

		inferredType, err := reconstruct(ctx, v.ReturnExpr)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeFun{ParamTypes: paramTypes, ReturnType: inferredType}, nil
	case *nodes.Application:
		inferredType, err := reconstruct(ctx, v.Function)

		if err != nil {
			return nil, err
		}

		inferredArgParamType := make([]nodes.StellaType, 0, len(v.Args))
		for _, arg := range v.Args {
			inferredArgType, err := reconstruct(ctx, arg)

			if err != nil {
				return nil, err
			}

			inferredArgParamType = append(inferredArgParamType, inferredArgType)
		}

		// get new fresh var
		newFresh := ctx.constraintHandler.GetFreshVar()

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(
				inferredType,
				&nodes.TypeFun{ParamTypes: inferredArgParamType, ReturnType: &newFresh},
				v,
			),
		)

		if inferredFunType, ok := inferredType.(*nodes.TypeFun); ok {
			return inferredFunType.ReturnType, nil
		}

		return &newFresh, nil
	case *nodes.Tuple:
		paramTypes := make([]nodes.StellaType, len(v.Exprs))

		for index, expr := range v.Exprs {
			type_, err := reconstruct(ctx, expr)

			if err != nil {
				return nil, err
			}

			paramTypes[index] = type_
		}

		return &nodes.TypeTuple{Types: paramTypes}, nil
	case *nodes.DotTuple:
		if v.Index > 2 || v.Index < 1 {
			err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
			err.AddAdditionalInfo("Type recostruction for tuples are not implemented")
			return nil, &err
		}
		inferredType, err := reconstruct(ctx, v.Subexpr)

		if err != nil {
			return nil, err
		}

		freshVar1 := ctx.constraintHandler.GetFreshVar()
		freshVar2 := ctx.constraintHandler.GetFreshVar()

		newTypes := make([]nodes.StellaType, 2)
		newTypes[0] = &freshVar1
		newTypes[1] = &freshVar2

		ctx.constraintHandler.AddEquationTypes(
			inferredType, &nodes.TypeTuple{Types: newTypes},
		)

		if v.Index == 1 {
			return &freshVar1, nil
		}

		return &freshVar2, nil

	case *nodes.Record:
		recordFields := make(map[nodes.StellaIdent]bool, len(v.Bindings))
		recordTypes := make([]nodes.RecordFieldType, len(v.Bindings))

		for index, binding := range v.Bindings {
			_, ok := recordFields[binding.Name]

			if ok {
				err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_RECORD_FIELDS)
				err.AddIfEmptyExpr(v)
				err.AddAdditionalInfo(fmt.Sprintf("Duplicate label: %s", &binding.Name))
				return nil, &err
			}

			recordFields[binding.Name] = true

			inferredType, err := reconstruct(ctx, binding.Rhs)

			if err != nil {
				return nil, err
			}

			recordTypes[index] = nodes.RecordFieldType{Label: binding.Name, Type_: inferredType}
		}

		return &nodes.TypeRecord{FieldTypes: recordTypes}, nil
	case *nodes.DotRecord:
		inferredType, err := reconstruct(ctx, v.Subexpr)

		if err != nil {
			return nil, err
		}

		if _, ok := inferredType.(*nodes.TypeRecord); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_RECORD)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a record", v.Subexpr))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return nil, &err
		}

		recordType, _ := inferredType.(*nodes.TypeRecord)

		for _, typeBinding := range recordType.FieldTypes {
			if typeBinding.Label.Name == v.Label.Name {
				return typeBinding.Type_, nil
			}
		}

		err_ := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_FIELD_ACCESS)
		err_.AddIfEmptyExpr(v)
		err_.AddAdditionalInfo(fmt.Sprintf("Unexpected label: %s", v.Label.String()))
		return nil, &err_
	case *nodes.TypeAsc:
		err := checkTypeConsistency(v.Type_)

		if err != nil {
			return nil, err
		}

		type_, err := reconstruct(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(type_, v.Type_, v.Expr_),
		)

		return v.Type_, nil
	case *nodes.Inl:
		leftType, err := reconstruct(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}
		freshVar := ctx.constraintHandler.GetFreshVar()
		return &nodes.TypeSum{Left: leftType, Right: &freshVar}, nil
	case *nodes.Inr:
		rightType, err := reconstruct(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}
		freshVar := ctx.constraintHandler.GetFreshVar()
		return &nodes.TypeSum{Left: &freshVar, Right: rightType}, nil
	case *nodes.List:
		// fresh var for type
		newFresh := ctx.constraintHandler.GetFreshVar()

		for _, expr := range v.Exprs {
			inferredType, err := reconstruct(ctx, expr)

			if err != nil {
				return nil, err
			}

			ctx.constraintHandler.AddEquation(
				constraint.NewEquation(&newFresh, inferredType, v),
			)
		}

		return &nodes.TypeList{Type_: &newFresh}, nil
	case *nodes.ConsList:
		headType, err := reconstruct(ctx, v.Head)

		if err != nil {
			return nil, err
		}

		listType := nodes.TypeList{Type_: headType}
		inferredTail, err := reconstruct(ctx, v.Tail)

		if err != nil {
			return nil, err
		}

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredTail, &listType, v),
		)

		return &listType, nil
	case *nodes.Head:
		inferredType, err := reconstruct(ctx, v.List)

		if err != nil {
			return nil, err
		}

		if listType, ok := inferredType.(*nodes.TypeList); ok {
			return listType.Type_, nil
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_LIST)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a list", v.List))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return nil, &err
		}
	case *nodes.Tail:
		inferredType, err := reconstruct(ctx, v.List)

		if err != nil {
			return nil, err
		}

		if listType, ok := inferredType.(*nodes.TypeList); ok {
			return listType, nil
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_LIST)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a list", v.List))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return nil, &err
		}
	case *nodes.IsEmpty:
		inferredType, err := reconstruct(ctx, v.List)

		if err != nil {
			return nil, err
		}

		if _, ok := inferredType.(*nodes.TypeList); ok {
			return &nodes.TypeBool{}, nil
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_LIST)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a list", v.List))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return nil, &err
		}
	case *nodes.Variant:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_VARIANT_TYPE)
		err.AddIfEmptyExpr(v)
		err.Freeze()
		return nil, &err
	case *nodes.Fix:
		inferredType, err := reconstruct(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		if funType, ok := inferredType.(*nodes.TypeFun); ok {

			if len(funType.ParamTypes) != 1 {
				err_ := NewTypeCheckErrorErrorType(ERROR_INCORRECT_NUMBER_OF_ARGUMENTS)
				err_.AddIfEmptyExpr(v.Expr_)
				err_.AddAdditionalInfo(fmt.Sprintf("Expected 1. Got %d", len(funType.ParamTypes)))
				err_.AddIfEmptyExpr(v)
				err_.Freeze()
				return nil, &err_
			}

			ctx.constraintHandler.AddEquation(
				constraint.NewEquation(funType.ParamTypes[0], funType.ReturnType, v.Expr_),
			)

			return funType.ReturnType, nil
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_FUNCTION)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a function", v.Expr_))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return nil, &err
		}
	case *nodes.Match:
		if len(v.Cases) == 0 {
			err := NewTypeCheckErrorErrorType(ERROR_ILLEGAL_EMPTY_MATCHING)
			err.AddIfEmptyExpr(node)
			return nil, &err
		}

		inferredType, err := reconstruct(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		err = checkPatternTypes(ctx, v.Cases, inferredType)
		if err != nil {
			return nil, err
		}

		err = checkExhaustiveness(v.Cases, inferredType)
		if err != nil {
			err.AddIfEmptyExpr(v)
			return nil, err
		}

		// Checking type
		inferredCaseType := ctx.constraintHandler.GetFreshVar()
		for _, case_ := range v.Cases {
			ctx.AddNewScope()
			patternToContext(ctx, case_.Pattern, inferredType)

			inferredType, err := reconstruct(ctx, case_.Expr_)

			if err != nil {
				ctx.RemoveLastScope()
				return nil, err
			}

			ctx.constraintHandler.AddEquation(
				constraint.NewEquation(&inferredCaseType, inferredType, case_.Expr_),
			)

			ctx.RemoveLastScope()
		}

		return &inferredCaseType, nil

	case *nodes.Let:
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()
		err := patternBindingsToContext(ctx, v.PatternBindings)
		if err != nil {
			if err.errorType == ERROR_DUPLICATE_LET_BINDING {
				err.AddIfEmptyExpr(node)
			}
			return nil, err
		}
		return reconstruct(ctx, v.Body)

	case *nodes.Sequence:
		inferredType, err := reconstruct(ctx, v.Expr1)

		if err != nil {
			return nil, err
		}

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredType, &nodes.TypeUnit{}, v.Expr1),
		)

		return reconstruct(ctx, v.Expr2)
	case *nodes.TerminatingSemicolon:
		return reconstruct(ctx, v.Expr_)
	case *nodes.ParenthesisedExpr:
		return reconstruct(ctx, v.Expr_)
	case *nodes.Ref:
		inferredType, err := reconstruct(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeRef{Type_: inferredType}, nil
	case *nodes.Deref:
		inferredType, err := reconstruct(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		freshVar := ctx.constraintHandler.GetFreshVar()

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredType, &nodes.TypeRef{Type_: &freshVar}, v.Expr_),
		)

		return &freshVar, nil
	case *nodes.Assign:
		lhs, err := reconstruct(ctx, v.Lhs)

		if err != nil {
			return nil, err
		}

		rhs, err := reconstruct(ctx, v.Rhs)

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(lhs, &nodes.TypeRef{Type_: rhs}, v),
		)

		return &nodes.TypeUnit{}, nil
	case *nodes.ConstMemory:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_REFERENCE_TYPE)
		err.AddIfEmptyExpr(v)
		err.Freeze()
		return nil, &err
	case *nodes.TypeCast:
		inferredType, err := reconstruct(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		ctx.constraintHandler.AddEquation(
			constraint.NewEquation(inferredType, v.Type_, v.Expr_),
		)

		return inferredType, nil
	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		err.AddAdditionalInfo(fmt.Sprintf("Not implemented type inference for %s", node))
		return nil, &err
	}
}

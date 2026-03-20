package typecheck

import (
	nodes "typechecker/internal/ast/nodes"
)

func infer(ctx *Context, node nodes.Node) (nodes.StellaType, *TypecheckError) {
	switch v := node.(type) {
	case *nodes.AProgram:
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		hasMain := false

		// Collect all declarations to context
		for _, declaration := range v.Declarations {
			declType, err := constructTypeFromDeclaration(&declaration)
			name := v.LanguageDecl.Name
			nameStella := nodes.StellaIdent{Name: name, Repr: name}

			if err != nil {
			}
			success := ctx.AddVar(nameStella, declType)

			if !success {
				var typeError = NewTypeCheckErrorErrorType(ERROR_DUPLICATE_FUNCTION_DECLARATION)
				typeError = typeError.AddIfEmptyFunctionName(nameStella).AddIfEmptyExpr(v)
				return nil, &typeError
			}
			hasMain = hasMain || (name == "main")
		}

		// Check declaration bodies
		for _, declaration := range v.Declarations {
			type_, _ := constructTypeFromDeclaration(&declaration)
			err := CheckType(ctx, declaration, type_)

			if err != nil {
				return nil, err
			}
		}

		if !hasMain {
			var typeError = NewTypeCheckErrorErrorType(ERROR_MISSING_MAIN)
			return nil, &typeError
		}

		return nil, nil

	case *nodes.ConstUnit:
		return &nodes.TypeUnit{Repr: "unit"}, nil
	case *nodes.ConstBool:
		return &nodes.TypeBool{Repr: "bool"}, nil
	case *nodes.ConstInt:
		return &nodes.TypeNat{Repr: "nat"}, nil
	case *nodes.Var:
		type_ := ctx.GetVarType(v.Name)

		if type_.IsEmpty() {
			var err = NewTypeCheckErrorErrorType(ERROR_UNDEFINED_VARIABLE).AddIfEmptyExpr(&v.Name)
			return nil, &err
		}
		return type_.Require(), nil
	case *nodes.If:
		err := CheckType(ctx, v.Condition, &nodes.TypeBool{Repr: "bool"})

		if err != nil {
			return nil, err
		}

		inferredType, err := infer(ctx, v.ThenExpr)

		if err != nil {
			return nil, err
		}

		err = CheckType(ctx, v.ElseExpr, inferredType)

		if err != nil {
			return nil, err
		}

		return inferredType, nil
	case *nodes.Succ:
		natType := nodes.TypeNat{Repr: "nat"}
		err := CheckType(ctx, v.N, &natType)

		if err != nil {
			err.AddIfEmptyExpr(node)
			return nil, err
		}

		return &natType, nil
	case *nodes.IsZero:
		natType := nodes.TypeNat{Repr: "nat"}
		err := CheckType(ctx, v.N, &natType)

		if err != nil {
			err.AddIfEmptyExpr(node)
			return nil, err
		}

		return &nodes.TypeBool{Repr: "bool"}, nil
	case *nodes.NatRec:
		natType := nodes.TypeNat{Repr: "nat"}
		err := CheckType(ctx, v.N, &natType)

		if err != nil {
			err.AddIfEmptyExpr(node)
			return nil, err
		}

		inferredType, err := infer(ctx, v.Initial)

		if err != nil {
			return nil, err
		}

		stepSubType := nodes.TypeFun{ParamTypes: []nodes.StellaType{inferredType}, ReturnType: inferredType}
		stepType := nodes.TypeFun{ParamTypes: []nodes.StellaType{&natType}, ReturnType: &stepSubType}

		err = CheckType(ctx, v.Step, &stepType)

		if err != nil {
			return nil, err
		}
		return inferredType, nil
	case *nodes.Abstraction:
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		err := addParametersToContext(ctx, v.Params)

		if err != nil {
			return nil, err
		}

		paramTypes := getTypeFromParameters(v.Params)

		inferredType, err := infer(ctx, v.ReturnExpr)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeFun{ParamTypes: paramTypes, ReturnType: inferredType, Repr: ""}, nil
	case *nodes.Application:
		inferredType, err := infer(ctx, v.Function)

		if err != nil {
			return nil, err
		}

		if _, ok := inferredType.(*nodes.TypeFun); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_FUNCTION)
			err = err.AddIfEmptyExpr(v.Function).AddIfEmptyActualType(inferredType)
			return nil, &err
		}

		inferredFunType, _ := inferredType.(*nodes.TypeFun)

		for _, pair := range Zip(v.Args, inferredFunType.ParamTypes) {
			expr := pair.First
			type_ := pair.Second

			err = CheckType(ctx, expr, type_)

			if err != nil {
				return nil, err
			}
		}

		return inferredFunType.ReturnType, nil
	case *nodes.Tuple:
		paramTypes := make([]nodes.StellaType, len(v.Exprs))

		for index, expr := range v.Exprs {
			type_, err := infer(ctx, expr)

			if err != nil {
				return nil, err
			}

			paramTypes[index] = type_
		}

		return &nodes.TypeTuple{Types: paramTypes}, nil
	case *nodes.DotTuple:
		inferredType, err := infer(ctx, v.Subexpr)

		if err != nil {
			return nil, err
		}

		if _, ok := inferredType.(*nodes.TypeTuple); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_TUPLE)
			err = err.AddIfEmptyExpr(v.Subexpr).AddIfEmptyActualType(inferredType)
			return nil, &err
		}

		tupleType, _ := inferredType.(*nodes.TypeTuple)

		if v.Index >= len(tupleType.Types) {
			// TODO show actual and expected tuple length
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TUPLE_LENGTH).AddIfEmptyExpr(v)
			return nil, &err
		}

	case *nodes.Record:
		recordFields := make(map[nodes.StellaIdent]bool, len(v.Bindings))
		recordTypes := make([]nodes.RecordFieldType, len(v.Bindings))

		for index, binding := range v.Bindings {
			_, ok := recordFields[binding.Name]

			if ok {
				// TODO show duplicated label
				err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_RECORD_FIELDS).AddIfEmptyExpr(v)
				return nil, &err
			}

			recordFields[binding.Name] = true

			inferredType, err := infer(ctx, binding.Rhs)

			if err != nil {
				return nil, err
			}

			recordTypes[index] = nodes.RecordFieldType{Label: binding.Name, Type_: inferredType, Repr: ""}
		}

		return &nodes.TypeRecord{FieldTypes: recordTypes}, nil
	case *nodes.DotRecord:
		inferredType, err := infer(ctx, v.Subexpr)

		if err != nil {
			return nil, err
		}

		if _, ok := inferredType.(*nodes.TypeRecord); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_RECORD)
			err = err.AddIfEmptyExpr(v.Subexpr).AddIfEmptyActualType(inferredType)
			return nil, &err
		}

		recordType, _ := inferredType.(*nodes.TypeRecord)

		for _, typeBinding := range recordType.FieldTypes {
			if typeBinding.Label.Name == v.Label.Name {
				return typeBinding.Type_, nil
			}
		}

		// TODO show expected label
		err_ := NewTypeCheckErrorErrorType(ERROR_MISSING_RECORD_FIELDS)
		err_ = err_.AddIfEmptyExpr(v)
		return nil, &err_
	case *nodes.TypeAsc:
		return v.Type_, CheckType(ctx, v.Expr_, v.Type_)
	case *nodes.Inl:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_SUM_TYPE)
		err.AddIfEmptyExpr(v)
		return nil, &err
	case *nodes.Inr:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_SUM_TYPE)
		err.AddIfEmptyExpr(v)
		return nil, &err
	case *nodes.List:
		if len(v.Exprs) == 0 {
			err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_LIST)
			err.AddIfEmptyExpr(v)
			return nil, &err
		}

		inferredType, err := infer(ctx, v.Exprs[0])
		index := 1

		for err != nil && err.IsAmbiguousError() && index < len(v.Exprs) {
			inferredType, err = infer(ctx, v.Exprs[index])
			index += 1
		}

		if err != nil {
			return nil, err
		}

		for _, expr := range v.Exprs {
			err := CheckType(ctx, expr, inferredType)

			if err != nil {
				return nil, err
			}
		}

		return &nodes.TypeList{Type_: inferredType, Repr: ""}, nil
	case *nodes.ConsList:
		headType, err := infer(ctx, v.Head)

		if err != nil && !err.IsAmbiguousError() {
			return nil, err
		}

		// Ambiguous head. Trying to get type from tail
		if err != nil {
			inferredType, err := infer(ctx, v.Tail)

			if err != nil {
				return nil, err
			}

			if _, ok := inferredType.(*nodes.TypeList); ok {
			} else {
				err := NewTypeCheckErrorErrorType(ERROR_NOT_A_LIST)
				err = err.AddIfEmptyExpr(v.Tail).AddIfEmptyActualType(inferredType)
				return nil, &err
			}

			listType, _ := inferredType.(*nodes.TypeList)

			err = CheckType(ctx, v.Head, listType.Type_)

			if err != nil {
				return nil, err
			}

			return listType, nil
		}

		listType := nodes.TypeList{Type_: headType, Repr: ""}
		err = CheckType(ctx, v.Tail, &listType)

		if err != nil {
			return nil, err
		}

		return &listType, nil
	case *nodes.Head:
		inferredType, err := infer(ctx, v.List)

		if err != nil {
			return nil, err
		}

		if listType, ok := inferredType.(*nodes.TypeList); ok {
			return listType.Type_, nil
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_LIST)
			err = err.AddIfEmptyExpr(v.List).AddIfEmptyActualType(inferredType)
			return nil, &err
		}
	case *nodes.Tail:
		inferredType, err := infer(ctx, v.List)

		if err != nil {
			return nil, err
		}

		if listType, ok := inferredType.(*nodes.TypeList); ok {
			return listType, nil
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_LIST)
			err = err.AddIfEmptyExpr(v.List).AddIfEmptyActualType(inferredType)
			return nil, &err
		}
	case *nodes.IsEmpty:
		inferredType, err := infer(ctx, v.List)

		if err != nil {
			return nil, err
		}

		if _, ok := inferredType.(*nodes.TypeList); ok {
			return &nodes.TypeBool{Repr: "bool"}, nil
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_LIST)
			err = err.AddIfEmptyExpr(v.List).AddIfEmptyActualType(inferredType)
			return nil, &err
		}
	case *nodes.Variant:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_VARIANT_TYPE)
		err.AddIfEmptyExpr(v)
		return nil, &err
	case *nodes.Fix:
		inferredType, err := infer(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		if listType, ok := inferredType.(*nodes.TypeFun); ok {

			if len(listType.ParamTypes) != 1 {
				// TODO show actual and expected number of parameters
				err_ := NewTypeCheckErrorErrorType(ERROR_INCORRECT_NUMBER_OF_ARGUMENTS).AddIfEmptyExpr(v.Expr_)

				return nil, &err_
			}

			if IsEqualType(listType.ParamTypes[0], listType.ReturnType) {
				var expectedFunType = nodes.TypeFun{ParamTypes: listType.ParamTypes, ReturnType: listType.ParamTypes[0]}

				err_ := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
				err_ = err_.AddIfEmptyActualType(inferredType).AddIfEmptyExpectedType(&expectedFunType)
				err_ = err_.AddIfEmptyExpr(v.Expr_)

				return nil, &err_
			}

			return listType.ReturnType, nil
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_FUNCTION)
			err = err.AddIfEmptyExpr(v.Expr_).AddIfEmptyActualType(inferredType)
			return nil, &err
		}
	case *nodes.Match:
		inferredType, err := infer(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		err = checkPatternTypes(v.Cases, inferredType)
		if err != nil {
			return nil, err
		}

		err = checkExhaustiveness(v.Cases, inferredType)
		if err != nil {
			return nil, err
		}

		// Try to get any inferred case type
		var inferredCaseType nodes.StellaType = nil
		for _, case_ := range v.Cases {
			ctx.AddNewScope()
			patternToContext(ctx, case_.Pattern, inferredType)

			if inferredCaseType == nil {
				inferredCaseType, err = infer(ctx, case_.Expr_)

				if err != nil && err.IsAmbiguousError() {
					continue
				} else if err != nil {
					ctx.RemoveLastScope()
					return nil, err
				}

				ctx.RemoveLastScope()
				break
			}
		}

		if inferredCaseType == nil {
			return nil, err
		}

		// Checking type
		for _, case_ := range v.Cases {
			ctx.AddNewScope()
			patternToContext(ctx, case_.Pattern, inferredType)

			err = CheckType(ctx, case_.Expr_, inferredCaseType)

			if err != nil {
				ctx.RemoveLastScope()
				return nil, err
			}

			ctx.RemoveLastScope()
		}

		return inferredCaseType, nil

	case *nodes.Let:
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()
		err := patternBindingsToContext(ctx, v.PatternBindings)
		if err != nil {
			if err.errorType == ERROR_DUPLICATE_LET_BINDING {
				err_ := err.AddIfEmptyExpr(node)
				return nil, &err_
			}
			return nil, err
		}
		return infer(ctx, v.Body)

	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		return nil, &err
	}
	err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
	return nil, &err
}

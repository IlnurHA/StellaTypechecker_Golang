package typecheck

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"

	"github.com/neocotic/go-optional"
)

func infer(ctx *Context, node nodes.Node) (nodes.StellaType, *TypecheckError) {
	switch v := node.(type) {
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
		err := CheckType(ctx, v.Condition, &nodes.TypeBool{})

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
		natType := nodes.TypeNat{}
		err := CheckType(ctx, v.N, &natType)

		if err != nil {
			err.AddIfEmptyExpr(node)
			return nil, err
		}

		return &natType, nil
	case *nodes.IsZero:
		natType := nodes.TypeNat{}
		err := CheckType(ctx, v.N, &natType)

		if err != nil {
			err.AddIfEmptyExpr(node)
			return nil, err
		}

		return &nodes.TypeBool{}, nil
	case *nodes.NatRec:
		natType := nodes.TypeNat{}
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
			err.AddIfEmptyExpr(v)
			return nil, err
		}

		paramTypes := getTypeFromParameters(v.Params)

		inferredType, err := infer(ctx, v.ReturnExpr)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeFun{ParamTypes: paramTypes, ReturnType: inferredType}, nil
	case *nodes.Application:
		inferredType, err := infer(ctx, v.Function)

		if err != nil {
			return nil, err
		}

		if _, ok := inferredType.(*nodes.TypeFun); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_FUNCTION)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a function", v.Function))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return nil, &err
		}

		inferredFunType, _ := inferredType.(*nodes.TypeFun)

		if len(v.Args) != len(inferredFunType.ParamTypes) {
			err := NewTypeCheckErrorErrorType(ERROR_INCORRECT_NUMBER_OF_ARGUMENTS)
			err.AddIfEmptyActualType(inferredFunType)
			err.AddAdditionalInfo(fmt.Sprintf("Expected %d. Got %d", len(inferredFunType.ParamTypes), len(v.Args)))
			err.AddIfEmptyExpr(v)
			err.Freeze()
			return nil, &err
		}

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
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a tuple", v.Subexpr))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return nil, &err
		}

		tupleType, _ := inferredType.(*nodes.TypeTuple)

		if v.Index > len(tupleType.Types) || v.Index <= 0 {
			err := NewTypeCheckErrorErrorType(ERROR_TUPLE_INDEX_OUT_OF_BOUNDS)
			err.AddIfEmptyExpr(v)
			err.AddAdditionalInfo(fmt.Sprintf("Actual length %d. Tried to access %d index", len(tupleType.Types), v.Index))
			err.AddIfEmptyActualType(tupleType)
			err.Freeze()
			return nil, &err
		}

		return tupleType.Types[v.Index-1], nil

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

			inferredType, err := infer(ctx, binding.Rhs)

			if err != nil {
				return nil, err
			}

			recordTypes[index] = nodes.RecordFieldType{Label: binding.Name, Type_: inferredType}
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

		err = CheckType(ctx, v.Expr_, v.Type_)

		if err != nil {
			return nil, err
		}

		return v.Type_, nil
	case *nodes.Inl:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_SUM_TYPE)
		err.AddIfEmptyExpr(v)
		err.Freeze()
		return nil, &err
	case *nodes.Inr:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_SUM_TYPE)
		err.AddIfEmptyExpr(v)
		err.Freeze()
		return nil, &err
	case *nodes.List:
		if len(v.Exprs) == 0 {
			err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_LIST)
			err.AddIfEmptyExpr(v)
			err.Freeze()
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

		return &nodes.TypeList{Type_: inferredType}, nil
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
				err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a list", v.Tail))
				err.AddIfEmptyExpr(v)
				err.AddIfEmptyActualType(inferredType)
				err.Freeze()
				return nil, &err
			}

			listType, _ := inferredType.(*nodes.TypeList)

			err = CheckType(ctx, v.Head, listType.Type_)

			if err != nil {
				return nil, err
			}

			return listType, nil
		}

		listType := nodes.TypeList{Type_: headType}
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
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a list", v.List))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
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
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a list", v.List))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return nil, &err
		}
	case *nodes.IsEmpty:
		inferredType, err := infer(ctx, v.List)

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
		if ctx.HasExtension(SUBTYPING) {
			subtype := optional.Empty[nodes.StellaType]()

			if v.Rhs.IsPresent() {
				inferredSubType, err := infer(ctx, v.Rhs.Require())
				if err != nil {
					return nil, nil
				}
				subtype = optional.Of(inferredSubType)
			}

			return &nodes.TypeVariant{FieldTypes: []nodes.VariantFieldType{{Label: v.Label, Type_: subtype}}}, nil
		}

		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_VARIANT_TYPE)
		err.AddIfEmptyExpr(v)
		err.Freeze()
		return nil, &err
	case *nodes.Fix:
		inferredType, err := infer(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		if listType, ok := inferredType.(*nodes.TypeFun); ok {

			if len(listType.ParamTypes) != 1 {
				err_ := NewTypeCheckErrorErrorType(ERROR_INCORRECT_NUMBER_OF_ARGUMENTS)
				err_.AddIfEmptyExpr(v.Expr_)
				err_.AddAdditionalInfo(fmt.Sprintf("Expected 1. Got %d", len(listType.ParamTypes)))
				err_.AddIfEmptyExpr(v)
				err_.Freeze()
				return nil, &err_
			}

			if !IsEqualType(listType.ParamTypes[0], listType.ReturnType) {
				var expectedFunType = nodes.TypeFun{ParamTypes: listType.ParamTypes, ReturnType: listType.ParamTypes[0]}

				err_ := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
				err_.AddIfEmptyActualType(inferredType)
				err_.AddIfEmptyExpectedType(&expectedFunType)
				err_.AddIfEmptyExpr(v.Expr_)

				return nil, &err_
			}

			return listType.ReturnType, nil
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
			err.AddIfEmptyExpr(v)
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
				err.AddIfEmptyExpr(node)
			}
			return nil, err
		}
		return infer(ctx, v.Body)

	case *nodes.Sequence:
		err := CheckType(ctx, v.Expr1, &nodes.TypeUnit{})

		if err != nil {
			return nil, err
		}

		return infer(ctx, v.Expr2)
	case *nodes.TerminatingSemicolon:
		return infer(ctx, v.Expr_)
	case *nodes.ParenthesisedExpr:
		return infer(ctx, v.Expr_)
	case *nodes.Ref:
		inferredType, err := infer(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeRef{Type_: inferredType}, nil
	case *nodes.Deref:
		inferredType, err := infer(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		if refType, ok := inferredType.(*nodes.TypeRef); ok {
			return refType.Type_, nil
		}

		typeErr := NewTypeCheckErrorErrorType(ERROR_NOT_A_REFERENCE)
		typeErr.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a reference", v.Expr_))
		typeErr.AddIfEmptyExpr(v)
		typeErr.AddIfEmptyActualType(inferredType)
		typeErr.Freeze()
		return nil, &typeErr
	case *nodes.Assign:
		lhs, err := infer(ctx, v.Lhs)

		if err != nil && err.IsAmbiguousError() {
			// If lhs is ambiguous try to get type from rhs
			rhs, err := infer(ctx, v.Rhs)

			if err != nil {
				return nil, err
			}

			err = CheckType(ctx, v.Lhs, &nodes.TypeRef{Type_: rhs})

			if err != nil {
				return nil, err
			}

			return &nodes.TypeUnit{}, nil
		}

		if err != nil {
			return nil, err
		}

		if refType, ok := lhs.(*nodes.TypeRef); ok {
			err = CheckType(ctx, v.Rhs, refType.Type_)

			if err != nil {
				return nil, err
			}

			return &nodes.TypeUnit{}, nil
		}

		typeErr := NewTypeCheckErrorErrorType(ERROR_NOT_A_REFERENCE)
		typeErr.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a reference", v.Lhs))
		typeErr.AddIfEmptyExpr(v)
		typeErr.AddIfEmptyActualType(lhs)
		typeErr.Freeze()
		return nil, &typeErr
	case *nodes.ConstMemory:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_REFERENCE_TYPE)
		err.AddIfEmptyExpr(v)
		err.Freeze()
		return nil, &err
	case *nodes.TypeCast:
		inferredType, err := infer(ctx, v.Expr_)

		if err != nil {
			return nil, err
		}

		err = CheckStellaType(ctx, inferredType, v.Type_)

		if err != nil {
			return nil, err
		}

		return v.Type_, nil
	case *nodes.Panic:
		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_PANIC_TYPE)
		err.AddIfEmptyExpr(v)
		err.Freeze()
		return nil, &err
	case *nodes.Throw:
		if ctx.GetExceptionType() == nil {
			err := NewTypeCheckErrorErrorType(ERROR_EXCEPTION_TYPE_NOT_DECLARED)
			err.AddIfEmptyExpr(v)
			return nil, &err
		}

		err := NewTypeCheckErrorErrorType(ERROR_AMBIGUOUS_THROW_TYPE)
		err.AddIfEmptyExpr(v)
		err.Freeze()
		return nil, &err
	case *nodes.TryWith:
		if ctx.GetExceptionType() == nil {
			err := NewTypeCheckErrorErrorType(ERROR_EXCEPTION_TYPE_NOT_DECLARED)
			err.AddIfEmptyExpr(v)
			return nil, &err
		}

		inferredType, err := infer(ctx, v.TryExpr)

		if err != nil && err.IsAmbiguousError() {
			// If lhs is ambiguous try to get type from rhs
			fallbackInferredType, err := infer(ctx, v.FallbackExpr)

			if err != nil {
				return nil, err
			}

			err = CheckType(ctx, v.TryExpr, fallbackInferredType)

			if err != nil {
				return nil, err
			}

			return fallbackInferredType, nil
		}

		if err != nil {
			return nil, err
		}

		err = CheckType(ctx, v.FallbackExpr, inferredType)

		if err != nil {
			return nil, err
		}

		return inferredType, nil
	case *nodes.TryCatch:
		exceptionType := ctx.GetExceptionType()

		if exceptionType == nil {
			err := NewTypeCheckErrorErrorType(ERROR_EXCEPTION_TYPE_NOT_DECLARED)
			err.AddIfEmptyExpr(v)
			return nil, &err
		}

		inferredType, err := infer(ctx, v.TryExpr)

		if err != nil && err.IsAmbiguousError() {
			// If lhs is ambiguous try to get type from rhs
			// Pattern handling
			ctx.AddNewScope()

			err = checkPatternType(v.Pattern, exceptionType)

			if err != nil {
				ctx.RemoveLastScope()
				return nil, err
			}

			patternToContext(ctx, v.Pattern, exceptionType)

			fallbackInferredType, err := infer(ctx, v.FallbackExpr)
			ctx.RemoveLastScope()

			if err != nil {
				return nil, err
			}

			err = CheckType(ctx, v.TryExpr, fallbackInferredType)

			if err != nil {
				return nil, err
			}

			return fallbackInferredType, nil
		}

		if err != nil {
			return nil, err
		}

		// Pattern handling
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		err = checkPatternType(v.Pattern, exceptionType)

		if err != nil {
			return nil, err
		}

		patternToContext(ctx, v.Pattern, exceptionType)

		err = CheckType(ctx, v.FallbackExpr, inferredType)

		if err != nil {
			return nil, err
		}

		return inferredType, nil

	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		err.AddAdditionalInfo(fmt.Sprintf("Not implemented type inference for %s", node))
		return nil, &err
	}
}

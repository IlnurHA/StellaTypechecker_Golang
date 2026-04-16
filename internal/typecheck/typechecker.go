package typecheck

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
)

func ParseProgram(node nodes.AProgram) *TypecheckError {
	context := NewContext(node.Extensions)

	context.AddNewScope()
	defer context.RemoveLastScope()

	hasMain := false

	// Collect all declarations to context
	for _, declaration := range node.Declarations {
		switch decl := declaration.(type) {
		case *nodes.FunctionDeclaration:
			declType, err := constructTypeFromDeclaration(&declaration)

			if err != nil {
				return err
			}

			name := decl.Name.String()
			nameStella := nodes.StellaIdent{Name: name, Repr: name}

			success := context.AddVar(nameStella, declType)

			if !success {
				var typeError = NewTypeCheckErrorErrorType(ERROR_DUPLICATE_FUNCTION_DECLARATION)
				typeError.AddAdditionalInfo(fmt.Sprintf("Duplicated function name: %s", &nameStella))
				// typeError.AddIfEmptyExpr(&node)
				return &typeError
			}

			if name == "main" {
				hasMain = true

				if len(declType.(*nodes.TypeFun).ParamTypes) != 1 {
					err := NewTypeCheckErrorErrorType(ERROR_INCORRECT_ARITY_OF_MAIN)
					return &err
				}
			}
		default:
			err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
			err.AddIfEmptyExpr(decl)
			err.AddAdditionalInfo(fmt.Sprintf("Not implemented parse program declaration switch for %s", decl))
			return &err
		}
	}

	if !hasMain {
		var typeError = NewTypeCheckErrorErrorType(ERROR_MISSING_MAIN)
		return &typeError
	}

	context.AddNewScope()
	// Check declaration bodies
	for _, declaration := range node.Declarations {
		type_, _ := constructTypeFromDeclaration(&declaration)
		err := CheckType(context, declaration, type_)

		if err != nil {
			return err
		}
	}
	context.RemoveLastScope()

	return nil
}

func CheckType(ctx *Context, node nodes.Node, expectedType nodes.StellaType) (err *TypecheckError) {

	defer func() {
		if err != nil {
			err.AddIfEmptyExpr(node)
			err.AddIfEmptyExpectedType(expectedType)
		}
	}()

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

		// Add current function for recursion
		ctx.AddVar(v.Name, expectedType)
		defer ctx.RemoveLastScope()

		// function parameters scope
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		err := addParametersToContext(ctx, v.Params)
		if err != nil {
			err.AddIfEmptyExpr(v)
			return err
		}

		// subdeclarations scope
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		// Collect all names
		for _, declaration := range v.Declarations {

			switch decl := declaration.(type) {
			case *nodes.FunctionDeclaration:
				type_, err := constructTypeFromDeclaration(&declaration)

				if err != nil {
					return err
				}

				success := ctx.AddVar(decl.Name, type_)

				if !success {
					err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_FUNCTION_DECLARATION)
					err.AddIfEmptyFunctionName(v.Name)
					// err.AddIfEmptyExpr(decl)
					return &err
				}
			default:
				err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
				err.AddIfEmptyExpr(decl)
				err.AddAdditionalInfo(fmt.Sprintf("Not implemented check type function declaration switch for %s", decl))
				return &err
			}
		}

		// Check types
		for _, decl := range v.Declarations {
			type_, err := constructTypeFromDeclaration(&decl)

			if err != nil {
				return err
			}

			err = CheckType(ctx, decl, type_)

			if err != nil {
				return err
			}
		}

		// body scope
		ctx.AddNewScope()

		err = CheckType(ctx, v.Expr, v.ReturnType.OrElse(&nodes.TypeUnit{}))

		if err != nil {
			return err
		}

		return nil
	case *nodes.ConstUnit:
		return CheckStellaType(&nodes.TypeUnit{}, expectedType)
	case *nodes.ConstBool:
		return CheckStellaType(&nodes.TypeBool{}, expectedType)
	case *nodes.ConstInt:
		return CheckStellaType(&nodes.TypeNat{}, expectedType)
	case *nodes.Var:
		// println("Expr:", v.Name.Name)
		type_ := ctx.GetVarType(v.Name)

		if type_.IsEmpty() {
			var err = NewTypeCheckErrorErrorType(ERROR_UNDEFINED_VARIABLE)
			err.AddAdditionalInfo(fmt.Sprintf("Undefined variable: %s", v.String()))
			return &err
		}
		return CheckStellaType(type_.Require(), expectedType)
	case *nodes.If:
		err := CheckType(ctx, v.Condition, &nodes.TypeBool{})

		if err != nil {
			return err
		}

		err = CheckType(ctx, v.ThenExpr, expectedType)

		if err != nil {
			return err
		}

		return CheckType(ctx, v.ElseExpr, expectedType)
	case *nodes.Succ:
		natType := nodes.TypeNat{}

		err := CheckStellaType(&natType, expectedType)

		if err != nil {
			return err
		}

		return CheckType(ctx, v.N, &natType)
	case *nodes.IsZero:
		natType := nodes.TypeNat{}

		err := CheckStellaType(&nodes.TypeBool{}, expectedType)

		if err != nil {
			return err
		}

		return CheckType(ctx, v.N, &natType)
	case *nodes.NatRec:
		natType := nodes.TypeNat{}
		err := CheckType(ctx, v.N, &natType)

		if err != nil {
			return err
		}

		err = CheckType(ctx, v.Initial, expectedType)

		if err != nil {
			return err
		}

		stepSubType := nodes.TypeFun{ParamTypes: []nodes.StellaType{expectedType}, ReturnType: expectedType}
		stepType := nodes.TypeFun{ParamTypes: []nodes.StellaType{&natType}, ReturnType: &stepSubType}

		return CheckType(ctx, v.Step, &stepType)
	case *nodes.Abstraction:
		// Unexpected lambda
		if _, ok := expectedType.(*nodes.TypeFun); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_LAMBDA)
			err.AddIfEmptyExpectedType(expectedType)
			return &err
		}

		//  ERROR_UNEXPECTED_NUMBER_OF_PARAMETERS_IN_LAMBDA
		typeFun, _ := expectedType.(*nodes.TypeFun)

		if len(typeFun.ParamTypes) != len(v.Params) {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_NUMBER_OF_PARAMETERS_IN_LAMBDA)
			err.AddIfEmptyExpectedType(expectedType)
			err.AddAdditionalInfo(fmt.Sprintf("Expected %d. Got %d", len(typeFun.ParamTypes), len(v.Params)))
			return &err
		}

		// ERROR_UNEXPECTED_TYPE_FOR_PARAMETER
		for _, paramPair := range Zip(v.Params, typeFun.ParamTypes) {
			err := CheckStellaType(paramPair.First.ParameterType, paramPair.Second)

			if err != nil {
				err.RewriteErrorType(ERROR_UNEXPECTED_TYPE_FOR_PARAMETER)
				err.AddAdditionalInfo(fmt.Sprintf("Unexpected type for parameter: %s", &paramPair.First.Name))
				return err
			}
		}

		ctx.AddNewScope()
		defer ctx.RemoveLastScope()

		err := addParametersToContext(ctx, v.Params)

		if err != nil {
			err.AddIfEmptyExpr(v)
			return err
		}

		return CheckType(ctx, v.ReturnExpr, typeFun.ReturnType)
	case *nodes.Application:
		inferredType, err := infer(ctx, v.Function)

		if err != nil {
			return err
		}

		if _, ok := inferredType.(*nodes.TypeFun); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_FUNCTION)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a function", v.Function))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return &err
		}

		inferredFunType, _ := inferredType.(*nodes.TypeFun)

		if len(v.Args) != len(inferredFunType.ParamTypes) {
			err := NewTypeCheckErrorErrorType(ERROR_INCORRECT_NUMBER_OF_ARGUMENTS)
			err.AddIfEmptyActualType(inferredFunType)
			err.AddAdditionalInfo(fmt.Sprintf("Expected %d. Got %d", len(inferredFunType.ParamTypes), len(v.Args)))
			err.AddIfEmptyExpr(v)
			err.Freeze()
			return &err
		}

		for _, pair := range Zip(v.Args, inferredFunType.ParamTypes) {
			expr := pair.First
			type_ := pair.Second

			err = CheckType(ctx, expr, type_)

			if err != nil {
				return err
			}
		}

		return CheckStellaType(inferredFunType.ReturnType, expectedType)
	case *nodes.Tuple:
		// ERROR_UNEXPECTED_TUPLE
		if _, ok := expectedType.(*nodes.TypeTuple); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TUPLE)
			err.AddIfEmptyExpectedType(expectedType)
			return &err
		}

		// ERROR_UNEXPECTED_TUPLE_LENGTH
		tupleType, _ := expectedType.(*nodes.TypeTuple)
		if len(tupleType.Types) != len(v.Exprs) {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TUPLE_LENGTH)
			err.AddIfEmptyExpectedType(expectedType)
			err.AddAdditionalInfo(fmt.Sprintf("Expected tuple of length %d. Got %d", len(tupleType.Types), len(v.Exprs)))
			return &err
		}

		for _, exprTypePair := range Zip(v.Exprs, tupleType.Types) {
			err := CheckType(ctx, exprTypePair.First, exprTypePair.Second)

			if err != nil {
				return err
			}
		}

		return nil
	case *nodes.DotTuple:
		inferredType, err := infer(ctx, v.Subexpr)

		if err != nil {
			return err
		}

		if _, ok := inferredType.(*nodes.TypeTuple); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_TUPLE)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a tuple", v.Subexpr))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return &err
		}

		tupleType, _ := inferredType.(*nodes.TypeTuple)

		if v.Index > len(tupleType.Types) || v.Index <= 0 {
			err := NewTypeCheckErrorErrorType(ERROR_TUPLE_INDEX_OUT_OF_BOUNDS)
			err.AddIfEmptyExpr(v)
			err.AddAdditionalInfo(fmt.Sprintf("Actual length %d. Tried to access %d index", len(tupleType.Types), v.Index))
			err.AddIfEmptyActualType(tupleType)
			err.Freeze()
			return &err
		}

		return CheckStellaType(tupleType.Types[v.Index-1], expectedType)

	case *nodes.Record:
		// ERROR_UNEXPECTED_RECORD
		if _, ok := expectedType.(*nodes.TypeRecord); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_RECORD)
			err.AddIfEmptyExpectedType(expectedType)
			return &err
		}

		// ERROR_MISSING_RECORD_FIELDS
		recordType, _ := expectedType.(*nodes.TypeRecord)

		notFoundFields := make([]*nodes.StellaIdent, 0)
		for _, expectedField := range recordType.FieldTypes {
			// fmt.Printf("DEBUG: Expected Field: %s\n", &expectedField.Label)
			isFound := false

			for _, actualField := range v.Bindings {
				// fmt.Printf("DEBUG: Actual Field: %s\n", &actualField.Name)
				if expectedField.Label.Equal(&actualField.Name) {
					isFound = true
				}
			}

			if !isFound {
				notFoundFields = append(notFoundFields, &expectedField.Label)
			}
		}

		if len(notFoundFields) != 0 {
			err := NewTypeCheckErrorErrorType(ERROR_MISSING_RECORD_FIELDS)
			err.AddAdditionalInfo(fmt.Sprintf("Missing the following fields: %s", showList(notFoundFields)))
			return &err
		}

		// ERROR_DUPLICATE_RECORD_FIELDS
		recordFields := make(map[nodes.StellaIdent]bool, len(v.Bindings))

		for _, binding := range v.Bindings {
			if _, ok := recordFields[binding.Name]; ok {
				err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_RECORD_FIELDS)
				err.AddAdditionalInfo(fmt.Sprintf("Duplicate label: %s", &binding.Name))
				return &err
			}
			recordFields[binding.Name] = true
		}

		// ERROR_UNEXPECTED_RECORD_FIELDS
		if len(recordType.FieldTypes) < len(v.Bindings) {
			unexpectedFields := make([]*nodes.StellaIdent, 0)
			for _, label := range v.Bindings {
				isFound := false

				for _, expectedLabel := range recordType.FieldTypes {
					if expectedLabel.Label.Equal(&label.Name) {
						isFound = true
						break
					}
				}
				if !isFound {
					unexpectedFields = append(unexpectedFields, &label.Name)
				}
			}

			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_RECORD_FIELDS)
			err.AddAdditionalInfo(fmt.Sprintf("Unexpected fields: %s", showList(unexpectedFields)))
			return &err
		}

		for _, expectedFieldType := range recordType.FieldTypes {
			for _, binding := range v.Bindings {
				if binding.Name.Equal(&expectedFieldType.Label) {
					err := CheckType(ctx, binding.Rhs, expectedFieldType.Type_)

					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	case *nodes.DotRecord:
		inferredType, err := infer(ctx, v.Subexpr)

		if err != nil {
			return err
		}

		if _, ok := inferredType.(*nodes.TypeRecord); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_RECORD)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a record", v.Subexpr))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return &err
		}

		recordType, _ := inferredType.(*nodes.TypeRecord)

		for _, typeBinding := range recordType.FieldTypes {
			if typeBinding.Label.Name == v.Label.Name {
				return CheckStellaType(typeBinding.Type_, expectedType)
			}
		}

		err_ := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_FIELD_ACCESS)
		err_.AddIfEmptyExpr(v)
		err_.AddAdditionalInfo(fmt.Sprintf("Unexpected label \"%s\" for record of type %s", v.Label.String(), recordType))
		err_.Freeze()
		return &err_
	case *nodes.TypeAsc:
		if err := checkTypeConsistency(v.Type_); err != nil {
			return err
		}
		if err := CheckStellaType(v.Type_, expectedType); err != nil {
			return err
		}
		return CheckType(ctx, v.Expr_, expectedType)
	case *nodes.Inl:
		// ERROR_UNEXPECTED_INJECTION
		if _, ok := expectedType.(*nodes.TypeSum); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_INJECTION)
			return &err
		}

		sumType, _ := expectedType.(*nodes.TypeSum)
		return CheckType(ctx, v.Expr_, sumType.Left)
	case *nodes.Inr:
		// ERROR_UNEXPECTED_INJECTION
		if _, ok := expectedType.(*nodes.TypeSum); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_INJECTION)
			return &err
		}

		sumType, _ := expectedType.(*nodes.TypeSum)
		return CheckType(ctx, v.Expr_, sumType.Right)
	case *nodes.List:
		// ERROR_UNEXPECTED_LIST
		if _, ok := expectedType.(*nodes.TypeList); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_LIST)
			return &err
		}

		listType, _ := expectedType.(*nodes.TypeList)

		for _, expr := range v.Exprs {
			err := CheckType(ctx, expr, listType.Type_)

			if err != nil {
				return err
			}
		}

		return nil
	case *nodes.ConsList:
		// ERROR_UNEXPECTED_LIST
		if _, ok := expectedType.(*nodes.TypeList); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_LIST)
			return &err
		}

		listType, _ := expectedType.(*nodes.TypeList)

		err := CheckType(ctx, v.Head, listType.Type_)

		if err != nil {
			return err
		}

		return CheckType(ctx, v.Tail, listType.Type_)
	case *nodes.Head:
		listType := nodes.TypeList{Type_: expectedType}
		return CheckType(ctx, v.List, &listType)
	case *nodes.Tail:
		// ERROR_UNEXPECTED_LIST
		if _, ok := expectedType.(*nodes.TypeList); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_LIST)
			return &err
		}
		return CheckType(ctx, v.List, expectedType)
	case *nodes.IsEmpty:
		inferredType, err := infer(ctx, v.List)

		if err != nil {
			return err
		}

		if _, ok := inferredType.(*nodes.TypeList); ok {
			return CheckStellaType(&nodes.TypeBool{}, expectedType)
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_NOT_A_LIST)
			err.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a list", v.List))
			err.AddIfEmptyExpr(v)
			err.AddIfEmptyActualType(inferredType)
			err.Freeze()
			return &err
		}
	case *nodes.Variant:
		// ERROR_UNEXPECTED_VARIANT
		if _, ok := expectedType.(*nodes.TypeVariant); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_VARIANT)
			err.AddAdditionalInfo(fmt.Sprintf("Expected %s. Got variant type", expectedType))
			return &err
		}

		variantType, _ := expectedType.(*nodes.TypeVariant)

		for _, expectedVariantField := range variantType.FieldTypes {
			if !expectedVariantField.Label.Equal(&v.Label) {
				continue
			}

			if expectedVariantField.Type_.IsEmpty() && v.Rhs.IsEmpty() {
				return nil
			}

			if expectedVariantField.Type_.IsPresent() && v.Rhs.IsPresent() {
				return CheckType(ctx, v.Rhs.Require(), expectedVariantField.Type_.Require())
			}

			if expectedVariantField.Type_.IsEmpty() {
				err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_DATA_FOR_NULLARY_LABEL)
				return &err
			}

			err := NewTypeCheckErrorErrorType(ERROR_MISSING_DATA_FOR_LABEL)
			return &err
		}
		err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_VARIANT_LABEL)
		err.AddAdditionalInfo(fmt.Sprintf("Unexpected label: %s", &v.Label))
		return &err
	case *nodes.Fix:
		listType := nodes.TypeFun{ParamTypes: []nodes.StellaType{expectedType}, ReturnType: expectedType}
		return CheckType(ctx, v.Expr_, &listType)
	case *nodes.Match:
		if len(v.Cases) == 0 {
			err := NewTypeCheckErrorErrorType(ERROR_ILLEGAL_EMPTY_MATCHING)
			err.AddIfEmptyExpr(node)
			return &err
		}

		inferredType, err := infer(ctx, v.Expr_)

		if err != nil {
			return err
		}

		err = checkPatternTypes(v.Cases, inferredType)
		if err != nil {
			return err
		}

		err = checkExhaustiveness(v.Cases, inferredType)
		if err != nil {
			err.AddIfEmptyExpr(v)
			err.Freeze()
			return err
		}

		// Checking type
		for _, case_ := range v.Cases {
			ctx.AddNewScope()
			patternToContext(ctx, case_.Pattern, inferredType)

			err = CheckType(ctx, case_.Expr_, expectedType)

			if err != nil {
				ctx.RemoveLastScope()
				return err
			}

			ctx.RemoveLastScope()
		}

		return nil

	case *nodes.Let:
		ctx.AddNewScope()
		defer ctx.RemoveLastScope()
		err := patternBindingsToContext(ctx, v.PatternBindings)
		if err != nil {
			if err.errorType == ERROR_DUPLICATE_LET_BINDING {
				err.AddIfEmptyExpr(node)
			}
			return err
		}
		return CheckType(ctx, v.Body, expectedType)

	case *nodes.Sequence:
		err := CheckType(ctx, v.Expr1, &nodes.TypeUnit{})

		if err != nil {
			return err
		}

		return CheckType(ctx, v.Expr2, expectedType)
	case *nodes.TerminatingSemicolon:
		return CheckType(ctx, v.Expr_, expectedType)
	case *nodes.ParenthesisedExpr:
		return CheckType(ctx, v.Expr_, expectedType)
	case *nodes.Ref:
		// ERROR_UNEXPECTED_REFERENCE
		if _, ok := expectedType.(*nodes.TypeRef); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_REFERENCE)
			err.AddAdditionalInfo(fmt.Sprintf("Expected %s. Got reference type", expectedType))
			return &err
		}

		refType, _ := expectedType.(*nodes.TypeRef)

		return CheckType(ctx, v.Expr_, refType.Type_)
	case *nodes.Deref:
		return CheckType(ctx, v.Expr_, &nodes.TypeRef{Type_: expectedType})
	case *nodes.Assign:
		// Check if expected type is unit
		err := CheckStellaType(&nodes.TypeUnit{}, expectedType)

		if err != nil {
			return err
		}

		// Get assign
		lhs, err := infer(ctx, v.Lhs)

		if err != nil && err.IsAmbiguousError() {
			// If lhs is ambiguous try to get type from rhs
			rhs, err := infer(ctx, v.Rhs)

			if err != nil {
				return err
			}

			return CheckType(ctx, v.Lhs, &nodes.TypeRef{Type_: rhs})
		}

		if err != nil {
			return err
		}

		if refType, ok := lhs.(*nodes.TypeRef); ok {
			return CheckType(ctx, v.Rhs, refType.Type_)
		}

		typeErr := NewTypeCheckErrorErrorType(ERROR_NOT_A_REFERENCE)
		typeErr.AddAdditionalInfo(fmt.Sprintf("Expression %s is not a reference", v.Lhs))
		typeErr.AddIfEmptyExpr(v)
		typeErr.AddIfEmptyActualType(lhs)
		typeErr.Freeze()
		return &typeErr
	case *nodes.ConstMemory:
		// ERROR_UNEXPECTED_MEMORY_ADDRESS
		if _, ok := expectedType.(*nodes.TypeRef); ok {
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_MEMORY_ADDRESS)
			err.AddAdditionalInfo(fmt.Sprintf("Expected %s. Got reference type", expectedType))
			return &err
		}

		return nil
	case *nodes.TypeCast:
		err := CheckStellaType(v.Type_, expectedType)

		if err != nil {
			return err
		}

		return CheckType(ctx, v.Expr_, v.Type_)

	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		err.AddAdditionalInfo(fmt.Sprintf("Not implemented check type switch for %T", node))
		return &err
	}
}

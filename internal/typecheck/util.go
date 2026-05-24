package typecheck

import (
	"fmt"
	"strings"
	nodes "typechecker/internal/ast/nodes"
)

func addParametersToContext(ctx *Context, paramDecls []nodes.ParameterDeclaration) ([]nodes.ParameterDeclaration, *TypecheckError) {

	newParamDeclarations := make([]nodes.ParameterDeclaration, 0, len(paramDecls))
	for _, param := range paramDecls {

		type_ := param.ParameterType

		if ctx.HasExtension(TYPE_RECONSTRUCTION) {
			var err *TypecheckError
			type_, err = transformTypeAuto(ctx, type_)

			if err != nil {
				return nil, err
			}
		}

		err := checkTypeConsistency(type_)

		if err != nil {
			return nil, err
		}

		res := ctx.AddVar(param.Name, type_)
		newParamDeclarations = append(newParamDeclarations, nodes.ParameterDeclaration{Name: param.Name, ParameterType: type_})
		// println("Tried to add", param.Name.String(), "with type", param.ParameterType.String(), "; success:", res)

		if !res {
			err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_FUNCTION_PARAMETER)
			return nil, &err
		}
	}
	return newParamDeclarations, nil
}

func constructTypeFromDeclaration(ctx *Context, declaration *nodes.Declaration) (nodes.StellaType, *TypecheckError) {
	declarationCheck := *declaration
	switch v := declarationCheck.(type) {
	case *nodes.FunctionDeclaration:
		paramTypes := make([]nodes.StellaType, len(v.Params))
		returnType := v.ReturnType.OrElse(&nodes.TypeUnit{})

		if ctx.HasExtension(TYPE_RECONSTRUCTION) {
			var err *TypecheckError
			returnType, err = transformTypeAuto(ctx, returnType)

			if err != nil {
				return nil, err
			}
		}

		for index, param := range v.Params {
			type_ := param.ParameterType

			if ctx.HasExtension(TYPE_RECONSTRUCTION) {
				var err *TypecheckError
				type_, err = transformTypeAuto(ctx, type_)

				if err != nil {
					return nil, err
				}
			}

			err := checkTypeConsistency(type_)

			if err != nil {
				return nil, err
			}

			paramTypes[index] = type_
		}

		err := checkTypeConsistency(returnType)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeFun{ParamTypes: paramTypes, ReturnType: returnType}, nil
	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		err.AddIfEmptyExpr(*declaration)
		err.AddAdditionalInfo(fmt.Sprintf("Not implemented constructTypeFromDeclaration switch for %s", *declaration))

		return nil, &err
	}
}

func getTypeFromParameters(paramDecls []nodes.ParameterDeclaration) []nodes.StellaType {
	paramTypes := make([]nodes.StellaType, len(paramDecls))

	for index, decl := range paramDecls {
		paramTypes[index] = decl.ParameterType
	}

	return paramTypes
}

func showList[T fmt.Stringer](elems []T) string {
	var builder strings.Builder

	builder.WriteString("[")
	for index, elem := range elems {
		if index != 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(elem.String())
	}
	builder.WriteString("]")

	return builder.String()
}

func checkTypeConsistency(type_ nodes.StellaType) *TypecheckError {
	switch t := type_.(type) {
	case *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit, *nodes.TypeBot, *nodes.TypeTop:
		return nil
	case *nodes.TypeFun:
		for _, param := range t.ParamTypes {
			err := checkTypeConsistency(param)

			if err != nil {
				return err
			}
		}

		return checkTypeConsistency(t.ReturnType)
	case *nodes.TypeList:
		return checkTypeConsistency(t.Type_)
	case *nodes.TypeParens:
		return checkTypeConsistency(t.Type_)
	case *nodes.TypeRef:
		return checkTypeConsistency(t.Type_)
	case *nodes.TypeSum:
		err := checkTypeConsistency(t.Left)

		if err != nil {
			return err
		}

		return checkTypeConsistency(t.Right)
	case *nodes.TypeTuple:
		for _, type_ := range t.Types {
			err := checkTypeConsistency(type_)

			if err != nil {
				return err
			}
		}

		return nil
	case *nodes.TypeRecord:
		// Check for record field duplicate
		labels := make(map[nodes.StellaIdent]bool)
		for _, recordFieldType := range t.FieldTypes {
			if _, ok := labels[recordFieldType.Label]; ok {
				err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_RECORD_TYPE_FIELDS)
				err.AddAdditionalInfo(fmt.Sprintf("Label duplicate: %s", &recordFieldType.Label))
				return &err
			}
			labels[recordFieldType.Label] = true
		}
		return nil
	case *nodes.TypeVariant:
		// Check for variant field duplicate
		labels := make(map[nodes.StellaIdent]bool)
		for _, recordFieldType := range t.FieldTypes {
			if _, ok := labels[recordFieldType.Label]; ok {
				err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_VARIANT_TYPE_FIELDS)
				err.AddAdditionalInfo(fmt.Sprintf("Label duplicate: %s", &recordFieldType.Label))
				return &err
			}
			labels[recordFieldType.Label] = true
		}
		return nil
	case *nodes.TypeVar:
		return nil
	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		err.AddAdditionalInfo(fmt.Sprintf("Unimplemented check type consistency for %s", type_.String()))
		return &err
	}
}

// Pair defines a generic struct to hold two values of potentially different types.
type Pair[T, U any] struct {
	First  T
	Second U
}

// Zip combines two slices into a slice of Pairs.
// The length of the result is determined by the length of the shorter input slice.
func Zip[T, U any](ts []T, us []U) []Pair[T, U] {
	minLen := len(ts)
	if len(us) < minLen {
		minLen = len(us)
	}

	pairs := make([]Pair[T, U], minLen)
	for i := 0; i < minLen; i++ {
		pairs[i] = Pair[T, U]{First: ts[i], Second: us[i]}
	}
	return pairs
}

func transformTypeAuto(ctx *Context, type_ nodes.StellaType) (nodes.StellaType, *TypecheckError) {
	switch t := type_.(type) {
	case *nodes.TypeAuto:
		freshVar := ctx.constraintHandler.GetFreshVar()
		return &freshVar, nil
	case *nodes.TypeVar, *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit:
		return t, nil
	case *nodes.TypeForAll:
		newType, err := transformTypeAuto(ctx, t.Type_)
		if err != nil {
			return nil, err
		}
		return &nodes.TypeForAll{Type_: newType, Types: t.Types}, nil
	case *nodes.TypeFun:
		newReturnType, err := transformTypeAuto(ctx, t.ReturnType)
		if err != nil {
			return nil, err
		}

		newParamTypes := make([]nodes.StellaType, 0, len(t.ParamTypes))
		for _, t_ := range t.ParamTypes {
			newParamType, err := transformTypeAuto(ctx, t_)
			if err != nil {
				return nil, err
			}
			newParamTypes = append(newParamTypes, newParamType)
		}

		return &nodes.TypeFun{ParamTypes: newParamTypes, ReturnType: newReturnType}, nil
	case *nodes.TypeList:
		newType, err := transformTypeAuto(ctx, t.Type_)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeList{Type_: newType}, nil
	case *nodes.TypeParens:
		newType, err := transformTypeAuto(ctx, t.Type_)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeParens{Type_: newType}, nil
	case *nodes.TypeRecord:
		newFieldTypes := make([]nodes.RecordFieldType, 0, len(t.FieldTypes))
		for _, t_ := range t.FieldTypes {
			newType, err := transformTypeAuto(ctx, t_.Type_)
			if err != nil {
				return nil, err
			}
			newFieldTypes = append(newFieldTypes, nodes.RecordFieldType{Label: t_.Label, Type_: newType})
		}
		return &nodes.TypeRecord{FieldTypes: newFieldTypes}, nil
	case *nodes.TypeRef:
		newType, err := transformTypeAuto(ctx, t.Type_)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeRef{Type_: newType}, nil
	case *nodes.TypeSum:

		newLeft, err := transformTypeAuto(ctx, t.Left)

		if err != nil {
			return nil, err
		}

		newRight, err := transformTypeAuto(ctx, t.Right)

		if err != nil {
			return nil, err
		}

		return &nodes.TypeSum{
			Left:  newLeft,
			Right: newRight,
		}, nil
	case *nodes.TypeTuple:
		newTypes := make([]nodes.StellaType, 0, len(t.Types))
		for _, t_ := range t.Types {
			newType, err := transformTypeAuto(ctx, t_)

			if err != nil {
				return nil, err
			}

			newTypes = append(newTypes, newType)
		}

		return &nodes.TypeTuple{Types: newTypes}, nil
	}

	err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
	err.AddAdditionalInfo(fmt.Sprintf("Unimplemented check type consistency for %s", type_.String()))
	return nil, &err
}

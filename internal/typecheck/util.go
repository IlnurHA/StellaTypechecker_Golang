package typecheck

import (
	"errors"
	nodes "typechecker/internal/ast/nodes"
)

func addParametersToContext(ctx *Context, paramDecls []nodes.ParameterDeclaration) *TypecheckError {
	for _, param := range paramDecls {
		res := ctx.AddVar(param.Name, param.ParameterType)

		if !res {
			err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_FUNCTION_PARAMETER).AddIfEmptyExpr(&param)
			return &err
		}
	}
	return nil
}

func constructTypeFromDeclaration(declaration *nodes.Declaration) (nodes.StellaType, error) {
	declarationCheck := *declaration
	switch v := declarationCheck.(type) {
	case *nodes.FunctionDeclaration:
		paramTypes := make([]nodes.StellaType, len(v.Params))
		returnType := v.ReturnType.OrElse(&nodes.TypeUnit{Repr: "unit"})

		for index, param := range v.Params {
			paramTypes[index] = param.ParameterType
		}

		return &nodes.TypeFun{ParamTypes: paramTypes, ReturnType: returnType, Repr: ""}, nil
	default:
		return nil, errors.New("Not Implemented")
	}
}

func getTypeFromParameters(paramDecls []nodes.ParameterDeclaration) []nodes.StellaType {
	paramTypes := make([]nodes.StellaType, len(paramDecls))

	for index, decl := range paramDecls {
		paramTypes[index] = decl.ParameterType
	}

	return paramTypes
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

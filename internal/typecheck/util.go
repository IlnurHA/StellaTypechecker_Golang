package typecheck

import (
	"fmt"
	"strings"
	nodes "typechecker/internal/ast/nodes"
)

func addParametersToContext(ctx *Context, paramDecls []nodes.ParameterDeclaration) *TypecheckError {
	for _, param := range paramDecls {
		res := ctx.AddVar(param.Name, param.ParameterType)

		if !res {
			err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_FUNCTION_PARAMETER)
			err.AddIfEmptyExpr(&param)
			return &err
		}
	}
	return nil
}

func constructTypeFromDeclaration(declaration *nodes.Declaration) (nodes.StellaType, *TypecheckError) {
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
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		err.AddIfEmptyExpr(*declaration)
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

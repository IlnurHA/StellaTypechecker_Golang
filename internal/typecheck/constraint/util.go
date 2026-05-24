package constraint

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
)

func needReconstruction(type_ nodes.StellaType) bool {
	switch t := type_.(type) {
	case *nodes.TypeVar:
		return t.Generated
	case *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit:
		return false
	case *nodes.TypeForAll:
		return needReconstruction(t.Type_)
	case *nodes.TypeFun:
		result := needReconstruction(t.ReturnType)

		for _, t_ := range t.ParamTypes {
			result = result || needReconstruction(t_)
		}

		return result
	case *nodes.TypeList:
		return needReconstruction(t.Type_)
	case *nodes.TypeParens:
		return needReconstruction(t.Type_)
	case *nodes.TypeRecord:
		result := false
		for _, t_ := range t.FieldTypes {
			result = result || needReconstruction(t_.Type_)
		}
		return result
	case *nodes.TypeRef:
		return needReconstruction(t.Type_)
	case *nodes.TypeSum:
		return needReconstruction(t.Left) || needReconstruction(t.Right)
	case *nodes.TypeTuple:
		result := false

		for _, t_ := range t.Types {
			result = result || needReconstruction(t_)
		}

		return result
	}

	return false
}

func freeVariables(type_ nodes.StellaType) []nodes.StellaIdent {
	switch t := type_.(type) {
	case *nodes.TypeVar:
		if t.Generated {
			return []nodes.StellaIdent{t.Name}
		}
	case *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit:
		return []nodes.StellaIdent{}
	case *nodes.TypeForAll:
		return freeVariables(t.Type_)
	case *nodes.TypeFun:
		result := freeVariables(t.ReturnType)

		for _, t_ := range t.ParamTypes {
			result = append(result, freeVariables(t_)...)
		}

		return result
	case *nodes.TypeList:
		return freeVariables(t.Type_)
	case *nodes.TypeParens:
		return freeVariables(t.Type_)
	case *nodes.TypeRecord:
		result := make([]nodes.StellaIdent, 0)
		for _, t_ := range t.FieldTypes {
			result = append(result, freeVariables(t_.Type_)...)
		}
		return result
	case *nodes.TypeRef:
		return freeVariables(t.Type_)
	case *nodes.TypeSum:
		result := freeVariables(t.Left)
		result = append(result, freeVariables(t.Right)...)
		return result
	case *nodes.TypeTuple:
		result := make([]nodes.StellaIdent, 0)

		for _, t_ := range t.Types {
			result = append(result, freeVariables(t_)...)
		}

		return result
	}

	panic(fmt.Sprintf("Unexpected type to get free variables: %s\n", type_.String()))
}

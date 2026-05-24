package constraint

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
)

type Equation struct {
	lhs  nodes.StellaType
	rhs  nodes.StellaType
	expr nodes.Node
}

func NewEquation(lhs, rhs nodes.StellaType, expr nodes.Node) Equation {
	return Equation{lhs: lhs, rhs: rhs, expr: expr}
}

func (e *Equation) isConstantOnly() bool {
	return !needReconstruction(e.lhs) && !needReconstruction(e.rhs)
}

func (e *Equation) apply(var_ *nodes.TypeVar, rhs nodes.StellaType) {
	e.lhs = applyToType(var_, rhs, e.lhs)
	e.rhs = applyToType(var_, rhs, e.rhs)
}

// func (e *Equation) Validate() *ConstraintError {
// 	result := IsCompatible(e.lhs, e.rhs)

// 	if result {
// 		return nil
// 	}

// 	return &ConstraintError{ErrorType: UNEXPECTED_TYPE, Expr: e.expr}
// }

// func IsCompatible(lhs, rhs nodes.StellaType) bool {
// 	if _, ok := lhs.(*nodes.TypeVar); ok {
// 		return true
// 	}

// 	if _, ok := rhs.(*nodes.TypeVar); ok {
// 		return true
// 	}

// 	switch rType := rhs.(type) {
// 	case *nodes.TypeBool:
// 		_, ok := lhs.(*nodes.TypeBool)
// 		return ok
// 	case *nodes.TypeNat:
// 		_, ok := lhs.(*nodes.TypeNat)
// 		return ok
// 	case *nodes.TypeUnit:
// 		_, ok := lhs.(*nodes.TypeUnit)
// 		return ok
// 	case *nodes.TypeBot:
// 		_, ok := lhs.(*nodes.TypeBot)
// 		return ok
// 	case *nodes.TypeTop:
// 		_, ok := lhs.(*nodes.TypeTop)
// 		return ok
// 	case *nodes.TypeList:
// 		if lType, ok := lhs.(*nodes.TypeList); ok {
// 			return IsCompatible(lType.Type_, rType.Type_)
// 		}
// 		return false
// 	case *nodes.TypeRef:
// 		if lType, ok := lhs.(*nodes.TypeRef); ok {
// 			return IsCompatible(lType.Type_, rType.Type_)
// 		}
// 		return false
// 	case *nodes.TypeParens:
// 		if lType, ok := lhs.(*nodes.TypeParens); ok {
// 			return IsCompatible(lType.Type_, rType.Type_)
// 		}
// 		return false
// 	case *nodes.TypeSum:
// 		result := true
// 		if lType, ok := lhs.(*nodes.TypeSum); ok {
// 			result = result && IsCompatible(lType.Left, rType.Left) && IsCompatible(lType.Right, lType.Right)
// 		}
// 		return result
// 	case *nodes.TypeTuple:
// 		// Assuming pairs

// 		result := true
// 		if lType, ok := lhs.(*nodes.TypeTuple); ok {
// 			result = result && IsCompatible(lType.Types[0], rType.Types[0]) && IsCompatible(lType.Types[1], lType.Types[1])
// 		}
// 		return result
// 	case *nodes.TypeRecord:
// 		result := true
// 		if lType, ok := lhs.(*nodes.TypeRecord); ok {
// 			for _, rFieldType := range rType.FieldTypes {
// 				for _, lFieldType := range lType.FieldTypes {
// 					if lFieldType.Label.Equal(&rFieldType.Label) {
// 						result = result && IsCompatible(lFieldType.Type_, rFieldType.Type_)
// 					}
// 				}
// 			}
// 		}
// 		return result
// 	case *nodes.TypeFun:
// 		result := true
// 		if lType, ok := lhs.(*nodes.TypeFun); ok {
// 			if len(rType.ParamTypes) != len(lType.ParamTypes) {
// 				return false
// 			}

// 			for index := range rType.ParamTypes {
// 				result = result && IsCompatible(lType.ParamTypes[index], rType.ParamTypes[index])
// 			}
// 			return result && IsCompatible(lType.ReturnType, rType.ReturnType)
// 		}
// 	}

// 	return false
// }

func applyToType(var_ *nodes.TypeVar, rhs, applyTo nodes.StellaType) nodes.StellaType {
	switch v := applyTo.(type) {
	case *nodes.TypeVar:
		if v.Name.Equal(&var_.Name) && v.Generated == var_.Generated {
			return rhs
		}
		return v
	case *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit, *nodes.TypeBot, *nodes.TypeTop:
		return v
	case *nodes.TypeFun:
		newParamTypes := make([]nodes.StellaType, 0, len(v.ParamTypes))

		for _, paramType := range v.ParamTypes {
			newParamTypes = append(newParamTypes, applyToType(var_, rhs, paramType))
		}

		newReturnType := applyToType(var_, rhs, v.ReturnType)

		return &nodes.TypeFun{ParamTypes: newParamTypes, ReturnType: newReturnType}
	case *nodes.TypeList:
		return &nodes.TypeList{Type_: applyToType(var_, rhs, v.Type_)}
	case *nodes.TypeParens:
		return applyToType(var_, rhs, v.Type_)
	case *nodes.TypeRecord:
		newFields := make([]nodes.RecordFieldType, 0, len(v.FieldTypes))

		for _, fieldType := range v.FieldTypes {
			newFieldType := nodes.RecordFieldType{Label: fieldType.Label, Type_: applyToType(var_, rhs, fieldType.Type_)}
			newFields = append(newFields, newFieldType)
		}

		return &nodes.TypeRecord{FieldTypes: newFields}
	case *nodes.TypeRef:
		return &nodes.TypeRef{Type_: applyToType(var_, rhs, v.Type_)}
	case *nodes.TypeSum:
		return &nodes.TypeSum{Left: applyToType(var_, rhs, v.Left), Right: applyToType(var_, rhs, v.Right)}
	case *nodes.TypeTuple:
		newTypes := make([]nodes.StellaType, 0, len(v.Types))

		for _, type_ := range v.Types {
			newTypes = append(newTypes, applyToType(var_, rhs, type_))
		}

		return &nodes.TypeTuple{Types: newTypes}
	default:
		panic(fmt.Sprintf("Unexpected type for application %s", applyTo.String()))
	}
}

// func resolveEquation(eq Equation) ([]Equation, *ConstraintError) {
// 	// Assuming that there is no equations with unmatched types (e.g. there is no Ref X = { Y, Z })

// 	// Checking if type var on top (if it is then it should be substituted before resolvement)
// 	if _, ok := eq.lhs.(*nodes.TypeVar); ok {
// 		return make([]Equation, 0), nil
// 	}
// 	if _, ok := eq.rhs.(*nodes.TypeVar); ok {
// 		return make([]Equation, 0), nil
// 	}

// 	newEquations, err := resolveType(eq.lhs, eq.rhs)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, equation := range newEquations {
// 		err := equation.Validate()

// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return newEquations, nil
// }

// func resolveType(lhs, rhs nodes.StellaType) ([]Equation, *ConstraintError) {
// 	// Assuming that there is no TypeVar on top level
// 	// Expected validation after return
// 	newEquations := make([]Equation, 0)
// 	switch lhsType := lhs.(type) {
// 	case *nodes.TypeFun:
// 		if rhsType, ok := rhs.(*nodes.TypeFun); ok {
// 			if len(lhsType.ParamTypes) != len(rhsType.ParamTypes) {
// 				return nil, &ConstraintError{"Unexpected number of parameters"}
// 			}

// 			for index := range lhsType.ParamTypes {
// 				newEquation := NewEquation(lhsType.ParamTypes[index], rhsType.ParamTypes[index])
// 				newEquations = append(newEquations, newEquation)
// 			}

// 			newEquation := NewEquation(lhsType.ReturnType, rhsType.ReturnType)
// 			newEquations = append(newEquations, newEquation)

// 			return newEquations, nil
// 		}
// 	case *nodes.TypeList:
// 		if rhsType, ok := rhs.(*nodes.TypeList); ok {
// 			newEquation := NewEquation(lhsType.Type_, rhsType.Type_)
// 			newEquations = append(newEquations, newEquation)
// 			return newEquations, nil
// 		}
// 	case *nodes.TypeRef:
// 		if rhsType, ok := rhs.(*nodes.TypeRef); ok {
// 			newEquations = append(newEquations, NewEquation(lhsType.Type_, rhsType.Type_))
// 			return newEquations, nil
// 		}
// 	case *nodes.TypeSum:
// 		if rhsType, ok := rhs.(*nodes.TypeSum); ok {
// 			newEquation1 := NewEquation(lhsType.Left, rhsType.Left)
// 			newEquation2 := NewEquation(lhsType.Right, rhsType.Right)

// 			newEquations = append(newEquations, newEquation1, newEquation2)

// 			return newEquations, nil
// 		}
// 	case *nodes.TypeTuple:
// 		// Expected only pairs
// 		if rhsType, ok := rhs.(*nodes.TypeTuple); ok {
// 			newEquation1 := NewEquation(lhsType.Types[0], rhsType.Types[0])
// 			newEquation2 := NewEquation(lhsType.Types[1], rhsType.Types[1])

// 			newEquations = append(newEquations, newEquation1, newEquation2)

// 			return newEquations, nil
// 		}
// 	case *nodes.TypeRecord:
// 		if rhsType, ok := rhs.(*nodes.TypeRecord); ok {

// 			newFieldTypes := make([]nodes.RecordFieldType, 0)
// 			gatherFieldTypes := make(map[nodes.StellaIdent]nodes.StellaType)

// 			// new equations for types from both sides
// 			for _, lhsFieldType := range lhsType.FieldTypes {
// 				if _, ok := gatherFieldTypes[lhsFieldType.Label]; !ok {
// 					gatherFieldTypes[lhsFieldType.Label] = lhsFieldType.Type_
// 				}
// 				for _, rhsFieldType := range rhsType.FieldTypes {
// 					if _, ok := gatherFieldTypes[rhsFieldType.Label]; !ok {
// 						gatherFieldTypes[rhsFieldType.Label] = lhsFieldType.Type_
// 					}
// 					if lhsFieldType.Label.Equal(&rhsFieldType.Label) {
// 						newEquations = append(newEquations, NewEquation(lhsFieldType.Type_, rhsFieldType.Type_))
// 					}
// 				}
// 			}

// 			// Unite all types
// 			for label, type_ := range gatherFieldTypes {
// 				newFieldTypes = append(newFieldTypes, nodes.RecordFieldType{Label: label, Type_: type_})
// 			}

// 			// TODO create new type var
// 			// newEquations = append(newEquations, NewEquation(&nodes.TypeRecord{FieldTypes: newFieldTypes}, TypeVar))

// 			return newEquations, nil
// 		}
// 	}

// 	message := fmt.Sprintf("Unexpected types to resolve: '%s' and '%s'", lhs.String(), rhs.String())
// 	return nil, &ConstraintError{Message: message}
// }

func IsInfiniteType(lhs *nodes.TypeVar, rhs nodes.StellaType) bool {
	switch rType := rhs.(type) {
	case *nodes.TypeVar:
		return lhs.Name.Equal(&rType.Name)
	case *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit, *nodes.TypeBot, *nodes.TypeTop:
		return false
	case *nodes.TypeList:
		return IsInfiniteType(lhs, rType.Type_)
	case *nodes.TypeRef:
		return IsInfiniteType(lhs, rType.Type_)
	case *nodes.TypeParens:
		return IsInfiniteType(lhs, rType.Type_)
	case *nodes.TypeSum:
		return IsInfiniteType(lhs, rType.Left) || IsInfiniteType(lhs, rType.Right)
	case *nodes.TypeTuple:
		// Assuming pairs
		return IsInfiniteType(lhs, rType.Types[0]) || IsInfiniteType(lhs, rType.Types[1])
	case *nodes.TypeRecord:
		result := false
		for _, fieldType := range rType.FieldTypes {
			result = result || IsInfiniteType(lhs, fieldType.Type_)
		}
		return result
	case *nodes.TypeFun:
		result := false
		for _, type_ := range rType.ParamTypes {
			result = result || IsInfiniteType(lhs, type_)
		}
		return result || IsInfiniteType(lhs, rType.ReturnType)
	}

	panic(fmt.Sprintf("Unexpected type to check for infinite type %s", rhs.String()))
}

// Constraints
// Pairs: {X1, X2} = {X3, X4}, X1 = {X2, X3}
// Records: X1 = {a: X2}, {b: X3} = {c: X4}
// Sum: X1 = X2 + X3, X1 + X2 = X3 + X4
// Ref: X1 = Ref X2, Ref X1 = Ref X2
// Int, Real, Unit as TT: X1 = TT

func (eq *Equation) String() string {
	return fmt.Sprintf("%s = %s", eq.lhs.String(), eq.rhs.String())
}

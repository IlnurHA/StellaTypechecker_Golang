package universaltypes

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"

	"github.com/neocotic/go-optional"
)

func ChangeTypeVar(from *nodes.StellaIdent, to nodes.StellaType, type_ nodes.StellaType) nodes.StellaType {
	switch t := type_.(type) {
	case *nodes.TypeVar:
		if t.Generated || !t.Name.Equal(from) {
			return t
		}
		return to
	case *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit, *nodes.TypeBot, *nodes.TypeTop:
		return t
	case *nodes.TypeFun:
		newParams := make([]nodes.StellaType, 0, len(t.ParamTypes))
		for _, param := range t.ParamTypes {
			newParams = append(newParams, ChangeTypeVar(from, to, param))
		}

		return &nodes.TypeFun{ParamTypes: newParams, ReturnType: ChangeTypeVar(from, to, t.ReturnType)}
	case *nodes.TypeList:
		return &nodes.TypeList{Type_: ChangeTypeVar(from, to, t.Type_)}
	case *nodes.TypeParens:
		return &nodes.TypeParens{Type_: ChangeTypeVar(from, to, t.Type_)}
	case *nodes.TypeRef:
		return &nodes.TypeRef{Type_: ChangeTypeVar(from, to, t.Type_)}
	case *nodes.TypeSum:
		return &nodes.TypeSum{
			Left:  ChangeTypeVar(from, to, t.Left),
			Right: ChangeTypeVar(from, to, t.Right),
		}
	case *nodes.TypeTuple:
		newTypes := make([]nodes.StellaType, 0, len(t.Types))

		for _, type_ := range t.Types {
			newTypes = append(newTypes, ChangeTypeVar(from, to, type_))
		}

		return &nodes.TypeTuple{Types: newTypes}
	case *nodes.TypeRecord:
		newFieldTypes := make([]nodes.RecordFieldType, 0, len(t.FieldTypes))
		for _, recordFieldType := range t.FieldTypes {
			newFieldTypes = append(newFieldTypes,
				nodes.RecordFieldType{
					Label: recordFieldType.Label,
					Type_: ChangeTypeVar(from, to, recordFieldType.Type_),
				},
			)
		}
		return &nodes.TypeRecord{FieldTypes: newFieldTypes}
	case *nodes.TypeVariant:
		newFieldTypes := make([]nodes.VariantFieldType, 0, len(t.FieldTypes))
		for _, variantFieldType := range t.FieldTypes {
			newType := optional.Empty[nodes.StellaType]()
			if variantFieldType.Type_.IsPresent() {
				newType = optional.Of(ChangeTypeVar(from, to, variantFieldType.Type_.Require()))
			}
			newFieldTypes = append(newFieldTypes,
				nodes.VariantFieldType{
					Label: variantFieldType.Label,
					Type_: newType,
				},
			)
		}
		return &nodes.TypeVariant{FieldTypes: newFieldTypes}
	case *nodes.TypeForAll:
		// Check if we cannot further overwrite 'from' typevar
		for _, generic := range t.Types {
			if generic.Equal(from) {
				return t
			}
		}
		// Check if we need to perform alpha conversion first
		for _, freevar := range freeVariables(to) {
			for i, generic := range t.Types {
				if generic.Equal(&freevar) {
					newName := nodes.StellaIdent{Name: generic.Repr + "_"}
					newTypeVar := nodes.TypeVar{Name: newName, Generated: false}
					subt := ChangeTypeVar(&freevar, &newTypeVar, t.Type_)
					newTypes := append(t.Types[:i], newName)
					newTypes = append(newTypes, t.Types[i+1:]...)
					t = &nodes.TypeForAll{Types: newTypes, Type_: subt}
					break
				}
			}
		}

		return &nodes.TypeForAll{
			Types: t.Types,
			Type_: ChangeTypeVar(from, to, t.Type_),
		}
	default:
		panic(fmt.Sprintf("Unimplemented changing type vars for %s\n", type_.String()))
	}
}

func freeVariables(type_ nodes.StellaType) []nodes.StellaIdent {
	switch t := type_.(type) {
	case *nodes.TypeVar:
		return []nodes.StellaIdent{t.Name}
	case *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit, *nodes.TypeBot, *nodes.TypeTop:
		return []nodes.StellaIdent{}
	case *nodes.TypeForAll:
		vars := freeVariables(t.Type_)
		freeVars := make([]nodes.StellaIdent, 0)

		for _, var_ := range vars {

			isFound := false
			for _, type_ := range t.Types {
				if type_.Equal(&var_) {
					isFound = true
					break
				}
			}

			if !isFound {
				freeVars = append(freeVars, var_)
			}
		}
		return freeVars
	case *nodes.TypeFun:
		result := allTypeForAllVars(t.ReturnType)

		for _, t_ := range t.ParamTypes {
			result = append(result, allTypeForAllVars(t_)...)
		}

		return result
	case *nodes.TypeList:
		return allTypeForAllVars(t.Type_)
	case *nodes.TypeParens:
		return allTypeForAllVars(t.Type_)
	case *nodes.TypeRecord:
		result := make([]nodes.StellaIdent, 0)
		for _, t_ := range t.FieldTypes {
			result = append(result, allTypeForAllVars(t_.Type_)...)
		}
		return result
	case *nodes.TypeRef:
		return allTypeForAllVars(t.Type_)
	case *nodes.TypeSum:
		result := allTypeForAllVars(t.Left)
		result = append(result, allTypeForAllVars(t.Right)...)
		return result
	case *nodes.TypeTuple:
		result := make([]nodes.StellaIdent, 0)

		for _, t_ := range t.Types {
			result = append(result, allTypeForAllVars(t_)...)
		}

		return result
	case *nodes.TypeVariant:
		result := make([]nodes.StellaIdent, 0)
		for _, t_ := range t.FieldTypes {
			if t_.Type_.IsEmpty() {
				continue
			}
			result = append(result, allTypeForAllVars(t_.Type_.Require())...)
		}
		return result
	}

	panic(fmt.Sprintf("Unexpected type to get free variables: %s\n", type_.String()))
}

func allTypeForAllVars(type_ nodes.StellaType) []nodes.StellaIdent {
	switch t := type_.(type) {
	case *nodes.TypeVar:
		return []nodes.StellaIdent{}
	case *nodes.TypeBool, *nodes.TypeNat, *nodes.TypeUnit, *nodes.TypeBot, *nodes.TypeTop:
		return []nodes.StellaIdent{}
	case *nodes.TypeForAll:
		result := t.Types[:]
		return append(result, allTypeForAllVars(t.Type_)...)
	case *nodes.TypeFun:
		result := allTypeForAllVars(t.ReturnType)

		for _, t_ := range t.ParamTypes {
			result = append(result, allTypeForAllVars(t_)...)
		}

		return result
	case *nodes.TypeList:
		return allTypeForAllVars(t.Type_)
	case *nodes.TypeParens:
		return allTypeForAllVars(t.Type_)
	case *nodes.TypeRecord:
		result := make([]nodes.StellaIdent, 0)
		for _, t_ := range t.FieldTypes {
			result = append(result, allTypeForAllVars(t_.Type_)...)
		}
		return result
	case *nodes.TypeRef:
		return allTypeForAllVars(t.Type_)
	case *nodes.TypeSum:
		result := allTypeForAllVars(t.Left)
		result = append(result, allTypeForAllVars(t.Right)...)
		return result
	case *nodes.TypeTuple:
		result := make([]nodes.StellaIdent, 0)

		for _, t_ := range t.Types {
			result = append(result, allTypeForAllVars(t_)...)
		}

		return result
	case *nodes.TypeVariant:
		result := make([]nodes.StellaIdent, 0)
		for _, t_ := range t.FieldTypes {
			if t_.Type_.IsEmpty() {
				continue
			}
			result = append(result, allTypeForAllVars(t_.Type_.Require())...)
		}
		return result
	}

	panic(fmt.Sprintf("Unexpected type to get free variables: %s\n", type_.String()))
}

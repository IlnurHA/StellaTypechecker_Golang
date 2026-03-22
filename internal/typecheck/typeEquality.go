package typecheck

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
)

func CheckStellaType(actual nodes.StellaType, expected nodes.StellaType) (err *TypecheckError) {
	defer func() {
		if err != nil {
			err.OverwriteActualType(actual)
			err.OverwriteExpectedType(expected)
		}
	}()
	// fmt.Printf("DEBUG: Checking actual %T against expected %T\n", actual, expected)
	switch lt := actual.(type) {
	case *nodes.TypeBool:
		_, ok := expected.(*nodes.TypeBool)

		if ok {
			return nil
		}

		err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
		return &err
	case *nodes.TypeNat:
		_, ok := expected.(*nodes.TypeNat)

		if ok {
			return nil
		}

		err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
		return &err
	case *nodes.TypeUnit:
		_, ok := expected.(*nodes.TypeUnit)

		if ok {
			return nil
		}

		err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
		return &err
	case *nodes.TypeList:
		rt, ok := expected.(*nodes.TypeList)

		if !ok {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			return &err
		}

		return CheckStellaType(lt.Type_, rt.Type_)
	case *nodes.TypeRef:
		rt, ok := expected.(*nodes.TypeRef)

		if !ok {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			return &err
		}

		return CheckStellaType(lt.Type_, rt.Type_)
	case *nodes.TypeSum:
		rt, ok := expected.(*nodes.TypeSum)

		if !ok {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			return &err
		}

		err := CheckStellaType(lt.Left, rt.Left)

		if err != nil {
			return err
		}

		return CheckStellaType(lt.Right, rt.Right)
	case *nodes.TypeFun:
		rt, ok := expected.(*nodes.TypeFun)

		if !ok {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			return &err
		}

		if len(rt.ParamTypes) != len(lt.ParamTypes) {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			err.AddAdditionalInfo(fmt.Sprintf("Expected %d params. Got %d", len(rt.ParamTypes), len(lt.ParamTypes)))
			return &err
		}

		res := CheckStellaType(lt.ReturnType, rt.ReturnType)

		index := 0
		for res == nil && index < len(rt.ParamTypes) {
			res = CheckStellaType(lt.ParamTypes[index], rt.ParamTypes[index])
			index++
		}

		return res
	case *nodes.TypeTuple:
		rt, ok := expected.(*nodes.TypeTuple)

		if !ok || len(rt.Types) != len(lt.Types) {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			return &err
		}

		var res *TypecheckError = nil
		index := 0
		for res == nil && index < len(rt.Types) {
			res = CheckStellaType(lt.Types[index], rt.Types[index])
			index++
		}

		return res
	case *nodes.TypeRecord:
		rt, ok := expected.(*nodes.TypeRecord)

		if !ok {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			return &err
		}

		notFoundFields := make([]*nodes.StellaIdent, 0)

		for _, rtField := range rt.FieldTypes {
			isFound := false

			for _, ltField := range lt.FieldTypes {
				if rtField.Label.Equal(&ltField.Label) {
					err := CheckStellaType(ltField.Type_, rtField.Type_)
					if err != nil {
						return err
					}
					isFound = true
					break
				}
			}

			if !isFound {
				notFoundFields = append(notFoundFields, &rtField.Label)
			}
		}

		if len(notFoundFields) != 0 {
			err := NewTypeCheckErrorErrorType(ERROR_MISSING_RECORD_FIELDS)
			err.AddAdditionalInfo(fmt.Sprintf("Missing record fields: %s", showList(notFoundFields)))
			return &err
		}

		if len(rt.FieldTypes) < len(lt.FieldTypes) {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_RECORD_FIELDS)
			return &err
		}

		return nil
	case *nodes.TypeVariant:
		rt, ok := expected.(*nodes.TypeVariant)

		if !ok || len(rt.FieldTypes) != len(lt.FieldTypes) {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			return &err
		}

		for _, rtField := range rt.FieldTypes {
			isFound := false

			for _, ltField := range lt.FieldTypes {
				if rtField.Label.Equal(&ltField.Label) {

					if rtField.Type_.IsEmpty() && ltField.Type_.IsPresent() {
						err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_DATA_FOR_NULLARY_LABEL)
						return &err
					}

					if rtField.Type_.IsPresent() && ltField.Type_.IsEmpty() {
						err := NewTypeCheckErrorErrorType(ERROR_MISSING_DATA_FOR_LABEL)
						return &err
					}

					if rtField.Type_.IsEmpty() && ltField.Type_.IsEmpty() {
						return nil
					}

					err := CheckStellaType(ltField.Type_.Require(), rtField.Type_.Require())
					if err != nil {
						return err
					}
					isFound = true
					break
				}
			}

			if !isFound {
				err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
				return &err
			}
		}

		return nil
	case *nodes.TypeParens:
		rt, ok := expected.(*nodes.TypeParens)

		if !ok {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
			return &err
		}

		return CheckStellaType(lt.Type_, rt.Type_)
	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		err.AddAdditionalInfo(fmt.Sprintf("Not implemented CheckStellaType for %s", actual))
		return &err
	}
}

func IsEqualType(actual nodes.StellaType, expected nodes.StellaType) bool {
	return CheckStellaType(actual, expected) == nil
}

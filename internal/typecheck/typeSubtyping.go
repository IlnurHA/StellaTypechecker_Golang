package typecheck

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
)

func SubtypeOf(left nodes.StellaType, right nodes.StellaType) (err *TypecheckError) {
	defer func() {
		if err != nil {
			err.AddIfEmptyActualType(right)
			err.AddIfEmptyExpectedType(left)
		}
	}()

	if _, ok := right.(*nodes.TypeTop); ok {
		return nil
	}
	if _, ok := left.(*nodes.TypeBot); ok {
		return nil
	}
	if rt, ok := right.(*nodes.TypeParens); ok {
		return SubtypeOf(left, rt.Type_)
	}
	if lt, ok := left.(*nodes.TypeParens); ok {
		return SubtypeOf(lt.Type_, right)
	}

	switch rt := right.(type) {
	case *nodes.TypeBool:
		if _, ok := left.(*nodes.TypeBool); ok {
			return nil
		}
	case *nodes.TypeNat:
		if _, ok := left.(*nodes.TypeNat); ok {
			return nil
		}
	case *nodes.TypeUnit:
		if _, ok := left.(*nodes.TypeNat); ok {
			return nil
		}
	case *nodes.TypeBot:
	case *nodes.TypeList:
		if lt, ok := left.(*nodes.TypeList); ok {
			return SubtypeOf(lt.Type_, rt.Type_)
		}
	case *nodes.TypeRef:
		if lt, ok := left.(*nodes.TypeRef); ok {
			err := SubtypeOf(lt.Type_, rt.Type_)
			if err != nil {
				return err
			}
			return SubtypeOf(rt.Type_, lt.Type_)
		}
	case *nodes.TypeSum:
		if lt, ok := left.(*nodes.TypeSum); ok {
			err := SubtypeOf(lt.Left, rt.Left)

			if err != nil {
				return nil
			}

			return SubtypeOf(lt.Right, rt.Right)
		}
	case *nodes.TypeRecord:
		if lt, ok := left.(*nodes.TypeRecord); ok {
			// Check for missing fields
			notFoundFields := make([]*nodes.StellaIdent, 0)
			for _, rft := range rt.FieldTypes {
				found := false

				for _, lft := range lt.FieldTypes {
					if lft.Label == rft.Label {
						found = true
					}
				}

				if !found {
					notFoundFields = append(notFoundFields, &rft.Label)
				}
			}

			if len(notFoundFields) != 0 {
				err := NewTypeCheckErrorErrorType(ERROR_MISSING_RECORD_FIELDS)
				err.AddAdditionalInfo(fmt.Sprintf("Missing record fields: %s", showList(notFoundFields)))
				return &err
			}

			// Check subtypes for field types
			for _, rft := range rt.FieldTypes {
				for _, lft := range lt.FieldTypes {
					if lft.Label == rft.Label {
						err := SubtypeOf(lft.Type_, rft.Type_)
						if err != nil {
							return err
						}
					}
				}
			}

			return nil
		}
	case *nodes.TypeFun:
		if lt, ok := left.(*nodes.TypeFun); ok {
			if len(rt.ParamTypes) != len(lt.ParamTypes) {
				err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION)
				err.AddAdditionalInfo(fmt.Sprintf("Expected %d params. Got %d", len(rt.ParamTypes), len(lt.ParamTypes)))
				return &err
			}

			res := SubtypeOf(lt.ReturnType, rt.ReturnType)

			index := 0
			for res == nil && index < len(rt.ParamTypes) {
				res = SubtypeOf(rt.ParamTypes[index], lt.ParamTypes[index])
				index++
			}

			return res
		}
	case *nodes.TypeTuple:
		if lt, ok := left.(*nodes.TypeTuple); ok {
			if len(rt.Types) > len(lt.Types) {
				err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_TUPLE_LENGTH)
				err.AddAdditionalInfo(fmt.Sprintf("Expected at least %d params. Got %d", len(rt.Types), len(lt.Types)))
				return &err
			}

			var res *TypecheckError = nil

			index := 0
			for res == nil && index < len(rt.Types) {
				res = SubtypeOf(lt.Types[index], rt.Types[index])
				index++
			}

			return res
		}
	case *nodes.TypeVariant:
		if lt, ok := left.(*nodes.TypeVariant); ok {
			// Check for extra fields
			extraFields := make([]*nodes.StellaIdent, 0)
			for _, lft := range lt.FieldTypes {
				found := false

				for _, rft := range rt.FieldTypes {
					if lft.Label == rft.Label {
						found = true
					}
				}

				if !found {
					extraFields = append(extraFields, &lft.Label)
				}
			}

			if len(extraFields) != 0 {
				err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_VARIANT_LABEL)
				err.AddAdditionalInfo(fmt.Sprintf("Extra variant fields: %s", showList(extraFields)))
				return &err
			}

			// Check subtypes for field types
			for _, lft := range lt.FieldTypes {
				for _, rft := range rt.FieldTypes {
					if lft.Label == rft.Label {
						if lft.Type_.IsEmpty() && rft.Type_.IsPresent() {
							err := NewTypeCheckErrorErrorType(ERROR_MISSING_DATA_FOR_LABEL)
							return &err
						}

						if lft.Type_.IsPresent() && rft.Type_.IsEmpty() {
							err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_DATA_FOR_NULLARY_LABEL)
							return &err
						}

						if lft.Type_.IsEmpty() && rft.Type_.IsEmpty() {
							continue
						}
						err := SubtypeOf(lft.Type_.Require(), rft.Type_.Require())
						if err != nil {
							return err
						}
					}
				}
			}

			return nil
		}
	}
	typeErr := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_SUBTYPE)
	return &typeErr

}

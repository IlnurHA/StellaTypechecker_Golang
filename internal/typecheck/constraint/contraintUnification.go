package constraint

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
)

type Substitution struct {
	substitutions map[nodes.StellaIdent]nodes.StellaType
}

func (s *Substitution) ApplyToAll(var_ *nodes.TypeVar, rhs nodes.StellaType) {
	for key, value := range s.substitutions {
		s.substitutions[key] = applyToType(var_, rhs, value)
	}
}

func (s *Substitution) String() string {
	message := fmt.Sprintln("Substitutions:")
	for key, value := range s.substitutions {
		message += fmt.Sprintf("%s : %s\n", key.String(), value.String())
	}
	return message
}

func (s *Substitution) HasAmbiguousTypes() *ConstraintError {
	for key, value := range s.substitutions {
		if len(freeVariables(value)) != 0 {
			return &ConstraintError{
				ErrorType:      AMBIGUOUS_TYPE,
				Lhs:            nil,
				Rhs:            nil,
				Expr:           nil,
				AdditionalInfo: fmt.Sprintf("Ambiguous type occured: %s = %s", key.String(), value.String()),
			}
		}
	}
	return nil
}

func Unify(constraints *ConstraintHandler) (*ConstraintError, *Substitution) {
	sub := Substitution{substitutions: make(map[nodes.StellaIdent]nodes.StellaType)}
	res := unify(constraints, &sub)
	return res, &sub
}

func unify(constraints *ConstraintHandler, substitution *Substitution) *ConstraintError {
	// Returns equation that failed

	for !constraints.IsEmpty() {
		constraint := constraints.PopEquation()
		expr := constraint.expr

		// TODO proper check for the same type
		if constraint.lhs == constraint.rhs {
			continue
		}

		if lVar, ok := constraint.lhs.(*nodes.TypeVar); ok {
			if hasLabel(&lVar.Name, freeVariables(constraint.rhs)) {
				return &ConstraintError{
					ErrorType:      INFINITE_TYPE,
					Lhs:            nil,
					Rhs:            nil,
					Expr:           constraint.expr,
					AdditionalInfo: fmt.Sprintf("Infinite type occured: %s = %s", lVar.String(), constraint.rhs.String()),
				}
			}

			if v, ok := substitution.substitutions[lVar.Name]; ok {
				constraints.AddEquationTypes(
					constraint.rhs, v,
				)
			}

			constraints.ApplyToAll(lVar, constraint.rhs)
			substitution.ApplyToAll(lVar, constraint.rhs)
			substitution.substitutions[lVar.Name] = constraint.rhs
			continue
		}

		if rVar, ok := constraint.rhs.(*nodes.TypeVar); ok {
			if hasLabel(&rVar.Name, freeVariables(constraint.lhs)) {
				return &ConstraintError{
					ErrorType:      INFINITE_TYPE,
					Lhs:            nil,
					Rhs:            nil,
					Expr:           constraint.expr,
					AdditionalInfo: fmt.Sprintf("Infinite type occured: %s = %s", rVar.String(), constraint.lhs.String()),
				}
			}

			if v, ok := substitution.substitutions[rVar.Name]; ok {
				constraints.AddEquationTypes(
					constraint.lhs, v,
				)
			}

			constraints.ApplyToAll(rVar, constraint.lhs)
			substitution.ApplyToAll(rVar, constraint.lhs)
			substitution.substitutions[rVar.Name] = constraint.lhs
			continue
		}

		if lFun, ok := constraint.lhs.(*nodes.TypeFun); ok {
			if rFun, ok := constraint.rhs.(*nodes.TypeFun); ok {
				eq := unifyFun(constraints, lFun, rFun, expr)

				if eq != nil {
					return eq
				}
				continue
			}

		}

		if lRec, ok := constraint.lhs.(*nodes.TypeRecord); ok {
			if rRec, ok := constraint.rhs.(*nodes.TypeRecord); ok {
				eq := unifyRecord(constraints, lRec, rRec, expr)

				if eq != nil {
					return eq
				}
				continue
			}
		}

		if lTuple, ok := constraint.lhs.(*nodes.TypeTuple); ok {
			if rTuple, ok := constraint.rhs.(*nodes.TypeTuple); ok {
				eq := unifyTuple(constraints, lTuple, rTuple, expr)

				if eq != nil {
					return eq
				}
				continue
			}
		}

		if lSum, ok := constraint.lhs.(*nodes.TypeSum); ok {
			if rSum, ok := constraint.rhs.(*nodes.TypeSum); ok {
				eq := unifySum(constraints, lSum, rSum, expr)

				if eq != nil {
					return eq
				}
				continue
			}
		}

		if lRef, ok := constraint.lhs.(*nodes.TypeRef); ok {
			if rRef, ok := constraint.rhs.(*nodes.TypeRef); ok {
				eq := unifyRef(constraints, lRef, rRef, expr)

				if eq != nil {
					return eq
				}
				continue
			}
		}

		return &ConstraintError{
			ErrorType: UNEXPECTED_TYPE,
			Lhs:       constraint.lhs,
			Rhs:       constraint.rhs,
			Expr:      constraint.expr,
		}
	}
	return nil
}

func unifyFun(constraints *ConstraintHandler, lhs *nodes.TypeFun, rhs *nodes.TypeFun, expr nodes.Node) *ConstraintError {
	if len(lhs.ParamTypes) != len(rhs.ParamTypes) {
		return &ConstraintError{
			ErrorType: UNEXPECTED_NUMBER_OF_PARAMETERS,
			Lhs:       lhs,
			Rhs:       rhs,
			Expr:      expr,
		}
	}

	for index := range lhs.ParamTypes {
		constraints.AddEquation(
			NewEquation(lhs.ParamTypes[index], rhs.ParamTypes[index], expr),
		)
	}

	constraints.AddEquation(
		NewEquation(lhs.ReturnType, rhs.ReturnType, expr),
	)

	return nil
}

func unifyRecord(constraints *ConstraintHandler, lhs *nodes.TypeRecord, rhs *nodes.TypeRecord, expr nodes.Node) *ConstraintError {
	if len(lhs.FieldTypes) != len(rhs.FieldTypes) {
		var errorType ConstraintErrorType

		if len(lhs.FieldTypes) > len(rhs.FieldTypes) {
			errorType = EXTRA_LABEL
		} else {
			errorType = MISSING_LABEL
		}

		return &ConstraintError{
			ErrorType: errorType,
			Lhs:       lhs,
			Rhs:       rhs,
			Expr:      expr,
		}
	}

	for _, lFieldType := range lhs.FieldTypes {
		found := false

		for _, rFieldType := range rhs.FieldTypes {
			if lFieldType.Label.Equal(&rFieldType.Label) {
				found = true

				constraints.AddEquation(
					NewEquation(lFieldType.Type_, rFieldType.Type_, expr),
				)
			}
		}

		if !found {
			return &ConstraintError{
				ErrorType: EXTRA_LABEL,
				Lhs:       lhs,
				Rhs:       rhs,
				Expr:      expr,
			}
		}
	}

	return nil
}

func unifyTuple(constraints *ConstraintHandler, lhs *nodes.TypeTuple, rhs *nodes.TypeTuple, expr nodes.Node) *ConstraintError {
	if len(lhs.Types) != len(rhs.Types) {
		return &ConstraintError{
			ErrorType: UNEXPECTED_LENGTH,
			Lhs:       lhs,
			Rhs:       rhs,
			Expr:      expr,
		}
	}

	for index := range lhs.Types {
		constraints.AddEquation(
			NewEquation(lhs.Types[index], rhs.Types[index], expr),
		)
	}

	return nil
}

func unifySum(constraints *ConstraintHandler, lhs *nodes.TypeSum, rhs *nodes.TypeSum, expr nodes.Node) *ConstraintError {
	constraints.AddEquation(
		NewEquation(lhs.Left, rhs.Left, expr),
	)

	constraints.AddEquation(
		NewEquation(lhs.Right, rhs.Right, expr),
	)

	return nil
}

func unifyRef(constraints *ConstraintHandler, lhs *nodes.TypeRef, rhs *nodes.TypeRef, expr nodes.Node) *ConstraintError {
	constraints.AddEquation(
		NewEquation(lhs.Type_, rhs.Type_, expr),
	)
	return nil
}

func hasLabel(var_ *nodes.StellaIdent, vars []nodes.StellaIdent) bool {
	for _, varCheck := range vars {
		if var_.Equal(&varCheck) {
			return true
		}
	}
	return false
}

package constraint

import nodes "typechecker/internal/ast/nodes"

type ConstraintError struct {
	ErrorType      ConstraintErrorType
	Lhs            nodes.StellaType
	Rhs            nodes.StellaType
	Expr           nodes.Node
	AdditionalInfo string
}

type ConstraintErrorType = int

const (
	INFINITE_TYPE ConstraintErrorType = iota
	UNEXPECTED_TYPE
	MISSING_LABEL
	EXTRA_LABEL
	UNEXPECTED_LENGTH
	UNEXPECTED_NUMBER_OF_PARAMETERS
	AMBIGUOUS_TYPE
)

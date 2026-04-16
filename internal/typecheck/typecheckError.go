package typecheck

import (
	"fmt"
	"strings"

	nodes "typechecker/internal/ast/nodes"

	"github.com/neocotic/go-optional"
)

type TypecheckError struct {
	errorType    ErrorType
	expectedType optional.Optional[nodes.StellaType]
	actualType   optional.Optional[nodes.StellaType]

	pattern      optional.Optional[nodes.Pattern]
	expr         optional.Optional[nodes.Node]
	functionName optional.Optional[nodes.StellaIdent]

	additionalInfo optional.Optional[string]
	freezed        bool
}

func NewTypeCheckErrorFull(
	errorType ErrorType,
	expectedType optional.Optional[nodes.StellaType],
	actualType optional.Optional[nodes.StellaType],

	pattern optional.Optional[nodes.Pattern],
	expr optional.Optional[nodes.Node],
	functionName optional.Optional[nodes.StellaIdent],

	additionalInfo optional.Optional[string],
) TypecheckError {
	return TypecheckError{
		errorType:    errorType,
		expectedType: expectedType,
		actualType:   actualType,

		pattern:      pattern,
		expr:         expr,
		functionName: functionName,

		additionalInfo: additionalInfo,
		freezed:        false,
	}
}

func NewTypeCheckErrorErrorType(
	errorType ErrorType,
) TypecheckError {
	return TypecheckError{
		errorType:    errorType,
		expectedType: optional.Empty[nodes.StellaType](),
		actualType:   optional.Empty[nodes.StellaType](),

		pattern:        optional.Empty[nodes.Pattern](),
		expr:           optional.Empty[nodes.Node](),
		functionName:   optional.Empty[nodes.StellaIdent](),
		additionalInfo: optional.Empty[string](),
		freezed:        false,
	}
}

func (err *TypecheckError) IsAmbiguousError() bool {
	return err.errorType == ERROR_AMBIGUOUS_LIST ||
		err.errorType == ERROR_AMBIGUOUS_SUM_TYPE ||
		err.errorType == ERROR_AMBIGUOUS_VARIANT_TYPE
}

func (err *TypecheckError) AddIfEmptyActualType(actualType nodes.StellaType) {
	if err.actualType.IsEmpty() && !err.freezed {
		err.actualType = optional.Of(actualType)
	}
}

func (err *TypecheckError) AddIfEmptyExpectedType(expectedType nodes.StellaType) {
	if err.expectedType.IsEmpty() && !err.freezed {
		err.expectedType = optional.Of(expectedType)
	}
}

func (err *TypecheckError) AddIfEmptyPattern(pattern nodes.Pattern) {
	if err.pattern.IsEmpty() && !err.freezed {
		err.pattern = optional.Of(pattern)
	}
}

func (err *TypecheckError) AddIfEmptyExpr(expr nodes.Node) {
	if err.expr.IsEmpty() && !err.freezed {
		err.expr = optional.Of(expr)
	}
}

func (err *TypecheckError) AddIfEmptyFunctionName(functionName nodes.StellaIdent) {
	if err.functionName.IsEmpty() {
		err.functionName = optional.Of(functionName)
	}
}

func (err *TypecheckError) RewriteErrorType(newErrorType ErrorType) {
	err.errorType = newErrorType
}

func (err *TypecheckError) AddAdditionalInfo(info string) {
	err.additionalInfo = optional.Of(err.additionalInfo.OrElse("") + info + "\n")
}

func (err *TypecheckError) OverwriteActualType(actualType nodes.StellaType) {
	err.actualType = optional.Of(actualType)
}

func (err *TypecheckError) OverwriteExpectedType(expectedType nodes.StellaType) {
	err.expectedType = optional.Of(expectedType)
}

func (err *TypecheckError) OverwritePattern(pattern nodes.Pattern) {
	err.pattern = optional.Of(pattern)
}

func (err *TypecheckError) Freeze() {
	err.freezed = true
}

const (
	errorType        = "Type Error: %s\n"
	errorDescription = "%s\n"
	expectedType     = "Expected type:\n\t%s\n"
	actualType       = "Actual type:\n\t%s\n"
	inPattern        = "in pattern: %s\n"
	inExpr           = "in expression: %s\n"
	inFunction       = "in function: %s\n"
	additionalInfo   = "%s"
)

func (err *TypecheckError) String() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, errorType, err.errorType.String())

	if err.additionalInfo.IsPresent() {
		var info = err.additionalInfo.Require()
		fmt.Fprintf(&builder, additionalInfo, info)
	}

	if err.expectedType.IsPresent() {
		var node = err.expectedType.Require()
		fmt.Fprintf(&builder, expectedType, node.String())
	}

	if err.actualType.IsPresent() {
		var node = err.actualType.Require()
		fmt.Fprintf(&builder, actualType, node.String())
	}

	if err.pattern.IsPresent() {
		var node = err.pattern.Require()
		fmt.Fprintf(&builder, inPattern, node.String())
	}

	if err.expr.IsPresent() {
		var node = err.expr.Require()
		fmt.Fprintf(&builder, inExpr, node.String())
	}

	if err.functionName.IsPresent() {
		var node = err.functionName.Require()
		fmt.Fprintf(&builder, inFunction, node.String())
	}

	return builder.String()
}

type ErrorType int

const (
	ERROR_MISSING_MAIN ErrorType = iota
	ERROR_UNDEFINED_VARIABLE
	ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION
	ERROR_NOT_A_FUNCTION
	ERROR_NOT_A_TUPLE
	ERROR_NOT_A_RECORD
	ERROR_NOT_A_LIST
	ERROR_UNEXPECTED_LAMBDA
	ERROR_UNEXPECTED_TYPE_FOR_PARAMETER
	ERROR_UNEXPECTED_TUPLE
	ERROR_UNEXPECTED_RECORD
	ERROR_UNEXPECTED_VARIANT
	ERROR_UNEXPECTED_LIST
	ERROR_UNEXPECTED_INJECTION
	ERROR_MISSING_RECORD_FIELDS
	ERROR_UNEXPECTED_RECORD_FIELDS
	ERROR_UNEXPECTED_FIELD_ACCESS
	ERROR_UNEXPECTED_VARIANT_LABEL
	ERROR_TUPLE_INDEX_OUT_OF_BOUNDS
	ERROR_UNEXPECTED_TUPLE_LENGTH
	ERROR_AMBIGUOUS_SUM_TYPE
	ERROR_AMBIGUOUS_VARIANT_TYPE
	ERROR_AMBIGUOUS_LIST
	ERROR_ILLEGAL_EMPTY_MATCHING
	ERROR_NONEXHAUSTIVE_MATCH_PATTERNS
	ERROR_UNEXPECTED_PATTERN_FOR_TYPE
	ERROR_DUPLICATE_RECORD_FIELDS
	ERROR_DUPLICATE_RECORD_TYPE_FIELDS
	ERROR_DUPLICATE_VARIANT_TYPE_FIELDS
	ERROR_DUPLICATE_FUNCTION_DECLARATION

	ERROR_INCORRECT_ARITY_OF_MAIN
	ERROR_INCORRECT_NUMBER_OF_ARGUMENTS
	ERROR_UNEXPECTED_NUMBER_OF_PARAMETERS_IN_LAMBDA

	ERROR_DUPLICATE_RECORD_PATTERN_FIELDS

	ERROR_UNEXPECTED_DATA_FOR_NULLARY_LABEL
	ERROR_MISSING_DATA_FOR_LABEL
	ERROR_UNEXPECTED_NON_NULLARY_VARIANT_PATTERN
	ERROR_UNEXPECTED_NULLARY_VARIANT_PATTERN

	ERROR_DUPLICATE_FUNCTION_PARAMETER
	ERROR_DUPLICATE_LET_BINDING
	ERROR_DUPLICATE_TYPE_PARAMETER

	ERROR_EXCEPTION_TYPE_NOT_DECLARED
	ERROR_AMBIGUOUS_THROW_TYPE
	ERROR_AMBIGUOUS_REFERENCE_TYPE
	ERROR_AMBIGUOUS_PANIC_TYPE
	ERROR_NOT_A_REFERENCE

	ERROR_UNEXPECTED_MEMORY_ADDRESS

	ERROR_UNEXPECTED_REFERENCE

	ERROR_UNEXPECTED_SUBTYPE

	UNIMPLEMENTED
)

func (s ErrorType) String() string {
	return []string{
		"ERROR_MISSING_MAIN",
		"ERROR_UNDEFINED_VARIABLE",
		"ERROR_UNEXPECTED_TYPE_FOR_EXPRESSION",
		"ERROR_NOT_A_FUNCTION",
		"ERROR_NOT_A_TUPLE",
		"ERROR_NOT_A_RECORD",
		"ERROR_NOT_A_LIST",
		"ERROR_UNEXPECTED_LAMBDA",
		"ERROR_UNEXPECTED_TYPE_FOR_PARAMETER",
		"ERROR_UNEXPECTED_TUPLE",
		"ERROR_UNEXPECTED_RECORD",
		"ERROR_UNEXPECTED_VARIANT",
		"ERROR_UNEXPECTED_LIST",
		"ERROR_UNEXPECTED_INJECTION",
		"ERROR_MISSING_RECORD_FIELDS",
		"ERROR_UNEXPECTED_RECORD_FIELDS",
		"ERROR_UNEXPECTED_FIELD_ACCESS",
		"ERROR_UNEXPECTED_VARIANT_LABEL",
		"ERROR_TUPLE_INDEX_OUT_OF_BOUNDS",
		"ERROR_UNEXPECTED_TUPLE_LENGTH",
		"ERROR_AMBIGUOUS_SUM_TYPE",
		"ERROR_AMBIGUOUS_VARIANT_TYPE",
		"ERROR_AMBIGUOUS_LIST",
		"ERROR_ILLEGAL_EMPTY_MATCHING",
		"ERROR_NONEXHAUSTIVE_MATCH_PATTERNS",
		"ERROR_UNEXPECTED_PATTERN_FOR_TYPE",
		"ERROR_DUPLICATE_RECORD_FIELDS",
		"ERROR_DUPLICATE_RECORD_TYPE_FIELDS",
		"ERROR_DUPLICATE_VARIANT_TYPE_FIELDS",
		"ERROR_DUPLICATE_FUNCTION_DECLARATION",

		"ERROR_INCORRECT_ARITY_OF_MAIN",
		"ERROR_INCORRECT_NUMBER_OF_ARGUMENTS",
		"ERROR_UNEXPECTED_NUMBER_OF_PARAMETERS_IN_LAMBDA",

		"ERROR_DUPLICATE_RECORD_PATTERN_FIELDS",

		"ERROR_UNEXPECTED_DATA_FOR_NULLARY_LABEL",
		"ERROR_MISSING_DATA_FOR_LABEL",
		"ERROR_UNEXPECTED_NON_NULLARY_VARIANT_PATTERN",
		"ERROR_UNEXPECTED_NULLARY_VARIANT_PATTERN",

		"ERROR_DUPLICATE_FUNCTION_PARAMETER",
		"ERROR_DUPLICATE_LET_BINDING",
		"ERROR_DUPLICATE_TYPE_PARAMETER",

		"ERROR_EXCEPTION_TYPE_NOT_DECLARED",
		"ERROR_AMBIGUOUS_THROW_TYPE",
		"ERROR_AMBIGUOUS_REFERENCE_TYPE",
		"ERROR_AMBIGUOUS_PANIC_TYPE",
		"ERROR_NOT_A_REFERENCE",

		"ERROR_UNEXPECTED_MEMORY_ADDRESS",

		"ERROR_UNEXPECTED_REFERENCE",

		"ERROR_UNEXPECTED_SUBTYPE",

		"NOT_IMPLEMENTED",
	}[s]
}

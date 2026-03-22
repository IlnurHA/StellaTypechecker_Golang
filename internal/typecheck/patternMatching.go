package typecheck

import (
	"fmt"
	"reflect"
	nodes "typechecker/internal/ast/nodes"
)

func checkPatternType(pattern nodes.Pattern, expectedType nodes.StellaType) (err *TypecheckError) {
	defer func() {
		if err != nil {
			err.OverwritePattern(pattern)
		}
	}()

	switch p := pattern.(type) {
	case *nodes.PatternVar:
		return nil
	case *nodes.PatternInl:
		if v, ok := expectedType.(*nodes.TypeSum); ok {
			return checkPatternType(p.Pattern, v.Left)
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_PATTERN_FOR_TYPE)
			err.AddIfEmptyExpectedType(expectedType)
			return &err
		}
	case *nodes.PatternInr:
		if v, ok := expectedType.(*nodes.TypeSum); ok {
			return checkPatternType(p.Pattern, v.Right)
		} else {
			err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_PATTERN_FOR_TYPE)
			err.AddIfEmptyExpectedType(expectedType)
			return &err
		}
	case *nodes.PatternVariant:
		if v, ok := expectedType.(*nodes.TypeVariant); ok {
			for _, variantFieldType := range v.FieldTypes {
				if variantFieldType.Label.Equal(&p.Label) {
					if p.Pattern.IsEmpty() && variantFieldType.Type_.IsEmpty() {
						return nil
					} else if p.Pattern.IsPresent() && variantFieldType.Type_.IsPresent() {
						return checkPatternType(p.Pattern.Require(), variantFieldType.Type_.Require())
					}
					// Pattern and type mismatch
					// Error occured

					if variantFieldType.Type_.IsEmpty() {
						err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_NON_NULLARY_VARIANT_PATTERN)
						return &err
					}
					err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_NULLARY_VARIANT_PATTERN)
					return &err
				}
			}
		}
		err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_PATTERN_FOR_TYPE)
		err.AddIfEmptyExpectedType(expectedType)
		return &err
	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
		err.AddAdditionalInfo(fmt.Sprintf("Not implemented checkPatternType for %s", pattern))
		return &err

	}
}

func checkPatternTypes(cases []nodes.MatchCase, expectedType nodes.StellaType) *TypecheckError {
	patterns := make([]nodes.Pattern, len(cases))

	for index, case_ := range cases {
		patterns[index] = case_.Pattern
	}

	return checkPatternTypesHelper(patterns, expectedType)
}

func checkPatternTypesHelper(patterns []nodes.Pattern, expectedType nodes.StellaType) *TypecheckError {
	for _, pattern := range patterns {
		err := checkPatternType(pattern, expectedType)
		if err != nil {
			return err
		}
	}
	return nil
}

func patternToContext(ctx *Context, pattern nodes.Pattern, expectedType nodes.StellaType) bool {
	// Expected to pattern to be checked by checkPatternTypes
	// Thus, ignoring errors
	// Returns false if duplicate variable. true otherwise
	switch p := pattern.(type) {
	case *nodes.PatternVar:
		return ctx.AddVar(p.Name, expectedType)
	case *nodes.PatternInl:
		if v, ok := expectedType.(*nodes.TypeSum); ok {
			return patternToContext(ctx, p.Pattern, v.Left)
		}
		return false
	case *nodes.PatternInr:
		if v, ok := expectedType.(*nodes.TypeSum); ok {
			return patternToContext(ctx, p.Pattern, v.Right)
		}
		return false
	case *nodes.PatternVariant:
		if v, ok := expectedType.(*nodes.TypeVariant); ok {
			for _, variantFieldType := range v.FieldTypes {
				if variantFieldType.Label.Equal(&p.Label) {
					if p.Pattern.IsEmpty() && variantFieldType.Type_.IsEmpty() {
						return false
					} else if p.Pattern.IsPresent() && variantFieldType.Type_.IsPresent() {
						return patternToContext(ctx, p.Pattern.Require(), variantFieldType.Type_.Require())
					}
				}
			}
		}
		return false
	default:
		return false
	}
}

func patternBindingsToContext(ctx *Context, patternBindings []nodes.PatternBinding) *TypecheckError {
	for _, patternBinding := range patternBindings {
		inferredType, err := infer(ctx, patternBinding.Rhs)

		if err != nil {
			return err
		}

		if !patternToContext(ctx, patternBinding.Pattern, inferredType) {
			err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_LET_BINDING)
			return &err
		}
	}
	return nil
}

func checkExhaustiveness(cases []nodes.MatchCase, expectedType nodes.StellaType) *TypecheckError {
	patterns := make([]nodes.Pattern, len(cases))

	for index, case_ := range cases {
		patterns[index] = case_.Pattern
	}

	return checkExhaustivenessHelper(patterns, expectedType)
}

func checkExhaustivenessHelper(patterns []nodes.Pattern, expectedType nodes.StellaType) *TypecheckError {
	switch v := expectedType.(type) {
	case *nodes.TypeVariant:
		usedLabels := make(map[nodes.StellaIdent]bool, len(v.FieldTypes))

		for _, pattern := range patterns {
			if _, ok := checkForPattern[*nodes.PatternVar](pattern); ok {
				return nil
			}

			variantPattern, _ := checkForPattern[*nodes.PatternVariant](pattern)

			usedLabels[variantPattern.Label] = true
		}

		if len(usedLabels) != len(v.FieldTypes) {
			err := NewTypeCheckErrorErrorType(ERROR_NONEXHAUSTIVE_MATCH_PATTERNS)
			return &err
		}

		splittedPatterns := splitPatternsVariant(patterns)

		for _, field := range v.FieldTypes {
			if field.Type_.IsEmpty() {
				continue
			}

			if v, ok := splittedPatterns[field.Label]; ok {
				err := checkExhaustivenessHelper(v, field.Type_.Require())
				if err != nil {
					return err
				}
			}
		}

		return nil
	case *nodes.TypeSum:
		usedInl := false
		usedInr := false

		for _, pattern := range patterns {
			if _, ok := checkForPattern[*nodes.PatternVar](pattern); ok {
				return nil
			}

			if _, ok := checkForPattern[*nodes.PatternInl](pattern); ok {
				usedInl = true
			} else if _, ok := checkForPattern[*nodes.PatternInr](pattern); ok {
				usedInr = true
			}
		}

		if !usedInl || !usedInr {
			err := NewTypeCheckErrorErrorType(ERROR_NONEXHAUSTIVE_MATCH_PATTERNS)
			return &err
		}

		inlPatterns, inrPatterns := splitPatternsSum(patterns)

		err := checkExhaustivenessHelper(inlPatterns, v.Left)

		if err != nil {
			return err
		}

		return checkExhaustivenessHelper(inrPatterns, v.Right)
	case *nodes.TypeVar:
		return nil
	case *nodes.TypeNat:
		for _, pattern := range patterns {
			if _, ok := checkForPattern[*nodes.PatternVar](pattern); ok {
				return nil
			}
		}
		err := NewTypeCheckErrorErrorType(ERROR_NONEXHAUSTIVE_MATCH_PATTERNS)
		return &err
	case *nodes.TypeBool:
		usedFalse, usedTrue := false, false

		for _, pattern := range patterns {
			// println(pattern.String())
			if _, ok := checkForPattern[*nodes.PatternVar](pattern); ok {
				return nil
			}

			if _, ok := checkForPattern[*nodes.PatternTrue](pattern); ok {
				usedTrue = true
			}

			if _, ok := checkForPattern[*nodes.PatternFalse](pattern); ok {
				usedFalse = true
			}
		}

		if usedTrue && usedFalse {
			return nil
		}

		err := NewTypeCheckErrorErrorType(ERROR_NONEXHAUSTIVE_MATCH_PATTERNS)
		return &err
	case *nodes.TypeUnit:
		for _, pattern := range patterns {
			if _, ok := checkForPattern[*nodes.PatternVar](pattern); ok {
				return nil
			}

			if _, ok := checkForPattern[*nodes.PatternUnit](pattern); ok {
				return nil
			}
		}
		err := NewTypeCheckErrorErrorType(ERROR_NONEXHAUSTIVE_MATCH_PATTERNS)
		return &err
	default:
		for _, pattern := range patterns {
			if _, ok := checkForPattern[*nodes.PatternVar](pattern); ok {
				return nil
			}
		}
	}

	err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
	err.AddAdditionalInfo(fmt.Sprintf("Not implemented checkExhaustivenessHelper for %s", expectedType))
	return &err
}

func splitPatternsVariant(patterns []nodes.Pattern) map[nodes.StellaIdent][]nodes.Pattern {
	variants := make(map[nodes.StellaIdent][]nodes.Pattern)

	for _, pattern := range patterns {
		variantPattern, _ := checkForPattern[*nodes.PatternVariant](pattern)

		if variantPattern.Pattern.IsEmpty() {
			continue
		}

		if v, ok := variants[variantPattern.Label]; ok {
			variants[variantPattern.Label] = append(v, variantPattern.Pattern.Require())
		} else {
			variants[variantPattern.Label] = append(make([]nodes.Pattern, 0), variantPattern.Pattern.Require())
		}
	}
	return variants
}

func splitPatternsSum(patterns []nodes.Pattern) ([]nodes.Pattern, []nodes.Pattern) {
	inlPatterns := make([]nodes.Pattern, 0)
	inrPatterns := make([]nodes.Pattern, 0)

	for _, pattern := range patterns {
		if inlPattern, ok := checkForPattern[*nodes.PatternInl](pattern); ok {
			inlPatterns = append(inlPatterns, inlPattern.Pattern)
		} else if inrPattern, ok := checkForPattern[*nodes.PatternInr](pattern); ok {
			inrPatterns = append(inrPatterns, inrPattern.Pattern)
		}
	}

	return inlPatterns, inrPatterns
}

func checkForPattern[T nodes.Pattern](pattern nodes.Pattern) (T, bool) {
	// Check if pattern (or subpattern for patterns that describe inner pattern) is the same
	switch v := pattern.(type) {
	case *nodes.PatternAsc:
		return checkForPattern[T](v.Pattern)
	case *nodes.ParenthesisedPattern:
		return checkForPattern[T](v.Pattern)
	default:
		fmt.Printf("trying to reflect %s with value of %s \n", pattern, reflect.ValueOf(pattern))
		return reflect.TypeAssert[T](reflect.ValueOf(pattern))
	}
}

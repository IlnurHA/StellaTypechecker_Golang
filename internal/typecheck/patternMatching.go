package typecheck

import nodes "typechecker/internal/ast/nodes"

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
				if variantFieldType.Label.Equal(p.Label) {
					if p.Pattern.IsEmpty() && variantFieldType.Type_.IsEmpty() {
						return nil
					} else if p.Pattern.IsPresent() && variantFieldType.Type_.IsPresent() {
						return checkPatternType(p.Pattern.Require(), variantFieldType.Type_.Require())
					}
					// Pattern and type mismatch
					// Error occured
					break
				}
			}
		}
		err := NewTypeCheckErrorErrorType(ERROR_UNEXPECTED_PATTERN_FOR_TYPE)
		err.AddIfEmptyExpectedType(expectedType)
		return &err
	default:
		err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
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
				if variantFieldType.Label.Equal(p.Label) {
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
			if _, ok := pattern.(*nodes.PatternVar); ok {
				return nil
			}

			variantPattern, _ := pattern.(*nodes.PatternVariant)

			usedLabels[variantPattern.Label] = true
		}

		if len(usedLabels) != len(v.FieldTypes) {
			err := NewTypeCheckErrorErrorType(ERROR_NONEXHAUSTIVE_MATCH_PATTERNS)
			return &err
		}

		splittedPatterns := splitPatternsVariant(patterns, v)

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
			if _, ok := pattern.(*nodes.PatternVar); ok {
				return nil
			}

			if _, ok := pattern.(*nodes.PatternInl); ok {
				usedInl = true
			} else {
				usedInr = true
			}
		}

		if !usedInl || !usedInr {
			err := NewTypeCheckErrorErrorType(ERROR_NONEXHAUSTIVE_MATCH_PATTERNS)
			return &err
		}

		inlPatterns, inrPatterns := splitPatternsSum(patterns, v)

		err := checkExhaustivenessHelper(inlPatterns, v.Left)

		if err != nil {
			return err
		}

		return checkExhaustivenessHelper(inrPatterns, v.Right)
	case *nodes.TypeVar:
		return nil
	}

	err := NewTypeCheckErrorErrorType(UNIMPLEMENTED)
	return &err
}

func splitPatternsVariant(patterns []nodes.Pattern, expectedType *nodes.TypeVariant) map[nodes.StellaIdent][]nodes.Pattern {
	variants := make(map[nodes.StellaIdent][]nodes.Pattern)

	for _, pattern := range patterns {
		variantPattern, _ := pattern.(*nodes.PatternVariant)
		if v, ok := variants[variantPattern.Label]; ok {
			variants[variantPattern.Label] = append(v, pattern)
		} else {
			variants[variantPattern.Label] = append(make([]nodes.Pattern, 1), pattern)
		}
	}
	return variants
}

func splitPatternsSum(patterns []nodes.Pattern, expectedType *nodes.TypeSum) ([]nodes.Pattern, []nodes.Pattern) {
	inlPatterns := make([]nodes.Pattern, 0)
	inrPatterns := make([]nodes.Pattern, 0)

	for _, pattern := range patterns {
		if _, ok := pattern.(*nodes.PatternInl); ok {
			inlPatterns = append(inlPatterns, pattern)
		} else {
			inrPatterns = append(inrPatterns, pattern)
		}
	}

	return inlPatterns, inrPatterns
}

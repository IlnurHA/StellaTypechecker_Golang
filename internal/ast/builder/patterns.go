package ast

import (
	"strconv"
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/parser"

	"github.com/neocotic/go-optional"
)

func (v *ASTBuilder) VisitPatternVariant(ctx *parser.PatternVariantContext) interface{} {
	label := parseStellaIdent(ctx.GetLabel())
	pattern := optional.Empty[nodes.Pattern]()

	if ctx.GetPattern_() != nil {
		pattern = optional.Of(parsePattern(ctx.GetPattern_(), v))
	}

	return &nodes.PatternVariant{
		Label:   label,
		Pattern: pattern,
		Repr:    getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternInl(ctx *parser.PatternInlContext) interface{} {
	pattern := parsePattern(ctx.GetPattern_(), v)

	return &nodes.PatternInl{
		Pattern: pattern,
		Repr:    getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternInr(ctx *parser.PatternInrContext) interface{} {
	pattern := parsePattern(ctx.GetPattern_(), v)

	return &nodes.PatternInr{
		Pattern: pattern,
		Repr:    getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternTuple(ctx *parser.PatternTupleContext) interface{} {
	patterns := parseListOfPattern(ctx.GetPatterns(), v)
	return &nodes.PatternTuple{
		Patterns: patterns,
		Repr:     getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternRecord(ctx *parser.PatternRecordContext) interface{} {
	patterns := parseListOfLabelledPattern(ctx.GetPatterns(), v)
	return &nodes.PatternRecord{
		Patterns: patterns,
		Repr:     getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternList(ctx *parser.PatternListContext) interface{} {
	patterns := parseListOfPattern(ctx.GetPatterns(), v)
	return &nodes.PatternList{
		Patterns: patterns,
		Repr:     getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternCons(ctx *parser.PatternConsContext) interface{} {
	head := parsePattern(ctx.GetHead(), v)
	tail := parsePattern(ctx.GetTail(), v)

	return &nodes.PatternCons{
		Head: head,
		Tail: tail,
		Repr: getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternTrue(ctx *parser.PatternTrueContext) interface{} {
	return &nodes.PatternTrue{
		Repr: getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternFalse(ctx *parser.PatternFalseContext) interface{} {
	return &nodes.PatternFalse{
		Repr: getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternUnit(ctx *parser.PatternUnitContext) interface{} {
	return &nodes.PatternUnit{
		Repr: getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternInt(ctx *parser.PatternIntContext) interface{} {
	res, _ := strconv.Atoi(ctx.GetN().GetText())
	return &nodes.PatternInt{
		N:    res,
		Repr: getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternSucc(ctx *parser.PatternSuccContext) interface{} {
	pattern := parsePattern(ctx.GetPattern_(), v)
	return &nodes.PatternSucc{
		Pattern: pattern,
		Repr:    getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternVar(ctx *parser.PatternVarContext) interface{} {
	name := parseStellaIdent(ctx.GetName())
	return &nodes.PatternVar{
		Name: name,
		Repr: getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternAsc(ctx *parser.PatternAscContext) interface{} {
	pattern := parsePattern(ctx.GetPattern_(), v)
	type_ := parseType(ctx.GetType_(), v)
	return &nodes.PatternAsc{
		Pattern: pattern,
		Type_:   type_,
		Repr:    getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitPatternCastAs(ctx *parser.PatternCastAsContext) interface{} {
	pattern := parsePattern(ctx.GetPattern_(), v)
	type_ := parseType(ctx.GetType_(), v)
	return &nodes.PatternCastAs{
		Pattern: pattern,
		Type_:   type_,
		Repr:    getOriginalCode(ctx, v),
	}
}

func (v *ASTBuilder) VisitParenthesisedPattern(ctx *parser.ParenthesisedPatternContext) interface{} {
	pattern := parsePattern(ctx.GetPattern_(), v)
	return &nodes.ParenthesisedPattern{
		Pattern: pattern,
		Repr:    getOriginalCode(ctx, v),
	}
}

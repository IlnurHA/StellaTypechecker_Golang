package ast

import (
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/parser"

	"github.com/antlr4-go/antlr/v4"
	"github.com/neocotic/go-optional"
)

func parseExtensions(ctx *parser.ProgramContext, v *ASTBuilder) []nodes.Extension {
	total_extension_size := 0
	extensions := make([][]string, len(ctx.GetExtensions()))
	for _, extension := range ctx.GetExtensions() {
		var exts = extension.Accept(v).([]string)
		extensions = append(extensions, exts)
		total_extension_size += len(exts)
	}

	total_extensions := make([]nodes.Extension, total_extension_size)
	i := 0
	for _, exts := range extensions {
		for _, subexts := range exts {
			total_extensions[i] = nodes.Extension{Name: subexts}
			i++
		}
	}

	return total_extensions
}

func parseParameterDeclarations(params []parser.IParamDeclContext, v *ASTBuilder) []nodes.ParameterDeclaration {
	if params == nil {
		return []nodes.ParameterDeclaration{}
	}

	parameters := make([]nodes.ParameterDeclaration, len(params))
	for index, param := range params {
		parameters[index] = param.Accept(v).(nodes.ParameterDeclaration)
	}
	return parameters
}

func parseLocalDeclarations(decls []parser.IDeclContext, v *ASTBuilder) []nodes.Declaration {
	declarations := make([]nodes.Declaration, len(decls))
	for index, decl := range decls {
		declarations[index] = decl.Accept(v).(nodes.Declaration)
	}
	return declarations
}

func parseType(typeCtx parser.IStellatypeContext, v *ASTBuilder) nodes.StellaType {
	return (typeCtx.Accept(v)).(nodes.StellaType)
}

func parseListOfTypes(typesCtx []parser.IStellatypeContext, v *ASTBuilder) []nodes.StellaType {
	return parseList(typesCtx, parseType, v)
}

func parseListOfStellaIdent(identsCtx []antlr.Token) []nodes.StellaIdent {
	return map_(identsCtx, parseStellaIdent)
}

func parseStellaIdent(identCtx antlr.Token) nodes.StellaIdent {
	return nodes.StellaIdent{Name: identCtx.GetText(), Repr: identCtx.GetText()}
}

func parseDeclarations(declarationsCtx []parser.IDeclContext, v *ASTBuilder) []nodes.Declaration {
	declarations := make([]nodes.Declaration, len(declarationsCtx))

	for index, declarationCtx := range declarationsCtx {
		declarations[index] = declarationCtx.Accept(v).(nodes.Declaration)
	}

	return declarations
}

func parsePattern(patternCtx parser.IPatternContext, v *ASTBuilder) nodes.Pattern {
	return patternCtx.Accept(v).(nodes.Pattern)
}

func parseListOfPattern(patternsCtx []parser.IPatternContext, v *ASTBuilder) []nodes.Pattern {
	return parseList(patternsCtx, parsePattern, v)
}

func parseLabelledPattern(patternCtx parser.ILabelledPatternContext, v *ASTBuilder) nodes.LabelledPattern {
	label := parseStellaIdent(patternCtx.GetLabel())
	pattern := v.Visit(patternCtx.GetPattern_()).(nodes.Pattern)
	return nodes.LabelledPattern{Label: label, Pattern: pattern, Repr: patternCtx.GetText()}
}

func parseListOfLabelledPattern(patternsCtx []parser.ILabelledPatternContext, v *ASTBuilder) []nodes.LabelledPattern {
	return parseList(patternsCtx, parseLabelledPattern, v)
}

func parseRecordFieldType(typeCtx parser.IRecordFieldTypeContext, v *ASTBuilder) nodes.RecordFieldType {
	label := parseStellaIdent(typeCtx.GetLabel())
	type_ := parseType(typeCtx.GetType_(), v)

	return nodes.RecordFieldType{Label: label, Type_: type_, Repr: typeCtx.GetText()}
}

func parseListOfRecordFieldType(typesCtx []parser.IRecordFieldTypeContext, v *ASTBuilder) []nodes.RecordFieldType {

	return parseList(typesCtx, parseRecordFieldType, v)
}

func parseVariantFieldType(typeCtx parser.IVariantFieldTypeContext, v *ASTBuilder) nodes.VariantFieldType {
	label := parseStellaIdent(typeCtx.GetLabel())
	type_ := optional.Empty[nodes.StellaType]()

	if typeCtx.GetType_() != nil {
		type_ = optional.Of(parseType(typeCtx.GetType_(), v))
	}

	return nodes.VariantFieldType{Label: label, Type_: type_, Repr: typeCtx.GetText()}
}

func parseListOfVariantFieldType(typesCtx []parser.IVariantFieldTypeContext, v *ASTBuilder) []nodes.VariantFieldType {
	return parseList(typesCtx, parseVariantFieldType, v)
}

func parseExpr(exprCtx parser.IExprContext, v *ASTBuilder) nodes.Expr {
	// res := exprCtx.Accept(v)

	// if res == nil {
	// 	println("Got nil")
	// 	println("Expr:", exprCtx.GetText())
	// 	println("Rule: ", exprCtx.GetParser().GetRuleNames()[exprCtx.GetRuleIndex()])
	// 	println("Altarnative rule index", exprCtx.GetAltNumber())
	// }
	return exprCtx.Accept(v).(nodes.Expr)
}

func parseListOfExpr(exprsCtx []parser.IExprContext, v *ASTBuilder) []nodes.Expr {
	return parseList(exprsCtx, parseExpr, v)
}

func parsePatternBinding(patternBindingCtx parser.IPatternBindingContext, v *ASTBuilder) nodes.PatternBinding {
	pattern := parsePattern(patternBindingCtx.GetPat(), v)
	expr := parseExpr(patternBindingCtx.GetRhs(), v)

	return nodes.PatternBinding{Pattern: pattern, Rhs: expr, Repr: patternBindingCtx.GetText()}
}

func parseListOfPatternBinding(patternBindingsCtx []parser.IPatternBindingContext, v *ASTBuilder) []nodes.PatternBinding {
	return parseList(patternBindingsCtx, parsePatternBinding, v)
}

func parseBinding(bindingCtx parser.IBindingContext, v *ASTBuilder) nodes.Binding {
	name := parseStellaIdent(bindingCtx.GetName())
	expr := parseExpr(bindingCtx.GetRhs(), v)

	return nodes.Binding{Name: name, Rhs: expr, Repr: bindingCtx.GetText()}
}

func parseListOfBinding(bindingsCtx []parser.IBindingContext, v *ASTBuilder) []nodes.Binding {
	return parseList(bindingsCtx, parseBinding, v)
}

func parseMatchCase(matchCaseCtx parser.IMatchCaseContext, v *ASTBuilder) nodes.MatchCase {
	pattern := parsePattern(matchCaseCtx.GetPattern_(), v)
	expr := parseExpr(matchCaseCtx.GetExpr_(), v)

	return nodes.MatchCase{Pattern: pattern, Expr_: expr, Repr: matchCaseCtx.GetText()}
}

func parseListOfMatchCase(matchCaseCtx []parser.IMatchCaseContext, v *ASTBuilder) []nodes.MatchCase {
	return parseList(matchCaseCtx, parseMatchCase, v)
}

func parseList[T, V, B any](ctx []T, f func(T, B) V, builder B) []V {
	return map_(ctx, callWithBuilder(builder, f))
}

func callWithBuilder[T, V, B any](builder B, f func(T, B) V) func(T) V {
	return func(ctx T) V {
		return f(ctx, builder)
	}
}

func map_[T, V any](ctx []T, f func(T) V) []V {
	if ctx == nil {
		return []V{}
	}

	list := make([]V, len(ctx))
	for index, subctx := range ctx {
		list[index] = f(subctx)
	}
	return list
}

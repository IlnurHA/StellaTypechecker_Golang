package ast

import (
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/parser"
)

type ASTBuilder struct {
	*parser.BasestellaParserVisitor
}

func NewASTBuilder() *ASTBuilder {
	return &ASTBuilder{
		BasestellaParserVisitor: &parser.BasestellaParserVisitor{},
	}
}

func (v *ASTBuilder) VisitStart_Program(ctx *parser.Start_ProgramContext) interface{} {
	return ctx.Program().Accept(v)
}

func (v *ASTBuilder) VisitStart_Expr(ctx *parser.Start_ExprContext) interface{} {
	return parseExpr(ctx.Expr(), v)
}

func (v *ASTBuilder) VisitStart_Type(ctx *parser.Start_TypeContext) interface{} {
	return parseType(ctx.Stellatype(), v)
}

func (v *ASTBuilder) VisitLanguageCore(ctx *parser.LanguageCoreContext) interface{} {
	return nodes.LanguageDeclaration{Name: "core", Repr: ctx.GetText()}
}

func (v *ASTBuilder) VisitProgram(ctx *parser.ProgramContext) interface{} {
	languageDecl := ctx.LanguageDecl().Accept(v).(nodes.LanguageDeclaration)
	extensions := parseExtensions(ctx, v)
	declarations := parseDeclarations(ctx.GetDecls(), v)

	return nodes.AProgram{
		LanguageDecl: languageDecl,
		Extensions:   extensions,
		Declarations: declarations,
		Repr:         ctx.GetText(),
	}
}

func (v *ASTBuilder) VisitAnExtension(ctx *parser.AnExtensionContext) interface{} {

	var extensions = make([]string, len(ctx.GetExtensionNames()))

	for _, extensionTokens := range ctx.GetExtensionNames() {
		extensions = append(extensions, extensionTokens.GetText())
	}
	return extensions
}

func (v *ASTBuilder) VisitParamDecl(ctx *parser.ParamDeclContext) interface{} {
	type_ := parseType(ctx.GetParamType(), v)
	name := parseStellaIdent(ctx.GetName())
	return nodes.ParameterDeclaration{Name: name, ParameterType: type_}
}

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

func (v *ASTBuilder) VisitLanguageCore(ctx *parser.LanguageCoreContext) interface{} {
	return nodes.LanguageDeclaration{Name: "core"}
}

func (v *ASTBuilder) VisitProgram(ctx *parser.ProgramContext) interface{} {
	languageDecl := v.Visit(ctx.GetChildOfType(0, nil)).(nodes.LanguageDeclaration)
	extensions := parseExtensions(ctx, v)
	declarations := parseDeclarations(ctx.GetDecls(), v)

	return nodes.AProgram{LanguageDecl: languageDecl, Extensions: extensions, Declarations: declarations}
}

func (v *ASTBuilder) VisitAnExtension(ctx *parser.AnExtensionContext) interface{} {

	var extensions = make([]string, len(ctx.GetExtensionNames()))

	for _, extensionTokens := range ctx.GetExtensionNames() {
		extensions = append(extensions, extensionTokens.GetText())
	}
	return extensions
}

package ast

import (
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/parser"

	"github.com/neocotic/go-optional"
)

func (v *ASTBuilder) VisitDeclTypeAlias(ctx *parser.DeclTypeAliasContext) interface{} {
	name := parseStellaIdent(ctx.GetName())
	type_ := v.Visit(ctx.GetChildOfType(0, nil)).(nodes.StellaType)

	return nodes.TypeAliasDeclaration{Name: name, Atype: type_}
}

func (v *ASTBuilder) VisitDeclExceptionVariant(ctx *parser.DeclExceptionVariantContext) interface{} {
	name := parseStellaIdent(ctx.GetName())
	type_ := v.Visit(ctx.GetChildOfType(0, nil)).(nodes.StellaType)

	return nodes.ExceptionVariantDeclaration{Name: name, VariantType: type_}
}

func (v *ASTBuilder) VisitDeclExceptionType(ctx *parser.DeclExceptionTypeContext) interface{} {
	type_ := v.Visit(ctx.GetChildOfType(0, nil)).(nodes.StellaType)
	return nodes.ExceptionTypeDeclaration{ExceptionType: type_}
}

func (v *ASTBuilder) VisitDeclFun(ctx *parser.DeclFunContext) interface{} {
	name := parseStellaIdent(ctx.GetName())
	parameters := parseParameterDeclarations(ctx.GetParamDecls(), v)
	localDeclarations := parseLocalDeclarations(ctx.GetLocalDecls(), v)
	throwTypes := parseListOfTypes(ctx.GetThrowTypes(), v)

	returnType := optional.Empty[nodes.StellaType]()
	if ctx.GetReturnType() != nil {
		returnType = optional.Of(v.Visit(ctx.GetReturnType()).(nodes.StellaType))
	}
	expr := v.Visit(ctx.GetReturnExpr()).(nodes.Expr)

	return nodes.FunctionDeclaration{Name: name, Params: parameters, Declarations: localDeclarations, ReturnType: returnType, ThrowTypes: throwTypes, Expr: expr}
}

func (v *ASTBuilder) VisitDeclFunGeneric(ctx *parser.DeclFunGenericContext) interface{} {
	generics := parseListOfStellaIdent(ctx.GetGenerics())

	name := parseStellaIdent(ctx.GetName())
	parameters := parseParameterDeclarations(ctx.GetParamDecls(), v)
	localDeclarations := parseLocalDeclarations(ctx.GetLocalDecls(), v)
	throwTypes := parseListOfTypes(ctx.GetThrowTypes(), v)

	returnType := optional.Empty[nodes.StellaType]()
	if ctx.GetReturnType() != nil {
		returnType = optional.Of(v.Visit(ctx.GetReturnType()).(nodes.StellaType))
	}
	expr := v.Visit(ctx.GetReturnExpr()).(nodes.Expr)

	return nodes.GenericFunctionDeclaration{Name: name, Params: parameters, Declarations: localDeclarations, ReturnType: returnType, ThrowTypes: throwTypes, Expr: expr, Generics: generics}
}

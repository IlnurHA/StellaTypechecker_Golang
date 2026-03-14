package ast

import "github.com/neocotic/go-optional"

type ParameterDeclaration struct {
	name          StellaIdent
	parameterType StellaType
}

type Annotation struct {
	name string
}

type Declaration interface{ isDecl() }

type FunctionDeclaration struct {
	name         string
	params       []ParameterDeclaration
	returnType   optional.Optional[StellaType]
	throwTypes   optional.Optional[[]StellaType]
	declarations []Declaration
	expr         Expr
}

type GenericFunctionDeclaration struct {
	name         string
	params       []ParameterDeclaration
	returnType   optional.Optional[StellaType]
	throwTypes   optional.Optional[[]StellaType]
	declarations []Declaration
	expr         Expr

	annotations []Annotation
	generics    []StellaIdent
}

type TypeAliasDeclaration struct {
	name  string
	atype StellaType
}

type ExceptionTypeDeclaration struct {
	exceptionType StellaType
}

type ExceptionVariantDeclaration struct {
	name        string
	variantType StellaType
}

func (x *ParameterDeclaration) isNode() {
}
func (x *Annotation) isNode() {
}
func (x *FunctionDeclaration) isNode() {
}
func (x *FunctionDeclaration) isDecl() {
}
func (x *GenericFunctionDeclaration) isNode() {
}
func (x *GenericFunctionDeclaration) isDecl() {
}
func (x *TypeAliasDeclaration) isNode() {
}
func (x *TypeAliasDeclaration) isDecl() {
}
func (x *ExceptionTypeDeclaration) isNode() {
}
func (x *ExceptionTypeDeclaration) isDecl() {
}
func (x *ExceptionVariantDeclaration) isNode() {
}
func (x *ExceptionVariantDeclaration) isDecl() {
}

package ast

import "github.com/neocotic/go-optional"

type ParameterDeclaration struct {
	Name          StellaIdent
	ParameterType StellaType
}

type Annotation struct {
	Name string
}

type Declaration interface{ isDecl() }

type FunctionDeclaration struct {
	Name         StellaIdent
	Params       []ParameterDeclaration
	ReturnType   optional.Optional[StellaType]
	ThrowTypes   []StellaType
	Declarations []Declaration
	Expr         Expr
}

type GenericFunctionDeclaration struct {
	Name         StellaIdent
	Params       []ParameterDeclaration
	ReturnType   optional.Optional[StellaType]
	ThrowTypes   []StellaType
	Declarations []Declaration
	Expr         Expr

	Generics []StellaIdent
}

type TypeAliasDeclaration struct {
	Name  StellaIdent
	Atype StellaType
}

type ExceptionTypeDeclaration struct {
	ExceptionType StellaType
}

type ExceptionVariantDeclaration struct {
	Name        StellaIdent
	VariantType StellaType
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

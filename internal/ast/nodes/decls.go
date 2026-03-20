package ast

import "github.com/neocotic/go-optional"

type ParameterDeclaration struct {
	Name          StellaIdent
	ParameterType StellaType

	Repr string
}

type Annotation struct {
	Name string

	Repr string
}

type Declaration interface {
	Node
	isDecl()
}

type FunctionDeclaration struct {
	Name         StellaIdent
	Params       []ParameterDeclaration
	ReturnType   optional.Optional[StellaType]
	ThrowTypes   []StellaType
	Declarations []Declaration
	Expr         Expr

	Repr string
}

type GenericFunctionDeclaration struct {
	Name         StellaIdent
	Params       []ParameterDeclaration
	ReturnType   optional.Optional[StellaType]
	ThrowTypes   []StellaType
	Declarations []Declaration
	Expr         Expr

	Generics []StellaIdent

	Repr string
}

type TypeAliasDeclaration struct {
	Name  StellaIdent
	Atype StellaType

	Repr string
}

type ExceptionTypeDeclaration struct {
	ExceptionType StellaType

	Repr string
}

type ExceptionVariantDeclaration struct {
	Name        StellaIdent
	VariantType StellaType

	Repr string
}

// Marker methods
func (x *ParameterDeclaration) isNode()        {}
func (x *Annotation) isNode()                  {}
func (x *FunctionDeclaration) isNode()         {}
func (x *FunctionDeclaration) isDecl()         {}
func (x *GenericFunctionDeclaration) isNode()  {}
func (x *GenericFunctionDeclaration) isDecl()  {}
func (x *TypeAliasDeclaration) isNode()        {}
func (x *TypeAliasDeclaration) isDecl()        {}
func (x *ExceptionTypeDeclaration) isNode()    {}
func (x *ExceptionTypeDeclaration) isDecl()    {}
func (x *ExceptionVariantDeclaration) isNode() {}
func (x *ExceptionVariantDeclaration) isDecl() {}

// String methods returning Repr
func (x *ParameterDeclaration) String() string        { return x.Repr }
func (x *Annotation) String() string                  { return x.Repr }
func (x *FunctionDeclaration) String() string         { return x.Repr }
func (x *GenericFunctionDeclaration) String() string  { return x.Repr }
func (x *TypeAliasDeclaration) String() string        { return x.Repr }
func (x *ExceptionTypeDeclaration) String() string    { return x.Repr }
func (x *ExceptionVariantDeclaration) String() string { return x.Repr }

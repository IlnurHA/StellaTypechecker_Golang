package ast

import "github.com/neocotic/go-optional"

type RecordFieldType struct {
	Label StellaIdent
	Type_ StellaType

	Repr string
}

type VariantFieldType struct {
	Label StellaIdent
	Type_ optional.Optional[StellaType]

	Repr string
}

type StellaType interface {
	Node
	isType()
}

type TypeBool struct {
	Repr string
}
type TypeNat struct {
	Repr string
}

type TypeRef struct {
	Type_ StellaType

	Repr string
}
type TypeSum struct {
	Left  StellaType
	Right StellaType

	Repr string
}
type TypeFun struct {
	ParamTypes []StellaType
	ReturnType StellaType

	Repr string
}
type TypeForAll struct {
	Types []StellaIdent
	Type_ StellaType

	Repr string
}
type TypeRec struct {
	Var_  StellaIdent
	Type_ StellaType

	Repr string
}
type TypeTuple struct {
	Types []StellaType

	Repr string
}
type TypeRecord struct {
	FieldTypes []RecordFieldType

	Repr string
}
type TypeVariant struct {
	FieldTypes []VariantFieldType

	Repr string
}
type TypeList struct {
	Type_ StellaType

	Repr string
}

type TypeUnit struct {
	Repr string
}
type TypeTop struct {
	Repr string
}
type TypeBot struct {
	Repr string
}
type TypeAuto struct {
	Repr string
}
type TypeVar struct {
	Name StellaIdent

	Repr string
}
type TypeParens struct {
	Type_ StellaType

	Repr string
}

// Marker methods
func (x *RecordFieldType) isNode()  {}
func (x *VariantFieldType) isNode() {}
func (x *TypeBool) isNode()         {}
func (x *TypeBool) isType()         {}
func (x *TypeNat) isNode()          {}
func (x *TypeNat) isType()          {}
func (x *TypeRef) isNode()          {}
func (x *TypeRef) isType()          {}
func (x *TypeSum) isNode()          {}
func (x *TypeSum) isType()          {}
func (x *TypeFun) isNode()          {}
func (x *TypeFun) isType()          {}
func (x *TypeForAll) isNode()       {}
func (x *TypeForAll) isType()       {}
func (x *TypeRec) isNode()          {}
func (x *TypeRec) isType()          {}
func (x *TypeTuple) isNode()        {}
func (x *TypeTuple) isType()        {}
func (x *TypeRecord) isNode()       {}
func (x *TypeRecord) isType()       {}
func (x *TypeVariant) isNode()      {}
func (x *TypeVariant) isType()      {}
func (x *TypeList) isNode()         {}
func (x *TypeList) isType()         {}
func (x *TypeUnit) isNode()         {}
func (x *TypeUnit) isType()         {}
func (x *TypeTop) isNode()          {}
func (x *TypeTop) isType()          {}
func (x *TypeBot) isNode()          {}
func (x *TypeBot) isType()          {}
func (x *TypeAuto) isNode()         {}
func (x *TypeAuto) isType()         {}
func (x *TypeVar) isNode()          {}
func (x *TypeVar) isType()          {}
func (x *TypeParens) isNode()       {}
func (x *TypeParens) isType()       {}

// String methods returning Repr
func (x *RecordFieldType) String() string  { return x.Repr }
func (x *VariantFieldType) String() string { return x.Repr }
func (x *TypeBool) String() string         { return x.Repr }
func (x *TypeNat) String() string          { return x.Repr }
func (x *TypeRef) String() string          { return x.Repr }
func (x *TypeSum) String() string          { return x.Repr }
func (x *TypeFun) String() string          { return x.Repr }
func (x *TypeForAll) String() string       { return x.Repr }
func (x *TypeRec) String() string          { return x.Repr }
func (x *TypeTuple) String() string        { return x.Repr }
func (x *TypeRecord) String() string       { return x.Repr }
func (x *TypeVariant) String() string      { return x.Repr }
func (x *TypeList) String() string         { return x.Repr }
func (x *TypeUnit) String() string         { return x.Repr }
func (x *TypeTop) String() string          { return x.Repr }
func (x *TypeBot) String() string          { return x.Repr }
func (x *TypeAuto) String() string         { return x.Repr }
func (x *TypeVar) String() string          { return x.Repr }
func (x *TypeParens) String() string       { return x.Repr }

package ast

import (
	"fmt"
	"strings"

	"github.com/neocotic/go-optional"
)

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

func (field RecordFieldType) String() string {
	return fmt.Sprintf("%s : %s", field.Label, field.Type_)
}

func (field VariantFieldType) String() string {
	var builder strings.Builder

	builder.WriteString(field.Label.String())

	field.Type_.IfPresent(
		func(value StellaType) {
			fmt.Fprintf(&builder, " : %s", value.String())
		},
	)

	return builder.String()
}

func (t TypeBool) String() string {
	return "Bool"
}

func (t TypeNat) String() string {
	return "Nat"
}

func (t TypeRef) String() string {
	return fmt.Sprintf("&%s", t.Type_.String())
}

func (t TypeSum) String() string {
	return fmt.Sprintf("%s+%s", t.Left.String(), t.Right.String())
}

func (t TypeFun) String() string {
	params := ListToString(t.ParamTypes, ", ")
	return fmt.Sprintf("fn(%s) -> %s", params, t.ReturnType.String())
}

func (t TypeForAll) String() string {
	params := ListToString(t.Types, " ")
	return fmt.Sprintf("forall %s.%s", params, t.Type_.String())
}

func (t TypeRec) String() string {
	return fmt.Sprintf("µ %s.%s", t.Var_.String(), t.Type_.String())
}

func (t TypeTuple) String() string {
	return fmt.Sprintf("{ %s }", ListToString(t.Types, ", "))
}

func (t TypeRecord) String() string {
	return fmt.Sprintf("{ %s }", ListToString(t.FieldTypes, ", "))
}

func (t TypeVariant) String() string {
	return fmt.Sprintf("<| %s |>", ListToString(t.FieldTypes, ", "))
}

func (t TypeList) String() string {
	return fmt.Sprintf("[ %s ]", t.Type_)
}

func (t TypeUnit) String() string {
	return "Unit"
}

func (t TypeTop) String() string {
	return "Top"
}

func (t TypeBot) String() string {
	return "Bot"
}

func (t TypeAuto) String() string {
	return "auto"
}

func (t TypeVar) String() string {
	return t.Name.String()
}

func (t TypeParens) String() string {
	return fmt.Sprintf("(%s)", t.Type_)
}

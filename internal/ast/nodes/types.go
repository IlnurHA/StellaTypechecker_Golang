package ast

import (
	"fmt"
	"strings"

	"github.com/neocotic/go-optional"
)

type RecordFieldType struct {
	Label StellaIdent
	Type_ StellaType
}

type VariantFieldType struct {
	Label StellaIdent
	Type_ optional.Optional[StellaType]
}

type StellaType interface {
	isType()
	String() string
}

type TypeBool struct{}
type TypeNat struct{}

type TypeRef struct {
	Type_ StellaType
}
type TypeSum struct {
	Left  StellaType
	Right StellaType
}
type TypeFun struct {
	ParamTypes []StellaType
	ReturnType StellaType
}
type TypeForAll struct {
	Types []StellaIdent
	Type_ StellaType
}
type TypeRec struct {
	Var_  StellaIdent
	Type_ StellaType
}
type TypeTuple struct {
	Types []StellaType
}
type TypeRecord struct {
	FieldTypes []RecordFieldType
}
type TypeVariant struct {
	FieldTypes []VariantFieldType
}
type TypeList struct {
	Type_ StellaType
}

type TypeUnit struct{}
type TypeTop struct{}
type TypeBot struct{}
type TypeAuto struct{}
type TypeVar struct {
	Name StellaIdent
}
type TypeParens struct {
	Type_ StellaType
}

func (field *RecordFieldType) String() string {
	return fmt.Sprintf("%s : %s", field.Label, field.Type_)
}

func (field *VariantFieldType) String() string {
	var builder strings.Builder

	builder.WriteString(field.Label.String())

	field.Type_.IfPresent(
		func(value StellaType) {
			fmt.Fprintf(&builder, " : %s", value.String())
		},
	)

	return builder.String()
}

func (t *TypeBool) String() string {
	return "Bool"
}

func (t *TypeNat) String() string {
	return "Nat"
}

func (t *TypeRef) String() string {
	return fmt.Sprintf("&%s", t.Type_.String())
}

func (t *TypeSum) String() string {
	return fmt.Sprintf("%s+%s", t.Left.String(), t.Right.String())
}

func (t *TypeFun) String() string {
	params := ListToString(t.ParamTypes, ", ")
	return fmt.Sprintf("fn(%s) -> %s", params, t.ReturnType.String())
}

func (t *TypeForAll) String() string {
	pointerParams := make([]*StellaIdent, len(t.Types))

	for index, type_ := range t.Types {
		pointerParams[index] = &type_
	}

	params := ListToString(pointerParams, " ")
	return fmt.Sprintf("forall %s.%s", params, t.Type_.String())
}

func (t *TypeRec) String() string {
	return fmt.Sprintf("µ %s.%s", t.Var_.String(), t.Type_.String())
}

func (t *TypeTuple) String() string {
	return fmt.Sprintf("{ %s }", ListToString(t.Types, ", "))
}

func (t *TypeRecord) String() string {
	pointerParams := make([]*RecordFieldType, len(t.FieldTypes))

	for index, type_ := range t.FieldTypes {
		pointerParams[index] = &type_
	}
	return fmt.Sprintf("{ %s }", ListToString(pointerParams, ", "))
}

func (t *TypeVariant) String() string {
	pointerParams := make([]*VariantFieldType, len(t.FieldTypes))

	for index, type_ := range t.FieldTypes {
		pointerParams[index] = &type_
	}
	return fmt.Sprintf("<| %s |>", ListToString(pointerParams, ", "))
}

func (t *TypeList) String() string {
	return fmt.Sprintf("[ %s ]", t.Type_)
}

func (t *TypeUnit) String() string {
	return "Unit"
}

func (t *TypeTop) String() string {
	return "Top"
}

func (t *TypeBot) String() string {
	return "Bot"
}

func (t *TypeAuto) String() string {
	return "auto"
}

func (t *TypeVar) String() string {
	return t.Name.String()
}

func (t *TypeParens) String() string {
	return fmt.Sprintf("(%s)", t.Type_)
}

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

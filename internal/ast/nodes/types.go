package ast

import "github.com/neocotic/go-optional"

type RecordFieldType struct {
	Label StellaIdent
	Type_ StellaType
}

type VariantFieldType struct {
	Label StellaIdent
	Type_ optional.Optional[StellaType]
}

type StellaType interface{ isType() }

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

func (x *RecordFieldType) isNode() {
}
func (x *VariantFieldType) isNode() {
}
func (x *TypeBool) isNode() {
}
func (x *TypeBool) isType() {
}
func (x *TypeNat) isNode() {
}
func (x *TypeNat) isType() {
}
func (x *TypeRef) isNode() {
}
func (x *TypeRef) isType() {
}
func (x *TypeSum) isNode() {
}
func (x *TypeSum) isType() {
}
func (x *TypeFun) isNode() {
}
func (x *TypeFun) isType() {
}
func (x *TypeForAll) isNode() {
}
func (x *TypeForAll) isType() {
}
func (x *TypeRec) isNode() {
}
func (x *TypeRec) isType() {
}
func (x *TypeTuple) isNode() {
}
func (x *TypeTuple) isType() {
}
func (x *TypeRecord) isNode() {
}
func (x *TypeRecord) isType() {
}
func (x *TypeVariant) isNode() {
}
func (x *TypeVariant) isType() {
}
func (x *TypeList) isNode() {
}
func (x *TypeList) isType() {
}
func (x *TypeUnit) isNode() {
}
func (x *TypeUnit) isType() {
}
func (x *TypeTop) isNode() {
}
func (x *TypeTop) isType() {
}
func (x *TypeBot) isNode() {
}
func (x *TypeBot) isType() {
}
func (x *TypeAuto) isNode() {
}
func (x *TypeAuto) isType() {
}
func (x *TypeVar) isNode() {
}
func (x *TypeVar) isType() {
}
func (x *TypeParens) isNode() {
}
func (x *TypeParens) isType() {
}

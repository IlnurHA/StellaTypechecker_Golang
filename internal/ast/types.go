package ast

import "github.com/neocotic/go-optional"

type RecordFieldType struct {
	label StellaIdent
	type_ StellaType
}

type VariantFieldType struct {
	label StellaIdent
	type_ optional.Optional[StellaType]
}

type StellaType interface{ isType() }

type TypeBool struct{}
type TypeNat struct{}

type TypeRef struct {
	type_ StellaType
}
type TypeSum struct {
	left  StellaType
	right StellaType
}
type TypeFun struct {
	paramTypes []StellaType
	returnType StellaType
}
type TypeForAll struct {
	types []StellaIdent
	type_ StellaType
}
type TypeRec struct {
	var_  StellaIdent
	type_ StellaType
}
type TypeTuple struct {
	types []StellaType
}
type TypeRecord struct {
	fieldTypes []RecordFieldType
}
type TypeVariant struct {
	fieldTypes []VariantFieldType
}
type TypeList struct {
	type_ StellaType
}

type TypeUnit struct{}
type TypeTop struct{}
type TypeBot struct{}
type TypeAuto struct{}
type TypeVar struct {
	name StellaIdent
}
type TypeParens struct {
	type_ StellaIdent
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

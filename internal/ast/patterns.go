package ast

import "github.com/neocotic/go-optional"

type Pattern interface{ isPattern() }

type PatternVariant struct {
	label   StellaIdent
	pattern optional.Optional[Pattern]
}

type PatternInl struct {
	pattern Pattern
}

type PatternInr struct {
	pattern Pattern
}

type PatternTuple struct {
	patterns []Pattern
}

type PatternRecord struct {
	patterns []LabelledPattern
}

type PatternList struct {
	patterns []Pattern
}

type PatternCons struct {
	head Pattern
	tail Pattern
}

type PatternFalse struct{}
type PatternTrue struct{}
type PatternUnit struct{}
type PatternInt struct {
	n int
}
type PatternSucc struct {
	pattern Pattern
}
type PatternVar struct {
	name StellaIdent
}
type PatternAsc struct {
	pattern Pattern
	type_   StellaType
}
type PatternCastAs struct {
	pattern Pattern
	type_   StellaType
}
type ParenthesisedPattern struct {
	pattern Pattern
}

type LabelledPattern struct {
	label   StellaIdent
	pattern Pattern
}

func (x *PatternVariant) isNode() {
}
func (x *PatternVariant) isPattern() {
}
func (x *PatternInl) isNode() {
}
func (x *PatternInl) isPattern() {
}
func (x *PatternInr) isNode() {
}
func (x *PatternInr) isPattern() {
}
func (x *PatternTuple) isNode() {
}
func (x *PatternTuple) isPattern() {
}
func (x *PatternRecord) isNode() {
}
func (x *PatternRecord) isPattern() {
}
func (x *PatternList) isNode() {
}
func (x *PatternList) isPattern() {
}
func (x *PatternCons) isNode() {
}
func (x *PatternCons) isPattern() {
}
func (x *PatternFalse) isNode() {
}
func (x *PatternFalse) isPattern() {
}
func (x *PatternTrue) isNode() {
}
func (x *PatternTrue) isPattern() {
}
func (x *PatternUnit) isNode() {
}
func (x *PatternUnit) isPattern() {
}
func (x *PatternInt) isNode() {
}
func (x *PatternInt) isPattern() {
}
func (x *PatternSucc) isNode() {
}
func (x *PatternSucc) isPattern() {
}
func (x *PatternVar) isNode() {
}
func (x *PatternVar) isPattern() {
}
func (x *PatternAsc) isNode() {
}
func (x *PatternAsc) isPattern() {
}
func (x *PatternCastAs) isNode() {
}
func (x *PatternCastAs) isPattern() {
}
func (x *ParenthesisedPattern) isNode() {
}
func (x *ParenthesisedPattern) isPattern() {
}
func (x *LabelledPattern) isNode() {
}

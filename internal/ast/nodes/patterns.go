package ast

import "github.com/neocotic/go-optional"

type Pattern interface{ isPattern() }

type PatternVariant struct {
	Label   StellaIdent
	Pattern optional.Optional[Pattern]
}

type PatternInl struct {
	Pattern Pattern
}

type PatternInr struct {
	Pattern Pattern
}

type PatternTuple struct {
	Patterns []Pattern
}

type PatternRecord struct {
	Patterns []LabelledPattern
}

type PatternList struct {
	Patterns []Pattern
}

type PatternCons struct {
	Head Pattern
	Tail Pattern
}

type PatternFalse struct{}
type PatternTrue struct{}
type PatternUnit struct{}
type PatternInt struct {
	N int
}
type PatternSucc struct {
	Pattern Pattern
}
type PatternVar struct {
	Name StellaIdent
}
type PatternAsc struct {
	Pattern Pattern
	Type_   StellaType
}
type PatternCastAs struct {
	Pattern Pattern
	Type_   StellaType
}
type ParenthesisedPattern struct {
	Pattern Pattern
}

type LabelledPattern struct {
	Label   StellaIdent
	Pattern Pattern
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

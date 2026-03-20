package ast

import "github.com/neocotic/go-optional"

type Pattern interface {
	Node
	isPattern()
}

type PatternVariant struct {
	Label   StellaIdent
	Pattern optional.Optional[Pattern]

	Repr string
}

type PatternInl struct {
	Pattern Pattern

	Repr string
}

type PatternInr struct {
	Pattern Pattern

	Repr string
}

type PatternTuple struct {
	Patterns []Pattern

	Repr string
}

type PatternRecord struct {
	Patterns []LabelledPattern

	Repr string
}

type PatternList struct {
	Patterns []Pattern

	Repr string
}

type PatternCons struct {
	Head Pattern
	Tail Pattern

	Repr string
}

type PatternFalse struct {
	Repr string
}
type PatternTrue struct {
	Repr string
}
type PatternUnit struct {
	Repr string
}
type PatternInt struct {
	N int

	Repr string
}
type PatternSucc struct {
	Pattern Pattern

	Repr string
}
type PatternVar struct {
	Name StellaIdent

	Repr string
}
type PatternAsc struct {
	Pattern Pattern
	Type_   StellaType

	Repr string
}
type PatternCastAs struct {
	Pattern Pattern
	Type_   StellaType

	Repr string
}
type ParenthesisedPattern struct {
	Pattern Pattern

	Repr string
}

type LabelledPattern struct {
	Label   StellaIdent
	Pattern Pattern

	Repr string
}

// Marker methods
func (x *PatternVariant) isNode()          {}
func (x *PatternVariant) isPattern()       {}
func (x *PatternInl) isNode()              {}
func (x *PatternInl) isPattern()           {}
func (x *PatternInr) isNode()              {}
func (x *PatternInr) isPattern()           {}
func (x *PatternTuple) isNode()            {}
func (x *PatternTuple) isPattern()         {}
func (x *PatternRecord) isNode()           {}
func (x *PatternRecord) isPattern()        {}
func (x *PatternList) isNode()             {}
func (x *PatternList) isPattern()          {}
func (x *PatternCons) isNode()             {}
func (x *PatternCons) isPattern()          {}
func (x *PatternFalse) isNode()            {}
func (x *PatternFalse) isPattern()         {}
func (x *PatternTrue) isNode()             {}
func (x *PatternTrue) isPattern()          {}
func (x *PatternUnit) isNode()             {}
func (x *PatternUnit) isPattern()          {}
func (x *PatternInt) isNode()              {}
func (x *PatternInt) isPattern()           {}
func (x *PatternSucc) isNode()             {}
func (x *PatternSucc) isPattern()          {}
func (x *PatternVar) isNode()              {}
func (x *PatternVar) isPattern()           {}
func (x *PatternAsc) isNode()              {}
func (x *PatternAsc) isPattern()           {}
func (x *PatternCastAs) isNode()           {}
func (x *PatternCastAs) isPattern()        {}
func (x *ParenthesisedPattern) isNode()    {}
func (x *ParenthesisedPattern) isPattern() {}
func (x *LabelledPattern) isNode()         {}

// String methods
func (x *PatternVariant) String() string       { return x.Repr }
func (x *PatternInl) String() string           { return x.Repr }
func (x *PatternInr) String() string           { return x.Repr }
func (x *PatternTuple) String() string         { return x.Repr }
func (x *PatternRecord) String() string        { return x.Repr }
func (x *PatternList) String() string          { return x.Repr }
func (x *PatternCons) String() string          { return x.Repr }
func (x *PatternFalse) String() string         { return x.Repr }
func (x *PatternTrue) String() string          { return x.Repr }
func (x *PatternUnit) String() string          { return x.Repr }
func (x *PatternInt) String() string           { return x.Repr }
func (x *PatternSucc) String() string          { return x.Repr }
func (x *PatternVar) String() string           { return x.Repr }
func (x *PatternAsc) String() string           { return x.Repr }
func (x *PatternCastAs) String() string        { return x.Repr }
func (x *ParenthesisedPattern) String() string { return x.Repr }
func (x *LabelledPattern) String() string      { return x.Repr }

package ast

import "github.com/neocotic/go-optional"

type PatternBinding struct {
	Pattern Pattern
	Rhs     Expr

	Repr string
}

type Binding struct {
	Name StellaIdent
	Rhs  Expr

	Repr string
}

type MatchCase struct {
	Pattern Pattern
	Expr_   Expr

	Repr string
}

type Expr interface {
	isExpr()
	String() string
}

type DotRecord struct {
	Subexpr Expr
	Label   StellaIdent

	Repr string
}

type DotTuple struct {
	Subexpr Expr
	Index   int

	Repr string
}

type ConstBool struct {
	Value bool

	Repr string
}

type ConstUnit struct{ Repr string }

type ConstInt struct {
	Value int

	Repr string
}

type ConstMemory struct {
	Memory MemoryAddress

	Repr string
}

type Var struct {
	Name StellaIdent

	Repr string
}

type Panic struct {
	Repr string
}
type Throw struct {
	Expr_ Expr

	Repr string
}
type TryCatch struct {
	TryExpr      Expr
	Pattern      Pattern
	FallbackExpr Expr

	Repr string
}
type TryCastAs struct {
	TryExpr      Expr
	Type_        StellaType
	Pattern      Pattern
	Expr_        Expr
	FallbackExpr Expr

	Repr string
}
type TryWith struct {
	TryExpr      Expr
	FallbackExpr Expr

	Repr string
}

type Inl struct {
	Expr_ Expr

	Repr string
}
type Inr struct {
	Expr_ Expr
	Repr  string
}

type ConsList struct {
	Head Expr
	Tail Expr

	Repr string
}
type Head struct {
	List Expr

	Repr string
}
type IsEmpty struct {
	List Expr

	Repr string
}
type Tail struct {
	List Expr

	Repr string
}

type Succ struct {
	N Expr

	Repr string
}

type LogicNot struct {
	Expr_ Expr

	Repr string
}

type Pred struct {
	N Expr

	Repr string
}

type IsZero struct {
	N Expr

	Repr string
}

type Fix struct {
	Expr_ Expr

	Repr string
}

type NatRec struct {
	N       Expr
	Initial Expr
	Step    Expr

	Repr string
}

type Fold struct {
	Type_ StellaType
	Expr_ Expr

	Repr string
}

type Unfold struct {
	Type_ StellaType
	Expr_ Expr

	Repr string
}

type Application struct {
	Function Expr
	Args     []Expr

	Repr string
}
type TypeApplication struct {
	Function Expr
	Types    []StellaType

	Repr string
}

type Multiply struct {
	Left  Expr
	Right Expr

	Repr string
}

type Divide struct {
	Left  Expr
	Right Expr

	Repr string
}

type LogicAnd struct {
	Left  Expr
	Right Expr

	Repr string
}

type Ref struct {
	Expr_ Expr

	Repr string
}

type Deref struct {
	Expr_ Expr

	Repr string
}

type Add struct {
	Left  Expr
	Right Expr

	Repr string
}
type Subtract struct {
	Left  Expr
	Right Expr

	Repr string
}
type LogicOr struct {
	Left  Expr
	Right Expr

	Repr string
}

type TypeAsc struct {
	Expr_ Expr
	Type_ StellaType

	Repr string
}

type TypeCast struct {
	Expr_ Expr
	Type_ StellaType

	Repr string
}

type Abstraction struct {
	Params     []ParameterDeclaration
	ReturnExpr Expr

	Repr string
}

type Tuple struct {
	Exprs []Expr

	Repr string
}

type Record struct {
	Bindings []Binding

	Repr string
}

type Variant struct {
	Label StellaIdent
	Rhs   optional.Optional[Expr]

	Repr string
}

type Match struct {
	Expr_ Expr
	Cases []MatchCase

	Repr string
}

type List struct {
	Exprs []Expr

	Repr string
}

type LessThan struct {
	Left  Expr
	Right Expr

	Repr string
}
type LessThanOrEqual struct {
	Left  Expr
	Right Expr

	Repr string
}
type GreaterThan struct {
	Left  Expr
	Right Expr

	Repr string
}
type GreaterThanOrEqual struct {
	Left  Expr
	Right Expr

	Repr string
}
type Equal struct {
	Left  Expr
	Right Expr

	Repr string
}
type NotEqual struct {
	Left  Expr
	Right Expr

	Repr string
}

type Assign struct {
	Lhs Expr
	Rhs Expr

	Repr string
}

type If struct {
	Condition Expr
	ThenExpr  Expr
	ElseExpr  Expr

	Repr string
}

type Sequence struct {
	Expr1 Expr
	Expr2 Expr

	Repr string
}

type Let struct {
	PatternBindings []PatternBinding
	Body            Expr

	Repr string
}

type LetRec struct {
	PatternBindings []PatternBinding
	Body            Expr

	Repr string
}

type TypeAbstraction struct {
	Generics []StellaIdent
	Expr_    Expr

	Repr string
}

type ParenthesisedExpr struct {
	Expr_ Expr

	Repr string
}

type TerminatingSemicolon struct {
	Expr_ Expr

	Repr string
}

// Marker methods
func (x *PatternBinding) isNode()       {}
func (x *Binding) isNode()              {}
func (x *MatchCase) isNode()            {}
func (x *DotRecord) isNode()            {}
func (x *DotRecord) isExpr()            {}
func (x *DotTuple) isNode()             {}
func (x *DotTuple) isExpr()             {}
func (x *ConstBool) isNode()            {}
func (x *ConstBool) isExpr()            {}
func (x *ConstUnit) isNode()            {}
func (x *ConstUnit) isExpr()            {}
func (x *ConstInt) isNode()             {}
func (x *ConstInt) isExpr()             {}
func (x *ConstMemory) isNode()          {}
func (x *ConstMemory) isExpr()          {}
func (x *Var) isNode()                  {}
func (x *Var) isExpr()                  {}
func (x *Panic) isNode()                {}
func (x *Panic) isExpr()                {}
func (x *Throw) isNode()                {}
func (x *Throw) isExpr()                {}
func (x *TryCatch) isNode()             {}
func (x *TryCatch) isExpr()             {}
func (x *TryCastAs) isNode()            {}
func (x *TryCastAs) isExpr()            {}
func (x *TryWith) isNode()              {}
func (x *TryWith) isExpr()              {}
func (x *Inl) isNode()                  {}
func (x *Inl) isExpr()                  {}
func (x *Inr) isNode()                  {}
func (x *Inr) isExpr()                  {}
func (x *ConsList) isNode()             {}
func (x *ConsList) isExpr()             {}
func (x *Head) isNode()                 {}
func (x *Head) isExpr()                 {}
func (x *IsEmpty) isNode()              {}
func (x *IsEmpty) isExpr()              {}
func (x *Tail) isNode()                 {}
func (x *Tail) isExpr()                 {}
func (x *Succ) isNode()                 {}
func (x *Succ) isExpr()                 {}
func (x *LogicNot) isNode()             {}
func (x *LogicNot) isExpr()             {}
func (x *Pred) isNode()                 {}
func (x *Pred) isExpr()                 {}
func (x *IsZero) isNode()               {}
func (x *IsZero) isExpr()               {}
func (x *Fix) isNode()                  {}
func (x *Fix) isExpr()                  {}
func (x *NatRec) isNode()               {}
func (x *NatRec) isExpr()               {}
func (x *Fold) isNode()                 {}
func (x *Fold) isExpr()                 {}
func (x *Unfold) isNode()               {}
func (x *Unfold) isExpr()               {}
func (x *Application) isNode()          {}
func (x *Application) isExpr()          {}
func (x *TypeApplication) isNode()      {}
func (x *TypeApplication) isExpr()      {}
func (x *Multiply) isNode()             {}
func (x *Multiply) isExpr()             {}
func (x *Divide) isNode()               {}
func (x *Divide) isExpr()               {}
func (x *LogicAnd) isNode()             {}
func (x *LogicAnd) isExpr()             {}
func (x *Ref) isNode()                  {}
func (x *Ref) isExpr()                  {}
func (x *Add) isNode()                  {}
func (x *Add) isExpr()                  {}
func (x *Subtract) isNode()             {}
func (x *Subtract) isExpr()             {}
func (x *LogicOr) isNode()              {}
func (x *LogicOr) isExpr()              {}
func (x *TypeAsc) isNode()              {}
func (x *TypeAsc) isExpr()              {}
func (x *TypeCast) isNode()             {}
func (x *TypeCast) isExpr()             {}
func (x *Abstraction) isNode()          {}
func (x *Abstraction) isExpr()          {}
func (x *Tuple) isNode()                {}
func (x *Tuple) isExpr()                {}
func (x *Record) isNode()               {}
func (x *Record) isExpr()               {}
func (x *Variant) isNode()              {}
func (x *Variant) isExpr()              {}
func (x *Match) isNode()                {}
func (x *Match) isExpr()                {}
func (x *List) isNode()                 {}
func (x *List) isExpr()                 {}
func (x *LessThan) isNode()             {}
func (x *LessThan) isExpr()             {}
func (x *LessThanOrEqual) isNode()      {}
func (x *LessThanOrEqual) isExpr()      {}
func (x *GreaterThan) isNode()          {}
func (x *GreaterThan) isExpr()          {}
func (x *GreaterThanOrEqual) isNode()   {}
func (x *GreaterThanOrEqual) isExpr()   {}
func (x *Equal) isNode()                {}
func (x *Equal) isExpr()                {}
func (x *NotEqual) isNode()             {}
func (x *NotEqual) isExpr()             {}
func (x *Assign) isNode()               {}
func (x *Assign) isExpr()               {}
func (x *If) isNode()                   {}
func (x *If) isExpr()                   {}
func (x *Sequence) isNode()             {}
func (x *Sequence) isExpr()             {}
func (x *Let) isNode()                  {}
func (x *Let) isExpr()                  {}
func (x *LetRec) isNode()               {}
func (x *LetRec) isExpr()               {}
func (x *TypeAbstraction) isNode()      {}
func (x *TypeAbstraction) isExpr()      {}
func (x *ParenthesisedExpr) isNode()    {}
func (x *ParenthesisedExpr) isExpr()    {}
func (x *TerminatingSemicolon) isNode() {}
func (x *TerminatingSemicolon) isExpr() {}

// String methods
func (x *PatternBinding) String() string       { return x.Repr }
func (x *Binding) String() string              { return x.Repr }
func (x *MatchCase) String() string            { return x.Repr }
func (x *DotRecord) String() string            { return x.Repr }
func (x *DotTuple) String() string             { return x.Repr }
func (x *ConstBool) String() string            { return x.Repr }
func (x *ConstUnit) String() string            { return x.Repr }
func (x *ConstInt) String() string             { return x.Repr }
func (x *ConstMemory) String() string          { return x.Repr }
func (x *Var) String() string                  { return x.Repr }
func (x *Panic) String() string                { return x.Repr }
func (x *Throw) String() string                { return x.Repr }
func (x *TryCatch) String() string             { return x.Repr }
func (x *TryCastAs) String() string            { return x.Repr }
func (x *TryWith) String() string              { return x.Repr }
func (x *Inl) String() string                  { return x.Repr }
func (x *Inr) String() string                  { return x.Repr }
func (x *ConsList) String() string             { return x.Repr }
func (x *Head) String() string                 { return x.Repr }
func (x *IsEmpty) String() string              { return x.Repr }
func (x *Tail) String() string                 { return x.Repr }
func (x *Succ) String() string                 { return x.Repr }
func (x *LogicNot) String() string             { return x.Repr }
func (x *Pred) String() string                 { return x.Repr }
func (x *IsZero) String() string               { return x.Repr }
func (x *Fix) String() string                  { return x.Repr }
func (x *NatRec) String() string               { return x.Repr }
func (x *Fold) String() string                 { return x.Repr }
func (x *Unfold) String() string               { return x.Repr }
func (x *Application) String() string          { return x.Repr }
func (x *TypeApplication) String() string      { return x.Repr }
func (x *Multiply) String() string             { return x.Repr }
func (x *Divide) String() string               { return x.Repr }
func (x *LogicAnd) String() string             { return x.Repr }
func (x *Ref) String() string                  { return x.Repr }
func (x *Deref) String() string                { return x.Repr }
func (x *Add) String() string                  { return x.Repr }
func (x *Subtract) String() string             { return x.Repr }
func (x *LogicOr) String() string              { return x.Repr }
func (x *TypeAsc) String() string              { return x.Repr }
func (x *TypeCast) String() string             { return x.Repr }
func (x *Abstraction) String() string          { return x.Repr }
func (x *Tuple) String() string                { return x.Repr }
func (x *Record) String() string               { return x.Repr }
func (x *Variant) String() string              { return x.Repr }
func (x *Match) String() string                { return x.Repr }
func (x *List) String() string                 { return x.Repr }
func (x *LessThan) String() string             { return x.Repr }
func (x *LessThanOrEqual) String() string      { return x.Repr }
func (x *GreaterThan) String() string          { return x.Repr }
func (x *GreaterThanOrEqual) String() string   { return x.Repr }
func (x *Equal) String() string                { return x.Repr }
func (x *NotEqual) String() string             { return x.Repr }
func (x *Assign) String() string               { return x.Repr }
func (x *If) String() string                   { return x.Repr }
func (x *Sequence) String() string             { return x.Repr }
func (x *Let) String() string                  { return x.Repr }
func (x *LetRec) String() string               { return x.Repr }
func (x *TypeAbstraction) String() string      { return x.Repr }
func (x *ParenthesisedExpr) String() string    { return x.Repr }
func (x *TerminatingSemicolon) String() string { return x.Repr }

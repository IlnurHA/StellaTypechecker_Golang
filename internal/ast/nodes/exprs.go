package ast

import "github.com/neocotic/go-optional"

type PatternBinding struct {
	Pattern Pattern
	Rhs     Expr
}

type Binding struct {
	Name StellaIdent
	Rhs  Expr
}

type MatchCase struct {
	Pattern Pattern
	Expr_   Expr
}

type Expr interface{ isExpr() }

type DotRecord struct {
	Subexpr Expr
	Label   StellaIdent
}

type DotTuple struct {
	Subexpr Expr
	Index   int
}

type ConstBool struct {
	Value bool
}

type ConstUnit struct{}

type ConstInt struct {
	Value int
}

type ConstMemory struct {
	Memory MemoryAddress
}

type Var struct {
	Name StellaIdent
}

type Panic struct{}
type Throw struct{ Expr_ Expr }
type TryCatch struct {
	TryExpr      Expr
	Pattern      Pattern
	FallbackExpr Expr
}
type TryCastAs struct {
	TryExpr      Expr
	Type_        StellaType
	Pattern      Pattern
	Expr_        Expr
	FallbackExpr Expr
}
type TryWith struct {
	TryExpr      Expr
	FallbackExpr Expr
}

type Inl struct{ Expr_ Expr }
type Inr struct{ Expr_ Expr }

type ConsList struct {
	Head Expr
	Tail Expr
}
type Head struct {
	List Expr
}
type IsEmpty struct {
	List Expr
}
type Tail struct {
	List Expr
}

type Succ struct {
	N Expr
}

type LogicNot struct {
	Expr_ Expr
}

type Pred struct {
	N Expr
}

type IsZero struct {
	N Expr
}

type Fix struct {
	Expr_ Expr
}

type NatRec struct {
	N       Expr
	Initial Expr
	Step    Expr
}

type Fold struct {
	Type_ StellaType
	Expr_ Expr
}

type Unfold struct {
	Type_ StellaType
	Expr_ Expr
}

type Application struct {
	Function Expr
	Args     []Expr
}
type TypeApplication struct {
	Function Expr
	Types    []StellaType
}

type Multiply struct {
	Left  Expr
	Right Expr
}

type Divide struct {
	Left  Expr
	Right Expr
}

type LogicAnd struct {
	Left  Expr
	Right Expr
}

type Ref struct {
	Expr_ Expr
}

type Deref struct {
	Expr_ Expr
}

type Add struct {
	Left  Expr
	Right Expr
}
type Subtract struct {
	Left  Expr
	Right Expr
}
type LogicOr struct {
	Left  Expr
	Right Expr
}

type TypeAsc struct {
	Expr_ Expr
	Type_ StellaType
}

type TypeCast struct {
	Expr_ Expr
	Type_ StellaType
}

type Abstraction struct {
	Params     []ParameterDeclaration
	ReturnExpr Expr
}

type Tuple struct {
	Exprs []Expr
}

type Record struct {
	Bindings []Binding
}

type Variant struct {
	Label StellaIdent
	Rhs   optional.Optional[Expr]
}

type Match struct {
	Expr_ Expr
	Cases []MatchCase
}

type List struct {
	Exprs []Expr
}

type LessThan struct {
	Left  Expr
	Right Expr
}
type LessThanOrEqual struct {
	Left  Expr
	Right Expr
}
type GreaterThan struct {
	Left  Expr
	Right Expr
}
type GreaterThanOrEqual struct {
	Left  Expr
	Right Expr
}
type Equal struct {
	Left  Expr
	Right Expr
}
type NotEqual struct {
	Left  Expr
	Right Expr
}

type Assign struct {
	Lhs Expr
	Rhs Expr
}

type If struct {
	Condition Expr
	ThenExpr  Expr
	ElseExpr  Expr
}

type Sequence struct {
	Expr1 Expr
	Expr2 Expr
}

type Let struct {
	PatternBindings []PatternBinding
	Body            Expr
}

type LetRec struct {
	PatternBindings []PatternBinding
	Body            Expr
}

type TypeAbstraction struct {
	Generics []StellaIdent
	Expr_    Expr
}

type ParenthesisedExpr struct {
	Expr_ Expr
}

type TerminatingSemicolon struct {
	Expr_ Expr
}

func (x *PatternBinding) isNode() {
}
func (x *Binding) isNode() {
}
func (x *MatchCase) isNode() {
}
func (x *DotRecord) isNode() {
}
func (x *DotRecord) isExpr() {
}
func (x *DotTuple) isNode() {
}
func (x *DotTuple) isExpr() {
}
func (x *ConstBool) isNode() {
}
func (x *ConstBool) isExpr() {
}
func (x *ConstUnit) isNode() {
}
func (x *ConstUnit) isExpr() {
}
func (x *ConstInt) isNode() {
}
func (x *ConstInt) isExpr() {
}
func (x *ConstMemory) isNode() {
}
func (x *ConstMemory) isExpr() {
}
func (x *Var) isNode() {
}
func (x *Var) isExpr() {
}
func (x *Panic) isNode() {
}
func (x *Panic) isExpr() {
}
func (x *Throw) isNode() {
}
func (x *Throw) isExpr() {
}
func (x *TryCatch) isNode() {
}
func (x *TryCatch) isExpr() {
}
func (x *TryCastAs) isNode() {
}
func (x *TryCastAs) isExpr() {
}
func (x *TryWith) isNode() {
}
func (x *TryWith) isExpr() {
}
func (x *Inl) isNode() {
}
func (x *Inl) isExpr() {
}
func (x *Inr) isNode() {
}
func (x *Inr) isExpr() {
}
func (x *ConsList) isNode() {
}
func (x *ConsList) isExpr() {
}
func (x *Head) isNode() {
}
func (x *Head) isExpr() {
}
func (x *IsEmpty) isNode() {
}
func (x *IsEmpty) isExpr() {
}
func (x *Tail) isNode() {
}
func (x *Tail) isExpr() {
}
func (x *Succ) isNode() {
}
func (x *Succ) isExpr() {
}
func (x *LogicNot) isNode() {
}
func (x *LogicNot) isExpr() {
}
func (x *Pred) isNode() {
}
func (x *Pred) isExpr() {
}
func (x *IsZero) isNode() {
}
func (x *IsZero) isExpr() {
}
func (x *Fix) isNode() {
}
func (x *Fix) isExpr() {
}
func (x *NatRec) isNode() {
}
func (x *NatRec) isExpr() {
}
func (x *Fold) isNode() {
}
func (x *Fold) isExpr() {
}
func (x *Unfold) isNode() {
}
func (x *Unfold) isExpr() {
}
func (x *Application) isNode() {
}
func (x *Application) isExpr() {
}
func (x *TypeApplication) isNode() {
}
func (x *TypeApplication) isExpr() {
}
func (x *Multiply) isNode() {
}
func (x *Multiply) isExpr() {
}
func (x *Divide) isNode() {
}
func (x *Divide) isExpr() {
}
func (x *LogicAnd) isNode() {
}
func (x *LogicAnd) isExpr() {
}
func (x *Ref) isNode() {
}
func (x *Ref) isExpr() {
}
func (x *Add) isNode() {
}
func (x *Add) isExpr() {
}
func (x *Subtract) isNode() {
}
func (x *Subtract) isExpr() {
}
func (x *LogicOr) isNode() {
}
func (x *LogicOr) isExpr() {
}
func (x *TypeAsc) isNode() {
}
func (x *TypeAsc) isExpr() {
}
func (x *TypeCast) isNode() {
}
func (x *TypeCast) isExpr() {
}
func (x *Abstraction) isNode() {
}
func (x *Abstraction) isExpr() {
}
func (x *Tuple) isNode() {
}
func (x *Tuple) isExpr() {
}
func (x *Record) isNode() {
}
func (x *Record) isExpr() {
}
func (x *Variant) isNode() {
}
func (x *Variant) isExpr() {
}
func (x *Match) isNode() {
}
func (x *Match) isExpr() {
}
func (x *List) isNode() {
}
func (x *List) isExpr() {
}
func (x *LessThan) isNode() {
}
func (x *LessThan) isExpr() {
}
func (x *LessThanOrEqual) isNode() {
}
func (x *LessThanOrEqual) isExpr() {
}
func (x *GreaterThan) isNode() {
}
func (x *GreaterThan) isExpr() {
}
func (x *GreaterThanOrEqual) isNode() {
}
func (x *GreaterThanOrEqual) isExpr() {
}
func (x *Equal) isNode() {
}
func (x *Equal) isExpr() {
}
func (x *NotEqual) isNode() {
}
func (x *NotEqual) isExpr() {
}
func (x *Assign) isNode() {
}
func (x *Assign) isExpr() {
}
func (x *If) isNode() {
}
func (x *If) isExpr() {
}
func (x *Sequence) isNode() {
}
func (x *Sequence) isExpr() {
}
func (x *Let) isNode() {
}
func (x *Let) isExpr() {
}
func (x *LetRec) isNode() {
}
func (x *LetRec) isExpr() {
}
func (x *TypeAbstraction) isNode() {
}
func (x *TypeAbstraction) isExpr() {
}
func (x *ParenthesisedExpr) isNode() {
}
func (x *ParenthesisedExpr) isExpr() {
}
func (x *TerminatingSemicolon) isNode() {
}
func (x *TerminatingSemicolon) isExpr() {
}

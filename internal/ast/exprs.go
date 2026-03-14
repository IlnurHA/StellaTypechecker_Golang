package ast

import "github.com/neocotic/go-optional"

type PatternBinding struct {
	pattern Pattern
	rhs     Expr
}

type Binding struct {
	name StellaIdent
	rhs  Expr
}

type MatchCase struct {
	pattern Pattern
	expr    Expr
}

type Expr interface{ isExpr() }

type DotRecord struct {
	subexpr Expr
	label   StellaIdent
}

type DotTuple struct {
	subexpr Expr
	index   int
}

type ConstBool struct {
	value bool
}

type ConstUnit struct{}

type ConstInt struct {
	value int
}

type ConstMemory struct {
	memory MemoryAddress
}

type Var struct {
	name StellaIdent
}

type Panic struct{}
type Throw struct{ expr Expr }
type TryCatch struct {
	tryExpr      Expr
	pattern      Pattern
	fallbackExpr Expr
}
type TryCastAs struct {
	tryExpr      Expr
	type_        StellaType
	pattern      Pattern
	fallbackExpr Expr
}
type TryWith struct {
	tryExpr      Expr
	fallbackExpr Expr
}

type Inl struct{ expr Expr }
type Inr struct{ expr Expr }

type ConsList struct {
	head Expr
	tail Expr
}
type Head struct {
	list Expr
}
type IsEmpty struct {
	list Expr
}
type Tail struct {
	list Expr
}

type Succ struct {
	n Expr
}

type LogicNot struct {
	expr Expr
}

type Pred struct {
	n Expr
}

type IsZero struct {
	n Expr
}

type Fix struct {
	expr Expr
}

type NatRec struct {
	n       Expr
	initial Expr
	step    Expr
}

type Fold struct {
	type_ StellaType
	expr  Expr
}

type Unfold struct {
	type_ StellaType
	expr  Expr
}

type Application struct {
	expr Expr
	args []Expr
}
type TypeApplication struct {
	expr  Expr
	types []StellaType
}

type Multiply struct {
	left  Expr
	right Expr
}

type Divide struct {
	left  Expr
	right Expr
}

type LogicAnd struct {
	left  Expr
	right Expr
}

type Ref struct {
	expr Expr
}

type Add struct {
	left  Expr
	right Expr
}
type Subtract struct {
	left  Expr
	right Expr
}
type LogicOr struct {
	left  Expr
	right Expr
}

type TypeAsc struct {
	expr  Expr
	type_ StellaType
}

type TypeCast struct {
	expr  Expr
	type_ StellaType
}

type Abstraction struct {
	params     []ParameterDeclaration
	returnExpr Expr
}

type Tuple struct {
	exprs []Expr
}

type Record struct {
	bindings []Binding
}

type Variant struct {
	label StellaIdent
	rhs   optional.Optional[Expr]
}

type Match struct {
	expr  Expr
	cases []MatchCase
}

type List struct {
	exprs []Expr
}

type LessThan struct {
	left  Expr
	right Expr
}
type LessThanOrEqual struct {
	left  Expr
	right Expr
}
type GreaterThan struct {
	left  Expr
	right Expr
}
type GreaterThanOrEqual struct {
	left  Expr
	right Expr
}
type Equal struct {
	left  Expr
	right Expr
}
type NotEqual struct {
	left  Expr
	right Expr
}

type Assign struct {
	lhs Expr
	rhs Expr
}

type If struct {
	condition Expr
	thenExpr  Expr
	elseExpr  Expr
}

type Sequence struct {
	expr1 Expr
	expr2 Expr
}

type Let struct {
	patternBindings []PatternBinding
	body            Expr
}

type LetRec struct {
	patternBindings []PatternBinding
	body            Expr
}

type TypeAbstraction struct {
	generics []StellaType
	expr     Expr
}

type ParenthesisedExpr struct {
	expr Expr
}

type TerminatingSemicolon struct {
	expr Expr
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

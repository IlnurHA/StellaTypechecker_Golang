package ast

import (
	"strconv"
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/parser"

	"github.com/neocotic/go-optional"
)

func (v *ASTBuilder) VisitDotRecord(ctx *parser.DotRecordContext) interface{} {
	subexpr := parseExpr(ctx.GetExpr_(), v)
	label := parseStellaIdent(ctx.GetLabel())
	return nodes.DotRecord{Subexpr: subexpr, Label: label}
}

func (v *ASTBuilder) VisitDotTuple(ctx *parser.DotTupleContext) interface{} {
	subexpr := parseExpr(ctx.GetExpr_(), v)
	index, _ := strconv.Atoi(ctx.GetIndex().GetText())
	return nodes.DotTuple{Subexpr: subexpr, Index: index}
}

func (v *ASTBuilder) VisitConstTrue(ctx *parser.ConstTrueContext) interface{} {
	return nodes.ConstBool{Value: true}
}

func (v *ASTBuilder) VisitConstFalse(ctx *parser.ConstFalseContext) interface{} {
	return nodes.ConstBool{Value: false}
}

func (v *ASTBuilder) VisitConstUnit(ctx *parser.ConstUnitContext) interface{} {
	return nodes.ConstUnit{}
}

func (v *ASTBuilder) VisitConstInt(ctx *parser.ConstIntContext) interface{} {
	value, _ := strconv.Atoi(ctx.GetText())
	return nodes.ConstInt{Value: value}
}

func (v *ASTBuilder) VisitConstMemory(ctx *parser.ConstMemoryContext) interface{} {
	addr := nodes.MemoryAddress{Addr: ctx.GetMem().GetText()}
	return nodes.ConstMemory{Memory: addr}
}

func (v *ASTBuilder) VisitVar(ctx *parser.VarContext) interface{} {
	var_ := parseStellaIdent(ctx.GetName())
	return nodes.Var{Name: var_}
}

func (v *ASTBuilder) VisitPanic(ctx *parser.PanicContext) interface{} {
	return nodes.Panic{}
}

func (v *ASTBuilder) VisitThrow(ctx *parser.ThrowContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.Throw{Expr_: expr}
}

func (v *ASTBuilder) VisitTryCatch(ctx *parser.TryCatchContext) interface{} {
	tryExpr := parseExpr(ctx.GetTryExpr(), v)
	pattern := parsePattern(ctx.GetPat(), v)
	fallbackExpr := parseExpr(ctx.GetFallbackExpr(), v)

	return nodes.TryCatch{TryExpr: tryExpr, Pattern: pattern, FallbackExpr: fallbackExpr}
}

func (v *ASTBuilder) VisitTryCastAs(ctx *parser.TryCastAsContext) interface{} {
	tryExpr := parseExpr(ctx.GetTryExpr(), v)
	type_ := parseType(ctx.GetType_(), v)
	pattern := parsePattern(ctx.GetPattern_(), v)
	expr := parseExpr(ctx.GetExpr_(), v)
	fallbackExpr := parseExpr(ctx.GetFallbackExpr(), v)

	return nodes.TryCastAs{TryExpr: tryExpr, Type_: type_, Pattern: pattern, Expr_: expr, FallbackExpr: fallbackExpr}
}

func (v *ASTBuilder) VisitTryWith(ctx *parser.TryWithContext) interface{} {
	tryExpr := parseExpr(ctx.GetTryExpr(), v)
	fallbackExpr := parseExpr(ctx.GetFallbackExpr(), v)

	return nodes.TryWith{TryExpr: tryExpr, FallbackExpr: fallbackExpr}
}

func (v *ASTBuilder) VisitInl(ctx *parser.InlContext) interface{} {
	subexpr := parseExpr(ctx.GetExpr_(), v)
	return nodes.Inl{Expr_: subexpr}
}

func (v *ASTBuilder) VisitInr(ctx *parser.InrContext) interface{} {
	subexpr := parseExpr(ctx.GetExpr_(), v)
	return nodes.Inr{Expr_: subexpr}
}

func (v *ASTBuilder) VisitConsList(ctx *parser.ConsListContext) interface{} {
	head := parseExpr(ctx.GetHead(), v)
	tail := parseExpr(ctx.GetTail(), v)
	return nodes.ConsList{Head: head, Tail: tail}
}

func (v *ASTBuilder) VisitHead(ctx *parser.HeadContext) interface{} {
	list := parseExpr(ctx.GetList(), v)
	return nodes.Head{List: list}
}

func (v *ASTBuilder) VisitIsEmpty(ctx *parser.IsEmptyContext) interface{} {
	list := parseExpr(ctx.GetList(), v)
	return nodes.IsEmpty{List: list}
}

func (v *ASTBuilder) VisitTail(ctx *parser.TailContext) interface{} {
	list := parseExpr(ctx.GetList(), v)
	return nodes.Tail{List: list}
}

func (v *ASTBuilder) VisitSucc(ctx *parser.SuccContext) interface{} {
	n := parseExpr(ctx.GetN(), v)
	return nodes.Succ{N: n}
}

func (v *ASTBuilder) VisitLogicNot(ctx *parser.LogicNotContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.LogicNot{Expr_: expr}
}

func (v *ASTBuilder) VisitPred(ctx *parser.PredContext) interface{} {
	n := parseExpr(ctx.GetN(), v)
	return nodes.Pred{N: n}
}

func (v *ASTBuilder) VisitIsZero(ctx *parser.IsZeroContext) interface{} {
	n := parseExpr(ctx.GetN(), v)
	return nodes.IsZero{N: n}
}

func (v *ASTBuilder) VisitFix(ctx *parser.FixContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.Fix{Expr_: expr}
}

func (v *ASTBuilder) VisitNatRec(ctx *parser.NatRecContext) interface{} {
	n := parseExpr(ctx.GetN(), v)
	initial := v.Visit(ctx.GetInitial()).(nodes.Expr)
	step := v.Visit(ctx.GetStep()).(nodes.Expr)

	return nodes.NatRec{N: n, Initial: initial, Step: step}
}

func (v *ASTBuilder) VisitFold(ctx *parser.FoldContext) interface{} {
	type_ := parseType(ctx.GetType_(), v)
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.Fold{Type_: type_, Expr_: expr}
}

func (v *ASTBuilder) VisitUnfold(ctx *parser.UnfoldContext) interface{} {
	type_ := parseType(ctx.GetType_(), v)
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.Unfold{Type_: type_, Expr_: expr}
}

func (v *ASTBuilder) VisitApplication(ctx *parser.ApplicationContext) interface{} {
	fun := v.Visit(ctx.GetFun()).(nodes.Expr)
	args := make([]nodes.Expr, len(ctx.GetArgs()))

	for index, argCtx := range ctx.GetArgs() {
		args[index] = v.Visit(argCtx).(nodes.Expr)
	}

	return nodes.Application{Function: fun, Args: args}
}

func (v *ASTBuilder) VisitTypeApplication(ctx *parser.TypeApplicationContext) interface{} {
	fun := v.Visit(ctx.GetFun()).(nodes.Expr)
	types := parseListOfTypes(ctx.GetTypes(), v)

	return nodes.TypeApplication{Function: fun, Types: types}
}

func (v *ASTBuilder) VisitMultiply(ctx *parser.MultiplyContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.Multiply{Left: left, Right: right}
}

func (v *ASTBuilder) VisitDivide(ctx *parser.DivideContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.Divide{Left: left, Right: right}
}

func (v *ASTBuilder) VisitLogicAnd(ctx *parser.LogicAndContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.LogicAnd{Left: left, Right: right}
}

func (v *ASTBuilder) VisitRef(ctx *parser.RefContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.Ref{Expr_: expr}
}

func (v *ASTBuilder) VisitDeref(ctx *parser.DerefContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.Deref{Expr_: expr}
}

func (v *ASTBuilder) VisitAdd(ctx *parser.AddContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.Add{Left: left, Right: right}
}

func (v *ASTBuilder) VisitSubtract(ctx *parser.SubtractContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.Subtract{Left: left, Right: right}
}

func (v *ASTBuilder) VisitLogicOr(ctx *parser.LogicOrContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.LogicOr{Left: left, Right: right}
}

func (v *ASTBuilder) VisitTypeAsc(ctx *parser.TypeAscContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	type_ := parseType(ctx.GetType_(), v)
	return nodes.TypeAsc{Expr_: expr, Type_: type_}
}

func (v *ASTBuilder) VisitTypeCast(ctx *parser.TypeCastContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	type_ := parseType(ctx.GetType_(), v)
	return nodes.TypeCast{Expr_: expr, Type_: type_}
}

func (v *ASTBuilder) VisitAbstraction(ctx *parser.AbstractionContext) interface{} {
	params := parseParameterDeclarations(ctx.GetParamDecls(), v)
	returnExpr := parseExpr(ctx.GetReturnExpr(), v)

	return nodes.Abstraction{Params: params, ReturnExpr: returnExpr}
}

func (v *ASTBuilder) VisitTuple(ctx *parser.TupleContext) interface{} {
	return nodes.Tuple{Exprs: parseListOfExpr(ctx.GetExprs(), v)}
}

func (v *ASTBuilder) VisitRecord(ctx *parser.RecordContext) interface{} {
	return nodes.Record{Bindings: parseListOfBinding(ctx.GetBindings(), v)}
}

func (v *ASTBuilder) VisitVariant(ctx *parser.VariantContext) interface{} {
	label := parseStellaIdent(ctx.GetLabel())

	expr := optional.Empty[nodes.Expr]()

	if ctx.GetRhs() != nil {
		expr = optional.Of(parseExpr(ctx.GetRhs(), v))
	}

	return nodes.Variant{Label: label, Rhs: expr}
}

func (v *ASTBuilder) VisitMatch(ctx *parser.MatchContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	cases := parseListOfMatchCase(ctx.GetCases(), v)

	return nodes.Match{Expr_: expr, Cases: cases}
}

func (v *ASTBuilder) VisitList(ctx *parser.ListContext) interface{} {
	return nodes.List{Exprs: parseListOfExpr(ctx.GetExprs(), v)}
}

func (v *ASTBuilder) VisitLessThan(ctx *parser.LessThanContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.LessThan{Left: left, Right: right}
}

func (v *ASTBuilder) VisitLessThanOrEqual(ctx *parser.LessThanOrEqualContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.LessThanOrEqual{Left: left, Right: right}
}
func (v *ASTBuilder) VisitGreaterThan(ctx *parser.GreaterThanContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.GreaterThan{Left: left, Right: right}
}
func (v *ASTBuilder) VisitGreaterThanOrEqual(ctx *parser.GreaterThanOrEqualContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.GreaterThanOrEqual{Left: left, Right: right}
}
func (v *ASTBuilder) VisitEqual(ctx *parser.EqualContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.Equal{Left: left, Right: right}
}
func (v *ASTBuilder) VisitNotEqual(ctx *parser.NotEqualContext) interface{} {
	left := parseExpr(ctx.GetLeft(), v)
	right := parseExpr(ctx.GetRight(), v)
	return nodes.NotEqual{Left: left, Right: right}
}

func (v *ASTBuilder) VisitAssign(ctx *parser.AssignContext) interface{} {
	lhs := parseExpr(ctx.GetLhs(), v)
	rhs := parseExpr(ctx.GetRhs(), v)

	return nodes.Assign{Lhs: lhs, Rhs: rhs}
}

func (v *ASTBuilder) VisitIf(ctx *parser.IfContext) interface{} {
	condition := parseExpr(ctx.GetCondition(), v)
	thenExpr := parseExpr(ctx.GetThenExpr(), v)
	elseExpr := parseExpr(ctx.GetElseExpr(), v)

	return nodes.If{Condition: condition, ThenExpr: thenExpr, ElseExpr: elseExpr}
}

func (v *ASTBuilder) VisitLet(ctx *parser.LetContext) interface{} {
	patternBindings := parseListOfPatternBinding(ctx.GetPatternBindings(), v)
	body := parseExpr(ctx.GetBody(), v)

	return nodes.Let{PatternBindings: patternBindings, Body: body}
}

func (v *ASTBuilder) VisitLetRec(ctx *parser.LetRecContext) interface{} {
	patternBindings := parseListOfPatternBinding(ctx.GetPatternBindings(), v)
	body := parseExpr(ctx.GetBody(), v)

	return nodes.LetRec{PatternBindings: patternBindings, Body: body}
}

func (v *ASTBuilder) VisitTypeAbstraction(ctx *parser.TypeAbstractionContext) interface{} {
	generics := parseListOfStellaIdent(ctx.GetGenerics())
	body := parseExpr(ctx.GetExpr_(), v)

	return nodes.TypeAbstraction{Generics: generics, Expr_: body}
}

func (v *ASTBuilder) VisitParenthesisedExpr(ctx *parser.ParenthesisedExprContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.ParenthesisedExpr{Expr_: expr}
}

func (v *ASTBuilder) VisitTerminatingSemicolon(ctx *parser.TerminatingSemicolonContext) interface{} {
	expr := parseExpr(ctx.GetExpr_(), v)
	return nodes.TerminatingSemicolon{Expr_: expr}
}

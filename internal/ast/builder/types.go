package ast

import (
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/parser"
)

func (v *ASTBuilder) VisitTypeBool(ctx *parser.TypeBoolContext) interface{} {
	return nodes.TypeBool{}
}

func (v *ASTBuilder) VisitTypeNat(ctx *parser.TypeNatContext) interface{} {
	return nodes.TypeNat{}
}

func (v *ASTBuilder) VisitTypeRef(ctx *parser.TypeRefContext) interface{} {
	subtype := parseType(ctx.GetType_(), v)
	return nodes.TypeRef{Type_: subtype}
}

func (v *ASTBuilder) VisitTypeSum(ctx *parser.TypeSumContext) interface{} {
	left := parseType(ctx.GetLeft(), v)
	right := parseType(ctx.GetRight(), v)
	return nodes.TypeSum{Left: left, Right: right}
}

func (v *ASTBuilder) VisitTypeFun(ctx *parser.TypeFunContext) interface{} {
	paramTypes := parseListOfTypes(ctx.GetParamTypes(), v)
	returnType := parseType(ctx.GetReturnType(), v)
	return nodes.TypeFun{ParamTypes: paramTypes, ReturnType: returnType}
}

func (v *ASTBuilder) VisitTypeForAll(ctx *parser.TypeForAllContext) interface{} {
	typeNameList := parseListOfStellaIdent(ctx.GetTypes())
	type_ := parseType(ctx.GetType_(), v)

	return nodes.TypeForAll{Types: typeNameList, Type_: type_}
}

func (v *ASTBuilder) VisitTypeRec(ctx *parser.TypeRecContext) interface{} {
	var_ := parseStellaIdent(ctx.GetVar_())
	type_ := parseType(ctx.GetType_(), v)

	return nodes.TypeRec{Var_: var_, Type_: type_}
}

func (v *ASTBuilder) VisitTypeTuple(ctx *parser.TypeTupleContext) interface{} {
	types := parseListOfTypes(ctx.GetTypes(), v)
	return nodes.TypeTuple{Types: types}
}

func (v *ASTBuilder) VisitTypeRecord(ctx *parser.TypeRecordContext) interface{} {
	fieldTypes := parseListOfRecordFieldType(ctx.GetFieldTypes(), v)
	return nodes.TypeRecord{FieldTypes: fieldTypes}
}

func (v *ASTBuilder) VisitTypeVariant(ctx *parser.TypeVariantContext) interface{} {
	fieldTypes := parseListOfVariantFieldType(ctx.GetFieldTypes(), v)
	return nodes.TypeVariant{FieldTypes: fieldTypes}
}

func (v *ASTBuilder) VisitTypeList(ctx *parser.TypeListContext) interface{} {
	subtype := parseType(ctx.GetType_(), v)
	return nodes.TypeList{Type_: subtype}
}

func (v *ASTBuilder) VisitTypeUnit(ctx *parser.TypeUnitContext) interface{} {
	return nodes.TypeUnit{}
}

func (v *ASTBuilder) VisitTypeTop(ctx *parser.TypeTopContext) interface{} {
	return nodes.TypeTop{}
}

func (v *ASTBuilder) VisitTypeBottom(ctx *parser.TypeBottomContext) interface{} {
	return nodes.TypeBot{}
}

func (v *ASTBuilder) VisitTypeAuto(ctx *parser.TypeAutoContext) interface{} {
	return nodes.TypeAuto{}
}

func (v *ASTBuilder) VisitTypeVar(ctx *parser.TypeVarContext) interface{} {
	return nodes.TypeVar{Name: parseStellaIdent(ctx.GetName())}
}

func (v *ASTBuilder) VisitTypeParens(ctx *parser.TypeParensContext) interface{} {
	subtype := parseType(ctx.GetType_(), v)
	return nodes.TypeParens{Type_: subtype}
}

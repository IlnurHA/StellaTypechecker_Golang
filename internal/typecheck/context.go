package typecheck

import (
	nodes "typechecker/internal/ast/nodes"
	scope "typechecker/internal/typecheck/scope"

	"github.com/neocotic/go-optional"
)

type Context struct {
	scope *scope.ScopeStack
}

func NewContext() *Context {
	return &Context{scope: scope.NewScopeStack()}
}

func (ctx *Context) RemoveLastScope() {
	ctx.scope.RemoveLastScope()
}

func (ctx *Context) AddNewScope() {
	ctx.scope.AddNewScope()
}

func (ctx *Context) AddVar(variable nodes.StellaIdent, type_ nodes.StellaType) bool {
	return ctx.scope.AddVar(variable, type_)
}

func (ctx *Context) GetVarType(variable nodes.StellaIdent) optional.Optional[nodes.StellaType] {
	return ctx.scope.GetVarType(variable)
}

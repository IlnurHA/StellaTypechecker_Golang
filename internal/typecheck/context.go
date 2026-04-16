package typecheck

import (
	nodes "typechecker/internal/ast/nodes"
	scope "typechecker/internal/typecheck/scope"

	"github.com/neocotic/go-optional"
)

type ProgramExtension int

const (
	SUBTYPING ProgramExtension = iota
	AMBIGUOUS_TYPE_AS_BOT
)

type Context struct {
	scope      scope.ScopeStack
	extensions []ProgramExtension
}

func NewContext(extensions []nodes.Extension) *Context {
	return &Context{scope: scope.NewScopeStack(), extensions: parseExtensions(extensions)}
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

func (ctx *Context) HasExtension(extension ProgramExtension) bool {
	for _, ext := range ctx.extensions {
		if ext == extension {
			return true
		}
	}
	return false
}

func parseExtensions(exts []nodes.Extension) []ProgramExtension {
	programExt := make([]ProgramExtension, len(exts))

	for _, ext := range exts {
		switch ext.Name {
		case "#structural-subtyping":
			{
				programExt = append(programExt, SUBTYPING)
			}
		case "#ambiguous-type-as-bottom":
			{
				programExt = append(programExt, AMBIGUOUS_TYPE_AS_BOT)
			}
		}
	}
	return programExt
}

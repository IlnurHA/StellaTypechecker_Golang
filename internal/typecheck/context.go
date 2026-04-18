package typecheck

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
	exceptionhandler "typechecker/internal/typecheck/exceptionHandler"
	scope "typechecker/internal/typecheck/scope"

	"github.com/neocotic/go-optional"
)

type ProgramExtension int

const (
	SUBTYPING ProgramExtension = iota
	AMBIGUOUS_TYPE_AS_BOT
)

func (pe ProgramExtension) String() string {
	return []string{
		"#structural-subtyping",
		"#ambiguous-type-as-bottom",
	}[pe]
}

type Context struct {
	scope         scope.ScopeStack
	extensions    []ProgramExtension
	exceptionType exceptionhandler.ExceptionHandler
}

func NewContext(extensions []nodes.Extension) *Context {
	return &Context{scope: scope.NewScopeStack(), extensions: parseExtensions(extensions), exceptionType: nil}
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

func (ctx *Context) GetExtensions() []ProgramExtension {
	return ctx.extensions
}

func (ctx *Context) SetExceptionType(exceptionType nodes.StellaType) *TypecheckError {
	if ctx.exceptionType != nil {
		var err TypecheckError
		// Duplicate declaration
		if _, ok := ctx.exceptionType.(*exceptionhandler.RegularExceptionHandler); ok {
			err = NewTypeCheckErrorErrorType(ERROR_DUPLICATE_EXCEPTION_TYPE)
		} else {
			err = NewTypeCheckErrorErrorType(ERROR_CONFLICTING_EXCEPTION_DECLARATIONS)
			err.AddAdditionalInfo("Only one type of exceptions are permitted. Please, use only of either open variant type or exception type")
		}
		return &err
	}

	ctx.exceptionType = exceptionhandler.NewRegularExceptionHandler(exceptionType)
	return nil
}

func (ctx *Context) SetExceptionVariant(variant nodes.ExceptionVariantDeclaration) *TypecheckError {
	if ctx.exceptionType != nil {
		if variantHandler, ok := ctx.exceptionType.(*exceptionhandler.VariantExceptionHandler); ok {
			ok := variantHandler.AddExceptionVariant(variant)

			if !ok {
				err := NewTypeCheckErrorErrorType(ERROR_DUPLICATE_EXCEPTION_VARIANT)
				err.AddAdditionalInfo(fmt.Sprintf("Open variant with this label '%s' already exists", &variant.Name))
				return &err
			}

			return nil
		}

		err := NewTypeCheckErrorErrorType(ERROR_CONFLICTING_EXCEPTION_DECLARATIONS)
		return &err
	}

	ctx.exceptionType = exceptionhandler.NewVariantExceptionHandler(variant)
	return nil
}

func (ctx *Context) GetExceptionType() nodes.StellaType {
	return ctx.exceptionType.GetExceptionType()
}

func parseExtensions(exts []nodes.Extension) []ProgramExtension {
	programExt := make([]ProgramExtension, 0)

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

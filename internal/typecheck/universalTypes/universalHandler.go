package universaltypes

import (
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/typecheck/scope"

	"github.com/neocotic/go-optional"
)

type UniversalTypeHandler struct {
	stack             scope.BaseStack[nodes.StellaIdent, bool]
	newUniversalIndex int
}

func NewUniversalTypeHandler() UniversalTypeHandler {
	return UniversalTypeHandler{
		stack:             scope.NewBaseStack[nodes.StellaIdent, bool](),
		newUniversalIndex: 0,
	}
}

func (h *UniversalTypeHandler) RemoveLastScope() {
	h.stack.RemoveLastScope()
}

func (h *UniversalTypeHandler) AddNewScope() {
	h.stack.AddNewScope()
}

func (h *UniversalTypeHandler) AddVar(variable nodes.StellaIdent, value bool) bool {
	return h.stack.AddEntry(variable, value)
}

func (h *UniversalTypeHandler) SetVar(variable nodes.StellaIdent, value bool) bool {
	return h.stack.SetEntry(variable, value)
}

func (h *UniversalTypeHandler) GetEnrty(variable nodes.StellaIdent) optional.Optional[bool] {
	return h.stack.GetEntry(variable)
}

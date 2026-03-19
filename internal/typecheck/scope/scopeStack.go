package scope

import (
	"slices"
	nodes "typechecker/internal/ast/nodes"

	"github.com/neocotic/go-optional"
)

type ScopeStack struct {
	scopes []*Scope
}

func NewScopeStack() *ScopeStack {
	return &ScopeStack{scopes: []*Scope{NewScope()}}
}

func (ss *ScopeStack) getLastScope() *Scope {
	if len(ss.scopes) == 0 {
		return nil
	}

	return ss.scopes[len(ss.scopes)-1]
}

func (ss *ScopeStack) RemoveLastScope() {
	if len(ss.scopes) == 0 {
		return
	}

	ss.getLastScope().variables = nil
	ss.scopes = slices.Delete(ss.scopes, len(ss.scopes)-1, len(ss.scopes))
}

func (ss *ScopeStack) AddNewScope() {
	ss.scopes = append(ss.scopes, NewScope())
}

func (ss *ScopeStack) AddVar(variable nodes.StellaIdent, type_ nodes.StellaType) bool {
	return ss.getLastScope().AddVar(variable, type_)
}

func (ss *ScopeStack) GetVarType(variable nodes.StellaIdent) optional.Optional[nodes.StellaType] {
	for _, scope := range slices.Backward(ss.scopes) {
		varType := scope.GetVarType(variable)

		if varType.IsEmpty() {
			continue
		}

		return varType
	}

	return optional.Empty[nodes.StellaType]()
}

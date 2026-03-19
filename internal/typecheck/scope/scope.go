package scope

import (
	nodes "typechecker/internal/ast/nodes"

	"github.com/neocotic/go-optional"
)

type Scope struct {
	variables map[nodes.StellaIdent]nodes.StellaType
}

func NewScope() *Scope {
	return &Scope{variables: make(map[nodes.StellaIdent]nodes.StellaType)}
}

func (s *Scope) GetVarType(variable nodes.StellaIdent) optional.Optional[nodes.StellaType] {
	if value, ok := s.variables[variable]; ok {
		return optional.Of(value)
	} else {
		return optional.Empty[nodes.StellaType]()
	}
}

func (s *Scope) AddVar(variable nodes.StellaIdent, type_ nodes.StellaType) bool {
	if _, ok := s.variables[variable]; ok {
		return false
	}
	return true
}

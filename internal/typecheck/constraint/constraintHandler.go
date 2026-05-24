package constraint

import (
	"fmt"
	nodes "typechecker/internal/ast/nodes"
)

type ConstraintHandler struct {
	equations   []*Equation
	newVarIndex int
}

func NewConstraintHandler() ConstraintHandler {
	return ConstraintHandler{equations: make([]*Equation, 0), newVarIndex: 0}
}

func (c *ConstraintHandler) AddEquation(eq Equation) {
	c.equations = append(c.equations, &eq)
}

func (c *ConstraintHandler) AddEquationTypes(lhs nodes.StellaType, rhs nodes.StellaType) {
	c.AddEquation(Equation{lhs: lhs, rhs: rhs})
}

func (c *ConstraintHandler) IsEmpty() bool {
	return len(c.equations) == 0
}

func (c *ConstraintHandler) PopEquation() *Equation {
	returnEq := c.equations[0]

	c.equations = c.equations[1:]

	return returnEq
}

func (c *ConstraintHandler) ApplyToAll(var_ *nodes.TypeVar, rhs nodes.StellaType) {
	for _, eq_ := range c.equations {
		eq_.apply(var_, rhs)
	}
}

func (c *ConstraintHandler) GetFreshVar() nodes.TypeVar {
	c.newVarIndex += 1
	name := fmt.Sprintf("?X%d", c.newVarIndex)
	node := nodes.TypeVar{
		Name:      nodes.StellaIdent{Name: name, Repr: name},
		Generated: true,
	}
	return node
}

func (c *ConstraintHandler) String() string {
	message := fmt.Sprintf("Total new vars: %d\nWith the following equations:\n", c.newVarIndex)
	for _, eq := range c.equations {
		message += fmt.Sprintln(eq.String())
	}
	return message
}

package ast

type Node interface{ isNode() }

type StellaIdent struct{ Name string }
type Extension struct{ Name string }
type MemoryAddress struct{ Addr string }

type LanguageDeclaration struct{ Name string }

type AProgram struct {
	LanguageDecl LanguageDeclaration
	Extensions   []Extension
	Declarations []Declaration
}

func (x *StellaIdent) isNode() {
}
func (x *Extension) isNode() {
}
func (x *MemoryAddress) isNode() {
}
func (x *LanguageDeclaration) isNode() {
}
func (x *AProgram) isNode() {
}

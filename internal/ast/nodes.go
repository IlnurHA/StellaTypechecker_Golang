package ast

type Node interface{ isNode() }

type StellaIdent struct{ name string }
type Extension struct{ name string }
type MemoryAddress struct{ addr string }

type LanguageDeclaration struct{ name string }

type AProgram struct {
	languageDecl LanguageDeclaration
	extensions   []Extension
	declarations []Declaration
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

package ast

type Node interface {
	isNode()
	String() string
}

type StellaIdent struct {
	Name, Repr string
}
type Extension struct {
	Name, Repr string
}
type MemoryAddress struct {
	Addr, Repr string
}

type LanguageDeclaration struct {
	Name, Repr string
}

type AProgram struct {
	LanguageDecl LanguageDeclaration
	Extensions   []Extension
	Declarations []Declaration

	Repr string
}

// Marker methods
func (x *StellaIdent) isNode()         {}
func (x *Extension) isNode()           {}
func (x *MemoryAddress) isNode()       {}
func (x *LanguageDeclaration) isNode() {}
func (x *AProgram) isNode()            {}

// String methods returning Repr
func (x *StellaIdent) String() string         { return x.Repr }
func (x *Extension) String() string           { return x.Repr }
func (x *MemoryAddress) String() string       { return x.Repr }
func (x *LanguageDeclaration) String() string { return x.Repr }
func (x *AProgram) String() string            { return x.Repr }

// Eq
func (x *StellaIdent) Equal(o *StellaIdent) bool {
	return o.Name == x.Name
}

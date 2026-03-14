package main

import (
	"fmt"
	"os"

	"typechecker/internal/parser"

	"github.com/antlr4-go/antlr"
)

func main() {
	input, _ := os.ReadFile("testdata/sample.input")
	is := antlr.NewInputStream(string(input))

	lexer := parser.NewstellaLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewstellaParser(stream)

	// Parse starting from your top rule
	tree := p.Program()

	// If you have a custom visitor that builds your own AST (in internal/ast)
	visitor := ast.NewBuilder()
	result := visitor.Visit(tree)
	fmt.Printf("AST: %+v\n", result)
}

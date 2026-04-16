package utils

import (
	"fmt"
	"os"
	astbuilder "typechecker/internal/ast/builder"
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/parser"
	typechecker "typechecker/internal/typecheck"

	"github.com/antlr4-go/antlr/v4"
)

func Typecheck(program string, exitOnError bool) {
	is := antlr.NewInputStream(program)

	lexer := parser.NewstellaLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewstellaParser(stream)

	tree := p.Start_Program()
	builder := astbuilder.NewASTBuilder(p)
	builtAST := tree.Accept(builder).(nodes.AProgram)

	err := typechecker.ParseProgram(builtAST)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Typecheck error:")
		fmt.Fprintln(os.Stderr, err.String())

		if exitOnError {
			os.Exit(1)
		}
	} else {
		fmt.Println("Program is well-typed")
	}
}

func TypecheckFromFile(filePath string, exitOnError bool) {
	fmt.Printf("Typechecking %s...\n", filePath)
	input, _ := os.ReadFile(filePath)
	Typecheck(string(input), exitOnError)
}

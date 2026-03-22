package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	astbuilder "typechecker/internal/ast/builder"
	nodes "typechecker/internal/ast/nodes"
	"typechecker/internal/parser"
	typechecker "typechecker/internal/typecheck"

	"github.com/antlr4-go/antlr/v4"
)

func typecheck(filePath string) {
	fmt.Printf("Typechecking %s...\n", filePath)
	input, _ := os.ReadFile(filePath)
	is := antlr.NewInputStream(string(input))

	lexer := parser.NewstellaLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewstellaParser(stream)

	tree := p.Start_Program()
	builder := astbuilder.NewASTBuilder(p)
	builtAST := tree.Accept(builder).(nodes.AProgram)

	err := typechecker.ParseProgram(builtAST)

	if err != nil {
		fmt.Println("Typecheck error:")
		fmt.Println(err.String())
	} else {
		fmt.Println("Program is well-typed")
	}
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil // Check if it is a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

func getFiles(dirPath string) ([]string, error) {
	testPaths := make([]string, 0)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			testPaths = append(testPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the path %v: %v\n", dirPath, err)
	}
	return testPaths, nil
}

func getTestPaths(dirPath string) ([]string, error) {
	exists, err := dirExists(dirPath)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("Directory not found")
	}
	return getFiles(dirPath)
}

func main() {
	dirPathPtr := flag.String("dirPath", "", "Path to get tests from")

	flag.Parse()

	if *dirPathPtr == "" {
		fmt.Println("Expected dir path")
	}

	files, err := getTestPaths(*dirPathPtr)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		typecheck(file)
	}
}

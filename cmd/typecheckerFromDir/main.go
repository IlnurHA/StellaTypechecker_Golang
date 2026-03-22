package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil // Check if it is a directory
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
	filePathPtr := flag.String("filePath", "", "Path to source code on stella")

	flag.Parse()

	if *dirPathPtr == "" && *filePathPtr == "" {
		var builder strings.Builder
		builder.WriteString("Expected either dir or file path\n")
		builder.WriteString("Usage:\n")
		builder.WriteString("\t-dirPath string\n")
		builder.WriteString("\t\tPath to get tests from\n")
		builder.WriteString("\t-filePath string\n")
		builder.WriteString("\t\tPath to source code on stella")
		fmt.Println(builder.String())
		return
	}

	if *filePathPtr != "" {
		exists, err := fileExists(*filePathPtr)

		if err != nil {
			fmt.Println(err)
			return
		}

		if !exists {
			fmt.Printf("File does not exist: %s\n", *filePathPtr)
			return
		}

		typecheck(*filePathPtr)
	}

	if *dirPathPtr != "" {
		files, err := getTestPaths(*dirPathPtr)

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, file := range files {
			typecheck(file)
		}
	}
}

//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type methodSpec struct {
	recvType string
	methods  []string
}

func main() {
	// Read mapping
	mapping, err := readMapping("struct_interfaces.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading mapping: %v\n", err)
		os.Exit(1)
	}

	// Process all .go files in the ast directory
	files, err := filepath.Glob("ast/*.go")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if err := processFile(file, mapping); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", file, err)
		}
	}
}

func readMapping(filename string) (map[string][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	mapping := make(map[string][]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		structName := strings.TrimSpace(parts[0])
		ifaces := strings.Split(parts[1], ",")
		for i, iface := range ifaces {
			ifaces[i] = strings.TrimSpace(iface)
		}
		mapping[structName] = ifaces
	}
	return mapping, scanner.Err()
}

func processFile(filename string, mapping map[string][]string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Collect existing methods per type
	existing := make(map[string]map[string]bool)
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) != 1 {
			continue
		}
		recv := funcDecl.Recv.List[0].Type
		var typeName string
		switch t := recv.(type) {
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				typeName = ident.Name
			}
		case *ast.Ident:
			typeName = t.Name
		}
		if typeName != "" {
			if existing[typeName] == nil {
				existing[typeName] = make(map[string]bool)
			}
			existing[typeName][funcDecl.Name.Name] = true
		}
	}
	fmt.Printf("Printing %v", existing)

	// Find struct types and add missing methods
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structName := typeSpec.Name.Name
			ifaces, ok := mapping[structName]
			if !ok {
				continue
			}
			// Build method declarations
			var methods []*ast.FuncDecl
			for _, iface := range ifaces {
				methodName := "is" + iface
				if existing[structName][methodName] {
					continue
				}
				methods = append(methods, createMarkerMethod(structName, methodName))
			}
			if len(methods) == 0 {
				continue
			}
			// Insert methods after the struct definition
			// We'll append them at the end of the file for simplicity.
			// A more precise insertion would place them after the struct.
			// We'll collect them and add to node.Decls later.
			for _, m := range methods {
				node.Decls = append(node.Decls, m)
			}
		}
	}

	// Write back the file
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		return err
	}
	return os.WriteFile(filename, buf.Bytes(), 0644)
}

func createMarkerMethod(typeName, methodName string) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("x")},
					Type: &ast.StarExpr{
						X: ast.NewIdent(typeName),
					},
				},
			},
		},
		Name: ast.NewIdent(methodName),
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{},
	}
}

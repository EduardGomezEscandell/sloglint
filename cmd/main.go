package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"example.com/sloglint"
)

func main() {
	v := sloglint.Visitor{FileSet: token.NewFileSet()}
	// v.Debug = true
	for _, filePath := range os.Args[1:] {
		if filePath == "--" { // to be able to run this like "go run main.go -- input.go"
			continue
		}

		f, err := parser.ParseFile(v.FileSet, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}

		ast.Walk(&v, f)
	}
}

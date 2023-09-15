package sloglint

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Visitor struct {
	Analyser
	FileSet *token.FileSet
}

type arg int

const (
	argKey arg = iota
	argValue
)

func (a arg) next() arg {
	return arg((int(a) + 1) % 2)
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	expr, ok := node.(*ast.CallExpr)
	if !ok {
		return v
	}

	f, ok := expr.Fun.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	id, ok := f.X.(*ast.Ident)
	if !ok {
		// Some other selector, for instance getObject().method()
		return v
	}

	packageName := id.Name
	funcName := f.Sel.Name

	if packageName != "slog" {
		return v
	}

	if expr.Args == nil {
		return v
	}

	// We let the compiler check the first argument
	expecting := argKey
	for i, arg := range expr.Args[1:] {
		if expecting == argValue {
			expecting = expecting.next()
			continue
		}

		if keyType := v.ExprReturns(arg); keyType.isKnown() && !keyType.isType("string") {
			fmt.Printf("%s: function call to '%s.%s': argument #%d should be of type string, but is %s.\n",
				v.FileSet.Position(node.Pos()), packageName, funcName, i+1, keyType.Name)
		}

		expecting = expecting.next()
	}

	if expecting != argValue {
		// Well-formed call: correct number of arguments
		return v
	}

	fmt.Printf("%s: function call to '%s.%s': last key has no associated value.\n",
		v.FileSet.Position(node.Pos()), packageName, funcName)

	return v
}

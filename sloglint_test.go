package sloglint_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"

	"example.com/sloglint"
	"github.com/stretchr/testify/require"
)

func TestAnalyser(t *testing.T) {
	testCases := map[string]struct {
		want map[int]string
	}{
		"string": {want: map[int]string{
			6: "{typename string}",
		}},
		"int": {want: map[int]string{
			6: "{typename int}",
		}},
		"int alias": {want: map[int]string{
			6: "{typename int}",
		}},
		"string alias": {want: map[int]string{
			6: "{typename string}",
		}},
		"int typedef": {want: map[int]string{
			6: "{typename int}",
		}},
		"string typedef": {want: map[int]string{
			6: "{typename string}",
		}},
		"undeclared": {want: map[int]string{
			6: "{unknown undeclaredVariable}",
		}},
		"MyType": {want: map[int]string{
			9: "{struct MyType}",
		}},
		"MyType.string": {want: map[int]string{
			9: "{typename string}",
		}},
		"http.Client": {want: map[int]string{
			11: "{unknown http.Client}",
		}},
		"function call string": {want: map[int]string{
			6: "{typename string}",
		}},
		"function call int": {want: map[int]string{
			6: "{typename int}",
		}},
		"function call empty return": {want: map[int]string{
			6: "{empty }",
		}},
		"external function call": {want: map[int]string{
			9: "{unknown log.Flags}",
		}},
		"function literal int": {want: map[int]string{
			6: "{typename int}",
		}},
		"function literal string": {want: map[int]string{
			6: "{typename string}",
		}},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			v := testVisitor{fset: token.NewFileSet(), m: make(map[int]string)}
			v.Analyser.Debug = true

			path := filepath.Join("testdata", t.Name(), "main.go")

			f, err := parser.ParseFile(v.fset, path, nil, 0)
			require.NoErrorf(t, err, "Setup: Failed to parse file %s", path)

			ast.Walk(&v, f)

			require.Equal(t, tc.want, v.m)
		})
	}
}

type testVisitor struct {
	sloglint.Analyser
	fset *token.FileSet

	m map[int]string
}

func (v *testVisitor) Visit(node ast.Node) ast.Visitor {
	// 1 Filter irrelevant fields

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

	if f.X.(*ast.Ident).Name != "slog" {
		return v
	}

	for i, arg := range expr.Args[1:] {
		if i%2 == 1 {
			continue
		}

		v.m[v.fset.Position(node.Pos()).Line] = string(v.ExprReturns(arg).String())
	}

	return v
}

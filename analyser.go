package sloglint

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"runtime"
	"strings"
)

type Analyser struct {
	Debug       bool
	debugIndent int
}

func (a *Analyser) logDebug(node any, ret *object) func() {
	if !a.Debug {
		return func() {}
	}

	a.debugIndent++

	buffer := make([]uintptr, 1)
	runtime.Callers(2, buffer)
	f, _ := runtime.CallersFrames(buffer).Next()
	tmp := strings.Split(f.Function, "(*Analyser).")
	callerName := tmp[len(tmp)-1]

	fmt.Fprintf(os.Stderr,
		strings.Repeat("  ", a.debugIndent)+"%s: (%T) %#v\n",
		callerName, node, node)

	return func() {
		fmt.Fprintf(os.Stderr,
			strings.Repeat("  ", a.debugIndent)+"%s^ %s\n",
			callerName, ret)

		a.debugIndent--
	}
}

func (a *Analyser) ExprReturns(expr ast.Expr) (obj object) {
	defer a.logDebug(expr, &obj)()

	switch e := expr.(type) {
	case *ast.Ident:
		return a.identType(e)
	case *ast.BasicLit:
		return a.basicLitType(e.Kind)
	case *ast.CallExpr:
		return a.ExprReturns(e.Fun)
	case *ast.SelectorExpr:
		return a.selectorExprReturns(e)
	case *ast.FuncLit:
		return a.ExprReturns(e.Type)
	case *ast.FuncType:
		return a.funcTypeReturns(e)
	}

	return object{Type: unknown, Name: "expression"}
}

func (a *Analyser) selectorExprReturns(sel *ast.SelectorExpr) (obj object) {
	defer a.logDebug(sel, &obj)()

	parent := a.ExprReturns(sel.X)
	child := sel.Sel.Name

	if parent.Type != structure {
		return object{Type: unknown, Name: fmt.Sprintf("%s.%s", parent.Name, child)}
	}

	x, ok := parent.Children[child]
	if !ok {
		return object{Type: unknown, Name: "Illegal struct access"}
	}

	return x
}

func (a *Analyser) declReturns(decl any) (obj object) {
	defer a.logDebug(decl, &obj)()
	switch d := decl.(type) {
	case ast.Expr:
		return a.ExprReturns(d)
	case *ast.ValueSpec:
		return a.ExprReturns(d.Type)
	case *ast.FuncDecl:
		return a.funcTypeReturns(d.Type)
	case *ast.Field:
		return a.ExprReturns(d.Type)
	case *ast.TypeSpec:
		return a.typeSpecIdentifier(d)
	}

	return object{Type: unknown}
}

func (a *Analyser) typeSpecIdentifier(spec *ast.TypeSpec) (obj object) {
	defer a.logDebug(spec, &obj)()

	switch t := spec.Type.(type) {
	case *ast.StructType:
		return a.structTypeDeinition(t, spec.Name.Name)
	case *ast.Ident:
		return a.identType(t)
	}

	return object{Type: unknown}
}

func (a *Analyser) structTypeDeinition(s *ast.StructType, name string) (obj object) {
	defer a.logDebug(s, &obj)()

	children := make(map[string]object)
	for _, ch := range s.Fields.List {
		child := a.ExprReturns(ch.Type)
		name := child.Name
		if ch.Names != nil {
			// Non-anonimous child
			name = ch.Names[0].Name
		}

		children[name] = child
	}

	return object{Type: structure, Name: name, Children: children}
}

func (a *Analyser) funcTypeReturns(fun *ast.FuncType) (obj object) {
	defer a.logDebug(fun, &obj)()

	if fun.Results == nil {
		return object{Type: empty}
	}

	if fun.Results.NumFields() == 0 {
		return object{Type: empty}
	}

	if fun.Results.NumFields() > 1 {
		return object{Type: tuple}
	}

	return a.ExprReturns(fun.Results.List[0].Type)
}

func (a *Analyser) identType(ident *ast.Ident) (obj object) {
	defer a.logDebug(ident, &obj)()

	if ident.Obj != nil {
		return a.declReturns(ident.Obj.Decl)
	}

	for _, tName := range basicTypes {
		if ident.Name != tName {
			continue
		}

		return object{
			Name: tName,
			Type: typeName,
		}
	}

	return object{
		Name: ident.Name,
		Type: unknown,
	}
}

func (a *Analyser) basicLitType(t token.Token) (obj object) {
	if t == token.IDENT {
		return object{Type: identifier} // A variable or function name (e.g: main)
	}

	if v, ok := basicTypes[t]; ok {
		return object{Type: typeName, Name: v}
	}

	return object{Type: unknown}
}

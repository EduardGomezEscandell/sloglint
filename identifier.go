package sloglint

import (
	"fmt"
	"go/token"
)

type identType string

const (
	unknown    identType = "unknown"    // Could not be determined
	typeName             = "typename"   // Name of a basic type. e.g: "string"
	structure            = "struct"     // Struct
	identifier           = "identifier" // Name of a variable, func, struct, module, etc.
	tuple                = "tuple"      // Return object of a string with multiple return types
	empty                = "empty"      // Return object of a function with no return type.
)

type object struct {
	Type     identType
	Name     string
	Children map[string]object
}

func (id object) String() string {
	return fmt.Sprintf("{%s %s}", id.Type, id.Name)
}

func (i object) isType(typename string) bool {
	return (i.Type == typeName || i.Type == structure) && i.Name == typename
}

func (i object) isKnown() bool {
	return i.Type != unknown
}

var basicTypes = map[token.Token]string{
	token.INT:    "int",
	token.FLOAT:  "float",
	token.IMAG:   "imag",
	token.CHAR:   "char",
	token.STRING: "string",
}

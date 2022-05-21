package output

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// NewVoidType func
func NewVoidType() *graphql.Scalar {
	return graphql.NewScalar(graphql.ScalarConfig{
		Name: "void",
		ParseValue: func(value interface{}) interface{} {
			return nil
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			return nil
		},
		Serialize: func(value interface{}) interface{} {
			return nil
		},
	})
}

func coerceInt64(value interface{}) interface{} {
	switch value := value.(type) {
	case int64:
		return value
	case *int64:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	}

	// If the value cannot be transformed into an int, return nil instead of '0'
	// to denote 'no integer found'
	return nil
}

var Int64 = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Int64",
	Description: "Int64",
	Serialize:   coerceInt64,
	ParseValue:  coerceInt64,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if intValue, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
				return intValue
			}
		}
		return nil
	},
})

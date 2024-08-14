//go:generate stringer -type=Token -trimprefix=Token

package semantic

import (
	"github.com/iancoleman/strcase"
)

type Token int

const (
	TokenComment Token = iota
	TokenMethod
	TokenMacro
	TokenVariable
	TokenString
	TokenOperator
	TokenTypeParameter
	TokenKeyword
	TokenProperty
	TokenFunction
	TokenParameter
	TokenNone
)

func Tokens() []string {
	var tokens []string
	for i := range len(_Token_index) - 2 {
		tokens = append(tokens, strcase.ToLowerCamel(Token(i).String()))
	}
	return tokens
}

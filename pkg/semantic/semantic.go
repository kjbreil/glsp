package semantic

import (
	"github.com/kjbreil/loc-macro/pkg/editreader"
)

type Semantic struct {
	Location *editreader.Range
	Token    Token
}

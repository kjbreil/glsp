package semantic

import "github.com/kjbreil/glsp/pkg/location"

type Semantic struct {
	Location *location.Range
	Token    Token
}

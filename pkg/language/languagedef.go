package language

import (
	"github.com/kjbreil/glsp/pkg/commands"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"io"
)

type LanguageDef interface {
	Parse(uri protocol.DocumentUri, r io.Reader) (File, error)
	ID() string
	Commands() []commands.Command
}

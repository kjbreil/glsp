package language

import (
	"github.com/kjbreil/glsp/pkg/commands"
	"github.com/kjbreil/glsp/pkg/completion"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"io"
)

type LanguageFunctions struct {
	GetFile func(uri protocol.DocumentUri) (*Language, File)
}

// LanguageDef is the interface that a language must implement to be supported by glsp.
type LanguageDef interface {
	Init(functions *LanguageFunctions)
	Parse(uri protocol.DocumentUri, r io.Reader) (File, error)
	ID() string
	Commands() []commands.Command
	Completions() completion.Completions
	On() *LanguageOn
}

type LanguageOn struct {
	SaveFn func(f File) error
}

func (l *LanguageOn) Save(f File) error {
	if l.SaveFn != nil {
		return l.SaveFn(f)
	}
	return nil
}

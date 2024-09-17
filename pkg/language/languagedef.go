package language

import (
	"github.com/kjbreil/glsp"
	"github.com/kjbreil/glsp/pkg/commands"
	"github.com/kjbreil/glsp/pkg/completion"
	"github.com/kjbreil/glsp/pkg/uri"
	"io"
)

type LanguageFunctions struct {
	GetFile   func(uri uri.DocumentURI) (*Language, File)
	GetSchema func(path string) string
	Notify    glsp.NotifyFunc
}

// LanguageDef is the interface that a language must implement to be supported by glsp.
type LanguageDef interface {
	Init(functions *LanguageFunctions)
	Parse(uri uri.DocumentURI, r io.Reader) (File, error)
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

package language

import (
	"github.com/kjbreil/glsp/pkg/hover"
	"github.com/kjbreil/glsp/pkg/location"
	"github.com/kjbreil/glsp/pkg/problems"
	"github.com/kjbreil/glsp/pkg/semantic"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"io"
	"sync"
)

type Language struct {
	files map[string]File
	mu    sync.Mutex
	def   LanguageDef
}

func (l *Language) CreateFile(uri string, r io.Reader) (File, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	file, err := l.def.Parse(uri, r)
	if err != nil {
		return nil, err
	}
	l.files[uri] = file

	return file, nil
}

func (l *Language) GetFromUri(uri protocol.DocumentUri) File {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.files[uri]
}

func (l *Language) DeleteUri(uri protocol.DocumentUri) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.files, uri)
}

func (l *Language) On() *LanguageOn {
	return l.def.On()
}

type File interface {
	Hover(point location.Point) *hover.Hover
	Replace(text string, r *location.Range)
	Problems() *problems.Problems
	Uri() protocol.DocumentUri
	Path() string
	Reset(s string)
	Semantics() *semantic.Semantics
	CodeActions(r *location.Range) ([]protocol.CodeAction, error)
}

package language

import (
	"github.com/kjbreil/glsp/pkg/hover"
	"github.com/kjbreil/glsp/pkg/location"
	"github.com/kjbreil/glsp/pkg/problems"
	"github.com/kjbreil/glsp/pkg/semantic"
	"github.com/kjbreil/glsp/pkg/uri"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"io"
	"sync"
)

type Language struct {
	files map[uri.DocumentURI]File
	mu    sync.Mutex
	def   LanguageDef
}

func (l *Language) CreateFile(u uri.DocumentURI, r io.Reader) (File, error) {
	// u, err := uri.ParseDocumentURI(uriString)
	// if err != nil {
	// 	return nil, err
	// }
	l.mu.Lock()
	defer l.mu.Unlock()
	file, err := l.def.Parse(u, r)
	if err != nil {
		return nil, err
	}
	l.files[u] = file

	return file, nil
}

func (l *Language) GetFromUri(u uri.DocumentURI) File {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.files[u]
}

func (l *Language) DeleteUri(u uri.DocumentURI) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.files, u)
}

func (l *Language) On() *LanguageOn {
	return l.def.On()
}

type File interface {
	Hover(point location.Point) *hover.Hover
	Replace(text string, r *location.Range)
	Problems() *problems.Problems
	Uri() uri.DocumentURI
	Path() string
	Reset(s string)
	Semantics() *semantic.Semantics
	CodeActions(r *location.Range) ([]protocol.CodeAction, error)
}

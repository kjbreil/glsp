package language

import (
	"github.com/kjbreil/glsp/pkg/hover"
	"github.com/kjbreil/glsp/pkg/problems"
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

type File interface {
	Hover(point Point) hover.Hover
	Replace(text string, r *Range)
	Problems() *problems.Problems
	Uri() protocol.DocumentUri
	Reset(s string)
}

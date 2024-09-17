package language

import (
	"context"
	"errors"
	"github.com/kjbreil/glsp"
	"github.com/kjbreil/glsp/pkg/commands"
	"github.com/kjbreil/glsp/pkg/uri"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"github.com/sourcegraph/jsonrpc2"
	"io"
	"sync"
)

type Languages struct {
	languages       map[string]*Language
	fileLanguageIDs map[uri.DocumentURI]string
	commands        *commands.Commands

	mu   sync.Mutex
	conn *jsonrpc2.Conn
}

func NewLanguages() *Languages {
	return &Languages{
		languages:       make(map[string]*Language),
		fileLanguageIDs: make(map[uri.DocumentURI]string),
		commands:        commands.New(),
		mu:              sync.Mutex{},
	}
}

func (l *Languages) notify(method string, params interface{}) error {

	return l.conn.Notify(context.Background(), method, params)

}

func (l *Languages) AddLanguage(lang LanguageDef) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, c := range lang.Commands() {
		l.commands.Register(c.Name, c.Fn)
	}

	lang.Init(&LanguageFunctions{
		GetFile:   l.GetFromUri,
		GetSchema: l.GetSchema,
		Notify: func(method string, params any) {
			l.notify(method, params)
		},
	})

	l.languages[lang.ID()] = &Language{
		files: make(map[uri.DocumentURI]File),
		mu:    sync.Mutex{},
		def:   lang,
	}
}

func (l *Languages) Get(id string) *Language {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.languages[id]
}

var (
	ErrLanguageNotFound = errors.New("language not found")
)

func (l *Languages) CreateFile(u uri.DocumentURI, langID string, r io.Reader) (File, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	lang := l.languages[langID]
	if lang == nil {
		return nil, ErrLanguageNotFound
	}
	file, err := lang.CreateFile(u, r)
	if err != nil {
		return nil, err
	}
	l.fileLanguageIDs[u] = langID
	return file, nil
}

func (l *Languages) GetFromUri(uri uri.DocumentURI) (*Language, File) {
	l.mu.Lock()
	defer l.mu.Unlock()
	langID, ok := l.fileLanguageIDs[uri]
	if !ok {
		return nil, nil
	}
	lang := l.languages[langID]
	if lang == nil {
		return nil, nil
	}
	return lang, lang.GetFromUri(uri)
}

func (l *Languages) GetSchema(path string) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	for f := range l.fileLanguageIDs {
		if f.IsPath(path) {
			return f.Schema()
		}
	}
	return "file"
}

func (l *Languages) DeleteUri(u uri.DocumentURI) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	langID, ok := l.fileLanguageIDs[u]
	if !ok {
		return nil
	}
	delete(l.fileLanguageIDs, u)
	lang := l.languages[langID]
	if lang == nil {
		return nil
	}
	lang.DeleteUri(u)

	return nil
}

func (l *Languages) CommandProvider() *protocol.ExecuteCommandOptions {
	return l.commands.Provider()
}

func (l *Languages) CommandsExecute(context *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
	return l.commands.Execute(context, params)
}

func (l *Languages) Languages(yield func(LanguageDef) bool) {
	for _, lang := range l.languages {
		if !yield(lang.def) {
			return
		}
	}
}

func (l *Languages) SetConn(conn *jsonrpc2.Conn) {
	l.conn = conn
}

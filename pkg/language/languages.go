package language

import (
	"errors"
	"github.com/kjbreil/glsp"
	"github.com/kjbreil/glsp/pkg/commands"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"io"
	"sync"
)

type Languages struct {
	languages       map[string]*Language
	fileLanguageIDs map[string]string
	commands        *commands.Commands

	mu sync.Mutex
}

func NewLanguages() *Languages {
	return &Languages{
		languages:       make(map[string]*Language),
		fileLanguageIDs: make(map[string]string),
		commands:        commands.New(),
		mu:              sync.Mutex{},
	}
}

func (l *Languages) AddLanguage(lang LanguageDef) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, c := range lang.Commands() {
		l.commands.Register(c.Name, c.Fn)
	}

	l.languages[lang.ID()] = &Language{
		files: make(map[string]File),
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

func (l *Languages) CreateFile(uri protocol.DocumentUri, langID string, r io.Reader) (File, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	lang := l.languages[langID]
	if lang == nil {
		return nil, ErrLanguageNotFound
	}
	file, err := lang.CreateFile(uri, r)
	if err != nil {
		return nil, err
	}
	l.fileLanguageIDs[uri] = langID
	return file, nil
}

func (l *Languages) GetFromUri(uri protocol.DocumentUri) (*Language, File) {
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

func (l *Languages) DeleteUri(uri protocol.DocumentUri) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	langID, ok := l.fileLanguageIDs[uri]
	if !ok {
		return nil
	}
	delete(l.fileLanguageIDs, uri)
	lang := l.languages[langID]
	if lang == nil {
		return nil
	}
	lang.DeleteUri(uri)

	return nil
}

func (l *Languages) CommandProvider() *protocol.ExecuteCommandOptions {
	return l.commands.Provider()
}

func (l *Languages) CommandsExecute(context *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
	return l.commands.Execute(context, params)
}

package server

import (
	"context"
	"github.com/kjbreil/glsp"
	"github.com/kjbreil/glsp/internal/helpers"
	"github.com/kjbreil/glsp/pkg/language"
	"github.com/kjbreil/glsp/pkg/semantic"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	glspserv "github.com/kjbreil/glsp/server"
	"github.com/sourcegraph/jsonrpc2"
	"log/slog"
)

type Server struct {
	languages *language.Languages
	handler   protocol.Handler

	logger *slog.Logger

	languageServerName string
	server             *glspserv.Server
	serverType         ServerType
	ctx                context.Context
}

type ServerType int

const (
	ServerTypeStdio ServerType = iota
	ServerTypeTcp
)

func New(opts ...func(server *Server)) *Server {
	s := &Server{
		languageServerName: "generic_lsp",
		languages:          language.NewLanguages(),
		logger:             slog.Default(),
		serverType:         ServerTypeStdio,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Run(ctx context.Context) error {
	s.ctx = ctx

	s.handler.Initialize = s.initialize
	s.handler.Initialized = s.initialized
	s.handler.Shutdown = s.shutdown
	s.handler.SetTrace = s.setTrace
	s.handler.TextDocumentDidOpen = s.textDocumentDidOpen
	s.handler.TextDocumentDidChange = s.textDocumentDidChange
	s.handler.TextDocumentDidSave = s.textDocumentDidSave
	s.handler.TextDocumentDidClose = s.textDocumentDidClose
	s.handler.TextDocumentSemanticTokensFull = s.textDocumentSemanticTokensFull
	s.handler.TextDocumentCompletion = s.textDocumentCompletion
	s.handler.TextDocumentHover = s.textDocumentHover
	//s.handler.TextDocumentDefinition = s.textDocumentDefinition
	s.handler.TextDocumentCodeAction = s.textDocumentCodeAction
	s.handler.WorkspaceExecuteCommand = s.languages.CommandsExecute

	s.server = glspserv.NewServer(
		&s.handler,
		s.languageServerName,
		false,
		func(conn *jsonrpc2.Conn) {
			s.languages.SetConn(conn)
		},
	)

	s.server.Log = s.logger.WithGroup("GLSP")
	// s.log = s.server.Log

	//done := make(chan error, 1)

	switch s.serverType {
	case ServerTypeStdio:
		err := s.server.RunStdio()
		if err != nil {
			return err
		}
	case ServerTypeTcp:
		err := s.server.RunTCP("localhost:8080")
		//s.languages.SetConn(s.server.Conn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (s *Server) shutdown(ctx *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (s *Server) setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (s *Server) initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := s.handler.CreateServerCapabilities()
	capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes:     semantic.Tokens(),
			TokenModifiers: nil,
		},
		Range: nil,
		Full:  helpers.Ptr(true),
	}
	capabilities.ExecuteCommandProvider = s.languages.CommandProvider()
	capabilities.TextDocumentSync = &protocol.TextDocumentSyncOptions{
		OpenClose: helpers.Ptr(true),
		Change:    helpers.Ptr(protocol.TextDocumentSyncKindIncremental),
		Save:      &protocol.SaveOptions{IncludeText: helpers.Ptr(true)},
		// WillSaveWaitUntil: ptr(true),
	}
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
		TriggerCharacters:       []string{"@", "®"},
		AllCommitCharacters:     nil,
		ResolveProvider:         nil,
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo:   &protocol.InitializeResultServerInfo{Name: s.languageServerName},
	}, nil
}

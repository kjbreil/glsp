package server

import (
	"context"
	"github.com/sourcegraph/jsonrpc2"
	"log/slog"
	"time"

	"github.com/kjbreil/glsp"
)

var DefaultTimeout = time.Minute

//
// Server
//

type Server struct {
	Handler     glsp.Handler
	LogBaseName string
	Debug       bool

	ctx    context.Context
	cancel context.CancelFunc

	Log              *slog.Logger
	Timeout          time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	StreamTimeout    time.Duration
	WebSocketTimeout time.Duration
	Conn             *jsonrpc2.Conn
	onConnect        func(conn *jsonrpc2.Conn)
}

func NewServer(handler glsp.Handler, logName string, debug bool, fn func(conn *jsonrpc2.Conn)) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		Handler:          handler,
		LogBaseName:      logName,
		Debug:            debug,
		ctx:              ctx,
		cancel:           cancel,
		Log:              slog.Default(),
		Timeout:          DefaultTimeout,
		ReadTimeout:      DefaultTimeout,
		WriteTimeout:     DefaultTimeout,
		StreamTimeout:    DefaultTimeout,
		WebSocketTimeout: DefaultTimeout,
		onConnect:        fn,
	}
}

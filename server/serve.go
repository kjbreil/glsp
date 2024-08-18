package server

import (
	"io"
	"log/slog"

	"github.com/gorilla/websocket"
)

// See: https://github.com/sourcegraph/go-langserver/blob/master/main.go#L179

func (s *Server) ServeStream(stream io.ReadWriteCloser, log *slog.Logger) {
	if log == nil {
		log = s.Log
	}
	log.Info("new stream connection")
	s.Conn = s.newStreamConnection(stream)
	s.onConnect(s.Conn)
	<-s.Conn.DisconnectNotify()
	log.Info("stream connection closed")
}

func (s *Server) ServeWebSocket(socket *websocket.Conn, log *slog.Logger) {
	if log == nil {
		log = s.Log
	}
	log.Info("new web socket connection")
	s.Conn = s.newWebSocketConnection(socket)
	<-s.Conn.DisconnectNotify()
	log.Info("web socket connection closed")
}

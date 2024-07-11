package server

import (
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (s *Server) RunWebSocket(address string) error {
	mux := http.NewServeMux()
	upgrader := websocket.Upgrader{CheckOrigin: func(request *http.Request) bool { return true }}

	var connectionCount uint64

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		connection, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			s.Log.Warn("error upgrading HTTP to web socket: %s", "err", err.Error())
			http.Error(writer, errors.Wrap(err, "could not upgrade to web socket").Error(), http.StatusBadRequest)
			return
		}

		log := s.Log.With("scope", "websocket", "id", atomic.AddUint64(&connectionCount, 1))
		defer func() {
			err = connection.Close()
			if err != nil {
				log.Error("connection.Close failed", "err", err.Error())
			}
			log.Info("web socket connection closed")
		}()
		s.ServeWebSocket(connection, log)
	})

	listener, err := s.newNetworkListener("tcp", address)
	if err != nil {
		return err
	}

	server := http.Server{
		Handler:      http.TimeoutHandler(mux, s.Timeout, ""),
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}

	s.Log.Info("listening for web socket connections", "address", address)
	err = server.Serve(*listener)
	return errors.Wrap(err, "WebSocket")
}

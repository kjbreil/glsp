package server

import (
	"github.com/kjbreil/glsp/pkg/language"
	"log/slog"
)

func WithLogger(logger *slog.Logger) func(*Server) {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithTcp() func(*Server) {
	return func(s *Server) {
		s.serverType = ServerTypeTcp
	}
}

func WithLanguage(lang language.LanguageDef) func(*Server) {
	return func(s *Server) {
		s.languages.AddLanguage(lang)
	}
}

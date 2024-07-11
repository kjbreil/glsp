package server

import (
	"errors"
	"os"
)

func (s *Server) RunStdio() error {
	s.Log.Info("reading from stdin, writing to stdout")
	s.ServeStream(Stdio{}, nil)
	return nil
}

type Stdio struct{}

// ([io.Reader] interface)
func (Stdio) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

// ([io.Writer] interface)
func (Stdio) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

// ([io.Closer] interface)
func (Stdio) Close() error {
	return errors.Join(os.Stdin.Close(), os.Stdout.Close())
}

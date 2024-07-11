package server

import (
	"fmt"
	"log/slog"
)

type JSONRPCLogger struct {
	log *slog.Logger
}

// ([jsonrpc2.Logger] interface)
func (j *JSONRPCLogger) Printf(format string, v ...any) {
	j.log.Debug(fmt.Sprintf(format, v...))
}

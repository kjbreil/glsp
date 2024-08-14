package server

import (
	"github.com/kjbreil/glsp"
	"github.com/kjbreil/glsp/pkg/language"
	"github.com/kjbreil/glsp/pkg/problems"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
)

func (s *Server) publishDiagnostics(ctx *glsp.Context, file language.File, maxLevel problems.ProblemLevel) {
	if file == nil {
		return
	}
	diagnostics := file.Problems().ProtocolDiagnostics(maxLevel)
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         file.Uri(),
		Diagnostics: diagnostics,
	})
}

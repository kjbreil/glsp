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

	//diagnostics := []protocol.Diagnostic{}
	//for _, err := range file.Problems().Slice() {
	//	if err.Level <= maxLevel {
	//		diagnostics = append(diagnostics, protocol.Diagnostic{
	//			Range:    err.Location.ProtocolRange(),
	//			Severity: problemLevelToSeverity(err.Level),
	//			Message:  err.Error().Error(),
	//		})
	//	}
	//
	//}
	//
	//uri, _ := pathToDocumentURI(file.Path)
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         file.Uri(),
		Diagnostics: diagnostics,
	})
}

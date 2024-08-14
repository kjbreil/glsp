package server

import (
	"errors"
	"github.com/kjbreil/glsp"
	"github.com/kjbreil/glsp/pkg/language"
	"github.com/kjbreil/glsp/pkg/problems"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"strings"
)

var (
	ErrFileNotOpened = errors.New("file not opened")
	ErrConfigIssue   = errors.New("configuration issue")
)

func (s *Server) textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {

	file, err := s.languages.CreateFile(params.TextDocument.URI, params.TextDocument.LanguageID, strings.NewReader(params.TextDocument.Text))
	if err != nil {
		return err
	}
	s.publishDiagnostics(ctx, file, problems.ProblemLevelNone)
	return nil
}

func (s *Server) textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	_, file := s.languages.GetFromUri(params.TextDocument.URI)
	if file == nil {
		return ErrFileNotOpened
	}

	for _, ct := range params.ContentChanges {
		switch cc := ct.(type) {
		case protocol.TextDocumentContentChangeEvent:

			file.Replace(cc.Text, language.ProtocolRange(cc.Range))

			s.publishDiagnostics(ctx, file, problems.ProblemLevelNone)

		case protocol.TextDocumentContentChangeEventWhole:
			return ErrConfigIssue
		}

	}
	return nil
}

func (s *Server) textDocumentDidSave(ctx *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	_, file := s.languages.GetFromUri(params.TextDocument.URI)
	if file == nil {
		return ErrFileNotOpened
	}

	file.Reset(*params.Text)

	s.publishDiagnostics(ctx, file, problems.ProblemLevelNone)

	//if s.syncFiles {
	//	if len(file.Errors()) > 0 {
	//		return fmt.Errorf("cannot sync file with errors")
	//	}
	//
	//	instanceFileLoc, err := instance.InstallLocation(file.Path, s.instance.Storeman)
	//	if err != nil {
	//		return nil
	//	}
	//
	//	err = s.instance.WriteFile(instanceFileLoc, file.Bytes(false))
	//	// err = file.WriteFile(instanceFileLoc)
	//	if err != nil {
	//		return err
	//	}
	//	d, err := s.option.SingleDSS(file.Path)
	//	if err == nil && d != nil {
	//		text, err := d.SIL.Marshal(false)
	//		if err != nil {
	//			return err
	//		}
	//		timestamp := time.Now().Format("20060102150405")
	//		err = s.instance.WriteInbox(fmt.Sprintf("%s_dss.sql", timestamp), text)
	//		if err != nil {
	//			return err
	//		}
	//
	//	}
	//}

	return nil
}

func (s *Server) textDocumentDidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	return s.languages.DeleteUri(params.TextDocument.URI)
}

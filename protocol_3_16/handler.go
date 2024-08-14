package protocol

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/kjbreil/glsp"
)

type Handler struct {
	// Base Protocol
	CancelRequest CancelRequestFunc
	Progress      ProgressFunc

	// General Messages
	Initialize  InitializeFunc
	Initialized InitializedFunc
	Shutdown    ShutdownFunc
	Exit        ExitFunc
	LogTrace    LogTraceFunc
	SetTrace    SetTraceFunc

	// Window
	WindowWorkDoneProgressCancel WindowWorkDoneProgressCancelFunc

	// Workspace
	WorkspaceDidChangeWorkspaceFolders WorkspaceDidChangeWorkspaceFoldersFunc
	WorkspaceDidChangeConfiguration    WorkspaceDidChangeConfigurationFunc
	WorkspaceDidChangeWatchedFiles     WorkspaceDidChangeWatchedFilesFunc
	WorkspaceSymbol                    WorkspaceSymbolFunc
	WorkspaceExecuteCommand            WorkspaceExecuteCommandFunc
	WorkspaceWillCreateFiles           WorkspaceWillCreateFilesFunc
	WorkspaceDidCreateFiles            WorkspaceDidCreateFilesFunc
	WorkspaceWillRenameFiles           WorkspaceWillRenameFilesFunc
	WorkspaceDidRenameFiles            WorkspaceDidRenameFilesFunc
	WorkspaceWillDeleteFiles           WorkspaceWillDeleteFilesFunc
	WorkspaceDidDeleteFiles            WorkspaceDidDeleteFilesFunc
	WorkspaceSemanticTokensRefresh     WorkspaceSemanticTokensRefreshFunc

	// Text Document Synchronization
	TextDocumentDidOpen           TextDocumentDidOpenFunc
	TextDocumentDidChange         TextDocumentDidChangeFunc
	TextDocumentWillSave          TextDocumentWillSaveFunc
	TextDocumentWillSaveWaitUntil TextDocumentWillSaveWaitUntilFunc
	TextDocumentDidSave           TextDocumentDidSaveFunc
	TextDocumentDidClose          TextDocumentDidCloseFunc

	// Language Features
	TextDocumentCompletion              TextDocumentCompletionFunc
	CompletionItemResolve               CompletionItemResolveFunc
	TextDocumentHover                   TextDocumentHoverFunc
	TextDocumentSignatureHelp           TextDocumentSignatureHelpFunc
	TextDocumentDeclaration             TextDocumentDeclarationFunc
	TextDocumentDefinition              TextDocumentDefinitionFunc
	TextDocumentTypeDefinition          TextDocumentTypeDefinitionFunc
	TextDocumentImplementation          TextDocumentImplementationFunc
	TextDocumentReferences              TextDocumentReferencesFunc
	TextDocumentDocumentHighlight       TextDocumentDocumentHighlightFunc
	TextDocumentDocumentSymbol          TextDocumentDocumentSymbolFunc
	TextDocumentCodeAction              TextDocumentCodeActionFunc
	CodeActionResolve                   CodeActionResolveFunc
	TextDocumentCodeLens                TextDocumentCodeLensFunc
	CodeLensResolve                     CodeLensResolveFunc
	TextDocumentDocumentLink            TextDocumentDocumentLinkFunc
	DocumentLinkResolve                 DocumentLinkResolveFunc
	TextDocumentColor                   TextDocumentColorFunc
	TextDocumentColorPresentation       TextDocumentColorPresentationFunc
	TextDocumentFormatting              TextDocumentFormattingFunc
	TextDocumentRangeFormatting         TextDocumentRangeFormattingFunc
	TextDocumentOnTypeFormatting        TextDocumentOnTypeFormattingFunc
	TextDocumentRename                  TextDocumentRenameFunc
	TextDocumentPrepareRename           TextDocumentPrepareRenameFunc
	TextDocumentFoldingRange            TextDocumentFoldingRangeFunc
	TextDocumentSelectionRange          TextDocumentSelectionRangeFunc
	TextDocumentPrepareCallHierarchy    TextDocumentPrepareCallHierarchyFunc
	CallHierarchyIncomingCalls          CallHierarchyIncomingCallsFunc
	CallHierarchyOutgoingCalls          CallHierarchyOutgoingCallsFunc
	TextDocumentSemanticTokensFull      TextDocumentSemanticTokensFullFunc
	TextDocumentSemanticTokensFullDelta TextDocumentSemanticTokensFullDeltaFunc
	TextDocumentSemanticTokensRange     TextDocumentSemanticTokensRangeFunc
	TextDocumentLinkedEditingRange      TextDocumentLinkedEditingRangeFunc
	TextDocumentMoniker                 TextDocumentMonikerFunc

	initialized bool
	lock        sync.Mutex

	contextMap map[int]*glsp.Context
}

func (h *Handler) registerCtx(context *glsp.Context) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.contextMap == nil {
		h.contextMap = make(map[int]*glsp.Context)
	}
	if context.ID.Num == 0 {
		return
	}

	h.contextMap[int(context.ID.Num)] = context
}

func (h *Handler) StopContext(context *glsp.Context) {

	if context.ID.Num == 0 {
		return
	}

	h.StopID(int(context.ID.Num))
}

func (h *Handler) StopID(id int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.contextMap == nil {
		return
	}
	if id == 0 {
		return
	}

	if ctx, ok := h.contextMap[id]; ok {
		if ctx.Cancel != nil {
			ctx.Cancel()
		}
	}

	delete(h.contextMap, id)
}

// ([glsp.Handler] interface)
func (h *Handler) Handle(context *glsp.Context) (r any, validMethod bool, validParams bool, err error) {
	if !h.IsInitialized() && (context.Method != MethodInitialize) {
		return nil, true, true, errors.New("server not initialized")
	}

	h.registerCtx(context)

	switch context.Method {
	// Base Protocol

	case MethodCancelRequest:
		if h.CancelRequest != nil {
			validMethod = true
			var params CancelParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.CancelRequest(context, &params) })
			}
		}

	case MethodProgress:
		if h.Progress != nil {
			validMethod = true
			var params ProgressParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.Progress(context, &params) })
			}
		}

	// General Messages

	case MethodInitialize:
		if h.Initialize != nil {
			validMethod = true
			var params InitializeParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				if r, err = h.Initialize(context, &params); err == nil {
					h.SetInitialized(true)
				}
			}
		}

	case MethodInitialized:
		if h.Initialized != nil {
			validMethod = true
			var params InitializedParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.Initialized(context, &params) })
			}
		}

	case MethodShutdown:
		h.SetInitialized(false)
		if h.Shutdown != nil {
			validMethod = true
			validParams = true
			err = cancelErr(context.Context, func() error { return h.Shutdown(context) })
		}

	case MethodExit:
		// Note that the server will close the connection after we handle it here
		if h.Exit != nil {
			validMethod = true
			validParams = true
			err = cancelErr(context.Context, func() error { return h.Exit(context) })
		}

	case MethodLogTrace:
		if h.LogTrace != nil {
			validMethod = true
			var params LogTraceParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.LogTrace(context, &params) })
			}
		}

	case MethodSetTrace:
		if h.SetTrace != nil {
			validMethod = true
			var params SetTraceParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.SetTrace(context, &params) })
			}
		}

	// Window

	case MethodWindowWorkDoneProgressCancel:
		if h.WindowWorkDoneProgressCancel != nil {
			validMethod = true
			var params WorkDoneProgressCancelParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.WindowWorkDoneProgressCancel(context, &params) })
			}
		}

	// Workspace

	case MethodWorkspaceDidChangeWorkspaceFolders:
		if h.WorkspaceDidChangeWorkspaceFolders != nil {
			validMethod = true
			var params DidChangeWorkspaceFoldersParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.WorkspaceDidChangeWorkspaceFolders(context, &params) })
			}
		}

	case MethodWorkspaceDidChangeConfiguration:
		if h.WorkspaceDidChangeConfiguration != nil {
			validMethod = true
			var params DidChangeConfigurationParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.WorkspaceDidChangeConfiguration(context, &params) })
			}
		}

	case MethodWorkspaceDidChangeWatchedFiles:
		if h.WorkspaceDidChangeWatchedFiles != nil {
			validMethod = true
			var params DidChangeWatchedFilesParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.WorkspaceDidChangeWatchedFiles(context, &params) })
			}
		}

	case MethodWorkspaceSymbol:
		if h.WorkspaceSymbol != nil {
			validMethod = true
			var params WorkspaceSymbolParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.WorkspaceSymbol(context, &params) })
			}
		}

	case MethodWorkspaceExecuteCommand:
		if h.WorkspaceExecuteCommand != nil {
			validMethod = true
			var params ExecuteCommandParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.WorkspaceExecuteCommand(context, &params) })
			}
		}

	case MethodWorkspaceWillCreateFiles:
		if h.WorkspaceWillCreateFiles != nil {
			validMethod = true
			var params CreateFilesParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.WorkspaceWillCreateFiles(context, &params) })
			}
		}

	case MethodWorkspaceDidCreateFiles:
		if h.WorkspaceDidCreateFiles != nil {
			validMethod = true
			var params CreateFilesParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.WorkspaceDidCreateFiles(context, &params) })
			}
		}

	case MethodWorkspaceWillRenameFiles:
		if h.WorkspaceWillRenameFiles != nil {
			validMethod = true
			var params RenameFilesParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.WorkspaceWillRenameFiles(context, &params) })
			}
		}

	case MethodWorkspaceDidRenameFiles:
		if h.WorkspaceDidRenameFiles != nil {
			validMethod = true
			var params RenameFilesParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.WorkspaceDidRenameFiles(context, &params) })
			}
		}

	case MethodWorkspaceWillDeleteFiles:
		if h.WorkspaceWillDeleteFiles != nil {
			validMethod = true
			var params DeleteFilesParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.WorkspaceWillDeleteFiles(context, &params) })
			}
		}

	case MethodWorkspaceDidDeleteFiles:
		if h.WorkspaceDidDeleteFiles != nil {
			validMethod = true
			var params DeleteFilesParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.WorkspaceDidDeleteFiles(context, &params) })
			}
		}

	// Text Document Synchronization

	case MethodTextDocumentDidOpen:
		if h.TextDocumentDidOpen != nil {
			validMethod = true
			var params DidOpenTextDocumentParams

			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.TextDocumentDidOpen(context, &params) })
			}
		}

	case MethodTextDocumentDidChange:
		if h.TextDocumentDidChange != nil {
			validMethod = true
			var params DidChangeTextDocumentParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.TextDocumentDidChange(context, &params) })
			}
		}

	case MethodTextDocumentWillSave:
		if h.TextDocumentWillSave != nil {
			validMethod = true
			var params WillSaveTextDocumentParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.TextDocumentWillSave(context, &params) })
			}
		}

	case MethodTextDocumentWillSaveWaitUntil:
		if h.TextDocumentWillSaveWaitUntil != nil {
			validMethod = true
			var params WillSaveTextDocumentParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentWillSaveWaitUntil(context, &params) })
			}
		}

	case MethodTextDocumentDidSave:
		if h.TextDocumentDidSave != nil {
			validMethod = true
			var params DidSaveTextDocumentParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.TextDocumentDidSave(context, &params) })
			}
		}

	case MethodTextDocumentDidClose:
		if h.TextDocumentDidClose != nil {
			validMethod = true
			var params DidCloseTextDocumentParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				err = cancelErr(context.Context, func() error { return h.TextDocumentDidClose(context, &params) })
			}
		}

	// Language Features

	case MethodTextDocumentCompletion:
		if h.TextDocumentCompletion != nil {
			validMethod = true
			var params CompletionParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentCompletion(context, &params) })
			}
		}

	case MethodCompletionItemResolve:
		if h.CompletionItemResolve != nil {
			validMethod = true
			var params CompletionItem
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.CompletionItemResolve(context, &params) })
			}
		}

	case MethodTextDocumentHover:
		if h.TextDocumentHover != nil {
			validMethod = true
			var params HoverParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentHover(context, &params) })
			}
		}

	case MethodTextDocumentSignatureHelp:
		if h.TextDocumentSignatureHelp != nil {
			validMethod = true
			var params SignatureHelpParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentSignatureHelp(context, &params) })
			}
		}

	case MethodTextDocumentDeclaration:
		if h.TextDocumentDeclaration != nil {
			validMethod = true
			var params DeclarationParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentDeclaration(context, &params) })
			}
		}

	case MethodTextDocumentDefinition:
		if h.TextDocumentDefinition != nil {
			validMethod = true
			var params DefinitionParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentDefinition(context, &params) })
			}
		}

	case MethodTextDocumentTypeDefinition:
		if h.TextDocumentTypeDefinition != nil {
			validMethod = true
			var params TypeDefinitionParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentTypeDefinition(context, &params) })
			}
		}

	case MethodTextDocumentImplementation:
		if h.TextDocumentImplementation != nil {
			validMethod = true
			var params ImplementationParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentImplementation(context, &params) })
			}
		}

	case MethodTextDocumentReferences:
		if h.TextDocumentReferences != nil {
			validMethod = true
			var params ReferenceParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentReferences(context, &params) })
			}
		}

	case MethodTextDocumentDocumentHighlight:
		if h.TextDocumentDocumentHighlight != nil {
			validMethod = true
			var params DocumentHighlightParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentDocumentHighlight(context, &params) })
			}
		}

	case MethodTextDocumentDocumentSymbol:
		if h.TextDocumentDocumentSymbol != nil {
			validMethod = true
			var params DocumentSymbolParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentDocumentSymbol(context, &params) })
			}
		}

	case MethodTextDocumentCodeAction:
		if h.TextDocumentCodeAction != nil {
			validMethod = true
			var params CodeActionParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentCodeAction(context, &params) })
			}
		}

	case MethodCodeActionResolve:
		if h.CodeActionResolve != nil {
			validMethod = true
			var params CodeAction
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.CodeActionResolve(context, &params) })
			}
		}

	case MethodTextDocumentCodeLens:
		if h.TextDocumentCodeLens != nil {
			validMethod = true
			var params CodeLensParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentCodeLens(context, &params) })
			}
		}

	case MethodCodeLensResolve:
		if h.TextDocumentDidClose != nil {
			validMethod = true
			var params CodeLens
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.CodeLensResolve(context, &params) })
			}
		}

	case MethodTextDocumentDocumentLink:
		if h.TextDocumentDocumentLink != nil {
			validMethod = true
			var params DocumentLinkParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentDocumentLink(context, &params) })
			}
		}

	case MethodDocumentLinkResolve:
		if h.DocumentLinkResolve != nil {
			validMethod = true
			var params DocumentLink
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.DocumentLinkResolve(context, &params) })
			}
		}

	case MethodTextDocumentColor:
		if h.TextDocumentColor != nil {
			validMethod = true
			var params DocumentColorParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentColor(context, &params) })
			}
		}

	case MethodTextDocumentColorPresentation:
		if h.TextDocumentColorPresentation != nil {
			validMethod = true
			var params ColorPresentationParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentColorPresentation(context, &params) })
			}
		}

	case MethodTextDocumentFormatting:
		if h.TextDocumentFormatting != nil {
			validMethod = true
			var params DocumentFormattingParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentFormatting(context, &params) })
			}
		}

	case MethodTextDocumentRangeFormatting:
		if h.TextDocumentRangeFormatting != nil {
			validMethod = true
			var params DocumentRangeFormattingParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentRangeFormatting(context, &params) })
			}
		}

	case MethodTextDocumentOnTypeFormatting:
		if h.TextDocumentOnTypeFormatting != nil {
			validMethod = true
			var params DocumentOnTypeFormattingParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentOnTypeFormatting(context, &params) })
			}
		}

	case MethodTextDocumentRename:
		if h.TextDocumentRename != nil {
			validMethod = true
			var params RenameParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentRename(context, &params) })
			}
		}

	case MethodTextDocumentPrepareRename:
		if h.TextDocumentPrepareRename != nil {
			validMethod = true
			var params PrepareRenameParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentPrepareRename(context, &params) })
			}
		}

	case MethodTextDocumentFoldingRange:
		if h.TextDocumentFoldingRange != nil {
			validMethod = true
			var params FoldingRangeParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentFoldingRange(context, &params) })
			}
		}

	case MethodTextDocumentSelectionRange:
		if h.TextDocumentSelectionRange != nil {
			validMethod = true
			var params SelectionRangeParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentSelectionRange(context, &params) })
			}
		}

	case MethodTextDocumentPrepareCallHierarchy:
		if h.TextDocumentPrepareCallHierarchy != nil {
			validMethod = true
			var params CallHierarchyPrepareParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentPrepareCallHierarchy(context, &params) })
			}
		}

	case MethodCallHierarchyIncomingCalls:
		if h.CallHierarchyIncomingCalls != nil {
			validMethod = true
			var params CallHierarchyIncomingCallsParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.CallHierarchyIncomingCalls(context, &params) })
			}
		}

	case MethodCallHierarchyOutgoingCalls:
		if h.CallHierarchyOutgoingCalls != nil {
			validMethod = true
			var params CallHierarchyOutgoingCallsParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.CallHierarchyOutgoingCalls(context, &params) })
			}
		}

	case MethodTextDocumentSemanticTokensFull:
		if h.TextDocumentSemanticTokensFull != nil {
			validMethod = true
			var params SemanticTokensParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentSemanticTokensFull(context, &params) })
			}
		}

	case MethodTextDocumentSemanticTokensFullDelta:
		if h.TextDocumentSemanticTokensFullDelta != nil {
			validMethod = true
			var params SemanticTokensDeltaParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentSemanticTokensFullDelta(context, &params) })
			}
		}

	case MethodTextDocumentSemanticTokensRange:
		if h.TextDocumentSemanticTokensRange != nil {
			validMethod = true
			var params SemanticTokensRangeParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentSemanticTokensRange(context, &params) })
			}
		}

	case MethodWorkspaceSemanticTokensRefresh:
		if h.WorkspaceSemanticTokensRefresh != nil {
			validMethod = true
			validParams = true
			err = cancelErr(context.Context, func() error { return h.WorkspaceSemanticTokensRefresh(context) })
		}

	case MethodTextDocumentLinkedEditingRange:
		if h.TextDocumentLinkedEditingRange != nil {
			validMethod = true
			var params LinkedEditingRangeParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentLinkedEditingRange(context, &params) })
			}
		}

	case MethodTextDocumentMoniker:
		if h.TextDocumentMoniker != nil {
			validMethod = true
			var params MonikerParams
			if err = json.Unmarshal(context.Params, &params); err == nil {
				validParams = true
				r, err = cancelRtnErr(context.Context, func() (any, error) { return h.TextDocumentMoniker(context, &params) })
			}
		}
	}

	h.StopContext(context)

	return
}

func (h *Handler) IsInitialized() bool {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.initialized
}

func (h *Handler) SetInitialized(initialized bool) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.initialized = initialized
}

func (h *Handler) cancelRequest(context *glsp.Context, params *CancelParams) error {

	if params != nil {
		if v, ok := params.ID.Value.(int32); ok {
			h.StopID(int(v))
		}
	}
	return nil
}

func (h *Handler) CreateServerCapabilities() ServerCapabilities {
	if h.CancelRequest == nil {
		h.CancelRequest = h.cancelRequest
	}

	var capabilities ServerCapabilities

	if (h.TextDocumentDidOpen != nil) || (h.TextDocumentDidClose != nil) {
		if _, ok := capabilities.TextDocumentSync.(*TextDocumentSyncOptions); !ok {
			capabilities.TextDocumentSync = &TextDocumentSyncOptions{}
		}
		capabilities.TextDocumentSync.(*TextDocumentSyncOptions).OpenClose = &True
	}

	if h.TextDocumentDidChange != nil {
		if _, ok := capabilities.TextDocumentSync.(*TextDocumentSyncOptions); !ok {
			capabilities.TextDocumentSync = &TextDocumentSyncOptions{}
		}
		// This can be overriden to TextDocumentSyncKindFull
		value := TextDocumentSyncKindIncremental
		capabilities.TextDocumentSync.(*TextDocumentSyncOptions).Change = &value
	}

	if h.TextDocumentWillSave != nil {
		if _, ok := capabilities.TextDocumentSync.(*TextDocumentSyncOptions); !ok {
			capabilities.TextDocumentSync = &TextDocumentSyncOptions{}
		}
		capabilities.TextDocumentSync.(*TextDocumentSyncOptions).WillSave = &True
	}

	if h.TextDocumentWillSaveWaitUntil != nil {
		if _, ok := capabilities.TextDocumentSync.(*TextDocumentSyncOptions); !ok {
			capabilities.TextDocumentSync = &TextDocumentSyncOptions{}
		}
		capabilities.TextDocumentSync.(*TextDocumentSyncOptions).WillSaveWaitUntil = &True
	}

	if h.TextDocumentDidSave != nil {
		if _, ok := capabilities.TextDocumentSync.(*TextDocumentSyncOptions); !ok {
			capabilities.TextDocumentSync = &TextDocumentSyncOptions{}
		}
		capabilities.TextDocumentSync.(*TextDocumentSyncOptions).Save = &True
	}

	if h.TextDocumentCompletion != nil {
		capabilities.CompletionProvider = &CompletionOptions{}
	}

	if h.TextDocumentHover != nil {
		capabilities.HoverProvider = true
	}

	if h.TextDocumentSignatureHelp != nil {
		capabilities.SignatureHelpProvider = &SignatureHelpOptions{}
	}

	if h.TextDocumentDeclaration != nil {
		capabilities.DeclarationProvider = true
	}

	if h.TextDocumentDefinition != nil {
		capabilities.DefinitionProvider = true
	}

	if h.TextDocumentTypeDefinition != nil {
		capabilities.TypeDefinitionProvider = true
	}

	if h.TextDocumentImplementation != nil {
		capabilities.ImplementationProvider = true
	}

	if h.TextDocumentReferences != nil {
		capabilities.ReferencesProvider = true
	}

	if h.TextDocumentDocumentHighlight != nil {
		capabilities.DocumentHighlightProvider = true
	}

	if h.TextDocumentDocumentSymbol != nil {
		capabilities.DocumentSymbolProvider = true
	}

	if h.TextDocumentCodeAction != nil {
		capabilities.CodeActionProvider = true
	}

	if h.TextDocumentCodeLens != nil {
		capabilities.CodeLensProvider = &CodeLensOptions{}
	}

	if h.TextDocumentDocumentLink != nil {
		capabilities.DocumentLinkProvider = &DocumentLinkOptions{}
	}

	if h.TextDocumentColor != nil {
		capabilities.ColorProvider = true
	}

	if h.TextDocumentFormatting != nil {
		capabilities.DocumentFormattingProvider = true
	}

	if h.TextDocumentRangeFormatting != nil {
		capabilities.DocumentRangeFormattingProvider = true
	}

	if h.TextDocumentOnTypeFormatting != nil {
		capabilities.DocumentOnTypeFormattingProvider = &DocumentOnTypeFormattingOptions{}
	}

	if h.TextDocumentRename != nil {
		capabilities.RenameProvider = true
	}

	if h.TextDocumentFoldingRange != nil {
		capabilities.FoldingRangeProvider = true
	}

	if h.WorkspaceExecuteCommand != nil {
		capabilities.ExecuteCommandProvider = &ExecuteCommandOptions{}
	}

	if h.TextDocumentSelectionRange != nil {
		capabilities.SelectionRangeProvider = true
	}

	if h.TextDocumentLinkedEditingRange != nil {
		capabilities.LinkedEditingRangeProvider = true
	}

	if h.TextDocumentPrepareCallHierarchy != nil {
		capabilities.CallHierarchyProvider = true
	}

	if h.TextDocumentSemanticTokensFull != nil {
		if _, ok := capabilities.SemanticTokensProvider.(*SemanticTokensOptions); !ok {
			capabilities.SemanticTokensProvider = &SemanticTokensOptions{}
		}
		if h.TextDocumentSemanticTokensFullDelta != nil {
			capabilities.SemanticTokensProvider.(*SemanticTokensOptions).Full = &SemanticDelta{}
			capabilities.SemanticTokensProvider.(*SemanticTokensOptions).Full.(*SemanticDelta).Delta = &True
		} else {
			capabilities.SemanticTokensProvider.(*SemanticTokensOptions).Full = true
		}
	}

	if h.TextDocumentSemanticTokensRange != nil {
		if _, ok := capabilities.SemanticTokensProvider.(*SemanticTokensOptions); !ok {
			capabilities.SemanticTokensProvider = &SemanticTokensOptions{}
		}
		capabilities.SemanticTokensProvider.(*SemanticTokensOptions).Range = true
	}

	// TODO: self.TextDocumentSemanticTokensRefresh?

	if h.TextDocumentMoniker != nil {
		capabilities.MonikerProvider = true
	}

	if h.WorkspaceSymbol != nil {
		capabilities.WorkspaceSymbolProvider = true
	}

	if h.WorkspaceDidCreateFiles != nil {
		if capabilities.Workspace == nil {
			capabilities.Workspace = &ServerCapabilitiesWorkspace{}
		}
		if capabilities.Workspace.FileOperations == nil {
			capabilities.Workspace.FileOperations = &ServerCapabilitiesWorkspaceFileOperations{}
		}
		capabilities.Workspace.FileOperations.DidCreate = &FileOperationRegistrationOptions{
			Filters: []FileOperationFilter{},
		}
	}

	if h.WorkspaceWillCreateFiles != nil {
		if capabilities.Workspace == nil {
			capabilities.Workspace = &ServerCapabilitiesWorkspace{}
		}
		if capabilities.Workspace.FileOperations == nil {
			capabilities.Workspace.FileOperations = &ServerCapabilitiesWorkspaceFileOperations{}
		}
		capabilities.Workspace.FileOperations.WillCreate = &FileOperationRegistrationOptions{
			Filters: []FileOperationFilter{},
		}
	}

	if h.WorkspaceDidRenameFiles != nil {
		capabilities.RenameProvider = true
		if capabilities.Workspace == nil {
			capabilities.Workspace = &ServerCapabilitiesWorkspace{}
		}
		if capabilities.Workspace.FileOperations == nil {
			capabilities.Workspace.FileOperations = &ServerCapabilitiesWorkspaceFileOperations{}
		}
		capabilities.Workspace.FileOperations.DidRename = &FileOperationRegistrationOptions{
			Filters: []FileOperationFilter{},
		}
	}

	if h.WorkspaceWillRenameFiles != nil {
		capabilities.RenameProvider = true
		if capabilities.Workspace == nil {
			capabilities.Workspace = &ServerCapabilitiesWorkspace{}
		}
		if capabilities.Workspace.FileOperations == nil {
			capabilities.Workspace.FileOperations = &ServerCapabilitiesWorkspaceFileOperations{}
		}
		capabilities.Workspace.FileOperations.WillRename = &FileOperationRegistrationOptions{
			Filters: []FileOperationFilter{},
		}
	}

	if h.WorkspaceDidDeleteFiles != nil {
		if capabilities.Workspace == nil {
			capabilities.Workspace = &ServerCapabilitiesWorkspace{}
		}
		if capabilities.Workspace.FileOperations == nil {
			capabilities.Workspace.FileOperations = &ServerCapabilitiesWorkspaceFileOperations{}
		}
		capabilities.Workspace.FileOperations.DidDelete = &FileOperationRegistrationOptions{
			Filters: []FileOperationFilter{},
		}
	}

	if h.WorkspaceWillDeleteFiles != nil {
		if capabilities.Workspace == nil {
			capabilities.Workspace = &ServerCapabilitiesWorkspace{}
		}
		if capabilities.Workspace.FileOperations == nil {
			capabilities.Workspace.FileOperations = &ServerCapabilitiesWorkspaceFileOperations{}
		}
		capabilities.Workspace.FileOperations.WillDelete = &FileOperationRegistrationOptions{
			Filters: []FileOperationFilter{},
		}
	}

	return capabilities
}

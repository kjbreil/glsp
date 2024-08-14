package completion

import (
	"github.com/kjbreil/glsp/internal/helpers"
	"github.com/kjbreil/glsp/pkg/markdown"
	protocol316 "github.com/kjbreil/glsp/protocol_3_16"
	"strings"
)

type Completion struct {
	Markdown   markdown.Markdown
	Detail     string
	Trigger    string
	InsertText string
}

func (c *Completion) Protocol() protocol316.CompletionItem {
	if strings.Contains(c.InsertText, "$1") {
		return protocol316.CompletionItem{
			Detail: c.Markdown.Detail(),
			Documentation: protocol316.MarkupContent{
				Kind:  protocol316.MarkupKindMarkdown,
				Value: c.Markdown.String(),
			},
			Kind:             ptr(protocol316.CompletionItemKindFunction),
			Label:            c.Trigger,
			InsertText:       &c.InsertText,
			InsertTextFormat: helpers.Ptr(protocol316.InsertTextFormatSnippet),
		}
	}

	return protocol316.CompletionItem{
		Detail: c.Markdown.Detail(),
		Documentation: protocol316.MarkupContent{
			Kind:  protocol316.MarkupKindMarkdown,
			Value: c.Markdown.String(),
		},
		Kind:       ptr(protocol316.CompletionItemKindFunction),
		Label:      c.Trigger,
		InsertText: &c.InsertText,
	}
}

func ptr[T any](v T) *T {
	return &v
}

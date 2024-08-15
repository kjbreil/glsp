package hover

import (
	"github.com/kjbreil/glsp/pkg/editreader"
	"github.com/kjbreil/glsp/pkg/markdown"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
)

type Hover struct {
	Markdown markdown.Markdown

	CharRange editreader.CharRange
}

func (h *Hover) Protocol() *protocol.Hover {
	return &protocol.Hover{
		Contents: h.Markdown.String(),
		Range:    h.CharRange.ProtocolRange(),
	}
}

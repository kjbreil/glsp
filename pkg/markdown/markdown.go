package markdown

import (
	"strings"
)

type Markdown struct {
	Main        string
	Sub         string
	Description string
	Example     string
	Language    string
}

func (md *Markdown) String() string {
	sb := strings.Builder{}

	sb.WriteString("### ")
	sb.WriteString(md.Main)
	if md.Sub != "" {
		sb.WriteRune(' ')
		sb.WriteString(md.Sub)
	}

	if md.Description != "" {
		sb.WriteString("\n")
		sb.WriteString(md.Description)
	}

	if md.Example != "" {
		sb.WriteString("\n```")
		if md.Language != "" {

			sb.WriteString(md.Language)
		}
		sb.WriteRune('\n')
		sb.WriteString(md.Example)
		sb.WriteString("\n```\n")
	}

	return sb.String()
}

func (md *Markdown) Detail() *string {
	return &md.Description
}

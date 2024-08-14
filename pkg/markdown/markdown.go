package markdown

import "fmt"

var language = "sms"

type Markdown struct {
	Macro       string
	Command     string
	Description string
	Example     string
}

func (md *Markdown) String() string {
	if md.Command != "" {
		return fmt.Sprintf("### %s %s\n%s\n```%s\n%s\n```", md.Macro, md.Command, md.Description, language, md.Example)

	}
	return fmt.Sprintf("### %s\n%s\n```%s\n%s\n```", md.Macro, md.Description, language, md.Example)
}

func (md *Markdown) Detail() *string {
	return &md.Description
}

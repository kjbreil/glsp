package completion

import protocol "github.com/kjbreil/glsp/protocol_3_16"

type Completions struct {
	c []Completion
}

func NewCompletions(c ...Completion) Completions {
	return Completions{c: c}
}

func (c *Completions) Len() int {
	return len(c.c)
}

func (c *Completions) Slice() []Completion {
	return c.c
}

func (c *Completions) Index(i int) *Completion {
	if i < 0 || i >= len(c.c) {
		return &Completion{}
	}
	return &c.c[i]
}

func (c *Completions) Combine(comps Completions) {
	c.c = append(c.c, comps.c...)
}

func (c *Completions) Add(cmp Completion) {
	c.c = append(c.c, cmp)
}

func (c *Completions) Protocol() []protocol.CompletionItem {
	items := make([]protocol.CompletionItem, len(c.c))
	for i, cmp := range c.c {
		items[i] = cmp.Protocol()
	}
	return items
}

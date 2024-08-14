package commands

import (
	"fmt"
	"github.com/kjbreil/glsp"
	protocol "github.com/kjbreil/glsp/protocol_3_16"
	"sync"
)

var trueObject bool = true

type Command struct {
	Name string
	Fn   func(ctx *glsp.Context, params *protocol.ExecuteCommandParams) (any, error)
}

type Commands struct {
	c map[string]func(ctx *glsp.Context, params *protocol.ExecuteCommandParams) (any, error)
	m sync.Mutex
}

func New() *Commands {
	commands := &Commands{
		c: make(map[string]func(ctx *glsp.Context, params *protocol.ExecuteCommandParams) (any, error)),
		m: sync.Mutex{},
	}
	// Register commands here
	return commands
}

func (c *Commands) Provider() *protocol.ExecuteCommandOptions {
	return &protocol.ExecuteCommandOptions{
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
			WorkDoneProgress: &trueObject,
		},
		Commands: c.commands(),
	}
}

func (c *Commands) Register(name string, handler func(ctx *glsp.Context, params *protocol.ExecuteCommandParams) (any, error)) {
	c.m.Lock()
	defer c.m.Unlock()
	c.c[name] = handler
}

func (c *Commands) Execute(ctx *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
	c.m.Lock()
	defer c.m.Unlock()
	handler, ok := c.c[params.Command]
	if !ok {
		return nil, fmt.Errorf("command '%s' not found", params.Command)
	}
	return handler(ctx, params)
}

func (c *Commands) commands() []string {
	c.m.Lock()
	defer c.m.Unlock()
	cmds := make([]string, 0, len(c.c))
	for name := range c.c {
		cmds = append(cmds, name)
	}
	return cmds
}

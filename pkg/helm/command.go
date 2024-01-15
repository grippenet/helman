package helm

import (
	"fmt"
	"strings"
)

// Command represents a Helm command line to be executed
type Command struct {
	CheckKubeContext string // KubeContext to check before applying command
	Args             []string
}

func NewCommand(args []string) *Command {
	return &Command{Args: args}
}

func (c *Command) AddArg(args ...string) {
	c.Args = append(c.Args, args...)
}

func (c *Command) String() string {
	var b strings.Builder
	for _, a := range c.Args {
		b.WriteString(" ")
		b.WriteString(a)
	}
	if c.CheckKubeContext != "" {
		b.WriteString(fmt.Sprintf(" [context=%s]", c.CheckKubeContext))
	}
	return b.String()
}

func (c *Command) Run() error {
	fmt.Println(c.String())
	return nil
}

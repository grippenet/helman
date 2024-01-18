package helm

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
)

// Command represents a Helm command line to be executed
type Command struct {
	CheckKubeContext string // KubeContext to check before applying command
	AskForDryRun     bool
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

	if c.AskForDryRun {
		b.WriteString(" [dry-run]")
	}
	return b.String()
}

func (c *Command) Run() error {
	fmt.Println(c.String())
	if c.AskForDryRun {
		r := PromptYesNoCancel("Do you want to run with --dry-run before ?")
		fmt.Println(r)
	}
	fmt.Println("I'm not complete so I wont do anything")
	return nil
}

func PromptYesNoCancel(label string) string {
	completer := func(d prompt.Document) []prompt.Suggest {
		s := []prompt.Suggest{
			{Text: "Y", Description: "Yes"},
			{Text: "N", Description: "No"},
			{Text: "C", Description: "Cancel"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
	fmt.Println(label, " (Yes/No)")
	t := prompt.Input("> ", completer)
	return t
}

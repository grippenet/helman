package helm

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ErrHelmNotFound = errors.New("helm binary has not been found")

// Command represents a Helm command line to be executed
type Command struct {
	HelmBinary       string
	CheckKubeContext string // KubeContext to check before applying command
	AskForDryRun     bool
	Args             []string
	Dir              string
}

func NewCommand(args []string) *Command {
	return &Command{Args: args}
}

func (c *Command) AddArg(args ...string) {
	c.Args = append(c.Args, args...)
}

func (c *Command) String() string {
	var b strings.Builder

	bin, err := FindHelm()
	if err != nil {
		fmt.Println("Unable to find helm binary ", err)
	} else {
		b.WriteString(bin)
		b.WriteString(" ")
	}

	b.WriteString(c.BuildArgs(nil))

	if c.CheckKubeContext != "" {
		b.WriteString(fmt.Sprintf(" [context=%s]", c.CheckKubeContext))
	}

	if c.AskForDryRun {
		b.WriteString(" [dry-run]")
	}
	return b.String()
}

func (c *Command) BuildArgs(extra []string) string {
	var b strings.Builder
	b.WriteString(" ")
	for _, a := range c.Args {
		b.WriteString(" ")
		b.WriteString(a)
	}
	if len(extra) > 0 {
		for _, a := range extra {
			b.WriteString(" ")
			b.WriteString(a)
		}
	}
	return b.String()
}

func (c *Command) RunCommand(extra []string) error {
	args := make([]string, 0, len(c.Args)+len(extra))
	args = append(args, c.Args...)
	args = append(args, extra...)
	cmd := exec.Command(c.HelmBinary, args...)
	cmd.Dir = c.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Command) Run() error {
	var err error
	c.HelmBinary, err = FindHelm()
	if err != nil {
		return ErrHelmNotFound
	}
	fmt.Println(c.String())
	if c.AskForDryRun {
		r := PromptYesNoCancel("Do you want to run with --dry-run before ?")
		if r == PromptCancel {
			return errors.New("Run cancelled")
		}
		if r == PromptYes {
			err := c.RunCommand([]string{"--dry-run"})
			if err != nil {
				return err
			}
			r = PromptYesNoCancel("Do you want to run it for good now ?")
			if r != PromptYes {
				return nil
			}
		}
	}
	return c.RunCommand(nil)
}

func FindHelm() (string, error) {
	path, err := exec.LookPath("helm")
	if err != nil {

		return "", err
	}
	return path, err
}

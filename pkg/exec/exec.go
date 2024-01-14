package exec

import (
	"context"
	"os"
	"os/exec"
	"time"
)

type CmdOpts struct {
	Env        []string
	InheritEnv bool
	Timeout    time.Duration
}

func DefaultCmdOpts() *CmdOpts {
	return &defaultCmdOpts
}

var defaultCmdOpts = CmdOpts{InheritEnv: true}

func RunCommand(args []string, opts *CmdOpts) (string, error) {

	if opts == nil {
		opts = DefaultCmdOpts()
	}

	command := args[0]
	var aa []string
	if len(args) > 1 {
		aa = args[1:]
	}

	var cmd *exec.Cmd

	if opts.Timeout != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, command, aa...)
	} else {
		cmd = exec.Command(command, aa...)
	}

	if opts.InheritEnv {
		cmd.Env = append(os.Environ(), opts.Env...)
	} else {
		cmd.Env = opts.Env
	}

	out, err := cmd.Output()
	return string(out), err
}

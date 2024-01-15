package commands

import (
	"errors"
	"fmt"

	"github.com/grippenet/helman/pkg/helm"
	"github.com/spf13/cobra"
)

const (
	OpResolve = "-config"
	OpCommand = "-show"
	OpRun     = "-run"
)

// Helm Commands

type HelmCommand struct {
	Name    string
	Short   string
	Aliases []string
}

func (helmCmd *HelmCommand) RunE(cmd *cobra.Command, cmdArgs []string) error {
	manager := getManager()

	args := make([]string, 0)
	op := OpRun
	for _, a := range cmdArgs {
		if a == OpResolve || a == OpCommand {
			op = a
		} else {
			args = append(args, a)
		}
	}
	if len(args) < 2 {
		return errors.New("at least target and stage names are expected as argument")
	}

	resolved, err := manager.ResolveCommand(helmCmd.Name, args[0], args[1], args[2:])
	if err != nil {
		return err
	}

	if op == OpResolve {
		resolved.Print()
		return nil
	}

	hc, err := manager.CreateHelmCommand(resolved)

	if err != nil {
		return err
	}

	if op == OpCommand {
		fmt.Println(hc.String())
	}

	if op == OpRun {
		fmt.Println("--Run")
		return hc.Run()
	}
	return nil
}

func (helmCmd *HelmCommand) createCommand() *cobra.Command {

	return &cobra.Command{
		DisableFlagParsing: true,
		Short:              helmCmd.Short,
		Long: `
			-subcommand (prefixed by a '-') : -show,-resolve. If missing uses the '-run'
			target: Name of the target in config
			stage: stage name (aka env)

			Available subcommands:
			-run (default): Run the helm command with resolved command from config
			-show : Dont run the helm command but show it
			-config : Only resolve the target configuration (real file name, arguments, ...)
		`,
		Aliases: helmCmd.Aliases,
		Use:     fmt.Sprintf("%s [-show|-config] target stage (... extra args passed to helm)", helmCmd.Name),

		Args: cobra.MinimumNArgs(2),
		RunE: helmCmd.RunE,
	}
}

var HelmCommandTemplates = []*HelmCommand{
	&HelmCommand{Name: helm.CommandInstall, Short: "`helm install` for the target & stage"},
	&HelmCommand{Name: helm.CommandUpgrade, Short: "`helm upgrade` for a target & stage"},
	&HelmCommand{Name: helm.CommandTemplate, Short: "`helm template` for a target & stage"},
	&HelmCommand{Name: helm.CommandDiff, Short: "`helm diff upgrade` for a target & stage"},
	&HelmCommand{Name: helm.CommandShowValues, Short: "`helm show values` for a target & stage"},
}

//install|upgrade|template|verify|diff

func init() {

	for _, hc := range HelmCommandTemplates {
		rootCmd.AddCommand(hc.createCommand())
	}

}

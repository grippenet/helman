package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Helm Commands

type HelmCommand struct {
	Name    string
	Short   string
	Aliases []string
}

func (helmCmd *HelmCommand) RunE(cmd *cobra.Command, args []string) error {
	manager := getManager()
	hc, err := manager.ChartCommand(helmCmd.Name, args[0], args[1], args[2:])
	if err != nil {
		return err
	}
	return hc.Run()
}

func (helmCmd *HelmCommand) createCommand() *cobra.Command {

	return &cobra.Command{
		DisableFlagParsing: true,
		Short:              helmCmd.Short,
		Long: `
			target: Name of the target in config
			stage: stage name (aka env)
		`,
		Aliases: helmCmd.Aliases,
		Use:     fmt.Sprintf("%s target stage ...", helmCmd.Name),

		Args: cobra.MinimumNArgs(2),
		RunE: helmCmd.RunE,
	}
}

func createCommand(name string, short string) *cobra.Command {
	hc := HelmCommand{Name: name, Short: short}
	return hc.createCommand()
}

//install|upgrade|template|verify|diff

func init() {

	rootCmd.AddCommand(
		createCommand("install", "Run helm install for the target & stage"),
		createCommand("upgrade", "Run helm upgrade for a target & stage"),
		createCommand("template", "Run helm upgrade for a target & stage"),
		createCommand("diff", "Run helm diff upgrade for a target & stage"),
		createCommand("show-values", "Run helm show values for a target & stage"),
	)
}

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package commands

import (
	"fmt"
	"os"

	"github.com/grippenet/helman/pkg/config"
	"github.com/grippenet/helman/pkg/helm"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helman",
	Short: "Helm manager",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func getManager() *helm.HelmManager {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("helman: Unable to load config : %s", err)
		os.Exit(2)
	}
	manager := helm.NewHelmManager(config)
	return manager
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package commands

import "github.com/spf13/cobra"

var ManCmd = &cobra.Command{
	Use:   ".man",
	Short: "Manager commands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	rootCmd.AddCommand(ManCmd)
}

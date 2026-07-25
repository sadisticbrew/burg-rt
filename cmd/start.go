package cmd

import (
	"burg/internal/operations"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <container-id>",
	Short: "Start a created container",
	Long: `Execute the user defined process in a created container.
The container must have been previously created with 'burg create'.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return operations.Start(stateRoot, args[0])
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

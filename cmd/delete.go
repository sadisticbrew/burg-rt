package cmd

import (
	"burg/internal/operations"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <container-id>",
	Short: "Delete a container",
	Long: `Delete any resources held by the container, including the state directory.
The container must be in a stopped state unless --force is used.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		return operations.Delete(stateRoot, args[0], force)
	},
}

func init() {
	deleteCmd.Flags().BoolP("force", "f", false, "forcibly delete the container if it is still running (uses SIGKILL)")

	rootCmd.AddCommand(deleteCmd)
}

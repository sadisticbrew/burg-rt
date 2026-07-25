package cmd

import (
	"burg/internal/operations"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <container-id>",
	Short: "Create and start a container",
	Long: `Create and immediately start a container from an OCI bundle.
This is a convenience command equivalent to 'burg create' followed by 'burg start'.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle, _ := cmd.Flags().GetString("bundle")
		pidFile, _ := cmd.Flags().GetString("pid-file")
		consoleSocket, _ := cmd.Flags().GetString("console-socket")
		detach, _ := cmd.Flags().GetBool("detach")

		return operations.Run(stateRoot, args[0], bundle, detach, pidFile, consoleSocket)
	},
}

func init() {
	runCmd.Flags().StringP("bundle", "b", ".", "path to the root of the OCI bundle directory")
	runCmd.Flags().String("pid-file", "", "file to write the container process id to")
	runCmd.Flags().String("console-socket", "", "path to an AF_UNIX socket for receiving the master pseudo-terminal")
	runCmd.Flags().BoolP("detach", "d", false, "detach from the container's process")

	rootCmd.AddCommand(runCmd)
}

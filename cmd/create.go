package cmd

import (
	"burg/internal/operations"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <container-id>",
	Short: "Create a container",
	Long: `Create an instance of a container from an OCI bundle directory.
The bundle must contain a config.json specification file. The container
will be created in a stopped state, ready to be started with 'burg start'.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle, _ := cmd.Flags().GetString("bundle")

		opts := operations.CreateOptions{
			ID:          args[0],
			Bundle:      bundle,
			Root:        stateRoot,
			Annotations: nil,
		}

		return operations.Create(opts)
	},
}

func init() {
	createCmd.Flags().StringP("bundle", "b", ".", "path to the root of the OCI bundle directory")

	rootCmd.AddCommand(createCmd)
}

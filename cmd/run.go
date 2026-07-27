package cmd

import (
	"burg/internal/operations"
	"log/slog"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <container-id>",
	Short: "Create and start a container",
	Long: `Create and immediately start a container from an OCI bundle.
This is a convenience command equivalent to 'burg create' followed by 'burg start'.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bundle, err := cmd.Flags().GetString("bundle")
		if err != nil {
			return err
		}

		opts := operations.CreateOptions{
			ID:     args[0],
			Bundle: bundle,
		}

		if err := operations.Create(opts); err != nil {
			if err == operations.ErrContainerExists {
				slog.Info("Container already exists. Starting.")
			} else {
				return err
			}
		} else {
			slog.Info("Container created successfully.")
		}

		if err := operations.Start(opts.Root, opts.ID); err != nil {
			return err
		}

		slog.Info("Container started successfully.")
		return nil
	},
}

func init() {
	runCmd.Flags().StringP("bundle", "b", ".", "path to the root of the OCI bundle directory")
	runCmd.Flags().String("pid-file", "", "file to write the container process id to")
	runCmd.Flags().String("console-socket", "", "path to an AF_UNIX socket for receiving the master pseudo-terminal")
	runCmd.Flags().BoolP("detach", "d", false, "detach from the container's process")

	rootCmd.AddCommand(runCmd)
}

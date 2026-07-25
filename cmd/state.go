package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"burg/internal/operations"

	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{
	Use:   "state <container-id>",
	Short: "Output the state of a container",
	Long: `Output the state of a container in JSON format as defined by the OCI
runtime specification. The state includes the container's ID, status,
process ID, bundle path, and annotations.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := operations.GetState(stateRoot, args[0])
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal state: %w", err)
		}

		_, err = fmt.Fprintln(os.Stdout, string(data))
		return err
	},
}

func init() {
	rootCmd.AddCommand(stateCmd)
}

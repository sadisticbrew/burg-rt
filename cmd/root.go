package cmd

import (
	"burg/internal/state"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	stateRoot string
	logPath   string
	logFormat string
	debug     bool
)

var rootCmd = &cobra.Command{
	Use:           "burg",
	Short:         "A minimal OCI container runtime",
	Long:          "burg - A minimal OCI-ish container runtime and eBPF security policy engine built from scratch in Go.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       "0.0.1",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "burg: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	defaultRoot := state.DefaultRoot()

	rootCmd.PersistentFlags().StringVar(&stateRoot, "root", defaultRoot, "root directory for storage of container state")
	rootCmd.PersistentFlags().StringVar(&logPath, "log", "", "set the log file path where internal debug information is written")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", "set the format used by logs ('text' or 'json')")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output for logging")
}

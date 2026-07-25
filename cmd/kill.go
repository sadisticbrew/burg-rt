package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"burg/internal/operations"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var killCmd = &cobra.Command{
	Use:   "kill <container-id> [signal]",
	Short: "Send a signal to the container's init process",
	Long: `Send the specified signal (default: SIGTERM) to the container's init process.
Signals may be specified by name (SIGKILL, KILL) or by number (9).`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		containerID := args[0]
		all, _ := cmd.Flags().GetBool("all")

		sig := unix.SIGTERM
		if len(args) > 1 {
			parsed, err := parseSignal(args[1])
			if err != nil {
				return err
			}
			sig = parsed
		}

		return operations.Kill(stateRoot, containerID, sig, all)
	},
}

func init() {
	killCmd.Flags().BoolP("all", "a", false, "send the specified signal to all processes inside the container")

	rootCmd.AddCommand(killCmd)
}

// parseSignal converts a signal string to a unix.Signal.
// Accepts formats: "SIGKILL", "KILL", "kill", "9".
func parseSignal(raw string) (unix.Signal, error) {
	// Try numeric first
	if num, err := strconv.Atoi(raw); err == nil {
		if num < 1 {
			return 0, fmt.Errorf("invalid signal: %s", raw)
		}
		return unix.Signal(num), nil
	}

	// Normalize: uppercase and ensure SIG prefix
	name := strings.ToUpper(raw)
	if !strings.HasPrefix(name, "SIG") {
		name = "SIG" + name
	}

	sig, ok := signalMap[name]
	if !ok {
		return 0, fmt.Errorf("unknown signal: %q", raw)
	}
	return sig, nil
}

var signalMap = map[string]unix.Signal{
	"SIGABRT":   unix.SIGABRT,
	"SIGALRM":   unix.SIGALRM,
	"SIGBUS":    unix.SIGBUS,
	"SIGCHLD":   unix.SIGCHLD,
	"SIGCONT":   unix.SIGCONT,
	"SIGFPE":    unix.SIGFPE,
	"SIGHUP":    unix.SIGHUP,
	"SIGILL":    unix.SIGILL,
	"SIGINT":    unix.SIGINT,
	"SIGIO":     unix.SIGIO,
	"SIGIOT":    unix.SIGIOT,
	"SIGKILL":   unix.SIGKILL,
	"SIGPIPE":   unix.SIGPIPE,
	"SIGPOLL":   unix.SIGPOLL,
	"SIGPROF":   unix.SIGPROF,
	"SIGPWR":    unix.SIGPWR,
	"SIGQUIT":   unix.SIGQUIT,
	"SIGSEGV":   unix.SIGSEGV,
	"SIGSTKFLT": unix.SIGSTKFLT,
	"SIGSTOP":   unix.SIGSTOP,
	"SIGSYS":    unix.SIGSYS,
	"SIGTERM":   unix.SIGTERM,
	"SIGTRAP":   unix.SIGTRAP,
	"SIGTSTP":   unix.SIGTSTP,
	"SIGTTIN":   unix.SIGTTIN,
	"SIGTTOU":   unix.SIGTTOU,
	"SIGURG":    unix.SIGURG,
	"SIGUSR1":   unix.SIGUSR1,
	"SIGUSR2":   unix.SIGUSR2,
	"SIGVTALRM": unix.SIGVTALRM,
	"SIGWINCH":  unix.SIGWINCH,
	"SIGXCPU":   unix.SIGXCPU,
	"SIGXFSZ":   unix.SIGXFSZ,
}

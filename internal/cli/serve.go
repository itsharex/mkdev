package cli

import "github.com/spf13/cobra"

// newServeCmd is an alias for the TUI. Kept for backward compatibility with
// the original walking-skeleton command. Both `mkdev` (no args) and
// `mkdev serve` launch the same TUI; the proxy runs in-process for as long
// as the TUI is open.
func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Launch the mkdev TUI (alias for the default command)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return launchTUI(cmd)
		},
	}
}

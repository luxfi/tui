// Copyright (C) 2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package cli provides the "tui" subcommand for the lux CLI.
// It launches the interactive terminal UI for monitoring and managing
// the Lux blockchain stack.
package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxfi/tui/app"
	"github.com/spf13/cobra"
)

// NewCmd returns the "tui" command for the lux CLI.
func NewCmd() *cobra.Command {
	var endpoint string

	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Interactive terminal UI for Lux",
		Long: `Launch the interactive terminal UI for monitoring and managing
the Lux blockchain stack.

The TUI provides a dashboard with tabs for:
  - Dashboard:   Network overview and health
  - Nodes:       Node status and management
  - Chains:      Chain status and block heights
  - Validators:  Validator set and uptime
  - Logs:        Live log streaming

NAVIGATION:

  Tab/Arrow   Navigate between tabs
  1-5         Jump to tab by number
  r           Refresh data
  q           Quit

Examples:
  lux tui
  lux tui --endpoint http://localhost:9640`,
		RunE: func(_ *cobra.Command, _ []string) error {
			model := app.NewApp()
			if endpoint != "" {
				model.SetEndpoint(endpoint)
			}

			p := tea.NewProgram(
				model,
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)

			if _, err := p.Run(); err != nil {
				return fmt.Errorf("run TUI: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&endpoint, "endpoint", "", "Lux node API endpoint (default: http://127.0.0.1:9630)")

	// Suppress usage on error -- the TUI handles its own display
	cmd.SilenceUsage = true

	return cmd
}

func init() {
	// Ensure TUI output goes to the right place when embedded in lux CLI
	_ = os.Stderr
}

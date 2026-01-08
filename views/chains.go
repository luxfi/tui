// Copyright (C) 2024-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxfi/tui/styles"
)

// ChainsModel represents the chains view
type ChainsModel struct {
	width    int
	height   int
	chains   []ChainStatus
	selected int
}

// NewChainsModel creates a new chains model
func NewChainsModel() ChainsModel {
	return ChainsModel{}
}

// SetSize sets the view dimensions
func (m *ChainsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// UpdateData updates the chains data
func (m *ChainsModel) UpdateData(chains []ChainStatus) {
	m.chains = chains
	if m.selected >= len(chains) {
		m.selected = len(chains) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

// Update handles messages
func (m ChainsModel) Update(msg tea.Msg) (ChainsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.chains)-1 {
				m.selected++
			}
		}
	}
	return m, nil
}

// View renders the chains view
func (m ChainsModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if len(m.chains) == 0 {
		return styles.CenterText("No chains running", m.width)
	}

	var rows []string

	// Header
	header := fmt.Sprintf("%-15s %-20s %-10s %-12s %-10s %-10s %-10s",
		"Name", "Chain ID", "Type", "Height", "Status", "Validators", "TPS")
	rows = append(rows, styles.TableHeaderStyle.Render(header))

	// Rows
	for i, chain := range m.chains {
		status := styles.FormatStatus(chain.Status)

		row := fmt.Sprintf("%-15s %-20s %-10s %-12d %-10s %-10d %-10.2f",
			styles.Truncate(chain.Name, 15),
			styles.Truncate(chain.ChainID, 20),
			chain.Type,
			chain.Height,
			status,
			chain.Validators,
			chain.TPS)

		if i == m.selected {
			rows = append(rows, styles.SelectedRowStyle.Render(row))
		} else {
			rows = append(rows, styles.TableRowStyle.Render(row))
		}
	}

	table := strings.Join(rows, "\n")

	// Details panel for selected chain
	details := ""
	if m.selected < len(m.chains) {
		chain := m.chains[m.selected]
		details = m.renderChainDetails(chain)
	}

	return styles.SectionTitleStyle.Render("Chains") + "\n" +
		table + "\n\n" +
		details + "\n" +
		styles.FooterStyle.Render("j/k or Up/Down: Navigate | d: Deploy new chain")
}

// renderChainDetails renders details for a selected chain
func (m ChainsModel) renderChainDetails(chain ChainStatus) string {
	details := fmt.Sprintf(`
%s: %s
%s: %s
%s: %s
%s: %s
%s: %s
%s: %s`,
		styles.LabelStyle.Render("Name"), styles.ValueStyle.Render(chain.Name),
		styles.LabelStyle.Render("Chain ID"), styles.ValueStyle.Render(chain.ChainID),
		styles.LabelStyle.Render("Subnet ID"), styles.ValueStyle.Render(chain.SubnetID),
		styles.LabelStyle.Render("Type"), styles.ValueStyle.Render(chain.Type),
		styles.LabelStyle.Render("Height"), styles.ValueStyle.Render(fmt.Sprintf("%d", chain.Height)),
		styles.LabelStyle.Render("TPS"), styles.ValueStyle.Render(fmt.Sprintf("%.2f", chain.TPS)))

	return styles.Box("Chain Details", details, m.width/2)
}

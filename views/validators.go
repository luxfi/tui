// Copyright (C) 2024-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxfi/tui/styles"
)

// ValidatorsModel represents the validators view
type ValidatorsModel struct {
	width      int
	height     int
	validators []ValidatorStatus
	selected   int
}

// NewValidatorsModel creates a new validators model
func NewValidatorsModel() ValidatorsModel {
	return ValidatorsModel{}
}

// SetSize sets the view dimensions
func (m *ValidatorsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// UpdateData updates the validators data
func (m *ValidatorsModel) UpdateData(validators []ValidatorStatus) {
	m.validators = validators
	if m.selected >= len(validators) {
		m.selected = len(validators) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

// Update handles messages
func (m ValidatorsModel) Update(msg tea.Msg) (ValidatorsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.validators)-1 {
				m.selected++
			}
		}
	}
	return m, nil
}

// View renders the validators view
func (m ValidatorsModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if len(m.validators) == 0 {
		return styles.CenterText("No validators", m.width)
	}

	var rows []string

	// Header
	header := fmt.Sprintf("%-50s %-15s %-10s %-12s",
		"Node ID", "Stake", "Uptime", "Connected")
	rows = append(rows, styles.TableHeaderStyle.Render(header))

	// Rows
	for i, v := range m.validators {
		connected := styles.ErrorStyle.Render("No")
		if v.Connected {
			connected = styles.SuccessStyle.Render("Yes")
		}

		row := fmt.Sprintf("%-50s %-15s %-10s %-12s",
			styles.Truncate(v.NodeID, 50),
			v.Stake,
			v.Uptime,
			connected)

		if i == m.selected {
			rows = append(rows, styles.SelectedRowStyle.Render(row))
		} else {
			rows = append(rows, styles.TableRowStyle.Render(row))
		}
	}

	table := strings.Join(rows, "\n")

	return styles.SectionTitleStyle.Render("Validators") + "\n" +
		table + "\n\n" +
		styles.FooterStyle.Render("j/k or Up/Down: Navigate")
}

// Copyright (C) 2024-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxfi/tui/styles"
)

// NodesModel represents the nodes view
type NodesModel struct {
	width    int
	height   int
	nodes    []NodeStatus
	selected int
}

// NewNodesModel creates a new nodes model
func NewNodesModel() NodesModel {
	return NodesModel{}
}

// SetSize sets the view dimensions
func (m *NodesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// UpdateData updates the nodes data
func (m *NodesModel) UpdateData(nodes []NodeStatus) {
	m.nodes = nodes
	if m.selected >= len(nodes) {
		m.selected = len(nodes) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

// Update handles messages
func (m NodesModel) Update(msg tea.Msg) (NodesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.nodes)-1 {
				m.selected++
			}
		}
	}
	return m, nil
}

// View renders the nodes view
func (m NodesModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if len(m.nodes) == 0 {
		return styles.CenterText("No nodes running", m.width)
	}

	var rows []string

	// Header
	header := fmt.Sprintf("%-20s %-20s %-12s %-10s %-10s %-10s",
		"Name", "Node ID", "Status", "Version", "Uptime", "Staking")
	rows = append(rows, styles.TableHeaderStyle.Render(header))

	// Rows
	for i, node := range m.nodes {
		staking := "No"
		if node.Staking {
			staking = styles.SuccessStyle.Render("Yes")
		}

		status := styles.FormatStatus(node.Status)

		row := fmt.Sprintf("%-20s %-20s %-12s %-10s %-10s %-10s",
			styles.Truncate(node.Name, 20),
			styles.Truncate(node.NodeID, 20),
			status,
			node.Version,
			node.Uptime,
			staking)

		if i == m.selected {
			rows = append(rows, styles.SelectedRowStyle.Render(row))
		} else {
			rows = append(rows, styles.TableRowStyle.Render(row))
		}
	}

	table := strings.Join(rows, "\n")

	// Details panel for selected node
	details := ""
	if m.selected < len(m.nodes) {
		node := m.nodes[m.selected]
		details = m.renderNodeDetails(node)
	}

	return styles.SectionTitleStyle.Render("Nodes") + "\n" +
		table + "\n\n" +
		details + "\n" +
		styles.FooterStyle.Render("j/k or Up/Down: Navigate")
}

// renderNodeDetails renders details for a selected node
func (m NodesModel) renderNodeDetails(node NodeStatus) string {
	details := fmt.Sprintf(`
%s: %s
%s: %s
%s: %s
%s: %s
%s: %s`,
		styles.LabelStyle.Render("Name"), styles.ValueStyle.Render(node.Name),
		styles.LabelStyle.Render("Node ID"), styles.ValueStyle.Render(node.NodeID),
		styles.LabelStyle.Render("Status"), styles.FormatStatus(node.Status),
		styles.LabelStyle.Render("Version"), styles.ValueStyle.Render(node.Version),
		styles.LabelStyle.Render("Connected Peers"), styles.ValueStyle.Render(fmt.Sprintf("%d", node.Connected)))

	return styles.Box("Node Details", details, m.width/2)
}

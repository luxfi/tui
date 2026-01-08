// Copyright (C) 2024-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luxfi/tui/styles"
)

// DashboardModel represents the dashboard view
type DashboardModel struct {
	width  int
	height int

	// Summary data
	totalNodes      int
	healthyNodes    int
	totalChains     int
	activeChains    int
	totalValidators int

	// Recent data
	nodes      []NodeStatus
	chains     []ChainStatus
	validators []ValidatorStatus
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel() DashboardModel {
	return DashboardModel{}
}

// SetSize sets the view dimensions
func (m *DashboardModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// UpdateData updates the dashboard with new data
func (m *DashboardModel) UpdateData(nodes []NodeStatus, chains []ChainStatus, validators []ValidatorStatus) {
	m.nodes = nodes
	m.chains = chains
	m.validators = validators

	// Calculate summaries
	m.totalNodes = len(nodes)
	m.healthyNodes = 0
	for _, n := range nodes {
		if n.Status == "healthy" {
			m.healthyNodes++
		}
	}

	m.totalChains = len(chains)
	m.activeChains = 0
	for _, c := range chains {
		if c.Status == "active" {
			m.activeChains++
		}
	}

	m.totalValidators = len(validators)
}

// Update handles messages
func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	return m, nil
}

// View renders the dashboard
func (m DashboardModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var sections []string

	// Summary cards
	sections = append(sections, m.renderSummary())

	// Nodes overview
	sections = append(sections, m.renderNodesOverview())

	// Chains overview
	sections = append(sections, m.renderChainsOverview())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderSummary renders the summary cards
func (m DashboardModel) renderSummary() string {
	cardWidth := (m.width - 10) / 3
	if cardWidth < 20 {
		cardWidth = 20
	}

	card := func(title, value, status string) string {
		titleRendered := styles.LabelStyle.Render(title)
		valueRendered := styles.ValueStyle.Render(value)
		statusRendered := styles.FormatStatus(status)

		content := lipgloss.JoinVertical(
			lipgloss.Center,
			titleRendered,
			valueRendered,
			statusRendered,
		)

		return styles.BoxStyle.
			Width(cardWidth).
			Align(lipgloss.Center).
			Render(content)
	}

	nodesStatus := "healthy"
	if m.healthyNodes < m.totalNodes {
		nodesStatus = "warning"
	}

	chainsStatus := "active"
	if m.activeChains < m.totalChains {
		chainsStatus = "warning"
	}

	cards := lipgloss.JoinHorizontal(
		lipgloss.Top,
		card("Nodes", fmt.Sprintf("%d/%d", m.healthyNodes, m.totalNodes), nodesStatus),
		"  ",
		card("Chains", fmt.Sprintf("%d/%d", m.activeChains, m.totalChains), chainsStatus),
		"  ",
		card("Validators", fmt.Sprintf("%d", m.totalValidators), "active"),
	)

	return styles.SectionTitleStyle.Render("Overview") + "\n" + cards
}

// renderNodesOverview renders the nodes overview
func (m DashboardModel) renderNodesOverview() string {
	if len(m.nodes) == 0 {
		return styles.SectionTitleStyle.Render("Nodes") + "\n" + styles.MutedStyle.Render("No nodes")
	}

	var rows []string
	header := fmt.Sprintf("%-15s %-12s %-10s %-10s",
		"Name", "Status", "Version", "Uptime")
	rows = append(rows, styles.TableHeaderStyle.Render(header))

	for _, node := range m.nodes {
		status := styles.FormatStatus(node.Status)
		row := fmt.Sprintf("%-15s %-12s %-10s %-10s",
			styles.Truncate(node.Name, 15),
			status,
			node.Version,
			node.Uptime)
		rows = append(rows, styles.TableRowStyle.Render(row))
	}

	table := strings.Join(rows, "\n")
	return styles.SectionTitleStyle.Render("Nodes") + "\n" + table
}

// renderChainsOverview renders the chains overview
func (m DashboardModel) renderChainsOverview() string {
	if len(m.chains) == 0 {
		return styles.SectionTitleStyle.Render("Chains") + "\n" + styles.MutedStyle.Render("No chains")
	}

	var rows []string
	header := fmt.Sprintf("%-15s %-10s %-12s %-10s",
		"Name", "Type", "Height", "Status")
	rows = append(rows, styles.TableHeaderStyle.Render(header))

	for _, chain := range m.chains {
		status := styles.FormatStatus(chain.Status)
		row := fmt.Sprintf("%-15s %-10s %-12d %-10s",
			styles.Truncate(chain.Name, 15),
			chain.Type,
			chain.Height,
			status)
		rows = append(rows, styles.TableRowStyle.Render(row))
	}

	table := strings.Join(rows, "\n")
	return styles.SectionTitleStyle.Render("Chains") + "\n" + table
}

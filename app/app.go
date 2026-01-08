// Copyright (C) 2024-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package app provides the main application model for the Lux TUI.
package app

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxfi/tui/styles"
	"github.com/luxfi/tui/views"
)

// Tab represents a tab in the TUI
type Tab int

const (
	TabDashboard Tab = iota
	TabNodes
	TabChains
	TabValidators
	TabLogs
)

// App is the main application model
type App struct {
	// Current active tab
	activeTab Tab

	// Tab names
	tabs []string

	// Window dimensions
	width  int
	height int

	// Loading state
	loading bool
	spinner spinner.Model

	// Error state
	err error

	// View models
	dashboard  views.DashboardModel
	nodes      views.NodesModel
	chains     views.ChainsModel
	validators views.ValidatorsModel
	logs       views.LogsModel

	// Config
	apiEndpoint string

	// Auto-refresh
	refreshInterval time.Duration
}

// NewApp creates a new application model
func NewApp() *App {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	return &App{
		activeTab: TabDashboard,
		tabs: []string{
			"Dashboard",
			"Nodes",
			"Chains",
			"Validators",
			"Logs",
		},
		loading:         true,
		spinner:         s,
		dashboard:       views.NewDashboardModel(),
		nodes:           views.NewNodesModel(),
		chains:          views.NewChainsModel(),
		validators:      views.NewValidatorsModel(),
		logs:            views.NewLogsModel(),
		apiEndpoint:     "http://127.0.0.1:9630",
		refreshInterval: 5 * time.Second,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.spinner.Tick,
		a.fetchData(),
		a.tickRefresh(),
	)
}

// Update handles messages
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit

		case "tab", "right", "l":
			a.activeTab = Tab((int(a.activeTab) + 1) % len(a.tabs))

		case "shift+tab", "left", "h":
			a.activeTab = Tab((int(a.activeTab) - 1 + len(a.tabs)) % len(a.tabs))

		case "1":
			a.activeTab = TabDashboard
		case "2":
			a.activeTab = TabNodes
		case "3":
			a.activeTab = TabChains
		case "4":
			a.activeTab = TabValidators
		case "5":
			a.activeTab = TabLogs

		case "r":
			// Manual refresh
			cmds = append(cmds, a.fetchData())
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.updateViewSizes()

	case spinner.TickMsg:
		var cmd tea.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case refreshMsg:
		cmds = append(cmds, a.fetchData(), a.tickRefresh())

	case dataMsg:
		a.loading = false
		a.updateFromData(msg)

	case errMsg:
		a.loading = false
		a.err = msg.err
	}

	// Update active view
	switch a.activeTab {
	case TabDashboard:
		var cmd tea.Cmd
		a.dashboard, cmd = a.dashboard.Update(msg)
		cmds = append(cmds, cmd)
	case TabNodes:
		var cmd tea.Cmd
		a.nodes, cmd = a.nodes.Update(msg)
		cmds = append(cmds, cmd)
	case TabChains:
		var cmd tea.Cmd
		a.chains, cmd = a.chains.Update(msg)
		cmds = append(cmds, cmd)
	case TabValidators:
		var cmd tea.Cmd
		a.validators, cmd = a.validators.Update(msg)
		cmds = append(cmds, cmd)
	case TabLogs:
		var cmd tea.Cmd
		a.logs, cmd = a.logs.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

// View renders the application
func (a *App) View() string {
	if a.width == 0 {
		return "Loading..."
	}

	// Build the view
	var content string

	// Header with tabs
	header := a.renderHeader()

	// Main content
	if a.loading {
		content = styles.CenterText(a.spinner.View()+" Loading...", a.width)
	} else if a.err != nil {
		content = styles.ErrorStyle.Render("Error: " + a.err.Error())
	} else {
		switch a.activeTab {
		case TabDashboard:
			content = a.dashboard.View()
		case TabNodes:
			content = a.nodes.View()
		case TabChains:
			content = a.chains.View()
		case TabValidators:
			content = a.validators.View()
		case TabLogs:
			content = a.logs.View()
		}
	}

	// Footer
	footer := a.renderFooter()

	// Combine
	return header + "\n" + content + "\n" + footer
}

// renderHeader renders the tab header
func (a *App) renderHeader() string {
	var tabs string
	for i, tab := range a.tabs {
		style := styles.TabStyle
		if Tab(i) == a.activeTab {
			style = styles.ActiveTabStyle
		}
		tabs += style.Render(tab) + " "
	}

	title := styles.TitleStyle.Render("Lux Blockchain Dashboard")
	return title + "\n" + tabs
}

// renderFooter renders the footer
func (a *App) renderFooter() string {
	help := "Tab/Arrow: Navigate | 1-5: Jump to tab | r: Refresh | q: Quit"
	return styles.FooterStyle.Render(help)
}

// updateViewSizes updates the sizes of all views
func (a *App) updateViewSizes() {
	// Reserve space for header and footer
	contentHeight := a.height - 6 // Header (2) + Footer (1) + margins

	a.dashboard.SetSize(a.width, contentHeight)
	a.nodes.SetSize(a.width, contentHeight)
	a.chains.SetSize(a.width, contentHeight)
	a.validators.SetSize(a.width, contentHeight)
	a.logs.SetSize(a.width, contentHeight)
}

// Messages

type refreshMsg struct{}

type dataMsg struct {
	nodes      []views.NodeStatus
	chains     []views.ChainStatus
	validators []views.ValidatorStatus
}

type errMsg struct {
	err error
}

// Commands

func (a *App) tickRefresh() tea.Cmd {
	return tea.Tick(a.refreshInterval, func(time.Time) tea.Msg {
		return refreshMsg{}
	})
}

func (a *App) fetchData() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual data fetching via netrunner API
		// For now, return mock data
		return dataMsg{
			nodes: []views.NodeStatus{
				{Name: "node-1", Status: "healthy", Version: "v1.21.0", Uptime: "24h"},
				{Name: "node-2", Status: "healthy", Version: "v1.21.0", Uptime: "24h"},
				{Name: "node-3", Status: "bootstrapping", Version: "v1.21.0", Uptime: "1h"},
			},
			chains: []views.ChainStatus{
				{Name: "P-Chain", Type: "Platform", Height: 12345, Status: "active"},
				{Name: "X-Chain", Type: "DAG", Height: 67890, Status: "active"},
				{Name: "C-Chain", Type: "EVM", Height: 11111, Status: "active"},
			},
			validators: []views.ValidatorStatus{
				{NodeID: "NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg", Stake: "2000 LUX", Uptime: "99.9%", Connected: true},
				{NodeID: "NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ", Stake: "2000 LUX", Uptime: "99.8%", Connected: true},
			},
		}
	}
}

func (a *App) updateFromData(data dataMsg) {
	a.dashboard.UpdateData(data.nodes, data.chains, data.validators)
	a.nodes.UpdateData(data.nodes)
	a.chains.UpdateData(data.chains)
	a.validators.UpdateData(data.validators)
}

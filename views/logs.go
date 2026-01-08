// Copyright (C) 2024-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luxfi/tui/styles"
)

// LogsModel represents the logs view
type LogsModel struct {
	width    int
	height   int
	viewport viewport.Model
	logs     []LogEntry
	filter   string
	level    string // Filter by log level
}

// NewLogsModel creates a new logs model
func NewLogsModel() LogsModel {
	return LogsModel{
		level: "all",
	}
}

// SetSize sets the view dimensions
func (m *LogsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport = viewport.New(width, height-4)
}

// AddLog adds a log entry
func (m *LogsModel) AddLog(entry LogEntry) {
	m.logs = append(m.logs, entry)
	// Keep last 1000 logs
	if len(m.logs) > 1000 {
		m.logs = m.logs[len(m.logs)-1000:]
	}
	m.updateViewport()
}

// Update handles messages
func (m LogsModel) Update(msg tea.Msg) (LogsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.level = "all"
			m.updateViewport()
		case "2":
			m.level = "error"
			m.updateViewport()
		case "3":
			m.level = "warn"
			m.updateViewport()
		case "4":
			m.level = "info"
			m.updateViewport()
		case "5":
			m.level = "debug"
			m.updateViewport()
		case "c":
			m.logs = nil
			m.updateViewport()
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the logs view
func (m LogsModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Filter bar
	filterBar := m.renderFilterBar()

	// Viewport with logs
	m.updateViewport()

	return styles.SectionTitleStyle.Render("Logs") + "\n" +
		filterBar + "\n" +
		m.viewport.View() + "\n" +
		styles.FooterStyle.Render("1-5: Filter level | c: Clear | PgUp/PgDn: Scroll")
}

// renderFilterBar renders the filter bar
func (m LogsModel) renderFilterBar() string {
	levels := []string{"all", "error", "warn", "info", "debug"}
	var parts []string

	for i, l := range levels {
		style := styles.TabStyle
		if l == m.level {
			style = styles.ActiveTabStyle
		}
		parts = append(parts, style.Render(fmt.Sprintf("%d:%s", i+1, l)))
	}

	return strings.Join(parts, " ")
}

// updateViewport updates the viewport content
func (m *LogsModel) updateViewport() {
	var lines []string

	for _, log := range m.logs {
		// Filter by level
		if m.level != "all" && log.Level != m.level {
			continue
		}

		// Format log line
		levelStyle := m.levelStyle(log.Level)
		line := fmt.Sprintf("%s %s %s: %s",
			styles.MutedStyle.Render(log.Timestamp.Format("15:04:05")),
			levelStyle.Render(fmt.Sprintf("%-5s", log.Level)),
			styles.LabelStyle.Render(log.Source),
			log.Message)
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		lines = append(lines, styles.MutedStyle.Render("No logs to display"))
	}

	m.viewport.SetContent(strings.Join(lines, "\n"))
	m.viewport.GotoBottom()
}

// levelStyle returns the style for a log level
func (m LogsModel) levelStyle(level string) styles.LogLevelStyle {
	switch level {
	case "error":
		return styles.LogLevelStyle{Style: styles.ErrorStyle}
	case "warn":
		return styles.LogLevelStyle{Style: styles.WarningStyle}
	case "info":
		return styles.LogLevelStyle{Style: styles.SuccessStyle}
	case "debug":
		return styles.LogLevelStyle{Style: styles.LabelStyle}
	default:
		return styles.LogLevelStyle{Style: styles.LabelStyle}
	}
}

// Mock some logs for demonstration
func (m *LogsModel) AddMockLogs() {
	now := time.Now()
	mockLogs := []LogEntry{
		{Timestamp: now.Add(-10 * time.Second), Level: "info", Source: "node-1", Message: "Node started successfully"},
		{Timestamp: now.Add(-8 * time.Second), Level: "info", Source: "node-2", Message: "Connected to network"},
		{Timestamp: now.Add(-6 * time.Second), Level: "warn", Source: "node-3", Message: "Low peer count: 5"},
		{Timestamp: now.Add(-4 * time.Second), Level: "info", Source: "chain", Message: "Block 12345 accepted"},
		{Timestamp: now.Add(-2 * time.Second), Level: "debug", Source: "consensus", Message: "Vote received from validator"},
	}

	for _, log := range mockLogs {
		m.AddLog(log)
	}
}

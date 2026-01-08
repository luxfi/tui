// Copyright (C) 2024-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package styles provides styling for the Lux TUI using lipgloss.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#3B82F6") // Blue
	Success   = lipgloss.Color("#10B981") // Green
	Warning   = lipgloss.Color("#F59E0B") // Amber
	Error     = lipgloss.Color("#EF4444") // Red
	Muted     = lipgloss.Color("#6B7280") // Gray
	Text      = lipgloss.Color("#F9FAFB") // Light text
	Border    = lipgloss.Color("#374151") // Dark gray
)

// Component styles
var (
	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	// Tab styles
	TabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(Muted)

	ActiveTabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(Text).
			Background(Primary).
			Bold(true)

	// Footer style
	FooterStyle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	// Warning style
	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	// Spinner style
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(Primary)

	// Box style
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1, 2)

	// Table header style
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Primary).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(Border)

	// Table row style
	TableRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Selected row style
	SelectedRowStyle = lipgloss.NewStyle().
				Background(Primary).
				Foreground(Text).
				Bold(true)

	// Status styles
	StatusHealthy = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	StatusUnhealthy = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	StatusPending = lipgloss.NewStyle().
			Foreground(Warning)

	// Label style
	LabelStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// Value style
	ValueStyle = lipgloss.NewStyle().
			Foreground(Text).
			Bold(true)

	// Section title style
	SectionTitleStyle = lipgloss.NewStyle().
				Foreground(Secondary).
				Bold(true).
				MarginTop(1).
				MarginBottom(1)

	// MutedStyle for subdued text rendering
	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)
)

// Helper functions

// CenterText centers text within a given width
func CenterText(text string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(text)
}

// StatusStyle returns the appropriate style for a status string
func StatusStyle(status string) lipgloss.Style {
	switch status {
	case "healthy", "active", "connected":
		return StatusHealthy
	case "unhealthy", "error", "disconnected":
		return StatusUnhealthy
	default:
		return StatusPending
	}
}

// FormatStatus formats a status with appropriate styling
func FormatStatus(status string) string {
	return StatusStyle(status).Render(status)
}

// Truncate truncates a string to a maximum length
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// LogLevelStyle is a wrapper for log level styling
type LogLevelStyle struct {
	Style lipgloss.Style
}

// Render renders text with this style
func (l LogLevelStyle) Render(text string) string {
	return l.Style.Render(text)
}

// Box creates a boxed content
func Box(title, content string, width int) string {
	titleRendered := SectionTitleStyle.Render(title)
	boxContent := lipgloss.NewStyle().Width(width - 4).Render(content)
	return BoxStyle.Width(width).Render(titleRendered + "\n" + boxContent)
}

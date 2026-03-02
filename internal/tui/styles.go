package tui

import "github.com/charmbracelet/lipgloss"

const (
	headerHeight   = 1
	composerHeight = 3 // border top + input line + border bottom
)

var (
	// Header: "◆ Luna  v0.1.0"
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")). // purple
			Padding(0, 1)

	// Version badge in header
	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")). // muted
			Padding(0, 1)

	// Role label: "You" / "Luna"
	roleLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

	// Divider line under role label
	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	// Composer border
	composerStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)

	// Spinner label style
	spinnerLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))
)

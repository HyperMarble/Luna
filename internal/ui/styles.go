package ui

import "github.com/charmbracelet/lipgloss"

var (
	logoStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("81")).
			Bold(true)

	chatStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	chatTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")).
			Bold(true)

	secondaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))

	composerStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69")).
			Padding(0, 1)

	compactNoticeStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("52"))
)

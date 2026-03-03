package tui

import "github.com/charmbracelet/lipgloss"

const headerHeight = 1

var (
	// Header: "◆ Luna  v0.1.0"
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Padding(0, 1)

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)

	// User message: "> message" pill
	userPillStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	userPromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// Luna response: "● text"
	responseBulletStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("208")) // orange like Claude Code

	responseTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	// Thinking: "* Word…"
	thinkingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Italic(true)

	// Composer divider
	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	// Input prompt "> "
	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))

	// Welcome screen
	mascotStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")) // white

	welcomeTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("255")).
				Padding(0, 1)

	welcomeSubStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	welcomePathStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	// Slash command picker
	pickerCmdStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99"))

	pickerSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("99"))

	pickerDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	pickerSelectedDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))
)

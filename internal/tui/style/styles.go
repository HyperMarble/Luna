package style

import "github.com/charmbracelet/lipgloss"

var Mascot = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

var WelcomeTitle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("255")).
	Padding(0, 1)

var WelcomeVersion = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("250"))

var WelcomeSub = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

var WelcomePath = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

var UserPill = lipgloss.NewStyle().
	Foreground(lipgloss.Color("252")).
	Background(lipgloss.Color("236")).
	Padding(0, 1)

var ResponseBullet = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))

var ResponseText = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

var Thinking = lipgloss.NewStyle().
	Foreground(lipgloss.Color("208")).
	Italic(true)

var Divider = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

var InputPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

var PickerCmd = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))

var PickerDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

var PickerSelected = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))

var PickerSelectedDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))

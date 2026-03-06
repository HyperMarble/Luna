package style

import "github.com/charmbracelet/lipgloss"

// Indian flag saffron (#FF9933) for the top blocks, green (#138808) for bottom.
var MascotTop = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9933"))
var MascotMid = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
var MascotBot = lipgloss.NewStyle().Foreground(lipgloss.Color("#138808"))

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

var PickerCmd = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9933"))

var PickerDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

var PickerSelected = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF9933"))

var PickerSelectedDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9933"))

// BadgeFree is shown next to free-tier providers in the model picker.
var BadgeFree = lipgloss.NewStyle().Foreground(lipgloss.Color("#138808"))

// BadgeLocked is shown next to paid providers without a saved API key.
var BadgeLocked = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

// BadgeUnlocked is shown next to paid providers that have a saved API key.
var BadgeUnlocked = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9933"))

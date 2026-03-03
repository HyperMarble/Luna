package tui

import "github.com/charmbracelet/lipgloss"

// ── Header ───────────────────────────────────────────────────────────────────

var headerStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("99")). // purple
	Padding(0, 1)

var versionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")). // muted grey
	Padding(0, 1)

// ── Welcome box ──────────────────────────────────────────────────────────────

var mascotStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("255")) // white

var welcomeTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("255")).
	Padding(0, 1)

var welcomeSubStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241"))

var welcomePathStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("238"))

// ── Messages ─────────────────────────────────────────────────────────────────

// userPillStyle wraps user messages in a grey pill: "> message"
var userPillStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("252")).
	Background(lipgloss.Color("236")).
	Padding(0, 1)

// responseBulletStyle colours the "●" orange, matching Claude Code's style.
var responseBulletStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("208"))

var responseTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("252"))

// thinkingStyle renders the animated "* Thinking…" indicator in orange italic.
var thinkingStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("208")).
	Italic(true)

// ── Composer (input area) ────────────────────────────────────────────────────

var dividerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("238"))

var inputPromptStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241"))

// ── Slash command picker ─────────────────────────────────────────────────────

var pickerCmdStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("99"))

var pickerDescStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241"))

var pickerSelectedStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("99"))

var pickerSelectedDescStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("99"))

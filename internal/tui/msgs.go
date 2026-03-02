package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// UserSubmitMsg is fired when the user presses Enter in the composer.
type UserSubmitMsg struct {
	Text string
}

// LunaStubMsg is a placeholder Luna response (no agent yet).
type LunaStubMsg struct {
	Text string
}

// SpinnerTickMsg drives the spinner animation.
type SpinnerTickMsg spinner.TickMsg

// windowSizeMsg is sent by Bubble Tea on terminal resize.
// We alias tea.WindowSizeMsg for clarity inside Update.
type windowSizeMsg = tea.WindowSizeMsg

package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Message represents a single conversation turn.
type Message struct {
	Role    string // "user" | "luna"
	Content string // raw text (markdown for luna)
}

// Model is the entire application state.
type Model struct {
	width    int
	height   int
	viewport viewport.Model
	input    textinput.Model
	spinner  spinner.Model
	messages []Message
	thinking bool // true while waiting for Luna response
	ready    bool // true once viewport is initialized with real dimensions
}

// NewModel returns the initial model.
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Ask Luna..."
	ti.Focus()
	ti.CharLimit = 2000

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return Model{
		input:   ti,
		spinner: sp,
	}
}

// Init starts the text input blink cursor.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Input returns the text input model (for testing).
func (m Model) Input() textinput.Model { return m.input }

// Messages returns the conversation history (for testing).
func (m Model) Messages() []Message { return m.messages }

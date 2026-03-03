package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var thinkingVerbs = []string{
	"Thinking", "Analyzing", "Processing", "Computing",
	"Reasoning", "Calculating", "Reviewing", "Working",
}

// Message represents a single conversation turn.
type Message struct {
	Role    string // "user" | "luna"
	Content string
}

// Model is the entire application state.
type Model struct {
	width       int
	input       textinput.Model
	spinner     spinner.Model
	messages    []Message
	thinking  bool
	verbIdx   int
	pickerIdx int // selected index in slash command picker
}

// NewModel returns the initial model.
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Ask Luna..."
	ti.Focus()
	ti.CharLimit = 2000
	ti.Prompt = ""

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return Model{
		input:   ti,
		spinner: sp,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Input() textinput.Model { return m.input }
func (m Model) Messages() []Message    { return m.messages }

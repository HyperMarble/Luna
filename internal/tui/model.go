package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// thinkingVerbs cycles through action words shown while Luna is processing a response.
var thinkingVerbs = []string{
	"Thinking", "Analyzing", "Processing", "Computing",
	"Reasoning", "Calculating", "Reviewing", "Working",
}

// Message is a single turn in the conversation.
type Message struct {
	Role    string // "user" | "luna"
	Content string
}

// Model holds the full application state managed by the Bubble Tea runtime.
type Model struct {
	width     int             // terminal width, updated on resize
	height    int             // terminal height, updated on resize
	viewport  viewport.Model  // scrollable conversation area
	ready     bool            // true once viewport is initialised with real dimensions
	input     textinput.Model // the composer text field
	spinner   spinner.Model   // dot spinner shown while thinking
	messages  []Message       // conversation history
	thinking  bool            // true while waiting for Luna's response
	verbIdx   int             // index into thinkingVerbs for the thinking animation
	pickerIdx int             // highlighted row in the slash command picker
}

// NewModel returns the initial application state with a focused input and dot spinner.
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

// Init starts the cursor blink when the program launches.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Input exposes the text input (used in tests).
func (m Model) Input() textinput.Model { return m.input }

// Messages exposes the conversation history (used in tests).
func (m Model) Messages() []Message { return m.messages }

package model

import (
	"context"
	"strings"

	"github.com/HyperMarble/Luna/internal/agent"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/HyperMarble/Luna/internal/tui/events"
	tuilayout "github.com/HyperMarble/Luna/internal/tui/layout"
	"github.com/HyperMarble/Luna/internal/tui/slash"
	"github.com/HyperMarble/Luna/internal/tui/types"
	"github.com/HyperMarble/Luna/internal/tui/view"
)

// UI is the main Bubble Tea model and state owner.
type UI struct {
	width  int
	height int
	input  textinput.Model
	layout tuilayout.UI

	spinner   spinner.Model
	messages  []types.Message
	thinking  bool
	verbIdx   int
	pickerIdx int
	agent     agent.Service
}

// New returns the initial UI model.
func New() UI {
	ti := textinput.New()
	ti.Placeholder = "Ask Luna..."
	ti.Focus()
	ti.CharLimit = 2000
	ti.Prompt = ""

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return UI{
		input:   ti,
		spinner: sp,
		layout:  tuilayout.Compute(80),
		agent:   agent.New(nil),
	}
}

// Init starts the cursor blink when the program launches.
func (m UI) Init() tea.Cmd { return textinput.Blink }

// Input exposes the text input (used in tests).
func (m UI) Input() textinput.Model { return m.input }

// Messages exposes the conversation history (used in tests).
func (m UI) Messages() []types.Message { return m.messages }

// Update routes all Bubble Tea messages and mutates model state.
func (m UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout = tuilayout.Compute(msg.Width)
		m.input.Width = max(10, m.layout.ComposerWidth-4)
		cmds = append(cmds, tea.ClearScreen)

	case tea.KeyMsg:
		cmd, done := m.onKey(msg)
		if done {
			return m, cmd
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case events.UserSubmitMsg:
		m.messages = append(m.messages, types.Message{Role: "user", Content: msg.Text})
		m.thinking = true

	case events.AgentResponseMsg:
		m.thinking = false
		m.messages = append(m.messages, types.Message{Role: "luna", Content: msg.Text})

	case spinner.TickMsg:
		if m.thinking {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			m.verbIdx++
			cmds = append(cmds, spinCmd)
		}
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)

	return m, tea.Batch(append(cmds, inputCmd)...)
}

// View renders the full UI.
func (m UI) View() string {
	return view.Render(view.State{
		Width:       m.width,
		Height:      m.height,
		Layout:      m.layout,
		Input:       m.input,
		Messages:    m.messages,
		Thinking:    m.thinking,
		VerbIdx:     m.verbIdx,
		PickerIndex: m.pickerIdx,
	})
}

func (m *UI) onKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c":
		return tea.Quit, true
	case "up", "down":
		m.movePicker(msg.String())
	case "tab":
		m.completePicker()
	case "esc":
		m.dismissPicker()
	case "enter":
		return m.onEnter()
	}

	if msg.String() != "up" && msg.String() != "down" {
		m.pickerIdx = 0
	}
	return nil, false
}

func (m *UI) movePicker(dir string) {
	if !strings.HasPrefix(m.input.Value(), "/") {
		return
	}
	filtered := slash.Filtered(m.input.Value())
	if len(filtered) == 0 {
		return
	}
	if dir == "up" {
		m.pickerIdx = max(0, m.pickerIdx-1)
	} else {
		m.pickerIdx = min(len(filtered)-1, m.pickerIdx+1)
	}
}

func (m *UI) completePicker() {
	if !strings.HasPrefix(m.input.Value(), "/") {
		return
	}
	filtered := slash.Filtered(m.input.Value())
	if len(filtered) > 0 {
		m.input.SetValue(filtered[m.pickerIdx].Name)
		m.input.CursorEnd()
	}
}

func (m *UI) dismissPicker() {
	if strings.HasPrefix(m.input.Value(), "/") {
		m.input.SetValue("")
		m.pickerIdx = 0
	}
}

func (m *UI) onEnter() (tea.Cmd, bool) {
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		return nil, false
	}
	if strings.HasPrefix(text, "/") {
		return m.executeSlash(text)
	}
	return m.submitText(text), false
}

func (m *UI) submitText(text string) tea.Cmd {
	m.input.SetValue("")
	m.messages = append(m.messages, types.Message{Role: "user", Content: text})
	m.thinking = true
	return tea.Batch(agentResponseCmd(m.agent, text), m.spinner.Tick)
}

func (m *UI) executeSlash(text string) (tea.Cmd, bool) {
	filtered := slash.Filtered(text)
	if len(filtered) > 0 && text != filtered[m.pickerIdx].Name {
		text = filtered[m.pickerIdx].Name
	}
	m.input.SetValue("")
	m.pickerIdx = 0
	return m.handleSlash(text), true
}

func (m *UI) handleSlash(cmd string) tea.Cmd {
	switch cmd {
	case "/exit":
		m.messages = append(m.messages, types.Message{Role: "user", Content: "/exit"})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: "Goodbye! Thanks for using Luna. 🌙"})
		return tea.Quit
	case "/clear":
		m.thinking = false
		m.messages = []types.Message{{Role: "luna", Content: "Conversation cleared."}}
		return nil
	case "/help":
		m.messages = append(m.messages, types.Message{Role: "user", Content: "/help"})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: helpText()})
		return nil
	case "/model":
		m.messages = append(m.messages, types.Message{Role: "user", Content: "/model"})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: "**Model:** claude-sonnet-4-6\n\nModel switching coming soon."})
		return nil
	case "/plugins":
		m.messages = append(m.messages, types.Message{Role: "user", Content: "/plugins"})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: "No plugins installed yet.\n\nPlugin system coming soon."})
		return nil
	default:
		m.messages = append(m.messages, types.Message{Role: "user", Content: cmd})
		m.messages = append(m.messages, types.Message{Role: "luna", Content: "Unknown command: `" + cmd + "`\n\nType `/help` to see available commands."})
		return nil
	}
}

func helpText() string {
	return `**Available commands**

| Command | Description |
|---------|-------------|
| ` + "`/help`" + ` | Show this help |
| ` + "`/clear`" + ` | Clear the conversation |
| ` + "`/model`" + ` | Show the current model |
| ` + "`/plugins`" + ` | Manage plugins |
| ` + "`/exit`" + ` | Exit Luna |`
}

func agentResponseCmd(svc agent.Service, text string) tea.Cmd {
	return func() tea.Msg {
		resp, err := svc.Run(context.Background(), agent.Request{Prompt: text})
		if err != nil {
			return events.AgentResponseMsg{Text: "Agent error: " + err.Error()}
		}
		return events.AgentResponseMsg{Text: resp.Text}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

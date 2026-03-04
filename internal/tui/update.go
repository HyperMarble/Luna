package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Main dispatcher ───────────────────────────────────────────────────────────

// Update is the Bubble Tea message handler — routes each message type to its handler.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - 4

		// Fixed chrome: header + composer only. Welcome box lives inside the viewport.
		headerH := lipgloss.Height(renderHeader(msg.Width))
		const composerH = 2 // divider line + input line (base, picker excluded)
		vpH := max(1, msg.Height-headerH-composerH)

		if !m.ready {
			m.viewport = viewport.New(msg.Width, vpH)
			m.viewport.YPosition = headerH
			m.viewport.SetContent(viewportContent(m.messages, m.thinking, m.verbIdx, vpH, msg.Width))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = vpH
		}

	case tea.KeyMsg:
		// onKey may return early (e.g. ctrl+c, slash commands).
		newM, cmd, done := m.onKey(msg)
		if done {
			return newM, cmd
		}
		m = newM
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case UserSubmitMsg:
		m.messages = append(m.messages, Message{Role: "user", Content: msg.Text})
		m.thinking = true
		if m.ready {
			m.viewport.SetContent(viewportContent(m.messages, m.thinking, m.verbIdx, m.viewport.Height, m.width))
			m.viewport.GotoBottom()
		}

	case LunaStubMsg:
		m.thinking = false
		m.messages = append(m.messages, Message{Role: "luna", Content: msg.Text})
		if m.ready {
			m.viewport.SetContent(viewportContent(m.messages, m.thinking, m.verbIdx, m.viewport.Height, m.width))
			m.viewport.GotoBottom()
		}

	case spinner.TickMsg:
		if m.thinking {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			m.verbIdx++
			if m.ready {
				m.viewport.SetContent(viewportContent(m.messages, m.thinking, m.verbIdx, m.viewport.Height, m.width))
			}
			cmds = append(cmds, spinCmd)
		}
	}

	// Always refresh the text input so cursor blink and typing work.
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)

	// Pass events to the viewport so mouse wheel / PgUp / PgDn scrolling work.
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)

	return m, tea.Batch(append(cmds, inputCmd, vpCmd)...)
}

// ── Keyboard handling ─────────────────────────────────────────────────────────

// onKey routes keyboard input to the appropriate sub-handler.
// Returns done=true when the caller should return immediately
// (quit or slash command — both manage their own final state).
func (m Model) onKey(msg tea.KeyMsg) (Model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit, true

	case "up", "down":
		m = m.movePicker(msg.String())

	case "tab":
		m = m.completePicker()

	case "esc":
		m = m.dismissPicker()

	case "enter":
		return m.onEnter()
	}

	// Reset picker highlight on every key except navigation.
	if msg.String() != "up" && msg.String() != "down" {
		m.pickerIdx = 0
	}

	return m, nil, false
}

// ── Slash picker navigation ───────────────────────────────────────────────────

// movePicker moves the picker highlight up or down when the slash menu is open.
func (m Model) movePicker(dir string) Model {
	if !strings.HasPrefix(m.input.Value(), "/") {
		return m
	}
	filtered := filteredCommands(m.input.Value())
	if len(filtered) == 0 {
		return m
	}
	if dir == "up" {
		m.pickerIdx = max(0, m.pickerIdx-1)
	} else {
		m.pickerIdx = min(len(filtered)-1, m.pickerIdx+1)
	}
	return m
}

// completePicker fills the input with the currently highlighted command on Tab.
func (m Model) completePicker() Model {
	if !strings.HasPrefix(m.input.Value(), "/") {
		return m
	}
	filtered := filteredCommands(m.input.Value())
	if len(filtered) > 0 {
		m.input.SetValue(filtered[m.pickerIdx].Name)
		m.input.CursorEnd()
	}
	return m
}

// dismissPicker clears the input and closes the picker on Esc.
func (m Model) dismissPicker() Model {
	if strings.HasPrefix(m.input.Value(), "/") {
		m.input.SetValue("")
		m.pickerIdx = 0
	}
	return m
}

// ── Enter key ─────────────────────────────────────────────────────────────────

// onEnter handles the Enter key: submits a message or executes a slash command.
func (m Model) onEnter() (Model, tea.Cmd, bool) {
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		return m, nil, false
	}
	if strings.HasPrefix(text, "/") {
		return m.executeSlash(text)
	}
	return m.submitText(text)
}

// submitText appends the user message and queues the stub response + spinner.
func (m Model) submitText(text string) (Model, tea.Cmd, bool) {
	m.input.SetValue("")
	m.messages = append(m.messages, Message{Role: "user", Content: text})
	m.thinking = true
	if m.ready {
		m.viewport.SetContent(viewportContent(m.messages, m.thinking, m.verbIdx, m.viewport.Height, m.width))
		m.viewport.GotoBottom()
	}
	return m, tea.Batch(stubResponseCmd(text), m.spinner.Tick), false
}

// executeSlash resolves the selected command from the picker and dispatches it.
// Always returns done=true — slash commands manage the input state themselves.
func (m Model) executeSlash(text string) (Model, tea.Cmd, bool) {
	filtered := filteredCommands(text)
	if len(filtered) > 0 && text != filtered[m.pickerIdx].Name {
		text = filtered[m.pickerIdx].Name // complete partial input to selected
	}
	m.input.SetValue("")
	m.pickerIdx = 0
	newM, cmd := m.handleSlashCommand(text)
	if newM.ready {
		newM.viewport.SetContent(viewportContent(newM.messages, newM.thinking, newM.verbIdx, newM.viewport.Height, newM.width))
		newM.viewport.GotoBottom()
	}
	return newM, cmd, true
}

// ── Slash command router ──────────────────────────────────────────────────────

// handleSlashCommand routes a fully resolved slash command to its implementation.
func (m Model) handleSlashCommand(cmd string) (Model, tea.Cmd) {
	switch cmd {
	case "/exit":
		return m.cmdExit()
	case "/clear":
		return m.cmdClear(), nil
	case "/help":
		return m.cmdHelp(), nil
	case "/model":
		return m.cmdModel(), nil
	case "/plugins":
		return m.cmdPlugins(), nil
	default:
		return m.cmdUnknown(cmd), nil
	}
}

// ── Slash command implementations ─────────────────────────────────────────────

// cmdExit shows a goodbye message then quits the program.
func (m Model) cmdExit() (Model, tea.Cmd) {
	m.messages = append(m.messages, Message{Role: "user", Content: "/exit"})
	m.messages = append(m.messages, Message{Role: "luna", Content: "Goodbye! Thanks for using Luna. 🌙"})
	return m, tea.Quit
}

// cmdClear resets the conversation to an empty state.
func (m Model) cmdClear() Model {
	m.thinking = false
	m.messages = []Message{{Role: "luna", Content: "Conversation cleared."}}
	return m
}

// cmdHelp appends the help table to the conversation.
func (m Model) cmdHelp() Model {
	m.messages = append(m.messages, Message{Role: "user", Content: "/help"})
	m.messages = append(m.messages, Message{Role: "luna", Content: helpText()})
	return m
}

// cmdModel shows the currently active AI model.
func (m Model) cmdModel() Model {
	m.messages = append(m.messages, Message{Role: "user", Content: "/model"})
	m.messages = append(m.messages, Message{Role: "luna", Content: "**Model:** claude-sonnet-4-6\n\nModel switching coming soon."})
	return m
}

// cmdPlugins shows the plugin status.
func (m Model) cmdPlugins() Model {
	m.messages = append(m.messages, Message{Role: "user", Content: "/plugins"})
	m.messages = append(m.messages, Message{Role: "luna", Content: "No plugins installed yet.\n\nPlugin system coming soon."})
	return m
}

// cmdUnknown handles any slash command that isn't registered.
func (m Model) cmdUnknown(cmd string) Model {
	m.messages = append(m.messages, Message{Role: "user", Content: cmd})
	m.messages = append(m.messages, Message{Role: "luna", Content: "Unknown command: `" + cmd + "`\n\nType `/help` to see available commands."})
	return m
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// helpText returns the markdown table rendered by /help.
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

// stubResponseCmd returns a command that immediately fires a stub Luna response.
// Replace this with the real Claude API call once the agent is wired in.
func stubResponseCmd(_ string) tea.Cmd {
	return func() tea.Msg {
		return LunaStubMsg{Text: "I'm Luna. Agent coming soon."}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

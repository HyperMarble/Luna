package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.input.Width = msg.Width - 4

	case tea.KeyMsg:
		input := m.input.Value()
		showingPicker := strings.HasPrefix(input, "/")

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "up":
			if showingPicker {
				filtered := filteredCommands(input)
				if len(filtered) > 0 {
					m.pickerIdx = max(0, m.pickerIdx-1)
				}
				break
			}

		case "down":
			if showingPicker {
				filtered := filteredCommands(input)
				if len(filtered) > 0 {
					m.pickerIdx = min(len(filtered)-1, m.pickerIdx+1)
				}
				break
			}

		case "tab":
			if showingPicker {
				filtered := filteredCommands(input)
				if len(filtered) > 0 {
					m.input.SetValue(filtered[m.pickerIdx].Name)
					m.input.CursorEnd()
				}
				break
			}

		case "esc":
			if showingPicker {
				m.input.SetValue("")
				m.pickerIdx = 0
				break
			}

		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" {
				break
			}

			// If picker is showing, select highlighted command
			if strings.HasPrefix(text, "/") {
				filtered := filteredCommands(text)
				if len(filtered) > 0 && text != filtered[m.pickerIdx].Name {
					text = filtered[m.pickerIdx].Name
				}
				m.input.SetValue("")
				m.pickerIdx = 0
				return m.handleSlashCommand(text)
			}

			m.input.SetValue("")
			m.messages = append(m.messages, Message{Role: "user", Content: text})
			m.thinking = true
			cmds = append(cmds, stubResponseCmd(text), m.spinner.Tick)
		}

		// Reset picker index when input changes
		if msg.String() != "up" && msg.String() != "down" {
			m.pickerIdx = 0
		}

	case UserSubmitMsg:
		m.messages = append(m.messages, Message{Role: "user", Content: msg.Text})
		m.thinking = true

	case LunaStubMsg:
		m.thinking = false
		m.messages = append(m.messages, Message{Role: "luna", Content: msg.Text})

	case SpinnerTickMsg:
		if m.thinking {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			m.verbIdx++
			cmds = append(cmds, spinCmd)
		}
	}

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) handleSlashCommand(cmd string) (tea.Model, tea.Cmd) {
	switch cmd {
	case "/exit":
		m.messages = append(m.messages, Message{Role: "user", Content: "/exit"})
		m.messages = append(m.messages, Message{Role: "luna", Content: "Goodbye! Thanks for using Luna. 🌙"})
		return m, tea.Quit

	case "/clear":
		m.thinking = false
		m.messages = []Message{
			{Role: "luna", Content: "Conversation cleared."},
		}

	case "/help":
		m.messages = append(m.messages, Message{Role: "user", Content: "/help"})
		m.messages = append(m.messages, Message{Role: "luna", Content: helpText()})

	case "/model":
		m.messages = append(m.messages, Message{Role: "user", Content: "/model"})
		m.messages = append(m.messages, Message{Role: "luna", Content: "**Model:** claude-sonnet-4-6\n\nModel switching coming soon."})

	case "/plugins":
		m.messages = append(m.messages, Message{Role: "user", Content: "/plugins"})
		m.messages = append(m.messages, Message{Role: "luna", Content: "No plugins installed yet.\n\nPlugin system coming soon."})

	default:
		m.messages = append(m.messages, Message{Role: "user", Content: cmd})
		m.messages = append(m.messages, Message{Role: "luna", Content: "Unknown command: `" + cmd + "`\n\nType `/help` to see available commands."})
	}

	return m, nil
}

func helpText() string {
	return `**Available commands**

| Command | Description |
|---------|-------------|
| ` + "`/help`" + ` | Show this help |
| ` + "`/clear`" + ` | Clear conversation |
| ` + "`/model`" + ` | Show current model |
| ` + "`/plugins`" + ` | Manage plugins |
| ` + "`/exit`" + ` | Exit Luna |`
}

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

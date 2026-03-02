package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - 6 // account for composer padding + border

		viewportHeight := msg.Height - headerHeight - composerHeight - 2
		if !m.ready {
			m.viewport = viewport.New(msg.Width, viewportHeight)
			m.viewport.SetContent("")
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = viewportHeight
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			text := m.input.Value()
			if text == "" {
				break
			}
			m.messages = append(m.messages, Message{Role: "user", Content: text})
			m.input.SetValue("")
			m.thinking = true
			m.viewport.SetContent(renderMessages(m.messages, m.thinking, m.spinner.View(), m.width))
			m.viewport.GotoBottom()
			cmds = append(cmds, stubResponseCmd(text), m.spinner.Tick)

		}

	case UserSubmitMsg:
		m.messages = append(m.messages, Message{Role: "user", Content: msg.Text})
		m.thinking = true

	case LunaStubMsg:
		m.thinking = false
		m.messages = append(m.messages, Message{Role: "luna", Content: msg.Text})
		m.viewport.SetContent(renderMessages(m.messages, false, "", m.width))
		m.viewport.GotoBottom()

	case SpinnerTickMsg:
		if m.thinking {
			var spinCmd tea.Cmd
			m.spinner, spinCmd = m.spinner.Update(msg)
			m.viewport.SetContent(renderMessages(m.messages, m.thinking, m.spinner.View(), m.width))
			cmds = append(cmds, spinCmd)
		}
	}

	// Always update input
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	// Always update viewport
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

// stubResponseCmd fires a fake Luna response after receiving input.
func stubResponseCmd(_ string) tea.Cmd {
	return func() tea.Msg {
		return LunaStubMsg{Text: "I'm Luna. Agent coming soon."}
	}
}

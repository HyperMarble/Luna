package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = max(24, msg.Width-6)
		return m, inputCmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.input.Value() == "" {
				return m, tea.Quit
			}
		case "enter":
			prompt := m.input.Value()
			if prompt != "" {
				m.appendUserPrompt(prompt)
				m.input.SetValue("")
			}
			return m, inputCmd
		}
	}

	return m, inputCmd
}

func (m *Model) appendUserPrompt(prompt string) {
	m.conversation = append(m.conversation, "You: "+prompt)
	m.conversation = append(m.conversation, "Luna: Got it.")
}

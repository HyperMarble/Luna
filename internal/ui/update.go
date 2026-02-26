package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update spinner
	m.spinner, cmd = m.spinner.Update(msg)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.waitingForConfirm {
			switch msg.String() {
			case "y", "Y":
				m.waitingForConfirm = false
				m.state = "writing"
				m.messages = append(m.messages, "Confirmed!")
			case "n", "N":
				m.waitingForConfirm = false
				m.state = "idle"
				m.messages = append(m.messages, "Cancelled")
			}
			return m, cmd
		}

		switch msg.Type {
		case tea.KeyEnter:
			input := m.input.Value()
			if input != "" {
				if input == "exit" || input == "quit" {
					return m, tea.Quit
				}
				m.messages = append(m.messages, input)
				m.processInput(input)
			}
			m.input.Model.SetValue("")
		}
	}

	m.input, _ = m.input.Update(msg)
	return m, cmd
}

func (m *Model) processInput(input string) {
	m.state = "reading"

	if contains(input, "tax") || contains(input, "itr") {
		m.tree.AddRead("form16.pdf")
		m.tree.AddRead("26as.csv")
		m.state = "processing"

		m.waitingForConfirm = true
		m.confirmMsg = "Create itr1.json with computed tax?"
	} else if contains(input, "ingest") {
		m.tree.AddRead("form16.pdf")
		m.state = "idle"
	} else if contains(input, "generate") {
		m.tree.AddWrite("itr1.json")
		m.waitingForConfirm = true
		m.confirmMsg = "Create itr1.json?"
	} else {
		m.tree.AddRead("data.csv")
		m.state = "idle"
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

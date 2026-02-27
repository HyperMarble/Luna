package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width  int
	height int

	input textinput.Model

	conversation []string
}

func NewModel() *Model {
	ti := textinput.New()
	ti.Placeholder = "Ask Luna..."
	ti.Focus()
	ti.Width = 72

	return &Model{
		input: ti,
		conversation: []string{
			"Luna: Ask anything. I will keep responses clear and short.",
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

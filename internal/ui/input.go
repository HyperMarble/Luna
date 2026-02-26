package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	textinput.Model
}

func NewInput() Input {
	ti := textinput.New()
	ti.Placeholder = "Ask Luna..."
	ti.Focus()
	ti.Width = 60

	return Input{ti}
}

func (i Input) Init() tea.Cmd {
	return textinput.Blink
}

func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	newModel, cmd := i.Model.Update(msg)
	i.Model = newModel
	return i, cmd
}

func (i Input) View() string {
	return i.Model.View()
}

func (i Input) Value() string {
	return i.Model.Value()
}

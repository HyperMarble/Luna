package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	input   Input
	spinner Spinner
	tree    Tree
	output  Output

	state             string
	messages          []string
	waitingForConfirm bool
	confirmMsg        string
}

func NewModel() *Model {
	return &Model{
		input:   NewInput(),
		spinner: NewSpinner(),
		tree:    NewTree(),
		output:  NewOutput(),

		state:             "idle",
		messages:          []string{},
		waitingForConfirm: false,
		confirmMsg:        "",
	}
}

func (m *Model) Init() tea.Cmd {
	return m.spinner.Init()
}

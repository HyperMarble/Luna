package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Spinner struct {
	spinner spinner.Model
}

func NewSpinner() Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Spinner{spinner: s}
}

func (m Spinner) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Spinner) Update(msg tea.Msg) (Spinner, tea.Cmd) {
	newSpinner, cmd := m.spinner.Update(msg)
	m.spinner = newSpinner
	return m, cmd
}

func (m Spinner) View() string {
	return m.spinner.View()
}

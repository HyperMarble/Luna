package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Spinner struct {
	frames []string
	index  int
	active bool
}

func NewSpinner() Spinner {
	return Spinner{
		frames: []string{"◐", "◑", "◒", "◓", "◑", "◒"},
		index:  0,
		active: false,
	}
}

func (s *Spinner) Start() {
	s.active = true
}

func (s *Spinner) Stop() {
	s.active = false
}

func (s *Spinner) Tick() {
	if s.active {
		s.index = (s.index + 1) % len(s.frames)
	}
}

func (s Spinner) View() string {
	if !s.active {
		return ""
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return style.Render(s.frames[s.index])
}

func (s *Spinner) DelayedTick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(100 * time.Millisecond)
		return spinnerTickMsg{}
	}
}

type spinnerTickMsg struct{}

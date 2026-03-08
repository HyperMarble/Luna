package tui_test

import (
	"testing"

	"github.com/HyperMarble/Luna/internal/tui"
	tea "charm.land/bubbletea/v2"
)

func TestNewModel_InputFocused(t *testing.T) {
	m := tui.NewModel()
	if !m.Input().Focused() {
		t.Fatal("expected input to be focused on init")
	}
}

func TestNewModel_EmptyMessages(t *testing.T) {
	m := tui.NewModel()
	if len(m.Messages()) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(m.Messages()))
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	m := tui.NewModel()
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	_ = m2
}

package tui_test

import (
	"testing"

	"github.com/hak/luna/internal/tui"
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

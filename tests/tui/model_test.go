package tui_test

import (
	"testing"

	"github.com/HyperMarble/Luna/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
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

func TestUpdate_SubmitAddsUserMessage(t *testing.T) {
	m := tui.NewModel()

	// Simulate window size first (viewport needs dimensions)
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Type "hello" then press enter
	m3, _ := m2.(tui.Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")})
	m4, _ := m3.(tui.Model).Update(tea.KeyMsg{Type: tea.KeyEnter})

	msgs := m4.(tui.Model).Messages()
	if len(msgs) == 0 {
		t.Fatal("expected at least 1 message after submit")
	}
	if msgs[0].Role != "user" {
		t.Fatalf("expected role 'user', got %q", msgs[0].Role)
	}
}

func TestUpdate_StubResponseAddsLunaMessage(t *testing.T) {
	m := tui.NewModel()
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m3, _ := m2.(tui.Model).Update(tui.LunaStubMsg{Text: "hello from luna"})

	msgs := m3.(tui.Model).Messages()
	if len(msgs) == 0 {
		t.Fatal("expected at least 1 message")
	}
	if msgs[0].Role != "luna" {
		t.Fatalf("expected role 'luna', got %q", msgs[0].Role)
	}
}

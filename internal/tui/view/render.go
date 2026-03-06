package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"

	"github.com/HyperMarble/Luna/internal/agent"
	tuilayout "github.com/HyperMarble/Luna/internal/tui/layout"
	"github.com/HyperMarble/Luna/internal/tui/types"
)

// State is the immutable input required for rendering.
type State struct {
	Width    int
	Height   int
	Layout   tuilayout.UI
	Input    textinput.Model
	Messages []types.Message
	Thinking bool
	VerbIdx  int

	// Slash command picker.
	PickerIndex int

	// Model picker tree.
	ModelPickerOpen    bool
	ModelPickerState   int // 0=providers, 1=models, 2=apikey, 3=custommodel
	ModelPickerProvIdx int // Highlighted provider row.
	ModelPickerModIdx  int // Highlighted model row within expanded provider.
	ExpandedProv       int // Index of expanded provider (-1 = none).
	APIKeyInput        textinput.Model
	APIKeyProvider     agent.ProviderInfo
	CustomModelInput   textinput.Model
	ActiveModel        string
}

// Render composes the full TUI layout from top to bottom.
func Render(s State) string {
	layout := s.Layout
	if layout.Width == 0 {
		layout = tuilayout.Compute(s.Width)
	}

	var b strings.Builder
	b.WriteString(renderWelcomeBox(layout.WelcomeWidth))
	b.WriteString("\n")

	if s.ModelPickerOpen {
		b.WriteString(renderModelPicker(s))
	} else {
		b.WriteString(renderMessages(s.Messages))
		b.WriteString(renderThinking(s.Thinking, s.VerbIdx))
	}

	b.WriteString(renderComposer(layout.ComposerWidth, s.Input, s.PickerIndex))
	return padToHeight(b.String(), s.Height)
}

func padToHeight(content string, height int) string {
	if height <= 0 {
		return content
	}
	lines := strings.Count(content, "\n") + 1
	if lines >= height {
		return content
	}
	return content + strings.Repeat("\n", height-lines)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

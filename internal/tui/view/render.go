package view

import (
	"strings"

	"charm.land/bubbles/v2/textinput"

	"github.com/HyperMarble/Luna/internal/agent"
	tuilayout "github.com/HyperMarble/Luna/internal/tui/layout"
	"github.com/HyperMarble/Luna/internal/tui/types"
)

// State is the immutable input required for rendering.
type State struct {
	Width      int
	Height     int
	Layout     tuilayout.UI
	Input      textinput.Model
	Messages   []types.Message
	Thinking   bool
	VerbIdx    int
	BodyView   string
	FooterView string

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

	body := s.BodyView
	if body == "" {
		body = RenderBodyContent(s)
	}
	footer := s.FooterView
	if footer == "" {
		footer = RenderFooter(s)
	}

	if s.Height <= 0 {
		return joinRegions(body, footer)
	}
	return joinRegions(body, footer)
}

func RenderBodyContent(s State) string {
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
		b.WriteString(renderMessages(s.Messages, layout.ComposerWidth))
		b.WriteString(renderThinking(s.Thinking, s.VerbIdx))
	}

	return b.String()
}

func RenderFooter(s State) string {
	layout := s.Layout
	if layout.Width == 0 {
		layout = tuilayout.Compute(s.Width)
	}
	return renderComposer(layout.ComposerWidth, s.Input, s.PickerIndex)
}

func joinRegions(top, bottom string) string {
	switch {
	case top == "":
		return bottom
	case bottom == "":
		return top
	default:
		return top + "\n" + bottom
	}
}

func FitBodyTop(s State, height int) string {
	return fitRegion(RenderBodyContent(s), height, false, 0)
}

func fitRegion(content string, height int, keepTail bool, scrollOffset int) string {
	if height <= 0 {
		return ""
	}

	lines := splitLines(content)
	if len(lines) > height {
		if keepTail {
			maxOffset := max(0, len(lines)-height)
			scrollOffset = max(0, min(maxOffset, scrollOffset))
			end := len(lines) - scrollOffset
			start := max(0, end-height)
			lines = lines[start:end]
		} else {
			lines = lines[:height]
		}
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

// SplitLines splits content on newlines, normalising \r\n and stripping a
// trailing newline. Always returns at least one element.
func SplitLines(content string) []string { return splitLines(content) }

func splitLines(content string) []string {
	if content == "" {
		return []string{""}
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.TrimSuffix(content, "\n")
	if content == "" {
		return []string{""}
	}
	return strings.Split(content, "\n")
}

func renderedLineCount(content string) int {
	if content == "" {
		return 0
	}
	return len(splitLines(content))
}

func RenderedLineCount(content string) int {
	return renderedLineCount(content)
}


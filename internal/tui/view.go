package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// View renders the full terminal UI. Pure function — no side effects.
func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := renderHeader(m.width)
	composer := renderComposer(m.input.View(), m.width)

	return fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), composer)
}

// renderHeader renders the top bar: "◆ Luna   v0.1.0"
func renderHeader(width int) string {
	title := headerStyle.Render("◆ Luna")
	version := versionStyle.Render("v0.1.0")
	gap := strings.Repeat(" ", max(0, width-lipgloss.Width(title)-lipgloss.Width(version)))
	return title + gap + version
}

// renderComposer wraps the text input in a bordered box.
func renderComposer(inputView string, width int) string {
	return composerStyle.Width(width - 2).Render(inputView)
}

// renderMessages builds the full viewport content string from all messages.
// Called from Update — never from View directly.
func renderMessages(messages []Message, thinking bool, spinnerView string, width int) string {
	var sb strings.Builder

	for _, msg := range messages {
		role := roleLabelStyle.Render(roleLabel(msg.Role))
		divider := dividerStyle.Render(strings.Repeat("─", lipgloss.Width(role)))
		sb.WriteString(role + "\n")
		sb.WriteString(divider + "\n")

		if msg.Role == "luna" {
			rendered, err := glamour.Render(msg.Content, "dark")
			if err != nil {
				rendered = msg.Content
			}
			sb.WriteString(rendered)
		} else {
			sb.WriteString("  " + msg.Content + "\n")
		}
		sb.WriteString("\n")
	}

	if thinking {
		role := roleLabelStyle.Render("Luna")
		divider := dividerStyle.Render(strings.Repeat("─", lipgloss.Width(role)))
		sb.WriteString(role + "\n")
		sb.WriteString(divider + "\n")
		sb.WriteString(spinnerLabelStyle.Render("  "+spinnerView+" thinking...") + "\n")
	}

	return sb.String()
}

// roleLabel returns the display name for a role.
func roleLabel(role string) string {
	switch role {
	case "luna":
		return "Luna"
	default:
		return "You"
	}
}

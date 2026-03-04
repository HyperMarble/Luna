package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// mascot is the ASCII block art shown in the welcome box.
var mascot = []string{
	"  ████ ████  ",
	"  ████ ████  ",
	"███       ███",
	"██   ^ ^   ██",
	"██         ██",
	"███       ███",
	"  ████ ████  ",
	"  ████ ████  ",
}

// ── Root renderer ─────────────────────────────────────────────────────────────

// View composes the full TUI layout from top to bottom.
func (m Model) View() string {
	if !m.ready {
		return "\n  Loading..."
	}
	var b strings.Builder
	b.WriteString(renderHeader(m.width) + "\n")
	b.WriteString(m.viewport.View())
	b.WriteString("\n")
	b.WriteString(renderComposer(m.width, m.input, m.pickerIdx))
	return b.String()
}

// viewportContent builds the string set as the viewport's content.
// When there are no messages, the welcome box is shown centered inside the viewport.
// Once conversation starts, the welcome box is replaced by the message history.
func viewportContent(messages []Message, thinking bool, verbIdx int, vpHeight, width int) string {
	if len(messages) == 0 && !thinking {
		return centeredWelcome(vpHeight, width)
	}
	return renderMessages(messages) + renderThinking(thinking, verbIdx)
}

// centeredWelcome vertically centres the welcome box inside the viewport area.
func centeredWelcome(vpHeight, width int) string {
	box := renderWelcomeBox(width)
	boxH := lipgloss.Height(box)
	pad := max(0, (vpHeight-boxH)/2)
	return strings.Repeat("\n", pad) + box
}

// ── Sections ──────────────────────────────────────────────────────────────────

// renderHeader renders the top bar with the logo on the left and version on the right.
func renderHeader(width int) string {
	title := headerStyle.Render("◆ Luna")
	version := versionStyle.Render("v0.1.0")
	gap := strings.Repeat(" ", max(0, width-lipgloss.Width(title)-lipgloss.Width(version)))
	return title + gap + version
}

// renderWelcomeBox renders the bordered welcome section — mascot, title, and
// workspace path. It stays pinned at the top for the entire session.
func renderWelcomeBox(width int) string {
	if width == 0 {
		width = 80
	}
	inner := width - 2 // subtract left+right border characters

	var b strings.Builder
	b.WriteString("\n")
	for _, line := range mascot {
		b.WriteString(centerText(mascotStyle.Render(line), inner) + "\n")
	}
	b.WriteString("\n")
	b.WriteString(centerText(welcomeTitleStyle.Render("Luna"), inner) + "\n")
	b.WriteString(centerText(welcomeSubStyle.Render("AI agent for Chartered Accountants"), inner) + "\n")
	b.WriteString(centerText(welcomePathStyle.Render(workspacePath()), inner) + "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("238")).
		Width(inner).
		Render(b.String())
}

// renderMessages renders the full conversation history in order.
func renderMessages(messages []Message) string {
	var b strings.Builder
	for _, msg := range messages {
		if msg.Role == "user" {
			b.WriteString(renderUserMsg(msg.Content))
		} else {
			b.WriteString(renderLunaMsg(msg.Content))
		}
	}
	return b.String()
}

// renderUserMsg renders a single user message as a grey pill: "> text"
func renderUserMsg(content string) string {
	return userPillStyle.Render("> "+content) + "\n\n"
}

// renderLunaMsg renders a single Luna response with an orange bullet and
// glamour markdown formatting for the body.
func renderLunaMsg(content string) string {
	rendered, err := glamour.Render(content, "dark")
	if err != nil {
		rendered = content + "\n"
	}

	// Split into first line (shown inline with bullet) and the rest.
	lines := strings.SplitN(strings.TrimLeft(rendered, "\n"), "\n", 2)
	bullet := responseBulletStyle.Render("● ")

	var b strings.Builder
	if len(lines) > 0 {
		b.WriteString(bullet + responseTextStyle.Render(strings.TrimSpace(lines[0])) + "\n")
		if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
			b.WriteString(lines[1])
		}
	}
	b.WriteString("\n")
	return b.String()
}

// renderThinking renders the animated "* Verb…" indicator while Luna is processing.
// Returns an empty string when not thinking.
func renderThinking(thinking bool, verbIdx int) string {
	if !thinking {
		return ""
	}
	verb := thinkingVerbs[verbIdx%len(thinkingVerbs)]
	return thinkingStyle.Render("* "+verb+"…") + "\n\n"
}

// renderComposer renders the bottom input area: divider, prompt, text field,
// and the slash command picker (when active).
func renderComposer(width int, input textinput.Model, pickerIdx int) string {
	var b strings.Builder
	b.WriteString(dividerStyle.Render(strings.Repeat("─", width)) + "\n")
	b.WriteString(inputPromptStyle.Render("> ") + input.View())

	if val := input.Value(); strings.HasPrefix(val, "/") {
		b.WriteString("\n" + renderPicker(filteredCommands(val), pickerIdx, width))
	}
	return b.String()
}

// renderPicker renders the slash command dropdown below the input.
// The selected row is highlighted with an arrow; others are indented.
func renderPicker(cmds []slashCommand, idx int, width int) string {
	if len(cmds) == 0 {
		return ""
	}
	col1 := width / 2
	var b strings.Builder
	b.WriteString("\n")
	for i, c := range cmds {
		if i == idx {
			label := pickerSelectedStyle.Render("› " + c.Name)
			gap := strings.Repeat(" ", max(0, col1-lipgloss.Width(label)))
			b.WriteString(label + gap + pickerSelectedDescStyle.Render(c.Desc) + "\n")
		} else {
			label := pickerCmdStyle.Render("  " + c.Name)
			gap := strings.Repeat(" ", max(0, col1-lipgloss.Width(label)))
			b.WriteString(label + gap + pickerDescStyle.Render(c.Desc) + "\n")
		}
	}
	b.WriteString("\n")
	return b.String()
}

// ── Utilities ─────────────────────────────────────────────────────────────────

// centerText pads s with leading spaces so it appears centred within width.
func centerText(s string, width int) string {
	pad := max(0, (width-lipgloss.Width(s))/2)
	return strings.Repeat(" ", pad) + s
}

// workspacePath returns the working directory relative to home ("~/..."),
// or the absolute path when the directory is outside home (e.g. external drives).
func workspacePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "~"
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return cwd
	}
	rel, err := filepath.Rel(home, cwd)
	if err != nil || strings.HasPrefix(rel, "..") {
		return cwd
	}
	return "~/" + rel
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

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

func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString(renderHeader(m.width) + "\n")
	sb.WriteString(renderWelcomeBox(m.width) + "\n")

	for _, msg := range m.messages {
		if msg.Role == "user" {
			sb.WriteString(userPillStyle.Render("> "+msg.Content) + "\n\n")
		} else {
			rendered, err := glamour.Render(msg.Content, "dark")
			if err != nil {
				rendered = msg.Content + "\n"
			}
			lines := strings.SplitN(strings.TrimLeft(rendered, "\n"), "\n", 2)
			bullet := responseBulletStyle.Render("● ")
			if len(lines) > 0 {
				sb.WriteString(bullet + responseTextStyle.Render(strings.TrimSpace(lines[0])) + "\n")
				if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
					sb.WriteString(lines[1])
				}
			}
			sb.WriteString("\n")
		}
	}

	if m.thinking {
		verb := thinkingVerbs[m.verbIdx%len(thinkingVerbs)]
		sb.WriteString(thinkingStyle.Render("* "+verb+"…") + "\n\n")
	}

	sb.WriteString(dividerStyle.Render(strings.Repeat("─", m.width)) + "\n")
	sb.WriteString(inputPromptStyle.Render("> ") + m.input.View())

	inputVal := m.input.Value()
	if strings.HasPrefix(inputVal, "/") {
		sb.WriteString("\n" + renderPicker(filteredCommands(inputVal), m.pickerIdx, m.width))
	}

	return sb.String()
}

func renderWelcomeBox(width int) string {
	if width == 0 {
		width = 80
	}
	inner := width - 2 // subtract border chars

	var content strings.Builder
	content.WriteString("\n")
	for _, line := range mascot {
		content.WriteString(centerText(mascotStyle.Render(line), inner) + "\n")
	}
	content.WriteString("\n")
	content.WriteString(centerText(welcomeTitleStyle.Render("Luna"), inner) + "\n")
	content.WriteString(centerText(welcomeSubStyle.Render("AI agent for Chartered Accountants"), inner) + "\n")
	content.WriteString(centerText(welcomePathStyle.Render(workspacePath()), inner) + "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("238")).
		Width(inner).
		Render(content.String())
}

func renderPicker(cmds []slashCommand, idx int, width int) string {
	if len(cmds) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\n")
	for i, c := range cmds {
		col1Width := width / 2
		if i == idx {
			marker := pickerSelectedStyle.Render("› " + c.Name)
			gap := strings.Repeat(" ", max(0, col1Width-lipgloss.Width(marker)))
			desc := pickerSelectedDescStyle.Render(c.Desc)
			sb.WriteString(marker + gap + desc + "\n")
		} else {
			name := pickerCmdStyle.Render("  " + c.Name)
			gap := strings.Repeat(" ", max(0, col1Width-lipgloss.Width(name)))
			desc := pickerDescStyle.Render(c.Desc)
			sb.WriteString(name + gap + desc + "\n")
		}
	}
	sb.WriteString("\n")
	return sb.String()
}

func renderHeader(width int) string {
	title := headerStyle.Render("◆ Luna")
	version := versionStyle.Render("v0.1.0")
	gap := strings.Repeat(" ", max(0, width-lipgloss.Width(title)-lipgloss.Width(version)))
	return title + gap + version
}

func centerText(s string, width int) string {
	w := lipgloss.Width(s)
	pad := max(0, (width-w)/2)
	return strings.Repeat(" ", pad) + s
}

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

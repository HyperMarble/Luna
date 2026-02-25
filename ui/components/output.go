package components

import (
	"github.com/charmbracelet/lipgloss"
)

type Output struct {
	content string
	style   lipgloss.Style
	border  lipgloss.Style
}

func NewOutput() Output {
	return Output{
		content: "",
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")),
		border: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")),
	}
}

func (o *Output) Set(content string) {
	o.content = content
}

func (o Output) View() string {
	if o.content == "" {
		return ""
	}

	lines := ""
	for i, line := range splitIntoLines(o.content) {
		if i == 0 {
			lines += o.style.Render(line)
		} else {
			lines += "\n  " + o.style.Render(line)
		}
	}

	return o.border.Height(0).Render(lines)
}

func splitIntoLines(s string) []string {
	var lines []string
	current := ""
	for _, ch := range s {
		current += string(ch)
		if ch == '\n' {
			lines = append(lines, current)
			current = ""
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

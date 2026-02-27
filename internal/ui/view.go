package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	width := m.width
	if width < 56 {
		width = 56
	}
	composer := renderComposer(m.input.View(), width)
	footer := renderFooter(lipgloss.Width(composer))
	return renderBanner() + "\n" + composer + "\n" + footer
}

func renderBanner() string {
	art := `██╗░░░░░██╗░░░██╗███╗░░██╗░█████╗░
██║░░░░░██║░░░██║████╗░██║██╔══██╗
██║░░░░░██║░░░██║██╔██╗██║███████║
██║░░░░░██║░░░██║██║╚████║██╔══██║
███████╗╚██████╔╝██║░╚███║██║░░██║
╚══════╝░╚═════╝░╚═╝░░╚══╝╚═╝░░╚═╝`

	label := "AI CA agent"
	pad := max(0, bannerWidth(art)-len([]rune(label)))
	return art + "\n" + strings.Repeat(" ", pad) + secondaryStyle.Render(label)
}

func renderComposer(input string, width int) string {
	return composerStyle.Width(width).Render(input)
}

func renderFooter(width int) string {
	left := "Ctrl+C: quit  |  " + shortDir()
	leftLine := secondaryStyle.Width(width).Align(lipgloss.Left).Render(left)
	return leftLine
}

func shortDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "~/"
	}
	base := filepath.Base(cwd)
	if base == "." || base == string(filepath.Separator) {
		return "~/"
	}
	return "~/" + base
}

func bannerWidth(banner string) int {
	maxWidth := 0
	for _, line := range strings.Split(banner, "\n") {
		lineWidth := len([]rune(line))
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}
	return maxWidth
}

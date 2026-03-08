package view

import (
	"os"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/HyperMarble/Luna/internal/tui/style"
)

// mascot lines — 0-1 top blocks (saffron), 2-5 middle (white), 6-7 bottom blocks (green)
var mascot = []string{
	"  ████ ████  ", // 0 — top saffron
	"  ████ ████  ", // 1 — top saffron
	"███       ███", // 2 — mid white
	"██   ^ ^   ██", // 3 — mid white
	"██         ██", // 4 — mid white
	"███       ███", // 5 — mid white
	"  ████ ████  ", // 6 — bot green
	"  ████ ████  ", // 7 — bot green
}

// RenderWelcomeBox is the exported entry point used by Init().
func RenderWelcomeBox(width int) string { return renderWelcomeBox(width) }

func renderWelcomeBox(width int) string {
	if width == 0 {
		width = 80
	}
	boxWidth := min(width, 72)
	if width >= 40 {
		boxWidth = max(40, boxWidth)
	}
	inner := max(1, boxWidth-2)

	var b strings.Builder
	b.WriteString("\n")
	for i, line := range mascot {
		var s lipgloss.Style
		switch {
		case i <= 1:
			s = style.MascotTop
		case i >= 6:
			s = style.MascotBot
		default:
			s = style.MascotMid
		}
		b.WriteString(leftText(s.Render(line)) + "\n")
	}
	b.WriteString("\n")
	title := style.WelcomeTitle.Render("Luna")
	version := style.WelcomeVersion.Render("v0.0.1")
	b.WriteString(leftText(title+" "+version) + "\n")
	b.WriteString(leftText(style.WelcomeSub.Render("AI agent for Chartered Accountants")) + "\n")
	b.WriteString(leftText(style.WelcomePath.Render(workspacePath())) + "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("238")).
		Width(inner).
		Render(b.String())

	return box
}

func leftText(s string) string { return "  " + s }

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

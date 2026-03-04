package view

import (
	"strings"

	"github.com/charmbracelet/glamour"

	"github.com/HyperMarble/Luna/internal/tui/style"
	"github.com/HyperMarble/Luna/internal/tui/types"
)

func renderMessages(messages []types.Message) string {
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

func renderUserMsg(content string) string {
	return style.UserPill.Render("> "+content) + "\n\n"
}

func renderLunaMsg(content string) string {
	rendered, err := glamour.Render(content, "dark")
	if err != nil {
		rendered = content + "\n"
	}
	lines := strings.SplitN(strings.TrimLeft(rendered, "\n"), "\n", 2)
	bullet := style.ResponseBullet.Render("● ")

	var b strings.Builder
	if len(lines) > 0 {
		b.WriteString(bullet + style.ResponseText.Render(strings.TrimSpace(lines[0])) + "\n")
		if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
			b.WriteString(lines[1])
		}
	}
	b.WriteString("\n")
	return b.String()
}

func renderThinking(thinking bool, verbIdx int) string {
	if !thinking {
		return ""
	}
	verb := types.ThinkingVerbs[verbIdx%len(types.ThinkingVerbs)]
	return style.Thinking.Render("* "+verb+"…") + "\n\n"
}

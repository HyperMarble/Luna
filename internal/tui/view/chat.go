package view

import (
	"strings"

	"github.com/charmbracelet/glamour"

	"github.com/HyperMarble/Luna/internal/tui/style"
	"github.com/HyperMarble/Luna/internal/tui/types"
)

func renderMessages(messages []types.Message, width int) string {
	var b strings.Builder
	for _, msg := range messages {
		if msg.Role == "user" {
			b.WriteString(renderUserMsg(msg.Content))
		} else {
			b.WriteString(renderLunaMsg(msg.Content, width))
		}
	}
	return b.String()
}

// RenderThinking is exported for use in tests / external packages.
func RenderThinking(thinking bool, idx int) string { return renderThinking(thinking, idx) }

func renderUserMsg(content string) string {
	return style.UserPill.Render("> "+content) + "\n\n"
}

// maxMsgWidth caps message width for readability (mirrors crush's maxTextWidth).
const maxMsgWidth = 120

func renderLunaMsg(content string, width int) string {
	// Cap width for readability, same as crush's cappedMessageWidth pattern.
	wrapWidth := width - 4 // account for bullet + padding
	if wrapWidth <= 0 {
		wrapWidth = 76
	}
	if wrapWidth > maxMsgWidth {
		wrapWidth = maxMsgWidth
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(wrapWidth),
	)
	var rendered string
	if err == nil {
		rendered, err = r.Render(content)
	}
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

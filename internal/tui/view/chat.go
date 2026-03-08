package view

import (
	"strings"

	"charm.land/lipgloss/v2"
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
func RenderThinking(thinking bool, idx int) string { return renderThinking(thinking, idx, 0) }

func renderUserMsg(content string) string {
	return style.UserPill.Render("> "+content) + "\n"
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

	// Trim leading and trailing blank lines glamour adds.
	rendered = strings.TrimSpace(rendered)
	lines := strings.SplitN(rendered, "\n", 2)
	bullet := style.ResponseBullet.Render("● ")

	var b strings.Builder
	if len(lines) > 0 {
		b.WriteString(bullet + style.ResponseText.Render(strings.TrimSpace(lines[0])) + "\n")
		if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
			b.WriteString(strings.Trim(lines[1], "\n") + "\n")
		}
	}
	return b.String()
}

// saffron shimmer palette: peak → dim
var saffronShades = []string{
	"#FFD580", // peak (bright gold)
	"#FF9933", // saffron
	"#CC7A29", // mid
	"#7a4510", // dim
}

func renderThinking(thinking bool, verbIdx int, wordIdx int) string {
	if !thinking {
		return ""
	}
	word := types.ThinkingVerbs[wordIdx%len(types.ThinkingVerbs)]
	return "  " + shimmerWord(word+"…", verbIdx) + "\n\n"
}

// shimmerWord sweeps a saffron bright spot left-to-right across a single word.
func shimmerWord(word string, tick int) string {
	runes := []rune(word)
	n := len(runes)
	if n == 0 {
		return ""
	}

	// pos travels 0 → n+4 then wraps; +4 pause at end before next word swap
	pos := tick % (n + 4)

	var b strings.Builder
	for i, ch := range runes {
		dist := pos - i
		if dist < 0 {
			dist = -dist
		}
		var hex string
		if dist < len(saffronShades) {
			hex = saffronShades[dist]
		} else {
			hex = "#4a2d0a" // base dim saffron
		}
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(hex)).Render(string(ch)))
	}
	return b.String()
}

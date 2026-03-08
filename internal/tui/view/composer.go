package view

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"

	"github.com/HyperMarble/Luna/internal/tui/slash"
	"github.com/HyperMarble/Luna/internal/tui/style"
	"github.com/HyperMarble/Luna/internal/tui/types"
)

func thinkingVerb(verbIdx int) string {
	return types.ThinkingVerbs[verbIdx%len(types.ThinkingVerbs)]
}

func renderComposer(width int, input textinput.Model, pickerIdx int) string {
	return renderComposerFull(width, input, pickerIdx, false, 0)
}

func renderComposerThinking(width int, verbIdx int) string {
	return renderComposerFull(width, textinput.Model{}, 0, true, verbIdx)
}

func renderComposerFull(width int, input textinput.Model, pickerIdx int, thinking bool, verbIdx int) string {
	var b strings.Builder
	dividerWidth := max(0, width-2)
	b.WriteString(style.Divider.Render(strings.Repeat("─", dividerWidth)) + "\n")
	if thinking {
		verb := thinkingVerb(verbIdx)
		b.WriteString(style.Thinking.Render("  * " + verb + "…"))
		return b.String()
	}
	b.WriteString(style.InputPrompt.Render("> ") + input.View())
	if val := input.Value(); strings.HasPrefix(val, "/") {
		b.WriteString("\n" + renderPicker(slash.Filtered(val), pickerIdx, width))
	}
	return b.String()
}

func renderPicker(cmds []slash.Command, idx int, width int) string {
	if len(cmds) == 0 {
		return ""
	}
	col1 := width / 2
	var b strings.Builder
	b.WriteString("\n")
	for i, c := range cmds {
		if i == idx {
			label := style.PickerSelected.Render("› " + c.Name)
			gap := strings.Repeat(" ", max(0, col1-lipgloss.Width(label)))
			b.WriteString(label + gap + style.PickerSelectedDesc.Render(c.Desc) + "\n")
		} else {
			label := style.PickerCmd.Render("  " + c.Name)
			gap := strings.Repeat(" ", max(0, col1-lipgloss.Width(label)))
			b.WriteString(label + gap + style.PickerDesc.Render(c.Desc) + "\n")
		}
	}
	b.WriteString("\n")
	return b.String()
}

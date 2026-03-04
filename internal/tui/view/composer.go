package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/HyperMarble/Luna/internal/tui/slash"
	"github.com/HyperMarble/Luna/internal/tui/style"
)

func renderComposer(width int, input textinput.Model, pickerIdx int) string {
	var b strings.Builder
	b.WriteString(style.Divider.Render(strings.Repeat("─", width)) + "\n")
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

package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var s string

	// Banner
	s += banner()
	s += "\n\n"

	// Messages (user inputs)
	for _, msg := range m.messages {
		s += fmt.Sprintf("  %s\n\n", msg)
	}

	// Tree (files being read/written)
	if m.tree.View() != "" {
		s += m.tree.View()
		s += "\n"
	}

	// Spinner (when processing)
	if m.state == "processing" {
		s += m.spinner.View() + " Processing...\n\n"
	}

	// Output
	if m.output.View() != "" {
		s += m.output.View()
		s += "\n"
	}

	// Confirmation
	if m.waitingForConfirm {
		s += lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Render("⚠ " + m.confirmMsg + " [y/n]: ")
	} else {
		// Input
		s += m.input.View()
	}

	return s
}

func banner() string {
	lunaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	moonStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D49A6A"))
	shadowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8B5A2B"))

	luna := `██╗     ██╗   ██╗███╗   ██╗ █████╗ 
			 ██║     ██║   ██║████╗  ██║██╔══██╗
			 ██║     ██║   ██║██╔██╗ ██║███████║
			 ██║     ██║   ██║██║╚██╗██║██╔══██║
			 ███████╗╚██████╔╝██║ ╚████║██║  ██║
			 ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝  ╚═╝`

	moon := `    ████    
  ████████  
 ███` + shadowStyle.Render("████") + `███ 
███` + shadowStyle.Render("██████") + `█
██` + shadowStyle.Render("████████") + `
██` + shadowStyle.Render("████████") + `
███` + shadowStyle.Render("██████") + `█
 ███` + shadowStyle.Render("████") + `███ 
  ████████  
    ████    `

	lunaLines := strings.Split(luna, "\n")
	moonLines := strings.Split(moon, "\n")

	var result string
	for i := 0; i < len(lunaLines) && i < len(moonLines); i++ {
		result += lunaStyle.Render(lunaLines[i]) + "   " + moonStyle.Render(moonLines[i]) + "\n"
	}

	result += "  " + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("AI CA Agent")

	return result
}

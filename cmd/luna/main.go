package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hak/luna/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.NewModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

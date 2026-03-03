package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hak/luna/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "luna: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/HyperMarble/Luna/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// module is the fully qualified path used by go install to fetch the latest binary.
const module = "github.com/HyperMarble/Luna/cmd/luna@latest"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "update" {
		runUpdate()
		return
	}

	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "luna: %v\n", err)
		os.Exit(1)
	}
}

func runUpdate() {
	cmd := exec.Command("go", "install", module)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "luna: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Luna updated successfully!")
}

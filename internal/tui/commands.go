package tui

import "strings"

type slashCommand struct {
	Name string
	Desc string
}

var allCommands = []slashCommand{
	{"/help", "Show all available commands"},
	{"/clear", "Clear conversation and return to welcome"},
	{"/model", "Show current AI model"},
	{"/plugins", "Manage plugins"},
	{"/exit", "Exit Luna"},
}

// filteredCommands returns commands matching the current input prefix.
func filteredCommands(input string) []slashCommand {
	if input == "/" {
		return allCommands
	}
	var out []slashCommand
	for _, c := range allCommands {
		if strings.HasPrefix(c.Name, input) {
			out = append(out, c)
		}
	}
	return out
}

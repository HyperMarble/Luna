package tui

import "strings"

// slashCommand is a single entry in the slash command menu.
type slashCommand struct {
	Name string // e.g. "/help"
	Desc string // shown in the picker alongside the name
}

// allCommands is the full list of available slash commands.
var allCommands = []slashCommand{
	{"/help", "Show all available commands"},
	{"/clear", "Clear the conversation"},
	{"/model", "Show the current AI model"},
	{"/plugins", "Manage plugins"},
	{"/exit", "Exit Luna"},
}

// filteredCommands returns commands whose name starts with the current input.
// If input is exactly "/" all commands are returned.
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

package slash

import "strings"

// Command is a single entry in the slash command menu.
type Command struct {
	Name string
	Desc string
}

// All is the full list of available slash commands.
var All = []Command{
	{"/help", "Show all available commands"},
	{"/clear", "Clear the conversation"},
	{"/model", "Show the current AI model"},
	{"/plugins", "Manage plugins"},
	{"/exit", "Exit Luna"},
}

// Filtered returns commands whose name starts with input.
func Filtered(input string) []Command {
	if input == "/" {
		return All
	}
	var out []Command
	for _, c := range All {
		if strings.HasPrefix(c.Name, input) {
			out = append(out, c)
		}
	}
	return out
}

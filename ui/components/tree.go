package components

import (
	"github.com/charmbracelet/lipgloss"
)

type FileItem struct {
	Name     string
	Type     string // "read" or "write"
	Children []FileItem
}

type Tree struct {
	items []FileItem
	style lipgloss.Style
}

func NewTree() Tree {
	return Tree{
		items: []FileItem{},
		style: lipgloss.NewStyle().Foreground(lipgloss.Color("86")),
	}
}

func (t *Tree) AddRead(filename string) {
	t.items = append(t.items, FileItem{Name: filename, Type: "read"})
}

func (t *Tree) AddWrite(filename string) {
	t.items = append(t.items, FileItem{Name: filename, Type: "write"})
}

func (t *Tree) Clear() {
	t.items = []FileItem{}
}

func (t Tree) View() string {
	if len(t.items) == 0 {
		return ""
	}

	var output string
	for i, item := range t.items {
		prefix := "├"
		if i == len(t.items)-1 {
			prefix = "└"
		}

		icon := "📄"
		if item.Type == "write" {
			icon = "📝"
		}

		output += t.style.Render(prefix+"─ "+icon+" "+item.Name) + "\n"
	}
	return output
}

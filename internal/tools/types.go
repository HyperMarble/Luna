package tools

import "context"

// ToolName is the stable identifier for a tool.
type ToolName string

const (
	ToolRGSearch  ToolName = "rg_search"
	ToolReadFile  ToolName = "read_file"
	ToolGlob      ToolName = "glob"
	ToolGrep      ToolName = "grep"
	ToolWebSearch ToolName = "web_search"
	ToolWebFetch  ToolName = "web_fetch"
	ToolEdit      ToolName = "edit"
	ToolMultiEdit ToolName = "multi_edit"
)

// Request is a generic tool request payload.
type Request struct {
	Input map[string]any
}

// Result is a generic tool response payload.
type Result struct {
	Output map[string]any
}

// Tool is the runtime contract every Luna tool implements.
type Tool interface {
	Name() ToolName
	Run(context.Context, Request) (Result, error)
}

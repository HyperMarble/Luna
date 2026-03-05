package tools

import "context"

type stubTool struct {
	name ToolName
}

// NewStub returns a placeholder implementation for a tool name.
func NewStub(name ToolName) Tool {
	return stubTool{name: name}
}

func (t stubTool) Name() ToolName {
	return t.name
}

func (t stubTool) Run(_ context.Context, _ Request) (Result, error) {
	return Result{}, ErrNotImplemented
}

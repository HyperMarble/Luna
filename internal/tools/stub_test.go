package tools

import "testing"

func TestStubTool_RunNotImplemented(t *testing.T) {
	t.Parallel()

	tool := NewStub(ToolReadFile)
	if _, err := tool.Run(t.Context(), Request{}); err != ErrNotImplemented {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}

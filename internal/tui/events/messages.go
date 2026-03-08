package events

// UserSubmitMsg is sent when the user submits a message programmatically.
type UserSubmitMsg struct{ Text string }

// AgentResponseMsg carries a complete text response (non-streaming fallback).
type AgentResponseMsg struct{ Text string }

// LunaStubMsg is kept as an alias for compatibility with current tests/callers.
type LunaStubMsg = AgentResponseMsg

// AgentChunkMsg carries a single streamed token from the agent.
type AgentChunkMsg struct{ Text string }

// AgentDoneMsg signals that the stream has ended.
type AgentDoneMsg struct{ Err error }

// ModelChangedMsg is sent when the user selects a new provider/model.
type ModelChangedMsg struct {
	Provider string
	ModelID  string
	Label    string
}

// SaveAPIKeyMsg is returned by the tea.Cmd that writes an API key to disk.
// Err is non-nil if the write failed.
type SaveAPIKeyMsg struct {
	EnvKey string
	Value  string
	Err    error
}

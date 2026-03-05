package events

// UserSubmitMsg is sent when the user submits a message programmatically.
type UserSubmitMsg struct{ Text string }

// AgentResponseMsg carries a text response returned by the agent service.
type AgentResponseMsg struct{ Text string }

// LunaStubMsg is kept as an alias for compatibility with current tests/callers.
type LunaStubMsg = AgentResponseMsg

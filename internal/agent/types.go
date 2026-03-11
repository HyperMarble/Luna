package agent

// Request represents a single agent prompt request.
type Request struct {
	Prompt string
}

// Response represents a single agent text response.
type Response struct {
	Text string
}

// EventType identifies a runtime event emitted while handling a request.
type EventType string

const (
	EventToolStart EventType = "tool_start"
	EventToolEnd   EventType = "tool_end"
)

// Event reports non-text activity such as tool execution.
type Event struct {
	Type   EventType
	Name   string
	Detail string
}

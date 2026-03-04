package types

// ThinkingVerbs cycles through action words shown while Luna is processing.
var ThinkingVerbs = []string{
	"Thinking", "Analyzing", "Processing", "Computing",
	"Reasoning", "Calculating", "Reviewing", "Working",
}

// Message is a single turn in the conversation.
type Message struct {
	Role    string // "user" | "luna"
	Content string
}

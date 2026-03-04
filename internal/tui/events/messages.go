package events

// UserSubmitMsg is sent when the user submits a message programmatically.
type UserSubmitMsg struct{ Text string }

// LunaStubMsg carries a stub response from Luna until the real agent is wired.
type LunaStubMsg struct{ Text string }

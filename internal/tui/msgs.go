package tui

// UserSubmitMsg is sent when the user submits a message programmatically (e.g. tests).
type UserSubmitMsg struct{ Text string }

// LunaStubMsg carries a stub response from Luna until the real agent is wired in.
type LunaStubMsg struct{ Text string }

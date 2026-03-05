package agent

import "errors"

var (
	// ErrEmptyPrompt is returned when a request has no prompt text.
	ErrEmptyPrompt = errors.New("empty prompt")
)

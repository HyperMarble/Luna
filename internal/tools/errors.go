package tools

import "errors"

var (
	// ErrToolNotFound is returned when a tool name is not registered.
	ErrToolNotFound = errors.New("tool not found")
	// ErrNotImplemented is returned by stub tools.
	ErrNotImplemented = errors.New("tool not implemented")
)

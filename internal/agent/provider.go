package agent

import "context"

// Provider generates responses for agent requests.
type Provider interface {
	Generate(context.Context, Request) (Response, error)
}

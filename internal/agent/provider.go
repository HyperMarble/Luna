package agent

import "context"

// Provider generates responses for agent requests.
type Provider interface {
	Generate(context.Context, Request) (Response, error)
	// StreamGenerate calls onChunk for each token. Returns when stream ends.
	StreamGenerate(context.Context, Request, func(string)) error
}

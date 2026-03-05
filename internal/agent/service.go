package agent

import (
	"context"
	"strings"
)

// Service is the app-facing agent interface.
type Service interface {
	Run(context.Context, Request) (Response, error)
}

type service struct {
	provider Provider
}

// New creates an agent service. If provider is nil, a stub provider is used.
func New(provider Provider) Service {
	if provider == nil {
		provider = NewStubProvider()
	}
	return &service{provider: provider}
}

// Run validates the request and dispatches to the configured provider.
func (s *service) Run(ctx context.Context, req Request) (Response, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		return Response{}, ErrEmptyPrompt
	}
	return s.provider.Generate(ctx, req)
}

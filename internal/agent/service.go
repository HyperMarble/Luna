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

// New creates an agent service. If provider is nil, it auto-detects from
// environment variables (ANTHROPIC_API_KEY, OPENAI_API_KEY, GEMINI_API_KEY,
// GROQ_API_KEY, or LUNA_PROVIDER=ollama). Falls back to stub if none set.
func New(provider Provider) Service {
	if provider == nil {
		provider = AutoDetectProvider()
	}
	return &service{provider: provider}
}

// NewWithModel creates an agent service for a specific provider and model.
func NewWithModel(providerName, modelID string) Service {
	return &service{provider: ProviderForModel(ProviderName(providerName), modelID)}
}

// Run validates the request and dispatches to the configured provider.
func (s *service) Run(ctx context.Context, req Request) (Response, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		return Response{}, ErrEmptyPrompt
	}
	return s.provider.Generate(ctx, req)
}

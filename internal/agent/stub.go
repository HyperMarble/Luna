package agent

import "context"

const stubText = "I'm Luna. Agent coming soon."

type stubProvider struct{}

// NewStubProvider returns a no-op provider used before LLM wiring.
func NewStubProvider() Provider {
	return stubProvider{}
}

// Generate returns a fixed placeholder response.
func (stubProvider) Generate(_ context.Context, _ Request) (Response, error) {
	return Response{Text: stubText}, nil
}

package agent

import (
	"context"
	"strings"
	"time"
)

const stubText = "I'm Luna. Agent coming soon."

type stubProvider struct{}

// NewStubProvider returns a no-op provider used before LLM wiring.
func NewStubProvider() Provider {
	return stubProvider{}
}

func (stubProvider) Generate(_ context.Context, _ Request) (Response, error) {
	return Response{Text: stubText}, nil
}

func (stubProvider) StreamGenerate(_ context.Context, _ Request, onChunk func(string)) error {
	for _, word := range strings.Fields(stubText) {
		onChunk(word + " ")
		time.Sleep(60 * time.Millisecond)
	}
	return nil
}

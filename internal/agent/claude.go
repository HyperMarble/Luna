package agent

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// claudeProvider calls the Anthropic Messages API.
type claudeProvider struct {
	client anthropic.Client
	model  string
}

// NewClaudeProvider returns a Provider backed by the Anthropic API wrapper.
// Reads ANTHROPIC_API_KEY from the environment.
func NewClaudeProvider(model string) Provider {
	if model == "" {
		model = "claude-sonnet-4-6"
	}
	return &claudeProvider{client: anthropic.NewClient(), model: model}
}

func (p *claudeProvider) messageParams(prompt string) anthropic.MessageNewParams {
	return anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: 4096,
		System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
		Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(prompt))},
	}
}

func (p *claudeProvider) Generate(ctx context.Context, req Request) (Response, error) {
	msg, err := p.client.Messages.New(ctx, p.messageParams(req.Prompt))
	if err != nil {
		return Response{}, fmt.Errorf("claude: %w", err)
	}
	var text string
	for _, block := range msg.Content {
		if tb, ok := block.AsAny().(anthropic.TextBlock); ok {
			text += tb.Text
		}
	}
	return Response{Text: text}, nil
}

func (p *claudeProvider) StreamGenerate(ctx context.Context, req Request, onChunk func(string)) error {
	stream := p.client.Messages.NewStreaming(ctx, p.messageParams(req.Prompt))
	for stream.Next() {
		event := stream.Current()
		if cb, ok := event.AsAny().(anthropic.ContentBlockDeltaEvent); ok {
			if delta, ok := cb.Delta.AsAny().(anthropic.TextDelta); ok {
				onChunk(delta.Text)
			}
		}
	}
	return stream.Err()
}

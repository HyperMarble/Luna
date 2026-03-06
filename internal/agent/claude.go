package agent

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// claudeProvider calls the Anthropic Messages API.
type claudeProvider struct {
	client anthropic.Client
}

// NewClaudeProvider returns a Provider backed by the Anthropic API wrapper.
// Reads ANTHROPIC_API_KEY from the environment.
func NewClaudeProvider() Provider {
	return &claudeProvider{client: anthropic.NewClient()}
}

func (p *claudeProvider) Generate(ctx context.Context, req Request) (Response, error) {
	msg, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: 4096,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt)),
		},
	})
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

const systemPrompt = `You are Luna, an AI assistant for Chartered Accountants in India.
You help with GST filing, income tax returns, TDS compliance, bank reconciliation, and financial analysis.
Be concise and accurate. Use Indian financial terminology. Format numbers in Indian style (lakhs, crores).`

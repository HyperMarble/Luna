package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// openAIProvider calls any OpenAI-compatible chat completions endpoint.
// Works with: OpenAI, Gemini (v1beta/openai), Groq, Ollama, etc.
type openAIProvider struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
	system  string
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func newOpenAIProvider(baseURL, apiKey, model, system string) Provider {
	return &openAIProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		system:  system,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *openAIProvider) Generate(ctx context.Context, req Request) (Response, error) {
	messages := []openAIMessage{}
	if p.system != "" {
		messages = append(messages, openAIMessage{Role: "system", Content: p.system})
	}
	messages = append(messages, openAIMessage{Role: "user", Content: req.Prompt})

	body, err := json.Marshal(openAIRequest{Model: p.model, Messages: messages})
	if err != nil {
		return Response{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return Response{}, fmt.Errorf("openai: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("openai: read body: %w", err)
	}

	var out openAIResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return Response{}, fmt.Errorf("openai: decode: %w", err)
	}
	if out.Error != nil {
		return Response{}, fmt.Errorf("openai: %s", out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return Response{}, fmt.Errorf("openai: empty response")
	}
	return Response{Text: out.Choices[0].Message.Content}, nil
}

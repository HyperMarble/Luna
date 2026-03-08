package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	Stream   bool            `json:"stream,omitempty"`
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

type openAIStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
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

func (p *openAIProvider) messages(prompt string) []openAIMessage {
	msgs := []openAIMessage{}
	if p.system != "" {
		msgs = append(msgs, openAIMessage{Role: "system", Content: p.system})
	}
	return append(msgs, openAIMessage{Role: "user", Content: prompt})
}

func (p *openAIProvider) doRequest(ctx context.Context, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}
	return p.client.Do(req)
}

func (p *openAIProvider) Generate(ctx context.Context, req Request) (Response, error) {
	body, err := json.Marshal(openAIRequest{Model: p.model, Messages: p.messages(req.Prompt)})
	if err != nil {
		return Response{}, err
	}
	resp, err := p.doRequest(ctx, body)
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
		return Response{}, fmt.Errorf("openai: empty response (status %d): %s", resp.StatusCode, string(raw))
	}
	return Response{Text: out.Choices[0].Message.Content}, nil
}

func (p *openAIProvider) StreamGenerate(ctx context.Context, req Request, onChunk func(string)) error {
	body, err := json.Marshal(openAIRequest{Model: p.model, Messages: p.messages(req.Prompt), Stream: true})
	if err != nil {
		return err
	}
	resp, err := p.doRequest(ctx, body)
	if err != nil {
		return fmt.Errorf("openai: %w", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk openAIStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) > 0 {
			if text := chunk.Choices[0].Delta.Content; text != "" {
				onChunk(text)
			}
		}
	}
	return scanner.Err()
}

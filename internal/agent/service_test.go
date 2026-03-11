package agent

import (
	"context"
	"fmt"
	"strings"
	"testing"

	lunatools "github.com/HyperMarble/Luna/internal/tools"
)

func TestService_Run_EmptyPrompt(t *testing.T) {
	t.Parallel()

	svc := New(nil)
	_, err := svc.Run(t.Context(), Request{Prompt: "   "})
	if err != ErrEmptyPrompt {
		t.Fatalf("expected ErrEmptyPrompt, got %v", err)
	}
}

func TestService_Run_UsesProvider(t *testing.T) {
	t.Parallel()

	svc := New(testProvider{})
	resp, err := svc.Run(t.Context(), Request{Prompt: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Text != "ok: hello" {
		t.Fatalf("unexpected response text: %q", resp.Text)
	}
}

func TestStubProvider_Generate(t *testing.T) {
	t.Parallel()

	resp, err := NewStubProvider().Generate(context.Background(), Request{Prompt: "anything"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Text != stubText {
		t.Fatalf("unexpected stub text: %q", resp.Text)
	}
}

type testProvider struct{}

func (testProvider) Generate(_ context.Context, req Request) (Response, error) {
	return Response{Text: "ok: " + req.Prompt}, nil
}

func (testProvider) StreamGenerate(_ context.Context, req Request, onChunk func(string)) error {
	onChunk("ok: " + req.Prompt)
	return nil
}

type recordingProvider struct {
	responses       []string
	prompts         []string
	summaryResponse string
	defaultResponse string
}

func (p *recordingProvider) Generate(_ context.Context, req Request) (Response, error) {
	p.prompts = append(p.prompts, req.Prompt)
	if strings.Contains(req.Prompt, "Summarize the older Luna conversation below") {
		return Response{Text: p.summaryResponse}, nil
	}
	if len(p.responses) > 0 {
		resp := p.responses[0]
		p.responses = p.responses[1:]
		return Response{Text: resp}, nil
	}
	return Response{Text: p.defaultResponse}, nil
}

func (p *recordingProvider) StreamGenerate(_ context.Context, req Request, onChunk func(string)) error {
	p.prompts = append(p.prompts, req.Prompt)
	resp := p.defaultResponse
	if len(p.responses) > 0 {
		resp = p.responses[0]
		p.responses = p.responses[1:]
	}
	onChunk(resp)
	return nil
}

type scriptedProvider struct {
	responses []string
	index     int
}

func (p *scriptedProvider) Generate(_ context.Context, _ Request) (Response, error) {
	resp := p.responses[p.index]
	p.index++
	return Response{Text: resp}, nil
}

func (p *scriptedProvider) StreamGenerate(_ context.Context, req Request, onChunk func(string)) error {
	onChunk(req.Prompt)
	return nil
}

type testTool struct {
	name lunatools.ToolName
	out  lunatools.Result
}

func (t testTool) Name() lunatools.ToolName { return t.name }

func (t testTool) Run(_ context.Context, _ lunatools.Request) (lunatools.Result, error) {
	return t.out, nil
}

func TestService_Run_UsesWebTools(t *testing.T) {
	t.Parallel()

	registry := lunatools.NewRegistry()
	registry.Register(testTool{
		name: lunatools.ToolWebSearch,
		out: lunatools.Result{Output: map[string]any{
			"results": []map[string]any{{"title": "ITR-1", "url": "https://incometax.gov.in"}},
		}},
	})
	registry.Register(testTool{
		name: lunatools.ToolWebFetch,
		out:  lunatools.Result{Output: map[string]any{"content": "Official ITR page"}},
	})

	provider := &scriptedProvider{responses: []string{
		`<tool_call>{"tool":"web_search","input":{"query":"ITR forms"}}</tool_call>`,
		`<final>Use ITR-1 from the official Income Tax portal.</final>`,
	}}

	svc := &service{provider: provider, tools: registry}
	resp, err := svc.Run(t.Context(), Request{Prompt: "Find the official ITR forms"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Text != "Use ITR-1 from the official Income Tax portal." {
		t.Fatalf("unexpected response: %q", resp.Text)
	}
}

func TestService_Run_IncludesConversationContext(t *testing.T) {
	t.Parallel()

	provider := &recordingProvider{
		responses: []string{"First answer", "Second answer"},
	}

	svc := New(provider)
	if _, err := svc.Run(t.Context(), Request{Prompt: "What is section 80C?"}); err != nil {
		t.Fatalf("unexpected error on first run: %v", err)
	}
	if _, err := svc.Run(t.Context(), Request{Prompt: "Explain it again briefly"}); err != nil {
		t.Fatalf("unexpected error on second run: %v", err)
	}

	if len(provider.prompts) != 2 {
		t.Fatalf("expected 2 prompts, got %d", len(provider.prompts))
	}
	secondPrompt := provider.prompts[1]
	if !strings.Contains(secondPrompt, "User: What is section 80C?") {
		t.Fatalf("expected previous user turn in prompt, got %q", secondPrompt)
	}
	if !strings.Contains(secondPrompt, "Assistant: First answer") {
		t.Fatalf("expected previous assistant turn in prompt, got %q", secondPrompt)
	}
	if !strings.Contains(secondPrompt, "Current user request:\nExplain it again briefly") {
		t.Fatalf("expected current user request in prompt, got %q", secondPrompt)
	}
}

func TestService_Reset_ClearsConversationContext(t *testing.T) {
	t.Parallel()

	provider := &recordingProvider{
		responses: []string{"First answer", "Second answer"},
	}

	svc := New(provider)
	if _, err := svc.Run(t.Context(), Request{Prompt: "Remember this"}); err != nil {
		t.Fatalf("unexpected error on first run: %v", err)
	}

	svc.Reset()

	if _, err := svc.Run(t.Context(), Request{Prompt: "New conversation"}); err != nil {
		t.Fatalf("unexpected error on second run: %v", err)
	}

	if len(provider.prompts) != 2 {
		t.Fatalf("expected 2 prompts, got %d", len(provider.prompts))
	}
	if strings.Contains(provider.prompts[1], "Remember this") {
		t.Fatalf("expected reset conversation to drop old context, got %q", provider.prompts[1])
	}
}

func TestService_Run_CompactsOlderConversation(t *testing.T) {
	t.Parallel()

	longAnswer := strings.Repeat("assistant details ", 80)
	provider := &recordingProvider{
		defaultResponse: longAnswer,
		summaryResponse: "Summary of earlier turns",
	}

	svc := New(provider)
	for i := range 7 {
		prompt := fmt.Sprintf("Question %d %s", i+1, strings.Repeat("x", 1200))
		if _, err := svc.Run(t.Context(), Request{Prompt: prompt}); err != nil {
			t.Fatalf("unexpected error on run %d: %v", i+1, err)
		}
	}

	if _, err := svc.Run(t.Context(), Request{Prompt: "What should I remember now?"}); err != nil {
		t.Fatalf("unexpected error on final run: %v", err)
	}

	foundSummaryCall := false
	for _, prompt := range provider.prompts {
		if strings.Contains(prompt, "Summarize the older Luna conversation below") {
			foundSummaryCall = true
			break
		}
	}
	if !foundSummaryCall {
		t.Fatalf("expected summarization call, prompts: %#v", provider.prompts)
	}

	finalPrompt := ""
	for i := len(provider.prompts) - 1; i >= 0; i-- {
		if !strings.Contains(provider.prompts[i], "Summarize the older Luna conversation below") {
			finalPrompt = provider.prompts[i]
			break
		}
	}
	if finalPrompt == "" {
		t.Fatalf("expected a non-summary prompt, prompts: %#v", provider.prompts)
	}
	if !strings.Contains(finalPrompt, "Summary:\nSummary of earlier turns") {
		t.Fatalf("expected compacted summary in final prompt, got %q", finalPrompt)
	}
	if strings.Contains(finalPrompt, "Question 1 ") {
		t.Fatalf("expected oldest turn to be summarized out, got %q", finalPrompt)
	}
}

func TestCompactToolOutputForPrompt_WebFetch(t *testing.T) {
	t.Parallel()

	content := strings.Repeat("A", maxToolOutputChars+250)
	output := compactToolOutputForPrompt(string(lunatools.ToolWebFetch), map[string]any{
		"title":   "  Official page  ",
		"body":    "raw body should not survive",
		"content": content,
	})

	if _, ok := output["body"]; ok {
		t.Fatalf("expected raw body to be removed: %#v", output)
	}
	if got := output["title"]; got != "Official page" {
		t.Fatalf("unexpected title value: %#v", got)
	}
	if got, ok := output["content"].(string); !ok || !strings.Contains(got, "...[truncated]") {
		t.Fatalf("expected truncated content, got %#v", output["content"])
	}
}

package agent

import (
	"context"
	"testing"
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

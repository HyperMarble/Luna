package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type webFetchTool struct {
	client *http.Client
}

// NewWebFetch returns a tool that fetches the raw content of a URL.
func NewWebFetch() Tool {
	return &webFetchTool{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (t *webFetchTool) Name() ToolName { return ToolWebFetch }

func (t *webFetchTool) Run(ctx context.Context, req Request) (Result, error) {
	url, ok := req.Input["url"].(string)
	if !ok || url == "" {
		return Result{}, fmt.Errorf("web_fetch: missing url")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{}, fmt.Errorf("web_fetch: %w", err)
	}
	httpReq.Header.Set("User-Agent", "Luna/0.0.1")

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return Result{}, fmt.Errorf("web_fetch: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1 MB cap
	if err != nil {
		return Result{}, fmt.Errorf("web_fetch: read body: %w", err)
	}

	return Result{Output: map[string]any{
		"status": resp.StatusCode,
		"body":   string(body),
	}}, nil
}

package tools

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestWebFetchToolRunHTMLMarkdown(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(strings.NewReader(`
					<html>
						<head><title>GSTR-3B Guide</title><script>ignored()</script></head>
						<body>
							<nav>Navigation</nav>
							<main>
								<h1>GSTR-3B</h1>
								<p>File monthly return.</p>
							</main>
							<footer>Footer text</footer>
						</body>
					</html>
				`)),
				Request: req,
			}
			resp.Header.Set("Content-Type", "text/html; charset=utf-8")
			return resp, nil
		}),
	}

	tool := &webFetchTool{client: client}
	result, err := tool.Run(context.Background(), Request{Input: map[string]any{
		"url":    "https://example.com/gstr3b",
		"format": "markdown",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := result.Output["content"].(string)
	if !strings.Contains(content, "# GSTR-3B Guide") {
		t.Fatalf("expected title in markdown, got %q", content)
	}
	if !strings.Contains(content, "File monthly return.") {
		t.Fatalf("expected body text in markdown, got %q", content)
	}
	if strings.Contains(content, "Navigation") || strings.Contains(content, "Footer text") {
		t.Fatalf("expected noisy elements to be removed, got %q", content)
	}
}

func TestWebFetchToolRunJSON(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"name":"itr","version":1}`)),
				Request:    req,
			}
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		}),
	}

	tool := &webFetchTool{client: client}
	result, err := tool.Run(context.Background(), Request{Input: map[string]any{
		"url": "https://example.com/itr.json",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := result.Output["content"].(string)
	if !strings.Contains(content, "\"name\": \"itr\"") {
		t.Fatalf("expected pretty json output, got %q", content)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

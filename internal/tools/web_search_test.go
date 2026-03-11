package tools

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestWebSearchToolRun(t *testing.T) {
	t.Parallel()

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if got := req.URL.Query().Get("q"); got != "gst return site:gst.gov.in" {
				t.Fatalf("unexpected query: %q", got)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(strings.NewReader(`
					<html><body>
						<a class="result-link" href="//duckduckgo.com/l/?uddg=https%3A%2F%2Fgst.gov.in%2Fhelp">GSTR Help</a>
						<td class="result-snippet">Official GST return help.</td>
						<a class="result-link" href="https://tutorial.gst.gov.in/userguide/returns/GSTR3B.htm">GSTR-3B Manual</a>
						<td class="result-snippet">Portal manual for GSTR-3B.</td>
					</body></html>
				`)),
				Request: req,
			}, nil
		}),
	}

	tool := &webSearchTool{
		client:  client,
		baseURL: defaultDuckDuckGoLiteURL,
	}

	result, err := tool.Run(context.Background(), Request{Input: map[string]any{
		"query":       "gst return",
		"domains":     []string{"gst.gov.in"},
		"max_results": 2,
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Output["count"] != 2 {
		t.Fatalf("expected 2 results, got %#v", result.Output["count"])
	}

	results, ok := result.Output["results"].([]map[string]any)
	if !ok {
		t.Fatalf("results has unexpected type %T", result.Output["results"])
	}
	if results[0]["url"] != "https://gst.gov.in/help" {
		t.Fatalf("unexpected decoded url: %#v", results[0]["url"])
	}
}

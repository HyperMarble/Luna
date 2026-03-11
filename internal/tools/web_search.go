package tools

import (
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const defaultDuckDuckGoLiteURL = "https://lite.duckduckgo.com/lite/"

type webSearchTool struct {
	client  *http.Client
	baseURL string
}

type searchResult struct {
	Title    string
	URL      string
	Snippet  string
	Position int
}

var searchUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0",
}

// NewWebSearch returns a tool that searches the web via DuckDuckGo lite HTML.
func NewWebSearch() Tool {
	return &webSearchTool{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: defaultDuckDuckGoLiteURL,
	}
}

func (t *webSearchTool) Name() ToolName { return ToolWebSearch }

func (t *webSearchTool) Run(ctx context.Context, req Request) (Result, error) {
	query, ok := req.Input["query"].(string)
	if !ok || strings.TrimSpace(query) == "" {
		return Result{}, fmt.Errorf("web_search: missing query")
	}

	maxResults := parsePositiveInt(req.Input["max_results"], 10)
	if maxResults > 20 {
		maxResults = 20
	}

	searchQuery := strings.TrimSpace(query)
	if domains := parseStringList(req.Input["domains"]); len(domains) > 0 {
		searchQuery += " " + strings.Join(prefixDomains(domains), " ")
	}

	results, err := t.search(ctx, searchQuery, maxResults)
	if err != nil {
		return Result{}, err
	}

	items := make([]map[string]any, 0, len(results))
	for _, result := range results {
		items = append(items, map[string]any{
			"title":    result.Title,
			"url":      result.URL,
			"snippet":  result.Snippet,
			"position": result.Position,
		})
	}

	return Result{Output: map[string]any{
		"query":   query,
		"count":   len(items),
		"results": items,
	}}, nil
}

func (t *webSearchTool) search(ctx context.Context, query string, maxResults int) ([]searchResult, error) {
	searchURL, err := url.Parse(t.baseURL)
	if err != nil {
		return nil, fmt.Errorf("web_search: invalid base url: %w", err)
	}

	values := searchURL.Query()
	values.Set("q", query)
	searchURL.RawQuery = values.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("web_search: %w", err)
	}
	setBrowserHeaders(httpReq)

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("web_search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("web_search: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, fmt.Errorf("web_search: read body: %w", err)
	}

	results, err := parseDuckDuckGoLiteResults(string(body), maxResults)
	if err != nil {
		return nil, fmt.Errorf("web_search: %w", err)
	}
	return results, nil
}

func setBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", searchUserAgents[rand.IntN(len(searchUserAgents))])
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
}

func parseDuckDuckGoLiteResults(htmlContent string, maxResults int) ([]searchResult, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	var results []searchResult
	var current *searchResult

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if len(results) >= maxResults {
			return
		}
		if n.Type == html.ElementNode {
			switch {
			case n.Data == "a" && hasClass(n, "result-link"):
				if current != nil && current.URL != "" {
					current.Position = len(results) + 1
					results = append(results, *current)
					if len(results) >= maxResults {
						return
					}
				}
				current = &searchResult{Title: collapseWhitespace(nodeText(n))}
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						current.URL = cleanDuckDuckGoURL(attr.Val)
						break
					}
				}
			case n.Data == "td" && hasClass(n, "result-snippet") && current != nil:
				current.Snippet = collapseWhitespace(nodeText(n))
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
			if len(results) >= maxResults {
				return
			}
		}
	}

	walk(doc)

	if current != nil && current.URL != "" && len(results) < maxResults {
		current.Position = len(results) + 1
		results = append(results, *current)
	}

	return results, nil
}

func cleanDuckDuckGoURL(raw string) string {
	if strings.HasPrefix(raw, "//duckduckgo.com/l/?uddg=") || strings.HasPrefix(raw, "https://duckduckgo.com/l/?uddg=") {
		if _, encoded, ok := strings.Cut(raw, "uddg="); ok {
			if idx := strings.Index(encoded, "&"); idx >= 0 {
				encoded = encoded[:idx]
			}
			if decoded, err := url.QueryUnescape(encoded); err == nil {
				return decoded
			}
		}
	}
	return raw
}

func prefixDomains(domains []string) []string {
	out := make([]string, 0, len(domains))
	for _, domain := range domains {
		domain = strings.TrimSpace(domain)
		if domain == "" {
			continue
		}
		out = append(out, "site:"+domain)
	}
	return out
}

func parsePositiveInt(v any, fallback int) int {
	switch n := v.(type) {
	case int:
		if n > 0 {
			return n
		}
	case int64:
		if n > 0 {
			return int(n)
		}
	case float64:
		if n > 0 {
			return int(n)
		}
	}
	return fallback
}

func parseStringList(v any) []string {
	switch items := v.(type) {
	case []string:
		return slices.Clone(items)
	case []any:
		out := make([]string, 0, len(items))
		for _, item := range items {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			for _, item := range strings.Fields(attr.Val) {
				if item == class {
					return true
				}
			}
		}
	}
	return false
}

func nodeText(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			b.WriteString(node.Data)
			b.WriteByte(' ')
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(n)
	return b.String()
}

func collapseWhitespace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

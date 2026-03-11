package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/net/html"
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
	rawURL, ok := req.Input["url"].(string)
	if !ok || strings.TrimSpace(rawURL) == "" {
		return Result{}, fmt.Errorf("web_fetch: missing url")
	}
	format := strings.ToLower(strings.TrimSpace(stringOrDefault(req.Input["format"], "markdown")))
	if format == "" {
		format = "markdown"
	}
	if format != "text" && format != "markdown" && format != "html" {
		return Result{}, fmt.Errorf("web_fetch: invalid format %q", format)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return Result{}, fmt.Errorf("web_fetch: invalid url")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return Result{}, fmt.Errorf("web_fetch: unsupported scheme %q", parsedURL.Scheme)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return Result{}, fmt.Errorf("web_fetch: %w", err)
	}
	setBrowserHeaders(httpReq)

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return Result{}, fmt.Errorf("web_fetch: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if err != nil {
		return Result{}, fmt.Errorf("web_fetch: read body: %w", err)
	}
	content := string(body)
	if !utf8.ValidString(content) {
		return Result{}, fmt.Errorf("web_fetch: response is not valid utf-8")
	}

	contentType := resp.Header.Get("Content-Type")
	title := ""
	if looksLikeHTML(contentType, content) {
		var page htmlPage
		page, err = parseHTMLPage(content)
		if err != nil {
			return Result{}, fmt.Errorf("web_fetch: parse html: %w", err)
		}
		title = page.Title
		switch format {
		case "html":
			content = page.HTML
		case "text":
			content = page.Text
		default:
			content = page.Markdown(parsedURL.String())
		}
	} else {
		content = normalizeNonHTMLContent(contentType, content, format)
	}

	return Result{Output: map[string]any{
		"status":       resp.StatusCode,
		"url":          rawURL,
		"final_url":    resp.Request.URL.String(),
		"format":       format,
		"content_type": contentType,
		"title":        title,
		"content":      content,
		"body":         content,
	}}, nil
}

type htmlPage struct {
	Title string
	Text  string
	HTML  string
}

func parseHTMLPage(content string) (htmlPage, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return htmlPage{}, err
	}

	removeNoisyNodes(doc)
	page := htmlPage{
		Title: collapseWhitespace(findFirstText(doc, "title")),
		Text:  extractReadableText(doc),
	}
	page.HTML = renderBodyHTML(doc)
	return page, nil
}

func (p htmlPage) Markdown(sourceURL string) string {
	var b strings.Builder
	if p.Title != "" {
		b.WriteString("# ")
		b.WriteString(p.Title)
		b.WriteString("\n\n")
	}
	b.WriteString("Source: ")
	b.WriteString(sourceURL)
	if p.Text != "" {
		b.WriteString("\n\n")
		b.WriteString(p.Text)
	}
	return strings.TrimSpace(b.String())
}

func removeNoisyNodes(n *html.Node) {
	noisy := map[string]bool{
		"script":   true,
		"style":    true,
		"nav":      true,
		"header":   true,
		"footer":   true,
		"aside":    true,
		"noscript": true,
		"svg":      true,
		"iframe":   true,
	}

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		var toRemove []*html.Node
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && noisy[child.Data] {
				toRemove = append(toRemove, child)
				continue
			}
			walk(child)
		}
		for _, child := range toRemove {
			node.RemoveChild(child)
		}
	}
	walk(n)
}

func findFirstText(n *html.Node, tag string) string {
	var out string
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if out != "" {
			return
		}
		if node.Type == html.ElementNode && node.Data == tag {
			out = nodeText(node)
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(n)
	return out
}

func renderBodyHTML(doc *html.Node) string {
	body := findFirstNode(doc, "body")
	if body == nil {
		return ""
	}
	var b strings.Builder
	for child := body.FirstChild; child != nil; child = child.NextSibling {
		_ = html.Render(&b, child)
	}
	return strings.TrimSpace(b.String())
}

func findFirstNode(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		return n
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if found := findFirstNode(child, tag); found != nil {
			return found
		}
	}
	return nil
}

func extractReadableText(doc *html.Node) string {
	body := findFirstNode(doc, "body")
	if body == nil {
		body = doc
	}

	blockTags := map[string]bool{
		"article": true, "section": true, "div": true, "p": true, "li": true,
		"ul": true, "ol": true, "table": true, "tr": true, "td": true, "th": true,
		"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
		"main": true, "br": true, "pre": true,
	}

	var parts []string
	var current strings.Builder

	flush := func() {
		text := collapseWhitespace(current.String())
		if text != "" {
			parts = append(parts, text)
		}
		current.Reset()
	}

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			current.WriteString(node.Data)
			current.WriteByte(' ')
		}
		if node.Type == html.ElementNode && blockTags[node.Data] {
			flush()
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
		if node.Type == html.ElementNode && blockTags[node.Data] {
			flush()
		}
	}

	walk(body)
	flush()

	return strings.Join(parts, "\n\n")
}

func looksLikeHTML(contentType, body string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/xhtml+xml") || strings.Contains(strings.ToLower(body), "<html")
}

func normalizeNonHTMLContent(contentType, body, format string) string {
	contentType = strings.ToLower(contentType)
	if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/json") {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, []byte(body), "", "  "); err == nil {
			return pretty.String()
		}
	}
	return strings.TrimSpace(body)
}

func stringOrDefault(v any, fallback string) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fallback
}

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	lunatools "github.com/HyperMarble/Luna/internal/tools"
)

// Service is the app-facing agent interface.
type Service interface {
	Run(context.Context, Request) (Response, error)
	// Stream calls onChunk for each token as it arrives.
	Stream(context.Context, Request, func(string), func(Event)) error
	Reset()
}

type service struct {
	provider Provider
	tools    *lunatools.Registry
	memory   *conversationMemory
}

const maxToolIterations = 4

var (
	toolCallRE = regexp.MustCompile(`(?s)<tool_call>\s*(\{.*\})\s*</tool_call>`)
	finalRE    = regexp.MustCompile(`(?s)<final>\s*(.*?)\s*</final>`)
)

type toolDecision struct {
	Tool  string         `json:"tool"`
	Input map[string]any `json:"input"`
}

type toolTranscriptEntry struct {
	Name   string
	Input  map[string]any
	Output map[string]any
}

// New creates an agent service. If provider is nil, it auto-detects from
// environment variables (ANTHROPIC_API_KEY, OPENAI_API_KEY, GEMINI_API_KEY,
// GROQ_API_KEY, or LUNA_PROVIDER=ollama). Falls back to stub if none set.
func New(provider Provider) Service {
	if provider == nil {
		provider = AutoDetectProvider()
	}
	return &service{
		provider: provider,
		tools:    lunatools.NewDefaultRegistry(),
		memory:   newConversationMemory(),
	}
}

// NewWithModel creates an agent service for a specific provider and model.
func NewWithModel(providerName, modelID string) Service {
	return &service{
		provider: ProviderForModel(ProviderName(providerName), modelID),
		tools:    lunatools.NewDefaultRegistry(),
		memory:   newConversationMemory(),
	}
}

func (s *service) ensureMemory() *conversationMemory {
	if s.memory == nil {
		s.memory = newConversationMemory()
	}
	return s.memory
}

// Run validates the request and dispatches to the configured provider.
func (s *service) Run(ctx context.Context, req Request) (Response, error) {
	memory := s.ensureMemory()
	if strings.TrimSpace(req.Prompt) == "" {
		return Response{}, ErrEmptyPrompt
	}
	if shouldUseWebTools(req.Prompt) {
		resp, err := s.runWithTools(ctx, req, nil)
		if err == nil {
			memory.Remember(ctx, s.provider, req.Prompt, resp.Text)
		}
		return resp, err
	}
	resp, err := s.provider.Generate(ctx, Request{
		Prompt: buildConversationPrompt(memory.Snapshot(), req.Prompt),
	})
	if err != nil {
		return Response{}, err
	}
	memory.Remember(ctx, s.provider, req.Prompt, resp.Text)
	return resp, nil
}

// Stream validates the request and streams tokens via onChunk.
func (s *service) Stream(ctx context.Context, req Request, onChunk func(string), onEvent func(Event)) error {
	memory := s.ensureMemory()
	if strings.TrimSpace(req.Prompt) == "" {
		return ErrEmptyPrompt
	}
	if shouldUseWebTools(req.Prompt) {
		resp, err := s.runWithTools(ctx, req, onEvent)
		if err != nil {
			return err
		}
		memory.Remember(ctx, s.provider, req.Prompt, resp.Text)
		streamText(resp.Text, onChunk)
		return nil
	}

	var fullText strings.Builder
	err := s.provider.StreamGenerate(ctx, Request{
		Prompt: buildConversationPrompt(memory.Snapshot(), req.Prompt),
	}, func(chunk string) {
		fullText.WriteString(chunk)
		onChunk(chunk)
	})
	if err != nil {
		return err
	}
	s.memory.Remember(ctx, s.provider, req.Prompt, fullText.String())
	return nil
}

func (s *service) Reset() {
	s.ensureMemory().Reset()
}

func (s *service) runWithTools(ctx context.Context, req Request, onEvent func(Event)) (Response, error) {
	contextBlock := buildConversationContextBlock(s.ensureMemory().Snapshot())
	transcript := make([]toolTranscriptEntry, 0, maxToolIterations)
	for range maxToolIterations {
		decisionResp, err := s.provider.Generate(ctx, Request{
			Prompt: buildToolPlannerPrompt(req.Prompt, contextBlock, transcript),
		})
		if err != nil {
			return Response{}, err
		}

		if final := parseFinal(decisionResp.Text); final != "" {
			return Response{Text: final}, nil
		}

		call, ok := parseToolCall(decisionResp.Text)
		if !ok {
			return Response{Text: strings.TrimSpace(decisionResp.Text)}, nil
		}

		entry, err := s.executeToolCall(ctx, req.Prompt, call, onEvent)
		if err != nil {
			return Response{}, err
		}
		transcript = append(transcript, entry)
	}

	resp, err := s.provider.Generate(ctx, Request{Prompt: buildToolAnswerPrompt(req.Prompt, contextBlock, transcript)})
	if err != nil {
		return Response{}, err
	}
	return Response{Text: strings.TrimSpace(resp.Text)}, nil
}

func (s *service) executeToolCall(ctx context.Context, userPrompt string, call toolDecision, onEvent func(Event)) (toolTranscriptEntry, error) {
	name := strings.TrimSpace(call.Tool)
	if name == "" {
		return toolTranscriptEntry{}, fmt.Errorf("tool call missing tool name")
	}

	input := cloneMap(call.Input)
	switch name {
	case string(lunatools.ToolWebSearch):
		if _, ok := input["max_results"]; !ok {
			input["max_results"] = 5
		}
		if len(anyStrings(input["domains"])) == 0 {
			if domains := suggestedDomains(userPrompt); len(domains) > 0 {
				input["domains"] = domains
			}
		}
	case string(lunatools.ToolWebFetch):
		if _, ok := input["format"]; !ok {
			input["format"] = "markdown"
		}
	default:
		return toolTranscriptEntry{}, fmt.Errorf("unsupported tool %q", name)
	}

	tool, err := s.tools.MustGet(lunatools.ToolName(name))
	if err != nil {
		return toolTranscriptEntry{}, err
	}

	if onEvent != nil {
		onEvent(Event{Type: EventToolStart, Name: name, Detail: toolActivityLabel(name)})
	}
	start := time.Now()
	result, err := tool.Run(ctx, lunatools.Request{Input: input})
	if onEvent != nil {
		onEvent(Event{Type: EventToolEnd, Name: name, Detail: time.Since(start).Round(100 * time.Millisecond).String()})
	}
	if err != nil {
		return toolTranscriptEntry{}, err
	}

	return toolTranscriptEntry{
		Name:   name,
		Input:  input,
		Output: compactToolOutputForPrompt(name, result.Output),
	}, nil
}

func buildToolPlannerPrompt(userPrompt, contextBlock string, transcript []toolTranscriptEntry) string {
	var b strings.Builder
	b.WriteString(toolPlannerPrompt)
	if contextBlock != "" {
		b.WriteString("\n\nConversation context:\n")
		b.WriteString(contextBlock)
	}
	b.WriteString("\n\nCurrent user request:\n")
	b.WriteString(userPrompt)
	if len(transcript) > 0 {
		b.WriteString("\n\nTool transcript:\n")
		for i, entry := range transcript {
			b.WriteString(fmt.Sprintf("\n[%d] %s\n", i+1, entry.Name))
			b.WriteString("input:\n")
			b.WriteString(mustJSON(entry.Input))
			b.WriteString("\noutput:\n")
			b.WriteString(truncateForPrompt(mustJSON(entry.Output), 12000))
			b.WriteString("\n")
		}
	}
	return b.String()
}

func buildToolAnswerPrompt(userPrompt, contextBlock string, transcript []toolTranscriptEntry) string {
	var b strings.Builder
	b.WriteString("Answer the user using the verified tool results below. Be concise, accurate, and mention the official portal or source when relevant.\n\n")
	if contextBlock != "" {
		b.WriteString("Conversation context:\n")
		b.WriteString(contextBlock)
		b.WriteString("\n\n")
	}
	b.WriteString("Current user request:\n")
	b.WriteString(userPrompt)
	b.WriteString("\n\nTool transcript:\n")
	for i, entry := range transcript {
		b.WriteString(fmt.Sprintf("\n[%d] %s\n", i+1, entry.Name))
		b.WriteString("input:\n")
		b.WriteString(mustJSON(entry.Input))
		b.WriteString("\noutput:\n")
		b.WriteString(truncateForPrompt(mustJSON(entry.Output), 12000))
		b.WriteString("\n")
	}
	return b.String()
}

func parseToolCall(text string) (toolDecision, bool) {
	match := toolCallRE.FindStringSubmatch(text)
	if len(match) != 2 {
		return toolDecision{}, false
	}
	var decision toolDecision
	if err := json.Unmarshal([]byte(strings.TrimSpace(match[1])), &decision); err != nil {
		return toolDecision{}, false
	}
	if decision.Input == nil {
		decision.Input = make(map[string]any)
	}
	return decision, true
}

func parseFinal(text string) string {
	match := finalRE.FindStringSubmatch(text)
	if len(match) != 2 {
		return ""
	}
	return strings.TrimSpace(match[1])
}

func streamText(text string, onChunk func(string)) {
	const chunkSize = 32
	for len(text) > 0 {
		if len(text) <= chunkSize {
			onChunk(text)
			return
		}
		onChunk(text[:chunkSize])
		text = text[chunkSize:]
	}
}

func shouldUseWebTools(prompt string) bool {
	lower := strings.ToLower(prompt)
	keywords := []string{
		"search", "web", "latest", "current", "today", "official", "portal",
		"form", "forms", "manual", "instruction", "notification", "circular",
		"itr", "income tax", "gst", "gstr", "tds", "traces", "mca", "roc",
		"llp", "due date", "section", "rule", "finance act",
	}
	for _, keyword := range keywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

func suggestedDomains(prompt string) []string {
	lower := strings.ToLower(prompt)
	switch {
	case strings.Contains(lower, "itr"), strings.Contains(lower, "income tax"), strings.Contains(lower, "26as"), strings.Contains(lower, "ais"), strings.Contains(lower, "tis"):
		return []string{"incometax.gov.in"}
	case strings.Contains(lower, "gst"), strings.Contains(lower, "gstr"), strings.Contains(lower, "itc"), strings.Contains(lower, "eway"), strings.Contains(lower, "e-way"):
		return []string{"gst.gov.in", "cbic-gst.gov.in"}
	case strings.Contains(lower, "tds"), strings.Contains(lower, "traces"), strings.Contains(lower, "24q"), strings.Contains(lower, "26q"), strings.Contains(lower, "27q"), strings.Contains(lower, "27eq"):
		return []string{"tdscpc.gov.in", "protean-tinpan.com"}
	case strings.Contains(lower, "mca"), strings.Contains(lower, "roc"), strings.Contains(lower, "aoc-4"), strings.Contains(lower, "mgt-7"), strings.Contains(lower, "llp"):
		return []string{"mca.gov.in"}
	default:
		return []string{"incometax.gov.in", "gst.gov.in", "mca.gov.in", "cbic.gov.in", "cbdt.gov.in", "tdscpc.gov.in"}
	}
}

func toolActivityLabel(name string) string {
	switch name {
	case string(lunatools.ToolWebSearch):
		return "Searching the web"
	case string(lunatools.ToolWebFetch):
		return "Reading source page"
	default:
		return "Using tool"
	}
}

func mustJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(data)
}

func compactToolOutputForPrompt(name string, output map[string]any) map[string]any {
	compact := cloneMap(output)
	switch name {
	case string(lunatools.ToolWebSearch):
		if results, ok := compact["results"].([]map[string]any); ok {
			compact["results"] = compactSearchResults(results)
		} else if results, ok := compact["results"].([]any); ok {
			compact["results"] = compactSearchResultValues(results)
		}
	case string(lunatools.ToolWebFetch):
		delete(compact, "body")
		if content, ok := compact["content"].(string); ok {
			compact["content"] = truncateForPrompt(strings.TrimSpace(content), maxToolOutputChars)
		}
		if title, ok := compact["title"].(string); ok {
			compact["title"] = strings.TrimSpace(title)
		}
	}
	return compact
}

func compactSearchResults(results []map[string]any) []map[string]any {
	compacted := make([]map[string]any, 0, len(results))
	for _, result := range results {
		row := cloneMap(result)
		if title, ok := row["title"].(string); ok {
			row["title"] = truncateForPrompt(strings.TrimSpace(title), 160)
		}
		if snippet, ok := row["snippet"].(string); ok {
			row["snippet"] = truncateForPrompt(strings.TrimSpace(snippet), 320)
		}
		compacted = append(compacted, row)
	}
	return compacted
}

func compactSearchResultValues(results []any) []map[string]any {
	compacted := make([]map[string]any, 0, len(results))
	for _, item := range results {
		row, ok := item.(map[string]any)
		if !ok {
			continue
		}
		compacted = append(compacted, compactSearchResults([]map[string]any{row})...)
	}
	return compacted
}

func truncateForPrompt(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit] + "\n...[truncated]"
}

func cloneMap(in map[string]any) map[string]any {
	if in == nil {
		return make(map[string]any)
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func anyStrings(v any) []string {
	switch items := v.(type) {
	case []string:
		return items
	case []any:
		out := make([]string, 0, len(items))
		for _, item := range items {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

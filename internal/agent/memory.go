package agent

import (
	"context"
	"strings"
	"sync"
)

const (
	recentTurnsWindow  = 6
	maxHistoryChars    = 12000
	maxSummaryChars    = 3000
	maxToolOutputChars = 6000
)

type conversationTurn struct {
	User      string
	Assistant string
}

type promptContext struct {
	Summary string
	Turns   []conversationTurn
}

type conversationMemory struct {
	mu      sync.Mutex
	summary string
	turns   []conversationTurn
}

func newConversationMemory() *conversationMemory {
	return &conversationMemory{}
}

func (m *conversationMemory) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.summary = ""
	m.turns = nil
}

func (m *conversationMemory) Snapshot() promptContext {
	m.mu.Lock()
	defer m.mu.Unlock()

	turns := append([]conversationTurn(nil), m.turns...)
	return promptContext{
		Summary: m.summary,
		Turns:   turns,
	}
}

func (m *conversationMemory) Remember(ctx context.Context, provider Provider, userPrompt, assistantText string) {
	userPrompt = strings.TrimSpace(userPrompt)
	assistantText = strings.TrimSpace(assistantText)
	if userPrompt == "" || assistantText == "" {
		return
	}

	m.mu.Lock()
	m.turns = append(m.turns, conversationTurn{
		User:      userPrompt,
		Assistant: assistantText,
	})

	if !m.shouldCompactLocked() {
		m.mu.Unlock()
		return
	}

	baseSummary := m.summary
	oldTurns := append([]conversationTurn(nil), m.turns[:len(m.turns)-recentTurnsWindow]...)
	m.mu.Unlock()

	summary := summarizeConversationTurns(ctx, provider, baseSummary, oldTurns)

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.turns) <= recentTurnsWindow || !hasTurnPrefix(m.turns, oldTurns) {
		return
	}

	m.summary = summary
	m.turns = append([]conversationTurn(nil), m.turns[len(oldTurns):]...)
}

func (m *conversationMemory) shouldCompactLocked() bool {
	if len(m.turns) <= recentTurnsWindow {
		return false
	}
	return historyChars(m.summary, m.turns) > maxHistoryChars
}

func summarizeConversationTurns(ctx context.Context, provider Provider, summary string, turns []conversationTurn) string {
	fallback := fallbackConversationSummary(summary, turns)
	if provider == nil || len(turns) == 0 {
		return fallback
	}

	resp, err := provider.Generate(ctx, Request{
		Prompt: buildSummaryPrompt(summary, turns),
	})
	if err != nil {
		return fallback
	}

	out := strings.TrimSpace(resp.Text)
	if out == "" {
		return fallback
	}
	return truncateForPrompt(out, maxSummaryChars)
}

func fallbackConversationSummary(summary string, turns []conversationTurn) string {
	var b strings.Builder
	if strings.TrimSpace(summary) != "" {
		b.WriteString(strings.TrimSpace(summary))
		b.WriteString("\n")
	}
	for _, turn := range turns {
		if user := strings.TrimSpace(turn.User); user != "" {
			b.WriteString("User: ")
			b.WriteString(user)
			b.WriteString("\n")
		}
		if assistant := strings.TrimSpace(turn.Assistant); assistant != "" {
			b.WriteString("Assistant: ")
			b.WriteString(assistant)
			b.WriteString("\n")
		}
	}
	return truncateForPrompt(strings.TrimSpace(b.String()), maxSummaryChars)
}

func buildConversationPrompt(ctx promptContext, userPrompt string) string {
	contextBlock := buildConversationContextBlock(ctx)
	if contextBlock == "" {
		return userPrompt
	}

	var b strings.Builder
	b.WriteString("Conversation context:\n")
	b.WriteString(contextBlock)
	b.WriteString("\n\nCurrent user request:\n")
	b.WriteString(strings.TrimSpace(userPrompt))
	b.WriteString("\n\nAnswer the current user request. Use the conversation context only when it is relevant.")
	return b.String()
}

func buildConversationContextBlock(ctx promptContext) string {
	var b strings.Builder
	if summary := strings.TrimSpace(ctx.Summary); summary != "" {
		b.WriteString("Summary:\n")
		b.WriteString(summary)
		b.WriteString("\n")
	}
	if len(ctx.Turns) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString("Recent conversation:\n")
		for _, turn := range ctx.Turns {
			if user := strings.TrimSpace(turn.User); user != "" {
				b.WriteString("User: ")
				b.WriteString(user)
				b.WriteString("\n")
			}
			if assistant := strings.TrimSpace(turn.Assistant); assistant != "" {
				b.WriteString("Assistant: ")
				b.WriteString(assistant)
				b.WriteString("\n")
			}
		}
	}
	return strings.TrimSpace(b.String())
}

func buildSummaryPrompt(summary string, turns []conversationTurn) string {
	var b strings.Builder
	b.WriteString("Summarize the older Luna conversation below for future turns.\n")
	b.WriteString("Preserve facts, deadlines, cited sections, URLs, user preferences, and unresolved follow-ups.\n")
	b.WriteString("Write one compact plain-English summary only. Do not add bullets or commentary.\n")
	b.WriteString("Keep it under 2200 characters.\n")
	if strings.TrimSpace(summary) != "" {
		b.WriteString("\nExisting summary:\n")
		b.WriteString(strings.TrimSpace(summary))
		b.WriteString("\n")
	}
	b.WriteString("\nOlder conversation:\n")
	for _, turn := range turns {
		if user := strings.TrimSpace(turn.User); user != "" {
			b.WriteString("User: ")
			b.WriteString(user)
			b.WriteString("\n")
		}
		if assistant := strings.TrimSpace(turn.Assistant); assistant != "" {
			b.WriteString("Assistant: ")
			b.WriteString(assistant)
			b.WriteString("\n")
		}
	}
	return strings.TrimSpace(b.String())
}

func historyChars(summary string, turns []conversationTurn) int {
	total := len(summary)
	for _, turn := range turns {
		total += len(turn.User) + len(turn.Assistant)
	}
	return total
}

func hasTurnPrefix(turns, prefix []conversationTurn) bool {
	if len(prefix) > len(turns) {
		return false
	}
	for i := range prefix {
		if turns[i] != prefix[i] {
			return false
		}
	}
	return true
}

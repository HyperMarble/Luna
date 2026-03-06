# Luna — Dev Journal

---

## February 24, 2026

**Project start**

- Initial Luna CLI framework
- Basic command handling scaffold

---

## February 25, 2026

**Initial components**

- Input, output, tree, and spinner components
- Luna skill for OpenCode

---

## February 26, 2026

**Project scaffold**

- Moved UI into `internal/` folder structure
- Added roadmap (`docs/roadmap.md`)
- Added install scripts
- Updated README with project description

---

## February 27, 2026

**Minimal TUI shell**

- Basic Bubble Tea shell with ASCII header and text composer
- Single-input REPL loop, no styling

---

## March 2, 2026

**Full TUI rewrite in Bubble Tea**

- Replaced old component-based UI with clean Bubble Tea Elm architecture
- Added `Model`, `Update`, `View` with proper message routing
- Spinner for thinking state, glamour for markdown rendering in responses
- Lipgloss styles for user pill, response bullet, thinking indicator
- Wired TUI into `cmd/luna/main.go`

---

## March 3, 2026

**Slash command picker + persistent welcome**

- Added `/help`, `/clear`, `/model`, `/plugins`, `/exit` slash commands
- Slash command picker renders inline above composer, filters as you type, tab to complete
- Persistent welcome box visible until first message is sent
- Added `luna update` sub-command (`go install ./cmd/luna@latest`)
- Refactored view and update into small focused functions

---

## March 4, 2026

**TUI modularisation + welcome box**

- Broke monolithic TUI into sub-packages: `layout/`, `style/`, `slash/`, `view/`, `events/`, `types/`
- Extracted `renderWelcomeBox()` — shows centered on startup, disappears on first message
- Added responsive layout via `tuilayout.Compute(width)` — composer and welcome box scale to terminal width
- Enabled `tea.WithAltScreen()` for proper fullscreen rendering
- Added scrollable viewport for message history
- Integration tests added: `tests/tui/model_test.go`

---

## March 5, 2026

**Agent service + tools scaffolding**

- Added `internal/agent/` — `Service` interface, `StubProvider` (returns hardcoded text), `Request`/`Response` types
- Scaffolded `internal/tools/` — `Registry`, `Tool` interface, stub implementations for all planned tools
- Wired agent service into TUI — chat input now goes through `svc.Run()` instead of returning hardcoded string
- Added tools implementation plan to `docs/`
- Switched to source-available no-resale license

---

## March 6, 2026

**Multi-provider LLM support + model picker UI**

- Added `internal/config/config.go` — persists API keys to `~/.luna/config.toml`, injects via `os.Setenv` immediately so no restart needed
- Added `internal/agent/claude.go` — Anthropic native SDK provider (claude-sonnet-4-6 default)
- Added `internal/agent/openai.go` — generic OpenAI-compatible HTTP provider (works for OpenAI, Gemini, Groq, Cerebras, OpenRouter, Ollama)
- Added `internal/agent/providers.go` — full provider registry with free-first ordering (Cerebras, Groq, OpenRouter free → Anthropic, OpenAI, Gemini, Ollama paid)
- Added `internal/tools/web_fetch.go` — real HTTP fetch tool (30s timeout, 1MB cap)
- Added badge styles to `style/styles.go` — `[free]` green, `[API key]` muted, `[unlocked]` saffron, `[local]` muted
- Added `SaveAPIKeyMsg` to `events/messages.go`
- Rewrote `internal/tui/view/modelpicker.go` — provider tree with inline model expansion, API key dialog (asterisk input), custom model dialog for OpenRouter
- Updated `view/render.go` State to carry all picker fields
- Rewrote `internal/tui/model/ui.go` — three-state picker machine (`providers → models → apikey/custommodel`), key saved to disk + immediately unlocks provider
- Updated `cmd/luna/main.go` — calls `config.Load()` on startup

**Models added:**
- Cerebras: Llama 3.3 70B, Llama 4 Scout 17B, Llama 3.1 8B, OpenAI GPT OSS 120B *(default)*, Qwen 3 235B Instruct, Z.ai GLM 4.7
- Groq: GPT OSS 120B, GPT OSS 20B, Qwen 3 32B, Llama 4 Scout, Kimi K2, Llama 3.3 70B, Llama 3.1 8B
- OpenRouter: Llama 3.3 70B (free), Gemini 2.0 Flash (free), Claude Sonnet 4.6, + Custom model
- Anthropic: Claude Opus 4.6, Claude Sonnet 4.6, Claude Haiku 4.5
- OpenAI: GPT-4o, GPT-4o Mini, o3
- Gemini: Gemini 2.5 Pro, Gemini 2.0 Flash
- Ollama: Llama 3.2, Mistral, Code Llama

---

## Build & Run

```bash
go build ./...           # verify compiles
go test ./...            # run all tests
go install ./cmd/luna/   # install binary → ~/go/bin/luna
luna                     # run
```

**Hackathon:** TinyFish $2M Pre-Accelerator — deadline March 29, 2026

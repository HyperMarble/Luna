# Luna Tasks

## Completed
- [x] Basic Bubble Tea TUI (welcome box, composer, slash picker)
- [x] Indian flag mascot colors (saffron top, green bottom)
- [x] Agent service + stub provider
- [x] Tool registry + stub tools
- [x] Claude provider (Anthropic SDK)
- [x] Multi-provider support (OpenAI, Gemini, Groq, Cerebras, OpenRouter, Ollama)
- [x] /model tree picker in TUI
- [x] web_fetch tool (real HTTP)
- [x] Provider-agnostic conversation memory and rolling summary in agent service
- [x] Compact tool transcripts before feeding them back into models
- [x] Reset agent memory on `/clear`
- [x] Add tests for context compaction and reset flow

## In Progress
- [ ] Implement real tools: read_file, glob, grep, rg_search, edit, multi_edit
- [ ] Wire tools into Claude agent loop (tool calling)
- [ ] TinyFish integration for web_fetch (public reference data)

## Backlog
- [ ] Streaming responses (replace blocking svc.Run with Program.Send chunks)
- [ ] Config file (~/.luna/config.toml)
- [ ] Engine: CSV parser (bank statements)
- [ ] Engine: GSTR-2A JSON parser
- [ ] Engine: GST computation (CGST/SGST/IGST)
- [ ] Engine: ITC reconciliation
- [ ] Viewport for scrollable chat history

## Review
- Verified with `GOCACHE=/Volumes/Hak_SSD/Luna/.gocache go test ./...`

# Lessons Learned

## 2026-03-06

### Assumed product scope without asking
- Jumped to portal automation (DSC, OTP, GSP) without understanding what the user meant
- Rule: When user says "fetch data", ask what kind before designing a solution

### Changed product direction repeatedly
- Kept pivoting the architecture mid-session based on tangents
- Rule: Lock in the product definition first, then build. Don't revisit unless user asks.

### Coded before researching
- Implemented web_fetch and providers without checking the actual SDKs/APIs first
- Rule: Research (web search, docs) before writing any integration code

### Didn't follow AGENTS.md / CLAUDE.md
- Skipped gofmt, skipped tasks/todo.md, skipped plan mode for multi-step work
- Rule: Read AGENTS.md and CLAUDE.md at session start, follow without exception

## 2026-03-11

### Added provider-specific optimization when the product needed uniform behavior
- Introduced Anthropic-only prompt caching even though the intended Luna behavior was model-agnostic.
- Rule: For runtime behavior that the user expects to be identical across models, prefer provider-agnostic compaction first and only add provider-specific optimizations when explicitly requested.

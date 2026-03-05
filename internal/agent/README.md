# Agent Package

`internal/agent` contains Luna's agent runtime contracts.

## Current State

- `Service`: app-facing interface used by TUI.
- `Provider`: pluggable backend interface for model providers.
- `stubProvider`: default provider until real LLM integration is wired.

## Next Step

Replace `stubProvider` with a real provider implementation (for example
Claude/OpenAI) and keep the `Service` interface stable.

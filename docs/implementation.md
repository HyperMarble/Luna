# Luna Tools Implementation

This document defines the tools layer for Luna's agent runtime.

## Scope

Luna uses a simple coding-agent-style toolset:

- `rg_search`
- `read_file`
- `glob`
- `grep`
- `web_search`
- `web_fetch`
- `edit`
- `multi_edit` (experimental)
- `sub-agent delegation` (single task stable, multiple tasks experimental)

No extra command abstraction is required at this stage.

## Design Principles

1. Keep tools primitive and composable.
2. Prefer deterministic, structured inputs/outputs.
3. Keep user-facing behavior reviewable before writes.
4. Keep UX consistent with terminal-agent flows.
5. Keep implementation modular and easy to swap.

## Tool Contracts

### `rg_search`

Text search wrapper optimized for speed.

- Input: `query`, optional `path`, optional `include`, optional `literal`.
- Output: matched file paths with snippets.
- Notes:
  - Primary implementation can shell out to `rg` if present.
  - Fallback can use internal Go search for environments without `rg`.

### `read_file`

Read file content with safety limits.

- Input: `path`, optional `offset`, optional `limit`.
- Output: content preview, line count, truncation metadata.

### `glob`

Find files by pattern.

- Input: `pattern`, optional `path`.
- Output: matching paths.

### `grep`

Pattern-based content search.

- Input: `pattern`, optional `path`, optional `include`.
- Output: matches with file + line snippets.

### `web_search`

Search external web sources.

- Input: query string.
- Output: result list (title, url, snippet).

### `web_fetch`

Fetch a specific URL content.

- Input: url, optional format.
- Output: raw or normalized page content.

### `edit`

Single-file targeted content edit.

- Input: file path + operation payload.
- Output: operation summary + preview.

### `multi_edit` (experimental)

Multiple edits in one call.

- Input: ordered edit operations.
- Output: per-edit result summary.
- Status: experimental until rollback/retry semantics are fully stable.

### Sub-agent Delegation

Main agent can delegate tasks to sub-agents:

- Stable: one delegated task at a time.
- Experimental: multiple delegated tasks in parallel.

## Review UX Rules

### Text changes

For text-like outputs, show operational diff preview:

- Action header (example: `Write(path)` / `Edit(path)`).
- Summary (lines changed, sections touched).
- Compact before/after snippet.

### JSON changes

Render JSON changes in table view by default:

- Columns: `Field | Before | After`.
- Nested keys are flattened path-style (example: `invoice.total_tax`).
- Long values are truncated in compact mode.

### Expand behavior

- Default view is compact/truncated.
- `ctrl+o` toggles expanded view for fuller value display.

## Safety + Control

1. Workspace scoping: all local file tools run inside active workspace root.
2. Risky operations require confirmation.
3. All write/edit operations produce review metadata for display.
4. Fail safely with explicit error messages.

## Implementation Order

1. `read_file`
2. `glob`
3. `grep`
4. `rg_search`
5. `edit`
6. `multi_edit` (experimental)
7. `web_search`
8. `web_fetch`
9. sub-agent delegation (single, then multiple experimental)

## Folder Plan

```text
internal/tools/
  registry.go
  types.go
  read_file.go
  glob.go
  grep.go
  rg_search.go
  edit.go
  multi_edit.go
  web_search.go
  web_fetch.go
  delegate.go
```

## Non-Goals (for now)

1. No heavy domain-specific tool explosion.
2. No full GUI-first review layer before TUI behavior is stable.
3. No broad plugin system before core tool reliability.

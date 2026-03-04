# Luna Development Guide

## Build/Test/Lint Commands

- **Build**: `go build ./...` or `go run ./cmd/luna`
- **Test**: `go test ./...`
- **Run single test**: `go test ./tests/tui -run TestName`
- **Format**: `gofmt -w .`

## Code Style Guidelines

- **Imports**: Use standard Go import grouping and formatting.
- **Formatting**: Always format Go code after edits.
- **Naming**: Standard Go conventions.
  - Exported: `PascalCase`
  - Unexported: `camelCase`
- **Error handling**: Return errors explicitly and wrap with `fmt.Errorf` when
  context is needed.
- **Context**: Pass `context.Context` as the first parameter for operations that
  can block or call external systems.
- **Interfaces**: Keep interfaces small and define them in consuming packages.
- **Structs**: Group related fields together; keep models cohesive.
- **Constants**: Group related constants in `const` blocks.
- **JSON tags**: Use `snake_case`.
- **File permissions**: Use octal notation (`0o755`, `0o644`).
- **Log messages**: Start with a capital letter.

## Formatting Rules

- Always format Go code you write.
- Preferred order:
  1. `gofumpt -w .` (if available)
  2. `goimports -w .` (if available)
  3. `gofmt -w .`

## Comments

- Comments on their own lines should start with a capital letter and end with a
  period.

## Testing Guidelines

- Use focused tests when possible before full suite.
- Run full suite (`go test ./...`) before finalizing changes.
- Use table-driven tests for validation-heavy logic.

## TUI Development Rules

- Before changing TUI code, read this guide.
- Keep a single render pipeline.
- Keep layout centralized (see `internal/tui/layout.go`).
- Avoid mixing competing layout/render paths.
- Prefer small, focused files:
  - `view.go` (orchestration)
  - `view_welcome.go`
  - `view_chat.go`
  - `view_composer.go`

## UI Development Instructions

### General Guidelines

- Never use commands to send messages when you can directly mutate children or
  state.
- Keep things simple; do not overcomplicate.
- Create files if needed to separate logic; do not nest models.
- Never do IO or expensive work in `Update`; always use a `tea.Cmd`.
- Never change model state inside a command; send messages and update state in
  the main loop.
- Use `github.com/charmbracelet/x/ansi` for ANSI-aware string manipulation.
  Do not manipulate ANSI strings at byte level.
  - Useful helpers: `ansi.Cut`, `ansi.StringWidth`, `ansi.Strip`,
    `ansi.Truncate`.

### Architecture

#### Main Model (`model/ui.go`)

Keep most logic and state in the main model. This is where:
- Message routing happens.
- Focus and UI state are managed.
- Layout calculations are performed.
- Dialogs are orchestrated.

#### Components Should Be Dumb

Components should not handle Bubble Tea messages directly. Instead:
- Expose methods for state changes.
- Return `tea.Cmd` from methods when side effects are needed.
- Handle their own rendering via `Render(width int) string`.

#### Chat Logic (`model/chat.go`)

Most chat-related logic belongs here. Individual chat items in `chat/` should
be simple renderers that cache their output and invalidate when data changes
(see `cachedMessageItem` in `chat/messages.go`).

### Key Patterns

#### Composition Over Inheritance

Use struct embedding for shared behaviors. See `chat/messages.go` for examples
of reusable embedded structs for highlighting, caching, and focus.

#### Interfaces

- List item interfaces are in `list/item.go`.
- Chat message interfaces are in `chat/messages.go`.
- Dialog interface is in `dialog/dialog.go`.

#### Styling

- All styles are defined in `styles/styles.go`.
- Access styles via `*common.Common` passed to components.
- Use semantic color fields rather than hardcoded colors.

#### Dialogs

- Implement the dialog interface in `dialog/dialog.go`.
- Return message types from `Update()` to signal actions to the main model.
- Use the overlay system for managing dialog lifecycle.

### File Organization

- `model/` - Main UI model and major components (chat, sidebar, etc.).
- `chat/` - Chat message item types and renderers.
- `dialog/` - Dialog implementations.
- `list/` - Generic list component with lazy rendering.
- `common/` - Shared utilities and the Common struct.
- `styles/` - All style definitions.
- `anim/` - Animation system.
- `logo/` - Logo rendering.

### Common Gotchas

- Always account for padding/borders in width calculations.
- Use `tea.Batch()` when returning multiple commands.
- Pass `*common.Common` to components that need styles or app access.

## Committing

- Use semantic commit prefixes:
  - `fix:`, `feat:`, `refactor:`, `chore:`, `docs:`, `test:`
- Keep commit messages concise and single-line unless extra context is required.

# Luna — Current Status

## What Works
- Binary builds and installs: `go install ./cmd/luna/`
- Binary lives at `~/go/bin/luna` — PATH is set in `~/.zshrc`
- TUI launches with header + welcome box + composer

## TUI Architecture (Current)

```
main.go
  ↓ fmt.Print(tui.Welcome(width))   ← printed BEFORE Bubble Tea starts (plain stdout)
  ↓ tea.NewProgram(tui.NewModel())  ← Bubble Tea only owns the composer (2 lines)

View() renders:
  [* Thinking…]        ← only when waiting for response
  ────────────────     ← divider
  > [input field]      ← text input

Messages are printed above with tea.Println() as they arrive.
Terminal scroll buffer is the history — no viewport needed.
```

**Key files:**
- `cmd/luna/main.go` — entry point, prints welcome before TUI
- `internal/tui/model.go` — Model struct, NewModel, Init
- `internal/tui/view.go` — View(), render functions, Welcome()
- `internal/tui/update.go` — Update(), keyboard, slash commands
- `internal/tui/styles.go` — all lipgloss styles
- `internal/tui/commands.go` — slash command list + filtering
- `internal/tui/msgs.go` — UserSubmitMsg, LunaStubMsg

## What's Broken / Not Done

### TUI rendering still has issues
- User reports rendering is still broken after all changes
- Root cause not confirmed yet — needs testing with new build
- Approach tried: viewport (broken), tea.Println from WindowSizeMsg (broken), fmt.Print before program (latest attempt)
- The official Bubble Tea v1 chat example uses viewport + fills terminal height exactly — may need to revisit that

### No real AI yet
- `stubResponseCmd` always returns `"I'm Luna. Agent coming soon."`
- `internal/agent/` doesn't exist yet
- No Claude API connection

### Missing entirely
- `internal/agent/` — Claude API streaming loop
- `internal/tools/` — ingest, validate, compute_gst, etc.
- `internal/engine/` — parser, graph, rules, compute, validation
- `internal/config/` — config loading

## Dependency on PATH

`~/.zshrc` has `export PATH="$HOME/go/bin:$PATH"` added.
Must open a **new terminal tab** after changes for `luna` to resolve.

## Build Commands

```bash
go build ./...           # verify it compiles
go install ./cmd/luna/   # install binary globally
luna                     # run (new terminal tab)
```

## Tech Stack
- Go 1.25.7
- `github.com/charmbracelet/bubbletea v1.3.10`
- `github.com/charmbracelet/bubbles v1.0.0`
- `github.com/charmbracelet/lipgloss v1.1.1`
- `github.com/charmbracelet/glamour v0.10.0`
- Module: `github.com/HyperMarble/Luna`

## Next Steps (in order)

1. **Fix TUI rendering** — confirm the fmt.Print approach works visually
2. **Wire Claude API** — `internal/agent/agent.go` with streaming
3. **Connect agent to TUI** — replace stubResponseCmd with real API call
4. **Build engine** — graph, parser, GST rules
5. **Build tools** — ingest, validate, compute_gst, reconcile_bank

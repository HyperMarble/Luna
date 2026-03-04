# TUI Folder Map

## Root (`internal/tui/`)

- `api.go`: Public Bubble Tea API (`tui.NewModel`, `tui.Model`).

## Folders

- `events/`: Internal message/event types.
- `layout/`: Layout computation.
- `model/`: Main model (`model/ui.go`) owns state and message routing.
- `slash/`: Slash command definitions/filtering.
- `style/`: Shared style definitions.
- `types/`: Shared domain types for TUI modules.
- `view/`: Dumb render helpers (welcome/chat/composer).

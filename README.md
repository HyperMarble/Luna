<p align="center">
  <h1 align="center">Luna</h1>
  <p align="center">Terminal-first AI Agent for Chartered Accountants.</p>
</p>

<p align="center">
  <a href="https://github.com/HyperMarble/Luna"><img src="https://img.shields.io/badge/version-v0.1.0-4c1" alt="Version"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/go-1.25.7-00ADD8?logo=go" alt="Go Version"></a>
  <a href="https://github.com/HyperMarble/Luna/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Source--Available%20No--Resale-blue.svg" alt="License"></a>
</p>

## Version

- App: `v0.1.0` (current CLI/TUI baseline)
- Module: `github.com/HyperMarble/Luna`

## Overview

Luna is a local-first TUI Agent for CA operations agent such as document review, reconciliation preparation, validation-oriented workflows, and filing-ready structured outputs.

Roadmap and product context live in `docs/`.

## Current Capabilities

- Responsive terminal UI built with Bubble Tea.
- Chat-style interaction with slash commands.
- Agent package scaffold (`internal/agent`) with stub provider path.
- Tools package scaffold (`internal/tools`) with registry + typed tool names.
- Test coverage for TUI and new package contracts.

## Package Map

```text
.
├── AGENTS.md                 # Luna Development Guide
├── cmd/
│   └── luna/                 # CLI entrypoint
├── docs/                     # PRD, overview, status, plans
├── internal/
│   ├── agent/                # agent service/provider contracts
│   ├── tools/                # tool registry + tool contracts
│   └── tui/                  # terminal UI architecture
└── tests/
    └── tui/                  # UI behavior tests
```

## Tech Stack

- Go `1.25.7`
- `github.com/charmbracelet/bubbletea v1.3.10`
- `github.com/charmbracelet/bubbles v1.0.0`
- `github.com/charmbracelet/lipgloss v1.1.1-0.20250404203927-76690c660834`
- `github.com/charmbracelet/glamour v0.10.0`

## Requirements

- Go `1.25+`
- POSIX shell or PowerShell

## Installation

### Install from source

```bash
git clone https://github.com/HyperMarble/Luna.git
cd Luna
go install ./cmd/luna
```

### Update local install

```bash
git pull
go install ./cmd/luna
```

## Usage

```bash
luna
```

## Development

### Build

```bash
go build ./...
```

### Test

```bash
go test ./...
```

### Format

```bash
gofmt -w .
```

## Docs

- `docs/OVERVIEW.md` - product overview.
- `docs/PRD.md` - product requirements.
- `docs/STATUS.md` - current implementation status.
- `docs/plans/` - implementation plans.
- `docs/implementation.md` - tools implementation plan.

## Contributing

Read `AGENTS.md` before making changes.

## License

Luna Source-Available License v1.0 (No Resale). See `LICENSE`.

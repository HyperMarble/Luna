# Luna

AI agent for Chartered Accountants.

## Overview

Luna is a terminal-first AI agent for Chartered Accountants, similar in spirit to
what coding agents are for developers. It helps with computation,
reconciliation, validation, and filing preparation.

Project docs and roadmap live in `docs/`.

## Features

- Natural language interface for CA workflows.
- Terminal UI optimized for fast operator workflows.
- Tax/GST/TDS support (incremental build-out).
- Validation and reconciliation pipeline (incremental build-out).
- JSON-oriented outputs for filing workflows.

## Project Structure

```text
.
├── AGENTS.md          # Luna Development Guide
├── cmd/
│   └── luna/          # executable entrypoint
├── docs/              # overview, PRD, status, plans
├── internal/
│   └── tui/           # terminal UI architecture
└── tests/
    └── tui/           # UI/model tests
```

Reference implementation (not Luna runtime code) lives in `crush/`.

## Installation

### Quick Install

```bash
git clone https://github.com/HyperMarble/Luna.git
cd Luna
go install ./cmd/luna
```

### Update

```bash
cd Luna
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

## Contributing

Contributions are welcome. Read `AGENTS.md` first for coding and UI rules.

## Requirements

- Go 1.25+

## License

MIT License.

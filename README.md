# The Project is highly in development contributions are welcome everything the roadmap, and all is inside the docs folder,

# Luna

AI Agent for Chartered Accountants

## Overview

Luna is an AI-powered terminal agent for Chartered Accountants, similar to how Claude Code works for developers. It helps CAs with tax computation, reconciliation, and filing through natural language commands.

## Features

- Natural language interface for CA workflows
- Document parsing (Form 16, 26AS, bank statements)
- Tax, TDS, and GST computation
- 26AS reconciliation
- ITR JSON generation
- Interactive terminal UI

## Installation

### Quick Install

```bash
git clone https://github.com/HyperMarble/Luna.git
cd Luna
go build -o luna
sudo cp luna /usr/local/bin/luna
```

### Update

```bash
cd Luna
git pull
go build -o luna
sudo cp luna /usr/local/bin/luna
```

## Usage

```bash
luna
```

### Commands

```
ingest <file>    Parse and ingest document
compute tax       Compute tax liability
compute tds       Compute TDS obligations
compute gst       Compute GST liability
reconcile 26as   Match 26AS with books
generate itr      Generate ITR JSON
status            Show client status
help              Show available commands
exit              Exit Luna
```

## Example

```
luna> compute tax for rajesh

Reading...
  form16.pdf
  26as.csv

Processing...
  Computing tax...

Writing...
  itr1.json

Confirm? [y/n]: y
```

## Requirements

- Go 1.21 or later

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT License

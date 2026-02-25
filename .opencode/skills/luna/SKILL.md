---
name: Luna CA Agent
description: AI Agent for Chartered Accountants - helps with tax computation, reconciliation, and filing
---

# Luna - AI CA Agent

Luna is an AI-powered CLI agent for Chartered Accountants, similar to how Claude Code works for developers.

## Project Structure

```
Luna/
├── main.go           # Entry point, CLI handling
├── go.mod            # Go module
└── (more files coming)
```

## How to Work with Luna

### Current CLI Commands
- `init <client>` - Create client workspace
- `ingest <file>` - Parse document
- `compute <tax|tds|gst>` - Compute tax/TDS/GST
- `reconcile <26as|gstr>` - Reconcile statements
- `generate <itr|gstr>` - Generate JSON
- `status` - Show status
- `help` - Show help

### How to Add New Features

1. **Add command handler in main.go**
   - Add new case in `handleCommand` function
   - Create new `cmd<Name>` function

2. **Build and test**
   ```bash
   go build -o luna
   ./luna
   ```

3. **Current status** - Basic CLI framework only, no real functionality yet

### Next Features to Build

1. Workspace creation (init command)
2. Document parsing (ingest command) - PDF, CSV, Excel
3. Tax computation engine
4. 26AS reconciliation
5. ITR JSON generation

### Important Notes

- This is a Go project
- Target: Generate JSON/XML for government portals (ITR, GST, TDS)
- End user: Chartered Accountants in India
- Output should be reviewed by CA before uploading to government portals
